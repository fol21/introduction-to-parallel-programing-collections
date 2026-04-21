package ch4

import (
	"errors"
	"fmt"
	"sort"
	"sync"
)

// HelloMessages simulates pth_hello output using goroutines.
func HelloMessages(threadCount int) ([]string, error) {
	if threadCount <= 0 {
		return nil, errors.New("threadCount must be positive")
	}

	out := make(chan string, threadCount)
	var wg sync.WaitGroup

	for rank := 0; rank < threadCount; rank++ {
		wg.Add(1)
		go func(r int) {
			defer wg.Done()
			out <- fmt.Sprintf("Hello from thread %d of %d", r, threadCount)
		}(rank)
	}

	wg.Wait()
	close(out)

	lines := make([]string, 0, threadCount+1)
	lines = append(lines, "Hello from the main thread")
	threadLines := make([]string, 0, threadCount)
	for line := range out {
		threadLines = append(threadLines, line)
	}
	sort.Strings(threadLines)
	lines = append(lines, threadLines...)
	return lines, nil
}
