package main

import (
	"fmt"
	"math/rand"

	evidence "github.com/nick-fedchik/articles/Epistemic_Circuit_Theory/src"
)

func simpleTheory() evidence.NormTheory {
	guard := evidence.SignedLiteral{Atom: 1}
	return evidence.NormTheory{Atoms: 4, Norms: []evidence.NormRule{
		{ID: "permit", Version: "1", Action: "run", Effect: evidence.Permit,
			Conditions: []evidence.SignedLiteral{{Atom: 0}}, Exception: &guard},
		{ID: "ban", Version: "1", Action: "run", Effect: evidence.Prohibit},
	}}
}

func setPriority(theory *evidence.NormTheory, mode int) {
	if mode == 1 {
		theory.Priorities = []evidence.NormPriority{{Higher: "permit", Lower: "ban"}}
	} else if mode == 2 {
		theory.Priorities = []evidence.NormPriority{{Higher: "ban", Lower: "permit"}}
	}
}

func randomTheory(seed int64) evidence.NormTheory {
	rng := rand.New(rand.NewSource(seed))
	theory := evidence.NormTheory{Atoms: 8}
	for i := 0; i < 8; i++ {
		head := 4 + i%4
		premises, seen := []evidence.SignedLiteral{}, map[int]bool{}
		count := 1 + rng.Intn(2)
		for len(premises) < count {
			code := rng.Intn(2 * head)
			if !seen[code] {
				seen[code] = true
				premises = append(premises, evidence.SignedLiteral{Atom: code / 2, Negative: code%2 != 0})
			}
		}
		theory.Evidence = append(theory.Evidence, evidence.SignedRule{
			ID: fmt.Sprintf("e-%d", i), Version: "1", Premises: premises,
			Conclusion: evidence.SignedLiteral{Atom: head, Negative: rng.Intn(2) != 0},
		})
	}
	for i := 0; i < 6; i++ {
		rule := evidence.NormRule{ID: fmt.Sprintf("n-%d", i), Version: "1", Action: fmt.Sprintf("action-%d", i%2), Effect: evidence.Permit}
		if rng.Intn(2) != 0 {
			rule.Effect = evidence.Prohibit
		}
		for n := rng.Intn(4); n > 0; n-- {
			rule.Conditions = append(rule.Conditions, evidence.SignedLiteral{Atom: rng.Intn(8), Negative: rng.Intn(2) != 0})
		}
		if rng.Intn(2) != 0 {
			guard := evidence.SignedLiteral{Atom: rng.Intn(8), Negative: rng.Intn(2) != 0}
			rule.Exception = &guard
		}
		theory.Norms = append(theory.Norms, rule)
	}
	for i := range theory.Norms {
		for j := i + 1; j < len(theory.Norms); j++ {
			if theory.Norms[i].Action == theory.Norms[j].Action && rng.Intn(3) == 0 {
				theory.Priorities = append(theory.Priorities, evidence.NormPriority{Higher: theory.Norms[i].ID, Lower: theory.Norms[j].ID})
			}
		}
	}
	return theory
}

func (s *study) randomStage() error {
	for seed := int64(73000); seed < 73100; seed++ {
		theory := randomTheory(seed)
		p, err := evidence.CompileNormative(theory)
		if err != nil {
			return err
		}
		rng := rand.New(rand.NewSource(seed + 100000))
		for kind := 0; kind < 5; kind++ {
			facts := make([]evidence.State, theory.Atoms)
			for i := range facts {
				if kind < 4 {
					facts[i] = evidence.State(kind)
				} else {
					facts[i] = evidence.State(rng.Intn(4))
				}
			}
			world, err := p.NewWorld(facts)
			if err != nil {
				return err
			}
			for i := range world.EvidenceEnabled {
				world.EvidenceEnabled[i] = rng.Intn(10) != 0
			}
			for i := range world.NormEnabled {
				world.NormEnabled[i] = rng.Intn(8) != 0
			}
			for _, action := range []string{"action-0", "action-1"} {
				if _, err := s.compare(p, world, action, "generated", fmt.Sprintf("seed=%d/world=%d", seed, kind), kind == 4); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
