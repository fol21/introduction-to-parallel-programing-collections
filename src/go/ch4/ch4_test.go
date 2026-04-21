package ch4

import (
	"fmt"
	"math"
	"reflect"
	"testing"
)

func TestHelloMessages(t *testing.T) {
	lines, err := HelloMessages(4)
	if err != nil {
		t.Fatalf("HelloMessages returned error: %v", err)
	}
	if len(lines) != 5 {
		t.Fatalf("unexpected line count: got %d want %d", len(lines), 5)
	}
	if lines[0] != "Hello from the main thread" {
		t.Fatalf("unexpected main line: %q", lines[0])
	}

	wantThreadLines := []string{
		fmt.Sprintf("Hello from thread %d of %d", 0, 4),
		fmt.Sprintf("Hello from thread %d of %d", 1, 4),
		fmt.Sprintf("Hello from thread %d of %d", 2, 4),
		fmt.Sprintf("Hello from thread %d of %d", 3, 4),
	}
	if !reflect.DeepEqual(lines[1:], wantThreadLines) {
		t.Fatalf("unexpected thread lines: got %v want %v", lines[1:], wantThreadLines)
	}
}

func TestPiEstimates(t *testing.T) {
	const n = int64(1_000_000)
	mutexPi, err := PiEstimateMutex(4, n)
	if err != nil {
		t.Fatalf("PiEstimateMutex returned error: %v", err)
	}
	chanPi, err := PiEstimateChannel(4, n)
	if err != nil {
		t.Fatalf("PiEstimateChannel returned error: %v", err)
	}
	serialPi, err := SerialPi(n)
	if err != nil {
		t.Fatalf("SerialPi returned error: %v", err)
	}
	ref := PiReference()

	if math.Abs(mutexPi-ref) > 2e-6 {
		t.Fatalf("mutex estimate too far from reference: got %.10f ref %.10f", mutexPi, ref)
	}
	if math.Abs(chanPi-ref) > 2e-6 {
		t.Fatalf("channel estimate too far from reference: got %.10f ref %.10f", chanPi, ref)
	}
	if math.Abs(serialPi-ref) > 2e-6 {
		t.Fatalf("serial estimate too far from reference: got %.10f ref %.10f", serialPi, ref)
	}
}

func TestBarriers(t *testing.T) {
	busyElapsed, err := RunBusyBarrier(4, 10)
	if err != nil {
		t.Fatalf("RunBusyBarrier returned error: %v", err)
	}
	condElapsed, err := RunCondBarrier(4, 10)
	if err != nil {
		t.Fatalf("RunCondBarrier returned error: %v", err)
	}
	if busyElapsed < 0 || condElapsed < 0 {
		t.Fatalf("elapsed durations must be non-negative")
	}
}
