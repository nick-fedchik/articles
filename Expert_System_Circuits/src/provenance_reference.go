package evidence

import (
	"errors"
	"sort"
)

// ReferenceSupport walks the original expression with memoized states and a single root support set.
func ReferenceSupport(root *Expr, sources []Source, inputs []State) (Result, error) {
	if len(sources) != len(inputs) {
		return Result{}, errors.New("support reference: input count mismatch")
	}
	if _, err := Compile(root, sources); err != nil {
		return Result{}, err
	}
	for _, state := range inputs {
		if !state.Valid() {
			return Result{}, errors.New("support reference: invalid input state")
		}
	}
	states := map[*Expr]State{}
	var evaluate func(*Expr) (State, error)
	evaluate = func(expr *Expr) (State, error) {
		if state, exists := states[expr]; exists {
			return state, nil
		}
		var state State
		if expr.Op == Input {
			state = inputs[expr.Input]
		} else {
			a, err := evaluate(expr.Left)
			if err != nil {
				return Unknown, err
			}
			b, err := evaluate(expr.Right)
			if err != nil {
				return Unknown, err
			}
			state, err = ReferencePriority(a, b)
			if err != nil {
				return Unknown, err
			}
			if expr.Op == Select {
				condition, err := evaluate(expr.Condition)
				if err != nil {
					return Unknown, err
				}
				if condition == True {
					state = b
				} else if condition == False {
					state = a
				} else {
					state, err = ReferencePriority(state, Both)
					if err != nil {
						return Unknown, err
					}
				}
			}
		}
		states[expr] = state
		return state, nil
	}
	state, err := evaluate(root)
	if err != nil {
		return Result{}, err
	}
	visited := map[*Expr]bool{}
	type key struct{ id, version string }
	support := map[key]Source{}
	var walk func(*Expr)
	walk = func(expr *Expr) {
		if visited[expr] {
			return
		}
		visited[expr] = true
		if expr.Op == Input {
			source := sources[expr.Input]
			support[key{source.ID, source.Version}] = source
			return
		}
		if expr.Op == All {
			walk(expr.Left)
			walk(expr.Right)
			return
		}
		walk(expr.Condition)
		switch states[expr.Condition] {
		case True:
			walk(expr.Right)
		case False:
			walk(expr.Left)
		default:
			walk(expr.Left)
			walk(expr.Right)
		}
	}
	walk(root)
	keys := []key{}
	for identity := range support {
		keys = append(keys, identity)
	}
	sort.Slice(keys, func(i, j int) bool {
		if keys[i].id != keys[j].id {
			return keys[i].id < keys[j].id
		}
		return keys[i].version < keys[j].version
	})
	result := Result{State: state}
	degradation, zero := 0.0, false
	for _, identity := range keys {
		source := support[identity]
		result.Sources = append(result.Sources, source)
		if source.Weight == 0 {
			zero = true
		} else {
			degradation += 1/source.Weight - 1
		}
	}
	result.Weight = 1 / (1 + degradation)
	if zero {
		result.Weight = 0
	}
	return result, nil
}
