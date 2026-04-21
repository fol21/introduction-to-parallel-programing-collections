package ch3

import (
	"errors"
	"sort"
	"sync"
)

// MPIBlockOddEvenSort simulates mpi_odd_even.c using block-distributed local arrays.
func MPIBlockOddEvenSort(global []int, commSz int) ([]int, error) {
	if commSz <= 0 {
		return nil, errors.New("commSz must be positive")
	}
	if len(global)%commSz != 0 {
		return nil, errors.New("global length must be divisible by commSz")
	}
	localN := len(global) / commSz
	locals := make([][]int, commSz)
	for rank := 0; rank < commSz; rank++ {
		start := rank * localN
		locals[rank] = append([]int(nil), global[start:start+localN]...)
		sort.Ints(locals[rank])
	}

	for phase := 0; phase < commSz; phase++ {
		startRank := 0
		if phase%2 == 1 {
			startRank = 1
		}

		type pairResult struct {
			leftRank  int
			rightRank int
			leftData  []int
			rightData []int
		}
		out := make(chan pairResult, commSz/2+1)
		var wg sync.WaitGroup

		for left := startRank; left+1 < commSz; left += 2 {
			right := left + 1
			leftCopy := append([]int(nil), locals[left]...)
			rightCopy := append([]int(nil), locals[right]...)

			wg.Add(1)
			go func(lRank, rRank int, lBlock, rBlock []int) {
				defer wg.Done()
				merged := make([]int, 0, len(lBlock)+len(rBlock))
				i, j := 0, 0
				for i < len(lBlock) && j < len(rBlock) {
					if lBlock[i] <= rBlock[j] {
						merged = append(merged, lBlock[i])
						i++
					} else {
						merged = append(merged, rBlock[j])
						j++
					}
				}
				merged = append(merged, lBlock[i:]...)
				merged = append(merged, rBlock[j:]...)
				n := len(lBlock)
				out <- pairResult{
					leftRank:  lRank,
					rightRank: rRank,
					leftData:  append([]int(nil), merged[:n]...),
					rightData: append([]int(nil), merged[n:]...),
				}
			}(left, right, leftCopy, rightCopy)
		}

		wg.Wait()
		close(out)
		for part := range out {
			locals[part.leftRank] = part.leftData
			locals[part.rightRank] = part.rightData
		}
	}

	result := make([]int, len(global))
	for rank := 0; rank < commSz; rank++ {
		copy(result[rank*localN:(rank+1)*localN], locals[rank])
	}
	return result, nil
}
