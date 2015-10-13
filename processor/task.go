package processor

type Task interface {
	Handle(pid int, resultChan chan<- interface{}) //pid send for debug
}

type Collector interface {
	Handle(result interface{}) error
}