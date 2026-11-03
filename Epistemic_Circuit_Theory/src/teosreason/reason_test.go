package teosreason

import (
	"fmt"
	"reflect"
	"sort"
	"testing"

	evidence "github.com/nick-fedchik/articles/Epistemic_Circuit_Theory/src"
	"github.com/nick-fedchik/articles/Epistemic_Circuit_Theory/src/teoscompose"
)

const (
	HG   = "heat_gen_increased(d)"
	CD   = "cooling_degraded(d)"
	ME1  = "measurement_error(ts1,d)"
	ME2  = "measurement_error(ir2,d)"
	AIR  = "airflow_reduced(d)"
	FAN  = "fan_speed_reduced(d)"
	PWR  = "power_increased(d)"
	LOAD = "load_increased(d)"
	TEMP = "overheat(d)"
)

var (
	U = evidence.Unknown
	T = evidence.True
	F = evidence.False
	B = evidence.Both
)

// expected is written by hand from protocol §6, never copied from program output.
type expected struct {
	pre        evidence.State
	preRoots   []string
	status     map[string]evidence.State
	support    map[string][]string
	refutedBy  map[string][]string
	conflicts  map[string][]string
	unmet      map[string][]string
	compatible []string
	excluded   []string
	contested  []string
	direction  string
	next       []string
	unused     []string
	label      string
}

func run(t *testing.T, name string) Answer {
	t.Helper()
	q, m, obs := Scenario(name)
	a, trace, err := Reason(q, m, obs)
	if err != nil {
		t.Fatal(err)
	}
	if len(trace) == 0 {
		t.Fatal("empty trace")
	}
	return Canonical(a)
}

