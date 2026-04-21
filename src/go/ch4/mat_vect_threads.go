package ch4

import "errors"

// MatVectThreads computes y=A*x where A is row-major m-by-n.
// Work is split by contiguous row blocks across threadCount workers.
func MatVectThreads(A []float64, m, n int, x []float64, threadCount int) ([]float64, error) {
	if threadCount <= 0 {
		return nil, errors.New("threadCount must be positive")
	}
	if m <= 0 || n <= 0 {
		return nil, errors.New("m and n must be positive")
	}
	if len(A) != m*n {
		return nil, errors.New("matrix size mismatch")
	}
	if len(x) != n {
		return nil, errors.New("vector size mismatch")
	}
	if m%threadCount != 0 {
		return nil, errors.New("m must be divisible by threadCount")
	}

	y := make([]float64, m)
	localM := m / threadCount
	done := make(chan struct{}, threadCount)

	for rank := 0; rank < threadCount; rank++ {
		go func(r int) {
			first := r * localM
			last := first + localM
			for i := first; i < last; i++ {
				sum := 0.0
				row := i * n
				for j := 0; j < n; j++ {
					sum += A[row+j] * x[j]
				}
				y[i] = sum
			}
			done <- struct{}{}
		}(rank)
	}

	for i := 0; i < threadCount; i++ {
		<-done
	}
	return y, nil
}
