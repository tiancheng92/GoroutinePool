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

func (p *Pool) Init(goroutineNumber int, taskNumber int) {
	p.GoroutineNumber = goroutineNumber
	p.Wg.Add(taskNumber + 1)
	p.Task = make(chan func() error, taskNumber)
}

func (p *Pool) InitForInfinite(goroutineNumber int) {
	p.GoroutineNumber = goroutineNumber
	p.Wg.Add(1)
	p.Task = make(chan func() error)
}

func (p *Pool) Start() {
	var once sync.Once
	for i := 0; i < p.GoroutineNumber; i++ {
		go func() {
			for {
				task, ok := <-p.Task
				if !ok {
					break
				}
				once.Do(func() {
					p.Wg.Done()
				})
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

func (p *Pool) Stop() {
	close(p.Task)
}

func (p *Pool) AddTask(task func() error) {
	p.Task <- task
}

func (p *Pool) AddTaskForInfinite(task func() error) {
	p.Wg.Add(1)
	p.Task <- task
}

func (p *Pool) SetFinishCallback(fun func()) {
	p.FinishCallback = fun
}

func (p *Pool) SetHandleError(fun func(error)) {
	p.HandleError = fun
}
