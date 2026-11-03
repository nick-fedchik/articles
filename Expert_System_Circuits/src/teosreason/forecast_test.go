package teosreason

import (
	"fmt"
	"reflect"
	"sort"
	"testing"

	evidence "github.com/nick-fedchik/articles/Expert_System_Circuits/src"
)

const (
	BD = "batch_delayed(p)"
	RD = "reserve_depleted(p)"
	EA = "expedite_available(p)"
	QS = "qa_staff_shortage(p)"
	BE = "buffer_exhausted(p)"
	RW = "rework_window_lost(p)"
	AH = "assembly_halted(p)"
	QH = "qa_hold(p)"
	SD = "schedule_delayed(p)"
)

// forecastExpected is written by hand from protocol 02 (§3.5), never copied from output.
type forecastExpected struct {
	pre         evidence.State
	preRoots    []string
	target      evidence.State
	targetRoots []string
	status      map[string]evidence.State
	roots       map[string][]string
	derivedBy   map[string][]string
	conflicts   map[string][]string
	unmet       map[string][]string
	active      []string
	blocked     []string
	contested   []string
	missing     []string
	unused      []string
	assumptions []string
	label       string
}

func runForecast(t *testing.T, name string) Answer {
	t.Helper()
	q, m, obs := ForecastScenario(name)
	a, trace, err := Reason(q, m, obs)
	if err != nil {
		t.Fatal(err)
	}
	if len(trace) == 0 {
		t.Fatal("empty trace")
	}
	return Canonical(a)
}

func checkForecast(t *testing.T, a Answer, e forecastExpected) {
	t.Helper()
	eq := func(what string, got, want any) {
		t.Helper()
		if !reflect.DeepEqual(got, want) {
			t.Errorf("%s: got %v, want %v", what, got, want)
		}
	}
	eq("presupposition", a.Presupposition, e.pre)
	eq("presupposition roots", norm(a.PresuppositionRoots), norm(e.preRoots))
	eq("target status", a.TargetStatus, e.target)
	eq("target roots", norm(a.TargetRoots), norm(e.targetRoots))

	got := map[string]DerivedResult{}
	for _, d := range a.Derived {
		got[d.Atom.Key()] = d
	}
	var gotKeys, wantKeys []string
	for k := range got {
		gotKeys = append(gotKeys, k)
	}
	for k := range e.status {
		wantKeys = append(wantKeys, k)
	}
	eq("derived atoms", norm(gotKeys), norm(wantKeys))
	for k, s := range e.status {
		d := got[k]
		eq(k+" status", d.Status, s)
		eq(k+" roots", norm(d.Roots), norm(e.roots[k]))
		eq(k+" derivedBy", norm(d.DerivedBy), norm(e.derivedBy[k]))
		eq(k+" conflicts", norm(d.Conflicts), norm(e.conflicts[k]))
		eq(k+" unmet", norm(d.Unmet), norm(e.unmet[k]))
	}
	eq("active rules", norm(a.ActiveRules), norm(e.active))
	eq("blocked rules", norm(a.BlockedRules), norm(e.blocked))
	eq("contested rules", norm(a.ContestedRules), norm(e.contested))
	eq("missing facts", norm(a.MissingFacts), norm(e.missing))
	eq("unused observations", norm(a.UnusedObservations), norm(e.unused))
	eq("assumptions", norm(a.Assumptions), norm(e.assumptions))
	eq("label", a.Label, e.label)
}

var fS1 = forecastExpected{
	pre: T, preRoots: []string{"sup1"},
	target: T, targetRoots: []string{"erp1", "pln3", "qa4", "sup1"},
	status:      map[string]evidence.State{BE: T, RW: T, AH: T, QH: T, SD: T},
	roots:       map[string][]string{BE: {"erp1", "sup1"}, RW: {"erp1", "sup1"}, AH: {"erp1", "pln3", "sup1"}, QH: {"erp1", "qa4", "sup1"}, SD: {"erp1", "pln3", "qa4", "sup1"}},
	derivedBy:   map[string][]string{BE: {"R1"}, RW: {"R3"}, AH: {"R2"}, QH: {"R4"}, SD: {"R5", "R6"}},
	unmet:       map[string][]string{SD: {"R7"}},
	active:      []string{"R1", "R2", "R3", "R4", "R5", "R6"},
	blocked:     []string{"R7"},
	assumptions: []string{"A1", "A2", "A3"},
	label:       LabelForecastDelayed,
}

