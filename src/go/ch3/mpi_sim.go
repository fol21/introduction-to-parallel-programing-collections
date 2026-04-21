package ch3

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// ManyMsgsResult stores elapsed times for two communication strategies.
type ManyMsgsResult struct {
	SingleMessagePerSend time.Duration
	BulkMessage          time.Duration
	Received             []float64
}

// MPIHelloMessages simulates mpi_hello.c using rank goroutines and channels.
func MPIHelloMessages(commSz int) ([]string, error) {
	if commSz <= 0 {
		return nil, errors.New("commSz must be positive")
	}

	rootInbox := make(chan string, commSz-1)
	var wg sync.WaitGroup

	for rank := 1; rank < commSz; rank++ {
		wg.Add(1)
		go func(r int) {
			defer wg.Done()
			rootInbox <- fmt.Sprintf("Greetings from process %d of %d!", r, commSz)
		}(rank)
	}

	lines := make([]string, 0, commSz)
	lines = append(lines, fmt.Sprintf("Greetings from process %d of %d!", 0, commSz))

	pending := make(map[string]struct{}, commSz-1)
	for rank := 1; rank < commSz; rank++ {
		pending[fmt.Sprintf("Greetings from process %d of %d!", rank, commSz)] = struct{}{}
	}

	for len(pending) > 0 {
		msg := <-rootInbox
		if _, ok := pending[msg]; ok {
			lines = append(lines, msg)
			delete(pending, msg)
		}
	}

	wg.Wait()
	close(rootInbox)
	return lines, nil
}

// MPIOutputMessages simulates each rank independently printing one output line.
func MPIOutputMessages(commSz int) ([]string, error) {
	if commSz <= 0 {
		return nil, errors.New("commSz must be positive")
	}

	out := make(chan string, commSz)
	var wg sync.WaitGroup
	for rank := 0; rank < commSz; rank++ {
		wg.Add(1)
		go func(r int) {
			defer wg.Done()
			out <- fmt.Sprintf("Proc %d of %d > Does anyone have a toothpick?", r, commSz)
		}(rank)
	}

	wg.Wait()
	close(out)

	lines := make([]string, 0, commSz)
	for line := range out {
		lines = append(lines, line)
	}
	return lines, nil
}

// MPIManyMsgs simulates sending n doubles from rank 0 to rank 1.
func MPIManyMsgs(n int) (*ManyMsgsResult, error) {
	if n <= 0 {
		return nil, errors.New("n must be positive")
	}

	sendSingle := make(chan float64)
	recvSingleDone := make(chan []float64, 1)

	startSingle := time.Now()
	go func() {
		recv := make([]float64, 0, n)
		for i := 0; i < n; i++ {
			recv = append(recv, <-sendSingle)
		}
		recvSingleDone <- recv
	}()
	for i := 0; i < n; i++ {
		sendSingle <- float64(i)
	}
	receivedSingle := <-recvSingleDone
	singleElapsed := time.Since(startSingle)
	close(sendSingle)

	sendBulk := make(chan []float64, 1)
	recvBulkDone := make(chan []float64, 1)
	payload := make([]float64, n)
	for i := range payload {
		payload[i] = float64(i)
	}

	startBulk := time.Now()
	go func() {
		data := <-sendBulk
		copyOut := make([]float64, len(data))
		copy(copyOut, data)
		recvBulkDone <- copyOut
	}()
	sendBulk <- payload
	receivedBulk := <-recvBulkDone
	bulkElapsed := time.Since(startBulk)
	close(sendBulk)

	if len(receivedBulk) != len(receivedSingle) {
		return nil, errors.New("communication mismatch between single and bulk")
	}

	return &ManyMsgsResult{
		SingleMessagePerSend: singleElapsed,
		BulkMessage:          bulkElapsed,
		Received:             receivedBulk,
	}, nil
}

// MPIVectorAddBlock simulates block-distributed vector add with scatter/gather.
func MPIVectorAddBlock(x, y []float64, commSz int) ([]float64, error) {
	if commSz <= 0 {
		return nil, errors.New("commSz must be positive")
	}
	if len(x) != len(y) {
		return nil, errors.New("x and y must have the same length")
	}
	if len(x)%commSz != 0 {
		return nil, errors.New("vector length must be divisible by commSz")
	}

	n := len(x)
	localN := n / commSz
	result := make([]float64, n)

	type block struct {
		rank int
		data []float64
	}
	gather := make(chan block, commSz)
	var wg sync.WaitGroup

	for rank := 0; rank < commSz; rank++ {
		start := rank * localN
		end := start + localN
		localX := append([]float64(nil), x[start:end]...)
		localY := append([]float64(nil), y[start:end]...)

		wg.Add(1)
		go func(r int, lx, ly []float64) {
			defer wg.Done()
			localZ := make([]float64, len(lx))
			for i := range lx {
				localZ[i] = lx[i] + ly[i]
			}
			gather <- block{rank: r, data: localZ}
		}(rank, localX, localY)
	}

	wg.Wait()
	close(gather)

	for part := range gather {
		start := part.rank * localN
		copy(result[start:start+localN], part.data)
	}
	return result, nil
}

