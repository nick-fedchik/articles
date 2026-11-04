package evidence

import (
	"math"
	"reflect"
	"slices"
	"testing"
)

func cacheFixture(t *testing.T, mode CleanCacheMode, enabled bool) (*CleanProgram, *CleanWitnessCache) {
	t.Helper()
	p, err := CompileCleanSupport(5, []SignedRule{
		{"derive", "1", []SignedLiteral{{Atom: 0}}, SignedLiteral{Atom: 2}},
		{"opposition", "1", []SignedLiteral{{Atom: 1}}, SignedLiteral{Atom: 0, Negative: true}},
	})
	if err != nil {
		t.Fatal(err)
	}
	cache, err := NewCleanWitnessCache(p, SignedWorld{1, []State{True, Unknown, Unknown, Unknown, Unknown}, []bool{true, enabled}}, mode)
	if err != nil {
		t.Fatal(err)
	}
	return p, cache
}

func checkCache(t *testing.T, p *CleanProgram, cache *CleanWitnessCache, query SignedLiteral) CleanCacheResult {
	t.Helper()
	result, err := cache.Resolve(query)
	if err != nil {
		t.Fatal(err)
	}
	world, err := cache.World()
	if err != nil {
		t.Fatal(err)
	}
	if err := p.VerifyCleanAnswer(world, query, result.Answer); err != nil {
		t.Fatal(err)
	}
	fresh, _, err := p.EvaluateClean(world, query)
	if err != nil || fresh.Status != result.Answer.Status {
		t.Fatalf("cache/fresh mismatch: %+v %v", result, err)
	}
	return result
}

func TestCleanCacheUnrelatedUpdatesAndOutsideOpposition(t *testing.T) {
	query := SignedLiteral{Atom: 2}
	for _, mode := range []CleanCacheMode{CleanRecompute, CleanRevision, CleanCone} {
		p, cache := cacheFixture(t, mode, true)
		initial := checkCache(t, p, cache, query)
		if initial.Answer.Status != CleanSupported || initial.Reused {
			t.Fatal(initial)
		}
		if err := cache.Update(HornUpdate{Facts: []HornChange{{Index: 6, Present: true}}}); err != nil {
			t.Fatal(err)
		}
		unrelated := checkCache(t, p, cache, query)
		if unrelated.Reused != (mode == CleanCone) {
			t.Fatalf("mode %s: %+v", mode, unrelated)
		}
		if err := cache.Update(HornUpdate{Facts: []HornChange{{Index: 2, Present: true}}}); err != nil {
			t.Fatal(err)
		}
		conflict := checkCache(t, p, cache, query)
		if conflict.Reused || conflict.Answer.Status != ConflictDependent {
			t.Fatal(conflict)
		}
	}
}

func TestCleanCacheDisabledProducersAndABA(t *testing.T) {
	query := SignedLiteral{Atom: 2}
	p, cache := cacheFixture(t, CleanCone, false)
	if err := cache.Update(HornUpdate{Facts: []HornChange{{Index: 2, Present: true}}}); err != nil {
		t.Fatal(err)
	}
	checkCache(t, p, cache, query)
	if err := cache.Update(HornUpdate{Rules: []HornChange{{Index: 1, Present: true}}}); err != nil {
		t.Fatal(err)
	}
	changed := checkCache(t, p, cache, query)
	if changed.Reused || changed.Invalidation != "RULE_EPOCH" || changed.Answer.Status != ConflictDependent {
		t.Fatal(changed)
	}
	p, cache = cacheFixture(t, CleanCone, true)
	checkCache(t, p, cache, query)
	for _, present := range []bool{true, false} {
		if err := cache.Update(HornUpdate{Facts: []HornChange{{Index: 2, Present: present}}}); err != nil {
			t.Fatal(err)
		}
	}
	restored := checkCache(t, p, cache, query)
	if restored.Reused || restored.Invalidation != "FACT_EPOCH" || restored.Answer.Status != CleanSupported {
		t.Fatal(restored)
	}
}

