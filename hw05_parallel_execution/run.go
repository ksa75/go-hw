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
		m = -m
	}
	done := make(chan struct{})
	taskPool := make(chan Task)
	defer close(done)
	defer close(taskPool)

	var errTasksCount int64

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func(wg2 *sync.WaitGroup) {
		defer wg2.Done()
		for _, task := range tasks {
			if atomic.LoadInt64(&errTasksCount) >= int64(m) {
				// fmt.Println("заканчиваем передачу по ошибке")
				for i := 1; i <= n; i++ {
					done <- struct{}{}
				}
				return
			}
			taskPool <- task
		}
		// fmt.Println("заканчиваем передачу нормально")
		for i := 1; i <= n; i++ {
			done <- struct{}{}
		}
	}(&wg)

	worker := func(done chan struct{}) {
		wg.Add(1)
		go func(wg1 *sync.WaitGroup) {
			defer wg1.Done()
			for {
				select {
				case <-done:
					// fmt.Println("worker: done")
					return
				case fn, ok := <-taskPool:
					_ = ok
					// fmt.Println("processed ", ok)
					if errmsg := fn(); errmsg != nil {
						// fmt.Println(errmsg)
						atomic.AddInt64(&errTasksCount, 1)
					}
				}
			}
		}(&wg)
	}

	// обрабатываем очередь n воркерами
	for i := 1; i <= n; i++ {
		worker(done)
	}

	wg.Wait()

	// fmt.Printf("сигналов осталось %v\n", len(done))
	// fmt.Printf("задач осталось %v\n", len(taskPool))

	// fmt.Println(errTasksCount)
	if errTasksCount >= int64(m) {
		return fmt.Errorf("неправильное количество: %w", ErrErrorsLimitExceeded)
	}
	return nil
}
