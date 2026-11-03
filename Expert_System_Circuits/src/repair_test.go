package evidence

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"testing"
)

func bruteRepair(problem RepairProblem) (int64, bool) {
	best, found := int64(math.MaxInt64), false
	for mask := uint64(0); mask < uint64(1)<<uint(len(problem.Actions)); mask++ {
		cost, valid := int64(0), true
		for i, action := range problem.Actions {
			if mask&(uint64(1)<<uint(i)) != 0 {
				valid = valid && action.Verified && !action.Weakens
				cost += action.Cost
			}
		}
		for _, pair := range problem.Conflicts {
			if mask&(uint64(1)<<uint(pair[0])) != 0 && mask&(uint64(1)<<uint(pair[1])) != 0 {
				valid = false
			}
		}
		for _, cut := range problem.Cuts {
			satisfied := false
			for _, index := range cut {
				satisfied = satisfied || mask&(uint64(1)<<uint(index)) != 0
			}
			valid = valid && satisfied
		}
		if valid && cost < best {
			best, found = cost, true
		}
	}
	return best, found
}

func TestSharedRepairAndAdmissibility(t *testing.T) {
	problem := RepairProblem{
		Actions: []Action{
			{"shared", 10, true, false},
			{"left", 6, true, false},
			{"right", 6, true, false},
			{"weaken", 1, true, true},
		},
		Cuts: [][]int{{0, 1, 3}, {0, 2, 3}},
	}
	got, err := FindRepair(context.Background(), problem, 1000)
	if err != nil || got.Cost != 10 || got.Mask != 1 {
		t.Fatalf("got %+v, %v", got, err)
	}
}

func TestRepairsMatchBruteForceAcrossTopologies(t *testing.T) {
	rng := rand.New(rand.NewSource(2026))
	for n := 0; n < 300; n++ {
		count := 1 + rng.Intn(10)
		problem := RepairProblem{}
		for i := 0; i < count; i++ {
			problem.Actions = append(problem.Actions, Action{fmt.Sprintf("a%d", i), int64(1 + rng.Intn(30)), rng.Intn(8) != 0, rng.Intn(8) == 0})
		}
		for c := 0; c < 1+rng.Intn(7); c++ {
			cut := []int{}
			for i := 0; i < count; i++ {
				if rng.Intn(3) == 0 {
					cut = append(cut, i)
				}
			}
			if len(cut) == 0 {
				cut = append(cut, rng.Intn(count))
			}
			problem.Cuts = append(problem.Cuts, cut)
		}
		if count > 1 {
			problem.Conflicts = append(problem.Conflicts, [2]int{0, count - 1})
		}
		want, possible := bruteRepair(problem)
		got, err := FindRepair(context.Background(), problem, 100000)
		if !possible {
			if !errors.Is(err, ErrNoRepair) {
				t.Fatalf("case %d: want no repair, got %+v %v", n, got, err)
			}
		} else if err != nil || got.Cost != want {
			t.Fatalf("case %d: got %+v %v, want %d", n, got, err, want)
		}
	}
}

func TestRepairBudgetCancellationAndJointConflicts(t *testing.T) {
	problem := RepairProblem{Actions: []Action{{"a", 1, true, false}, {"b", 1, true, false}}, Cuts: [][]int{{0}, {1}}}
	if got, err := FindRepair(context.Background(), problem, 1); !errors.Is(err, ErrSearchBudget) || got != (Repair{}) {
		t.Fatalf("budget result must be explicit failure: %+v %v", got, err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := FindRepair(ctx, problem, 100); !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
	problem.Conflicts = [][2]int{{0, 1}}
	if _, err := FindRepair(context.Background(), problem, 100); !errors.Is(err, ErrNoRepair) {
		t.Fatalf("incompatible joint repair accepted: %v", err)
	}
	problem.Actions[0].Cost = math.MaxInt64
	if _, err := FindRepair(context.Background(), problem, 100); err == nil {
		t.Fatal("overflowing cost accepted")
	}
}

func TestExponentialFrontierDoesNotImplyPolynomialEnumeration(t *testing.T) {
	const pairs = 20
	const minimalRepairs = uint64(1) << pairs
	if minimalRepairs != 1048576 {
		t.Fatal("incorrect cardinality")
	}
	problem := RepairProblem{}
	for i := 0; i < pairs; i++ {
		problem.Actions = append(problem.Actions,
			Action{fmt.Sprintf("a%d", i), 1, true, false},
			Action{fmt.Sprintf("b%d", i), 1, true, false})
		problem.Cuts = append(problem.Cuts, []int{2 * i, 2*i + 1})
	}
	if _, err := FindRepair(context.Background(), problem, 100); !errors.Is(err, ErrSearchBudget) {
		t.Fatalf("expected explicit bounded search, got %v", err)
	}
}
