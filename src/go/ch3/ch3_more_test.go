package ch3

import (
	"reflect"
	"sort"
	"testing"
)

func TestBubbleSort(t *testing.T) {
	in := []int{5, 1, 4, 2, 8}
	BubbleSort(in)
	want := []int{1, 2, 4, 5, 8}
	if !reflect.DeepEqual(in, want) {
		t.Fatalf("unexpected bubble sort result: got %v want %v", in, want)
	}
}

func TestBubbleSortParallel(t *testing.T) {
	in := []int{10, 3, 9, 1, 8, 2, 7, 4, 6, 5}
	want := append([]int(nil), in...)
	sort.Ints(want)
	if err := BubbleSortParallel(in, 4); err != nil {
		t.Fatalf("BubbleSortParallel returned error: %v", err)
	}
	if !reflect.DeepEqual(in, want) {
		t.Fatalf("unexpected BubbleSortParallel result: got %v want %v", in, want)
	}
}

func TestMPIBlockOddEvenSort(t *testing.T) {
	in := []int{9, 8, 7, 6, 5, 4, 3, 2}
	got, err := MPIBlockOddEvenSort(in, 4)
	if err != nil {
		t.Fatalf("MPIBlockOddEvenSort returned error: %v", err)
	}
	want := []int{2, 3, 4, 5, 6, 7, 8, 9}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected distributed odd-even sort result: got %v want %v", got, want)
	}
}

func TestMPIMatVectTimeSim(t *testing.T) {
	res, err := MPIMatVectTimeSim(4, 4, 2, 42)
	if err != nil {
		t.Fatalf("MPIMatVectTimeSim returned error: %v", err)
	}
	if len(res.Y) != 4 {
		t.Fatalf("unexpected y length: got %d", len(res.Y))
	}
	if res.Elapsed < 0 {
		t.Fatalf("elapsed should be non-negative")
	}
}