func TestSparseCleanWitnessRejectsForgery(t *testing.T) {
	p, cache := cacheFixture(t, CleanCone, false)
	if err := cache.Update(HornUpdate{Facts: []HornChange{{Index: 2, Present: true}}}); err != nil {
		t.Fatal(err)
	}
	result := checkCache(t, p, cache, SignedLiteral{Atom: 2})
	world, err := cache.World()
	if err != nil {
		t.Fatal(err)
	}
	for _, mutate := range []func(*CleanWitnessCertificate){
		func(c *CleanWitnessCertificate) {
			c.Absence = slices.DeleteFunc(c.Absence, func(v CleanSupportValue) bool { return v.Code == 2 })
		},
		func(c *CleanWitnessCertificate) { c.Absence = append(c.Absence, c.Absence[0]) },
		func(c *CleanWitnessCertificate) { c.Absence[0].Present = !c.Absence[0].Present },
		func(c *CleanWitnessCertificate) { c.Proof.Rules = []int{999} },
		func(c *CleanWitnessCertificate) { c.Revision++ },
	} {
		forged := cloneCleanWitness(*result.Answer.Witness)
		mutate(&forged)
		if err := p.VerifyCleanWitness(world, SignedLiteral{Atom: 2}, forged); err == nil {
			t.Fatal("forged fragment accepted")
		}
	}
	result.Answer.Witness.Absence[0].Present = !result.Answer.Witness.Absence[0].Present
	if !checkCache(t, p, cache, SignedLiteral{Atom: 2}).Reused {
		t.Fatal("caller mutated the retained cache entry")
	}
	before, err := cache.World()
	if err != nil {
		t.Fatal(err)
	}
	if err := cache.Update(HornUpdate{Facts: []HornChange{{0, false}, {999, true}}}); err == nil {
		t.Fatal("invalid update accepted")
	}
	after, err := cache.World()
	if err != nil || !reflect.DeepEqual(before, after) {
		t.Fatal("rejected update partially applied", err)
	}
}

func TestCleanCacheUnusedAlternativeAndCapacity(t *testing.T) {
	rules, facts := cleanFixture()
	p, err := CompileCleanSupport(3, rules)
	if err != nil {
		t.Fatal(err)
	}
	cache, err := NewCleanWitnessCache(p, SignedWorld{1, facts, []bool{true, true}}, CleanCone)
	if err != nil {
		t.Fatal(err)
	}
	checkCache(t, p, cache, SignedLiteral{Atom: 2})
	if err := cache.Update(HornUpdate{Rules: []HornChange{{0, false}}}); err != nil {
		t.Fatal(err)
	}
	if !checkCache(t, p, cache, SignedLiteral{Atom: 2}).Reused {
		t.Fatal("unused positive alternative invalidated witness")
	}
	p, err = CompileCleanSupport(129, nil)
	if err != nil {
		t.Fatal(err)
	}
	world := SignedWorld{Revision: 1, Facts: make([]State, 129)}
	for i := range world.Facts {
		world.Facts[i] = True
	}
	cache, err = NewCleanWitnessCache(p, world, CleanCone)
	if err != nil {
		t.Fatal(err)
	}
	last := CleanCacheResult{}
	for atom := range world.Facts {
		last = checkCache(t, p, cache, SignedLiteral{Atom: atom})
	}
	if !last.Evicted || len(cache.entries) != 128 {
		t.Fatal("cache capacity not enforced")
	}
}

func TestCleanCacheOwnerAndEpochLimits(t *testing.T) {
	if _, err := NewCleanWitnessCache(nil, SignedWorld{}, CleanCone); err == nil {
		t.Fatal("nil owner accepted")
	}
	p, cache := cacheFixture(t, CleanCone, true)
	if _, err := NewCleanWitnessCache(p, SignedWorld{}, CleanCacheMode("other")); err == nil {
		t.Fatal("unknown mode accepted")
	}
	before, err := cache.World()
	if err != nil {
		t.Fatal(err)
	}
	cache.factEpochs[2] = math.MaxUint64
	if err := cache.Update(HornUpdate{Facts: []HornChange{{0, false}, {2, true}}}); err == nil {
		t.Fatal("epoch overflow accepted")
	}
	after, err := cache.World()
	if err != nil || !reflect.DeepEqual(before, after) {
		t.Fatal("overflow partially applied", err)
	}
	cache.factEpochs[2] = 1
	cache.clock = math.MaxUint64
	if _, err := cache.Resolve(SignedLiteral{Atom: 2}); err == nil {
		t.Fatal("cache clock overflow accepted")
	}
}

func FuzzCleanWitnessCache(f *testing.F) {
	f.Add([]byte{0, 3, 5, 8, 2, 1, 0, 6, 9})
	f.Add([]byte{1, 4, 7, 0, 2, 8})
	f.Fuzz(func(t *testing.T, data []byte) {
		p, cache := cacheFixture(t, CleanCone, true)
		for _, raw := range data[:min(len(data), 32)] {
			update := HornUpdate{}
			if raw%3 == 0 {
				update.Rules = []HornChange{{int(raw>>2) % 2, raw&1 != 0}}
			} else {
				update.Facts = []HornChange{{int(raw>>1) % 10, raw&1 != 0}}
			}
			if err := cache.Update(update); err != nil {
				t.Fatal(err)
			}
			checkCache(t, p, cache, SignedLiteral{Atom: 2})
		}
	})
}
