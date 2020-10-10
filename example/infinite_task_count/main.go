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

func main() {
	var indexChan = make(chan int)
	defer close(indexChan)

	p := GoroutinePool.New()
	defer p.Stop()

	p.SetCoroutinesCount(2).
		SetHandleError(HandleError)

	go func() {
		i := 1
		for {
			indexChan <- i
			i++
		}
	}()

	go func() {
		for {
			p.AddTask(func() error {
				return MockTask(<-indexChan)
			})
		}
	}()
	p.Start()
}
