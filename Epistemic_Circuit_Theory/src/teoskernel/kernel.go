// Package teoskernel implements K, a minimal unified semantics of the TEOS
// admission gate for ONE proposed action over a finite, acyclic, already
// extracted evidence snapshot. It composes existing repository semantics and
// does not redefine them:
//
//   - four-valued evidence (evidence.State, Belnap bits True=2, False=1);
//   - the signed normative policy of evidence.CompileNormative,
//     NormProgram.Evaluate and NormProgram.Verify (explicit exceptions, strict
//     acyclic priorities, verdicts ALLOW, DENY, DEFER_UNKNOWN, DEFER_CONFLICT,
//     DEFER_NO_BASIS);
//   - the interval-guard convention of teoscompose.IntervalGuard (contained
//     TRUE, disjoint FALSE, partial overlap BOTH, unit mismatch UNKNOWN);
//   - provenance as a set of verified source identities and freshness as a
//     nonnegative age bound.
//
// Stable API:
//
//	cert, err := teoskernel.Evaluate(spec, input)    // certificate, self-checked
//	err = teoskernel.Verify(spec, input, cert)       // independent recomputation
//	violated := teoskernel.CheckSafety(spec, input)  // property P; empty when P holds
//	report, err := teoskernel.Exhaustive(bounds)     // bounded model check
//
// Kernel K. Facts are joined into a four-valued pre-closure assignment; each
// guard atom receives the join of IntervalGuard states of all measurements of
// its quantity (UNKNOWN when none). The existing normative layer computes the
// signed evidence closure and the policy decision D. Admission postconditions
// are: every input item is fresh (Now-MaxAge <= t <= Now) and carries a
// nonempty set of sources included in TrustedSources; every guard has at least
// one measurement and every measurement of a guarded quantity has the guard
// unit and lies inside the closed window. The verdict equals D unless D is
// ALLOW and a postcondition fails: then it is DENY (some interval disjoint),
// else DEFER_CONFLICT (some partial overlap), else DEFER_UNKNOWN.
//
// Soundness theorem (paper proof, Lean 4 proof of an abstract model in
// src/formal/kernel, bounded exhaustive check in Exhaustive): for every
// well-formed Spec S, Input I and Certificate c, Verify(S, I, c) == nil and
// c.Verdict == ALLOW imply len(CheckSafety(S, I)) == 0. Conversely, because
// validated measurements satisfy Low <= High, P implies ALLOW (Lean theorem
// allow_iff_P), so every DENY or DEFER marks a violated clause of P.
//
// Scope limits: finite and acyclic theories, one action per Spec, no natural
// language extraction, no signature cryptography (TrustedSources is the set of
// identities verified upstream), no timing or WCET claims.
package teoskernel

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	evidence "github.com/nick-fedchik/articles/Epistemic_Circuit_Theory/src"
	"github.com/nick-fedchik/articles/Epistemic_Circuit_Theory/src/teoscompose"
)

// Verdicts are the existing normative decisions.
const (
	Allow    = evidence.NormAllow
	Deny     = evidence.NormDeny
	Unknown  = evidence.NormUnknown
	Conflict = evidence.NormConflict
	NoBasis  = evidence.NormNoBasis
)

// Guard binds an atom to a hard closed interval on one measured quantity.
// Every declared guard is a mandatory admission postcondition.
type Guard struct {
	Atom      int
	Quantity  string
	Unit      string
	Low, High float64
}

// Spec is the static part of K: one action, its signed normative theory, the
// interval guards, the verified source identities and the freshness bound.
type Spec struct {
	Action         string
	Theory         evidence.NormTheory
	Guards         []Guard
	TrustedSources []string
	MaxAge         time.Duration
}

// Fact is one signed evidence literal with provenance and a Unix-ns timestamp.
type Fact struct {
	Literal evidence.SignedLiteral
	Sources []string
	Time    int64
}

// Measurement is an interval-valued reading of a declared quantity.
type Measurement struct {
	Quantity  string
	Unit      string
	Low, High float64
	Sources   []string
	Time      int64
}

