package evidence

import (
	"errors"
	"math/bits"
	"unsafe"
)

type ProvenanceFormat string

const (
	ProvenanceDense    ProvenanceFormat = "dense"
	ProvenanceSparse   ProvenanceFormat = "sparse"
	ProvenanceInterned ProvenanceFormat = "interned"
	ProvenanceTrace    ProvenanceFormat = "root-trace"
)

var ErrProvenanceBudget = errors.New("provenance: managed-array memory budget exceeded before allocation")

type ProvenanceStats struct {
	Nodes, Sources, RootSources        int
	VectorIDsAllocated, MaskWords      int64
	InternedVectors, IndexSlots        int
	ManagedWorkingBytes, RetainedBytes int64
}

type ProvenanceSnapshot struct {
	program *Program
	format  ProvenanceFormat
	state   State
	mask    []uint64
	ids     []uint32
	stats   ProvenanceStats
}

type provenanceBudget struct{ used, limit int64 }

func (b *provenanceBudget) take(bytes int64) error {
	if bytes < 0 || bytes > b.limit-b.used {
		return ErrProvenanceBudget
	}
	b.used += bytes
	return nil
}

func validProvenanceFormat(format ProvenanceFormat) bool {
	return format == ProvenanceDense || format == ProvenanceSparse || format == ProvenanceInterned || format == ProvenanceTrace
}

func (p *Program) BuildProvenance(inputs []State, format ProvenanceFormat, maxManagedBytes int64) (*ProvenanceSnapshot, error) {
	if p == nil || len(p.code) == 0 || len(inputs) != len(p.sourceIndex) || !validProvenanceFormat(format) || maxManagedBytes < 0 {
		return nil, errors.New("provenance: malformed program, inputs, format or budget")
	}
	for _, state := range inputs {
		if !state.Valid() {
			return nil, errors.New("provenance: invalid input state")
		}
	}
	budget := &provenanceBudget{limit: maxManagedBytes}
	if err := budget.take(int64(unsafe.Sizeof(ProvenanceSnapshot{})) + int64(len(p.code))); err != nil {
		return nil, err
	}
	snapshot := &ProvenanceSnapshot{program: p, format: format,
		stats: ProvenanceStats{Nodes: len(p.code), Sources: len(p.sources)}}
	states := make([]State, len(p.code))
	words := (len(p.sources) + 63) / 64
	if format == ProvenanceTrace {
		if err := p.buildTrace(inputs, states, snapshot, budget, words); err != nil {
			return nil, err
		}
	} else if format == ProvenanceDense {
		cells := len(p.code) * words
		if err := budget.take(int64(cells) * 8); err != nil {
			return nil, err
		}
		masks := make([]uint64, cells)
		for i, inst := range p.code {
			dst := masks[i*words : (i+1)*words]
			if inst.op == Input {
				states[i] = inputs[inst.input]
				source := p.sourceIndex[inst.input]
				dst[source/64] = uint64(1) << uint(source%64)
				continue
			}
			left, right, state := provenanceBranches(inst, states)
			states[i] = state
			for w := range dst {
				dst[w] = masks[left*words+w] | masks[right*words+w]
				if inst.op == Select {
					dst[w] |= masks[inst.cond*words+w]
				}
			}
		}
		if err := budget.take(int64(words) * 8); err != nil {
			return nil, err
		}
		snapshot.mask = make([]uint64, words)
		copy(snapshot.mask, masks[p.root*words:(p.root+1)*words])
		snapshot.stats.MaskWords = int64(cells + words)
	} else {
		if err := p.buildSparse(inputs, states, snapshot, budget, format == ProvenanceInterned); err != nil {
			return nil, err
		}
	}
	snapshot.state = states[p.root]
	snapshot.stats.ManagedWorkingBytes = budget.used
	snapshot.stats.RetainedBytes = int64(unsafe.Sizeof(ProvenanceSnapshot{})) + int64(cap(snapshot.mask))*8 + int64(cap(snapshot.ids))*4
	if snapshot.mask != nil {
		for _, word := range snapshot.mask {
			snapshot.stats.RootSources += bits.OnesCount64(word)
		}
	} else {
		snapshot.stats.RootSources = len(snapshot.ids)
	}
	return snapshot, nil
}

func provenanceBranches(inst instruction, states []State) (int, int, State) {
	left, right := inst.left, inst.right
	state := priority(states[left], states[right])
	if inst.op == Select {
		switch states[inst.cond] {
		case True:
			left, state = right, states[right]
		case False:
			right, state = left, states[left]
		default:
			state = priority(state, Both)
		}
	}
	return left, right, state
}

func (p *Program) buildTrace(inputs, states []State, snapshot *ProvenanceSnapshot, budget *provenanceBudget, words int) error {
	if err := budget.take(int64(len(p.code))*(1+int64(unsafe.Sizeof(int(0)))) + int64(words)*8); err != nil {
		return err
	}
	for i, inst := range p.code {
		if inst.op == Input {
			states[i] = inputs[inst.input]
		} else {
			_, _, states[i] = provenanceBranches(inst, states)
		}
	}
	seen := make([]bool, len(p.code))
	stack := make([]int, 0, len(p.code))
	push := func(node int) {
		if !seen[node] {
			seen[node] = true
			stack = append(stack, node)
		}
	}
	snapshot.mask = make([]uint64, words)
	push(p.root)
	for len(stack) > 0 {
		node := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		inst := p.code[node]
		if inst.op == Input {
			source := p.sourceIndex[inst.input]
			snapshot.mask[source/64] |= uint64(1) << uint(source%64)
			continue
		}
		left, right, _ := provenanceBranches(inst, states)
		push(left)
		push(right)
		if inst.op == Select {
			push(inst.cond)
		}
	}
	snapshot.stats.MaskWords = int64(words)
	return nil
}

func (s *ProvenanceSnapshot) Stats() ProvenanceStats { return s.stats }

func (s *ProvenanceSnapshot) Result() (Result, error) {
	if s == nil || s.program == nil || !validProvenanceFormat(s.format) {
		return Result{}, errors.New("provenance: uninitialized snapshot")
	}
	out := Result{State: s.state, Sources: make([]Source, 0, s.stats.RootSources)}
	degradation, zero := 0.0, false
	add := func(index int) {
		source := s.program.sources[index]
		out.Sources = append(out.Sources, source)
		if source.Weight == 0 {
			zero = true
		} else {
			degradation += (1 - source.Weight) / source.Weight
		}
	}
	if s.mask != nil {
		for w, word := range s.mask {
			for word != 0 {
				add(64*w + bits.TrailingZeros64(word))
				word &= word - 1
			}
		}
	} else {
		for _, index := range s.ids {
			add(int(index))
		}
	}
	out.Weight = 1 / (1 + degradation)
	if zero {
		out.Weight = 0
	}
	return out, nil
}

func (s *ProvenanceSnapshot) WithSourceWeight(id, version string, weight float64) (*ProvenanceSnapshot, error) {
	if s == nil || s.program == nil {
		return nil, errors.New("provenance: uninitialized snapshot")
	}
	p, err := s.program.WithSourceWeight(id, version, weight)
	if err != nil {
		return nil, err
	}
	out := *s
	out.program = p
	return &out, nil
}
