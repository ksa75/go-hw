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
	var errTasksCount int64
	taskPool := make(chan Task)
	wg := sync.WaitGroup{}
	worker := func() {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range taskPool {
				if fn, ok := <-taskPool; ok {
					if errmsg := fn(); errmsg != nil {
						atomic.AddInt64(&errTasksCount, 1)
					}
				} else {
					return
				}
			}
		}()
	}

	for i := 1; i <= n; i++ {
		worker()
	}
	for _, task := range tasks {
		if atomic.LoadInt64(&errTasksCount) >= int64(m) {
			break
		}
		taskPool <- task
	}
	close(taskPool)
	wg.Wait()
	if errTasksCount >= int64(m) {
		return fmt.Errorf("неправильное количество: %w", ErrErrorsLimitExceeded)
	}
	return nil
}
