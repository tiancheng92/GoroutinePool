package main

import (
	"GoroutinePool"
	"context"
	"fmt"
	"time"
)

func MockTask(index int) error {
	time.Sleep(time.Second)
	fmt.Printf("task %d finish\n", index)
	return nil
}

func HandleError(err error) {
	fmt.Printf("发生错误：%s", err)
}

func FinishCallback() {
	fmt.Println("所有协程均执行完毕。")
}

func LimitedTaskCount() {
	var p GoroutinePool.Pool
	defer p.Stop()
	taskNumbers := 10
	p.Init(3, taskNumbers)

	for i := 1; i <= taskNumbers; i++ {
		index := i
		p.AddTask(func() error {
			return MockTask(index)
		})
	}
	p.SetHandleError(HandleError)
	p.SetFinishCallback(FinishCallback)
	p.Start()
}

func InfiniteTaskCount(ctx context.Context) {
	var indexChan = make(chan int)
	var p GoroutinePool.Pool

	p.InitForInfinite(2)
	p.SetHandleError(HandleError)

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
				p.AddTaskForInfinite(func() error {
					return MockTask(<-indexChan)
				})
			}
		}
	}()

	p.SetFinishCallback(func() {
		fmt.Println("已经达到设定的任务执行时长")
	})
	p.Start()
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	InfiniteTaskCount(ctx)

	LimitedTaskCount()
}
