package ch2

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"sync"
)

// Histogram stores bin boundaries and counts.
type Histogram struct {
	BinMaxes  []float64
	BinCounts []int
	MinMeas   float64
}

// GenData generates random values in the range [minMeas, maxMeas).
func GenData(minMeas, maxMeas float64, dataCount int, seed int64) ([]float64, error) {
	if dataCount < 0 {
		return nil, errors.New("dataCount must be non-negative")
	}
	if maxMeas <= minMeas {
		return nil, errors.New("maxMeas must be greater than minMeas")
	}

	rng := rand.New(rand.NewSource(seed))
	data := make([]float64, dataCount)
	for i := range data {
		data[i] = minMeas + (maxMeas-minMeas)*rng.Float64()
	}
	return data, nil
}

// GenBins computes each bin max and initializes all counts to zero.
func GenBins(minMeas, maxMeas float64, binCount int) ([]float64, []int, error) {
	if binCount <= 0 {
		return nil, nil, errors.New("binCount must be > 0")
	}
	if maxMeas <= minMeas {
		return nil, nil, errors.New("maxMeas must be greater than minMeas")
	}

	binWidth := (maxMeas - minMeas) / float64(binCount)
	binMaxes := make([]float64, binCount)
	binCounts := make([]int, binCount)

	for i := 0; i < binCount; i++ {
		binMaxes[i] = minMeas + float64(i+1)*binWidth
	}

	return binMaxes, binCounts, nil
}

// WhichBin finds the bin for a measurement using binary search.
func WhichBin(data float64, binMaxes []float64, minMeas float64) (int, error) {
	if len(binMaxes) == 0 {
		return -1, errors.New("binMaxes cannot be empty")
	}

	bottom := 0
	top := len(binMaxes) - 1

	for bottom <= top {
		mid := (bottom + top) / 2
		binMax := binMaxes[mid]
		binMin := minMeas
		if mid > 0 {
			binMin = binMaxes[mid-1]
		}

		if data >= binMax {
			bottom = mid + 1
		} else if data < binMin {
			top = mid - 1
		} else {
			return mid, nil
		}
	}

	return -1, fmt.Errorf("data=%f does not belong to any bin", data)
}

// BuildHistogramParallel counts measurements in parallel using goroutines, channels, and sync.
func BuildHistogramParallel(data []float64, binCount int, minMeas, maxMeas float64, workers int) (*Histogram, error) {
	if workers <= 0 {
		return nil, errors.New("workers must be > 0")
	}

	binMaxes, binCounts, err := GenBins(minMeas, maxMeas, binCount)
	if err != nil {
		return nil, err
	}

	jobs := make(chan float64)
	errCh := make(chan error, 1)
	var wg sync.WaitGroup
	var mu sync.Mutex

	worker := func() {
		defer wg.Done()
		for measurement := range jobs {
			bin, findErr := WhichBin(measurement, binMaxes, minMeas)
			if findErr != nil {
				select {
				case errCh <- findErr:
				default:
				}
				continue
			}

			mu.Lock()
			binCounts[bin]++
			mu.Unlock()
		}
	}

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go worker()
	}

	for _, measurement := range data {
		jobs <- measurement
	}
	close(jobs)
	wg.Wait()

	select {
	case countErr := <-errCh:
		return nil, countErr
	default:
	}

	return &Histogram{
		BinMaxes:  binMaxes,
		BinCounts: binCounts,
		MinMeas:   minMeas,
	}, nil
}

// Print writes a text histogram to the provided writer.
func (h *Histogram) Print(w io.Writer) {
	for i := range h.BinMaxes {
		binMin := h.MinMeas
		if i > 0 {
			binMin = h.BinMaxes[i-1]
		}
		fmt.Fprintf(w, "%.3f-%.3f:\t", binMin, h.BinMaxes[i])
		for j := 0; j < h.BinCounts[i]; j++ {
			fmt.Fprint(w, "X")
		}
		fmt.Fprintln(w)
	}
}
