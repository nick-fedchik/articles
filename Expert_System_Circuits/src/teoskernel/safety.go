package teoskernel

import (
	"fmt"
	"slices"

	evidence "github.com/nick-fedchik/articles/Expert_System_Circuits/src"
)

// Clause names of the safety property P for ALLOW.
const (
	ClauseWellFormed      = "P0_WELL_FORMED"
	ClausePositive        = "P1_POSITIVE_CLEAN_SUPPORT"
	ClauseProhibition     = "P2_NO_UNDEFEATED_PROHIBITION"
	ClauseDetermined      = "P3_DETERMINED_CONDITIONS_AND_EXCEPTIONS"
	ClauseInterval        = "P4_INTERVAL_SAME_UNIT_WITHIN_BOUNDS"
	ClauseFresh           = "P5_FRESH"
	ClauseVerifiedSources = "P6_VERIFIED_SOURCES"
)

func ruleValues(closure []evidence.State, rule evidence.NormRule) (evidence.State, evidence.State) {
	condition := evidence.True
	for _, literal := range rule.Conditions {
		condition, _ = evidence.TruthAnd(condition, literalValue(closure, literal))
	}
	exception := evidence.False
	if rule.Exception != nil {
		exception = literalValue(closure, *rule.Exception)
	}
	return condition, exception
}

func dominates(theory evidence.NormTheory, higher, lower string) bool {
	stack, seen := []string{higher}, map[string]bool{}
	for len(stack) > 0 {
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if seen[current] {
			continue
		}
		seen[current] = true
		for _, edge := range theory.Priorities {
			if edge.Higher == current {
				if edge.Lower == lower {
					return true
				}
				stack = append(stack, edge.Lower)
			}
		}
	}
	return false
}

// CheckSafety evaluates the safety property P of an ALLOW decision directly on
// (s, in), without any certificate. It returns the violated clauses; an empty
// result means P holds. P is:
//
//	P1 some permit r for the action has every condition literal cleanly
//	   supported (positive and no negative support in the least closure, and
//	   derived without any conflicting BOTH ancestor in Clean), exception
//	   absent or exactly FALSE, and no applicable rule strictly above r;
//	P2 every applicable prohibition (condition TRUE, exception FALSE) is
//	   strictly dominated by an applicable rule (and, since CompileNormative
//	   enforces a strict priority order on finite rules, by an applicable
//	   permit);
//	P3 every rule has condition FALSE, or exception TRUE, or condition TRUE and
//	   exception FALSE: no unknown or conflicting condition or exception;
//	P4 every guard has a measurement, and every measurement of its quantity
//	   has the identical unit and Low >= guard.Low, High <= guard.High;
//	P5 every relevant fact (on an atom referenced by a guard, evidence rule, or
//	   norm) and every measurement time t satisfies Now-MaxAge <= t <= Now;
//	P6 every relevant fact and every measurement has a nonempty source set
//	   included in TrustedSources.
func CheckSafety(s Spec, in Input) []string {
	if _, _, err := compile(s); err != nil {
		return []string{ClauseWellFormed}
	}
	if err := validateInput(s, in); err != nil {
		return []string{ClauseWellFormed}
	}
	return checkSafetyValidated(s, in)
}

func checkSafetyValidated(s Spec, in Input) []string {
	var violated []string
	atoms, _ := referenceAtoms(s, in)
	closure := leastClosure(s.Theory, atoms)
	clean := leastClean(s.Theory, atoms, closure)
	applicable := map[string]bool{}
	for _, rule := range s.Theory.Norms {
		condition, exception := ruleValues(closure, rule)
		applicable[rule.ID] = condition == evidence.True && exception == evidence.False
		if !(condition == evidence.False || exception == evidence.True || applicable[rule.ID]) {
			violated = append(violated, fmt.Sprintf("%s:%s", ClauseDetermined, rule.ID))
		}
	}
	defeated := func(id string, byPermitOnly bool) bool {
		for _, h := range s.Theory.Norms {
			if applicable[h.ID] && (!byPermitOnly || h.Effect == evidence.Permit) && dominates(s.Theory, h.ID, id) {
				return true
			}
		}
		return false
	}
	witness := false
	for _, rule := range s.Theory.Norms {
		if rule.Effect == evidence.Permit && applicable[rule.ID] && !defeated(rule.ID, false) {
			allClean := true
			for _, l := range rule.Conditions {
				positive, _ := literalBits(clean, l)
				allClean = allClean && positive
			}
			if allClean {
				witness = true
			}
		}
		if rule.Effect == evidence.Prohibit && applicable[rule.ID] && (!defeated(rule.ID, false) || !defeated(rule.ID, true)) {
			violated = append(violated, fmt.Sprintf("%s:%s", ClauseProhibition, rule.ID))
		}
	}
	if !witness {
		violated = append(violated, ClausePositive)
	}
	for gi, g := range s.Guards {
		count := 0
		for mi, m := range in.Measurements {
			if m.Quantity == g.Quantity {
				count++
				if m.Unit != g.Unit || m.Low < g.Low || m.High > g.High {
					violated = append(violated, fmt.Sprintf("%s:measurement:%d", ClauseInterval, mi))
				}
			}
		}
		if count == 0 {
			violated = append(violated, fmt.Sprintf("%s:guard:%d", ClauseInterval, gi))
		}
	}
	check := func(kind string, index int, sources []string, t int64) {
		if !freshBig(in.Now, t, s.MaxAge) {
			violated = append(violated, fmt.Sprintf("%s:%s:%d", ClauseFresh, kind, index))
		}
		if len(sources) == 0 || slices.ContainsFunc(sources, func(id string) bool { return !slices.Contains(s.TrustedSources, id) }) {
			violated = append(violated, fmt.Sprintf("%s:%s:%d", ClauseVerifiedSources, kind, index))
		}
	}
	for i, f := range in.Facts {
		if relevantAtom(s, f.Literal.Atom) {
			check("fact", i, f.Sources, f.Time)
		}
	}
	for i, m := range in.Measurements {
		check("measurement", i, m.Sources, m.Time)
	}
	return violated
}
