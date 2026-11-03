package evidence

import (
	"reflect"
	"slices"
	"testing"
)

func normFixture() NormTheory {
	exception := SignedLiteral{Atom: 1}
	return NormTheory{Atoms: 3, Norms: []NormRule{
		{ID: "permit", Version: "1", Action: "run", Effect: Permit, Conditions: []SignedLiteral{{Atom: 0}}, Exception: &exception},
		{ID: "ban", Version: "1", Action: "run", Effect: Prohibit},
	}}
}

func TestExplicitNegationDoesNotExplode(t *testing.T) {
	theory := NormTheory{
		Atoms: 3,
		Evidence: []SignedRule{
			{"positive", "1", []SignedLiteral{{Atom: 0}}, SignedLiteral{Atom: 1}},
			{"negative", "1", []SignedLiteral{{Atom: 0, Negative: true}}, SignedLiteral{Atom: 1, Negative: true}},
		},
		Norms: []NormRule{{ID: "permit", Version: "1", Action: "run", Effect: Permit, Conditions: []SignedLiteral{{Atom: 1}}}},
	}
	p, err := CompileNormative(theory)
	if err != nil {
		t.Fatal(err)
	}
	world, err := p.NewWorld([]State{Both, Unknown, Unknown})
	if err != nil {
		t.Fatal(err)
	}
	certificate, err := p.Evaluate(world, "run")
	if err != nil || !slices.Equal(certificate.Closure, []State{Both, Both, Unknown}) || certificate.Result.Decision != NormConflict {
		t.Fatal(certificate, err)
	}
}

func TestExceptionStatesAndPriorityPreserveEvidence(t *testing.T) {
	for _, guard := range []State{Unknown, False, True, Both} {
		theory := normFixture()
		theory.Priorities = []NormPriority{{"permit", "ban"}}
		p, err := CompileNormative(theory)
		if err != nil {
			t.Fatal(err)
		}
		world, err := p.NewWorld([]State{True, guard, Unknown})
		if err != nil {
			t.Fatal(err)
		}
		certificate, err := p.Evaluate(world, "run")
		if err != nil {
			t.Fatal(err)
		}
		want := map[State]NormDecision{Unknown: NormUnknown, False: NormAllow, True: NormDeny, Both: NormConflict}[guard]
		if certificate.Result.Decision != want || certificate.Result.ConditionEvidence != Both {
			t.Fatalf("exception=%v: %+v, want %s", guard, certificate.Result, want)
		}
	}
}

func TestPriorityCyclesAndUnknownReferencesRejected(t *testing.T) {
	for _, priorities := range [][]NormPriority{
		{{"permit", "ban"}, {"ban", "permit"}},
		{{"missing", "ban"}},
	} {
		theory := normFixture()
		theory.Priorities = priorities
		if _, err := CompileNormative(theory); err == nil {
			t.Fatal("invalid priority relation accepted")
		}
	}
}

func TestNormCertificateRejectsForgeryAndOldContext(t *testing.T) {
	p, err := CompileNormative(normFixture())
	if err != nil {
		t.Fatal(err)
	}
	world, err := p.NewWorld([]State{True, False, Unknown})
	if err != nil {
		t.Fatal(err)
	}
	certificate, err := p.Evaluate(world, "run")
	if err != nil {
		t.Fatal(err)
	}
	forged := cloneNormCertificate(certificate)
	forged.Closure[1] = True
	if err := p.Verify(world, "run", forged); err == nil {
		t.Fatal("hidden contrary evidence accepted")
	}
	forged = cloneNormCertificate(certificate)
	forged.Result.Decision = NormAllow
	if err := p.Verify(world, "run", forged); err == nil {
		t.Fatal("unjustified allowance accepted")
	}
	world.Revision++
	if err := p.Verify(world, "run", certificate); err == nil {
		t.Fatal("stale context accepted")
	}
}

func TestNormCacheChecksOppositionOutsidePositiveProof(t *testing.T) {
	theory := normFixture()
	p, err := CompileNormative(theory)
	if err != nil {
		t.Fatal(err)
	}
	world, err := p.NewWorld([]State{True, False, Unknown})
	if err != nil {
		t.Fatal(err)
	}
	world.NormEnabled[1] = false
	cache := NewNormCache()
	first, err := cache.Resolve(p, world, "run")
	if err != nil || first.Certificate.Result.Decision != NormAllow {
		t.Fatal(first, err)
	}
	first.Certificate.Closure[0] = False
	hit, err := cache.Resolve(p, world, "run")
	if err != nil || !hit.Reused || hit.Certificate.Closure[0] != True {
		t.Fatal("returned mutable data corrupted the cache", hit, err)
	}
	world.NormEnabled[1] = true
	world.Revision++
	result, err := cache.Resolve(p, world, "run")
	if err != nil || !result.Invalidated || result.Reused || result.Certificate.Result.Decision != NormConflict {
		t.Fatal("new prohibition did not invalidate the cached allowance", result, err)
	}
}

func TestNormativeReferenceExhaustiveSmallWorlds(t *testing.T) {
	for raw := 0; raw < 64; raw++ {
		for _, priorities := range [][]NormPriority{nil, {{"permit", "ban"}}, {{"ban", "permit"}}} {
			theory := normFixture()
			theory.Priorities = priorities
			p, err := CompileNormative(theory)
			if err != nil {
				t.Fatal(err)
			}
			world, err := p.NewWorld([]State{State(raw & 3), State((raw >> 2) & 3), State((raw >> 4) & 3)})
			if err != nil {
				t.Fatal(err)
			}
			actual, err := p.Evaluate(world, "run")
			if err != nil {
				t.Fatal(err)
			}
			closure, reference, err := p.Reference(world, "run")
			if err != nil || !slices.Equal(closure, actual.Closure) || !reflect.DeepEqual(reference, actual.Result) {
				t.Fatalf("world %d: %+v / %+v, %v", raw, actual.Result, reference, err)
			}
		}
	}
}
