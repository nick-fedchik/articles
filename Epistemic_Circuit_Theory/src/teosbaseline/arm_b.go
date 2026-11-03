package teosbaseline

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sync"

	evidence "github.com/nick-fedchik/articles/Epistemic_Circuit_Theory/src"
	"github.com/nick-fedchik/articles/Epistemic_Circuit_Theory/src/teoskernel"
)

// InstrumentMarginBar is the fixed datasheet accuracy B0 applies around a
// point reading. It equals the largest half-width of a healthy sensor in the
// generator, so B0 is not handicapped on healthy readings.
const InstrumentMarginBar = 0.3

// AuditLog is an append-only, hash-chained JSON-lines log. Arms B0 and B1
// write one record per decision; B0 never reads it back for decisions.
type AuditLog struct {
	mu    sync.Mutex
	lines [][]byte
	prev  string
}

// Append adds a record and returns the bytes written.
func (l *AuditLog) Append(record map[string]any) int {
	l.mu.Lock()
	defer l.mu.Unlock()
	record["seq"], record["prev"] = len(l.lines), l.prev
	line, err := json.Marshal(record)
	if err != nil {
		panic(err)
	}
	sum := sha256.Sum256(line)
	l.prev = hex.EncodeToString(sum[:])
	l.lines = append(l.lines, line)
	return len(line) + 1
}

// Len returns the number of records.
func (l *AuditLog) Len() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return len(l.lines)
}

// conventional Horn program over named boolean atoms; for B1 the atom space
// is doubled into pro (2a) and contra (2a+1) codes.
func hornProgram(atoms int, rules []evidence.HornRule) *evidence.HornProgram {
	p, err := evidence.CompileHorn(atoms, rules)
	if err != nil {
		panic(err)
	}
	return p
}

var (
	b0Horn = hornProgram(atomCount, []evidence.HornRule{
		{ID: "lock_engages_interlock", Version: "1", Premises: []int{AtomLock}, Conclusion: AtomInterlock}})
	b1Horn = hornProgram(2*atomCount, []evidence.HornRule{
		{ID: "lock_engages_interlock", Version: "1", Premises: []int{2 * AtomLock}, Conclusion: 2 * AtomInterlock}})
)

func closure(p *evidence.HornProgram, facts []bool) []bool {
	known, err := p.FullClosure(evidence.HornSnapshot{Revision: 1, Facts: facts, Enabled: []bool{true}})
	if err != nil {
		panic(err)
	}
	return known
}

func freshTime(now, t int64) bool { return t <= now && now-t <= int64(MaxAge) }

// Reading is a point value with its declared unit and timestamp.
type Reading struct {
	Value float64
	Unit  string
	Time  int64
}

// Flag is one boolean fact with the timestamp of the report it came from.
type Flag struct {
	Value bool
	Time  int64
}

// ConventionalView is the information B0 receives: a point reading with a
// unit, one boolean per fact (latest report wins; ties: later report), and
// the provenance of every report, which goes only to the audit log.
type ConventionalView struct {
	Reading    Reading
	HasReading bool
	Flags      map[int]Flag
	Provenance []string
}

// ProjectConventional builds the B0 view from the full observation.
func ProjectConventional(in teoskernel.Input) ConventionalView {
	v := ConventionalView{Flags: map[int]Flag{}}
	for _, m := range in.Measurements {
		if !v.HasReading || m.Time >= v.Reading.Time {
			v.Reading, v.HasReading = Reading{Value: (m.Low + m.High) / 2, Unit: m.Unit, Time: m.Time}, true
		}
		v.Provenance = append(v.Provenance, fmt.Sprint("measurement:", m.Sources))
	}
	for _, f := range in.Facts {
		if old, ok := v.Flags[f.Literal.Atom]; !ok || f.Time >= old.Time {
			v.Flags[f.Literal.Atom] = Flag{Value: !f.Literal.Negative, Time: f.Time}
		}
		v.Provenance = append(v.Provenance, fmt.Sprint(AtomNames[f.Literal.Atom], ":", f.Sources))
	}
	return v
}

// ArmB0 is the conventional engineered stack on conventional information.
func ArmB0(in teoskernel.Input, log *AuditLog) Outcome {
	v := ProjectConventional(in)
	// spec:B0:begin
	decide := func() (Decision, string) {
		if !v.HasReading || !freshTime(in.Now, v.Reading.Time) {
			return Deny, "pressure missing or stale"
		}
		for atom := range atomCount {
			if f, ok := v.Flags[atom]; ok && !freshTime(in.Now, f.Time) {
				return Deny, "stale fact " + AtomNames[atom]
			}
		}
		p, err := ParsePressure(v.Reading.Value, v.Reading.Unit)
		if err != nil {
			return Deny, err.Error()
		}
		if bar := p.Bar(); bar-InstrumentMarginBar < LimitLowBar || bar+InstrumentMarginBar > LimitHighBar {
			return Deny, "pressure outside window with instrument margin"
		}
		facts := make([]bool, atomCount)
		for atom, f := range v.Flags {
			facts[atom] = f.Value
		}
		k := closure(b0Horn, facts)
		rules := []struct {
			name   string
			when   bool
			effect Decision
		}{
			{"ban_under_lock", k[AtomLock], Deny},
			{"ban_interlock", k[AtomInterlock], Deny},
			{"hot_work_cleared", k[AtomHotWork] && k[AtomGasTest] && k[AtomAuthorized] && k[AtomDemand], Allow},
			{"ban_hot_work", k[AtomHotWork], Deny},
			{"open_when_safe", k[AtomAuthorized] && k[AtomDemand], Allow},
		}
		for _, r := range rules {
			if r.when {
				return r.effect, r.name
			}
		}
		return Deny, "default_deny"
	}
	// spec:B0:end
	d, reason := decide()
	size := log.Append(map[string]any{"arm": "B0", "decision": d, "reason": reason,
		"reading": v.Reading, "flags": v.Flags, "provenance": v.Provenance})
	return Outcome{Decision: d, Raw: reason, Audit: size}
}

