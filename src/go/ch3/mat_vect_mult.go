package ch3

import (
	"errors"
	"sync"
)

// MatVectMultParallel computes y = A*x for an m-by-n row-major matrix A.
func MatVectMultParallel(A []float64, m, n int, x []float64, workers int) ([]float64, error) {
	if m <= 0 || n <= 0 {
		return nil, errors.New("m and n must be positive")
	}
	if len(A) != m*n {
		return nil, errors.New("matrix size does not match m*n")
	}
	if len(x) != n {
		return nil, errors.New("vector size does not match n")
	}
	if workers <= 0 {
		return nil, errors.New("workers must be > 0")
	}

	y := make([]float64, m)
	jobs := make(chan int)
	var wg sync.WaitGroup

	worker := func() {
		defer wg.Done()
		for row := range jobs {
			offset := row * n
			sum := 0.0
			for col := 0; col < n; col++ {
				sum += A[offset+col] * x[col]
			}
			y[row] = sum
		}
	}

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go worker()
	}

	for row := 0; row < m; row++ {
		jobs <- row
	}
	close(jobs)
	wg.Wait()

	return y, nil
}
