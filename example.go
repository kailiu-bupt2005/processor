package main
import (
	"code.google.com/p/go/src/pkg/fmt"
	"huzznn/processor"
"runtime"
	"code.google.com/p/go/src/pkg/encoding/json"
)

type taskTest struct {
	Content string
	Weight int
	Point float64
}

func (task *taskTest)Handle(pid int)  {
	js, _ := json.Marshal(task)
	fmt.Printf("[%v] %v\r\n", pid, string(js))
}

func (task *taskTest)String() string {
	return fmt.Sprintf("c:%v w:%v p:%v", task.Content, task.Weight, task.Point)
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() * 2)
	processor := processor.NewStreamProcessor(3)

	for i := 0; i < 30; i++ {
		task := new(taskTest)
		task.Content = fmt.Sprintf("task %v", i + 1)
		task.Weight = i + 1;
		task.Point = float64(i) * 0.1
		processor.AddTask(task)
	}
	processor.FinishAdd()

	fmt.Println("task finish.\r\n")
}
