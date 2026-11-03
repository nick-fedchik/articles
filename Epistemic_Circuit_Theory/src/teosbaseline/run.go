package teosbaseline

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/nick-fedchik/articles/Epistemic_Circuit_Theory/src/teoskernel"
)

// Arms in reporting order. B1S is an exploratory sensitivity variant that was
// not preregistered.
var Arms = []string{"A", "B0", "B1", "C", "B1S"}

// Record is the per-scenario outcome of every arm.
type Record struct {
	Class      string
	Index      int
	Control    bool
	Variant    string
	PhysSafe   bool
	PHolds     bool
	Violations string
	Decision   []Decision
	Raw        []string
	Audit      []int
}

// Evaluator runs every arm on a scenario with per-arm audit logs.
type Evaluator struct {
	spec teoskernel.Spec
	logs map[string]*AuditLog
}

// NewEvaluator returns an evaluator with fresh audit logs.
func NewEvaluator() *Evaluator {
	return &Evaluator{spec: CaseSpec(), logs: map[string]*AuditLog{"B0": {}, "B1": {}, "B1S": {}}}
}

// Run evaluates one arm on one input.
func (e *Evaluator) Run(arm string, in teoskernel.Input) (Outcome, error) {
	switch arm {
	case "A":
		return ArmA(in), nil
	case "B0":
		return ArmB0(in, e.logs[arm]), nil
	case "B1":
		return ArmB1(in, e.logs[arm]), nil
	case "B1S":
		return ArmB1Strict(in, e.logs[arm]), nil
	case "C":
		out, _, err := ArmC(e.spec, in)
		return out, err
	}
	return Outcome{}, fmt.Errorf("teosbaseline: unknown arm %q", arm)
}

// Evaluate runs all arms on all scenarios and checks the kernel invariant
// ALLOW <=> P on every scenario.
func Evaluate(scenarios []Scenario) ([]Record, error) {
	e := NewEvaluator()
	out := make([]Record, 0, len(scenarios))
	for _, s := range scenarios {
		violated := teoskernel.CheckSafety(e.spec, s.Input)
		r := Record{Class: s.Class, Index: s.Index, Control: s.Control, Variant: s.Variant,
			PhysSafe: s.World.Safe(), PHolds: len(violated) == 0, Violations: strings.Join(violated, ";")}
		for _, arm := range Arms {
			o, err := e.Run(arm, s.Input)
			if err != nil {
				return nil, fmt.Errorf("%s/%d arm %s: %w", s.Class, s.Index, arm, err)
			}
			r.Decision = append(r.Decision, o.Decision)
			r.Raw = append(r.Raw, o.Raw)
			r.Audit = append(r.Audit, o.Audit)
		}
		if c := r.Decision[slices.Index(Arms, "C")]; (c == Allow) != r.PHolds {
			return nil, fmt.Errorf("%s/%d: kernel ALLOW=%v but P=%v", s.Class, s.Index, c == Allow, r.PHolds)
		}
		out = append(out, r)
	}
	return out, nil
}

// Group labels the fault, control or clean subset of a class.
func (r Record) Group() string {
	switch {
	case r.Class == "CLEAN":
		return "clean"
	case r.Control:
		return "control"
	}
	return "fault"
}

// Bound is a zero-or-more-event upper bound with its explicit denominator.
type Bound struct {
	Events, Trials         int
	RuleOfThree            float64 `json:",omitempty"`
	CPUpper95TwoSided      float64
	CPUpper95OneSided      float64
	Defined                bool
	RuleOfThreeApplicable  bool
	DenominatorDescription string
}

// NewBound computes exact Clopper-Pearson bounds (and 3/n for zero events).
func NewBound(events, trials int, description string) Bound {
	b := Bound{Events: events, Trials: trials, DenominatorDescription: description}
	if trials <= 0 {
		return b
	}
	b.Defined = true
	b.CPUpper95TwoSided, _ = ClopperPearsonUpper(events, trials, 0.025)
	b.CPUpper95OneSided, _ = ClopperPearsonUpper(events, trials, 0.05)
	if events == 0 {
		b.RuleOfThree, b.RuleOfThreeApplicable = RuleOfThree(trials), true
	}
	return b
}

