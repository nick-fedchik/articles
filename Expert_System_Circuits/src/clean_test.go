package evidence

import (
	"errors"
	"fmt"
	"reflect"
	"slices"
	"testing"
)

func cleanFixture() ([]SignedRule, []State) {
	return []SignedRule{
		{"tainted", "1", []SignedLiteral{{Atom: 0}}, SignedLiteral{Atom: 2}},
		{"alternative", "1", []SignedLiteral{{Atom: 1}}, SignedLiteral{Atom: 2}},
	}, []State{Both, True, Unknown}
}

func TestCleanAlternativesAndForgedCertificates(t *testing.T) {
	rules, facts := cleanFixture()
	p, err := CompileCleanSupport(3, rules)
	if err != nil {
		t.Fatal(err)
	}
	world := SignedWorld{1, facts, []bool{true, true}}
	query := SignedLiteral{Atom: 2}
	cert, work, err := p.EvaluateClean(world, query)
	if err != nil || cert.Status != CleanSupported || cert.Proof == nil ||
		!slices.Equal(cert.Proof.Rules, []int{1}) || work.Nodes != 6 {
		t.Fatalf("%+v %+v %v", cert, work, err)
	}
	if err := p.VerifyClean(world, query, cert); err != nil {
		t.Fatal(err)
	}
	ref, err := p.ReferenceClean(world, query, 1000)
	if err != nil || !ref.ExistsClean || ref.AllProofsClean || ref.Proofs != 2 {
		t.Fatalf("%+v %v", ref, err)
	}
	for _, mutate := range []func(*CleanCertificate){
		func(c *CleanCertificate) { c.Closure[1] = false },
		func(c *CleanCertificate) { c.Clean[0] = true },
		func(c *CleanCertificate) { c.Status = NotSupported },
		func(c *CleanCertificate) { c.Proof = nil },
		func(c *CleanCertificate) { c.Proof.Facts = []int{0}; c.Proof.Rules = []int{0} },
	} {
		fresh, _, err := p.EvaluateClean(world, query)
		if err != nil {
			t.Fatal(err)
		}
		mutate(&fresh)
		if err := p.VerifyClean(world, query, fresh); err == nil {
			t.Fatal("forged clean certificate accepted")
		}
	}
	world.Enabled[1] = false
	next, _, err := p.EvaluateClean(world, query)
	if err != nil || next.Status != ConflictDependent {
		t.Fatalf("%+v %v", next, err)
	}
	if err := p.VerifyClean(world, query, cert); err == nil {
		t.Fatal("stale context accepted")
	}
	world.Facts[0] = Unknown
	next, _, err = p.EvaluateClean(world, query)
	if err != nil || next.Status != NotSupported {
		t.Fatalf("%+v %v", next, err)
	}
}

func TestCleanLeavesAreNotCleanIntermediateAtoms(t *testing.T) {
	rules := []SignedRule{
		{"pro", "1", []SignedLiteral{{Atom: 0}}, SignedLiteral{Atom: 2}},
		{"contra", "1", []SignedLiteral{{Atom: 1}}, SignedLiteral{Atom: 2, Negative: true}},
		{"query", "1", []SignedLiteral{{Atom: 2}}, SignedLiteral{Atom: 3}},
	}
	p, err := CompileCleanSupport(4, rules)
	if err != nil {
		t.Fatal(err)
	}
	world := SignedWorld{1, []State{True, True, Unknown, Unknown}, []bool{true, true, true}}
	query := SignedLiteral{Atom: 3}
	cert, _, err := p.EvaluateClean(world, query)
	if err != nil || cert.Status != ConflictDependent {
		t.Fatalf("%+v %v", cert, err)
	}
	ref, err := p.ReferenceClean(world, query, 1000)
	if err != nil || !ref.LeavesOnly || ref.ExistsClean {
		t.Fatalf("%+v %v", ref, err)
	}
}

func TestCleanLayerDoesNotChangeNormativePolicy(t *testing.T) {
	theory := NormTheory{Atoms: 2, Evidence: []SignedRule{
		{"derive", "1", []SignedLiteral{{Atom: 0}}, SignedLiteral{Atom: 1}},
	}, Norms: []NormRule{{ID: "permit", Version: "1", Action: "run", Effect: Permit, Conditions: []SignedLiteral{{Atom: 1}}}}}
	norm, err := CompileNormative(theory)
	if err != nil {
		t.Fatal(err)
	}
	world, err := norm.NewWorld([]State{Both, Unknown})
	if err != nil {
		t.Fatal(err)
	}
	before, err := norm.Evaluate(world, "run")
	if err != nil || before.Result.Decision != NormAllow {
		t.Fatalf("%+v %v", before, err)
	}
	p, err := CompileCleanSupport(theory.Atoms, theory.Evidence)
	if err != nil {
		t.Fatal(err)
	}
	cert, _, err := p.EvaluateClean(SignedWorld{world.Revision, world.Facts, world.EvidenceEnabled}, SignedLiteral{Atom: 1})
	if err != nil || cert.Status != ConflictDependent {
		t.Fatalf("%+v %v", cert, err)
	}
	after, err := norm.Evaluate(world, "run")
	if err != nil || !reflect.DeepEqual(before, after) {
		t.Fatal("experimental clean support altered normative semantics", err)
	}
}

