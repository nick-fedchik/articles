package evidence

import "testing"

func TestAbsentExceptionIsNotUnknownException(t *testing.T) {
	theory := normFixture()
	theory.Norms = theory.Norms[:1]
	theory.Norms[0].Exception = nil
	p, err := CompileNormative(theory)
	if err != nil {
		t.Fatal(err)
	}
	world, err := p.NewWorld([]State{True, Unknown, Unknown})
	if err != nil {
		t.Fatal(err)
	}
	without, err := p.Evaluate(world, "run")
	if err != nil || without.Result.Decision != NormAllow || without.Result.Rules[0].HasException {
		t.Fatal("absent exception was treated as unknown", without, err)
	}
	guard := SignedLiteral{Atom: 1}
	theory.Norms[0].Exception = &guard
	withProgram, err := CompileNormative(theory)
	if err != nil {
		t.Fatal(err)
	}
	with, err := withProgram.Evaluate(world, "run")
	if err != nil || with.Result.Decision != NormUnknown || !with.Result.Rules[0].HasException ||
		with.Result.Rules[0].Exception != Unknown {
		t.Fatal("unknown declared exception was treated as absent", with, err)
	}
}

func TestNewProhibitionOutsidePositiveProofInvalidatesAllowance(t *testing.T) {
	theory := normFixture()
	theory.Norms = theory.Norms[:1]
	p, err := CompileNormative(theory)
	if err != nil {
		t.Fatal(err)
	}

	world, err := p.NewWorld([]State{True, False, Unknown})
	if err != nil {
		t.Fatal(err)
	}
	cache := NewNormCache()
	first, err := cache.Resolve(p, world, "run")
	if err != nil || first.Certificate.Result.Decision != NormAllow {
		t.Fatal(first, err)
	}
	theory.Norms = append(theory.Norms, NormRule{ID: "new-ban", Version: "1", Action: "run", Effect: Prohibit})
	changed, err := CompileNormative(theory)
	if err != nil {
		t.Fatal(err)
	}
	newWorld, err := changed.NewWorld(world.Facts)
	if err != nil {
		t.Fatal(err)
	}
	result, err := cache.Resolve(changed, newWorld, "run")
	if err != nil || !result.Invalidated || result.Reused || result.Certificate.Result.Decision != NormConflict ||
		result.Certificate.Result.ConditionEvidence != Both {
		t.Fatal("structurally new prohibition was hidden by unchanged positive support", result, err)
	}
	if err := changed.Verify(newWorld, "run", first.Certificate); err == nil {
		t.Fatal("old permit certificate accepted after a new opposing rule")
	}
}
