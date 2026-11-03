package evidence

import (
	"errors"
	"fmt"
	"math/rand"
)

var ProvenanceFamilies = []string{"local-select", "dense-overlap", "repeated-support", "balanced-all"}

type ProvenanceCase struct {
	Root    *Expr
	Sources []Source
	Inputs  []State
	Program *Program
}

func balancedProvenance(nodes []*Expr) *Expr {
	if len(nodes) == 1 {
		return nodes[0]
	}
	middle := len(nodes) / 2
	return &Expr{Op: All, Left: balancedProvenance(nodes[:middle]), Right: balancedProvenance(nodes[middle:])}
}

func GenerateProvenanceCase(family string, sourceCount int, regime string, seed int64) (ProvenanceCase, error) {
	valid := false
	for _, name := range ProvenanceFamilies {
		valid = valid || name == family
	}
	if !valid || sourceCount < 32 || sourceCount > 4096 || sourceCount&(sourceCount-1) != 0 ||
		regime != "resolved" && regime != "uncertain" {
		return ProvenanceCase{}, errors.New("provenance study: unknown family/regime or invalid power-of-two source count")
	}
	rng := rand.New(rand.NewSource(seed))
	permutation := rng.Perm(sourceCount)
	c := ProvenanceCase{Sources: make([]Source, sourceCount), Inputs: make([]State, sourceCount)}
	atoms := make([]*Expr, sourceCount)
	for i := range c.Sources {
		source, err := NewSource(fmt.Sprintf("public-provenance-%d", i), "1", []byte(fmt.Sprintf("public generated source %d", i)), .9)
		if err != nil {
			return ProvenanceCase{}, err
		}
		c.Sources[i] = source
		c.Inputs[i] = True
		if regime == "uncertain" {
			c.Inputs[i] = Both
		}
		atoms[i] = &Expr{Op: Input, Input: permutation[i]}
	}
	if family == "balanced-all" {
		c.Root = balancedProvenance(atoms)
	} else {
		groups := min(128, sourceCount/4)
		blockSize := sourceCount / groups
		blocks := make([]*Expr, groups)
		for i := range blocks {
			blocks[i] = balancedProvenance(atoms[i*blockSize : (i+1)*blockSize])
		}
		var selection func([]*Expr, int) *Expr
		selection = func(nodes []*Expr, begin int) *Expr {
			if len(nodes) == 1 {
				return nodes[0]
			}
			middle := len(nodes) / 2
			condition := atoms[(17*begin+31*(begin+len(nodes)))%len(atoms)]
			return &Expr{Op: Select, Condition: condition,
				Left: selection(nodes[:middle], begin), Right: selection(nodes[middle:], begin+middle)}
		}
		if family == "local-select" {
			c.Root = selection(blocks, 0)
		} else {
			prefix, suffix := make([]*Expr, groups), make([]*Expr, groups)
			prefix[0], suffix[groups-1] = blocks[0], blocks[groups-1]
			for i := 1; i < groups; i++ {
				prefix[i] = &Expr{Op: All, Left: prefix[i-1], Right: blocks[i]}
			}
			for i := groups - 2; i >= 0; i-- {
				suffix[i] = &Expr{Op: All, Left: blocks[i], Right: suffix[i+1]}
			}
			choices := []*Expr{}
			if family == "repeated-support" {
				for i := 0; i < groups-1; i++ {
					choices = append(choices, &Expr{Op: All, Left: prefix[i], Right: suffix[i+1]})
				}
			} else {
				choices = append(choices, suffix[1])
				for i := 1; i < groups-1; i++ {
					choices = append(choices, &Expr{Op: All, Left: prefix[i-1], Right: suffix[i+1]})
				}
				choices = append(choices, prefix[groups-2])
			}
			c.Root = selection(choices, 0)
		}
	}
	var err error
	c.Program, err = Compile(c.Root, c.Sources)
	if err != nil {
		return ProvenanceCase{}, err
	}
	used := make([]bool, sourceCount)
	for _, inst := range c.Program.code {
		if inst.op == Input {
			used[inst.input] = true
		}
	}
	for _, present := range used {
		if !present {
			return ProvenanceCase{}, errors.New("provenance study: an unused source inflated the declared universe")
		}
	}
	return c, nil
}
