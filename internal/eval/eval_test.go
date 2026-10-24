package eval

import "testing"

func TestRecallAtK(t *testing.T) {
	ret := []string{"A", "B", "C", "D"}
	if got := RecallAtK(ret, []string{"B", "D"}, 2); got != 0.5 {
		t.Fatalf("got %v", got)
	}
	if got := RecallAtK(ret, []string{"B", "D"}, 4); got != 1 {
		t.Fatalf("got %v", got)
	}
	if RecallAtK(ret, nil, 3) != 0 {
		t.Fatal("no relevant docs")
	}
}

func TestReciprocalRank(t *testing.T) {
	if got := ReciprocalRank([]string{"X", "Y", "A"}, []string{"A", "Y"}); got != 0.5 {
		t.Fatalf("got %v", got)
	}
	if ReciprocalRank([]string{"X"}, []string{"A"}) != 0 {
		t.Fatal("miss must be 0")
	}
}
