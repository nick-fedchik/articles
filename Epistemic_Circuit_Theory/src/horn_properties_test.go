package evidence

import (
	"slices"
	"sync"
	"testing"
)

func TestHornSearchAllFactAssignments(t *testing.T) {
	p, _ := hornFixture(t)
	for mask := 0; mask < 16; mask++ {
		snapshot := HornSnapshot{1, make([]bool, 4), []bool{true, true}}
		for i := range snapshot.Facts {
			snapshot.Facts[i] = mask&(1<<i) != 0
		}
		want := snapshot.Facts[2] || snapshot.Facts[0] && snapshot.Facts[1] || snapshot.Facts[3]
		for _, method := range HornMethods {
			got, err := p.Search(snapshot, 2, method, 100000)
			if err != nil || got.Exhausted || got.Proved != want {
				t.Fatalf("assignment %d %s: %+v %v", mask, method, got, err)
			}
		}
	}
}

func TestHornProgramSupportsIndependentConcurrentStates(t *testing.T) {
	c, err := GenerateHornCase("conjunction", 64, 1)
	if err != nil {
		t.Fatal(err)
	}
	var workers sync.WaitGroup
	for worker := 0; worker < 4; worker++ {
		workers.Add(1)
		go func(worker int) {
			defer workers.Done()
			state, err := NewHornState(c.Program, c.Snapshot)
			if err != nil {
				t.Error(err)
				return
			}
			if _, err := state.Update(HornUpdate{Facts: []HornChange{{c.Roots[worker], false}}}); err != nil {
				t.Error(err)
				return
			}
			full, err := c.Program.FullClosure(state.Snapshot())
			if err != nil || !slices.Equal(full, state.Known()) {
				t.Error("independent concurrent state disagrees with full recomputation", err)
			}
		}(worker)
	}
	workers.Wait()
}

func FuzzHornCertificate(f *testing.F) {
	f.Add(uint8(3), uint8(2))
	f.Fuzz(func(t *testing.T, facts, query uint8) {
		p, snapshot := hornFixture(t)
		for i := range snapshot.Facts {
			snapshot.Facts[i] = facts&(1<<i) != 0
		}
		q := int(query % 4)
		closure, err := p.FullClosure(snapshot)
		if err != nil {
			t.Fatal(err)
		}
		for _, method := range HornMethods {
			got, err := p.Search(snapshot, q, method, 100000)
			if err != nil || got.Exhausted || got.Proved != closure[q] {
				t.Fatalf("fuzz %s q=%d: %+v %v", method, q, got, err)
			}
		}
	})
}