// value is the four-valued state of an atom from pro/contra closure bits.
func value(pro, contra bool) evidence.State {
	switch {
	case pro && contra:
		return evidence.Both
	case pro:
		return evidence.True
	case contra:
		return evidence.False
	}
	return evidence.Unknown
}

// ArmB1 is the same conventional stack given the kernel's metadata (interval,
// pro/contra, source sets, timestamps, exception and priority knowledge),
// written directly from the derived STPA safety constraints SC-1..SC-9. It
// keeps typed unit conversion, which the kernel deliberately does not do.
func ArmB1(in teoskernel.Input, log *AuditLog) Outcome { return armB1(in, log, true) }

// ArmB1Strict is an exploratory sensitivity variant (not preregistered): B1
// that rejects (DEFER) every unit other than the declared bar instead of
// converting it.
func ArmB1Strict(in teoskernel.Input, log *AuditLog) Outcome { return armB1(in, log, false) }

func armB1(in teoskernel.Input, log *AuditLog, convert bool) Outcome {
	// spec:B1:begin
	decide := func() (Decision, string) {
		for i, m := range in.Measurements {
			if !freshTime(in.Now, m.Time) {
				return Defer, fmt.Sprintf("SC-1 stale measurement %d", i)
			}
			if !allTrusted(m.Sources) {
				return Defer, fmt.Sprintf("SC-6 unverified measurement %d", i)
			}
		}
		for i, f := range in.Facts {
			if f.Literal.Atom >= AtomUnreferencedNoise {
				continue
			}
			if !freshTime(in.Now, f.Time) {
				return Defer, fmt.Sprintf("SC-1 stale fact %d", i)
			}
			if !allTrusted(f.Sources) {
				return Defer, fmt.Sprintf("SC-6 unverified fact %d", i)
			}
		}
		if len(in.Measurements) == 0 {
			return Defer, "SC-3 pressure missing"
		}
		straddles := false
		for _, m := range in.Measurements {
			low, errLow := ParsePressure(m.Low, m.Unit)
			high, errHigh := ParsePressure(m.High, m.Unit)
			if errLow != nil || errHigh != nil || !convert && m.Unit != "bar" {
				return Defer, "SC-2 unknown unit " + m.Unit
			}
			if high.Bar() < LimitLowBar || low.Bar() > LimitHighBar {
				return Deny, "SC-3 pressure interval outside window"
			}
			straddles = straddles || low.Bar() < LimitLowBar || high.Bar() > LimitHighBar
		}
		if straddles {
			return Defer, "SC-3 pressure interval straddles window"
		}
		facts := make([]bool, 2*atomCount)
		for _, f := range in.Facts {
			code := 2 * f.Literal.Atom
			if f.Literal.Negative {
				code++
			}
			facts[code] = true
		}
		k := closure(b1Horn, facts)
		state := func(atom int) evidence.State { return value(k[2*atom], k[2*atom+1]) }
		for _, atom := range []int{AtomLock, AtomInterlock} {
			switch state(atom) {
			case evidence.True:
				return Deny, "SC-5 " + AtomNames[atom]
			case evidence.False:
			default:
				return Defer, "SC-4/SC-5 undetermined " + AtomNames[atom]
			}
		}
		for _, atom := range []int{AtomAuthorized, AtomDemand} {
			switch state(atom) {
			case evidence.False:
				return Deny, "SC-7/SC-8 missing " + AtomNames[atom]
			case evidence.True:
			default:
				return Defer, "SC-4/SC-8 undetermined " + AtomNames[atom]
			}
		}
		switch state(AtomHotWork) {
		case evidence.False:
			return Allow, "open_when_safe"
		case evidence.True:
			switch state(AtomGasTest) {
			case evidence.True:
				return Allow, "hot_work_cleared over ban_hot_work"
			case evidence.False:
				return Deny, "SC-9 ban_hot_work"
			}
		}
		return Defer, "SC-4/SC-9 undetermined hot work or gas test"
	}
	// spec:B1:end
	d, reason := decide()
	size := log.Append(map[string]any{"arm": "B1", "convert": convert, "decision": d, "reason": reason, "input": in})
	return Outcome{Decision: d, Raw: reason, Audit: size}
}

func allTrusted(sources []string) bool {
	if len(sources) == 0 {
		return false
	}
	for _, id := range sources {
		if !isTrusted(id) {
			return false
		}
	}
	return true
}

// ArmC is the kernel K over CaseSpec; DEFER_* verdicts map to DEFER and the
// audit size is the JSON size of the certificate.
func ArmC(spec teoskernel.Spec, in teoskernel.Input) (Outcome, teoskernel.Certificate, error) {
	cert, err := teoskernel.Evaluate(spec, in)
	if err != nil {
		return Outcome{}, cert, err
	}
	data, err := json.Marshal(cert)
	if err != nil {
		return Outcome{}, cert, err
	}
	d := Defer
	switch cert.Verdict {
	case teoskernel.Allow:
		d = Allow
	case teoskernel.Deny:
		d = Deny
	}
	return Outcome{Decision: d, Raw: string(cert.Verdict), Audit: len(data)}, cert, nil
}
