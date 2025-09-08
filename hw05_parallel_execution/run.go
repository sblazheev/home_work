package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func Run(tasks []Task, n, m int) error {
	wg := sync.WaitGroup{}
	lockShift := sync.Mutex{}
	lockErrorCount := sync.Mutex{}
	errorCount := 0
	shift := -1

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				if m > 0 {
					lockErrorCount.Lock()
					if errorCount >= m {
						lockErrorCount.Unlock()
						return
					}
					lockErrorCount.Unlock()
				}
				lockShift.Lock()
				shift++

				if shift < len(tasks) {
					taskIndex := shift
					lockShift.Unlock()
					er := tasks[taskIndex]()
					if er != nil && m > 0 {
						lockErrorCount.Lock()
						errorCount++
						lockErrorCount.Unlock()
					}
				} else {
					lockShift.Unlock()
					return
				}
			}
		}()
	}

	wg.Wait()

	if m > 0 && errorCount >= m {
		return ErrErrorsLimitExceeded
	}
	return nil
}
