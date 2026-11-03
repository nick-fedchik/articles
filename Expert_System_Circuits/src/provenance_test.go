package evidence

import (
	"errors"
	"math"
	"math/rand"
	"slices"
	"testing"
)

func provenanceFormats() []ProvenanceFormat {
	return []ProvenanceFormat{ProvenanceDense, ProvenanceSparse, ProvenanceInterned, ProvenanceTrace}
}

func TestProvenanceFormatsMatchIndependentExpression(t *testing.T) {
	rng := rand.New(rand.NewSource(20261104))
	sources := fixtureSources(t, 6)
	var tree func(int) *Expr
	tree = func(depth int) *Expr {
		if depth == 0 || rng.Intn(3) == 0 {
			return &Expr{Op: Input, Input: rng.Intn(len(sources))}
		}
		out := &Expr{Op: All, Left: tree(depth - 1), Right: tree(depth - 1)}
		if rng.Intn(2) == 0 {
			out.Op, out.Condition = Select, &Expr{Op: Input, Input: rng.Intn(len(sources))}
		}
		return out
	}
	for trial := 0; trial < 400; trial++ {
		root := tree(5)
		program, err := Compile(root, sources)
		if err != nil {
			t.Fatal(err)
		}
		inputs := make([]State, len(sources))
		for i := range inputs {
			inputs[i] = State(rng.Intn(4))
		}
		want, err := ReferenceSupport(root, sources, inputs)
		if err != nil {
			t.Fatal(err)
		}
		original, err := ReferenceEvaluate(root, sources, inputs)
		if err != nil || !equalResults(original, want) {
			t.Fatal("the new oracle does not match the previous recursive oracle", err)
		}
		for _, format := range provenanceFormats() {
			snapshot, err := program.BuildProvenance(inputs, format, 64<<20)
			if err != nil {
				t.Fatal(err)
			}
			got, err := snapshot.Result()
			if err != nil || !equalResults(got, want) {
				t.Fatalf("case=%d format=%s got=%+v want=%+v err=%v", trial, format, got, want, err)
			}
		}
	}
}

func TestProvenanceBudgetsAndRootStorage(t *testing.T) {
	sources := fixtureSources(t, 130)
	root := balancedExpression(0, len(sources))
	program, err := Compile(root, sources)
	if err != nil {
		t.Fatal(err)
	}
	inputs := make([]State, len(sources))
	for i := range inputs {
		inputs[i] = True
	}
	for _, format := range provenanceFormats() {
		snapshot, err := program.BuildProvenance(inputs, format, 64<<20)
		if err != nil {
			t.Fatal(err)
		}
		stats := snapshot.Stats()
		if stats.RootSources != 130 || stats.RetainedBytes >= stats.ManagedWorkingBytes {
			t.Fatal("root retained the full working arena", format, stats)
		}
		if _, err := program.BuildProvenance(inputs, format, stats.ManagedWorkingBytes-1); !errors.Is(err, ErrProvenanceBudget) {
			t.Fatal("memory budget ignored", format, err)
		}
		if _, err := program.BuildProvenance(inputs, format, stats.ManagedWorkingBytes); err != nil {
			t.Fatal("exact managed budget was insufficient", format, err)
		}
	}
}

func TestInterningSharesEqualSetsWithoutMutatingInputs(t *testing.T) {
	a, b, c := []uint32{1, 3}, []uint32{2, 4}, []uint32{1, 2}
	budget := &provenanceBudget{limit: 100000}
	pool, err := newProvenancePool(budget, 8, true)
	if err != nil {
		t.Fatal(err)
	}
	first, err := pool.union(a, b)
	if err != nil {
		t.Fatal(err)
	}
	second, err := pool.union(c, []uint32{3, 4})
	if err != nil || &first[0] != &second[0] || !slices.Equal(first, []uint32{1, 2, 3, 4}) ||
		!slices.Equal(a, []uint32{1, 3}) {
		t.Fatal("equal support was not shared exactly", err)
	}
}

func TestWeightOnlySnapshotAndVersionIdentity(t *testing.T) {
	sources := fixtureSources(t, 2)
	sources[1].ID = sources[0].ID
	sources[1].Version = "2"
	root := &Expr{Op: All, Left: &Expr{Op: Input}, Right: &Expr{Op: Input, Input: 1}}
	program, err := Compile(root, sources)
	if err != nil {
		t.Fatal(err)
	}
	for _, format := range provenanceFormats() {
		snapshot, err := program.BuildProvenance([]State{True, True}, format, 64<<20)
		if err != nil {
			t.Fatal(err)
		}
		original, err := snapshot.Result()
		if err != nil || len(original.Sources) != 2 {
			t.Fatal("different versions were merged", err)
		}
		updated, err := snapshot.WithSourceWeight(sources[0].ID, sources[0].Version, 0)
		if err != nil {
			t.Fatal(err)
		}
		zero, err := updated.Result()
		again, oldErr := snapshot.Result()
		if err != nil || oldErr != nil || zero.Weight != 0 || !equalResults(original, again) {
			t.Fatal("immutable weight update failed", err, oldErr)
		}
		if _, err := snapshot.WithSourceWeight(sources[0].ID, "missing", .5); err == nil {
			t.Fatal("unknown source version silently accepted")
		}
	}
}

func TestProvenanceMalformedInputs(t *testing.T) {
	p, err := Compile(&Expr{Op: Input}, fixtureSources(t, 1))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := p.BuildProvenance([]State{State(4)}, ProvenanceSparse, 1000); err == nil {
		t.Fatal("invalid logical state accepted")
	}
	if _, err := p.BuildProvenance([]State{True}, ProvenanceFormat("unknown"), 1000); err == nil {
		t.Fatal("unknown representation accepted")
	}
	snapshot, err := p.BuildProvenance([]State{True}, ProvenanceTrace, 10000)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := snapshot.WithSourceWeight("s0", "1", math.NaN()); err == nil {
		t.Fatal("invalid weight accepted")
	}
}
