package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	evidence "github.com/nick-fedchik/articles/Epistemic_Circuit_Theory/src"
)

func run(writer io.Writer) error {
	results := map[string]any{}
	checked := 0
	for a := evidence.Unknown; a <= evidence.Both; a++ {
		for b := evidence.Unknown; b <= evidence.Both; b++ {
			for c := evidence.Unknown; c <= evidence.Both; c++ {
				for d := evidence.Unknown; d <= evidence.Both; d++ {
					left, err := evidence.Priority(a, b)
					if err != nil {
						return err
					}
					right, err := evidence.Priority(c, d)
					if err != nil {
						return err
					}
					got, err := evidence.Priority(left, right)
					if err != nil {
						return err
					}
					want, err := evidence.ReferencePriority(a, b, c, d)
					if err != nil || got != want {
						return fmt.Errorf("priority reference mismatch: %v", err)
					}
					checked++
				}
			}
		}
	}
	results["01_semantics"] = map[string]any{"priority_tuples_checked": checked, "kind": "finite model check, not a safety certification"}
	source, err := evidence.NewSource("public-sensor", "1", []byte("public fixture"), .9)
	if err != nil {
		return err
	}
	leaf := &evidence.Expr{Op: evidence.Input}
	root := &evidence.Expr{Op: evidence.All, Left: leaf, Right: leaf}
	program, err := evidence.Compile(root, []evidence.Source{source})
	if err != nil {
		return err
	}
	value, err := program.Evaluate([]evidence.State{evidence.True})
	if err != nil {
		return err
	}
	results["02_provenance"] = map[string]any{"unique_sources": len(value.Sources), "weight": value.Weight, "compiled_nodes": program.Nodes()}
	results["03_runtime"] = "run go test -run '^$' -bench BenchmarkRuleBank -benchmem -count 3; scalar Go, no SIMD claim"
	repair, err := evidence.FindRepair(context.Background(), evidence.RepairProblem{
		Actions: []evidence.Action{{ID: "shared", Cost: 10, Verified: true}, {ID: "left", Cost: 6, Verified: true}, {ID: "right", Cost: 6, Verified: true}},
		Cuts:    [][]int{{0, 1}, {0, 2}},
	}, 1000)
	if err != nil {
		return err
	}
	results["04_repairs"] = repair
	plants := []evidence.PlantResult{}
	for _, width := range []float64{.88, 1.12} {
		plant, err := evidence.SimulateMemoPlant(width, 10000)
		if err != nil {
			return err
		}
		plants = append(plants, plant)
	}
	results["05_zero_ripple_plant"] = plants
	route, err := evidence.SelectCandidates([]evidence.Candidate{
		{ID: "public-standard", Version: "1", Score: .9, Radius: .02},
		{ID: "public-standard", Version: "2", Score: .89, Radius: .02},
	}, .5, .01)
	if err != nil {
		return err
	}
	bound, err := evidence.ZeroFailureUpperBound(100, .05)
	if err != nil {
		return err
	}
	results["06_routing"] = route
	results["zero_errors_in_100_upper_95pct"] = bound
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(results)
}

func main() {
	if err := run(os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
