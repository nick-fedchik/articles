package evidence

import (
	"context"
	"errors"
	"fmt"
	"math"
)

var (
	ErrSearchBudget = errors.New("repair: search budget exhausted; optimum not proved")
	ErrNoRepair     = errors.New("repair: no admissible repair")
)

type Action struct {
	ID                string
	Cost              int64
	Verified, Weakens bool
}

type RepairProblem struct {
	Actions   []Action
	Cuts      [][]int
	Conflicts [][2]int
}

type Repair struct {
	Mask    uint64
	Cost    int64
	Visited uint64
}

func (p RepairProblem) validate() error {
	if len(p.Actions) == 0 || len(p.Actions) > 63 {
		return errors.New("repair: action count outside [1,63]")
	}
	seen := map[string]bool{}
	var total int64
	for _, action := range p.Actions {
		if action.ID == "" || seen[action.ID] || action.Cost <= 0 || total > math.MaxInt64-action.Cost {
			return errors.New("repair: duplicate/empty ID, nonpositive cost or cost overflow")
		}
		seen[action.ID] = true
		total += action.Cost
	}
	for _, cut := range p.Cuts {
		if len(cut) == 0 {
			return errors.New("repair: empty cut")
		}
		for _, action := range cut {
			if action < 0 || action >= len(p.Actions) {
				return errors.New("repair: action index outside range")
			}
		}
	}
	for _, pair := range p.Conflicts {
		if pair[0] < 0 || pair[1] < 0 || pair[0] >= len(p.Actions) || pair[1] >= len(p.Actions) {
			return errors.New("repair: conflict index outside range")
		}
	}
	return nil
}

func FindRepair(ctx context.Context, problem RepairProblem, maxVisited uint64) (Repair, error) {
	if ctx == nil || maxVisited == 0 {
		return Repair{}, errors.New("repair: context and positive search budget required")
	}
	if err := problem.validate(); err != nil {
		return Repair{}, err
	}
	cuts := make([]uint64, len(problem.Cuts))
	for i, cut := range problem.Cuts {
		for _, k := range cut {
			action := problem.Actions[k]
			if action.Verified && !action.Weakens {
				cuts[i] |= uint64(1) << uint(k)
			}
		}
		if cuts[i] == 0 {
			return Repair{}, ErrNoRepair
		}
	}
	best := Repair{}
	found := false
	var visited uint64
	var search func(uint64, int64) error
	search = func(mask uint64, cost int64) error {
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("repair: %w", err)
		}
		if visited >= maxVisited {
			return ErrSearchBudget
		}
		visited++
		if found && cost >= best.Cost {
			return nil
		}
		for _, pair := range problem.Conflicts {
			if mask&(uint64(1)<<uint(pair[0])) != 0 && mask&(uint64(1)<<uint(pair[1])) != 0 {
				return nil
			}
		}
		for _, cut := range cuts {
			if mask&cut != 0 {
				continue
			}
			for i, action := range problem.Actions {
				bit := uint64(1) << uint(i)
				if cut&bit != 0 {
					if err := search(mask|bit, cost+action.Cost); err != nil {
						return err
					}
				}
			}
			return nil
		}
		best.Mask, best.Cost = mask, cost
		found = true
		return nil
	}
	if err := search(0, 0); err != nil {
		return Repair{}, err
	}
	if !found {
		return Repair{}, ErrNoRepair
	}
	best.Visited = visited
	return best, nil
}
