package ch4

import "testing"

func isSortedUnique(vals []int) bool {
	for i := 1; i < len(vals); i++ {
		if vals[i-1] >= vals[i] {
			return false
		}
	}
	return true
}

func TestSortedListBasicOps(t *testing.T) {
	l := &SortedList{}
	if !l.Insert(5) || !l.Insert(3) || !l.Insert(9) {
		t.Fatalf("expected inserts to succeed")
	}
	if l.Insert(5) {
		t.Fatalf("duplicate insert should fail")
	}
	if !l.Member(3) || l.Member(4) {
		t.Fatalf("member checks returned unexpected results")
	}
	if !l.Delete(5) || l.Delete(5) {
		t.Fatalf("delete checks returned unexpected results")
	}
	got := l.Snapshot()
	if len(got) != 2 || got[0] != 3 || got[1] != 9 {
		t.Fatalf("unexpected list snapshot: %v", got)
	}
}

func TestRunListWorkloadMutex(t *testing.T) {
	vals, counts, err := RunListWorkloadMutex(4, 2000, 200, 10000, 0.8, 0.1, 7)
	if err != nil {
		t.Fatalf("RunListWorkloadMutex returned error: %v", err)
	}
	if !isSortedUnique(vals) {
		t.Fatalf("final list should be sorted and unique")
	}
	if counts.Member+counts.Insert+counts.Delete != 2000 {
		t.Fatalf("unexpected op totals: %+v", counts)
	}
}

func TestRunListWorkloadRWMutex(t *testing.T) {
	vals, counts, err := RunListWorkloadRWMutex(4, 2400, 200, 10000, 0.8, 0.1, 11)
	if err != nil {
		t.Fatalf("RunListWorkloadRWMutex returned error: %v", err)
	}
	if !isSortedUnique(vals) {
		t.Fatalf("final list should be sorted and unique")
	}
	if counts.Member+counts.Insert+counts.Delete != 2400 {
		t.Fatalf("unexpected op totals: %+v", counts)
	}
}

func TestFineListBasicOps(t *testing.T) {
	l := NewFineList()
	if !l.Insert(7) || !l.Insert(1) || !l.Insert(5) {
		t.Fatalf("expected inserts to succeed")
	}
	if l.Insert(5) {
		t.Fatalf("duplicate insert should fail")
	}
	if !l.Member(1) || l.Member(3) {
		t.Fatalf("member checks returned unexpected results")
	}
	if !l.Delete(5) || l.Delete(5) {
		t.Fatalf("delete checks returned unexpected results")
	}
	got := l.Snapshot()
	if len(got) != 2 || got[0] != 1 || got[1] != 7 {
		t.Fatalf("unexpected list snapshot: %v", got)
	}
}

func TestRunListWorkloadFineMutex(t *testing.T) {
	vals, counts, err := RunListWorkloadFineMutex(4, 2800, 200, 10000, 0.8, 0.1, 13)
	if err != nil {
		t.Fatalf("RunListWorkloadFineMutex returned error: %v", err)
	}
	if !isSortedUnique(vals) {
		t.Fatalf("final list should be sorted and unique")
	}
	if counts.Member+counts.Insert+counts.Delete != 2800 {
		t.Fatalf("unexpected op totals: %+v", counts)
	}
}
