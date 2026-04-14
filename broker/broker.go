package broker

import "sync"

type Broker struct {
	topics map[string]*Topic
	mu     sync.Mutex
}

func NewBroker() *Broker {
	return &Broker{
		topics: make(map[string]*Topic),
	}
}

func (b *Broker) GetOrCreateTopic(name string) (*Topic, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if topic, exists := b.topics[name]; exists {
		return topic, nil
	}

	topic, err := NewTopic(name, 3)
	if err != nil {
		return nil, err
	}

	b.topics[name] = topic
	return topic, nil
}
