package hw05parallelexecution

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if m < 0 {
		return fmt.Errorf("макс колво ошибок - ноль: %w", ErrErrorsLimitExceeded)
	}
	taskPool := make(chan Task)

	var errTasksCount int64

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func(wg2 *sync.WaitGroup) {
		defer wg2.Done()
		for _, task := range tasks {
			if atomic.LoadInt64(&errTasksCount) >= int64(m) {
				close(taskPool)
				return
			}
			taskPool <- task
		}
		close(taskPool)
	}(&wg)

	worker := func() {
		wg.Add(1)
		go func(wg1 *sync.WaitGroup) {
			defer wg1.Done()
			for {
				if fn, ok := <-taskPool; ok {
					if errmsg := fn(); errmsg != nil {
						atomic.AddInt64(&errTasksCount, 1)
					}
				} else {
					return
				}
			}
		}(&wg)
	}
	for i := 1; i <= n; i++ {
		worker()
	}
	wg.Wait()
	if errTasksCount >= int64(m) {
		return fmt.Errorf("неправильное количество: %w", ErrErrorsLimitExceeded)
	}
	return nil
}
