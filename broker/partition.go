package broker

import "mini-kafka/storage"

type Partition struct {
	log *storage.Log
}

func NewPartition(path string) (*Partition, error) {
	log, err := storage.NewLog(path)
	if err != nil {
		return nil, err
	}

	return &Partition{log: log}, nil
}

func (p *Partition) Produce(msg []byte) (int64, error) {
	return p.log.Append(msg)
}

func (p *Partition) Consume(offset int64) ([]byte, int64, error) {
	return p.log.Read(offset)
}
