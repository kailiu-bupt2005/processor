package main
import (
	"fmt"
	p "huzznn/processor"
	"runtime"
//	"encoding/json"
	"errors"
	"encoding/json"
)

type taskTest struct {
	ID int
	Content string
	Weight int
	Point float64
}

func (task *taskTest)Handle(pid int, result chan<- interface{})  {
	task.ID = pid
	if result != nil {
		result <- task
	} else {
		js, _ := json.Marshal(task)
		fmt.Printf("[%v] %v\r\n", pid, string(js))
	}
}

func (task *taskTest)String() string {
	return fmt.Sprintf("c:%v w:%v p:%v", task.Content, task.Weight, task.Point)
}

type collecotTest struct {}

func (c *collecotTest)Handle(result <-chan interface{}) error {
	for {
		var ret interface{}
		var ok bool = true
		var task *taskTest = nil
		if ret, ok = <- result; !ok {
			fmt.Println("result chan close, so return")
			return nil
		}

		if task, ok = ret.(*taskTest); !ok {
			return errors.New("result chan input is not *taskTest");
		}

		fmt.Printf("task %v return\r\n", task.ID)
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() * 2)
	//有收集器
//	processor := p.NewProcessor(3, &collecotTest{})
//
//	for i := 0; i < 30; i++ {
//		task := new(taskTest)
//		task.Content = fmt.Sprintf("task %v", i + 1)
//		task.Weight = i + 1;
//		task.Point = float64(i) * 0.1
//		processor.AddTask(task)
//	}
//	processor.FinishAdd()

	//没有收集器
	processor := p.NewProcessor(3, nil)

	for i := 0; i < 30; i++ {
		task := new(taskTest)
		task.Content = fmt.Sprintf("task %v", i + 1)
		task.Weight = i + 1;
		task.Point = float64(i) * 0.1
		processor.AddTask(task)
	}
	processor.FinishAdd()

	fmt.Printf("task finish, error is %v.\r\n", processor.GetError())
}
