package evidence

import (
	"reflect"
	"testing"
)

func TestPriorityChangeInvalidatesWithoutWorldRevisionChange(t *testing.T) {
	theory := normFixture()
	theory.Priorities = []NormPriority{{"permit", "ban"}}
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
	theory.Priorities = []NormPriority{{"ban", "permit"}}
	changed, err := CompileNormative(theory)
	if err != nil {
		t.Fatal(err)
	}
	current, err := cache.Resolve(changed, world, "run")
	if err != nil || !current.Invalidated || current.Reused || current.Certificate.Result.Decision != NormDeny ||
		current.Certificate.Result.ConditionEvidence != Both {
		t.Fatal("priority change was hidden by unchanged positive support or world revision", current, err)
	}
}

func TestNormativeCacheDoesNotHideCertificateCorruption(t *testing.T) {
	p, err := CompileNormative(normFixture())
	if err != nil {
		t.Fatal(err)
	}
	world, err := p.NewWorld([]State{True, False, Unknown})
	if err != nil {
		t.Fatal(err)
	}
	cache := NewNormCache()
	if _, err := cache.Resolve(p, world, "run"); err != nil {
		t.Fatal(err)
	}
	entry := cache.entries["run"]
	entry.Result.Decision = NormAllow
	cache.entries["run"] = entry
	if _, err := cache.Resolve(p, world, "run"); err == nil {
		t.Fatal("corrupted certificate was accepted or silently replaced")
	}
}

func TestCompiledNormativeTheoryOwnsItsInputs(t *testing.T) {
	theory := normFixture()
	p, err := CompileNormative(theory)
	if err != nil {
		t.Fatal(err)
	}
	world, err := p.NewWorld([]State{True, False, Unknown})
	if err != nil {
		t.Fatal(err)
	}
	before, err := p.Evaluate(world, "run")
	if err != nil {
		t.Fatal(err)
	}
	theory.Norms[0].Conditions[0].Negative = true
	theory.Norms[0].Exception.Atom = 2
	after, err := p.Evaluate(world, "run")
	if err != nil || !reflect.DeepEqual(before, after) {
		t.Fatal("caller mutated an immutable compiled normative model", err)
	}
}
