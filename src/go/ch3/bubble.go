package ch3

import (
	"errors"
	"sync"
)

// BubbleSort sorts an integer slice in ascending order in-place.
func BubbleSort(a []int) {
	n := len(a)
	for listLength := n; listLength >= 2; listLength-- {
		for i := 0; i < listLength-1; i++ {
			if a[i] > a[i+1] {
				a[i], a[i+1] = a[i+1], a[i]
			}
		}
	}
}

// BubbleSortParallel performs parallel compare/swap phases to reach sorted order.
func BubbleSortParallel(a []int, workers int) error {
	if workers <= 0 {
		return errors.New("workers must be > 0")
	}
	n := len(a)
	if n < 2 {
		return nil
	}

	for phase := 0; phase < n; phase++ {
		start := 0
		if phase%2 == 1 {
			start = 1
		}
		jobs := make(chan int)
		var wg sync.WaitGroup

		worker := func() {
			defer wg.Done()
			for i := range jobs {
				if a[i] > a[i+1] {
					a[i], a[i+1] = a[i+1], a[i]
				}
			}
		}

		for w := 0; w < workers; w++ {
			wg.Add(1)
			go worker()
		}
		for i := start; i+1 < n; i += 2 {
			jobs <- i
		}
		close(jobs)
		wg.Wait()
	}
	return nil
}
