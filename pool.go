package GoroutinePool

import (
	"fmt"
	"sync"
)

type Pool struct {
	CoroutinesCount int               // 协程数量
	TasksCount      int               // 任务数量
	Task            chan func() error // 存储任务的管道
	FinishCallback  func()            // 任务全部完成后回调的函数
	HandleError     func(error)       // 任务的错误处理
	Wg              sync.WaitGroup    // 阻塞
	Once            sync.Once         // 用于确保Stop只执行一次
}

// New 新建Pool对象
func New() *Pool {
	p := new(Pool)
	p.Task = make(chan func() error)
	p.Wg.Add(1)
	return p
}

// SetCoroutinesCount 设定协程数量
func (p *Pool) SetCoroutinesCount(count int) *Pool {
	p.CoroutinesCount = count
	return p
}

// SetCoroutinesCount 设定协程数量
func (p *Pool) SetTasksCount(count int) *Pool {
	p.TasksCount = count
	p.Wg.Add(count - 1)
	return p
}

// SetFinishCallback 设定任务结束后的回调函数
func (p *Pool) SetFinishCallback(function func()) *Pool {
	p.FinishCallback = function
	return p
}

// SetHandleError 设定任务执行失败时的回调函数
func (p *Pool) SetHandleError(function func(error)) *Pool {
	p.HandleError = function
	return p
}

// AddTask 添加任务
func (p *Pool) AddTask(task func() error) {
	p.Task <- task
}

// Start 开始执行任务
func (p *Pool) Start() {
	for i := 0; i < p.CoroutinesCount; i++ {
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
				if p.TasksCount > 0 {
					p.Wg.Done()
				}
			}
		}()
	}
	p.Wg.Wait()
	if p.FinishCallback != nil {
		p.FinishCallback()
	}

}

// Stop 终止任务
func (p *Pool) Stop() {
	p.Once.Do(
		func() {
			if p.TasksCount == 0 {
				p.Wg.Done()
			}
			close(p.Task)
		})
}
