package evidence

import (
	"context"
	"errors"
)

type NormSATAnswer struct {
	Closure                       []State
	ConditionEvidence, Applicable State
	Surviving                     State
	Decision                      NormDecision
}

func (c *cnf) normOr(values ...int) int {
	out := c.constant(false)
	for _, value := range values {
		out = c.or(out, value)
	}
	return out
}

func (c *cnf) normAnd(values ...int) int {
	out := c.constant(true)
	for _, value := range values {
		out = c.and(out, value)
	}
	return out
}

func (p *NormProgram) SATReference(ctx context.Context, solver string, world NormWorld, action string) (NormSATAnswer, error) {
	if err := p.validateWorld(world); err != nil {
		return NormSATAnswer{}, err
	}
	if !p.actions[action] || p.theory.Atoms > 16 || len(p.theory.Evidence) > 128 || len(p.theory.Norms) > 32 {
		return NormSATAnswer{}, errors.New("normative SAT: unknown action or bounded oracle size exceeded")
	}
	c := &cnf{}
	literals := make([]int, 2*p.theory.Atoms)
	for i := range literals {
		literals[i] = c.variable()
	}
	bodies := make([]int, len(p.theory.Evidence))
	for i, rule := range p.theory.Evidence {
		body := c.constant(world.EvidenceEnabled[i])
		for _, premise := range rule.Premises {
			body = c.and(body, literals[premise.code()])
		}
		bodies[i] = body
	}
	for code, literal := range literals {
		mask := State(2)
		if code%2 != 0 {
			mask = 1
		}
		right := c.constant(world.Facts[code/2]&mask != 0)
		for i, rule := range p.theory.Evidence {
			if rule.Conclusion.code() == code {
				right = c.or(right, bodies[i])
			}
		}
		c.add(-literal, right)
		c.add(literal, -right)
	}
	state := func(literal SignedLiteral) symbolicState {
		result := symbolicState{literals[2*literal.Atom], literals[2*literal.Atom+1]}
		if literal.Negative {
			result.high, result.low = result.low, result.high
		}
		return result
	}
	active := make([]int, len(p.theory.Norms))
	rawPositive, rawNegative, applicablePositive, applicableNegative := []int{}, []int{}, []int{}, []int{}
	unknowns, conflicts := []int{}, []int{}
	for i, rule := range p.theory.Norms {
		active[i] = c.constant(false)
		if rule.Action != action {
			continue
		}
		condition := symbolicState{c.constant(true), c.constant(false)}
		for _, premise := range rule.Conditions {
			part := state(premise)
			condition.high = c.and(condition.high, part.high)
			condition.low = c.or(condition.low, part.low)
		}
		exception := symbolicState{c.constant(false), c.constant(true)}
		if rule.Exception != nil {
			exception = state(*rule.Exception)
		}
		enabled := c.constant(world.NormEnabled[i])
		conditionTrue, conditionFalse := c.is(condition, True), c.is(condition, False)
		exceptionTrue, exceptionFalse := c.is(exception, True), c.is(exception, False)
		relevant := c.normAnd(enabled, -conditionFalse, -exceptionTrue)
		unknowns = append(unknowns, c.and(relevant, c.normOr(c.is(condition, Unknown), c.and(conditionTrue, c.is(exception, Unknown)))))
		conflicts = append(conflicts, c.and(relevant, c.normOr(c.is(condition, Both), c.and(conditionTrue, c.is(exception, Both)))))
		raw := c.and(enabled, conditionTrue)
		active[i] = c.normAnd(enabled, conditionTrue, exceptionFalse)
		if rule.Effect == Permit {
			rawPositive, applicablePositive = append(rawPositive, raw), append(applicablePositive, active[i])
		} else {
			rawNegative, applicableNegative = append(rawNegative, raw), append(applicableNegative, active[i])
		}
	}
	positive, negative := []int{}, []int{}
	for i, rule := range p.theory.Norms {
		if rule.Action != action {
			continue
		}
		dominators := []int{}
		for j, other := range p.theory.Norms {
			if other.Action == action && referenceDominates(p.theory, other.ID, rule.ID) {
				dominators = append(dominators, active[j])
			}
		}
		survives := c.and(active[i], -c.normOr(dominators...))
		if rule.Effect == Permit {
			positive = append(positive, survives)
		} else {
			negative = append(negative, survives)
		}
	}
	pos, neg := c.normOr(positive...), c.normOr(negative...)
	unknown, conflict := c.normOr(unknowns...), c.normOr(conflicts...)
	unblocked := c.and(-unknown, -conflict)
	decisions := map[NormDecision]int{
		NormAllow:    c.normAnd(unblocked, pos, -neg),
		NormDeny:     c.normAnd(unblocked, -pos, neg),
		NormUnknown:  c.and(-conflict, unknown),
		NormConflict: c.or(conflict, c.normAnd(unblocked, pos, neg)),
		NormNoBasis:  c.normAnd(unblocked, -pos, -neg),
	}
	rawP, rawN := c.normOr(rawPositive...), c.normOr(rawNegative...)
	appP, appN := c.normOr(applicablePositive...), c.normOr(applicableNegative...)
	sat, values, err := c.solve(ctx, solver)
	if err != nil {
		return NormSATAnswer{}, err
	}
	if !sat {
		return NormSATAnswer{}, errors.New("normative SAT: unexpected UNSAT; signed contradictions are permitted")
	}
	combine := func(high, low int) State {
		result := Unknown
		if values[high] {
			result |= True
		}
		if values[low] {
			result |= False
		}
		return result
	}
	answer := NormSATAnswer{Closure: make([]State, p.theory.Atoms),
		ConditionEvidence: combine(rawP, rawN), Applicable: combine(appP, appN), Surviving: combine(pos, neg)}
	for i := range answer.Closure {
		answer.Closure[i] = combine(literals[2*i], literals[2*i+1])
	}
	selected := 0
	for decision, literal := range decisions {
		if values[literal] {
			answer.Decision = decision
			selected++
		}
	}
	if selected != 1 {
		return NormSATAnswer{}, errors.New("normative SAT: verdict formulas are not mutually exclusive and exhaustive")
	}
	return answer, nil
}
