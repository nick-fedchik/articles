package teoscompose

import (
	"crypto/sha256"
	"math"
	"time"

	evidence "github.com/nick-fedchik/articles/Expert_System_Circuits/src"
)

type Observation struct {
	Name       string `json:"name"`
	State      string `json:"state"`
	Permitted  bool   `json:"permitted"`
	References int    `json:"references"`
}

type Countermodel struct {
	Name               string `json:"name"`
	LedgerBalanced     bool   `json:"ledger_balanced"`
	SourceBytesIntact  bool   `json:"source_bytes_intact"`
	EntailedByQuote    bool   `json:"entailed_by_quote"`
	EmpiricallyCorrect bool   `json:"empirically_correct"`
	GatePermitted      bool   `json:"gate_permitted"`
}

type Report struct {
	Kind   string `json:"kind"`
	FanOut struct {
		Branches                int     `json:"branches"`
		PerBranchRating         float64 `json:"per_branch_rating"`
		RatingSum               float64 `json:"rating_sum"`
		MergedRating            float64 `json:"merged_rating"`
		MergedReferences        int     `json:"merged_references"`
		BalancedWithoutCreation bool    `json:"balanced_without_creation"`
		BalancedWithCreation    bool    `json:"balanced_with_creation"`
	} `json:"fan_out"`
	Gates         []Observation  `json:"gates"`
	Countermodels []Countermodel `json:"countermodels"`
	Example       struct {
		NominalK       float64 `json:"nominal_K"`
		StandardK      float64 `json:"standard_K"`
		ExpandedK      float64 `json:"expanded_K"`
		FlowScore      float64 `json:"flow_score_per_tick"`
		WeakLimitState string  `json:"limit_925_state"`
	} `json:"example"`
}

func studyOrigin(id, text string) (Origin, error) {
	document := []byte(text)
	source, err := evidence.NewSource(id, "1", document, 0.6)
	if err != nil {
		return Origin{}, err
	}
	return NewOrigin(source, document, 0, len(document), sha256.Sum256(document))
}

// Study executes synthetic, deterministic constructions; it reports neither
// hardware timing nor measured truth/calibration rates.
func Study() (Report, error) {
	report := Report{Kind: "synthetic deterministic operator constructions"}
	source, err := studyOrigin("measurement", "Load is 930 K, interval [920,940] K.")
	if err != nil {
		return report, err
	}
	normSource, err := studyOrigin("norm", "Synthetic N1 permits turbine operation for load in [0,950] K.")
	if err != nil {
		return report, err
	}
	signal, err := AssumeSignal("load observation", evidence.True, 0.6, source)
	if err != nil {
		return report, err
	}
	branches, err := FanOut(signal, 4)
	if err != nil {
		return report, err
	}
	merged := branches[0]
	for _, branch := range branches[1:] {
		merged, err = Merge(merged, branch)
		if err != nil {
			return report, err
		}
	}
	report.FanOut.Branches = len(branches)
	report.FanOut.PerBranchRating = signal.Rating()
	for _, branch := range branches {
		report.FanOut.RatingSum += branch.Rating()
	}
	report.FanOut.MergedRating = merged.Rating()
	report.FanOut.MergedReferences = len(merged.Provenance())
	report.FanOut.BalancedWithoutCreation, err = (Ledger{Input: 1, Output: uint64(len(branches))}).Balanced()
	if err != nil {
		return report, err
	}
	report.FanOut.BalancedWithCreation, err = (Ledger{Input: 1, Created: uint64(len(branches) - 1), Output: uint64(len(branches))}).Balanced()
	if err != nil {
		return report, err
	}
	roles := evidence.Roles{Actor: "controller", Action: "run", Target: "turbine"}
	rule := evidence.Rule{Pattern: roles, Sort: "K", Window: evidence.Interval{Low: 0, High: 950},
		MaxAge: time.Second, Modal: evidence.True}
	norm, err := NewNorm(rule, normSource)
	if err != nil {
		return report, err
	}
	for _, row := range []struct {
		name        string
		value       float64
		low, high   float64
		unit        string
		normPresent bool
	}{
		{"inside", 930, 920, 940, "K", true},
		{"boundary", 950, 940, 960, "K", true},
		{"outside", 970, 960, 980, "K", true},
		{"wrong_unit", 930, 920, 940, "Pa", true},
		{"fact_without_norm", 930, 920, 940, "K", false},
	} {
		m, err := NewMeasurement(row.value, evidence.Interval{Low: row.low, High: row.high}, row.unit)
		if err != nil {
			return report, err
		}
		fact, err := NewFact(merged, m, roles, 0)
		if err != nil {
			return report, err
		}
		n := &norm
		if !row.normPresent {
			n = nil
		}
		result, err := Gate(fact, n)
		if err != nil {
			return report, err
		}
		report.Gates = append(report.Gates, Observation{row.name, result.State.String(), result.Permitted, len(result.Provenance)})
	}
	balanced, err := (Ledger{Input: 1, Output: 1}).Balanced()
	if err != nil {
		return report, err
	}
	// Countermodel 1: the quote gives 930, yet the assumed semantic input says 900.
	// Countermodel 2: the quote is correctly read, but the synthetic world is 980.
	for _, row := range []struct {
		name    string
		nominal float64
	}{
		{"wrong_quote_interpretation_world_980_K", 900},
		{"faithful_quote_but_false_source_world_980_K", 930},
	} {
		m, err := NewMeasurement(row.nominal, evidence.Interval{Low: row.nominal - 10, High: row.nominal + 10}, "K")
		if err != nil {
			return report, err
		}
		assumed, err := AssumeSignal(row.name, evidence.True, 0.6, source)
		if err != nil {
			return report, err
		}
		fact, err := NewFact(assumed, m, roles, 0)
		if err != nil {
			return report, err
		}
		result, err := Gate(fact, &norm)
		if err != nil {
			return report, err
		}
		report.Countermodels = append(report.Countermodels, Countermodel{
			row.name, balanced, true, row.nominal == 930,
			m.Range().Low <= 980 && 980 <= m.Range().High, result.Permitted,
		})
	}
	report.Example.NominalK = 850 + 2*40
	report.Example.StandardK = math.Hypot(3, 2*2)
	report.Example.ExpandedK = 2 * report.Example.StandardK
	report.Example.FlowScore, err = FlowScore(2, 0.85, 0.95, 0.45, report.Example.ExpandedK, 20)
	if err != nil {
		return report, err
	}
	m, err := NewMeasurement(930, evidence.Interval{Low: 920, High: 940}, "K")
	if err != nil {
		return report, err
	}
	state, err := IntervalGuard(m, evidence.Interval{Low: 0, High: 925}, "K")
	report.Example.WeakLimitState = state.String()
	return report, err
}
