package insights

import (
	"math"

	"github.com/reclover/sinakhontao/internal/httpx"
	"testing"
)

func TestMeanSD(t *testing.T) {
	m, sd := MeanSD([]float64{5, 4, 4, 5, 3})
	if math.Abs(m-4.2) > 1e-9 || math.Abs(sd-0.83666) > 1e-4 {
		t.Fatalf("got %.4f %.4f", m, sd)
	}
	if m, sd := MeanSD(nil); m != 0 || sd != 0 {
		t.Fatal("empty input")
	}
	if _, sd := MeanSD([]float64{4}); sd != 0 {
		t.Fatal("single value has no spread")
	}
}

func TestLevel(t *testing.T) {
	for mean, want := range map[float64]string{4.6: "highest", 4.2: "high", 3.0: "moderate", 2.0: "low", 1.2: "lowest"} {
		if got := Level(mean); got != want {
			t.Errorf("%.1f: got %s want %s", mean, got, want)
		}
	}
}

func TestCSVSafe(t *testing.T) {
	if httpx.CSVSafe("=HYPERLINK(1)") != "'=HYPERLINK(1)" || httpx.CSVSafe("ปกติ") != "ปกติ" || httpx.CSVSafe("") != "" {
		t.Fatal("formula injection not neutralised")
	}
}
