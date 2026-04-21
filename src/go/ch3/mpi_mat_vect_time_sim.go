package ch3

import (
	"errors"
	"math/rand"
	"sync"
	"time"
)

// TimedMatVectResult carries the resulting y vector and elapsed wall time.
type TimedMatVectResult struct {
	Y       []float64
	Elapsed time.Duration
}

// MPIMatVectTimeSim simulates mpi_mat_vect_time.c behavior with generated data.
func MPIMatVectTimeSim(m, n, commSz int, seed int64) (*TimedMatVectResult, error) {
	if commSz <= 0 {
		return nil, errors.New("commSz must be positive")
	}
	if m <= 0 || n <= 0 {
		return nil, errors.New("m and n must be positive")
	}
	if m%commSz != 0 || n%commSz != 0 {
		return nil, errors.New("m and n must be divisible by commSz")
	}

	localM := m / commSz
	localN := n / commSz
	localA := make([][]float64, commSz)
	localX := make([][]float64, commSz)

	for rank := 0; rank < commSz; rank++ {
		rng := rand.New(rand.NewSource(seed + int64(rank)))
		blockA := make([]float64, localM*n)
		for i := range blockA {
			blockA[i] = rng.Float64()
		}
		blockX := make([]float64, localN)
		for i := range blockX {
			blockX[i] = rng.Float64()
		}
		localA[rank] = blockA
		localX[rank] = blockX
	}

	xGlobal := make([]float64, 0, n)
	for rank := 0; rank < commSz; rank++ {
		xGlobal = append(xGlobal, localX[rank]...)
	}

	y := make([]float64, m)
	var wg sync.WaitGroup
	start := time.Now()
	for rank := 0; rank < commSz; rank++ {
		wg.Add(1)
		go func(r int) {
			defer wg.Done()
			for i := 0; i < localM; i++ {
				offset := i * n
				sum := 0.0
				for j := 0; j < n; j++ {
					sum += localA[r][offset+j] * xGlobal[j]
				}
				y[r*localM+i] = sum
			}
		}(rank)
	}
	wg.Wait()
	elapsed := time.Since(start)

	return &TimedMatVectResult{Y: y, Elapsed: elapsed}, nil
}
