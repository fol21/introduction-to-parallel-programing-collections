package ch4

import (
	"reflect"
	"strings"
	"testing"
)

func TestMatVectThreads(t *testing.T) {
	A := []float64{
		1, 2, 3,
		4, 5, 6,
		7, 8, 9,
		10, 11, 12,
	}
	x := []float64{1, 0, 1}
	got, err := MatVectThreads(A, 4, 3, x, 2)
	if err != nil {
		t.Fatalf("MatVectThreads returned error: %v", err)
	}
	want := []float64{4, 10, 16, 22}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected product: got %v want %v", got, want)
	}
}

func TestTokenizeLinesSafe(t *testing.T) {
	lines := []string{
		"alpha beta gamma",
		"go routines and channels",
		"one\ttwo  three",
	}
	results, err := TokenizeLinesSafe(lines, 3)
	if err != nil {
		t.Fatalf("TokenizeLinesSafe returned error: %v", err)
	}
	if len(results) != len(lines) {
		t.Fatalf("unexpected results length: got %d want %d", len(results), len(lines))
	}

	byIndex := make(map[int]TokenizeResult, len(results))
	for _, r := range results {
		byIndex[r.LineIndex] = r
	}
	for i, line := range lines {
		r, ok := byIndex[i]
		if !ok {
			t.Fatalf("missing result for line index %d", i)
		}
		want := strings.Fields(line)
		if !reflect.DeepEqual(r.Tokens, want) {
			t.Fatalf("unexpected tokens for line %d: got %v want %v", i, r.Tokens, want)
		}
	}
}

func TestTokenizeLinesUnsafeReturnsAllLines(t *testing.T) {
	lines := []string{
		"a b c d",
		"1 2 3 4",
		"x y z",
		"hello world",
	}
	results, err := TokenizeLinesUnsafe(lines, 4)
	if err != nil {
		t.Fatalf("TokenizeLinesUnsafe returned error: %v", err)
	}
	if len(results) != len(lines) {
		t.Fatalf("unexpected results length: got %d want %d", len(results), len(lines))
	}
	for _, r := range results {
		if r.Line == "" {
			t.Fatalf("expected line text in every result")
		}
	}
}
