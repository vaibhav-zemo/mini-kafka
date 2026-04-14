package broker

import (
	"fmt"
	"mini-kafka/utils"
)

type Topic struct {
	name       string
	Partitions []*Partition
	groups     map[string]*ConsumerGroup
}

func NewTopic(name string, numPartitions int) (*Topic, error) {
	partitions := make([]*Partition, numPartitions)
	groups := make(map[string]*ConsumerGroup)

	for i := 0; i < numPartitions; i++ {
		path := fmt.Sprintf("data/%s-%d.log", name, i)

		p, err := NewPartition(path)
		if err != nil {
			return nil, err
		}

		partitions[i] = p
	}

	return &Topic{
		name:       name,
		Partitions: partitions,
		groups:     groups,
	}, nil
}

func (t *Topic) JoinGroup(groupID, consumerID string) {
	group, exists := t.groups[groupID]

	if !exists {
		group = &ConsumerGroup{
			ID:          groupID,
			assignments: make(map[int]*Consumer),
			offsets:     make(map[int]int64),
		}
		t.groups[groupID] = group
	}

	group.consumers = append(group.consumers, &Consumer{ID: consumerID})

	t.rebalance(group)
}

func (t *Topic) rebalance(group *ConsumerGroup) {
	nConsumers := len(group.consumers)
	nPartitions := len(t.Partitions)

	group.assignments = make(map[int]*Consumer)

	for i := 0; i < nPartitions; i++ {
		consumer := group.consumers[i%nConsumers]
		group.assignments[i] = consumer
	}
}

func (t *Topic) Produce(key string, msg []byte) (int, int64, error) {
	partitionIndex := 0

	if key != "" {
		partitionIndex = utils.Hash(key) % len(t.Partitions)
	}

	p := t.Partitions[partitionIndex]

	offset, err := p.Produce(msg)
	return partitionIndex, offset, err
}

func (t *Topic) Consume(partition int, offset int64) ([]byte, int64, error) {
	if partition >= len(t.Partitions) {
		return nil, 0, fmt.Errorf("invalid partition")
	}

	return t.Partitions[partition].Consume(offset)
}

func (t *Topic) ConsumeFromGroup(groupID, consumerID string) (string, int, int64, error) {
	group, exists := t.groups[groupID]
	if !exists {
		return "", 0, 0, fmt.Errorf("group not found")
	}

	for partition, consumer := range group.assignments {
		if consumer.ID == consumerID {
			offset := group.offsets[partition]

			msg, nextOffset, err := t.Partitions[partition].Consume(offset)
			if err != nil {
				continue
			}

			group.offsets[partition] = nextOffset

			return string(msg), partition, nextOffset, nil
		}
	}

	return "", 0, 0, fmt.Errorf("no partition assigned")
}
