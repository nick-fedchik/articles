// Package teoscompose implements a bounded TEOS operator fragment, not a complete
// TEOS/EARS/ASPIC+ semantics. Source integrity, assumed leaf meaning, evidence
// logic, interval applicability and normative authorization remain distinct.
package teoscompose

import (
	"crypto/sha256"
	"errors"
	"math"
	"sort"
	"strings"
	"time"

	evidence "github.com/nick-fedchik/articles/Epistemic_Circuit_Theory/src"
)

type Reference struct {
	ID, Version  string
	DocumentHash [sha256.Size]byte
	QuoteHash    [sha256.Size]byte
	Start, End   int
}

type Origin struct {
	ref      Reference
	verified bool
}

// NewOrigin checks bytes against supplied anchors; it does not authenticate the
// publisher, establish trust in an anchor, or interpret the quoted statement.
func NewOrigin(source evidence.Source, document []byte, start, end int, quoteHash [32]byte) (Origin, error) {
	checked, err := evidence.NewSource(source.ID, source.Version, document, source.Weight)
	if err != nil {
		return Origin{}, err
	}
	if checked.Digest != source.Digest {
		return Origin{}, errors.New("origin: document digest mismatch")
	}
	if err := evidence.VerifyQuote(document, start, end, quoteHash); err != nil {
		return Origin{}, err
	}
	return Origin{Reference{source.ID, source.Version, source.Digest, quoteHash, start, end}, true}, nil
}

func provenance(groups ...[]Origin) ([]Origin, error) {
	seen := map[Reference]Origin{}
	versions := map[[2]string][32]byte{}
	for _, group := range groups {
		for _, origin := range group {
			if !origin.verified {
				return nil, errors.New("provenance: unverified origin")
			}
			key := [2]string{origin.ref.ID, origin.ref.Version}
			if digest, ok := versions[key]; ok && digest != origin.ref.DocumentHash {
				return nil, errors.New("provenance: conflicting document for one source version")
			}
			versions[key] = origin.ref.DocumentHash
			seen[origin.ref] = origin
		}
	}
	out := make([]Origin, 0, len(seen))
	for _, origin := range seen {
		out = append(out, origin)
	}
	sort.Slice(out, func(i, j int) bool {
		a, b := out[i].ref, out[j].ref
		if a.ID != b.ID {
			return a.ID < b.ID
		}
		if a.Version != b.Version {
			return a.Version < b.Version
		}
		if a.Start != b.Start {
			return a.Start < b.Start
		}
		return a.End < b.End
	})
	return out, nil
}

type Signal struct {
	claim   string
	state   evidence.State
	rating  float64
	origins []Origin
}

// AssumeSignal makes the semantic leaf assumption explicit. Verified source
// bytes alone do not prove that claim follows from those bytes.
func AssumeSignal(claim string, state evidence.State, rating float64, origins ...Origin) (Signal, error) {
	if strings.TrimSpace(claim) == "" || !state.Valid() || !finite(rating) || rating < 0 || rating > 1 {
		return Signal{}, errors.New("signal: claim, four-valued state and rating in [0,1] required")
	}
	refs, err := provenance(origins)
	if err != nil {
		return Signal{}, err
	}
	return Signal{claim, state, rating, refs}, nil
}

func (s Signal) Claim() string         { return s.claim }
func (s Signal) State() evidence.State { return s.state }
func (s Signal) Rating() float64       { return s.rating }
func (s Signal) Provenance() []Reference {
	out := make([]Reference, len(s.origins))
	for i, origin := range s.origins {
		out[i] = origin.ref
	}
	return out
}

func valid(s Signal) bool { return s.claim != "" && s.state.Valid() }

func FanOut(s Signal, branches int) ([]Signal, error) {
	if !valid(s) || branches < 1 {
		return nil, errors.New("fan-out: valid signal and positive branch count required")
	}
	out := make([]Signal, branches)
	for i := range out {
		out[i] = s
		out[i].origins = append([]Origin(nil), s.origins...)
	}
	return out, nil
}

