package hw05parallelexecution

import (
	"errors"
	"fmt"
	"time"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	done := make(chan struct{}, n)
	fin := make(chan struct{}, len(tasks))
	taskPool := make(chan Task)
	waitSleepers := make(chan int, len(tasks))

	// var errTasksCount int32

	// заполняем канал
	go func(fin <-chan struct{}, taskPool chan Task) {
		for _, task := range tasks {
			for {
				select {
				case <-fin:
					return
				case taskPool <- task:
				}
			}
		}
	}(fin, taskPool)

	worker := func(done <-chan struct{}, taskPool chan Task, i int) bool {
		var err bool
		go func() {
			for {
				select {
				case <-done:
					fmt.Println("worker: done")
					return
				case fn, ok := <-taskPool:
					_ = ok
					fmt.Println("processed ", fn(), ok, i)
					// if i := fn(); i != nil {
					// fmt.Println("error")
					// if atomic.AddInt32(&errTasksCount, 1) == int32(m) {
					// err = true
					// close(taskPool)
					// }
					// }
					waitSleepers <- i
				}
			}
		}()
		return err
	}

	//запускаем воркеры
	for i := 1; i <= n; i++ {
		_ = worker(done, taskPool, i)
	}
	//снимаем блокировку завершения блока select
	time.Sleep(time.Millisecond * 100)
	for i := 1; i <= len(tasks); i++ {
		<-waitSleepers
	}
	//сигнализируем поставщику завершение работы
	for i := 1; i <= len(tasks); i++ {
		fin <- struct{}{}
	}
	//сигнализируем воркерам завершение работы
	for i := 1; i <= n; i++ {
		done <- struct{}{}
	}
	return nil
}
