package ch3

import (
	"errors"
	"sync"
)

// VectorSumParallel adds two vectors using worker goroutines.
func VectorSumParallel(x, y []float64, workers int) ([]float64, error) {
	if len(x) != len(y) {
		return nil, errors.New("x and y must have the same length")
	}
	if workers <= 0 {
		return nil, errors.New("workers must be > 0")
	}

	n := len(x)
	z := make([]float64, n)
	jobs := make(chan int)
	var wg sync.WaitGroup

	worker := func() {
		defer wg.Done()
		for i := range jobs {
			z[i] = x[i] + y[i]
		}
	}

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go worker()
	}

	for i := 0; i < n; i++ {
		jobs <- i
	}
	close(jobs)
	wg.Wait()

	return z, nil
}