// Merge is the knowledge join for the SAME proposition, with exact-reference
// deduplication. Max rating is an explicit no-reinforcement envelope, not the
// independent-source S-conorm of chapter 8 or a calibrated truth probability.
func Merge(a, b Signal) (Signal, error) {
	if !valid(a) || !valid(b) || a.claim != b.claim {
		return Signal{}, errors.New("merge: same proposition required")
	}
	refs, err := provenance(a.origins, b.origins)
	if err != nil {
		return Signal{}, err
	}
	return AssumeSignal(a.claim, a.state|b.state, math.Max(a.rating, b.rating), refs...)
}

func And(a, b Signal) (Signal, error) {
	if !valid(a) || !valid(b) {
		return Signal{}, errors.New("and: valid inputs required")
	}
	state, err := evidence.TruthAnd(a.state, b.state)
	if err != nil {
		return Signal{}, err
	}
	refs, err := provenance(a.origins, b.origins)
	if err != nil {
		return Signal{}, err
	}
	return AssumeSignal("("+a.claim+" AND "+b.claim+")", state, math.Min(a.rating, b.rating), refs...)
}

func Not(s Signal) (Signal, error) {
	if !valid(s) {
		return Signal{}, errors.New("not: valid input required")
	}
	state := ((s.state & 1) << 1) | ((s.state & 2) >> 1)
	return AssumeSignal("NOT("+s.claim+")", state, s.rating, s.origins...)
}

type Measurement struct {
	value  float64
	bounds evidence.Interval
	sort   string
}

func NewMeasurement(value float64, bounds evidence.Interval, unit string) (Measurement, error) {
	if !finite(value) || !finite(bounds.Low) || !finite(bounds.High) ||
		bounds.Low > value || value > bounds.High || strings.TrimSpace(unit) == "" {
		return Measurement{}, errors.New("measurement: finite nominal inside ordered interval and unit required")
	}
	return Measurement{value, bounds, unit}, nil
}

func (m Measurement) Value() float64           { return m.value }
func (m Measurement) Range() evidence.Interval { return m.bounds }
func (m Measurement) Unit() string             { return m.sort }

type FactPort struct {
	signal      Signal
	measurement Measurement
	fact        evidence.Fact
}

func NewFact(signal Signal, m Measurement, roles evidence.Roles, age time.Duration) (FactPort, error) {
	if !valid(signal) || m.sort == "" {
		return FactPort{}, errors.New("fact: initialized signal and measurement required")
	}
	fact := evidence.Fact{Roles: roles, Sort: m.sort, Range: m.bounds, Age: age}
	identity := evidence.Rule{Pattern: roles, Sort: m.sort, Window: m.bounds, MaxAge: time.Duration(math.MaxInt64), Modal: evidence.True}
	if _, err := evidence.EvaluateRule(identity, fact); err != nil {
		return FactPort{}, err
	}
	return FactPort{signal, m, fact}, nil
}

type NormPort struct {
	rule    evidence.Rule
	origins []Origin
}

// NewNorm represents one supplied permission rule, not conflict resolution,
// legal validity, revocation checking or a new normative policy.
func NewNorm(rule evidence.Rule, origins ...Origin) (NormPort, error) {
	roles := rule.Pattern
	for _, slot := range []*string{&roles.Actor, &roles.Action, &roles.Target} {
		if *slot == "*" {
			*slot = "witness"
		}
	}
	fact := evidence.Fact{Roles: roles, Sort: rule.Sort, Range: rule.Window}
	if _, err := evidence.EvaluateRule(rule, fact); err != nil {
		return NormPort{}, err
	}
	refs, err := provenance(origins)
	if err != nil {
		return NormPort{}, err
	}
	return NormPort{rule, refs}, nil
}

