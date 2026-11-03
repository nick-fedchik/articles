package teoskernel

import (
	"errors"
	"fmt"
	"math/big"
	"slices"
	"time"

	evidence "github.com/nick-fedchik/articles/Epistemic_Circuit_Theory/src"
)

// The functions in this file deliberately avoid the Evaluate code path: the
// guard comparison, freshness arithmetic, evidence closure and verdict mapping
// are recomputed from the definitions of K.

func referenceAtoms(s Spec, in Input) ([]evidence.State, []evidence.State) {
	atoms := make([]evidence.State, s.Theory.Atoms)
	for atom := range atoms {
		positive, negative := false, false
		for _, f := range in.Facts {
			if f.Literal.Atom == atom {
				positive = positive || !f.Literal.Negative
				negative = negative || f.Literal.Negative
			}
		}
		atoms[atom] = bits(positive, negative)
	}
	guards := make([]evidence.State, len(s.Guards))
	for gi, g := range s.Guards {
		for _, m := range in.Measurements {
			if m.Quantity == g.Quantity {
				guards[gi] |= guardState(g, m)
			}
		}
		atoms[g.Atom] = guards[gi]
	}
	return atoms, guards
}

func bits(positive, negative bool) evidence.State {
	var s evidence.State
	if positive {
		s |= evidence.True
	}
	if negative {
		s |= evidence.False
	}
	return s
}

// guardState restates the teoscompose.IntervalGuard convention directly.
func guardState(g Guard, m Measurement) evidence.State {
	switch {
	case m.Unit != g.Unit:
		return evidence.Unknown
	case m.High < g.Low || m.Low > g.High:
		return evidence.False
	case m.Low < g.Low || m.High > g.High:
		return evidence.Both
	default:
		return evidence.True
	}
}

func literalBits(states []evidence.State, l evidence.SignedLiteral) (bool, bool) {
	positive, negative := states[l.Atom]&evidence.True != 0, states[l.Atom]&evidence.False != 0
	if l.Negative {
		return negative, positive
	}
	return positive, negative
}

func literalValue(states []evidence.State, l evidence.SignedLiteral) evidence.State {
	return bits(literalBits(states, l))
}

// leastClosure is the naive least fixpoint of the signed evidence rules.
func leastClosure(theory evidence.NormTheory, facts []evidence.State) []evidence.State {
	closure := slices.Clone(facts)
	for changed := true; changed; {
		changed = false
		for _, rule := range theory.Evidence {
			holds := true
			for _, premise := range rule.Premises {
				positive, _ := literalBits(closure, premise)
				holds = holds && positive
			}
			if !holds {
				continue
			}
			bit := evidence.True
			if rule.Conclusion.Negative {
				bit = evidence.False
			}
			if closure[rule.Conclusion.Atom]&bit == 0 {
				closure[rule.Conclusion.Atom] |= bit
				changed = true
			}
		}
	}
	return closure
}

func freshBig(now, t int64, maxAge time.Duration) bool {
	n, ts := big.NewInt(now), big.NewInt(t)
	lower := new(big.Int).Sub(n, big.NewInt(int64(maxAge)))
	return lower.Cmp(ts) <= 0 && ts.Cmp(n) <= 0
}

type postconditions struct {
	fresh, verified, contained, outOfBounds, boundary bool
	violations                                        []string
}

func referencePostconditions(s Spec, in Input) postconditions {
	p := postconditions{fresh: true, verified: true, contained: true}
	check := func(kind string, index int, sources []string, t int64) {
		if !freshBig(in.Now, t, s.MaxAge) {
			p.fresh = false
			p.violations = append(p.violations, fmt.Sprintf("%s:%s:%d", codeStale, kind, index))
		}
		if len(sources) == 0 || slices.ContainsFunc(sources, func(id string) bool { return !slices.Contains(s.TrustedSources, id) }) {
			p.verified = false
			p.violations = append(p.violations, fmt.Sprintf("%s:%s:%d", codeUnverified, kind, index))
		}
	}
	for i, f := range in.Facts {
		check("fact", i, f.Sources, f.Time)
	}
	for i, m := range in.Measurements {
		check("measurement", i, m.Sources, m.Time)
	}
	for gi, g := range s.Guards {
		count := 0
		for mi, m := range in.Measurements {
			if m.Quantity != g.Quantity {
				continue
			}
			count++
			code := map[evidence.State]string{evidence.Unknown: codeUnit, evidence.False: codeOutOfBound, evidence.Both: codeBoundary}[guardState(g, m)]
			if code != "" {
				p.contained = false
				p.outOfBounds = p.outOfBounds || code == codeOutOfBound
				p.boundary = p.boundary || code == codeBoundary
				p.violations = append(p.violations, fmt.Sprintf("%s:measurement:%d", code, mi))
			}
		}
		if count == 0 {
			p.contained = false
			p.violations = append(p.violations, fmt.Sprintf("%s:guard:%d", codeMissing, gi))
		}
	}
	return p
}

// Verify accepts c only if every recorded intermediate result equals the
// independently recomputed semantics of K for (s, in).
func Verify(s Spec, in Input, c Certificate) error {
	program, err := compile(s)
	if err != nil {
		return err
	}
	if err := validateInput(s, in); err != nil {
		return err
	}
	specDigest, err := digest(s)
	if err != nil {
		return err
	}
	inputDigest, err := digest(in)
	if err != nil {
		return err
	}
	if c.SpecDigest != specDigest || c.InputDigest != inputDigest {
		return errors.New("teoskernel certificate: spec or input digest mismatch")
	}
	atoms, guards := referenceAtoms(s, in)
	if !slices.Equal(atoms, c.Atoms) || !slices.Equal(guards, c.Guards) {
		return errors.New("teoskernel certificate: evidence or guard assignment mismatch")
	}
	world := evidence.NormWorld{Revision: 1, Facts: atoms,
		EvidenceEnabled: make([]bool, len(s.Theory.Evidence)), NormEnabled: make([]bool, len(s.Theory.Norms))}
	for i := range world.EvidenceEnabled {
		world.EvidenceEnabled[i] = true
	}
	for i := range world.NormEnabled {
		world.NormEnabled[i] = true
	}
	if err := program.Verify(world, s.Action, c.Norm); err != nil {
		return fmt.Errorf("teoskernel certificate: %w", err)
	}
	if !slices.Equal(leastClosure(s.Theory, atoms), c.Norm.Closure) {
		return errors.New("teoskernel certificate: closure is not the least fixpoint")
	}
	p := referencePostconditions(s, in)
	if p.fresh != c.Fresh || p.verified != c.Verified || p.contained != c.Contained || !slices.Equal(p.violations, c.Violations) {
		return errors.New("teoskernel certificate: admission postcondition mismatch")
	}
	verdict := c.Norm.Result.Decision
	if verdict == Allow && p.outOfBounds {
		verdict = Deny
	} else if verdict == Allow && p.boundary {
		verdict = Conflict
	} else if verdict == Allow && (!p.fresh || !p.verified || !p.contained) {
		verdict = Unknown
	}
	if verdict != c.Verdict {
		return errors.New("teoskernel certificate: verdict mismatch")
	}
	return nil
}