// Input is the dynamic snapshot evaluated at time Now (Unix nanoseconds).
type Input struct {
	Now          int64
	Facts        []Fact
	Measurements []Measurement
}

// Certificate records every intermediate result checked by Verify.
type Certificate struct {
	SpecDigest  [32]byte
	InputDigest [32]byte
	// Atoms is the four-valued assignment before the evidence closure.
	Atoms []evidence.State
	// Guards is the joined interval-guard state per declared guard.
	Guards     []evidence.State
	Norm       evidence.NormCertificate
	Fresh      bool
	Verified   bool
	Contained  bool
	Violations []string
	Verdict    evidence.NormDecision
}

const (
	maxItems   = 4096
	maxSources = 64
)

// Violation codes of the admission postconditions.
const (
	codeStale      = "STALE"
	codeUnverified = "UNVERIFIED_SOURCE"
	codeUnit       = "UNIT_MISMATCH"
	codeOutOfBound = "GUARD_OUT_OF_BOUNDS"
	codeBoundary   = "GUARD_BOUNDARY"
	codeMissing    = "GUARD_MISSING"
)

func finite(v float64) bool { return !math.IsNaN(v) && !math.IsInf(v, 0) }

func validateSpec(s Spec) error {
	if strings.TrimSpace(s.Action) == "" || s.MaxAge < 0 || len(s.Guards) > s.Theory.Atoms ||
		len(s.TrustedSources) > maxItems {
		return errors.New("teoskernel: action, nonnegative max age and guard count required")
	}
	for _, rule := range s.Theory.Norms {
		if rule.Action != s.Action {
			return errors.New("teoskernel: every norm must govern the single spec action")
		}
	}
	for _, id := range s.TrustedSources {
		if strings.TrimSpace(id) == "" {
			return errors.New("teoskernel: empty trusted source identity")
		}
	}
	atoms, quantities := map[int]bool{}, map[string]bool{}
	for _, g := range s.Guards {
		if g.Atom < 0 || g.Atom >= s.Theory.Atoms || atoms[g.Atom] || quantities[g.Quantity] ||
			strings.TrimSpace(g.Quantity) == "" || strings.TrimSpace(g.Unit) == "" ||
			!finite(g.Low) || !finite(g.High) || g.Low > g.High {
			return errors.New("teoskernel: guard needs unique atom and quantity, unit and finite ordered bounds")
		}
		atoms[g.Atom], quantities[g.Quantity] = true, true
	}
	for _, rule := range s.Theory.Evidence {
		if atoms[rule.Conclusion.Atom] {
			return errors.New("teoskernel: evidence rules may not conclude a guard atom")
		}
	}
	return nil
}

func validateInput(s Spec, in Input) error {
	if len(in.Facts) > maxItems || len(in.Measurements) > maxItems {
		return errors.New("teoskernel: input size limit exceeded")
	}
	guarded, quantities := map[int]bool{}, map[string]bool{}
	for _, g := range s.Guards {
		guarded[g.Atom], quantities[g.Quantity] = true, true
	}
	for _, f := range in.Facts {
		if f.Literal.Atom < 0 || f.Literal.Atom >= s.Theory.Atoms || guarded[f.Literal.Atom] || len(f.Sources) > maxSources {
			return errors.New("teoskernel: fact literal outside non-guard atom range")
		}
	}
	for _, m := range in.Measurements {
		if !quantities[m.Quantity] || strings.TrimSpace(m.Unit) == "" || !finite(m.Low) || !finite(m.High) ||
			m.Low > m.High || len(m.Sources) > maxSources {
			return errors.New("teoskernel: measurement needs declared quantity, unit and finite ordered interval")
		}
	}
	return nil
}

func digest(v any) ([32]byte, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return [32]byte{}, err
	}
	return sha256.Sum256(data), nil
}

// fresh reports Now-MaxAge <= t <= Now without signed overflow.
func fresh(now, t int64, maxAge time.Duration) bool {
	return t <= now && uint64(now)-uint64(t) <= uint64(maxAge)
}