// Check 1 (Case 2): sufficient observations over two 3-step branches with shared primary sources.
func TestForecastCheck1Sufficient(t *testing.T) {
	checkForecast(t, runForecast(t, "S1"), fS1)
}

// Check 2 (Case 2): missing data at step 2 stops the chain and reports missing input facts.
func TestForecastCheck2MissingData(t *testing.T) {
	checkForecast(t, runForecast(t, "S2"), forecastExpected{
		pre: T, preRoots: []string{"sup1"},
		target:      U,
		status:      map[string]evidence.State{BE: T, RW: T, AH: U, QH: U, SD: U},
		roots:       map[string][]string{BE: {"erp1", "sup1"}, RW: {"erp1", "sup1"}},
		derivedBy:   map[string][]string{BE: {"R1"}, RW: {"R3"}},
		unmet:       map[string][]string{SD: {"R7"}},
		active:      []string{"R1", "R3"},
		blocked:     []string{"R7"},
		missing:     []string{EA, QS},
		assumptions: []string{"A1", "A2", "A3"},
		label:       LabelForecastUndetermined,
	})
}

// Check 3 (Case 2): contradictory reserve observations propagate conflict down the 3-step chains.
func TestForecastCheck3Conflict(t *testing.T) {
	a := runForecast(t, "S3")
	checkForecast(t, a, forecastExpected{
		pre: T, preRoots: []string{"sup1"},
		target:      U,
		status:      map[string]evidence.State{BE: U, RW: U, AH: U, QH: U, SD: U},
		conflicts:   map[string][]string{BE: {"R1"}, RW: {"R3"}, AH: {"R2"}, QH: {"R4"}, SD: {"R5", "R6", "R7"}},
		contested:   []string{"R1", "R2", "R3", "R4", "R5", "R6", "R7"},
		assumptions: []string{"A1", "A2", "A3"},
		label:       LabelForecastContested,
	})
}

// Check 4 (Case 2): relayed copies and shared primary sources across derived branches are deduplicated.
func TestForecastCheck4CopiesAndSharedRoots(t *testing.T) {
	if a, b := runForecast(t, "S4a"), runForecast(t, "S1"); !reflect.DeepEqual(a, b) {
		t.Errorf("relayed copy changed forecast answer:\n%+v\n%+v", a, b)
	}
	checkForecast(t, runForecast(t, "S4b"), forecastExpected{
		pre: T, preRoots: []string{"sup1"},
		target: T, targetRoots: []string{"erp1", "pln3", "sup1"},
		status:      map[string]evidence.State{BE: T, RW: T, AH: T, QH: T, SD: T},
		roots:       map[string][]string{BE: {"erp1", "sup1"}, RW: {"erp1", "sup1"}, AH: {"erp1", "pln3", "sup1"}, QH: {"erp1", "pln3", "sup1"}, SD: {"erp1", "pln3", "sup1"}},
		derivedBy:   map[string][]string{BE: {"R1"}, RW: {"R3"}, AH: {"R2"}, QH: {"R4"}, SD: {"R5", "R6"}},
		unmet:       map[string][]string{SD: {"R7"}},
		active:      []string{"R1", "R2", "R3", "R4", "R5", "R6"},
		blocked:     []string{"R7"},
		assumptions: []string{"A1", "A2", "A3"},
		label:       LabelForecastDelayed,
	})
}

