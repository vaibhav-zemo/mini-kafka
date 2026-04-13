package main

import (
	"encoding/json"
	"log"
	"mini-kafka/broker"
	"net/http"
	"strconv"
)

var b = broker.NewBroker()

type ProduceRequest struct {
	Topic   string `json:"topic"`
	Message string `json:"message"`
}

func produceHandler(w http.ResponseWriter, r *http.Request) {
	var req ProduceRequest
	json.NewDecoder(r.Body).Decode(&req)

	topic, err := b.GetOrCreateTopic(req.Topic)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	offset, err := topic.Partition.Produce([]byte(req.Message))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write([]byte("Offset: " + strconv.FormatInt(offset, 10)))
}

func consumeHandler(w http.ResponseWriter, r *http.Request) {
	topicName := r.URL.Query().Get("topic")
	offsetStr := r.URL.Query().Get("offset")

	offset, _ := strconv.ParseInt(offsetStr, 10, 64)

	topic, err := b.GetOrCreateTopic(topicName)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	msg, nextOffset, err := topic.Partition.Consume(offset)
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

func main() {
	http.HandleFunc("/produce", produceHandler)
	http.HandleFunc("/consume", consumeHandler)

	log.Println("Server running on :8080")
	http.ListenAndServe(":8080", nil)
}
