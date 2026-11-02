package evidence

import (
	"math"
	"slices"
	"sync"
	"testing"
)

func TestProvenanceSelectorChangesRequireRebuild(t *testing.T) {
	sources := fixtureSources(t, 3)
	root := &Expr{Op: Select, Condition: &Expr{Op: Input},
		Left: &Expr{Op: Input, Input: 1}, Right: &Expr{Op: Input, Input: 2}}
	p, err := Compile(root, sources)
	if err != nil {
		t.Fatal(err)
	}
	for _, format := range provenanceFormats() {
		for _, condition := range []State{True, False, Unknown, Both} {
			states := []State{condition, True, True}
			snapshot, err := p.BuildProvenance(states, format, 64<<20)
			if err != nil {
				t.Fatal(err)
			}
			got, err := snapshot.Result()
			want, refErr := ReferenceSupport(root, sources, states)
			if err != nil || refErr != nil || !equalResults(got, want) {
				t.Fatal("selector changed the source set incorrectly", format, condition, err, refErr)
			}
		}
	}
}

func TestProvenanceExactSourceDuplicatesAndSmallWeights(t *testing.T) {
	sources := fixtureSources(t, 2)
	sources[1] = sources[0]
	p, err := Compile(&Expr{Op: All, Left: &Expr{Op: Input}, Right: &Expr{Op: Input, Input: 1}}, sources)
	if err != nil {
		t.Fatal(err)
	}
	for _, format := range provenanceFormats() {
		snapshot, err := p.BuildProvenance([]State{True, True}, format, 64<<20)
		if err != nil {
			t.Fatal(err)
		}
		result, err := snapshot.Result()
		if err != nil || len(result.Sources) != 1 || math.Abs(result.Weight-.9) > 1e-14 {
			t.Fatal("an exact source duplicate was counted twice", format, result, err)
		}
		tiny, err := snapshot.WithSourceWeight(sources[0].ID, sources[0].Version, 1e-12)
		if err != nil {
			t.Fatal(err)
		}
		value, err := tiny.Result()
		if err != nil || math.Abs((value.Weight-1e-12)/1e-12) > 1e-12 {
			t.Fatal("absolute tolerance concealed a small-weight error", value, err)
		}
	}
}

func TestInterningResolvesHashBucketCollisionsExactly(t *testing.T) {
	hash := func(values []uint32) uint64 {
		out := uint64(14695981039346656037)
		for _, value := range values {
			for shift := uint(0); shift < 32; shift += 8 {
				out ^= uint64(byte(value >> shift))
				out *= 1099511628211
			}
		}
		return out
	}
	seen := map[uint64]uint32{}
	var first, second uint32
	for value := uint32(1); value < 100; value++ {
		slot := hash([]uint32{0, value}) & 31
		if previous, ok := seen[slot]; ok {
			first, second = previous, value
			break
		}
		seen[slot] = value
	}
	if first == 0 || second == 0 {
		t.Fatal("no deliberate bucket collision found")
	}
	pool, err := newProvenancePool(&provenanceBudget{limit: 10000}, 8, true)
	if err != nil {
		t.Fatal(err)
	}
	a, err := pool.union([]uint32{0}, []uint32{first})
	if err != nil {
		t.Fatal(err)
	}
	b, err := pool.union([]uint32{0}, []uint32{second})
	if err != nil || slices.Equal(a, b) || len(pool.entries) != 2 {
		t.Fatal("hash bucket collision merged different source sets", a, b, err)
	}
}

func TestConcurrentWeightSnapshotsDoNotMutateSharedSupport(t *testing.T) {
	c, err := GenerateProvenanceCase("repeated-support", 128, "uncertain", 81001)
	if err != nil {
		t.Fatal(err)
	}
	base, err := c.Program.BuildProvenance(c.Inputs, ProvenanceInterned, 64<<20)
	if err != nil {
		t.Fatal(err)
	}
	original, err := base.Result()
	if err != nil {
		t.Fatal(err)
	}
	var workers sync.WaitGroup
	for i := 0; i < 4; i++ {
		workers.Add(1)
		go func() {
			defer workers.Done()
			for step := 0; step < 20; step++ {
				updated, err := base.WithSourceWeight(c.Sources[0].ID, "1", .5)
				if err != nil {
					t.Error(err)
					return
				}
				if _, err := updated.Result(); err != nil {
					t.Error(err)
				}
			}
		}()
	}
	workers.Wait()
	current, err := base.Result()
	if err != nil || !equalResults(original, current) {
		t.Fatal("shared immutable support was changed", err)
	}
}

func FuzzProvenanceFormats(f *testing.F) {
	f.Add(uint8(2), uint8(0), uint8(3))
	f.Fuzz(func(t *testing.T, condition, left, right uint8) {
		sources := fixtureSources(t, 3)
		root := &Expr{Op: Select, Condition: &Expr{Op: Input},
			Left: &Expr{Op: Input, Input: 1}, Right: &Expr{Op: Input, Input: 2}}
		p, err := Compile(root, sources)
		if err != nil {
			t.Fatal(err)
		}
		states := []State{State(condition % 4), State(left % 4), State(right % 4)}
		want, err := ReferenceSupport(root, sources, states)
		if err != nil {
			t.Fatal(err)
		}
		for _, format := range provenanceFormats() {
			snapshot, err := p.BuildProvenance(states, format, 100000)
			if err != nil {
				t.Fatal(err)
			}
			got, err := snapshot.Result()
			if err != nil || !equalResults(got, want) {
				t.Fatal("format changed semantics", format, err)
			}
		}
	})
}
