package teosreason

import (
	"sort"

	evidence "github.com/nick-fedchik/articles/Epistemic_Circuit_Theory/src"
	"github.com/nick-fedchik/articles/Epistemic_Circuit_Theory/src/teoscompose"
)

// LimitWindow is the teaching overheating window [85, 1000] degC of protocol §2.
var LimitWindow = evidence.Interval{Low: 85, High: 1000}

const Unit = "degC"

// ProtocolModel returns schemes K1–K3 and rules R1–R5 of protocol §3 with
// assumptions A1–A3 active and A4 inactive.
func ProtocolModel() Model {
	return Model{
		Schemes: []Scheme{
			{ID: "K1", Hypothesis: A("heat_gen_increased", "?d"), Effect: A("overheat", "?d")},
			{ID: "K2", Hypothesis: A("cooling_degraded", "?d"), Effect: A("overheat", "?d")},
			{ID: "K3", Hypothesis: A("measurement_error", "?i", "?d"), Effect: A("overheat", "?d"), PerReportingRoot: "?i"},
		},
		Rules: []Rule{
			{ID: "R1", If: A("cooling_degraded", "?d"), Then: A("airflow_reduced", "?d"), Expected: evidence.True, Strength: Typical, Requires: []string{"A1"}},
			{ID: "R2", If: A("cooling_degraded", "?d"), Then: A("fan_speed_reduced", "?d"), Expected: evidence.True, Strength: Typical, Requires: []string{"A1"}},
			{ID: "R3", If: A("heat_gen_increased", "?d"), Then: A("power_increased", "?d"), Expected: evidence.True, Strength: Typical},
			{ID: "R3s", If: A("heat_gen_increased", "?d"), Then: A("power_increased", "?d"), Expected: evidence.True, Strength: Strong, Requires: []string{"A4"}},
			{ID: "R4", If: A("heat_gen_increased", "?d"), Then: A("load_increased", "?d"), Expected: evidence.True, Strength: Typical},
			{ID: "R5", If: A("measurement_error", "?i", "?d"), Then: A("overheat", "?d"), Expected: evidence.False, Strength: Strong, Requires: []string{"A2", "A3"}, ExcludeRootVar: "?i"},
		},
		Active: map[string]bool{"A1": true, "A2": true, "A3": true},
	}
}

func WhyOverheat() Query     { return Query{Why, A("overheat", "d")} }
func WhetherOverheat() Query { return Query{Whether, A("overheat", "d")} }

// Temperature converts one interval reading with ObserveInterval and panics
// only on programmer error in fixed scenario data.
func Temperature(lo, hi float64, relay, root string) Observation {
	m, err := teoscompose.NewMeasurement((lo+hi)/2, evidence.Interval{Low: lo, High: hi}, Unit)
	if err != nil {
		panic(err)
	}
	o, _, err := ObserveInterval(A("overheat", "d"), m, LimitWindow, Unit, relay, root)
	if err != nil {
		panic(err)
	}
	return o
}

func Fact(pred string, v evidence.State, root string) Observation {
	return Observation{Atom: A(pred, "d"), Value: v, Relay: root, Root: root}
}

