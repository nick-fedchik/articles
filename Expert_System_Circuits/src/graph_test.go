package evidence

import (
	"fmt"
	"math"
	"math/rand"
	"slices"
	"testing"
)

func fixtureSources(t testing.TB, count int) []Source {
	t.Helper()
	out := make([]Source, count)
	for i := range out {
		var err error
		out[i], err = NewSource(fmt.Sprintf("s%d", i), "1", []byte(fmt.Sprintf("public source %d", i)), .9)
		if err != nil {
			t.Fatal(err)
		}
	}
	return out
}

func equalResults(a, b Result) bool {
	if a.State != b.State || math.Abs(a.Weight-b.Weight) > 1e-12 || len(a.Sources) != len(b.Sources) {
		return false
	}
	for _, source := range a.Sources {
		if !slices.Contains(b.Sources, source) {
			return false
		}
	}
	return true
}

func TestCSEAndUniqueSourceSupport(t *testing.T) {
	a := &Expr{Op: Input, Input: 0}
	duplicate := &Expr{Op: All, Left: &Expr{Op: Input, Input: 0}, Right: &Expr{Op: Input, Input: 0}}
	root := &Expr{Op: All, Left: &Expr{Op: All, Left: a, Right: a}, Right: duplicate}
	p, err := Compile(root, fixtureSources(t, 1))
	if err != nil {
		t.Fatal(err)
	}
	got, err := p.Evaluate([]State{True})
	if err != nil || p.Nodes() != 3 || math.Abs(got.Weight-.9) > 1e-12 || len(got.Sources) != 1 {
		t.Fatalf("nodes=%d, result=%+v, err=%v", p.Nodes(), got, err)
	}
}

func TestMultiwordProvenance(t *testing.T) {
	root := &Expr{Op: Input, Input: 0}
	for i := 1; i < 130; i++ {
		root = &Expr{Op: All, Left: root, Right: &Expr{Op: Input, Input: i}}
	}
	inputs := make([]State, 130)
	for i := range inputs {
		inputs[i] = True
	}
	sources := fixtureSources(t, len(inputs))
	p, err := Compile(root, sources)
	if err != nil {
		t.Fatal(err)
	}
	got, err := p.Evaluate(inputs)
	want, refErr := ReferenceEvaluate(root, sources, inputs)
	if err != nil || refErr != nil || !equalResults(got, want) || len(got.Sources) != 130 {
		t.Fatalf("multiword mismatch: %v / %v", err, refErr)
	}
}

func TestCompilerMatchesIndependentReference(t *testing.T) {
	rng := rand.New(rand.NewSource(20261102))
	sources := fixtureSources(t, 6)
	var tree func(int) *Expr
	tree = func(depth int) *Expr {
		if depth == 0 || rng.Intn(3) == 0 {
			return &Expr{Op: Input, Input: rng.Intn(len(sources))}
		}
		expr := &Expr{Op: All, Left: tree(depth - 1), Right: tree(depth - 1)}
		if rng.Intn(3) == 0 {
			expr.Op = Select
			expr.Condition = &Expr{Op: Input, Input: rng.Intn(len(sources))}
		}
		return expr
	}
	for n := 0; n < 1000; n++ {
		root := tree(5)
		p, err := Compile(root, sources)
		if err != nil {
			t.Fatal(err)
		}
		inputs := make([]State, len(sources))
		for i := range inputs {
			inputs[i] = State(rng.Intn(4))
		}
		got, err := p.Evaluate(inputs)
		want, refErr := ReferenceEvaluate(root, sources, inputs)
		if err != nil || refErr != nil || !equalResults(got, want) {
			t.Fatalf("case %d: got %+v, want %+v (%v, %v)", n, got, want, err, refErr)
		}
	}
}

func TestResolvedSelectorUsesOnlyActiveSupport(t *testing.T) {
	root := &Expr{Op: Select, Condition: &Expr{Op: Input, Input: 0},
		Left: &Expr{Op: Input, Input: 1}, Right: &Expr{Op: Input, Input: 2}}
	p, err := Compile(root, fixtureSources(t, 3))
	if err != nil {
		t.Fatal(err)
	}
	for _, condition := range []State{True, False, Unknown, Both} {
		got, err := p.Evaluate([]State{condition, True, Both})
		wantCount := 3
		if condition == True || condition == False {
			wantCount = 2
		}
		if err != nil || len(got.Sources) != wantCount {
			t.Fatalf("condition %v: %+v, %v", condition, got, err)
		}
		if condition == False && got.State != True {
			t.Fatal("negative selector must not negate the selected rule")
		}
	}
}

func TestGraphRejectsInvalidInputs(t *testing.T) {
	sources := fixtureSources(t, 1)
	cycle := &Expr{Op: All, Right: &Expr{Op: Input}}
	cycle.Left = cycle
	for _, root := range []*Expr{nil, cycle, {Op: Input, Input: -1}, {Op: Op(99)}} {
		if _, err := Compile(root, sources); err == nil {
			t.Fatalf("accepted malformed graph %+v", root)
		}
		if _, err := ReferenceEvaluate(root, sources, []State{True}); err == nil {
			t.Fatal("reference accepted malformed graph")
		}
	}
	conflict := append(slices.Clone(sources), sources[0])
	conflict[1].Weight = .8
	if _, err := Compile(&Expr{Op: Input}, conflict); err == nil {
		t.Fatal("conflicting snapshot accepted")
	}
	p, _ := Compile(&Expr{Op: Input}, sources)
	if _, err := p.Evaluate([]State{State(4)}); err == nil {
		t.Fatal("invalid state accepted")
	}
}
