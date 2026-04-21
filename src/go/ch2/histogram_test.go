package ch2

import (
	"reflect"
	"testing"
)

func TestGenBins(t *testing.T) {
	binMaxes, binCounts, err := GenBins(0, 10, 5)
	if err != nil {
		t.Fatalf("GenBins returned error: %v", err)
	}

	expectedMaxes := []float64{2, 4, 6, 8, 10}
	expectedCounts := []int{0, 0, 0, 0, 0}

	if !reflect.DeepEqual(binMaxes, expectedMaxes) {
		t.Fatalf("unexpected bin maxes: got %v, want %v", binMaxes, expectedMaxes)
	}
	if !reflect.DeepEqual(binCounts, expectedCounts) {
		t.Fatalf("unexpected bin counts: got %v, want %v", binCounts, expectedCounts)
	}
}

func TestWhichBinBoundaries(t *testing.T) {
	binMaxes := []float64{2, 4, 6, 8, 10}
	minMeas := 0.0

	tests := []struct {
		value float64
		want  int
	}{
		{0.0, 0},
		{1.999, 0},
		{2.0, 1},
		{3.0, 1},
		{7.999, 3},
		{9.999, 4},
	}

	for _, tc := range tests {
		got, err := WhichBin(tc.value, binMaxes, minMeas)
		if err != nil {
			t.Fatalf("WhichBin(%v) unexpected error: %v", tc.value, err)
		}
		if got != tc.want {
			t.Fatalf("WhichBin(%v) = %d, want %d", tc.value, got, tc.want)
		}
	}

	if _, err := WhichBin(10.0, binMaxes, minMeas); err == nil {
		t.Fatalf("WhichBin(10.0) expected error for upper bound")
	}
}

func TestBuildHistogramParallel(t *testing.T) {
	data := []float64{0.5, 1.2, 2.1, 3.7, 3.9, 8.2, 9.9}

	h, err := BuildHistogramParallel(data, 5, 0, 10, 3)
	if err != nil {
		t.Fatalf("BuildHistogramParallel returned error: %v", err)
	}

	want := []int{2, 3, 0, 0, 2}
	if !reflect.DeepEqual(h.BinCounts, want) {
		t.Fatalf("unexpected counts: got %v, want %v", h.BinCounts, want)
	}
}

func TestBuildHistogramParallelWithInvalidMeasurement(t *testing.T) {
	data := []float64{0.5, 10.0}

	_, err := BuildHistogramParallel(data, 5, 0, 10, 2)
	if err == nil {
		t.Fatalf("expected error for measurement outside valid range")
	}
}