// Scenario returns the observations and model of protocol §6.
func Scenario(name string) (Query, Model, []Observation) {
	m := ProtocolModel()
	s1 := []Observation{
		Temperature(90, 94, "ts1", "ts1"),
		Temperature(88, 94, "ir2", "ir2"),
		Fact("airflow_reduced", evidence.True, "an3"),
		Fact("power_increased", evidence.False, "wm4"),
	}
	switch name {
	case "S1":
		return WhyOverheat(), m, s1
	case "S2":
		return WhyOverheat(), m, s1[:1]
	case "S3":
		return WhyOverheat(), m, append(clone(s1), Fact("airflow_reduced", evidence.False, "an5"))
	case "S4a":
		cp := Fact("airflow_reduced", evidence.True, "an3")
		cp.Relay = "dashboard"
		return WhyOverheat(), m, append(clone(s1), cp)
	case "S4b":
		return WhyOverheat(), m, []Observation{
			Temperature(90, 94, "ts1", "ts1"),
			Temperature(90, 94, "scada", "ts1"),
			s1[2], s1[3],
		}
	case "S5a":
		return WhetherOverheat(), m, s1
	case "S5b":
		return WhyOverheat(), m, []Observation{
			s1[0], s1[1],
			Fact("airflow_reduced", evidence.False, "an3"),
			Fact("power_increased", evidence.True, "wm4"),
		}
	case "S7a":
		var rules []Rule
		for _, r := range m.Rules {
			if r.ID != "R1" {
				rules = append(rules, r)
			}
		}
		m.Rules = rules
		return WhyOverheat(), m, s1
	case "S7b":
		m.Active["A4"] = true
		return WhyOverheat(), m, s1
	case "S7c":
		noise := Fact("fan_noise_changed", evidence.True, "op6")
		noise.Reliability = 1.0
		return WhyOverheat(), m, append(clone(s1), noise)
	case "S8":
		return WhyOverheat(), m, []Observation{
			Temperature(90, 94, "ts1", "ts1"),
			Temperature(83, 89, "ir2", "ir2"),
			s1[2], s1[3],
		}
	}
	panic("unknown scenario " + name)
}

var ScenarioNames = []string{"S1", "S2", "S3", "S4a", "S4b", "S5a", "S5b", "S7a", "S7b", "S7c", "S8"}

func clone(xs []Observation) []Observation { return append([]Observation(nil), xs...) }

// Renaming is a bijection over all non-variable symbols: predicates, constants,
// primary sources, relays, scheme, rule and assumption identifiers.
type Renaming map[string]string

func (r Renaming) s(x string) string {
	if v, ok := r[x]; ok {
		return v
	}
	return x
}

func (r Renaming) ss(xs []string) []string {
	if xs == nil {
		return nil
	}
	out := make([]string, len(xs))
	for i, x := range xs {
		out[i] = r.s(x)
	}
	return out
}

func (r Renaming) Atom(a Atom) Atom { return Atom{r.s(a.Pred), r.ss(a.Args)} }

// Key renames a ground or pattern atom key produced by Atom.Key.
func (r Renaming) key(k string, atoms map[string]Atom) string {
	if a, ok := atoms[k]; ok {
		return r.Atom(a).Key()
	}
	return r.s(k)
}

func (r Renaming) Input(q Query, m Model, obs []Observation) (Query, Model, []Observation) {
	q2 := Query{q.Kind, r.Atom(q.Target)}
	m2 := Model{Active: map[string]bool{}}
	for k, v := range m.Active {
		m2.Active[r.s(k)] = v
	}
	for _, s := range m.Schemes {
		m2.Schemes = append(m2.Schemes, Scheme{r.s(s.ID), r.Atom(s.Hypothesis), r.Atom(s.Effect), s.PerReportingRoot})
	}
	for _, x := range m.Rules {
		m2.Rules = append(m2.Rules, Rule{r.s(x.ID), r.Atom(x.If), r.Atom(x.Then), x.Expected, x.Strength, r.ss(x.Requires), x.ExcludeRootVar})
	}
	var o2 []Observation
	for _, o := range obs {
		o2 = append(o2, Observation{r.Atom(o.Atom), o.Value, r.s(o.Relay), r.s(o.Root), o.Reliability})
	}
	return q2, m2, o2
}

