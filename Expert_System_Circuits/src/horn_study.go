package evidence

import (
	"errors"
	"fmt"
	"math/rand"
)

var HornFamilies = []string{"distractors", "conjunction", "shared", "alternatives"}

type HornCase struct {
	Program  *HornProgram
	Snapshot HornSnapshot
	Query    int
	Roots    []int
}

func GenerateHornCase(family string, ruleCount int, seed int64) (HornCase, error) {
	if ruleCount < 32 || ruleCount > 9000 {
		return HornCase{}, errors.New("Horn study: rules outside [32,9000]")
	}
	valid := false
	for _, name := range HornFamilies {
		valid = valid || name == family
	}
	if !valid {
		return HornCase{}, errors.New("Horn study: unknown family")
	}
	rng := rand.New(rand.NewSource(seed))
	roots := 32
	if family == "conjunction" {
		roots = ruleCount + 1
	}
	atoms := roots + ruleCount
	rules := make([]HornRule, ruleCount)
	query := atoms - 1
	frontier := make([]int, roots)
	for i := range frontier {
		frontier[i] = i
	}
	for i := range rules {
		head := roots + i
		premises := []int{}
		switch family {
		case "distractors":
			parent := 0
			if i > 0 && i < 8 {
				parent = head - 1
			} else if i >= 8 {
				parent = 1 + i%31
				if i >= 39 {
					parent = head - 31
				}
			}
			premises = []int{parent}
			query = roots + 7
		case "conjunction":
			premises = []int{frontier[0], frontier[1]}
			frontier = append(frontier[2:], head)
		case "shared":
			premises = []int{i % roots}
			if i > 0 {
				premises = append(premises, head-1)
			}
			if i > 2 {
				premises = append(premises, roots+i/2)
			}
		case "alternatives":
			if i == ruleCount-3 {
				premises = []int{2}
			} else if i == ruleCount-2 {
				premises = []int{head - 1}
				query = head
			} else if i == ruleCount-1 {
				head = query
				premises = []int{roots + ruleCount - 4}
			} else if i == 0 {
				premises = []int{0}
			} else {
				premises = []int{i % roots, roots + i - 1}
			}
		}
		rules[i] = HornRule{fmt.Sprintf("rule-%d", i), "1", premises, head}
	}
	permutation := rng.Perm(atoms)
	facts := make([]bool, atoms)
	rootIDs := make([]int, roots)
	for i := 0; i < roots; i++ {
		facts[permutation[i]], rootIDs[i] = true, permutation[i]
	}
	if seed%5 == 4 {
		facts[permutation[0]] = false
		if family == "alternatives" {
			facts[permutation[2]] = false
		}
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
	enabled := make([]bool, ruleCount)
	for i := range enabled {
		enabled[i] = true
	}
	usedRoots := []int{}
	for _, root := range rootIDs {
		if len(p.consumers[root]) > 0 {
			usedRoots = append(usedRoots, root)
		}
	}
	return HornCase{p, HornSnapshot{1, facts, enabled}, permutation[query], usedRoots}, nil
}

func (c HornCase) AtomCount() int { return c.Program.atoms }

func (c HornCase) Expected() (bool, error) {
	if err := c.Program.validate(c.Snapshot, c.Query); err != nil {
		return false, err
	}
	known := append([]bool(nil), c.Snapshot.Facts...)
	changed := true
	for changed {
		changed = false
		for i, rule := range c.Program.rules {
			if !c.Snapshot.Enabled[i] || known[rule.Conclusion] {
				continue
			}
			holds := true
			for _, atom := range rule.Premises {
				holds = holds && known[atom]
			}
			if holds {
				known[rule.Conclusion], changed = true, true
			}
		}
	}
	return known[c.Query], nil
}
