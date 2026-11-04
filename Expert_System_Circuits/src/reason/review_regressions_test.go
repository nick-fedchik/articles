package reason

import (
	"strings"
	"testing"

	evidence "github.com/nick-fedchik/articles/Expert_System_Circuits/src"
)

func ob(atom Atom, v evidence.State, root string) Observation {
	return Observation{Atom: atom, Value: v, Relay: root, Root: root}
}

func TestAtomKeyIsInjective(t *testing.T) {
	a, b := A("p", "a", "b"), A("p", "a,b")
	if a.Key() == b.Key() {
		t.Fatalf("distinct atoms share key %q", a.Key())
	}
	ans, _, err := Reason(Query{Kind: Whether, Target: a}, Model{}, []Observation{ob(b, evidence.True, "sensor")})
	if err != nil {
		t.Fatal(err)
	}
	if ans.Presupposition != evidence.Unknown {
		t.Fatalf("observing %s changed %s to %s", b.Key(), a.Key(), ans.Presupposition)
	}
	for _, x := range []Atom{a, b, A("q"), A("r", `x\`, "(y)", ""), A(`a(b`, ",")} {
		got, ok := parseKey(x.Key())
		if !ok || got.Key() != x.Key() {
			t.Fatalf("parseKey(%q) = %v, %v", x.Key(), got, ok)
		}
	}
}

func TestWhetherJoinsObservationWithRuleConclusion(t *testing.T) {
	goal, condition := A("goal", "d"), A("condition", "d")
	m := Model{Premise: condition, Rules: []Rule{{
		ID: "r", Premises: []Premise{{Atom: condition, Expected: evidence.True}},
		Then: goal, Expected: evidence.False, Strength: Strong,
	}}}
	in := []Observation{ob(goal, evidence.True, "s1"), ob(condition, evidence.True, "s2")}
	w, _, err := Reason(Query{Kind: Whether, Target: goal}, m, in)
	if err != nil {
		t.Fatal(err)
	}
	f, _, err := Reason(Query{Kind: Forecast, Target: goal}, m, in)
	if err != nil {
		t.Fatal(err)
	}
	if w.Presupposition != evidence.Both || f.TargetStatus != evidence.Both {
		t.Fatalf("Whether=%s Forecast=%s, want BOTH for both", w.Presupposition, f.TargetStatus)
	}
}

func TestObservedIntermediateSatisfiesDownstreamRule(t *testing.T) {
	start, mid, target, failed := A("start"), A("intermediate"), A("target"), A("failed")
	m := Model{Premise: start, Rules: []Rule{
		{ID: "r1", Premises: []Premise{{Atom: start, Expected: evidence.True}, {Atom: failed, Expected: evidence.True}},
			Then: mid, Expected: evidence.True, Strength: Strong},
		{ID: "r2", Premises: []Premise{{Atom: start, Expected: evidence.True}, {Atom: mid, Expected: evidence.True}},
			Then: target, Expected: evidence.True, Strength: Strong},
	}}
	f, _, err := Reason(Query{Kind: Forecast, Target: target}, m, []Observation{
		ob(start, evidence.True, "start-s"), ob(failed, evidence.False, "fail-s"), ob(mid, evidence.True, "observed-b")})
	if err != nil {
		t.Fatal(err)
	}
	if f.TargetStatus != evidence.True {
		t.Fatalf("target=%s, want TRUE (active=%v blocked=%v)", f.TargetStatus, f.ActiveRules, f.BlockedRules)
	}
}

func TestSwappedVariablesBindByUnification(t *testing.T) {
	target := A("result", "a", "b")
	m := Model{Premise: A("go", "?y"), Rules: []Rule{{
		ID: "reverse", Premises: []Premise{{Atom: A("link", "?x", "?y"), Expected: evidence.True}},
		Then: A("result", "?y", "?x"), Expected: evidence.True, Strength: Strong,
	}}}
	f, _, err := Reason(Query{Kind: Forecast, Target: target}, m, []Observation{
		ob(A("go", "a"), evidence.True, "s1"), ob(A("link", "b", "a"), evidence.True, "s2")})
	if err != nil {
		t.Fatal(err)
	}
	if f.TargetStatus != evidence.True {
		t.Fatalf("result(a,b)=%s, want TRUE; derived=%v", f.TargetStatus, f.Derived)
	}
}

func TestCycleRejectedAtModelAcceptance(t *testing.T) {
	q, m, in := ForecastScenario("S1")
	m.Rules = append(m.Rules, Rule{
		ID: "cycle", Premises: []Premise{{Atom: A("schedule_delayed", "?p"), Expected: evidence.True}},
		Then: A("buffer_exhausted", "?p"), Expected: evidence.True, Strength: Strong,
	})
	for name, obs := range map[string][]Observation{"active": in, "no-observations": nil} {
		if _, _, err := Reason(q, m, obs); err == nil || !strings.Contains(err.Error(), "cycle") {
			t.Fatalf("%s: err=%v, want cycle rejection", name, err)
		}
	}
}

func TestOverheatWindowIsClosedAndBounded(t *testing.T) {
	cases := map[float64]evidence.State{85: evidence.True, 1000: evidence.True, 1001: evidence.False, 84.9: evidence.False}
	for temp, want := range cases {
		if got := Temperature(temp, temp, "s", "s").Value; got != want {
			t.Fatalf("T=%v: %s, want %s", temp, got, want)
		}
	}
}