// Answer renames an answer; hypothesis and atom keys are renamed structurally.
func (r Renaming) Answer(a Answer) Answer {
	atoms := map[string]Atom{}
	for _, h := range a.Hypotheses {
		atoms[h.Hypothesis.Key()] = h.Hypothesis
		for _, l := range h.Links {
			atoms[l.Atom.Key()] = l.Atom
		}
	}
	keys := func(xs []string) []string {
		if xs == nil {
			return nil
		}
		out := make([]string, len(xs))
		for i, x := range xs {
			out[i] = r.key(x, atoms)
		}
		return out
	}
	b := a
	b.Target = r.Atom(a.Target)
	b.PresuppositionRoots = r.ss(a.PresuppositionRoots)
	b.Hypotheses = nil
	for _, h := range a.Hypotheses {
		h2 := h
		h2.Hypothesis = r.Atom(h.Hypothesis)
		h2.Scheme = r.s(h.Scheme)
		h2.SupportRoots, h2.RefuteRoots = r.ss(h.SupportRoots), r.ss(h.RefuteRoots)
		h2.RefutedBy, h2.Conflicts, h2.Unmet = r.ss(h.RefutedBy), r.ss(h.Conflicts), r.ss(h.Unmet)
		h2.Links = nil
		for _, l := range h.Links {
			l2 := l
			l2.Hypothesis = r.key(l.Hypothesis, atoms)
			l2.Rule = r.s(l.Rule)
			l2.Atom = r.Atom(l.Atom)
			l2.ExcludedRoots, l2.Assumptions, l2.Roots = r.ss(l.ExcludedRoots), r.ss(l.Assumptions), r.ss(l.Roots)
			h2.Links = append(h2.Links, l2)
		}
		b.Hypotheses = append(b.Hypotheses, h2)
	}
	b.Compatible, b.Excluded, b.Contested = keys(a.Compatible), keys(a.Excluded), keys(a.Contested)
	if a.Direction != "" {
		b.Direction = r.key(a.Direction, atoms)
	}
	b.NextChecks = nil
	for _, c := range a.NextChecks {
		b.NextChecks = append(b.NextChecks, Check{r.key(c.Hypothesis, atoms), r.s(c.Rule), r.Atom(c.Atom), c.Expected, r.ss(c.ExcludedRoots)})
	}
	b.UnusedObservations = nil
	for _, u := range a.UnusedObservations {
		b.UnusedObservations = append(b.UnusedObservations, r.renameKeyText(u))
	}
	b.Assumptions = r.ss(a.Assumptions)
	return Canonical(b)
}

// renameKeyText renames an unused-observation key "pred(a,b)" token-wise.
func (r Renaming) renameKeyText(k string) string {
	open := -1
	for i, c := range k {
		if c == '(' {
			open = i
			break
		}
	}
	if open < 0 || k[len(k)-1] != ')' {
		return r.s(k)
	}
	var args []string
	if inner := k[open+1 : len(k)-1]; inner != "" {
		start := 0
		for i := 0; i <= len(inner); i++ {
			if i == len(inner) || inner[i] == ',' {
				args = append(args, inner[start:i])
				start = i + 1
			}
		}
	}
	return Atom{r.s(k[:open]), r.ss(args)}.Key()
}

func sortedCopy(xs []string) []string {
	if xs == nil {
		return nil
	}
	out := append([]string(nil), xs...)
	sort.Strings(out)
	return out
}

// Canonical sorts every set-valued field so answers compare as sets.
func Canonical(a Answer) Answer {
	b := a
	for _, xs := range []*[]string{&b.PresuppositionRoots, &b.Compatible, &b.Excluded, &b.Contested, &b.UnusedObservations, &b.Assumptions} {
		*xs = sortedCopy(*xs)
	}
	hs := append([]HypothesisResult(nil), a.Hypotheses...)
	for i := range hs {
		h := &hs[i]
		for _, xs := range []*[]string{&h.SupportRoots, &h.RefuteRoots, &h.RefutedBy, &h.Conflicts, &h.Unmet} {
			*xs = sortedCopy(*xs)
		}
		h.Links = append([]Link(nil), h.Links...)
		for j := range h.Links {
			l := &h.Links[j]
			l.Roots, l.ExcludedRoots, l.Assumptions = sortedCopy(l.Roots), sortedCopy(l.ExcludedRoots), sortedCopy(l.Assumptions)
		}
		sort.Slice(h.Links, func(i, j int) bool {
			return h.Links[i].Rule+"|"+h.Links[i].Atom.Key() < h.Links[j].Rule+"|"+h.Links[j].Atom.Key()
		})
	}
	sort.Slice(hs, func(i, j int) bool { return hs[i].Hypothesis.Key() < hs[j].Hypothesis.Key() })
	b.Hypotheses = hs
	cs := append([]Check(nil), a.NextChecks...)
	sort.Slice(cs, func(i, j int) bool {
		return cs[i].Hypothesis+"|"+cs[i].Rule+"|"+cs[i].Atom.Key() < cs[j].Hypothesis+"|"+cs[j].Rule+"|"+cs[j].Atom.Key()
	})
	b.NextChecks = cs
	return b
}
