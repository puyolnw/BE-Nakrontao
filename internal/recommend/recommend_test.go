package recommend

import (
	"math"
	"testing"
)

var places = []Place{
	{ID: 1, Category: "sacred", Lat: 16.2, Lon: 103.3, VisitMinutes: 45, Tags: []string{"shrine"}},
	{ID: 2, Category: "history", Lat: 16.201, Lon: 103.301, VisitMinutes: 30, Tags: []string{"archaeology"}},
	{ID: 3, Category: "food", Lat: 16.202, Lon: 103.302, VisitMinutes: 40},
	{ID: 4, Category: "history", Lat: 16.5, Lon: 103.9, VisitMinutes: 60}, // far away
}

func TestHaversine(t *testing.T) {
	d := Haversine(13.7563, 100.5018, 18.7883, 98.9853) // Bangkok–Chiang Mai ≈ 580 km
	if math.Abs(d-580) > 15 {
		t.Fatalf("got %.1f km", d)
	}
}

func TestBuildRespectsBudgetAndInterests(t *testing.T) {
	p := Build(places, Request{Interests: []string{"history"}, Minutes: 120, Mode: "walk"})
	if p.TotalMinutes > 120 {
		t.Fatalf("over budget: %d", p.TotalMinutes)
	}
	ids := map[int64]bool{}
	for _, s := range p.Stops {
		ids[s.Place.ID] = true
		if len(s.Reasons) == 0 {
			t.Fatal("every stop needs a reason")
		}
	}
	if ids[3] {
		t.Fatal("food was not requested")
	}
	if !ids[1] || !ids[2] {
		t.Fatalf("expected shrine and nearby history site, got %+v", ids)
	}
	if ids[4] {
		t.Fatal("far site cannot fit in a 2h walking plan")
	}
}

func TestBuildNoMatches(t *testing.T) {
	p := Build([]Place{{ID: 9, Category: "stay", VisitMinutes: 10}}, Request{Interests: []string{"food"}, Minutes: 60})
	if len(p.Stops) != 0 || p.Stops == nil {
		t.Fatalf("want empty non-nil stops, got %+v", p.Stops)
	}
}

func TestDriveIsFasterThanWalk(t *testing.T) {
	if travelMinutes(5, "drive") >= travelMinutes(5, "walk") {
		t.Fatal("drive should be faster")
	}
}