// IntervalGuard reuses the repository's existing interval boundary convention:
// contained TRUE, disjoint FALSE, overlap BOTH, unit mismatch UNKNOWN.
func IntervalGuard(m Measurement, window evidence.Interval, unit string) (evidence.State, error) {
	if m.sort == "" {
		return evidence.Unknown, errors.New("guard: initialized measurement required")
	}
	roles := evidence.Roles{Actor: "sensor", Action: "measure", Target: "load"}
	decision, err := evidence.EvaluateRule(
		evidence.Rule{Pattern: roles, Sort: unit, Window: window, MaxAge: time.Second, Modal: evidence.True},
		evidence.Fact{Roles: roles, Sort: m.sort, Range: m.bounds})
	return decision.State, err
}

type GateResult struct {
	State      evidence.State
	Permitted  bool
	Reasons    []string
	Provenance []Reference
}

// Gate composes the assumed fact signal with EvaluateRule/Priority. A positive
// analog score or a source hash is never substituted for explicit permission.
func Gate(fact FactPort, norm *NormPort) (GateResult, error) {
	if !valid(fact.signal) || fact.measurement.sort == "" {
		return GateResult{}, errors.New("gate: initialized fact required")
	}
	result := GateResult{State: evidence.Unknown, Provenance: fact.signal.Provenance()}
	if norm == nil {
		result.State, _ = evidence.Priority(fact.signal.state, evidence.Unknown)
		result.Reasons = []string{"MISSING_NORM"}
		return result, nil
	}
	decision, err := evidence.EvaluateRule(norm.rule, fact.fact)
	if err != nil {
		return GateResult{}, err
	}
	refs, err := provenance(fact.signal.origins, norm.origins)
	if err != nil {
		return GateResult{}, err
	}
	result.State, _ = evidence.Priority(fact.signal.state, decision.State)
	result.Provenance = Signal{origins: refs}.Provenance()
	for _, reason := range decision.Reasons {
		result.Reasons = append(result.Reasons, string(reason))
	}
	if len(fact.signal.origins) == 0 || len(norm.origins) == 0 {
		result.State, _ = evidence.Priority(result.State, evidence.Unknown)
		result.Reasons = append(result.Reasons, "MISSING_PROVENANCE")
	}
	result.Permitted = result.State == evidence.True
	return result, nil
}

// Ledger counts message/resource tokens, not evidence truth or intensive ratings.
type Ledger struct {
	Input, Created, Before, Output, Rejected, After uint64
}

func (l Ledger) Balanced() (bool, error) {
	sum := func(a, b, c uint64) (uint64, error) {
		if math.MaxUint64-a < b || math.MaxUint64-(a+b) < c {
			return 0, errors.New("ledger: counter overflow")
		}
		return a + b + c, nil
	}
	left, err := sum(l.Input, l.Created, l.Before)
	if err != nil {
		return false, err
	}
	right, err := sum(l.Output, l.Rejected, l.After)
	return left == right && err == nil, err
}

// FlowScore is equation (3.5)'s phenomenological scheduling score, not a verdict.
func FlowScore(gain, factRating, normRating, threshold, expanded, tolerance float64) (float64, error) {
	for _, v := range []float64{gain, factRating, normRating, threshold, expanded, tolerance} {
		if !finite(v) {
			return 0, errors.New("flow score: finite parameters required")
		}
	}
	if gain <= 0 || factRating < 0 || factRating > 1 || normRating < 0 || normRating > 1 ||
		threshold <= 0 || threshold >= 1 || expanded < 0 || tolerance <= 0 {
		return 0, errors.New("flow score: parameter outside domain")
	}
	score := gain * math.Max(0, normRating-threshold) * factRating * math.Max(0, 1-expanded/tolerance)
	if !finite(score) {
		return 0, errors.New("flow score: overflow")
	}
	return score, nil
}

func finite(v float64) bool { return !math.IsNaN(v) && !math.IsInf(v, 0) }
