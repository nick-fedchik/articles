package evidence

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"os/exec"
	"testing"
	"time"
)

func TestExternalSATOracle(t *testing.T) {
	solver, err := exec.LookPath("minisat")
	if err != nil {
		t.Skip("external minisat not installed; run cmd/formal with an explicit solver for required validation")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	sources := fixtureSources(t, 3)
	root := &Expr{Op: Select, Condition: &Expr{Op: Input},
		Left:  &Expr{Op: All, Left: &Expr{Op: Input, Input: 1}, Right: &Expr{Op: Input, Input: 1}},
		Right: &Expr{Op: All, Left: &Expr{Op: Input, Input: 1}, Right: &Expr{Op: Input, Input: 2}}}
	check, err := VerifyGraphSAT(ctx, solver, root, sources)
	if err != nil || !check.Equivalent {
		t.Fatalf("miter: %+v %v", check, err)
	}
	check, err = verifyGraphSAT(ctx, solver, root, sources, priorityLUT^1)
	if err != nil || check.Equivalent {
		t.Fatalf("mutation was not detected: %+v %v", check, err)
	}
	rng := rand.New(rand.NewSource(77))
	for n := 0; n < 20; n++ {
		p := RepairProblem{}
		for i := 0; i < 8; i++ {
			p.Actions = append(p.Actions, Action{fmt.Sprintf("a%d", i), int64(1 + rng.Intn(12)), true, false})
		}
		p.Actions[7].Weakens = true
		p.Conflicts = [][2]int{{0, 1}}
		for i := 0; i < 5; i++ {
			p.Cuts = append(p.Cuts, []int{rng.Intn(8), rng.Intn(8), rng.Intn(8)})
		}
		want, possible := bruteRepair(p)
		sat, err := FindRepairSAT(ctx, solver, p)
		if !possible {
			if !errors.Is(err, ErrNoRepair) {
				t.Fatal(err)
			}
		} else if err != nil || sat.Cost != want {
			t.Fatalf("SAT repair %d: %+v %v, want %d", n, sat, err, want)
		}
	}
}

func TestSATUnavailableIsNotAProof(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	c := &cnf{}
	c.add(c.constant(true))
	if _, _, err := c.solve(ctx, "missing-solver"); !errors.Is(err, context.Canceled) {
		t.Fatal("cancellation lost")
	}
	if _, _, err := c.solve(context.Background(), "/nonexistent/teos-solver"); err == nil {
		t.Fatal("missing external solver accepted")
	}
	check, err := VerifyGraphSAT(context.Background(), "/nonexistent/teos-solver", &Expr{Op: Input}, fixtureSources(t, 1))
	if err == nil || check.Equivalent {
		t.Fatal("solver failure produced a success-shaped equivalence claim")
	}
}
