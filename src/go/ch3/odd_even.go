package ch3

import (
	"errors"
	"math/rand"
	"sync"
)

// GenerateList creates n pseudo-random integers in [0, maxValue).
func GenerateList(n, maxValue int, seed int64) ([]int, error) {
	if n < 0 {
		return nil, errors.New("n must be non-negative")
	}
	if maxValue <= 0 {
		return nil, errors.New("maxValue must be positive")
	}

	rng := rand.New(rand.NewSource(seed))
	out := make([]int, n)
	for i := 0; i < n; i++ {
		out[i] = rng.Intn(maxValue)
	}
	return out, nil
}

// OddEvenSortParallel sorts the slice in-place with odd-even transposition sort.
func OddEvenSortParallel(a []int, workers int) error {
	if workers <= 0 {
		return errors.New("workers must be > 0")
	}

	n := len(a)
	for phase := 0; phase < n; phase++ {
		start := 0
		if phase%2 == 1 {
			start = 1
		}

		jobs := make(chan int)
		var wg sync.WaitGroup

		worker := func() {
			defer wg.Done()
			for left := range jobs {
				right := left + 1
				if a[left] > a[right] {
					a[left], a[right] = a[right], a[left]
				}
			}
		}

		for i := 0; i < workers; i++ {
			wg.Add(1)
			go worker()
		}

		for left := start; left+1 < n; left += 2 {
			jobs <- left
		}
		close(jobs)
		wg.Wait()
	}

	return nil
}
