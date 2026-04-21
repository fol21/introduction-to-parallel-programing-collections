package ch3

import (
	"math"
	"reflect"
	"sort"
	"testing"
)

func TestVectorSumParallel(t *testing.T) {
	x := []float64{1, 2, 3, 4}
	y := []float64{10, 20, 30, 40}

	z, err := VectorSumParallel(x, y, 3)
	if err != nil {
		t.Fatalf("VectorSumParallel returned error: %v", err)
	}

	want := []float64{11, 22, 33, 44}
	if !reflect.DeepEqual(z, want) {
		t.Fatalf("unexpected vector sum: got %v, want %v", z, want)
	}
}

func TestMatVectMultParallel(t *testing.T) {
	A := []float64{
		1, 2, 3,
		4, 5, 6,
	}
	x := []float64{1, 1, 1}

	y, err := MatVectMultParallel(A, 2, 3, x, 2)
	if err != nil {
		t.Fatalf("MatVectMultParallel returned error: %v", err)
	}

	want := []float64{6, 15}
	if !reflect.DeepEqual(y, want) {
		t.Fatalf("unexpected product: got %v, want %v", y, want)
	}
}

func TestTrapParallel(t *testing.T) {
	got, err := TrapParallel(0, 1, 10000, 4)
	if err != nil {
		t.Fatalf("TrapParallel returned error: %v", err)
	}

	want := 1.0 / 3.0
	if math.Abs(got-want) > 1e-5 {
		t.Fatalf("unexpected integral estimate: got %.8f, want %.8f", got, want)
	}
}

func TestOddEvenSortParallel(t *testing.T) {
	in := []int{9, 3, 6, 1, 8, 2, 7, 4, 5}
	want := append([]int(nil), in...)
	sort.Ints(want)

	err := OddEvenSortParallel(in, 3)
	if err != nil {
		t.Fatalf("OddEvenSortParallel returned error: %v", err)
	}

	if !reflect.DeepEqual(in, want) {
		t.Fatalf("unexpected sorted output: got %v, want %v", in, want)
	}
}

func TestGenerateListDeterministic(t *testing.T) {
	a, err := GenerateList(8, 100, 0)
	if err != nil {
		t.Fatalf("GenerateList returned error: %v", err)
	}
	b, err := GenerateList(8, 100, 0)
	if err != nil {
		t.Fatalf("GenerateList returned error: %v", err)
	}

	if !reflect.DeepEqual(a, b) {
		t.Fatalf("expected deterministic generation, got %v and %v", a, b)
	}
}
