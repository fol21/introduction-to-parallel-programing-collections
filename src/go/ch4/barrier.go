package ch4

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

// BusyBarrier mirrors the busy-wait style from pth_busy_bar.c.
type BusyBarrier struct {
	threadCount int
	counts      []atomic.Int32
}

func newBusyBarrier(threadCount, barrierCount int) (*BusyBarrier, error) {
	if threadCount <= 0 {
		return nil, errors.New("threadCount must be positive")
	}
	if barrierCount <= 0 {
		return nil, errors.New("barrierCount must be positive")
	}
	return &BusyBarrier{threadCount: threadCount, counts: make([]atomic.Int32, barrierCount)}, nil
}

func (b *BusyBarrier) wait(i int) {
	b.counts[i].Add(1)
	for b.counts[i].Load() < int32(b.threadCount) {
	}
}

// CondBarrier mirrors condition-variable based barrier synchronization.
type CondBarrier struct {
	threadCount int
	count       int
	generation  int
	mu          sync.Mutex
	cond        *sync.Cond
}

func newCondBarrier(threadCount int) (*CondBarrier, error) {
	if threadCount <= 0 {
		return nil, errors.New("threadCount must be positive")
	}
	b := &CondBarrier{threadCount: threadCount}
	b.cond = sync.NewCond(&b.mu)
	return b, nil
}

func (b *CondBarrier) wait() {
	b.mu.Lock()
	g := b.generation
	b.count++
	if b.count == b.threadCount {
		b.count = 0
		b.generation++
		b.cond.Broadcast()
		b.mu.Unlock()
		return
	}
	for g == b.generation {
		b.cond.Wait()
	}
	b.mu.Unlock()
}

// RunBusyBarrier executes barrierCount busy-wait barriers across threadCount goroutines.
func RunBusyBarrier(threadCount, barrierCount int) (time.Duration, error) {
	barrier, err := newBusyBarrier(threadCount, barrierCount)
	if err != nil {
		return 0, err
	}
	var wg sync.WaitGroup
	start := time.Now()
	for rank := 0; rank < threadCount; rank++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < barrierCount; i++ {
				barrier.wait(i)
			}
		}()
	}
	wg.Wait()
	return time.Since(start), nil
}

// RunCondBarrier executes barrierCount cond-var barriers across threadCount goroutines.
func RunCondBarrier(threadCount, barrierCount int) (time.Duration, error) {
	if barrierCount <= 0 {
		return 0, errors.New("barrierCount must be positive")
	}
	barrier, err := newCondBarrier(threadCount)
	if err != nil {
		return 0, err
	}
	var wg sync.WaitGroup
	start := time.Now()
	for rank := 0; rank < threadCount; rank++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < barrierCount; i++ {
				barrier.wait()
			}
		}()
	}
	wg.Wait()
	return time.Since(start), nil
}