type item struct {
	kind    string
	index   int
	sources []string
	time    int64
}

func items(in Input) []item {
	out := []item{}
	for i, f := range in.Facts {
		out = append(out, item{"fact", i, f.Sources, f.Time})
	}
	for i, m := range in.Measurements {
		out = append(out, item{"measurement", i, m.Sources, m.Time})
	}
	return out
}

func compile(s Spec) (*evidence.NormProgram, error) {
	if err := validateSpec(s); err != nil {
		return nil, err
	}
	return evidence.CompileNormative(s.Theory)
}

// Evaluate computes K(S, I) and returns a certificate that Verify accepts.
func Evaluate(s Spec, in Input) (Certificate, error) {
	program, err := compile(s)
	if err != nil {
		return Certificate{}, err
	}
	if err := validateInput(s, in); err != nil {
		return Certificate{}, err
	}
	c := Certificate{Fresh: true, Verified: true, Contained: true}
	if c.SpecDigest, err = digest(s); err != nil {
		return Certificate{}, err
	}
	if c.InputDigest, err = digest(in); err != nil {
		return Certificate{}, err
	}
	c.Atoms = make([]evidence.State, s.Theory.Atoms)
	for _, f := range in.Facts {
		bit := evidence.True
		if f.Literal.Negative {
			bit = evidence.False
		}
		c.Atoms[f.Literal.Atom] |= bit
	}
	trusted := map[string]bool{}
	for _, id := range s.TrustedSources {
		trusted[id] = true
	}
	for _, it := range items(in) {
		if !fresh(in.Now, it.time, s.MaxAge) {
			c.Fresh = false
			c.Violations = append(c.Violations, fmt.Sprintf("%s:%s:%d", codeStale, it.kind, it.index))
		}
		ok := len(it.sources) > 0
		for _, id := range it.sources {
			ok = ok && trusted[id]
		}
		if !ok {
			c.Verified = false
			c.Violations = append(c.Violations, fmt.Sprintf("%s:%s:%d", codeUnverified, it.kind, it.index))
		}
	}
	outOfBounds, boundary := false, false
	c.Guards = make([]evidence.State, len(s.Guards))
	for gi, g := range s.Guards {
		seen := false
		for mi, m := range in.Measurements {
			if m.Quantity != g.Quantity {
				continue
			}
			seen = true
			measurement, err := teoscompose.NewMeasurement(m.Low, evidence.Interval{Low: m.Low, High: m.High}, m.Unit)
			if err != nil {
				return Certificate{}, err
			}
			state, err := teoscompose.IntervalGuard(measurement, evidence.Interval{Low: g.Low, High: g.High}, g.Unit)
			if err != nil {
				return Certificate{}, err
			}
			c.Guards[gi] |= state
			code := ""
			switch state {
			case evidence.Unknown:
				code = codeUnit
			case evidence.False:
				code, outOfBounds = codeOutOfBound, true
			case evidence.Both:
				code, boundary = codeBoundary, true
			}
			if code != "" {
				c.Contained = false
				c.Violations = append(c.Violations, fmt.Sprintf("%s:measurement:%d", code, mi))
			}
		}
		if !seen {
			c.Contained = false
			c.Violations = append(c.Violations, fmt.Sprintf("%s:guard:%d", codeMissing, gi))
		}
		c.Atoms[g.Atom] = c.Guards[gi]
	}
	world, err := program.NewWorld(c.Atoms)
	if err != nil {
		return Certificate{}, err
	}
	if c.Norm, err = program.Evaluate(world, s.Action); err != nil {
		return Certificate{}, err
	}
	c.Verdict = c.Norm.Result.Decision
	if c.Verdict == Allow {
		switch {
		case outOfBounds:
			c.Verdict = Deny
		case boundary:
			c.Verdict = Conflict
		case !c.Fresh || !c.Verified || !c.Contained:
			c.Verdict = Unknown
		}
	}
	if err := Verify(s, in, c); err != nil {
		return Certificate{}, fmt.Errorf("teoskernel: self-check failed: %w", err)
	}
	return c, nil
}
