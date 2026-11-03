package evidence

import "testing"

func TestEnergyMinimumIsNotEntailment(t *testing.T) {
	models, modelsWithAFalse := 0, 0
	for _, a := range []bool{false, true} {
		for _, b := range []bool{false, true} {
			if a || b {
				models++
				if !a {
					modelsWithAFalse++
				}
			}
		}
	}
	if models != 3 || modelsWithAFalse != 1 {
		t.Fatal("incorrect counterexample for K=(A OR B)")
	}
}

func TestNonzeroCohomologyCanHaveConsistentGlobalSection(t *testing.T) {
	// Incidence of a triangle graph with no 2-cell, for the constant real sheaf.
	delta := [3][3]int{{-1, 1, 0}, {0, -1, 1}, {1, 0, -1}}
	for _, row := range delta {
		if row[0]+row[1]+row[2] != 0 {
			t.Fatal("the constant global section is not in the kernel")
		}
	}
	// First two rows are independent; the sum of all rows is zero: rank=2.
	minor := delta[0][0]*delta[1][1] - delta[0][1]*delta[1][0]
	if minor == 0 {
		t.Fatal("incidence rank is not two")
	}
	for column := 0; column < 3; column++ {
		if delta[0][column]+delta[1][column]+delta[2][column] != 0 {
			t.Fatal("missing row dependency")
		}
	}
	const rank, vertices, edges = 2, 3, 3
	if vertices-rank != 1 || edges-rank != 1 {
		t.Fatal("incorrect H0 or H1 dimension")
	}
}

func TestHarmonicPotentialIsNotConjunctiveEvidence(t *testing.T) {
	p, snapshot := hornFixture(t)
	snapshot.Facts[1] = false
	potential := (1.0 + 0.0) / 2
	result, err := p.Search(snapshot, 2, "heat", 100000)
	if err != nil || potential != .5 || result.Proved {
		t.Fatal("positive harmonic potential was confused with a Horn proof", err)
	}
}

func TestStudyUpdatesExcludeUnusedPrimaryFacts(t *testing.T) {
	for _, family := range HornFamilies {
		c, err := GenerateHornCase(family, 256, 1)
		if err != nil {
			t.Fatal(err)
		}
		for _, root := range c.Roots {
			if len(c.Program.consumers[root]) == 0 {
				t.Fatal("an unused root contaminated the update benchmark")
			}
		}
	}
}