func TestCleanProofInvalidatedByOppositionOutsidePositiveWitness(t *testing.T) {
	p, err := CompileCleanSupport(3, []SignedRule{
		{"derive", "1", []SignedLiteral{{Atom: 0}}, SignedLiteral{Atom: 2}},
		{"opposition", "1", []SignedLiteral{{Atom: 1}}, SignedLiteral{Atom: 0, Negative: true}},
	})
	if err != nil {
		t.Fatal(err)
	}
	world := SignedWorld{1, []State{True, Unknown, Unknown}, []bool{true, true}}
	query := SignedLiteral{Atom: 2}
	old, _, err := p.EvaluateClean(world, query)
	if err != nil || old.Status != CleanSupported || !slices.Equal(old.Proof.Facts, []int{0}) {
		t.Fatalf("%+v %v", old, err)
	}
	world.Facts[1], world.Revision = True, 2
	current, _, err := p.EvaluateClean(world, query)
	if err != nil || !current.Closure[query.code()] || current.Status != ConflictDependent {
		t.Fatalf("%+v %v", current, err)
	}
	if err := p.VerifyClean(world, query, old); err == nil {
		t.Fatal("opposition outside the old positive witness was ignored")
	}
}

func FuzzCleanSupport(f *testing.F) {
	f.Add([]byte{3, 0, 3, 2, 1, 0, 4, 9, 3, 5, 6, 2, 7, 1})
	f.Add([]byte{1, 2, 0, 3, 1, 0, 1, 2})
	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) < 8 {
			t.Skip()
		}
		at := func(i int) int { return int(data[i%len(data)]) }
		atoms := 1 + at(0)%5
		rules := []SignedRule{}
		for i := 0; i < at(1)%9; i++ {
			head := 1 + at(4*i+2)%(2*atoms-1)
			first := at(4*i+3) % head
			premises := []SignedLiteral{{Atom: first / 2, Negative: first%2 != 0}}
			if head > 1 && at(4*i+4)%2 != 0 {
				second := at(4*i+5) % head
				if second == first {
					second = (second + 1) % head
				}
				premises = append(premises, SignedLiteral{Atom: second / 2, Negative: second%2 != 0})
			}
			rules = append(rules, SignedRule{fmt.Sprintf("r%d", i), "1", premises, SignedLiteral{head / 2, head%2 != 0}})
		}
		p, err := CompileCleanSupport(atoms, rules)
		if err != nil {
			t.Fatal(err)
		}
		world := SignedWorld{Revision: 1, Facts: make([]State, atoms), Enabled: make([]bool, len(rules))}
		for i := range world.Facts {
			world.Facts[i] = State(at(i+3) % 4)
		}
		for i := range world.Enabled {
			world.Enabled[i] = at(i+2)%2 != 0
		}
		code := at(7) % (2 * atoms)
		query := SignedLiteral{code / 2, code%2 != 0}
		cert, _, err := p.EvaluateClean(world, query)
		if err != nil {
			t.Fatal(err)
		}
		if err := p.VerifyClean(world, query, cert); err != nil {
			t.Fatal(err)
		}
		ref, err := p.ReferenceClean(world, query, 100000)
		if errors.Is(err, ErrCleanProofBudget) {
			t.Skip("independent oracle exhausted its explicit budget")
		}
		if err != nil || ref.ExistsClean != (cert.Status == CleanSupported) || !slices.Equal(ref.Closure, cert.Closure) {
			t.Fatalf("reference mismatch: %+v %+v %v", cert, ref, err)
		}
	})
}

func TestCleanInputValidationAndEnumerationBudget(t *testing.T) {
	rules, facts := cleanFixture()
	p, err := CompileCleanSupport(3, rules)
	if err != nil {
		t.Fatal(err)
	}
	world := SignedWorld{1, facts, []bool{true, true}}
	if _, err := p.ReferenceClean(world, SignedLiteral{Atom: 2}, 1); !errors.Is(err, ErrCleanProofBudget) {
		t.Fatal("enumeration budget did not fail explicitly", err)
	}
	for _, invalid := range []SignedWorld{{0, facts, world.Enabled}, {1, []State{True}, world.Enabled},
		{1, []State{State(4), True, Unknown}, world.Enabled}, {1, facts, nil}} {
		if _, _, err := p.EvaluateClean(invalid, SignedLiteral{}); err == nil {
			t.Fatal("invalid world accepted")
		}
	}
	if _, _, err := p.EvaluateClean(world, SignedLiteral{Atom: 3}); err == nil {
		t.Fatal("invalid query accepted")
	}
	if _, err := CompileCleanSupport(0, nil); err == nil {
		t.Fatal("invalid model size accepted")
	}
	cyclic := append(slices.Clone(rules), SignedRule{"cycle", "1", []SignedLiteral{{Atom: 2}}, SignedLiteral{Atom: 0}})
	if _, err := CompileCleanSupport(3, cyclic); err == nil {
		t.Fatal("cycle accepted")
	}
}
