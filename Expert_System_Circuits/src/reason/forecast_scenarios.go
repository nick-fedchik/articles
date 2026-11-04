package reason

import (
	evidence "github.com/nick-fedchik/articles/Expert_System_Circuits/src"
)

// ForecastModel returns the supply-chain forecast model of protocol 02 (§3.2–§3.3)
// with rules R1–R7s, assumptions A1–A3 active by default, and A4 inactive.
func ForecastModel() Model {
	return Model{
		Premise: A("batch_delayed", "?p"),
		Derived: []Atom{
			A("buffer_exhausted", "?p"),
			A("rework_window_lost", "?p"),
			A("assembly_halted", "?p"),
			A("qa_hold", "?p"),
			A("schedule_delayed", "?p"),
		},
		Rules: []Rule{
			{
				ID: "R1",
				Premises: []Premise{
					{Atom: A("batch_delayed", "?p"), Expected: evidence.True},
					{Atom: A("reserve_depleted", "?p"), Expected: evidence.True},
				},
				Then:     A("buffer_exhausted", "?p"),
				Expected: evidence.True,
				Strength: Strong,
				Requires: []string{"A1"},
			},
			{
				ID: "R2",
				Premises: []Premise{
					{Atom: A("buffer_exhausted", "?p"), Expected: evidence.True},
					{Atom: A("expedite_available", "?p"), Expected: evidence.False},
				},
				Then:     A("assembly_halted", "?p"),
				Expected: evidence.True,
				Strength: Strong,
				Requires: []string{"A1", "A2"},
			},
			{
				ID: "R3",
				Premises: []Premise{
					{Atom: A("batch_delayed", "?p"), Expected: evidence.True},
					{Atom: A("reserve_depleted", "?p"), Expected: evidence.True},
				},
				Then:     A("rework_window_lost", "?p"),
				Expected: evidence.True,
				Strength: Strong,
				Requires: []string{"A3"},
			},
			{
				ID: "R4",
				Premises: []Premise{
					{Atom: A("rework_window_lost", "?p"), Expected: evidence.True},
					{Atom: A("qa_staff_shortage", "?p"), Expected: evidence.True},
				},
				Then:     A("qa_hold", "?p"),
				Expected: evidence.True,
				Strength: Strong,
				Requires: []string{"A2", "A3"},
			},
			{
				ID: "R5",
				Premises: []Premise{
					{Atom: A("assembly_halted", "?p"), Expected: evidence.True},
				},
				Then:     A("schedule_delayed", "?p"),
				Expected: evidence.True,
				Strength: Strong,
				Requires: []string{"A1"},
			},
			{
				ID: "R6",
				Premises: []Premise{
					{Atom: A("qa_hold", "?p"), Expected: evidence.True},
				},
				Then:     A("schedule_delayed", "?p"),
				Expected: evidence.True,
				Strength: Strong,
				Requires: []string{"A3"},
			},
			{
				ID: "R7",
				Premises: []Premise{
					{Atom: A("batch_delayed", "?p"), Expected: evidence.True},
					{Atom: A("reserve_depleted", "?p"), Expected: evidence.False},
				},
				Then:     A("schedule_delayed", "?p"),
				Expected: evidence.False,
				Strength: Typical,
				Requires: []string{"A1"},
			},
			{
				ID: "R7s",
				Premises: []Premise{
					{Atom: A("batch_delayed", "?p"), Expected: evidence.True},
					{Atom: A("reserve_depleted", "?p"), Expected: evidence.False},
				},
				Then:     A("schedule_delayed", "?p"),
				Expected: evidence.False,
				Strength: Strong,
				Requires: []string{"A1", "A4"},
			},
		},
		Active: map[string]bool{"A1": true, "A2": true, "A3": true},
	}
}

func ForecastSchedule() Query       { return Query{Kind: Forecast, Target: A("schedule_delayed", "p")} }
func WhetherBufferExhausted() Query { return Query{Kind: Whether, Target: A("buffer_exhausted", "p")} }

func SupplyFact(pred string, v evidence.State, root string) Observation {
	return Observation{Atom: A(pred, "p"), Value: v, Relay: root, Root: root}
}

// ForecastScenario returns the query, model and observations of protocol 02 (§3.5).
func ForecastScenario(name string) (Query, Model, []Observation) {
	m := ForecastModel()
	s1 := []Observation{
		SupplyFact("batch_delayed", evidence.True, "sup1"),
		SupplyFact("reserve_depleted", evidence.True, "erp1"),
		SupplyFact("expedite_available", evidence.False, "pln3"),
		SupplyFact("qa_staff_shortage", evidence.True, "qa4"),
	}
	switch name {
	case "S1":
		return ForecastSchedule(), m, s1
	case "S2":
		return ForecastSchedule(), m, s1[:2]
	case "S3":
		return ForecastSchedule(), m, append(clone(s1), SupplyFact("reserve_depleted", evidence.False, "aud2"))
	case "S4a":
		cp := SupplyFact("reserve_depleted", evidence.True, "erp1")
		cp.Relay = "portal"
		return ForecastSchedule(), m, append(clone(s1), cp)
	case "S4b":
		return ForecastSchedule(), m, []Observation{
			s1[0], s1[1], s1[2],
			SupplyFact("qa_staff_shortage", evidence.True, "pln3"),
		}
	case "S5a":
		return WhetherBufferExhausted(), m, s1
	case "S5b":
		return ForecastSchedule(), m, []Observation{
			s1[0], s1[1],
			SupplyFact("expedite_available", evidence.True, "pln3"),
			SupplyFact("qa_staff_shortage", evidence.False, "qa4"),
		}
	case "S7a-1":
		var rules []Rule
		for _, r := range m.Rules {
			if r.ID != "R2" {
				rules = append(rules, r)
			}
		}
		m.Rules = rules
		return ForecastSchedule(), m, s1
	case "S7a-2":
		var rules []Rule
		for _, r := range m.Rules {
			if r.ID != "R2" && r.ID != "R4" {
				rules = append(rules, r)
			}
		}
		m.Rules = rules
		return ForecastSchedule(), m, s1
	case "S7b-0":
		return ForecastSchedule(), m, []Observation{
			s1[0],
			SupplyFact("reserve_depleted", evidence.False, "erp1"),
			s1[2], s1[3],
		}
	case "S7b-1":
		m.Active["A4"] = true
		return ForecastSchedule(), m, []Observation{
			s1[0],
			SupplyFact("reserve_depleted", evidence.False, "erp1"),
			s1[2], s1[3],
		}
	case "S7c":
		van := SupplyFact("courier_color_blue", evidence.True, "car1")
		van.Reliability = 1.0
		return ForecastSchedule(), m, append(clone(s1), van)
	case "S8":
		return ForecastSchedule(), m, append(clone(s1), SupplyFact("batch_delayed", evidence.False, "car1"))
	}
	panic("unknown forecast scenario " + name)
}

var ForecastScenarioNames = []string{
	"S1", "S2", "S3", "S4a", "S4b", "S5a", "S5b", "S7a-1", "S7a-2", "S7b-0", "S7b-1", "S7c", "S8",
}
