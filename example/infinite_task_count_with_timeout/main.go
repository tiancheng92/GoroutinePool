package main

import (
	"GoroutinePool"
	"context"
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
	fmt.Println("已经达到设定的任务执行时长。")
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var indexChan = make(chan int)
	defer close(indexChan)

	p := GoroutinePool.New()
	defer p.Stop()

	p.SetCoroutinesCount(2).
		SetHandleError(HandleError).
		SetFinishCallback(FinishCallback)

	go func() {
		i := 1
		for {
			indexChan <- i
			i++
		}
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				break
			default:
				p.AddTask(func() error {
					return MockTask(<-indexChan)
				})
			}
		}
	}()
	p.Start()
}
