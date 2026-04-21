package ch4

import (
	"errors"
	"math/rand"
	"sort"
	"sync"
)

// SortedList stores unique integer keys in ascending order.
type SortedList struct {
	data []int
}

func (l *SortedList) Insert(value int) bool {
	idx := sort.SearchInts(l.data, value)
	if idx < len(l.data) && l.data[idx] == value {
		return false
	}
	l.data = append(l.data, 0)
	copy(l.data[idx+1:], l.data[idx:])
	l.data[idx] = value
	return true
}

func (l *SortedList) Member(value int) bool {
	idx := sort.SearchInts(l.data, value)
	return idx < len(l.data) && l.data[idx] == value
}

func (l *SortedList) Delete(value int) bool {
	idx := sort.SearchInts(l.data, value)
	if idx >= len(l.data) || l.data[idx] != value {
		return false
	}
	copy(l.data[idx:], l.data[idx+1:])
	l.data = l.data[:len(l.data)-1]
	return true
}

func (l *SortedList) Snapshot() []int {
	out := make([]int, len(l.data))
	copy(out, l.data)
	return out
}

// OpCounts tracks operation totals executed by worker goroutines.
type OpCounts struct {
	Member int
	Insert int
	Delete int
}

func (c *OpCounts) add(o OpCounts) {
	c.Member += o.Member
	c.Insert += o.Insert
	c.Delete += o.Delete
}

// RunListWorkloadMutex simulates pth_ll_one_mut style synchronization.
func RunListWorkloadMutex(threadCount, totalOps, initialKeys, maxKey int, searchPercent, insertPercent float64, seed int64) ([]int, OpCounts, error) {
	if threadCount <= 0 {
		return nil, OpCounts{}, errors.New("threadCount must be positive")
	}
	if totalOps < 0 || initialKeys < 0 {
		return nil, OpCounts{}, errors.New("totalOps and initialKeys must be non-negative")
	}
	if maxKey <= 0 {
		return nil, OpCounts{}, errors.New("maxKey must be positive")
	}
	if searchPercent < 0 || insertPercent < 0 || searchPercent+insertPercent > 1 {
		return nil, OpCounts{}, errors.New("invalid operation percentages")
	}

	list := &SortedList{}
	rng := rand.New(rand.NewSource(seed))
	for i := 0; i < initialKeys; i++ {
		list.Insert(rng.Intn(maxKey))
	}

	var mu sync.Mutex
	var countsMu sync.Mutex
	counts := OpCounts{}
	opsPerThread := totalOps / threadCount
	var wg sync.WaitGroup

	for rank := 0; rank < threadCount; rank++ {
		wg.Add(1)
		go func(r int) {
			defer wg.Done()
			localRng := rand.New(rand.NewSource(seed + int64(r+1)))
			local := OpCounts{}
			for i := 0; i < opsPerThread; i++ {
				which := localRng.Float64()
				key := localRng.Intn(maxKey)
				if which < searchPercent {
					mu.Lock()
					list.Member(key)
					mu.Unlock()
					local.Member++
				} else if which < searchPercent+insertPercent {
					mu.Lock()
					list.Insert(key)
					mu.Unlock()
					local.Insert++
				} else {
					mu.Lock()
					list.Delete(key)
					mu.Unlock()
					local.Delete++
				}
			}
			countsMu.Lock()
			counts.add(local)
			countsMu.Unlock()
		}(rank)
	}

	wg.Wait()
	return list.Snapshot(), counts, nil
}

// RunListWorkloadRWMutex simulates pth_ll_rwl style synchronization.
func RunListWorkloadRWMutex(threadCount, totalOps, initialKeys, maxKey int, searchPercent, insertPercent float64, seed int64) ([]int, OpCounts, error) {
	if threadCount <= 0 {
		return nil, OpCounts{}, errors.New("threadCount must be positive")
	}
	if totalOps < 0 || initialKeys < 0 {
		return nil, OpCounts{}, errors.New("totalOps and initialKeys must be non-negative")
	}
	if maxKey <= 0 {
		return nil, OpCounts{}, errors.New("maxKey must be positive")
	}
	if searchPercent < 0 || insertPercent < 0 || searchPercent+insertPercent > 1 {
		return nil, OpCounts{}, errors.New("invalid operation percentages")
	}

	list := &SortedList{}
	rng := rand.New(rand.NewSource(seed))
	for i := 0; i < initialKeys; i++ {
		list.Insert(rng.Intn(maxKey))
	}

	var rw sync.RWMutex
	var countsMu sync.Mutex
	counts := OpCounts{}
	opsPerThread := totalOps / threadCount
	var wg sync.WaitGroup

	for rank := 0; rank < threadCount; rank++ {
		wg.Add(1)
		go func(r int) {
			defer wg.Done()
			localRng := rand.New(rand.NewSource(seed + int64(r+1)))
			local := OpCounts{}
			for i := 0; i < opsPerThread; i++ {
				which := localRng.Float64()
				key := localRng.Intn(maxKey)
				if which < searchPercent {
					rw.RLock()
					list.Member(key)
					rw.RUnlock()
					local.Member++
				} else if which < searchPercent+insertPercent {
					rw.Lock()
					list.Insert(key)
					rw.Unlock()
					local.Insert++
				} else {
					rw.Lock()
					list.Delete(key)
					rw.Unlock()
					local.Delete++
				}
			}
			countsMu.Lock()
			counts.add(local)
			countsMu.Unlock()
		}(rank)
	}

	wg.Wait()
	return list.Snapshot(), counts, nil
}
