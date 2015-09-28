package processor

type Task interface {
	Handle(pid int, result chan<- interface{}) //pid send for debug
}

type Collector interface {
	Handle(result <-chan interface{}) error
}