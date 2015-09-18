package processor
import (
	"sync"
	"sync/atomic"
)


type Processor struct {
	concurrency int32
	inChan chan Task
	finishChan chan int32
	thatAll bool
	initOnce sync.Once
	finishOnce sync.Once
	taskNum int64
}

func NewStreamProcessor(concurrency int32) *Processor {
	p := &Processor{concurrency:concurrency}
	p.initOnce.Do(p.init)
	return p
}

func (p *Processor)init() {
	if p.concurrency <= 0 {
		p.concurrency = 100
	}
	p.inChan = make(chan Task, p.concurrency)
	p.finishChan = make(chan int32)
	for i := 0; i < int(p.concurrency); i++ {
		go p.work(i + 1)
	}
}

func (p *Processor)FinishAdd() {
	p.thatAll = true
	if atomic.LoadInt64(&p.taskNum) == 0 {
		close(p.inChan)
		close(p.finishChan)
		return
	}
	<- p.finishChan
}

func (p *Processor)AddTask(task Task) {
	p.initOnce.Do(p.init)
	if p.thatAll {
		panic("User told me it's finshed, Why add task again.")
	}
	atomic.AddInt64(&p.taskNum, 1)
	p.inChan <- task
}


func (p *Processor)work(pid int) {
	var task Task
	var ok bool
	for {
		if task, ok = <-p.inChan; !ok {
			return
		}
		task.Handle(pid)
		atomic.AddInt64(&p.taskNum, -1)
		if atomic.LoadInt64(&p.taskNum) <= 0 && p.thatAll  {
			close(p.inChan)
			p.finishOnce.Do(func() {
				p.finishChan <- 1
			})
		}
	}
}