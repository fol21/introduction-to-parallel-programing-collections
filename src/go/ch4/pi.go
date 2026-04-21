package ch4

import (
	"errors"
	"math"
	"sync"
)

// SerialPi estimates pi with n terms of the Gregory-Leibniz series.
func SerialPi(n int64) (float64, error) {
	if n <= 0 {
		return 0, errors.New("n must be positive")
	}

	sum := 0.0
	factor := 1.0
	for i := int64(0); i < n; i++ {
		sum += factor / float64(2*i+1)
		factor = -factor
	}
	return 4.0 * sum, nil
}

// PiEstimateMutex estimates pi using goroutines and a mutex-protected accumulation.
func PiEstimateMutex(threadCount int, n int64) (float64, error) {
	if threadCount <= 0 {
		return 0, errors.New("threadCount must be positive")
	}
	if n <= 0 {
		return 0, errors.New("n must be positive")
	}
	if n%int64(threadCount) != 0 {
		return 0, errors.New("n must be divisible by threadCount")
	}

	var sum float64
	var mu sync.Mutex
	var wg sync.WaitGroup

	myN := n / int64(threadCount)
	for rank := 0; rank < threadCount; rank++ {
		wg.Add(1)
		go func(r int) {
			defer wg.Done()
			myFirst := int64(r) * myN
			myLast := myFirst + myN
			factor := 1.0
			if myFirst%2 == 1 {
				factor = -1.0
			}
			local := 0.0
			for i := myFirst; i < myLast; i++ {
				local += factor / float64(2*i+1)
				factor = -factor
			}
			mu.Lock()
			sum += local
			mu.Unlock()
		}(rank)
	}

	wg.Wait()
	return 4.0 * sum, nil
}

// PiEstimateChannel estimates pi using channel-based reduction.
func PiEstimateChannel(threadCount int, n int64) (float64, error) {
	if threadCount <= 0 {
		return 0, errors.New("threadCount must be positive")
	}
	if n <= 0 {
		return 0, errors.New("n must be positive")
	}
	if n%int64(threadCount) != 0 {
		return 0, errors.New("n must be divisible by threadCount")
	}

	partials := make(chan float64, threadCount)
	myN := n / int64(threadCount)
	var wg sync.WaitGroup

	for rank := 0; rank < threadCount; rank++ {
		wg.Add(1)
		go func(r int) {
			defer wg.Done()
			myFirst := int64(r) * myN
			myLast := myFirst + myN
			factor := 1.0
			if myFirst%2 == 1 {
				factor = -1.0
			}
			local := 0.0
			for i := myFirst; i < myLast; i++ {
				local += factor / float64(2*i+1)
				factor = -factor
			}
			partials <- local
		}(rank)
	}

	wg.Wait()
	close(partials)

	sum := 0.0
	for p := range partials {
		sum += p
	}
	return 4.0 * sum, nil
}

// PiReference returns the value from the math library identity used by the C code.
func PiReference() float64 {
	return 4.0 * math.Atan(1.0)
}