// Check 5 (Case 2): different query kind over derived statement (S5a) and same input profile with different content (S5b).
func TestForecastCheck5SameNumbersDifferentContent(t *testing.T) {
	a := runForecast(t, "S5a")
	if a.Kind != Whether || a.Presupposition != T || !reflect.DeepEqual(a.PresuppositionRoots, []string{"erp1", "sup1"}) ||
		len(a.Derived) != 0 || len(a.ActiveRules) != 0 || a.Label != LabelWhetherState ||
		!reflect.DeepEqual(a.UnusedObservations, []string{EA, QS}) {
		t.Errorf("S5a Whether(buffer_exhausted) unexpected answer: %+v", a)
	}
	checkForecast(t, runForecast(t, "S5b"), forecastExpected{
		pre: T, preRoots: []string{"sup1"},
		target:      U,
		status:      map[string]evidence.State{BE: T, RW: T, AH: U, QH: U, SD: U},
		roots:       map[string][]string{BE: {"erp1", "sup1"}, RW: {"erp1", "sup1"}},
		derivedBy:   map[string][]string{BE: {"R1"}, RW: {"R3"}},
		unmet:       map[string][]string{SD: {"R7"}},
		active:      []string{"R1", "R3"},
		blocked:     []string{"R2", "R4", "R7"},
		assumptions: []string{"A1", "A3"},
		label:       LabelForecastUndetermined,
	})
	_, _, o1 := ForecastScenario("S1")
	_, _, o5 := ForecastScenario("S5b")
	p1 := characteristicProfile(o1, Answer{})
	p5 := characteristicProfile(o5, Answer{})
	if p1 != p5 {
		t.Fatalf("input profiles differ: %+v vs %+v", p1, p5)
	}
}

