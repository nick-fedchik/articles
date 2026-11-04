package evidence

import (
	"errors"
	"fmt"
	"sort"
)

// ReferencePriority deliberately uses a rank specification rather than the runtime LUT.
func ReferencePriority(states ...State) (State, error) {
	rank := map[State]int{True: 0, Both: 1, Unknown: 2, False: 3}
	out, highest := True, 0
	for _, state := range states {
		r, ok := rank[state]
		if !ok {
			return Unknown, errors.New("reference: invalid state")
		}
		if r > highest {
			highest, out = r, state
		}
	}
	return out, nil
}

func ReferenceRule(rule Rule, fact Fact) (State, error) {
	if err := rule.validate(); err != nil {
		return Unknown, err
	}
	if err := fact.validate(); err != nil {
		return Unknown, err
	}
	checks := []State{True, True, True, rule.Modal}
	p := []string{rule.Pattern.Actor, rule.Pattern.Action, rule.Pattern.Target}
	f := []string{fact.Roles.Actor, fact.Roles.Action, fact.Roles.Target}
	for i := range p {
		if p[i] != "*" && p[i] != f[i] {
			checks[0] = Unknown
		}
	}
	if checks[0] != True && p[0] != "*" && p[2] != "*" && p[0] != p[2] &&
		p[0] == f[2] && p[2] == f[0] && (p[1] == "*" || p[1] == f[1]) {
		checks[0] = False
	}
	if rule.Sort != fact.Sort {
		checks[1] = Unknown
	} else {
		contained := rule.Window.Low <= fact.Range.Low && rule.Window.High >= fact.Range.High
		intersects := rule.Window.Low <= fact.Range.High && fact.Range.Low <= rule.Window.High
		switch {
		case contained:
			checks[1] = True
		case intersects:
			checks[1] = Both
		default:
			checks[1] = False
		}
	}
	if fact.Age > rule.MaxAge {
		checks[2] = False
	}
	return ReferencePriority(checks...)
}

func ReferenceEvaluate(root *Expr, sources []Source, inputs []State) (Result, error) {
	if len(sources) != len(inputs) || len(sources) == 0 || len(sources) > MaxGraphNodes {
		return Result{}, errors.New("reference: input count mismatch")
	}
	snapshots := map[string]Source{}
	for i, source := range sources {
		if err := source.validate(); err != nil {
			return Result{}, err
		}
		if !inputs[i].Valid() {
			return Result{}, errors.New("reference: invalid state")
		}
		key := source.ID + "\x00" + source.Version
		if previous, ok := snapshots[key]; ok && previous != source {
			return Result{}, errors.New("reference: conflicting source snapshot")
		}
		snapshots[key] = source
	}
	type value struct {
		state   State
		sources map[string]Source
	}
	visits := 0
	active := map[*Expr]bool{}
	var visit func(*Expr, int) (value, error)
	visit = func(expr *Expr, depth int) (value, error) {
		visits++
		if expr == nil || active[expr] || depth > MaxGraphDepth || visits > 100000 {
			return value{}, errors.New("reference: nil, cycle, depth or visit limit")
		}
		active[expr] = true
		defer delete(active, expr)
		if expr.Op == Input {
			if expr.Input < 0 || expr.Input >= len(inputs) || expr.Left != nil || expr.Right != nil || expr.Condition != nil {
				return value{}, errors.New("reference: malformed input")
			}
			s := sources[expr.Input]
			return value{inputs[expr.Input], map[string]Source{s.ID + "\x00" + s.Version: s}}, nil
		}
		if expr.Op != All && expr.Op != Select || expr.Op == All && expr.Condition != nil {
			return value{}, errors.New("reference: malformed operator")
		}
		left, err := visit(expr.Left, depth+1)
		if err != nil {
			return value{}, err
		}
		right, err := visit(expr.Right, depth+1)
		if err != nil {
			return value{}, err
		}
		parts := []value{left, right}
		if expr.Op == Select {
			cond, err := visit(expr.Condition, depth+1)
			if err != nil {
				return value{}, err
			}
			switch cond.state {
			case True:
				parts = []value{right, cond}
			case False:
				parts = []value{left, cond}
			default:
				parts = []value{left, right, {state: Both, sources: cond.sources}}
			}
		}
		out := value{state: True, sources: map[string]Source{}}
		states := []State{}
		for _, part := range parts {
			// A resolved selector is a control dependency, not an asserted conjunct.
			if expr.Op != Select || len(parts) != 2 || len(states) == 0 {
				states = append(states, part.state)
			}
			for key, source := range part.sources {
				out.sources[key] = source
			}
		}
		out.state, err = ReferencePriority(states...)
		return out, err
	}
	valueOut, err := visit(root, 0)
	if err != nil {
		return Result{}, err
	}
	keys := []string{}
	for key := range valueOut.sources {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	out := Result{State: valueOut.state}
	sum, zero := 0.0, false
	for _, key := range keys {
		source := valueOut.sources[key]
		out.Sources = append(out.Sources, source)
		if source.Weight == 0 {
			zero = true
		} else {
			sum += 1/source.Weight - 1
		}
	}
	out.Weight = 1 / (1 + sum)
	if zero {
		out.Weight = 0
	}
	if len(out.Sources) == 0 {
		return Result{}, fmt.Errorf("reference: empty support")
	}
	return out, nil
}