// Row aggregates one arm on one class group. Unsafe permits are ALLOW with the
// ground truth unsafe; false refusals are DENY or DEFER with it safe.
type Row struct {
	Class, Group, Arm               string
	N, Allow, Deny, Defer           int
	PSafe, PhysSafe                 int
	UnsafeP, UnsafePhys             int
	FalseRefusalP, FalseRefusalPhys int
	Automatic                       int
	AuditMedianBytes, AuditMaxBytes int
	UnsafePhysPerAllow              Bound
	UnsafePhysPerScenario           Bound
}

func median(values []int) int {
	if len(values) == 0 {
		return 0
	}
	s := slices.Clone(values)
	slices.Sort(s)
	return s[len(s)/2]
}

// Aggregate builds rows in class order, groups fault/control/clean, arm order.
func Aggregate(records []Record) []Row {
	type key struct{ class, group string }
	var order []key
	buckets := map[key][]Record{}
	for _, r := range records {
		k := key{r.Class, r.Group()}
		if _, ok := buckets[k]; !ok {
			order = append(order, k)
		}
		buckets[k] = append(buckets[k], r)
	}
	var rows []Row
	for _, k := range order {
		for ai, arm := range Arms {
			row := Row{Class: k.class, Group: k.group, Arm: arm}
			var audits []int
			for _, r := range buckets[k] {
				row.N++
				d := r.Decision[ai]
				audits = append(audits, r.Audit[ai])
				row.AuditMaxBytes = max(row.AuditMaxBytes, r.Audit[ai])
				if r.PHolds {
					row.PSafe++
				}
				if r.PhysSafe {
					row.PhysSafe++
				}
				switch d {
				case Allow:
					row.Allow++
					if !r.PHolds {
						row.UnsafeP++
					}
					if !r.PhysSafe {
						row.UnsafePhys++
					}
				case Deny:
					row.Deny++
				default:
					row.Defer++
				}
				if d != Allow && r.PHolds {
					row.FalseRefusalP++
				}
				if d != Allow && r.PhysSafe {
					row.FalseRefusalPhys++
				}
			}
			row.Automatic = row.Allow + row.Deny
			row.AuditMedianBytes = median(audits)
			row.UnsafePhysPerAllow = NewBound(row.UnsafePhys, row.Allow, "physically unsafe among ALLOW decisions")
			row.UnsafePhysPerScenario = NewBound(row.UnsafePhys, row.N, "unsafe ALLOW among all scenarios of the group")
			rows = append(rows, row)
		}
	}
	return rows
}

// TimingRow is a host microbenchmark: median per-scenario evaluation time.
type TimingRow struct {
	Class, Arm      string
	Scenarios, Reps int
	MedianNs        int64
}

// Timing measures median wall-clock evaluation time per scenario (host
// microbenchmark, non-deterministic, includes audit serialization).
func Timing(scenarios []Scenario, perClass, reps int) ([]TimingRow, error) {
	byClass := map[string][]Scenario{}
	var classes []string
	for _, s := range scenarios {
		if _, ok := byClass[s.Class]; !ok {
			classes = append(classes, s.Class)
		}
		if len(byClass[s.Class]) < perClass {
			byClass[s.Class] = append(byClass[s.Class], s)
		}
	}
	var rows []TimingRow
	for _, class := range classes {
		for _, arm := range Arms {
			e := NewEvaluator()
			var ns []int
			for _, s := range byClass[class] {
				start := time.Now()
				for range reps {
					if _, err := e.Run(arm, s.Input); err != nil {
						return nil, err
					}
				}
				ns = append(ns, int(time.Since(start).Nanoseconds()/int64(reps)))
			}
			rows = append(rows, TimingRow{Class: class, Arm: arm, Scenarios: len(byClass[class]), Reps: reps,
				MedianNs: int64(median(ns))})
		}
	}
	return rows, nil
}
