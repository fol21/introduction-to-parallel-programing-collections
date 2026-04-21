package ch4

import (
	"errors"
	"math"
	"math/rand"
	"sync"
)

type fineNode struct {
	value int
	next  *fineNode
	mu    sync.Mutex
}

// FineList is a sorted unique linked list that uses lock-coupling.
type FineList struct {
	head *fineNode
}

func NewFineList() *FineList {
	// Sentinel node is never removed.
	return &FineList{head: &fineNode{value: math.MinInt}}
}

func (l *FineList) Insert(value int) bool {
	pred := l.head
	pred.mu.Lock()
	curr := pred.next
	if curr != nil {
		curr.mu.Lock()
	}

	for curr != nil && curr.value < value {
		pred.mu.Unlock()
		pred = curr
		curr = curr.next
		if curr != nil {
			curr.mu.Lock()
		}
	}

	if curr != nil && curr.value == value {
		curr.mu.Unlock()
		pred.mu.Unlock()
		return false
	}

	node := &fineNode{value: value, next: curr}
	pred.next = node
	if curr != nil {
		curr.mu.Unlock()
	}
	pred.mu.Unlock()
	return true
}

func (l *FineList) Member(value int) bool {
	pred := l.head
	pred.mu.Lock()
	curr := pred.next
	if curr != nil {
		curr.mu.Lock()
	}
	pred.mu.Unlock()

	for curr != nil && curr.value < value {
		next := curr.next
		if next != nil {
			next.mu.Lock()
		}
		curr.mu.Unlock()
		curr = next
	}

	if curr == nil {
		return false
	}
	found := curr.value == value
	curr.mu.Unlock()
	return found
}

func (l *FineList) Delete(value int) bool {
	pred := l.head
	pred.mu.Lock()
	curr := pred.next
	if curr != nil {
		curr.mu.Lock()
	}

	for curr != nil && curr.value < value {
		pred.mu.Unlock()
		pred = curr
		curr = curr.next
		if curr != nil {
			curr.mu.Lock()
		}
	}

	if curr == nil || curr.value != value {
		if curr != nil {
			curr.mu.Unlock()
		}
		pred.mu.Unlock()
		return false
	}

	pred.next = curr.next
	curr.mu.Unlock()
	pred.mu.Unlock()
	return true
}

func (l *FineList) Snapshot() []int {
	out := make([]int, 0)
	pred := l.head
	pred.mu.Lock()
	curr := pred.next
	if curr != nil {
		curr.mu.Lock()
	}
	pred.mu.Unlock()

	for curr != nil {
		out = append(out, curr.value)
		next := curr.next
		if next != nil {
			next.mu.Lock()
		}
		curr.mu.Unlock()
		curr = next
	}

	return out
}

// RunListWorkloadFineMutex simulates pth_ll_mult_mut style fine-grained node locking.
func RunListWorkloadFineMutex(threadCount, totalOps, initialKeys, maxKey int, searchPercent, insertPercent float64, seed int64) ([]int, OpCounts, error) {
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

	list := NewFineList()
	rng := rand.New(rand.NewSource(seed))
	for i := 0; i < initialKeys; i++ {
		list.Insert(rng.Intn(maxKey))
	}

	opsPerThread := totalOps / threadCount
	counts := OpCounts{}
	var countsMu sync.Mutex
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
					list.Member(key)
					local.Member++
				} else if which < searchPercent+insertPercent {
					list.Insert(key)
					local.Insert++
				} else {
					list.Delete(key)
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
