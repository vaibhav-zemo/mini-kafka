package storage

import (
	"encoding/json"
	"os"
	"strconv"
	"sync"
)

type OffsetEntry struct {
	Group     string `json:"group"`
	Topic     string `json:"topic"`
	Partition int    `json:"partition"`
	Offset    int64  `json:"offset"`
}

type OffsetManager struct {
	file *os.File
	mu   sync.Mutex

	// in-memory cache for fast lookup
	offsets map[string]int64
}

func NewOffsetManager(path string) (*OffsetManager, error) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}

	om := &OffsetManager{
		file:    file,
		offsets: make(map[string]int64),
	}

	om.load()

	return om, nil
}

func (om *OffsetManager) load() {
	om.file.Seek(0, 0)

	decoder := json.NewDecoder(om.file)

	for {
		var entry OffsetEntry
		if err := decoder.Decode(&entry); err != nil {
			break
		}

		key := om.key(entry.Group, entry.Topic, entry.Partition)
		om.offsets[key] = entry.Offset
	}
}

func (om *OffsetManager) Save(group, topic string, partition int, offset int64) error {
	om.mu.Lock()
	defer om.mu.Unlock()

	entry := OffsetEntry{
		Group:     group,
		Topic:     topic,
		Partition: partition,
		Offset:    offset,
	}

	data, _ := json.Marshal(entry)

	_, err := om.file.Write(append(data, '\n'))
	if err != nil {
		return err
	}

	key := om.key(group, topic, partition)
	om.offsets[key] = offset

	return nil
}

func (om *OffsetManager) Get(group, topic string, partition int) int64 {
	key := om.key(group, topic, partition)

	if val, ok := om.offsets[key]; ok {
		return val
	}

	return 0
}

func (om *OffsetManager) key(group, topic string, partition int) string {
	return group + ":" + topic + ":" + strconv.Itoa(partition)
}
