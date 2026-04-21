package ch4

import (
	"errors"
	"fmt"
	"sync"
)

// UnsyncRingMessages simulates pth_msg.c where message delivery is not synchronized.
func UnsyncRingMessages(threadCount int) ([]string, error) {
	if threadCount <= 0 {
		return nil, errors.New("threadCount must be positive")
	}

	messages := make([]string, threadCount)
	printed := make([]string, threadCount)
	var wg sync.WaitGroup

	for rank := 0; rank < threadCount; rank++ {
		wg.Add(1)
		go func(r int) {
			defer wg.Done()
			dest := (r + 1) % threadCount
			source := (r + threadCount - 1) % threadCount
			messages[dest] = fmt.Sprintf("Hello to %d from %d", dest, r)
			if messages[r] != "" {
				printed[r] = fmt.Sprintf("Thread %d > %s", r, messages[r])
			} else {
				printed[r] = fmt.Sprintf("Thread %d > No message from %d", r, source)
			}
		}(rank)
	}

	wg.Wait()
	return printed, nil
}

// SemaphoreRingMessages simulates pth_msg_sem.c using per-thread semaphores.
func SemaphoreRingMessages(threadCount int) ([]string, error) {
	if threadCount <= 0 {
		return nil, errors.New("threadCount must be positive")
	}

	messages := make([]string, threadCount)
	printed := make([]string, threadCount)
	semaphores := make([]chan struct{}, threadCount)
	for i := range semaphores {
		semaphores[i] = make(chan struct{}, 1)
	}

	var wg sync.WaitGroup
	for rank := 0; rank < threadCount; rank++ {
		wg.Add(1)
		go func(r int) {
			defer wg.Done()
			dest := (r + 1) % threadCount
			messages[dest] = fmt.Sprintf("Hello to %d from %d", dest, r)
			semaphores[dest] <- struct{}{}

			<-semaphores[r]
			printed[r] = fmt.Sprintf("Thread %d > %s", r, messages[r])
		}(rank)
	}

	wg.Wait()
	return printed, nil
}
