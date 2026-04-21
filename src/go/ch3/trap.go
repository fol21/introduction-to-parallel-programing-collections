package ch3

import (
	"errors"
	"sync"
)

func f(x float64) float64 {
	return x * x
}

// TrapParallel estimates integral(a,b) of f(x)=x^2 using n trapezoids.
func TrapParallel(a, b float64, n, workers int) (float64, error) {
	if n <= 0 {
		return 0, errors.New("n must be positive")
	}
	if workers <= 0 {
		return 0, errors.New("workers must be > 0")
	}

	h := (b - a) / float64(n)
	jobs := make(chan int)
	partials := make(chan float64, workers)
	var wg sync.WaitGroup

	worker := func() {
		defer wg.Done()
		local := 0.0
		for k := range jobs {
			local += f(a + float64(k)*h)
		}
		partials <- local
	}

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go worker()
	}

	for k := 1; k <= n-1; k++ {
		jobs <- k
	}
	close(jobs)
	wg.Wait()
	close(partials)

	integral := (f(a) + f(b)) / 2.0
	for part := range partials {
		integral += part
	}

	return integral * h, nil
}
