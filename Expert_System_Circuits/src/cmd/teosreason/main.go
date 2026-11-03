// Command teosreason prints the full inference trace and answer of one
// scenario from research protocol 01.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/nick-fedchik/articles/Expert_System_Circuits/src/teosreason"
)

func main() {
	caseNum := flag.Int("case", 1, "case number: 1 (overheating diagnosis) or 2 (supply forecast)")
	name := flag.String("scenario", "S1", "scenario name")
	flag.Parse()

	var (
		q     teosreason.Query
		m     teosreason.Model
		obs   []teosreason.Observation
		names = teosreason.ScenarioNames
	)
	if *caseNum == 2 {
		names = teosreason.ForecastScenarioNames
	}
	known := false
	for _, n := range names {
		known = known || n == *name
	}
	if !known {
		fmt.Fprintln(os.Stderr, "unknown scenario", *name)
		os.Exit(2)
	}
	if *caseNum == 2 {
		q, m, obs = teosreason.ForecastScenario(*name)
	} else {
		q, m, obs = teosreason.Scenario(*name)
	}
	a, trace, err := teosreason.Reason(q, m, obs)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	for _, s := range trace {
		fmt.Printf("%2d %-14s [%s] %s\n", s.N, s.Operator, s.Kind, s.Basis)
		for _, x := range s.Inputs {
			fmt.Printf("     <- %s\n", x)
		}
		for _, x := range s.Outputs {
			fmt.Printf("     -> %s\n", x)
		}
	}
	a = teosreason.Canonical(a)
	fmt.Printf("\nВІДПОВІДЬ (%s)\n", a.Label)
	if a.Kind == teosreason.Forecast {
		fmt.Printf("вхідна умова %s=%s за %v\n", a.Premise.Key(), a.Presupposition, a.PresuppositionRoots)
		fmt.Printf("цільовий висновок %s=%s за %v\n", a.Target.Key(), a.TargetStatus, a.TargetRoots)
		fmt.Printf("активні правила %v\nзаблоковані правила %v\nоскаржені правила %v\n", a.ActiveRules, a.BlockedRules, a.ContestedRules)
		fmt.Printf("відсутні вхідні дані %v\nприпущення %v\nневикористані спостереження %v\n", a.MissingFacts, a.Assumptions, a.UnusedObservations)
		return
	}
	fmt.Printf("передумова %s=%s за %v\n", a.Target.Key(), a.Presupposition, a.PresuppositionRoots)
	fmt.Printf("сумісні %v\nвиключені %v\nоскаржені %v\nнапрям %q\n", a.Compatible, a.Excluded, a.Contested, a.Direction)
	for _, c := range a.NextChecks {
		fmt.Printf("наступна перевірка %s (для %s, %s)\n", c.Atom.Key(), c.Hypothesis, c.Rule)
	}
	fmt.Printf("припущення %v\nневикористані спостереження %v\n", a.Assumptions, a.UnusedObservations)
}
