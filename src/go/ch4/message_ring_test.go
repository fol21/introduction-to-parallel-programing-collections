package ch4

import (
	"fmt"
	"strings"
	"testing"
)

func TestUnsyncRingMessages(t *testing.T) {
	lines, err := UnsyncRingMessages(8)
	if err != nil {
		t.Fatalf("UnsyncRingMessages returned error: %v", err)
	}
	if len(lines) != 8 {
		t.Fatalf("unexpected line count: got %d", len(lines))
	}

	for rank, line := range lines {
		source := (rank + 8 - 1) % 8
		wantMsg := fmt.Sprintf("Thread %d > Hello to %d from %d", rank, rank, source)
		wantNoMsg := fmt.Sprintf("Thread %d > No message from %d", rank, source)
		if line != wantMsg && line != wantNoMsg {
			t.Fatalf("unexpected line for rank %d: %q", rank, line)
		}
	}
}

func TestSemaphoreRingMessages(t *testing.T) {
	lines, err := SemaphoreRingMessages(8)
	if err != nil {
		t.Fatalf("SemaphoreRingMessages returned error: %v", err)
	}
	if len(lines) != 8 {
		t.Fatalf("unexpected line count: got %d", len(lines))
	}

	for rank, line := range lines {
		source := (rank + 8 - 1) % 8
		want := fmt.Sprintf("Thread %d > Hello to %d from %d", rank, rank, source)
		if line != want {
			t.Fatalf("unexpected line for rank %d: got %q want %q", rank, line, want)
		}
		if strings.Contains(line, "No message") {
			t.Fatalf("did not expect missing message with semaphore sync")
		}
	}
}
