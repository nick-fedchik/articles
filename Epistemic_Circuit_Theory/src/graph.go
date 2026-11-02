package evidence

import (
	"errors"
	"fmt"
	"math"
	"math/bits"
)

type Op uint8

const (
	Input Op = iota
	All
	Select
)

type Expr struct {
	Op                     Op
	Input                  int
	Left, Right, Condition *Expr
}

const (
	MaxGraphNodes = 10000
	MaxGraphDepth = 256
)

type instruction struct {
	op                       Op
	input, left, right, cond int
}

type Program struct {
	code        []instruction
	root        int
	sources     []Source
	sourceIndex []int
}

type Result struct {
	State   State
	Weight  float64
	Sources []Source
}

func sourceTable(sources []Source) ([]Source, []int, error) {
	if len(sources) == 0 || len(sources) > MaxGraphNodes {
		return nil, nil, errors.New("sources: count outside [1,10000]")
	}
	unique := []Source{}
	indices := make([]int, len(sources))
	seen := map[string]int{}
	for i, source := range sources {
		if err := source.validate(); err != nil {
			return nil, nil, fmt.Errorf("source %d: %w", i, err)
		}
		key := source.ID + "\x00" + source.Version
		index, exists := seen[key]
		if exists {
			if unique[index] != source {
				return nil, nil, fmt.Errorf("source %q version %q: conflicting snapshot", source.ID, source.Version)
			}
		} else {
			index = len(unique)
			seen[key] = index
			unique = append(unique, source)
		}
		indices[i] = index
	}
	return unique, indices, nil
}

func Compile(root *Expr, sources []Source) (*Program, error) {
	unique, indices, err := sourceTable(sources)
	if err != nil {
		return nil, err
	}
	p := &Program{sources: unique, sourceIndex: indices}
	active := map[*Expr]bool{}
	compiled := map[*Expr]int{}
	canonical := map[instruction]int{}
	var visit func(*Expr, int) (int, error)
	visit = func(expr *Expr, depth int) (int, error) {
		if expr == nil || depth > MaxGraphDepth {
			return 0, errors.New("compile: nil expression or depth limit")
		}
		if active[expr] {
			return 0, errors.New("compile: cyclic expression")
		}
		if index, exists := compiled[expr]; exists {
			return index, nil
		}
		if len(compiled)+len(active) >= MaxGraphNodes {
			return 0, errors.New("compile: node limit")
		}
		active[expr] = true
		inst := instruction{op: expr.Op}
		switch expr.Op {
		case Input:
			if expr.Input < 0 || expr.Input >= len(sources) || expr.Left != nil || expr.Right != nil || expr.Condition != nil {
				return 0, errors.New("compile: malformed input")
			}
			inst.input = expr.Input
		case All, Select:
			if expr.Op == All && expr.Condition != nil {
				return 0, errors.New("compile: condition on conjunction")
			}
			var err error
			inst.left, err = visit(expr.Left, depth+1)
			if err != nil {
				return 0, err
			}
			inst.right, err = visit(expr.Right, depth+1)
			if err != nil {
				return 0, err
			}
			if expr.Op == Select {
				inst.cond, err = visit(expr.Condition, depth+1)
				if err != nil {
					return 0, err
				}
			}
		default:
			return 0, errors.New("compile: unknown operator")
		}
		delete(active, expr)
		index, exists := canonical[inst]
		if !exists {
			index = len(p.code)
			canonical[inst] = index
			p.code = append(p.code, inst)
		}
		compiled[expr] = index
		return index, nil
	}
	p.root, err = visit(root, 0)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (p *Program) Nodes() int { return len(p.code) }

func (p *Program) Evaluate(inputs []State) (Result, error) {
	if p == nil || len(p.code) == 0 {
		return Result{}, errors.New("evaluate: uninitialized program")
	}
	if len(inputs) != len(p.sourceIndex) {
		return Result{}, errors.New("evaluate: input count mismatch")
	}
	for _, state := range inputs {
		if !state.Valid() {
			return Result{}, errors.New("evaluate: invalid input state")
		}
	}
	words := (len(p.sources) + 63) / 64
	masks := make([]uint64, len(p.code)*words)
	states := make([]State, len(p.code))
	for i, inst := range p.code {
		dst := masks[i*words : (i+1)*words]
		if inst.op == Input {
			states[i] = inputs[inst.input]
			source := p.sourceIndex[inst.input]
			dst[source/64] = uint64(1) << uint(source%64)
			continue
		}
		left, right := inst.left, inst.right
		states[i] = priority(states[left], states[right])
		if inst.op == Select {
			switch states[inst.cond] {
			case True:
				left = right
				states[i] = states[right]
			case False:
				right = left
				states[i] = states[left]
			default:
				states[i] = priority(states[i], Both)
			}
		}
		for w := range dst {
			dst[w] = masks[left*words+w] | masks[right*words+w]
			if inst.op == Select {
				dst[w] |= masks[inst.cond*words+w]
			}
		}
	}
	out := Result{State: states[p.root]}
	var degradation float64
	zero := false
	for w, mask := range masks[p.root*words : (p.root+1)*words] {
		for mask != 0 {
			k := 64*w + bits.TrailingZeros64(mask)
			source := p.sources[k]
			out.Sources = append(out.Sources, source)
			if source.Weight == 0 {
				zero = true
			} else {
				degradation += (1 - source.Weight) / source.Weight
			}
			mask &= mask - 1
		}
	}
	out.Weight = 1 / (1 + degradation)
	if zero {
		out.Weight = 0
	}
	if math.IsInf(degradation, 0) {
		return Result{}, errors.New("evaluate: degradation overflow")
	}
	return out, nil
}