func nextKeys(a Answer) []string {
	var out []string
	for _, c := range a.NextChecks {
		k := c.Atom.Key()
		if len(c.ExcludedRoots) > 0 {
			k += fmt.Sprintf(" без %v", c.ExcludedRoots)
		}
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

func norm(xs []string) []string {
	if len(xs) == 0 {
		return nil
	}
	out := append([]string(nil), xs...)
	sort.Strings(out)
	return out
}

func check(t *testing.T, a Answer, e expected) {
	t.Helper()
	eq := func(what string, got, want any) {
		t.Helper()
		if !reflect.DeepEqual(got, want) {
			t.Errorf("%s: got %v, want %v", what, got, want)
		}
	}
	eq("presupposition", a.Presupposition, e.pre)
	eq("presupposition roots", norm(a.PresuppositionRoots), norm(e.preRoots))
	got := map[string]HypothesisResult{}
	for _, h := range a.Hypotheses {
		got[h.Hypothesis.Key()] = h
	}
	var gotKeys, wantKeys []string
	for k := range got {
		gotKeys = append(gotKeys, k)
	}
	for k := range e.status {
		wantKeys = append(wantKeys, k)
	}
	eq("hypotheses", norm(gotKeys), norm(wantKeys))
	for k, s := range e.status {
		h := got[k]
		eq(k+" status", h.Status, s)
		eq(k+" support", norm(h.SupportRoots), norm(e.support[k]))
		eq(k+" refutedBy", norm(h.RefutedBy), norm(e.refutedBy[k]))
		eq(k+" conflicts", norm(h.Conflicts), norm(e.conflicts[k]))
		eq(k+" unmet", norm(h.Unmet), norm(e.unmet[k]))
	}
	eq("compatible", norm(a.Compatible), norm(e.compatible))
	eq("excluded", norm(a.Excluded), norm(e.excluded))
	eq("contested", norm(a.Contested), norm(e.contested))
	eq("direction", a.Direction, e.direction)
	eq("next", norm(nextKeys(a)), norm(e.next))
	eq("unused", norm(a.UnusedObservations), norm(e.unused))
	eq("label", a.Label, e.label)
}

var s1 = expected{
	pre: T, preRoots: []string{"ir2", "ts1"},
	status:     map[string]evidence.State{HG: U, CD: T, ME1: F, ME2: F},
	support:    map[string][]string{CD: {"an3"}},
	refutedBy:  map[string][]string{ME1: {"R5"}, ME2: {"R5"}},
	unmet:      map[string][]string{HG: {"R3"}},
	compatible: []string{CD, HG}, excluded: []string{ME1, ME2},
	direction: CD, next: []string{FAN, LOAD}, label: LabelWhy,
}

// Check 1: sufficient observations give an answer with the right status.
func TestCheck1Sufficient(t *testing.T) {
	a := run(t, "S1")
	check(t, a, s1)
	for _, h := range a.Hypotheses {
		for _, l := range h.Links {
			if h.Hypothesis.Key() == ME1 && l.Rule == "R5" && !reflect.DeepEqual(l.Roots, []string{"ir2"}) {
				t.Errorf("ME(ts1) refuted by roots %v, want [ir2]", l.Roots)
			}
		}
	}
}

// Check 2: missing data leaves several explanations and no hidden choice.
func TestCheck2MissingData(t *testing.T) {
	check(t, run(t, "S2"), expected{
		pre: T, preRoots: []string{"ts1"},
		status:     map[string]evidence.State{HG: U, CD: U, ME1: U},
		compatible: []string{CD, HG, ME1},
		next:       []string{AIR, FAN, PWR, LOAD, TEMP + " без [ts1]"},
		label:      LabelNoDirection,
	})
}

// Check 3: contradictory observations are kept and reported, not erased.
func TestCheck3Conflict(t *testing.T) {
	a := run(t, "S3")
	check(t, a, expected{
		pre: T, preRoots: []string{"ir2", "ts1"},
		status:     map[string]evidence.State{HG: U, CD: U, ME1: F, ME2: F},
		refutedBy:  map[string][]string{ME1: {"R5"}, ME2: {"R5"}},
		conflicts:  map[string][]string{CD: {"R1"}},
		unmet:      map[string][]string{HG: {"R3"}},
		compatible: []string{CD, HG}, excluded: []string{ME1, ME2}, contested: []string{CD},
		next:  []string{FAN, LOAD},
		label: LabelNoDirection,
	})
	for _, h := range a.Hypotheses {
		for _, l := range h.Links {
			if l.Rule == "R1" && (l.Observed != B || !reflect.DeepEqual(l.Roots, []string{"an3", "an5"})) {
				t.Errorf("R1 link observed %v roots %v, want BOTH [an3 an5]", l.Observed, l.Roots)
			}
		}
	}
}

// Check 4: copies of one primary source are not independent observations.
func TestCheck4Copies(t *testing.T) {
	if a, b := run(t, "S4a"), run(t, "S1"); !reflect.DeepEqual(a, b) {
		t.Errorf("relayed copy changed the answer:\n%+v\n%+v", a, b)
	}
	check(t, run(t, "S4b"), expected{
		pre: T, preRoots: []string{"ts1"},
		status:     map[string]evidence.State{HG: U, CD: T, ME1: U},
		support:    map[string][]string{CD: {"an3"}},
		unmet:      map[string][]string{HG: {"R3"}},
		compatible: []string{CD, HG, ME1},
		direction:  CD, next: []string{FAN, LOAD, TEMP + " без [ts1]"}, label: LabelWhy,
	})
}

type profile struct {
	observations, roots int
	values              string
	statuses            string
}

// characteristicProfile is what a rating-only processor would see: counts and
// values without statement identity.
func characteristicProfile(obs []Observation, a Answer) profile {
	roots := map[string]bool{}
	var vals, sts []string
	for _, o := range obs {
		roots[o.Root] = true
		vals = append(vals, o.Value.String())
	}
	for _, h := range a.Hypotheses {
		sts = append(sts, fmt.Sprintf("%s/%d", h.Status, len(h.SupportRoots)))
	}
	sort.Strings(vals)
	sort.Strings(sts)
	return profile{len(obs), len(roots), fmt.Sprint(vals), fmt.Sprint(sts)}
}

// Check 5: same numeric characteristics, different content, different answer.
func TestCheck5SameNumbersDifferentContent(t *testing.T) {
	a := run(t, "S5a")
	if a.Kind != Whether || a.Presupposition != T || !reflect.DeepEqual(a.PresuppositionRoots, []string{"ir2", "ts1"}) ||
		len(a.Hypotheses) != 0 || a.Direction != "" || a.Label != LabelWhetherState {
		t.Errorf("Whether answer: %+v", a)
	}
	check(t, run(t, "S5b"), expected{
		pre: T, preRoots: []string{"ir2", "ts1"},
		status:     map[string]evidence.State{HG: T, CD: U, ME1: F, ME2: F},
		support:    map[string][]string{HG: {"wm4"}},
		refutedBy:  map[string][]string{ME1: {"R5"}, ME2: {"R5"}},
		unmet:      map[string][]string{CD: {"R1"}},
		compatible: []string{CD, HG}, excluded: []string{ME1, ME2},
		direction: HG, next: []string{FAN, LOAD}, label: LabelWhy,
	})
	_, _, o1 := Scenario("S1")
	_, _, o5 := Scenario("S5b")
	p1, p5 := characteristicProfile(o1, run(t, "S1")), characteristicProfile(o5, run(t, "S5b"))
	if p1 != p5 {
		t.Fatalf("profiles differ, check is not discriminating: %+v vs %+v", p1, p5)
	}
}

func symbols(q Query, m Model, obs []Observation) []string {
	set := map[string]bool{}
	atom := func(a Atom) {
		set[a.Pred] = true
		for _, x := range a.Args {
			if x[0] != '?' {
				set[x] = true
			}
		}
	}
	atom(q.Target)
	for k := range m.Active {
		set[k] = true
	}
	for _, s := range m.Schemes {
		set[s.ID] = true
		atom(s.Hypothesis)
		atom(s.Effect)
	}
	for _, r := range m.Rules {
		set[r.ID] = true
		atom(r.If)
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

// Check 6: consistent renaming changes the answer only by the same renaming.
// The bijection reverses alphabetical order to expose name-based tie-breaking.
func TestCheck6RenamingEquivariance(t *testing.T) {
	for _, name := range ScenarioNames {
		q, m, obs := Scenario(name)
		syms := symbols(q, m, obs)
		r := Renaming{}
		for i, s := range syms {
			r[s] = fmt.Sprintf("s%03d", len(syms)-i)
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
}

// Added after the mutation analysis, not preregistered: S1–S8 contain no tie
// between supported hypotheses, so name-based tie-breaking escaped check 6.
func TestCheck6RenamingOnTiePostHoc(t *testing.T) {
	q, m := WhyOverheat(), ProtocolModel()
	obs := []Observation{
		Temperature(90, 94, "ts1", "ts1"), Temperature(88, 94, "ir2", "ir2"),
		Fact("airflow_reduced", T, "an3"), Fact("power_increased", T, "wm4"),
	}
	a, _, err := Reason(q, m, obs)
	if err != nil {
		t.Fatal(err)
	}
	if a.Direction != "" || a.Label != LabelNoDirection {
		t.Fatalf("tie produced direction %q", a.Direction)
	}
	syms := symbols(q, m, obs)
	r := Renaming{}
	for i, s := range syms {
		r[s] = fmt.Sprintf("s%03d", len(syms)-i)
	}
	q2, m2, o2 := r.Input(q, m, obs)
	b, _, err := Reason(q2, m2, o2)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(Canonical(b), r.Answer(a)) {
		t.Error("renaming changed the tie outcome")
	}
}

func permutations(xs []Observation, f func([]Observation)) {
	var rec func(int)
	rec = func(k int) {
		if k == len(xs) {
			f(append([]Observation(nil), xs...))
			return
		}
		for i := k; i < len(xs); i++ {
			xs[k], xs[i] = xs[i], xs[k]
			rec(k + 1)
			xs[k], xs[i] = xs[i], xs[k]
		}
	}
	rec(0)
}

func TestCheck6ObservationOrderInvariance(t *testing.T) {
	for _, name := range []string{"S1", "S3"} {
		q, m, obs := Scenario(name)
		base := run(t, name)
		n := 0
		permutations(append([]Observation(nil), obs...), func(p []Observation) {
			n++
			a, _, err := Reason(q, m, p)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(Canonical(a), base) {
				t.Errorf("%s: order %d changed the answer", name, n)
			}
		})
		t.Logf("%s: %d permutations", name, n)
	}
}

// Check 7: a rule or observation change alters a traceable transition; a high
// rating without a rule creates none.
func TestCheck7RulesGovernTransitions(t *testing.T) {
	check(t, run(t, "S7a"), expected{
		pre: T, preRoots: []string{"ir2", "ts1"},
		status:     map[string]evidence.State{HG: U, CD: U, ME1: F, ME2: F},
		refutedBy:  map[string][]string{ME1: {"R5"}, ME2: {"R5"}},
		unmet:      map[string][]string{HG: {"R3"}},
		compatible: []string{CD, HG}, excluded: []string{ME1, ME2},
		next: []string{FAN, LOAD}, unused: []string{AIR}, label: LabelNoDirection,
	})
	b := run(t, "S7b")
	check(t, b, expected{
		pre: T, preRoots: []string{"ir2", "ts1"},
		status:     map[string]evidence.State{HG: F, CD: T, ME1: F, ME2: F},
		support:    map[string][]string{CD: {"an3"}},
		refutedBy:  map[string][]string{HG: {"R3s"}, ME1: {"R5"}, ME2: {"R5"}},
		unmet:      map[string][]string{HG: {"R3"}},
		compatible: []string{CD}, excluded: []string{HG, ME1, ME2},
		direction: CD, next: []string{FAN}, label: LabelWhy,
	})
	if !reflect.DeepEqual(b.Assumptions, []string{"A1", "A2", "A3", "A4"}) {
		t.Errorf("S7b assumptions %v", b.Assumptions)
	}
	c, s := run(t, "S7c"), run(t, "S1")
	if !reflect.DeepEqual(c.UnusedObservations, []string{"fan_noise_changed(d)"}) {
		t.Errorf("S7c unused %v", c.UnusedObservations)
	}
	c.UnusedObservations = nil
	if !reflect.DeepEqual(c, s) {
		t.Errorf("high-reliability observation without a rule changed the answer")
	}
}

// S8: the unchanged IntervalGuard reports overlap as BOTH; the protocol maps it
// to UNKNOWN before reasoning.
func TestS8IntervalAdaptation(t *testing.T) {
	_, _, obs := Scenario("S8")
	if obs[1].Value != U {
		t.Fatalf("ir2 mapped to %v, want UNKNOWN", obs[1].Value)
	}
	m, err := evidenceMeasurement(83, 89)
	if err != nil {
		t.Fatal(err)
	}
	if _, raw, err := ObserveInterval(A("overheat", "d"), m, LimitWindow, Unit, "ir2", "ir2"); err != nil || raw != B {
		t.Fatalf("raw guard %v err %v, want BOTH", raw, err)
	}
	check(t, run(t, "S8"), expected{
		pre: T, preRoots: []string{"ts1"},
		status:     map[string]evidence.State{HG: U, CD: T, ME1: U},
		support:    map[string][]string{CD: {"an3"}},
		unmet:      map[string][]string{HG: {"R3"}},
		compatible: []string{CD, HG, ME1},
		direction:  CD, next: []string{FAN, LOAD, TEMP + " без [ts1]"}, label: LabelWhy,
	})
}

func evidenceMeasurement(lo, hi float64) (teoscompose.Measurement, error) {
	return teoscompose.NewMeasurement((lo+hi)/2, evidence.Interval{Low: lo, High: hi}, Unit)
}

func TestTraceMarksEveryContentStep(t *testing.T) {
	q, m, obs := Scenario("S1")
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
	for _, k := range []StepKind{HumanInput, Characteristic, Refinement, Abduction, Deduction, Matching, Selection, Formation} {
		if kinds[k] == 0 {
			t.Errorf("trace has no %s step", k)
		}
	}
}

func TestValidationRejectsMalformedInput(t *testing.T) {
	q, m, obs := Scenario("S1")
	bad := append([]Observation(nil), obs...)
	bad[0].Root = ""
	if _, _, err := Reason(q, m, bad); err == nil {
		t.Error("observation without primary source accepted")
	}
	m2 := ProtocolModel()
	m2.Rules = append(m2.Rules, Rule{ID: "Rx", If: A("cooling_degraded", "?d"), Then: A("x", "?z"), Expected: T})
	if _, _, err := Reason(q, m2, obs); err == nil {
		t.Error("rule with unbound variable accepted")
	}
	if _, _, err := Reason(Query{Why, A("overheat", "?d")}, m, obs); err == nil {
		t.Error("non-ground query accepted")
	}
}

// reference is written directly from protocol tables §3–§5 for the fixed
// vocabulary with ts1=TRUE. It does not call Reason or its helpers.
func reference(ir, air, fan, pwr, load evidence.State) expected {
	e := expected{status: map[string]evidence.State{}, support: map[string][]string{},
		refutedBy: map[string][]string{}, conflicts: map[string][]string{}, unmet: map[string][]string{}}
	e.preRoots = []string{"ts1"}
	if ir != U {
		e.preRoots = append(e.preRoots, "ir2")
	}
	switch ir {
	case F, B:
		e.pre, e.label = B, LabelContested
		return e
	}
	e.pre = T
	typical := func(h string, obs evidence.State, rule, atom, root string) {
		switch obs {
		case T:
			e.support[h] = append(e.support[h], root)
		case F:
			e.unmet[h] = append(e.unmet[h], rule)
		case B:
			e.conflicts[h] = append(e.conflicts[h], rule)
		case U:
			e.next = append(e.next, atom)
		}
	}
	typical(CD, air, "R1", AIR, "an3")
	typical(CD, fan, "R2", FAN, "fs6")
	typical(HG, pwr, "R3", PWR, "wm4")
	typical(HG, load, "R4", LOAD, "ld7")
	for _, h := range []string{CD, HG} {
		e.status[h] = U
		if len(e.support[h]) > 0 {
			e.status[h] = T
		}
		e.compatible = append(e.compatible, h)
		if len(e.conflicts[h]) > 0 {
			e.contested = append(e.contested, h)
		}
	}
	if ir == T {
		e.status[ME1], e.status[ME2] = F, F
		e.refutedBy[ME1], e.refutedBy[ME2] = []string{"R5"}, []string{"R5"}
		e.excluded = []string{ME1, ME2}
	} else {
		e.status[ME1] = U
		e.compatible = append(e.compatible, ME1)
		e.next = append(e.next, TEMP+" без [ts1]")
	}
	clean := func(h string) int {
		if e.status[h] != T || len(e.conflicts[h]) > 0 {
			return -1
		}
		return len(e.support[h])
	}
	switch c, h := clean(CD), clean(HG); {
	case c > h:
		e.direction = CD
	case h > c:
		e.direction = HG
	}
	e.label = LabelNoDirection
	if e.direction != "" {
		e.label = LabelWhy
	}
	return e
}

func TestExhaustiveAgainstReference(t *testing.T) {
	q, m := WhyOverheat(), ProtocolModel()
	vals := []evidence.State{U, F, T, B}
	n := 0
	for _, ir := range vals {
		for _, air := range vals {
			for _, fan := range vals {
				for _, pwr := range vals {
					for _, load := range vals {
						obs := []Observation{
							Fact("overheat", T, "ts1"), Fact("overheat", ir, "ir2"),
							Fact("airflow_reduced", air, "an3"), Fact("fan_speed_reduced", fan, "fs6"),
							Fact("power_increased", pwr, "wm4"), Fact("load_increased", load, "ld7"),
						}
						a, _, err := Reason(q, m, obs)
						if err != nil {
							t.Fatal(err)
						}
						e := reference(ir, air, fan, pwr, load)
						if e.pre != T {
							e.unused = []string{AIR, FAN, LOAD, PWR}
						}
						t.Run(fmt.Sprintf("ir=%v,air=%v,fan=%v,pwr=%v,load=%v", ir, air, fan, pwr, load), func(t *testing.T) {
							check(t, Canonical(a), e)
						})
						n++
					}
				}
			}
		}
	}
	if n != 1024 {
		t.Fatalf("enumerated %d cases, want 1024", n)
	}
}
