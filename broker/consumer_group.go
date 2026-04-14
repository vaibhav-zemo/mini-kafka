package broker

type Consumer struct {
	ID string
}

type ConsumerGroup struct {
	ID          string
	consumers   []*Consumer
	assignments map[int]*Consumer
	offsets     map[int]int64
}
