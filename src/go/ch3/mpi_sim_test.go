package ch3

import (
	"fmt"
	"math"
	"reflect"
	"sort"
	"testing"
)

func TestMPIHelloMessages(t *testing.T) {
	lines, err := MPIHelloMessages(4)
	if err != nil {
		t.Fatalf("MPIHelloMessages returned error: %v", err)
	}
	if len(lines) != 4 {
		t.Fatalf("expected 4 lines, got %d", len(lines))
	}
	if lines[0] != "Greetings from process 0 of 4!" {
		t.Fatalf("unexpected root line: %q", lines[0])
	}

	got := append([]string(nil), lines[1:]...)
	sort.Strings(got)
	want := []string{
		"Greetings from process 1 of 4!",
		"Greetings from process 2 of 4!",
		"Greetings from process 3 of 4!",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected message set: got %v, want %v", got, want)
	}
}

func TestMPIOutputMessages(t *testing.T) {
	lines, err := MPIOutputMessages(3)
	if err != nil {
		t.Fatalf("MPIOutputMessages returned error: %v", err)
	}
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}

	got := append([]string(nil), lines...)
	sort.Strings(got)
	want := []string{
		"Proc 0 of 3 > Does anyone have a toothpick?",
		"Proc 1 of 3 > Does anyone have a toothpick?",
		"Proc 2 of 3 > Does anyone have a toothpick?",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected output set: got %v, want %v", got, want)
	}
}

func TestMPIManyMsgs(t *testing.T) {
	res, err := MPIManyMsgs(16)
	if err != nil {
		t.Fatalf("MPIManyMsgs returned error: %v", err)
	}
	if len(res.Received) != 16 {
		t.Fatalf("expected 16 received values, got %d", len(res.Received))
	}
	for i, v := range res.Received {
		if v != float64(i) {
			t.Fatalf("unexpected payload at %d: got %v, want %v", i, v, float64(i))
		}
	}
}

func TestMPIVectorAddBlock(t *testing.T) {
	x := []float64{1, 2, 3, 4, 5, 6}
	y := []float64{6, 5, 4, 3, 2, 1}
	z, err := MPIVectorAddBlock(x, y, 3)
	if err != nil {
		t.Fatalf("MPIVectorAddBlock returned error: %v", err)
	}
	want := []float64{7, 7, 7, 7, 7, 7}
	if !reflect.DeepEqual(z, want) {
		t.Fatalf("unexpected sum: got %v, want %v", z, want)
	}
}

func TestMPIMatVectMultBlockRows(t *testing.T) {
	A := []float64{
		1, 2,
		3, 4,
		5, 6,
		7, 8,
	}
	x := []float64{1, 1}
	y, err := MPIMatVectMultBlockRows(A, 4, 2, x, 2)
	if err != nil {
		t.Fatalf("MPIMatVectMultBlockRows returned error: %v", err)
	}
	want := []float64{3, 7, 11, 15}
	if !reflect.DeepEqual(y, want) {
		t.Fatalf("unexpected y: got %v, want %v", y, want)
	}
}

func TestMPITrapVariants(t *testing.T) {
	manual, err := MPITrapManualSendRecv(0, 3, 1024, 4)
	if err != nil {
		t.Fatalf("MPITrapManualSendRecv returned error: %v", err)
	}
	reduce, err := MPITrapReduce(0, 3, 1024, 4)
	if err != nil {
		t.Fatalf("MPITrapReduce returned error: %v", err)
	}

	want := 9.0
	if math.Abs(manual-want) > 1e-4 {
		t.Fatalf("manual trap too far: got %.8f want %.8f", manual, want)
	}
	if math.Abs(reduce-want) > 1e-4 {
		t.Fatalf("reduce trap too far: got %.8f want %.8f", reduce, want)
	}
	if math.Abs(manual-reduce) > 1e-9 {
		t.Fatalf("manual and reduce results differ: %v vs %v", manual, reduce)
	}
}

func TestMPIValidationErrors(t *testing.T) {
	_, err := MPIVectorAddBlock([]float64{1, 2}, []float64{1}, 2)
	if err == nil {
		t.Fatalf("expected vector length mismatch error")
	}

	_, err = MPIMatVectMultBlockRows([]float64{1, 2, 3, 4}, 2, 2, []float64{1, 2}, 3)
	if err == nil {
		t.Fatalf("expected divisibility error")
	}

	_, err = MPITrapManualSendRecv(0, 1, 10, 3)
	if err == nil {
		t.Fatalf("expected trap divisibility error")
	}

	lines, err := MPIHelloMessages(2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lines) != 2 {
		t.Fatalf("unexpected line count: got %d", len(lines))
	}
	if lines[0] != "Greetings from process 0 of 2!" {
		t.Fatalf("unexpected line 0: %s", lines[0])
	}
	if lines[1] != fmt.Sprintf("Greetings from process %d of %d!", 1, 2) {
		t.Fatalf("unexpected line 1: %s", lines[1])
	}
}
