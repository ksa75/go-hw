package hw05parallelexecution

import (
	"errors"
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("no tasks", func(t *testing.T) {
		tasks := make([]Task, 0)
		workersCount := 3
		maxErrorsCount := 1
		err := Run(tasks, workersCount, maxErrorsCount)
		require.Truef(t, errors.Is(err, nil), "actual err - %v", err)
	})

	t.Run("'максимум 0 ошибок' - значит функция всегда будет возвращать ErrErrorsLimitExceeded", func(t *testing.T) {
		tasksCount := 8
		tasks := make([]Task, 0, tasksCount)
		//
		var runTasksCount int32
		//
		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}
		//
		workersCount := 3
		maxErrorsCount := -3
		err := Run(tasks, workersCount, maxErrorsCount)
		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
	})

	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)
		//
		var runTasksCount int32
		//
		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}
		//
		workersCount := 10
		maxErrorsCount := 23
		err := Run(tasks, workersCount, maxErrorsCount)
		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
	})

	t.Run("if errors were in less than M tasks at any order, then no error", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)
		workersCount := 10
		maxErrorsCount := 23

		var runTasksCount, errCount int64

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt64(&runTasksCount, 1)
				if int(atomic.LoadInt64(&errCount)) < maxErrorsCount-1 {
					if rand.Intn(2) == 1 {
						atomic.AddInt64(&errCount, 1)
						return err
					}
				}
				return nil
			})
		}
		err := Run(tasks, workersCount, maxErrorsCount)
		require.Truef(t, errors.Is(err, nil), "actual err - %v", err)
		require.LessOrEqual(t, runTasksCount, int64(tasksCount), "extra tasks were started")
		fmt.Printf("задач всего %v - выполнено %v - ошибок %v\n", tasksCount, runTasksCount, errCount)
	})

	t.Run("tasks without errors", func(t *testing.T) {
		tasksCount := 45
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep

			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 1

		start := time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.NoError(t, err)

		require.Truef(t, errors.Is(err, nil), "actual err - %v", err)
		require.Equal(t, int32(tasksCount), runTasksCount, "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
		fmt.Printf("задач всего %v - выполнено %v\n", tasksCount, runTasksCount)
		fmt.Printf("машинного времени затрачено %v - фактически прошло %v\n", sumTime, elapsedTime)
	})
}
