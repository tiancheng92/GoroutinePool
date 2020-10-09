package GoroutinePool

import (
	"fmt"
	"sync"
)

type Pool struct {
	GoroutineNumber int
	Task            chan func() error
	FinishCallback  func()
	HandleError     func(error)
	Wg              sync.WaitGroup
}

type PoolForInfinite struct {
	GoroutineNumber int
	Task            chan func() error
	FinishCallback  func()
	HandleError     func(error)
}

func (p *Pool) Init(goroutineNumber int, taskNumber int) {
	p.GoroutineNumber = goroutineNumber
	p.Wg.Add(taskNumber)
	p.Task = make(chan func() error, taskNumber)
}

func (p *PoolForInfinite) Init(goroutineNumber int) {
	p.GoroutineNumber = goroutineNumber
	p.Task = make(chan func() error)
}

func (p *Pool) Start() {
	for i := 0; i < p.GoroutineNumber; i++ {
		go func() {
			for {
				task, ok := <-p.Task
				if !ok {
					break
				}
				err := task()
				if err != nil {
					if p.HandleError != nil {
						p.HandleError(err)
					} else {
						fmt.Println(err)
					}
				}
				p.Wg.Done()
			}
		}()
	}
	p.Wg.Wait()
	if p.FinishCallback != nil {
		p.FinishCallback()
	}
}

func (p *PoolForInfinite) Start() {
	for i := 0; i < p.GoroutineNumber; i++ {
		go func() {
			for {
				task := <-p.Task
				err := task()
				if err != nil {
					if p.HandleError != nil {
						p.HandleError(err)
					} else {
						fmt.Println(err)
					}
				}
			}
		}()
	}
}

func (p *Pool) Stop() {
	close(p.Task)
}

func (p *Pool) AddTask(task func() error) {
	p.Task <- task
}

func (p *PoolForInfinite) AddTask(task func() error) {
	p.Task <- task
}

func (p *Pool) SetFinishCallback(fun func()) {
	p.FinishCallback = fun
}

func (p *Pool) SetHandleError(fun func(error)) {
	p.HandleError = fun
}

func (p *PoolForInfinite) SetHandleError(fun func(error)) {
	p.HandleError = fun
}
