package evidence

import (
	"errors"
	"fmt"
	"time"
)

const priorityLUT uint32 = 0xF4E45504

func lookupPriority(a, b State) State {
	return State((priorityLUT >> (2 * (4*uint(a) + uint(b)))) & 3)
}

type RuleBank struct {
	actors, actions, targets, sorts []string
	lower, upper                    []float64
	maxAge                          []time.Duration
	modal                           []State
}

func CompileRules(rules []Rule) (*RuleBank, error) {
	if len(rules) == 0 {
		return nil, errors.New("rule bank: empty")
	}
	bank := &RuleBank{}
	for i, rule := range rules {
		if err := rule.validate(); err != nil {
			return nil, fmt.Errorf("rule %d: %w", i, err)
		}
		bank.actors = append(bank.actors, rule.Pattern.Actor)
		bank.actions = append(bank.actions, rule.Pattern.Action)
		bank.targets = append(bank.targets, rule.Pattern.Target)
		bank.sorts = append(bank.sorts, rule.Sort)
		bank.lower = append(bank.lower, rule.Window.Low)
		bank.upper = append(bank.upper, rule.Window.High)
		bank.maxAge = append(bank.maxAge, rule.MaxAge)
		bank.modal = append(bank.modal, rule.Modal)
	}
	return bank, nil
}

func (bank *RuleBank) Evaluate(fact Fact, dst []State) error {
	if bank == nil || len(bank.modal) == 0 {
		return errors.New("rule bank: uninitialized")
	}
	if err := fact.validate(); err != nil {
		return err
	}
	if len(dst) != len(bank.modal) {
		return errors.New("rule bank: output length mismatch")
	}
	for i := range dst {
		pattern := Roles{bank.actors[i], bank.actions[i], bank.targets[i]}
		bind := bindRoles(pattern, fact.Roles)
		bound := True
		switch {
		case bank.sorts[i] != fact.Sort:
			bound = Unknown
		case fact.Range.High < bank.lower[i] || fact.Range.Low > bank.upper[i]:
			bound = False
		case fact.Range.Low < bank.lower[i] || fact.Range.High > bank.upper[i]:
			bound = Both
		}
		fresh := True
		if fact.Age > bank.maxAge[i] {
			fresh = False
		}
		dst[i] = lookupPriority(lookupPriority(bind, bound), lookupPriority(fresh, bank.modal[i]))
	}
	return nil
}

func evaluateValidatedRule(rule Rule, fact Fact) State {
	bind := bindRoles(rule.Pattern, fact.Roles)
	bound := True
	if rule.Sort != fact.Sort {
		bound = Unknown
	} else if fact.Range.High < rule.Window.Low || fact.Range.Low > rule.Window.High {
		bound = False
	} else if fact.Range.Low < rule.Window.Low || fact.Range.High > rule.Window.High {
		bound = Both
	}
	fresh := True
	if fact.Age > rule.MaxAge {
		fresh = False
	}
	return priority(priority(bind, bound), priority(fresh, rule.Modal))
}
