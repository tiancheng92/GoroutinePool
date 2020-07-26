package main

import (
	"GoroutinePool"
	"errors"
	"fmt"
	"time"
)

func CalRemainder(a int64, b int64) error {
	time.Sleep(time.Second)
	if b == 0 {
		return errors.New("The divisor cannot be zero\n")
	}
	fmt.Printf("%d %% %d = %d\n", a, b, a%b)
	return nil
}

func HandleError(err error) {
	fmt.Printf("发生错误：%s", err)
}

func FinishCallback() {
	fmt.Println("所有协程均执行完毕。")
}

func main() {
	var p GoroutinePool.Pool
	defer p.Stop()
	numbers := []int64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	p.Init(3, len(numbers))
	
	for index := range numbers {
		number := numbers[index]
		p.AddTask(func() error {
			return CalRemainder(number, 3)
		})
	}
	p.SetHandleError(HandleError)
	p.SetFinishCallback(FinishCallback)
	p.Start()
}
