package evidence

import (
	"math"
	"slices"
	"testing"
)

func TestAdaptiveModesAgreeOnFreshGraphs(t *testing.T) {
	for _, family := range AdaptiveFamilies {
		c, err := GenerateAdaptiveCase(family, 256, 2001)
		if err != nil {
			t.Fatal(err)
		}
		events, err := PlanAdaptiveUpdates(c, 3000)
		if err != nil {
			t.Fatal(err)
		}
		full, err := NewHornState(c.Program, c.Snapshot)
		if err != nil {
			t.Fatal(err)
		}
		partial, _ := NewHornState(c.Program, c.Snapshot)
		adaptive, _ := NewHornState(c.Program, c.Snapshot)
		for step, event := range events {
			if _, err := full.UpdateFull(event.Update); err != nil {
				t.Fatal(err)
			}
			if _, err := partial.Update(event.Update); err != nil {
				t.Fatal(err)
			}
			stats, err := adaptive.UpdateAdaptive(event.Update, .05)
			if err != nil || !slices.Equal(full.Known(), partial.Known()) || !slices.Equal(full.Known(), adaptive.Known()) {
				t.Fatalf("%s step=%d: %+v %v", family, step, stats, err)
			}
			limit := max(1, int(.05*float64(c.Program.atoms+c.Program.edges)))
			if stats.DiscoveryVisits > limit {
				t.Fatal("adaptive discovery exceeded its explicit limit")
			}
		}
	}
}

func TestAdaptiveRejectsInvalidInputsAtomically(t *testing.T) {
	p, snapshot := hornFixture(t)
	state, err := NewHornState(p, snapshot)
	if err != nil {
		t.Fatal(err)
	}
	for _, fraction := range []float64{0, -1, 2, math.NaN(), math.Inf(1)} {
		if _, err := state.UpdateAdaptive(HornUpdate{Facts: []HornChange{{0, false}}}, fraction); err == nil {
			t.Fatal("invalid discovery limit accepted")
		}
	}
	if !slices.Equal(snapshot.Facts, state.Snapshot().Facts) || snapshot.Revision != state.Snapshot().Revision {
		t.Fatal("invalid policy modified the state")
	}
	stats, err := state.UpdateAdaptive(HornUpdate{Facts: []HornChange{{0, false}}}, .001)
	if err != nil || stats.Mode != "adaptive-full" || state.Known()[2] {
		t.Fatal(stats, err)
	}
}
