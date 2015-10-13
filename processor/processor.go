package processor
import (
	"sync"
	"sync/atomic"
)


type Processor struct {
	concurrency int
	inChan chan Task
	initOnce sync.Once
	taskNum int64
	taskDoneNum int64
	resultCollector Collector
	resultChan chan interface{}
	resultError error
	finishWait *sync.WaitGroup
	thatAll bool
}

func NewProcessor(concurrency int, resultCollector Collector) *Processor {
	p := &Processor{concurrency:concurrency, resultCollector: resultCollector}
	p.initOnce.Do(p.init)
	return p
}

func (p *Processor)init() {
	p.finishWait = new(sync.WaitGroup)
	if p.concurrency <= 0 {
		p.concurrency = 100
	}
	p.inChan = make(chan Task, p.concurrency)

	if p.resultCollector != nil {
		p.resultChan = make(chan interface{})
		go p.collet()
	}

	for i := 0; i < int(p.concurrency); i++ {
		go p.work(i + 1)
	}
}

func (p *Processor)FinishAdd() {
	close(p.inChan)
	p.thatAll = true
	p.finishWait.Wait()
}

func (p *Processor)AddTask(task Task) {
	p.initOnce.Do(p.init)
	if p.thatAll {
		panic("User told me it's finished, Why add task again.")
	}
	atomic.AddInt64(&p.taskNum, 1)
	p.finishWait.Add(1)
	p.inChan <- task
}

func (p *Processor)collet() {
	for result := range p.resultChan {
		err := p.resultCollector.Handle(result)
		if err != nil {
			p.resultError = err
		}
		p.finishWait.Done()
	}
}

func (p *Processor)GetError() error {
	return p.resultError
}

func (p *Processor)work(pid int) {
	for task := range p.inChan{
		task.Handle(pid, p.resultChan)
		atomic.AddInt64(&p.taskDoneNum, 1)
		if p.resultCollector == nil {
			p.finishWait.Done()
		}
	}
}