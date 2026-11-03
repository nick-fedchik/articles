package evidence

import (
	"context"
	"os/exec"
	"slices"
	"testing"
	"time"
)

func TestNormativeSATReference(t *testing.T) {
	solver, err := exec.LookPath("minisat")
	if err != nil {
		t.Skip("external MiniSat absent; normative study command requires it explicitly")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	theory := normFixture()
	theory.Priorities = []NormPriority{{"permit", "ban"}}
	p, err := CompileNormative(theory)
	if err != nil {
		t.Fatal(err)
	}
	for _, premise := range []State{Unknown, False, True, Both} {
		for _, guard := range []State{Unknown, False, True, Both} {
			world, err := p.NewWorld([]State{premise, guard, Both})
			if err != nil {
				t.Fatal(err)
			}
			actual, err := p.Evaluate(world, "run")
			if err != nil {
				t.Fatal(err)
			}
			sat, err := p.SATReference(ctx, solver, world, "run")
			if err != nil || !slices.Equal(sat.Closure, actual.Closure) || sat.Decision != actual.Result.Decision ||
				sat.ConditionEvidence != actual.Result.ConditionEvidence || sat.Applicable != actual.Result.Applicable ||
				sat.Surviving != actual.Result.Surviving {
				t.Fatal(actual, sat, err)
			}
		}
	}
}

func TestNewExceptionInvalidatesSameRevisionCache(t *testing.T) {
	theory := normFixture()
	theory.Norms[0].Exception = nil
	p, err := CompileNormative(theory)
	if err != nil {
		t.Fatal(err)
	}
	world, err := p.NewWorld([]State{True, Unknown, Unknown})
	if err != nil {
		t.Fatal(err)
	}
	world.NormEnabled[1] = false
	cache := NewNormCache()
	first, err := cache.Resolve(p, world, "run")
	if err != nil || first.Certificate.Result.Decision != NormAllow {
		t.Fatal(first, err)
	}
	guard := SignedLiteral{Atom: 1}
	theory.Norms[0].Exception = &guard
	theory.Norms[0].Version = "2"
	changed, err := CompileNormative(theory)
	if err != nil {
		t.Fatal(err)
	}
	result, err := cache.Resolve(changed, world, "run")
	if err != nil || !result.Invalidated || result.Reused || result.Certificate.Result.Decision != NormUnknown {
		t.Fatal("new exception outside the positive proof did not invalidate allowance", result, err)
	}
}

func TestPriorityCannotHideUnknownLowerRuleException(t *testing.T) {
	theory := normFixture()
	theory.Norms[0].Exception = nil
	guard := SignedLiteral{Atom: 1}
	theory.Norms[1].Exception = &guard
	theory.Priorities = []NormPriority{{"permit", "ban"}}
	p, err := CompileNormative(theory)
	if err != nil {
		t.Fatal(err)
	}
	world, err := p.NewWorld([]State{True, Unknown, Unknown})
	if err != nil {
		t.Fatal(err)
	}
	certificate, err := p.Evaluate(world, "run")
	if err != nil || certificate.Result.Decision != NormUnknown {
		t.Fatal(certificate, err)
	}
}
