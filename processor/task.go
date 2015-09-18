package processor

type Task interface {
	Handle(pid int) //pid send for debug
}