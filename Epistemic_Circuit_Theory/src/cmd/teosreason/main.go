// Command teosreason prints the full inference trace and answer of one
// scenario from research protocol 01.
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/nick-fedchik/articles/Epistemic_Circuit_Theory/src/teosreason"
)

func main() {
	name := flag.String("scenario", "S1", "scenario: "+strings.Join(teosreason.ScenarioNames, ", "))
	flag.Parse()
	known := false
	for _, n := range teosreason.ScenarioNames {
		known = known || n == *name
	}
	if !known {
		fmt.Fprintln(os.Stderr, "unknown scenario", *name)
		os.Exit(2)
	}
	q, m, obs := teosreason.Scenario(*name)
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
	fmt.Printf("передумова %s=%s за %v\n", a.Target.Key(), a.Presupposition, a.PresuppositionRoots)
	fmt.Printf("сумісні %v\nвиключені %v\nоскаржені %v\nнапрям %q\n", a.Compatible, a.Excluded, a.Contested, a.Direction)
	for _, c := range a.NextChecks {
		fmt.Printf("наступна перевірка %s (для %s, %s)\n", c.Atom.Key(), c.Hypothesis, c.Rule)
	}
	fmt.Printf("припущення %v\nневикористані спостереження %v\n", a.Assumptions, a.UnusedObservations)
}
