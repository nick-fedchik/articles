package kernel

import (
	"errors"
	"fmt"
	"math/big"
	"slices"
	"time"

	evidence "github.com/nick-fedchik/articles/Expert_System_Circuits/src"
)

// The functions in this file deliberately avoid the Evaluate code path: the
// guard comparison, freshness arithmetic, evidence closure, clean-support
// fixpoint and verdict mapping are recomputed from the definitions of K.
//
// Shared TCB note (F-10): Verify and Evaluate share the static theory compiler
// compile(s) (which validates spec invariants and invokes
// evidence.CompileNormative) and NormProgram.Verify (which checks the normative
// policy observations and SHA-256 digests). All dynamic semantics—pre-closure
// atom and guard assignments, the least signed closure, the clean-support
// fixpoint Clean (Theorem B.7), freshness via arbitrary-precision big.Int
// arithmetic, provenance, interval containment, and the final verdict
// downgrade—are recomputed independently below without calling Evaluate or
// CleanProgram.EvaluateClean.

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

// guardState restates the compose.IntervalGuard convention directly.
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

// leastClean is the naive least fixpoint of clean support (Theorem B.7): a
// signed literal l is cleanly supported iff literalValue(closure, l) == True
// (purely supported in closure with no opposite support) and either l is
// directly supported in facts or derived by an evidence rule whose premises are
// all cleanly supported.
func leastClean(theory evidence.NormTheory, facts, closure []evidence.State) []evidence.State {
	clean := make([]evidence.State, theory.Atoms)
	for a := 0; a < theory.Atoms; a++ {
		if closure[a] == evidence.True && facts[a]&evidence.True != 0 {
			clean[a] |= evidence.True
		}
		if closure[a] == evidence.False && facts[a]&evidence.False != 0 {
			clean[a] |= evidence.False
		}
	}
	for changed := true; changed; {
		changed = false
		for _, rule := range theory.Evidence {
			if literalValue(closure, rule.Conclusion) != evidence.True {
				continue
			}
			holds := true
			for _, premise := range rule.Premises {
				positive, _ := literalBits(clean, premise)
				holds = holds && positive
			}
			if !holds {
				continue
			}
			bit := evidence.True
			if rule.Conclusion.Negative {
				bit = evidence.False
			}
			if clean[rule.Conclusion.Atom]&bit == 0 {
				clean[rule.Conclusion.Atom] |= bit
				changed = true
			}
		}
	}
	return clean
}

func referenceCleanWitness(s Spec, closure, clean []evidence.State) bool {
	applicable := map[string]bool{}
	for _, rule := range s.Theory.Norms {
		condition, exception := ruleValues(closure, rule)
		applicable[rule.ID] = condition == evidence.True && exception == evidence.False
	}
	for _, rule := range s.Theory.Norms {
		if rule.Effect != evidence.Permit || !applicable[rule.ID] {
			continue
		}
		defeated := false
		for _, h := range s.Theory.Norms {
			if applicable[h.ID] && dominates(s.Theory, h.ID, rule.ID) {
				defeated = true
				break
			}
		}
		if defeated {
			continue
		}
		allClean := true
		for _, l := range rule.Conditions {
			positive, _ := literalBits(clean, l)
			allClean = allClean && positive
		}
		if allClean {
			return true
		}
	}
	return false
}

func freshBig(now, t int64, maxAge time.Duration) bool {
	n, ts := big.NewInt(now), big.NewInt(t)
	lower := new(big.Int).Sub(n, big.NewInt(int64(maxAge)))
	return lower.Cmp(ts) <= 0 && ts.Cmp(n) <= 0
}

type postconditions struct {
	fresh, verified, contained, outOfBounds, boundary, cleanWitness bool
	violations                                                      []string
}

func referencePostconditions(s Spec, in Input, closure, clean []evidence.State, normDecision evidence.NormDecision) postconditions {
	p := postconditions{fresh: true, verified: true, contained: true, cleanWitness: true}
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
		if relevantAtom(s, f.Literal.Atom) {
			check("fact", i, f.Sources, f.Time)
		}
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
	if normDecision == Allow && !referenceCleanWitness(s, closure, clean) {
		p.cleanWitness = false
		p.violations = append(p.violations, codeConflictDependent)
	}
	return p
}

// Verify accepts c only if every recorded intermediate result equals the
// independently recomputed semantics of K for (s, in).
func Verify(s Spec, in Input, c Certificate) error {
	program, _, err := compile(s)
	if err != nil {
		return err
	}
	specDigest, err := digest(s)
	if err != nil {
		return err
	}
	return verifyCompiled(s, in, c, program, specDigest)
}

func verifyCompiled(s Spec, in Input, c Certificate, program *evidence.NormProgram, specDigest [32]byte) error {
	if err := validateInput(s, in); err != nil {
		return err
	}
	inputDigest, err := digest(in)
	if err != nil {
		return err
	}
	if c.SpecDigest != specDigest || c.InputDigest != inputDigest {
		return errors.New("kernel certificate: spec or input digest mismatch")
	}
	atoms, guards := referenceAtoms(s, in)
	if !slices.Equal(atoms, c.Atoms) || !slices.Equal(guards, c.Guards) {
		return errors.New("kernel certificate: evidence or guard assignment mismatch")
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
		return fmt.Errorf("kernel certificate: %w", err)
	}
	closure := leastClosure(s.Theory, atoms)
	if !slices.Equal(closure, c.Norm.Closure) {
		return errors.New("kernel certificate: closure is not the least fixpoint")
	}
	clean := leastClean(s.Theory, atoms, closure)
	if !slices.Equal(clean, c.Clean) {
		return errors.New("kernel certificate: clean support is not the least fixpoint")
	}
	p := referencePostconditions(s, in, closure, clean, c.Norm.Result.Decision)
	if p.fresh != c.Fresh || p.verified != c.Verified || p.contained != c.Contained || !slices.Equal(p.violations, c.Violations) {
		return errors.New("kernel certificate: admission postcondition mismatch")
	}
	verdict := c.Norm.Result.Decision
	if verdict == Allow && p.outOfBounds {
		verdict = Deny
	} else if verdict == Allow && (p.boundary || !p.cleanWitness) {
		verdict = Conflict
	} else if verdict == Allow && (!p.fresh || !p.verified || !p.contained) {
		verdict = Unknown
	}
	if verdict != c.Verdict {
		return errors.New("kernel certificate: verdict mismatch")
	}
	return nil
}
