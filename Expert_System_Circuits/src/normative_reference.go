package evidence

import (
	"errors"
	"reflect"
	"sort"
)

func referenceLiteral(states []State, literal SignedLiteral) (bool, bool) {
	positive, negative := states[literal.Atom]&2 != 0, states[literal.Atom]&1 != 0
	if literal.Negative {
		return negative, positive
	}
	return positive, negative
}

func referenceConjunction(states []State, conditions []SignedLiteral) State {
	positive, negative := true, false
	for _, literal := range conditions {
		pos, neg := referenceLiteral(states, literal)
		positive, negative = positive && pos, negative || neg
	}
	result := Unknown
	if positive {
		result += 2
	}
	if negative {
		result++
	}
	return result
}

func referenceDominates(theory NormTheory, higher, lower string) bool {
	queue, seen := []string{higher}, map[string]bool{}
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		if seen[current] {
			continue
		}
		seen[current] = true
		for _, edge := range theory.Priorities {
			if edge.Higher == current {
				if edge.Lower == lower {
					return true
				}
				queue = append(queue, edge.Lower)
			}
		}
	}
	return false
}

func referenceNormPolicy(theory NormTheory, world NormWorld, action string, closure []State) NormResult {
	result := NormResult{Action: action}
	applies := map[string]NormEffect{}
	hasUnknown, hasConflict := false, false
	for index, rule := range theory.Norms {
		if rule.Action != action {
			continue
		}
		condition := referenceConjunction(closure, rule.Conditions)
		exception := False
		if rule.Exception != nil {
			exception = referenceConjunction(closure, []SignedLiteral{*rule.Exception})
		}
		observation := NormObservation{ID: rule.ID, Condition: condition, Exception: exception, HasException: rule.Exception != nil}
		bit := State(1)
		if rule.Effect == Permit {
			bit = 2
		}
		if world.NormEnabled[index] && condition == True {
			result.ConditionEvidence |= bit
		}
		if !world.NormEnabled[index] {
			observation.Status = "DISABLED"
		} else if condition == False {
			observation.Status = "INACTIVE"
		} else if exception == True {
			observation.Status = "EXCEPTED"
		} else {
			switch condition {
			case Unknown:
				observation.Status = "UNKNOWN_CONDITION"
				hasUnknown = true
			case Both:
				observation.Status = "CONFLICTING_CONDITION"
				hasConflict = true
			case True:
				switch exception {
				case Unknown:
					observation.Status = "UNKNOWN_EXCEPTION"
					hasUnknown = true
				case Both:
					observation.Status = "CONFLICTING_EXCEPTION"
					hasConflict = true
				case False:
					observation.Status = "APPLICABLE"
					applies[rule.ID] = rule.Effect
					result.Applicable |= bit
				}
			}
		}
		result.Rules = append(result.Rules, observation)
	}
	for i := range result.Rules {
		id := result.Rules[i].ID
		effect, exists := applies[id]
		if !exists {
			continue
		}
		for other := range applies {
			if other != id && referenceDominates(theory, other, id) {
				result.Rules[i].DominatedBy = append(result.Rules[i].DominatedBy, other)
			}
		}
		sort.Strings(result.Rules[i].DominatedBy)
		if len(result.Rules[i].DominatedBy) == 0 {
			if effect == Permit {
				result.Surviving |= 2
			} else {
				result.Surviving |= 1
			}
		}
	}
	result.Decision = NormNoBasis
	if hasConflict {
		result.Decision = NormConflict
	} else if hasUnknown {
		result.Decision = NormUnknown
	} else {
		switch result.Surviving {
		case True:
			result.Decision = NormAllow
		case False:
			result.Decision = NormDeny
		case Both:
			result.Decision = NormConflict
		}
	}
	sort.Slice(result.Rules, func(i, j int) bool { return result.Rules[i].ID < result.Rules[j].ID })
	return result
}

func (p *NormProgram) Reference(w NormWorld, action string) ([]State, NormResult, error) {
	if err := p.validateWorld(w); err != nil {
		return nil, NormResult{}, err
	}
	if !p.actions[action] {
		return nil, NormResult{}, errors.New("normative reference: unknown action")
	}
	closure := append([]State(nil), w.Facts...)
	for iteration := 0; iteration <= 2*p.theory.Atoms; iteration++ {
		changed := false
		for i, rule := range p.theory.Evidence {
			if !w.EvidenceEnabled[i] {
				continue
			}
			holds := true
			for _, literal := range rule.Premises {
				positive, _ := referenceLiteral(closure, literal)
				holds = holds && positive
			}
			if holds {
				bit := State(2)
				if rule.Conclusion.Negative {
					bit = 1
				}
				old := closure[rule.Conclusion.Atom]
				closure[rule.Conclusion.Atom] |= bit
				changed = changed || old != closure[rule.Conclusion.Atom]
			}
		}
		if !changed {
			return closure, referenceNormPolicy(p.theory, w, action, closure), nil
		}
	}
	return nil, NormResult{}, errors.New("normative reference: monotone closure iteration limit")
}

func (p *NormProgram) Verify(w NormWorld, action string, certificate NormCertificate) error {
	context, err := p.context(w)
	if err != nil {
		return err
	}
	if !p.actions[action] || certificate.Result.Action != action || certificate.Theory != p.digest ||
		certificate.Context != context || len(certificate.Closure) != p.theory.Atoms {
		return errors.New("normative certificate: context, query, theory or closure size mismatch")
	}
	for _, state := range certificate.Closure {
		if !state.Valid() {
			return errors.New("normative certificate: invalid state")
		}
	}
	producers := make([][]int, 2*p.theory.Atoms)
	for i, rule := range p.theory.Evidence {
		producers[rule.Conclusion.code()] = append(producers[rule.Conclusion.code()], i)
	}
	for code, rules := range producers {
		literal := SignedLiteral{Atom: code / 2, Negative: code%2 != 0}
		expected, _ := referenceLiteral(w.Facts, literal)
		for _, index := range rules {
			if !w.EvidenceEnabled[index] {
				continue
			}
			holds := true
			for _, premise := range p.theory.Evidence[index].Premises {
				positive, _ := referenceLiteral(certificate.Closure, premise)
				holds = holds && positive
			}
			expected = expected || holds
		}
		actual, _ := referenceLiteral(certificate.Closure, literal)
		if actual != expected {
			return errors.New("normative certificate: signed closure equation does not hold")
		}
	}
	expectedResult := referenceNormPolicy(p.theory, w, action, certificate.Closure)
	if !reflect.DeepEqual(expectedResult, certificate.Result) {
		return errors.New("normative certificate: applicability, priority or decision mismatch")
	}
	return nil
}
