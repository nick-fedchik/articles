package evidence

import (
	"math/rand"
	"slices"
	"testing"
)

func hornFixture(t *testing.T) (*HornProgram, HornSnapshot) {
	t.Helper()
	p, err := CompileHorn(4, []HornRule{
		{"both", "1", []int{0, 1}, 2},
		{"alternative", "1", []int{3}, 2},
	})
	if err != nil {
		t.Fatal(err)
	}
	return p, HornSnapshot{1, []bool{true, true, false, false}, []bool{true, true}}
}

func TestHornSearchProofsAndMissingConjunct(t *testing.T) {
	p, snapshot := hornFixture(t)
	for _, method := range HornMethods {
		got, err := p.Search(snapshot, 2, method, 100000)
		if err != nil || !got.Proved || p.VerifyProof(snapshot, 2, got.Proof) != nil {
			t.Fatalf("%s: %+v %v", method, got, err)
		}
		missing := cloneHornSnapshot(snapshot)
		missing.Facts[1] = false
		got, err = p.Search(missing, 2, method, 100000)
		if err != nil || got.Proved || got.Exhausted {
			t.Fatalf("missing conjunct %s: %+v %v", method, got, err)
		}
	}
}

func TestHornCertificateRejectsForgeryQueryAndRevision(t *testing.T) {
	p, snapshot := hornFixture(t)
	got, err := p.Search(snapshot, 2, "backward", 100000)
	if err != nil {
		t.Fatal(err)
	}
	proof := got.Proof
	proof.Facts = []int{0}
	if err := p.VerifyProof(snapshot, 2, proof); err == nil {
		t.Fatal("proof without the second premise accepted")
	}
	if err := p.VerifyProof(snapshot, 3, got.Proof); err == nil {
		t.Fatal("proof of another query accepted")
	}
	snapshot.Revision++
	if err := p.VerifyProof(snapshot, 2, got.Proof); err == nil {
		t.Fatal("stale revision accepted")
	}
}

func TestHornBudgetDoesNotClaimFalsehood(t *testing.T) {
	p, snapshot := hornFixture(t)
	for _, method := range HornMethods {
		out, err := p.Search(snapshot, 2, method, 1)
		if err != nil || out.Proved || !out.Exhausted || out.Operations > 1 {
			t.Fatal(out, err)
		}
	}
}

func TestHornWithdrawalAndAlternativeSupport(t *testing.T) {
	p, snapshot := hornFixture(t)
	state, err := NewHornState(p, snapshot)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := state.Update(HornUpdate{Facts: []HornChange{{1, false}}}); err != nil || state.Known()[2] {
		t.Fatal("withdrawing the only proof left a stale conclusion", err)
	}
	if _, err := state.Update(HornUpdate{Facts: []HornChange{{3, true}}}); err != nil || !state.Known()[2] {
		t.Fatal("alternative source did not establish the conclusion", err)
	}
	if _, err := state.Update(HornUpdate{Rules: []HornChange{{1, false}}}); err != nil || state.Known()[2] {
		t.Fatal("disabling alternative rule left a stale conclusion", err)
	}
	before := state.Snapshot()
	if _, err := state.Update(HornUpdate{Facts: []HornChange{{0, false}, {99, true}}}); err == nil ||
		!slices.Equal(before.Facts, state.Snapshot().Facts) || before.Revision != state.Snapshot().Revision {
		t.Fatal("invalid update partially modified the ledger")
	}
}

func TestHornIncrementalAgainstIndependentFixedPoint(t *testing.T) {
	rng := rand.New(rand.NewSource(20261103))
	for trial := 0; trial < 100; trial++ {
		rules := []HornRule{}
		for i := 4; i < 40; i++ {
			a, b := rng.Intn(i), rng.Intn(i)
			premises := []int{a}
			if b != a {
				premises = append(premises, b)
			}
			rules = append(rules, HornRule{ID: string(rune('A' + i)), Version: "1", Premises: premises, Conclusion: i})
		}
		p, err := CompileHorn(40, rules)
		if err != nil {
			t.Fatal(err)
		}
		snapshot := HornSnapshot{Revision: 1, Facts: make([]bool, 40), Enabled: make([]bool, len(rules))}
		for i := range snapshot.Enabled {
			snapshot.Enabled[i] = true
		}
		state, err := NewHornState(p, snapshot)
		if err != nil {
			t.Fatal(err)
		}
		for step := 0; step < 20; step++ {
			if _, err := state.Update(HornUpdate{Facts: []HornChange{{rng.Intn(4), rng.Intn(2) == 0}}}); err != nil {
				t.Fatal(err)
			}
			s := state.Snapshot()
			want := append([]bool(nil), s.Facts...)
			changed := true
			for changed {
				changed = false
				for i, rule := range rules {
					if !s.Enabled[i] || want[rule.Conclusion] {
						continue
					}
					holds := true
					for _, atom := range rule.Premises {
						holds = holds && want[atom]
					}
					if holds {
						want[rule.Conclusion], changed = true, true
					}
				}
			}
			if !slices.Equal(want, state.Known()) {
				t.Fatalf("trial %d step %d: incremental result differs from the fixed point", trial, step)
			}
		}
	}
}

func TestHornRejectsCyclesAndMalformedRules(t *testing.T) {
	for _, rules := range [][]HornRule{
		{{"loop", "1", []int{0}, 0}},
		{{"repeat", "1", []int{0, 0}, 1}},
		{{"empty", "1", nil, 1}},
		{{"range", "1", []int{-1}, 1}},
	} {
		if _, err := CompileHorn(2, rules); err == nil {
			t.Fatal("malformed Horn theory accepted")
		}
	}
}
