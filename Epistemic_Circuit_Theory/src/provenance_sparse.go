package evidence

import (
	"errors"
	"slices"
	"unsafe"
)

type provenancePool struct {
	budget  *provenanceBudget
	entries [][]uint32
	slots   []uint32
	ids     int64
}

func newProvenancePool(budget *provenanceBudget, nodes int, intern bool) (*provenancePool, error) {
	pool := &provenancePool{budget: budget}
	if intern {
		capacity := 1
		for capacity < 4*nodes {
			capacity *= 2
		}
		if err := budget.take(int64(2*nodes)*int64(unsafe.Sizeof([]uint32{})) + int64(capacity)*4); err != nil {
			return nil, err
		}
		pool.entries = make([][]uint32, 0, 2*nodes)
		pool.slots = make([]uint32, capacity)
	}
	return pool, nil
}

func unionWalk(a, b []uint32, visit func(uint32)) {
	i, j := 0, 0
	for i < len(a) || j < len(b) {
		var value uint32
		if j == len(b) || i < len(a) && a[i] < b[j] {
			value = a[i]
			i++
		} else if i == len(a) || b[j] < a[i] {
			value = b[j]
			j++
		} else {
			value = a[i]
			i++
			j++
		}
		visit(value)
	}
}

func (p *provenancePool) union(a, b []uint32) ([]uint32, error) {
	if len(a) == 0 {
		return b, nil
	}
	if len(b) == 0 || len(a) == len(b) && &a[0] == &b[0] {
		return a, nil
	}
	if len(a) == 1 {
		if _, exists := slices.BinarySearch(b, a[0]); exists {
			return b, nil
		}
	}
	if len(b) == 1 {
		if _, exists := slices.BinarySearch(a, b[0]); exists {
			return a, nil
		}
	}
	count, hash := 0, uint64(14695981039346656037)
	hashing := len(p.slots) > 0
	unionWalk(a, b, func(value uint32) {
		count++
		if hashing {
			// Unsigned overflow is the defined FNV-1a hash arithmetic.
			for shift := uint(0); shift < 32; shift += 8 {
				hash ^= uint64(byte(value >> shift))
				hash *= 1099511628211
			}
		}
	})
	if count == len(a) {
		return a, nil
	}
	if count == len(b) {
		return b, nil
	}
	slot := 0
	if len(p.slots) > 0 {
		slot = int(hash & uint64(len(p.slots)-1))
		for p.slots[slot] != 0 {
			candidate := p.entries[p.slots[slot]-1]
			equal := len(candidate) == count
			at := 0
			if equal {
				unionWalk(a, b, func(value uint32) {
					equal = equal && candidate[at] == value
					at++
				})
			}
			if equal {
				return candidate, nil
			}
			slot = (slot + 1) & (len(p.slots) - 1)
		}
		if len(p.entries) == cap(p.entries) {
			return nil, errors.New("provenance: explicit interning capacity exhausted")
		}
	}
	if err := p.budget.take(int64(count) * 4); err != nil {
		return nil, err
	}
	out := make([]uint32, count)
	at := 0
	unionWalk(a, b, func(value uint32) { out[at], at = value, at+1 })
	p.ids += int64(count)
	if len(p.slots) > 0 {
		p.entries = append(p.entries, out)
		p.slots[slot] = uint32(len(p.entries))
	}
	return out, nil
}

func (p *Program) buildSparse(inputs, states []State, snapshot *ProvenanceSnapshot, budget *provenanceBudget, intern bool) error {
	if err := budget.take(int64(len(p.code))*int64(unsafe.Sizeof([]uint32{})) + int64(len(p.sources))*4); err != nil {
		return err
	}
	sets := make([][]uint32, len(p.code))
	singletons := make([]uint32, len(p.sources))
	for i := range singletons {
		singletons[i] = uint32(i)
	}
	pool, err := newProvenancePool(budget, len(p.code), intern)
	if err != nil {
		return err
	}
	for i, inst := range p.code {
		if inst.op == Input {
			states[i] = inputs[inst.input]
			source := p.sourceIndex[inst.input]
			sets[i] = singletons[source : source+1 : source+1]
			continue
		}
		left, right, state := provenanceBranches(inst, states)
		states[i] = state
		sets[i], err = pool.union(sets[left], sets[right])
		if err != nil {
			return err
		}
		if inst.op == Select {
			sets[i], err = pool.union(sets[i], sets[inst.cond])
			if err != nil {
				return err
			}
		}
	}
	if err := budget.take(int64(len(sets[p.root])) * 4); err != nil {
		return err
	}
	snapshot.ids = make([]uint32, len(sets[p.root]))
	copy(snapshot.ids, sets[p.root])
	snapshot.stats.VectorIDsAllocated = int64(len(singletons)+len(snapshot.ids)) + pool.ids
	snapshot.stats.InternedVectors, snapshot.stats.IndexSlots = len(pool.entries), len(pool.slots)
	return nil
}
