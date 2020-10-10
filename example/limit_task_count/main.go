package main

import (
	"GoroutinePool"
	"errors"
	"fmt"
	"time"
)

func MockTask(index int) error {
	time.Sleep(time.Second)
	fmt.Printf("task %d finish\n", index)
	if index%7 == 0 {
		return errors.New("mock error")
	}
	return nil
}

func HandleError(err error) {
	fmt.Printf("发生错误：%s\n", err)
}

func FinishCallback() {
	fmt.Println("所有协程均执行完毕。")
}

func main() {
	p := GoroutinePool.New()
	defer p.Stop()

	p.SetCoroutinesCount(3).
		SetHandleError(HandleError).
		SetFinishCallback(FinishCallback)

	go func() {
		for i := 1; i <= 10; i++ {
			index := i
			p.AddTask(func() error {
				return MockTask(index)
			})
		}
	}()

	p.Start()
}