// MPIMatVectMultBlockRows simulates block-row matrix distribution with a broadcast x.
func MPIMatVectMultBlockRows(A []float64, m, n int, x []float64, commSz int) ([]float64, error) {
	if commSz <= 0 {
		return nil, errors.New("commSz must be positive")
	}
	if m <= 0 || n <= 0 {
		return nil, errors.New("m and n must be positive")
	}
	if len(A) != m*n {
		return nil, errors.New("matrix size mismatch")
	}
	if len(x) != n {
		return nil, errors.New("x size mismatch")
	}
	if m%commSz != 0 || n%commSz != 0 {
		return nil, errors.New("m and n must be divisible by commSz")
	}

	localM := m / commSz
	y := make([]float64, m)
	type rows struct {
		rank int
		vec  []float64
	}
	gather := make(chan rows, commSz)
	var wg sync.WaitGroup

	for rank := 0; rank < commSz; rank++ {
		rowStart := rank * localM
		rowEnd := rowStart + localM
		localA := make([]float64, localM*n)
		copy(localA, A[rowStart*n:rowEnd*n])
		xCopy := append([]float64(nil), x...)

		wg.Add(1)
		go func(r int, lA, bx []float64) {
			defer wg.Done()
			localY := make([]float64, localM)
			for i := 0; i < localM; i++ {
				offset := i * n
				sum := 0.0
				for j := 0; j < n; j++ {
					sum += lA[offset+j] * bx[j]
				}
				localY[i] = sum
			}
			gather <- rows{rank: r, vec: localY}
		}(rank, localA, xCopy)
	}

	wg.Wait()
	close(gather)

	for part := range gather {
		start := part.rank * localM
		copy(y[start:start+localM], part.vec)
	}

	return y, nil
}

// MPITrapManualSendRecv simulates mpi_trap2 style explicit send/recv accumulation.
func MPITrapManualSendRecv(a, b float64, n, commSz int) (float64, error) {
	if commSz <= 0 {
		return 0, errors.New("commSz must be positive")
	}
	if n <= 0 || n%commSz != 0 {
		return 0, errors.New("n must be positive and divisible by commSz")
	}

	h := (b - a) / float64(n)
	localN := n / commSz
	toRoot := make(chan float64, commSz-1)
	var wg sync.WaitGroup

	trap := func(left, right float64, trapCount int, baseLen float64) float64 {
		estimate := (left*left + right*right) / 2.0
		for i := 1; i <= trapCount-1; i++ {
			x := left + float64(i)*baseLen
			estimate += x * x
		}
		return estimate * baseLen
	}

	rootLocalA := a
	rootLocalB := rootLocalA + float64(localN)*h
	total := trap(rootLocalA, rootLocalB, localN, h)

	for rank := 1; rank < commSz; rank++ {
		wg.Add(1)
		go func(r int) {
			defer wg.Done()
			localA := a + float64(r*localN)*h
			localB := localA + float64(localN)*h
			toRoot <- trap(localA, localB, localN, h)
		}(rank)
	}

	wg.Wait()
	close(toRoot)
	for v := range toRoot {
		total += v
	}
	return total, nil
}

// MPITrapReduce simulates mpi_trap3 style reduction.
func MPITrapReduce(a, b float64, n, commSz int) (float64, error) {
	if commSz <= 0 {
		return 0, errors.New("commSz must be positive")
	}
	if n <= 0 || n%commSz != 0 {
		return 0, errors.New("n must be positive and divisible by commSz")
	}

	h := (b - a) / float64(n)
	localN := n / commSz
	partials := make(chan float64, commSz)
	var wg sync.WaitGroup

	trap := func(left, right float64, trapCount int, baseLen float64) float64 {
		estimate := (left*left + right*right) / 2.0
		for i := 1; i <= trapCount-1; i++ {
			x := left + float64(i)*baseLen
			estimate += x * x
		}
		return estimate * baseLen
	}

	for rank := 0; rank < commSz; rank++ {
		wg.Add(1)
		go func(r int) {
			defer wg.Done()
			localA := a + float64(r*localN)*h
			localB := localA + float64(localN)*h
			partials <- trap(localA, localB, localN, h)
		}(rank)
	}

	wg.Wait()
	close(partials)

	total := 0.0
	for p := range partials {
		total += p
	}
	return total, nil
}