func forecastSymbols(q Query, m Model, obs []Observation) []string {
	set := map[string]bool{}
	atom := func(a Atom) {
		if a.Pred == "" {
			return
		}
		set[a.Pred] = true
		for _, x := range a.Args {
			if x[0] != '?' {
				set[x] = true
			}
		}
	}
	atom(q.Target)
	atom(m.Premise)
	for _, d := range m.Derived {
		atom(d)
	}
	for k := range m.Active {
		set[k] = true
	}
	for _, r := range m.Rules {
		set[r.ID] = true
		for _, p := range r.premises() {
			atom(p.Atom)
		}
		atom(r.Then)
		for _, a := range r.Requires {
			set[a] = true
		}
	}
	for _, o := range obs {
		atom(o.Atom)
		set[o.Root], set[o.Relay] = true, true
	}
	out := make([]string, 0, len(set))
	for k := range set {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// Check 6 (Case 2): consistent renaming equivariance and observation-order invariance.
func TestForecastCheck6RenamingAndOrderInvariance(t *testing.T) {
	for _, name := range ForecastScenarioNames {
		q, m, obs := ForecastScenario(name)
		syms := forecastSymbols(q, m, obs)
		r := Renaming{}
		for i, s := range syms {
			r[s] = fmt.Sprintf("z%03d", len(syms)-i)
		}
		a, _, err := Reason(q, m, obs)
		if err != nil {
			t.Fatal(err)
		}
		q2, m2, o2 := r.Input(q, m, obs)
		b, _, err := Reason(q2, m2, o2)
		if err != nil {
			t.Fatal(err)
		}
		if want, got := r.Answer(a), Canonical(b); !reflect.DeepEqual(got, want) {
			t.Errorf("%s: renamed answer differs\n got %+v\nwant %+v", name, got, want)
		}
	}

	for _, name := range []string{"S1", "S3"} {
		q, m, obs := ForecastScenario(name)
		base := runForecast(t, name)
		n := 0
		permutations(append([]Observation(nil), obs...), func(p []Observation) {
			n++
			a, _, err := Reason(q, m, p)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(Canonical(a), base) {
				t.Errorf("%s: permutation %d changed the answer", name, n)
			}
		})
		t.Logf("Forecast %s: %d permutations", name, n)
	}
}

// Check 7 (Case 2): removing intermediate rules, activating assumption A4, and high-reliability irrelevant observation.
func TestForecastCheck7RulesAndAssumptions(t *testing.T) {
	checkForecast(t, runForecast(t, "S7a-1"), forecastExpected{
		pre: T, preRoots: []string{"sup1"},
		target: T, targetRoots: []string{"erp1", "qa4", "sup1"},
		status:      map[string]evidence.State{BE: T, RW: T, AH: U, QH: T, SD: T},
		roots:       map[string][]string{BE: {"erp1", "sup1"}, RW: {"erp1", "sup1"}, QH: {"erp1", "qa4", "sup1"}, SD: {"erp1", "qa4", "sup1"}},
		derivedBy:   map[string][]string{BE: {"R1"}, RW: {"R3"}, QH: {"R4"}, SD: {"R6"}},
		unmet:       map[string][]string{SD: {"R7"}},
		active:      []string{"R1", "R3", "R4", "R6"},
		blocked:     []string{"R7"},
		unused:      []string{EA},
		assumptions: []string{"A1", "A2", "A3"},
		label:       LabelForecastDelayed,
	})

	checkForecast(t, runForecast(t, "S7a-2"), forecastExpected{
		pre: T, preRoots: []string{"sup1"},
		target:      U,
		status:      map[string]evidence.State{BE: T, RW: T, AH: U, QH: U, SD: U},
		roots:       map[string][]string{BE: {"erp1", "sup1"}, RW: {"erp1", "sup1"}},
		derivedBy:   map[string][]string{BE: {"R1"}, RW: {"R3"}},
		unmet:       map[string][]string{SD: {"R7"}},
		active:      []string{"R1", "R3"},
		blocked:     []string{"R7"},
		unused:      []string{EA, QS},
		assumptions: []string{"A1", "A3"},
		label:       LabelForecastUndetermined,
	})

	checkForecast(t, runForecast(t, "S7b-0"), forecastExpected{
		pre: T, preRoots: []string{"sup1"},
		target:      U,
		status:      map[string]evidence.State{BE: U, RW: U, AH: U, QH: U, SD: U},
		blocked:     []string{"R1", "R3"},
		unused:      []string{EA, QS},
		assumptions: []string{"A1"},
		label:       LabelForecastUndetermined,
	})

	checkForecast(t, runForecast(t, "S7b-1"), forecastExpected{
		pre: T, preRoots: []string{"sup1"},
		target: F, targetRoots: []string{"erp1", "sup1"},
		status:      map[string]evidence.State{BE: U, RW: U, AH: U, QH: U, SD: F},
		roots:       map[string][]string{SD: {"erp1", "sup1"}},
		derivedBy:   map[string][]string{SD: {"R7s"}},
		active:      []string{"R7s"},
		blocked:     []string{"R1", "R3"},
		unused:      []string{EA, QS},
		assumptions: []string{"A1", "A4"},
		label:       LabelForecastPreserved,
	})

	c, s := runForecast(t, "S7c"), runForecast(t, "S1")
	if !reflect.DeepEqual(c.UnusedObservations, []string{"courier_color_blue(p)"}) {
		t.Errorf("S7c unused %v", c.UnusedObservations)
	}
	c.UnusedObservations = nil
	if !reflect.DeepEqual(c, s) {
		t.Errorf("high-reliability observation without a rule changed forecast answer")
	}
}

// S8 (Case 2): contested input condition stops propagation under L4 (repeating Limit 1 of experiment 01).
func TestForecastS8ContestedPremise(t *testing.T) {
	checkForecast(t, runForecast(t, "S8"), forecastExpected{
		pre: B, preRoots: []string{"car1", "sup1"},
		target: U,
		unused: []string{EA, QS, RD},
		label:  LabelPremiseContested,
	})
}

func TestForecastTraceAndCycleRejection(t *testing.T) {
	q, m, obs := ForecastScenario("S1")
	_, trace, err := Reason(q, m, obs)
	if err != nil {
		t.Fatal(err)
	}
	kinds := map[StepKind]int{}
	for _, s := range trace {
		if s.Basis == "" || s.Operator == "" {
			t.Errorf("step %d without operator or basis", s.N)
		}
		kinds[s.Kind]++
	}
	for _, k := range []StepKind{HumanInput, Characteristic, Refinement, Deduction, Matching, Selection, Formation} {
		if kinds[k] == 0 {
			t.Errorf("forecast trace has no %s step", k)
		}
	}
	if kinds[Abduction] != 0 {
		t.Errorf("forecast trace unexpectedly invoked PROPOSE (%d times)", kinds[Abduction])
	}

	mCycle := ForecastModel()
	mCycle.Rules = append(mCycle.Rules, Rule{
		ID:       "Rcycle",
		Premises: []Premise{{Atom: A("schedule_delayed", "?p"), Expected: T}},
		Then:     A("buffer_exhausted", "?p"),
		Expected: T,
		Strength: Strong,
	})
	if _, _, err := Reason(q, mCycle, obs); err == nil {
		t.Error("cycle in forward rules accepted")
	}
}

// forecastReference is written directly from protocol 02 tables (§3.2–§3.3) for
// the fixed supply-chain vocabulary with sup1 reporting batch_delayed(p)=TRUE.
// It does not call Reason or its helpers.
func forecastReference(car, rd, ea, qs evidence.State, a4 bool) forecastExpected {
	e := forecastExpected{
		preRoots:  []string{"sup1"},
		status:    map[string]evidence.State{},
		roots:     map[string][]string{},
		derivedBy: map[string][]string{},
		conflicts: map[string][]string{},
		unmet:     map[string][]string{},
	}
	if car != U {
		e.preRoots = []string{"car1", "sup1"}
	}
	if car == F || car == B {
		e.pre = B
		e.unused = []string{EA, QS, RD}
		e.label = LabelPremiseContested
		return e
	}
	e.pre = T
	for _, k := range []string{BE, RW, AH, QH, SD} {
		e.status[k] = U
	}
	baseRoots := append([]string(nil), e.preRoots...)
	with := func(rs []string, extra ...string) []string {
		m := map[string]bool{}
		for _, r := range rs {
			m[r] = true
		}
		for _, r := range extra {
			m[r] = true
		}
		return sortedKeys(m)
	}
	assump := map[string]bool{}

	switch rd {
	case F:
		e.blocked = []string{"R1", "R3"}
		e.unused = []string{EA, QS}
		assump["A1"] = true
		if a4 {
			assump["A4"] = true
			e.active = []string{"R7s"}
			e.status[SD] = F
			e.roots[SD] = with(baseRoots, "erp1")
			e.derivedBy[SD] = []string{"R7s"}
			e.target = F
			e.targetRoots = e.roots[SD]
			e.label = LabelForecastPreserved
		} else {
			e.label = LabelForecastUndetermined
		}
		e.assumptions = sortedKeys(assump)
		return e

	case B:
		assump["A1"], assump["A3"] = true, true
		e.contested = append(e.contested, "R1", "R3", "R7")
		e.conflicts[BE] = []string{"R1"}
		e.conflicts[RW] = []string{"R3"}
		e.conflicts[SD] = append(e.conflicts[SD], "R7")
		if a4 {
			assump["A4"] = true
			e.contested = append(e.contested, "R7s")
			e.conflicts[SD] = append(e.conflicts[SD], "R7s")
		}
		if ea == T {
			e.blocked = append(e.blocked, "R2")
		} else {
			assump["A2"] = true
			e.contested = append(e.contested, "R2", "R5")
			e.conflicts[AH] = []string{"R2"}
			e.conflicts[SD] = append(e.conflicts[SD], "R5")
		}
		if qs == F {
			e.blocked = append(e.blocked, "R4")
		} else {
			assump["A2"] = true
			e.contested = append(e.contested, "R4", "R6")
			e.conflicts[QH] = []string{"R4"}
			e.conflicts[SD] = append(e.conflicts[SD], "R6")
		}
		e.assumptions = sortedKeys(assump)
		e.label = LabelForecastContested
		return e

	case U:
		assump["A1"], assump["A3"] = true, true
		if a4 {
			assump["A4"] = true
		}
		e.missing = append(e.missing, RD)
		switch ea {
		case T:
			e.blocked = append(e.blocked, "R2")
		case B:
			assump["A2"] = true
			e.contested = append(e.contested, "R2", "R5")
			e.conflicts[AH] = []string{"R2"}
			e.conflicts[SD] = append(e.conflicts[SD], "R5")
		case U:
			assump["A2"] = true
			e.missing = append(e.missing, EA)
		case F:
			assump["A2"] = true
		}
		switch qs {
		case F:
			e.blocked = append(e.blocked, "R4")
		case B:
			assump["A2"] = true
			e.contested = append(e.contested, "R4", "R6")
			e.conflicts[QH] = []string{"R4"}
			e.conflicts[SD] = append(e.conflicts[SD], "R6")
		case U:
			assump["A2"] = true
			e.missing = append(e.missing, QS)
		case T:
			assump["A2"] = true
		}
		e.assumptions = sortedKeys(assump)
		if len(e.conflicts[SD]) > 0 {
			e.label = LabelForecastContested
		} else {
			e.label = LabelForecastUndetermined
		}
		return e

	case T:
		assump["A1"], assump["A3"] = true, true
		r1Roots := with(baseRoots, "erp1")
		e.active = append(e.active, "R1", "R3")
		e.status[BE], e.roots[BE], e.derivedBy[BE] = T, r1Roots, []string{"R1"}
		e.status[RW], e.roots[RW], e.derivedBy[RW] = T, r1Roots, []string{"R3"}
		e.blocked = append(e.blocked, "R7")
		e.unmet[SD] = []string{"R7"}
		if a4 {
			e.blocked = append(e.blocked, "R7s")
		}
		sdRoots := map[string]bool{}
		switch ea {
		case F:
			assump["A2"] = true
			e.active = append(e.active, "R2", "R5")
			ahRoots := with(r1Roots, "pln3")
			e.status[AH], e.roots[AH], e.derivedBy[AH] = T, ahRoots, []string{"R2"}
			e.status[SD] |= T
			e.derivedBy[SD] = append(e.derivedBy[SD], "R5")
			for _, r := range ahRoots {
				sdRoots[r] = true
			}
		case T:
			e.blocked = append(e.blocked, "R2")
		case U:
			assump["A2"] = true
			e.missing = append(e.missing, EA)
		case B:
			assump["A2"] = true
			e.contested = append(e.contested, "R2", "R5")
			e.conflicts[AH] = []string{"R2"}
			e.conflicts[SD] = append(e.conflicts[SD], "R5")
		}
		switch qs {
		case T:
			assump["A2"] = true
			e.active = append(e.active, "R4", "R6")
			qhRoots := with(r1Roots, "qa4")
			e.status[QH], e.roots[QH], e.derivedBy[QH] = T, qhRoots, []string{"R4"}
			e.status[SD] |= T
			e.derivedBy[SD] = append(e.derivedBy[SD], "R6")
			for _, r := range qhRoots {
				sdRoots[r] = true
			}
		case F:
			e.blocked = append(e.blocked, "R4")
		case U:
			assump["A2"] = true
			e.missing = append(e.missing, QS)
		case B:
			assump["A2"] = true
			e.contested = append(e.contested, "R4", "R6")
			e.conflicts[QH] = []string{"R4"}
			e.conflicts[SD] = append(e.conflicts[SD], "R6")
		}
		if len(sdRoots) > 0 {
			e.roots[SD] = sortedKeys(sdRoots)
		}
		e.target = e.status[SD]
		e.targetRoots = e.roots[SD]
		e.assumptions = sortedKeys(assump)
		switch {
		case len(e.conflicts[SD]) > 0:
			e.label = LabelForecastContested
		case e.target == T:
			e.label = LabelForecastDelayed
		default:
			e.label = LabelForecastUndetermined
		}
		return e
	}
	panic("unreachable")
}

// Exhaustive 512-case check (§3.6): 4^4 input combinations x 2 assumption modes (A4 off/on).
func TestForecastExhaustiveAgainstReference(t *testing.T) {
	q := ForecastSchedule()
	vals := []evidence.State{U, F, T, B}
	n := 0
	for _, a4 := range []bool{false, true} {
		for _, car := range vals {
			for _, rd := range vals {
				for _, ea := range vals {
					for _, qs := range vals {
						m := ForecastModel()
						if a4 {
							m.Active["A4"] = true
						}
						obs := []Observation{
							SupplyFact("batch_delayed", T, "sup1"),
							SupplyFact("batch_delayed", car, "car1"),
							SupplyFact("reserve_depleted", rd, "erp1"),
							SupplyFact("expedite_available", ea, "pln3"),
							SupplyFact("qa_staff_shortage", qs, "qa4"),
						}
						a, _, err := Reason(q, m, obs)
						if err != nil {
							t.Fatal(err)
						}
						e := forecastReference(car, rd, ea, qs, a4)
						t.Run(fmt.Sprintf("a4=%v,car=%v,rd=%v,ea=%v,qs=%v", a4, car, rd, ea, qs), func(t *testing.T) {
							checkForecast(t, Canonical(a), e)
						})
						n++
					}
				}
			}
		}
	}
	if n != 512 {
		t.Fatalf("enumerated %d cases, want 512", n)
	}
}
