package rag

import "testing"

func TestCitedHits(t *testing.T) {
	hits := []Hit{{DocCode: "K-001"}, {DocCode: "K-003"}, {DocCode: "K-003", Lang: "zh"}, {DocCode: "K-004"}}
	got := citedHits("Bring flowers [K-003]. Dress modestly [K-003].", hits, false)
	if len(got) != 1 || got[0].DocCode != "K-003" {
		t.Fatalf("got %+v", got)
	}
	if got := citedHits("no refs", hits, false); len(got) != 0 || got == nil {
		t.Fatalf("want empty non-nil, got %+v", got)
	}
	if got := citedHits("quoted", hits, true); len(got) != 1 || got[0].DocCode != "K-001" {
		t.Fatalf("extractive should cite the quoted passage, got %+v", got)
	}
}
