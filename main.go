package main

import (
	"encoding/binary"
	"encoding/json"
	"io"
	"log"
	"mini-kafka/broker"
	"net"
	"net/http"
	"strconv"
)

var b = broker.NewBroker()

type ProduceRequest struct {
	Topic   string `json:"topic"`
	Key     string `json:"key"` // NEW
	Message string `json:"message"`
}

type JoinGroupRequest struct {
	Topic      string `json:"topic"`
	Group      string `json:"group"`
	ConsumerID string `json:"consumer_id"`
}

func produceHandler(w http.ResponseWriter, r *http.Request) {
	var req ProduceRequest
	json.NewDecoder(r.Body).Decode(&req)

	topic, err := b.GetOrCreateTopic(req.Topic)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	partition, offset, err := topic.Produce(req.Key, []byte(req.Message))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"partition": partition,
		"offset":    offset,
	})
}

func consumeHandler(w http.ResponseWriter, r *http.Request) {
	topicName := r.URL.Query().Get("topic")
	offsetStr := r.URL.Query().Get("offset")
	partitionStr := r.URL.Query().Get("partition")

	partition, _ := strconv.Atoi(partitionStr)
	offset, _ := strconv.ParseInt(offsetStr, 10, 64)

	topic, err := b.GetOrCreateTopic(topicName)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	msg, nextOffset, err := topic.Consume(partition, offset)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	resp := map[string]interface{}{
		"message":     string(msg),
		"next_offset": nextOffset,
	}

	json.NewEncoder(w).Encode(resp)
}

func joinGroupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var req JoinGroupRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Topic == "" || req.Group == "" || req.ConsumerID == "" {
		http.Error(w, "Missing fields", http.StatusBadRequest)
		return
	}

	topic, err := b.GetOrCreateTopic(req.Topic)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 👇 Core logic
	topic.JoinGroup(req.Group, req.ConsumerID)

	w.Write([]byte("Consumer joined group successfully"))
}

func consumeGroupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
		return
	}

	topicName := r.URL.Query().Get("topic")
	groupID := r.URL.Query().Get("group")
	consumerID := r.URL.Query().Get("consumer_id")

	if topicName == "" || groupID == "" || consumerID == "" {
		http.Error(w, "Missing query params", http.StatusBadRequest)
		return
	}

	topic, err := b.GetOrCreateTopic(topicName)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	msg, partition, nextOffset, err := topic.ConsumeFromGroup(groupID, consumerID)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":     msg,
		"partition":   partition,
		"next_offset": nextOffset,
	})
}

func webServer() {
	http.HandleFunc("/produce", produceHandler)
	http.HandleFunc("/consume", consumeHandler)
	http.HandleFunc("/join-group", joinGroupHandler)
	http.HandleFunc("/consume-group", consumeGroupHandler)

	log.Println("Server running on :8080")
	http.ListenAndServe(":8080", nil)
}

func main() {
	go webServer()

	ln, err := net.Listen("tcp", ":9092")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("TCP server running on :9092")

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	for {
		req, err := readRequest(conn)
		if err != nil {
			return
		}

		resp := processRequest(req)

		writeResponse(conn, resp)
	}
}

func readRequest(conn net.Conn) (map[string]interface{}, error) {
	lenBuf := make([]byte, 4)

	_, err := conn.Read(lenBuf)
	if err != nil {
		return nil, err
	}

	msgLen := binary.BigEndian.Uint32(lenBuf)

	data := make([]byte, msgLen)
	_, err = io.ReadFull(conn, data)
	if err != nil {
		return nil, err
	}

	var req map[string]interface{}
	json.Unmarshal(data, &req)

	return req, nil
}

func writeResponse(conn net.Conn, resp interface{}) {
	data, _ := json.Marshal(resp)

	lenBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(lenBuf, uint32(len(data)))

	conn.Write(lenBuf)
	conn.Write(data)
}

func processRequest(req map[string]interface{}) map[string]interface{} {
	switch req["type"] {

	case "produce":
		topicName := req["topic"].(string)
		key := req["key"].(string)
		message := req["message"].(string)

		topic, _ := b.GetOrCreateTopic(topicName)

		partition, offset, _ := topic.Produce(key, []byte(message))

		return map[string]interface{}{
			"status":    "ok",
			"partition": partition,
			"offset":    offset,
		}

	case "consume":
		topicName := req["topic"].(string)
		group := req["group"].(string)
		consumer := req["consumer_id"].(string)

		topic, _ := b.GetOrCreateTopic(topicName)

		msg, partition, offset, err := topic.ConsumeFromGroup(group, consumer)
		if err != nil {
			return map[string]interface{}{"error": err.Error()}
		}

		return map[string]interface{}{
			"message":   msg,
			"partition": partition,
			"offset":    offset,
		}
	}

	return map[string]interface{}{"error": "unknown type"}
}
