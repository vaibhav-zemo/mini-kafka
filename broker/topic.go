package broker

type Topic struct {
	name      string
	Partition *Partition
}

func NewTopic(name string) (*Topic, error) {
	p, err := NewPartition("data/" + name + ".log")
	if err != nil {
		return nil, err
	}

	return &Topic{
		name:      name,
		Partition: p,
	}, nil
}
