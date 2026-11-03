package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	evidence "github.com/nick-fedchik/articles/Expert_System_Circuits/src"
)

func run(ctx context.Context, solver string) error {
	version, err := evidence.IdentifySolver(ctx, solver)
	if err != nil {
		return fmt.Errorf("required external SAT solver: %w", err)
	}
	rng := rand.New(rand.NewSource(20261102))
	graphs, repairs, impossible := 0, 0, 0
	for n := 0; n < 30; n++ {
		sources := []evidence.Source{}
		for i := 0; i < 6; i++ {
			source, err := evidence.NewSource(fmt.Sprintf("public-%d", i), "1", []byte{byte(i)}, .9)
			if err != nil {
				return err
			}
			sources = append(sources, source)
		}
		var tree func(int) *evidence.Expr
		tree = func(depth int) *evidence.Expr {
			if depth == 0 || rng.Intn(3) == 0 {
				return &evidence.Expr{Op: evidence.Input, Input: rng.Intn(6)}
			}
			out := &evidence.Expr{Op: evidence.All, Left: tree(depth - 1), Right: tree(depth - 1)}
			if rng.Intn(2) == 0 {
				out.Op = evidence.Select
				out.Condition = &evidence.Expr{Op: evidence.Input, Input: rng.Intn(6)}
			}
			return out
		}
		check, err := evidence.VerifyGraphSAT(ctx, solver, tree(4), sources)
		if err != nil || !check.Equivalent {
			return fmt.Errorf("graph %d not proved equivalent: %+v %v", n, check, err)
		}
		graphs++
	}
	for n := 0; n < 60; n++ {
		problem := evidence.RepairProblem{}
		for i := 0; i < 10; i++ {
			problem.Actions = append(problem.Actions, evidence.Action{
				ID: fmt.Sprintf("a%d", i), Cost: int64(1 + rng.Intn(30)), Verified: rng.Intn(10) != 0, Weakens: rng.Intn(12) == 0})
		}
		for i := 0; i < 7; i++ {
			problem.Cuts = append(problem.Cuts, []int{rng.Intn(10), rng.Intn(10), rng.Intn(10)})
		}
		problem.Conflicts = [][2]int{{0, 1}, {3, 5}}
		got, searchErr := evidence.FindRepair(ctx, problem, 1000000)
		oracle, satErr := evidence.FindRepairSAT(ctx, solver, problem)
		if errors.Is(searchErr, evidence.ErrNoRepair) && errors.Is(satErr, evidence.ErrNoRepair) {
			impossible++
		} else if searchErr != nil || satErr != nil || got.Cost != oracle.Cost {
			return fmt.Errorf("repair %d mismatch: search=%+v/%v SAT=%+v/%v", n, got, searchErr, oracle, satErr)
		}
		repairs++
	}
	return json.NewEncoder(os.Stdout).Encode(map[string]any{
		"solver": version, "graph_miters_unsat": graphs, "repair_problems_checked": repairs,
		"impossible_problems": impossible,
		"scope":               "all state assignments of each selected finite graph; no universal graph theorem or floating-point proof",
	})
}

func main() {
	solver := flag.String("solver", "minisat", "external Minisat-compatible solver")
	flag.Parse()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	if err := run(ctx, *solver); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
