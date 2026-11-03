package evidence

import (
	"reflect"
	"testing"
)

func TestContradictedPremiseDoesNotInventOppositeConclusion(t *testing.T) {
	theory := NormTheory{
		Atoms:    2,
		Evidence: []SignedRule{{"r", "1", []SignedLiteral{{Atom: 0}}, SignedLiteral{Atom: 1}}},
		Norms: []NormRule{{ID: "permit", Version: "1", Action: "run", Effect: Permit,
			Conditions: []SignedLiteral{{Atom: 1}}}},
	}
	p, err := CompileNormative(theory)
	if err != nil {
		t.Fatal(err)
	}
	world, err := p.NewWorld([]State{Both, Unknown})
	if err != nil {
		t.Fatal(err)
	}
	certificate, err := p.Evaluate(world, "run")
	if err != nil || certificate.Closure[0] != Both || certificate.Closure[1] != True || certificate.Result.Decision != NormAllow {
		t.Fatal("declared signed-Horn semantics changed; it is not contradiction-taint propagation", certificate, err)
	}
}

func TestSignedLiteralCyclesAndMalformedWorldsRejected(t *testing.T) {
	theory := normFixture()
	theory.Evidence = []SignedRule{
		{"forward", "1", []SignedLiteral{{Atom: 0}}, SignedLiteral{Atom: 0, Negative: true}},
		{"back", "1", []SignedLiteral{{Atom: 0, Negative: true}}, SignedLiteral{Atom: 0}},
	}
	if _, err := CompileNormative(theory); err == nil {
		t.Fatal("recursive signed fragment accepted")
	}
	p, err := CompileNormative(normFixture())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := p.NewWorld([]State{State(4), Unknown, Unknown}); err == nil {
		t.Fatal("invalid logical state accepted")
	}
	world, err := p.NewWorld([]State{True, False, Unknown})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := p.Evaluate(world, "missing-action"); err == nil {
		t.Fatal("unknown action was silently mapped to a successful refusal")
	}
}

func FuzzNormativeWorld(f *testing.F) {
	f.Add(uint8(2), uint8(1), uint8(1), true)
	f.Fuzz(func(t *testing.T, premise, guard, order uint8, negated bool) {
		theory := normFixture()
		theory.Norms[0].Conditions[0].Negative = negated
		if order%3 == 1 {
			theory.Priorities = []NormPriority{{"permit", "ban"}}
		} else if order%3 == 2 {
			theory.Priorities = []NormPriority{{"ban", "permit"}}
		}
		p, err := CompileNormative(theory)
		if err != nil {
			t.Fatal(err)
		}
		world, err := p.NewWorld([]State{State(premise % 4), State(guard % 4), Unknown})
		if err != nil {
			t.Fatal(err)
		}
		cache := NewNormCache()
		for step := 0; step < 3; step++ {
			_, reference, err := p.Reference(world, "run")
			if err != nil {
				t.Fatal(err)
			}
			result, err := cache.Resolve(p, world, "run")
			if err != nil || !reflect.DeepEqual(reference, result.Certificate.Result) {
				t.Fatal("normative cache disagrees with current context", result, err)
			}
			if step == 1 {
				world.Revision++
				world.Facts[1] = State((int(world.Facts[1]) + 1) % 4)
			}
		}
	})
}
