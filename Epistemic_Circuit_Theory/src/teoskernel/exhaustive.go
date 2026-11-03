package teoskernel

import (
	"errors"
	"fmt"
	"runtime"
	"sync"
	"time"

	evidence "github.com/nick-fedchik/articles/Epistemic_Circuit_Theory/src"
)

// Bounds selects the finite domain enumerated by Exhaustive. The atoms are
// fixed: 0 is the guard atom of quantity "pressure" [0,10] bar, 1 and 2 are
// evidence atoms; atom 2 may be derived by one signed evidence rule.
type Bounds struct {
	// MaxNorms is 1 or 2; every pair of norms is combined with every strict
	// priority (none, first over second, second over first).
	MaxNorms int
	// MaxConditions is 0, 1 or 2 condition literals per norm.
	MaxConditions int
	// Exceptions enables an optional exception literal per norm.
	Exceptions bool
}

// Report summarizes an exhaustive run. Every counter has an explicit
// denominator: Cases.
type Report struct {
	Bounds             Bounds
	Theories           int
	Cases              int
	Verdicts           map[evidence.NormDecision]int
	AllowWithSafetyOK  int
	Counterexamples    []string
	ForgeriesAttempted int
	ForgeriesRejected  int
	PFalseCases        int
	ElapsedSeconds     float64 `json:"-"`
}

var literalPool = func() []evidence.SignedLiteral {
	var out []evidence.SignedLiteral
	for atom := 0; atom < 3; atom++ {
		out = append(out, evidence.SignedLiteral{Atom: atom}, evidence.SignedLiteral{Atom: atom, Negative: true})
	}
	return out
}()

func conditionSets(max int) [][]evidence.SignedLiteral {
	out := [][]evidence.SignedLiteral{nil}
	for i := range literalPool {
		if max >= 1 {
			out = append(out, []evidence.SignedLiteral{literalPool[i]})
		}
		for j := i + 1; j < len(literalPool) && max >= 2; j++ {
			out = append(out, []evidence.SignedLiteral{literalPool[i], literalPool[j]})
		}
	}
	return out
}

func normPool(b Bounds) []evidence.NormRule {
	exceptions := []*evidence.SignedLiteral{nil}
	if b.Exceptions {
		for i := range literalPool {
			exceptions = append(exceptions, &literalPool[i])
		}
	}
	var out []evidence.NormRule
	for _, effect := range []evidence.NormEffect{evidence.Permit, evidence.Prohibit} {
		for _, conditions := range conditionSets(b.MaxConditions) {
			for _, exception := range exceptions {
				out = append(out, evidence.NormRule{Version: "1", Action: "open", Effect: effect, Conditions: conditions, Exception: exception})
			}
		}
	}
	return out
}

var evidencePool = [][]evidence.SignedRule{
	nil,
	{{ID: "e", Version: "1", Premises: []evidence.SignedLiteral{{Atom: 1}}, Conclusion: evidence.SignedLiteral{Atom: 2}}},
	{{ID: "e", Version: "1", Premises: []evidence.SignedLiteral{{Atom: 1, Negative: true}}, Conclusion: evidence.SignedLiteral{Atom: 2}}},
	{{ID: "e", Version: "1", Premises: []evidence.SignedLiteral{{Atom: 1}}, Conclusion: evidence.SignedLiteral{Atom: 2, Negative: true}}},
	{{ID: "e", Version: "1", Premises: []evidence.SignedLiteral{{Atom: 0}}, Conclusion: evidence.SignedLiteral{Atom: 2}}},
}

// Theories enumerates the specifications of the bounded domain.
func Theories(b Bounds) ([]Spec, error) {
	if b.MaxNorms < 1 || b.MaxNorms > 2 || b.MaxConditions < 0 || b.MaxConditions > 2 {
		return nil, errors.New("teoskernel: unsupported exhaustive bounds")
	}
	pool := normPool(b)
	var theories [][]evidence.NormRule
	var priorities [][]evidence.NormPriority
	for _, a := range pool {
		a.ID = "n0"
		theories, priorities = append(theories, []evidence.NormRule{a}), append(priorities, nil)
		if b.MaxNorms < 2 {
			continue
		}
		for _, c := range pool {
			c.ID = "n1"
			for _, p := range [][]evidence.NormPriority{nil, {{Higher: "n0", Lower: "n1"}}, {{Higher: "n1", Lower: "n0"}}} {
				theories, priorities = append(theories, []evidence.NormRule{a, c}), append(priorities, p)
			}
		}
	}
	var specs []Spec
	for _, rules := range evidencePool {
		for i := range theories {
			specs = append(specs, Spec{
				Action:         "open",
				Theory:         evidence.NormTheory{Atoms: 3, Evidence: rules, Norms: theories[i], Priorities: priorities[i]},
				Guards:         []Guard{{Atom: 0, Quantity: "pressure", Unit: "bar", Low: 0, High: 10}},
				TrustedSources: []string{"s1", "s2"},
				MaxAge:         5 * time.Second,
			})
		}
	}
	return specs, nil
}

const exhaustiveNow = int64(1_000_000_000_000)

