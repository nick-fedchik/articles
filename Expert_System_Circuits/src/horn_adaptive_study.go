package evidence

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"sort"
)

var AdaptiveFamilies = []string{"clusters", "fanout", "layered", "mixed"}

func GenerateAdaptiveCase(family string, ruleCount int, seed int64) (HornCase, error) {
	if ruleCount < 128 || ruleCount > 9000 {
		return HornCase{}, errors.New("adaptive study: rule count outside [128,9000]")
	}
	valid := false
	for _, name := range AdaptiveFamilies {
		valid = valid || family == name
	}
	if !valid {
		return HornCase{}, errors.New("adaptive study: unknown family")
	}
	rng := rand.New(rand.NewSource(seed))
	roots := max(32, ruleCount/8)
	atoms := roots + ruleCount
	rules := make([]HornRule, ruleCount)
	for i := range rules {
		premises := []int{i % roots}
		switch family {
		case "clusters":
			group := i % 16
			width := roots / 16
			premises[0] = group*width + (i/16)%width
			if i >= 16 && rng.Intn(4) != 0 {
				premises = append(premises, roots+i-16)
			}
		case "fanout":
			if i > 0 {
				premises = append(premises, roots)
			}
			if i > 1 && rng.Intn(3) == 0 {
				premises = append(premises, roots+rng.Intn(i))
			}
		case "layered":
			if i >= 32 {
				begin := (i/32 - 1) * 32
				premises = append(premises, roots+begin+rng.Intn(32), roots+begin+rng.Intn(32))
			}
		case "mixed":
			if i > 0 {
				for n := 0; n < 1+rng.Intn(3); n++ {
					if rng.Intn(4) != 0 {
						premises = append(premises, roots+rng.Intn(i))
					}
				}
			}
		}
		sort.Ints(premises)
		unique := premises[:0]
		for _, atom := range premises {
			if len(unique) == 0 || unique[len(unique)-1] != atom {
				unique = append(unique, atom)
			}
		}
		rules[i] = HornRule{fmt.Sprintf("adaptive-%d", i), "1", unique, roots + i}
	}
	permutation := rng.Perm(atoms)
	facts, rootIDs := make([]bool, atoms), []int{}
	for i := 0; i < roots; i++ {
		facts[permutation[i]] = true
		rootIDs = append(rootIDs, permutation[i])
	}
	for i := range rules {
		rules[i].Conclusion = permutation[rules[i].Conclusion]
		for j := range rules[i].Premises {
			rules[i].Premises[j] = permutation[rules[i].Premises[j]]
		}
	}
	rng.Shuffle(len(rules), func(i, j int) { rules[i], rules[j] = rules[j], rules[i] })
	p, err := CompileHorn(atoms, rules)
	if err != nil {
		return HornCase{}, err
	}
	used := []int{}
	for _, atom := range rootIDs {
		if len(p.consumers[atom]) > 0 {
			used = append(used, atom)
		}
	}
	enabled := make([]bool, ruleCount)
	for i := range enabled {
		enabled[i] = true
	}
	return HornCase{p, HornSnapshot{1, facts, enabled}, permutation[atoms-1], used}, nil
}

type PlannedHornUpdate struct {
	Percent int
	Kind    string
	Update  HornUpdate
}

func PlanAdaptiveUpdates(c HornCase, seed int64) ([]PlannedHornUpdate, error) {
	if err := c.Program.validate(c.Snapshot, c.Query); err != nil {
		return nil, err
	}
	if len(c.Roots) == 0 {
		return nil, errors.New("adaptive study: no used primary facts")
	}
	snapshot := cloneHornSnapshot(c.Snapshot)
	rng := rand.New(rand.NewSource(seed))
	out := []PlannedHornUpdate{}
	for round := 0; round < 5; round++ {
		for _, percent := range []int{1, 5, 10, 50, 100} {
			kind, limit := "facts", len(c.Roots)
			if round == 2 || round == 4 {
				kind, limit = "rules", len(snapshot.Enabled)
			}
			order := rng.Perm(limit)
			count := int(math.Ceil(float64(limit*percent) / 100))
			if round == 1 {
				hub := 0
				for i, atom := range c.Roots {
					if len(c.Program.consumers[atom]) > len(c.Program.consumers[c.Roots[hub]]) {
						hub = i
					}
				}
				for i, index := range order {
					if index == hub {
						order[0], order[i] = order[i], order[0]
						break
					}
				}
				kind = "hub-facts"
			}
			update := HornUpdate{}
			for _, index := range order[:count] {
				if kind == "rules" {
					value := !snapshot.Enabled[index]
					update.Rules = append(update.Rules, HornChange{index, value})
					snapshot.Enabled[index] = value
				} else {
					atom := c.Roots[index]
					value := !snapshot.Facts[atom]
					update.Facts = append(update.Facts, HornChange{atom, value})
					snapshot.Facts[atom] = value
				}
			}
			out = append(out, PlannedHornUpdate{percent, kind, update})
		}
	}
	return out, nil
}

func (c HornCase) DerivedQueries(count int, seed int64) ([]int, error) {
	if count < 1 || c.Program == nil {
		return nil, errors.New("Horn study: valid program and positive query count required")
	}
	candidates := []int{}
	for atom, rules := range c.Program.producers {
		if len(rules) > 0 {
			candidates = append(candidates, atom)
		}
	}
	if count > len(candidates) {
		return nil, errors.New("Horn study: insufficient distinct derived atoms")
	}
	rng := rand.New(rand.NewSource(seed))
	rng.Shuffle(len(candidates), func(i, j int) { candidates[i], candidates[j] = candidates[j], candidates[i] })
	return candidates[:count], nil
}
