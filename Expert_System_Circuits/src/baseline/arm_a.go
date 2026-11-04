package baseline

import (
	"slices"

	"github.com/nick-fedchik/articles/Expert_System_Circuits/src/kernel"
)

// Decision is the normalized outcome of an arm: ALLOW, DENY or DEFER.
type Decision string

const (
	Allow Decision = "ALLOW"
	Deny  Decision = "DENY"
	Defer Decision = "DEFER"
)

// Outcome is an arm decision with its raw verdict and the bytes the arm
// appends to its certificate or audit log for this decision.
type Outcome struct {
	Decision Decision
	Raw      string
	Audit    int
}

// Reliability weights of arm A: a typical hand-tuned prior per source class.
const (
	naiveTrusted   = 0.95
	naiveUntrusted = 0.6
	naiveRequired  = 0.8
	naiveActive    = 0.5
)

// NaiveView is the information arm A receives: one scalar confidence per
// fact name and one point pressure reading with an untyped number.
type NaiveView struct {
	Pressure   float64
	HasReading bool
	Confidence [atomCount]float64
}

// ProjectNaive builds the arm A view. Each report contributes the reliability
// of its most reliable source; positive and negative support are combined by
// noisy-OR (independent reports, duplicates counted) into
// c = 0.5 + 0.5*(NOR(pro) - NOR(contra)); no report gives 0.5. The point
// reading is the midpoint of the most recently stamped measurement, read as a
// plain number in the configured unit (bar).
func ProjectNaive(in kernel.Input) NaiveView {
	var v NaiveView
	var missPro, missCon [atomCount]float64
	for a := range missPro {
		missPro[a], missCon[a] = 1, 1
	}
	for _, f := range in.Facts {
		w := naiveUntrusted
		if slices.ContainsFunc(f.Sources, isTrusted) {
			w = naiveTrusted
		}
		if f.Literal.Negative {
			missCon[f.Literal.Atom] *= 1 - w
		} else {
			missPro[f.Literal.Atom] *= 1 - w
		}
	}
	for a := range v.Confidence {
		v.Confidence[a] = 0.5 + 0.5*((1-missPro[a])-(1-missCon[a]))
	}
	latest := int64(0)
	for _, m := range in.Measurements {
		if !v.HasReading || m.Time >= latest {
			v.Pressure, latest, v.HasReading = (m.Low+m.High)/2, m.Time, true
		}
	}
	return v
}

// ArmA is the naive gate: a flat deny-first chain of point comparisons.
func ArmA(in kernel.Input) Outcome {
	v := ProjectNaive(in)
	c := v.Confidence
	// spec:A:begin
	decide := func() Decision {
		if c[AtomLock] > naiveActive || c[AtomInterlock] > naiveActive {
			return Deny
		}
		if c[AtomHotWork] > naiveActive && c[AtomGasTest] < naiveRequired {
			return Deny
		}
		if v.HasReading && v.Pressure >= LimitLowBar && v.Pressure <= LimitHighBar &&
			c[AtomAuthorized] >= naiveRequired && c[AtomDemand] >= naiveRequired {
			return Allow
		}
		return Deny
	}
	// spec:A:end
	d := decide()
	return Outcome{Decision: d, Raw: string(d)}
}