// Inputs enumerates the input snapshots of the bounded domain: four-valued
// facts on atoms 1 and 2, seven measurement patterns and seven freshness or
// provenance perturbations of the first input item.
func Inputs() []Input {
	factVariants := func(atom int) [][]Fact {
		pos := Fact{Literal: evidence.SignedLiteral{Atom: atom}, Sources: []string{"s1"}, Time: exhaustiveNow}
		neg := Fact{Literal: evidence.SignedLiteral{Atom: atom, Negative: true}, Sources: []string{"s2"}, Time: exhaustiveNow}
		return [][]Fact{nil, {pos}, {neg}, {pos, neg}}
	}
	m := func(unit string, low, high float64) Measurement {
		return Measurement{Quantity: "pressure", Unit: unit, Low: low, High: high, Sources: []string{"s1"}, Time: exhaustiveNow}
	}
	measurements := [][]Measurement{nil, {m("bar", 2, 3)}, {m("bar", 9, 11)}, {m("bar", 11, 12)},
		{m("kPa", 2, 3)}, {m("bar", 2, 3), m("bar", 11, 12)}, {m("bar", 0, 10)}}
	var out []Input
	for _, a := range factVariants(1) {
		for _, b := range factVariants(2) {
			for _, ms := range measurements {
				for meta := 0; meta < 7; meta++ {
					in := Input{Now: exhaustiveNow, Facts: append(append([]Fact{}, a...), b...), Measurements: append([]Measurement{}, ms...)}
					var sources *[]string
					var t *int64
					switch {
					case len(in.Facts) > 0:
						f := &in.Facts[0]
						f.Sources = append([]string{}, f.Sources...)
						sources, t = &f.Sources, &f.Time
					case len(in.Measurements) > 0:
						mm := &in.Measurements[0]
						mm.Sources = append([]string{}, mm.Sources...)
						sources, t = &mm.Sources, &mm.Time
					default:
						if meta > 0 {
							continue
						}
					}
					switch meta {
					case 1:
						*t = exhaustiveNow - int64(5*time.Second) - 1
					case 2:
						*t = exhaustiveNow - int64(5*time.Second)
					case 3:
						*t = exhaustiveNow + 1
					case 4:
						*sources = []string{"rogue"}
					case 5:
						*sources = nil
					case 6:
						*sources = []string{"s1", "rogue"}
					}
					out = append(out, in)
				}
			}
		}
	}
	return out
}

// Exhaustive checks, for every (spec, input) of the bounded domain, that
// Evaluate succeeds, Verify accepts the honest certificate, an ALLOW verdict
// satisfies CheckSafety, and two forged ALLOW certificates are rejected
// whenever the honest verdict is not ALLOW.
func Exhaustive(b Bounds) (Report, error) {
	start := time.Now()
	specs, err := Theories(b)
	if err != nil {
		return Report{}, err
	}
	inputs := Inputs()
	report := Report{Bounds: b, Theories: len(specs), Verdicts: map[evidence.NormDecision]int{}}
	var mu sync.Mutex
	var firstErr error
	jobs := make(chan int)
	var wg sync.WaitGroup
	for w := 0; w < runtime.GOMAXPROCS(0); w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			local := Report{Verdicts: map[evidence.NormDecision]int{}}
			var localErr error
			for i := range jobs {
				for _, in := range inputs {
					if err := checkCase(specs[i], in, &local); err != nil && localErr == nil {
						localErr = fmt.Errorf("spec %d: %w", i, err)
					}
				}
			}
			mu.Lock()
			defer mu.Unlock()
			if localErr != nil && firstErr == nil {
				firstErr = localErr
			}
			report.Cases += local.Cases
			report.AllowWithSafetyOK += local.AllowWithSafetyOK
			report.ForgeriesAttempted += local.ForgeriesAttempted
			report.ForgeriesRejected += local.ForgeriesRejected
			report.PFalseCases += local.PFalseCases
			for k, v := range local.Verdicts {
				report.Verdicts[k] += v
			}
			if len(report.Counterexamples) < 10 {
				report.Counterexamples = append(report.Counterexamples, local.Counterexamples...)
			}
		}()
	}
	for i := range specs {
		jobs <- i
	}
	close(jobs)
	wg.Wait()
	report.ElapsedSeconds = time.Since(start).Seconds()
	return report, firstErr
}

func checkCase(s Spec, in Input, r *Report) error {
	c, err := Evaluate(s, in)
	if err != nil {
		return err
	}
	r.Cases++
	r.Verdicts[c.Verdict]++
	if err := Verify(s, in, c); err != nil {
		return fmt.Errorf("honest certificate rejected: %w", err)
	}
	violated := CheckSafety(s, in)
	if len(violated) > 0 {
		r.PFalseCases++
	}
	if c.Verdict == Allow {
		if len(violated) > 0 {
			if len(r.Counterexamples) < 10 {
				r.Counterexamples = append(r.Counterexamples, fmt.Sprintf("%+v %+v -> %v", s, in, violated))
			}
		} else {
			r.AllowWithSafetyOK++
		}
		return nil
	}
	forged := c
	forged.Verdict = Allow
	r.ForgeriesAttempted++
	if Verify(s, in, forged) != nil {
		r.ForgeriesRejected++
	}
	forged.Norm.Result.Decision = Allow
	forged.Fresh, forged.Verified, forged.Contained, forged.Violations = true, true, true, nil
	r.ForgeriesAttempted++
	if Verify(s, in, forged) != nil {
		r.ForgeriesRejected++
	}
	return nil
}
