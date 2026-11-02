package evidence

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"math"
	"sort"
)

type HornRule struct {
	ID, Version string
	Premises    []int
	Conclusion  int
}

type HornProgram struct {
	atoms     int
	rules     []HornRule
	producers [][]int
	consumers [][]int
	topology  []int
	digest    [32]byte
	edges     int
}

type HornSnapshot struct {
	Revision uint64
	Facts    []bool
	Enabled  []bool
}

type HornProof struct {
	Theory   [32]byte
	Revision uint64
	Query    int
	Facts    []int
	Rules    []int
}

func CompileHorn(atoms int, rules []HornRule) (*HornProgram, error) {
	if atoms < 1 || atoms > 20000 || len(rules) > 10000 {
		return nil, errors.New("Horn: atom/rule size limit")
	}
	p := &HornProgram{atoms: atoms, producers: make([][]int, atoms), consumers: make([][]int, atoms)}
	seen := map[string]bool{}
	adj := make([][]int, atoms)
	degree := make([]int, atoms)
	hash := sha256.New()
	fmt.Fprintf(hash, "atoms:%d;", atoms)
	for i, rule := range rules {
		key := rule.ID + "\x00" + rule.Version
		if rule.ID == "" || rule.Version == "" || seen[key] || rule.Conclusion < 0 ||
			rule.Conclusion >= atoms || len(rule.Premises) == 0 || len(rule.Premises) > 64 {
			return nil, errors.New("Horn: invalid rule identity or arity")
		}
		seen[key] = true
		rule.Premises = append([]int(nil), rule.Premises...)
		sort.Ints(rule.Premises)
		for j, atom := range rule.Premises {
			if atom < 0 || atom >= atoms || j > 0 && atom == rule.Premises[j-1] {
				return nil, errors.New("Horn: invalid or repeated premise")
			}
			p.consumers[atom] = append(p.consumers[atom], i)
			adj[atom] = append(adj[atom], rule.Conclusion)
			degree[rule.Conclusion]++
			p.edges++
		}
		p.rules = append(p.rules, rule)
		p.producers[rule.Conclusion] = append(p.producers[rule.Conclusion], i)
		fmt.Fprintf(hash, "%q:%q:%v:%d;", rule.ID, rule.Version, rule.Premises, rule.Conclusion)
	}
	copy(p.digest[:], hash.Sum(nil))
	queue := []int{}
	for atom, d := range degree {
		if d == 0 {
			queue = append(queue, atom)
		}
	}
	for len(queue) > 0 {
		atom := queue[0]
		queue = queue[1:]
		p.topology = append(p.topology, atom)
		for _, next := range adj[atom] {
			degree[next]--
			if degree[next] == 0 {
				queue = append(queue, next)
			}
		}
	}
	if len(p.topology) != atoms {
		return nil, errors.New("Horn: cyclic theory outside the declared fragment")
	}
	return p, nil
}

func (p *HornProgram) validate(snapshot HornSnapshot, query int) error {
	if p == nil || p.atoms == 0 || len(snapshot.Facts) != p.atoms || len(snapshot.Enabled) != len(p.rules) ||
		snapshot.Revision == 0 || query < 0 || query >= p.atoms {
		return errors.New("Horn: invalid program, snapshot or query")
	}
	return nil
}

func (p *HornProgram) BaseVolume() int { return p.atoms + len(p.rules) + p.edges }

func (p *HornProgram) Rule(index int) (HornRule, error) {
	if p == nil || index < 0 || index >= len(p.rules) {
		return HornRule{}, errors.New("Horn: invalid rule index")
	}
	rule := p.rules[index]
	rule.Premises = append([]int(nil), rule.Premises...)
	return rule, nil
}

func (p *HornProgram) FullClosure(snapshot HornSnapshot) ([]bool, error) {
	if err := p.validate(snapshot, 0); err != nil {
		return nil, err
	}
	known, _ := p.fullClosure(snapshot)
	return known, nil
}

func (p *HornProgram) fullClosure(snapshot HornSnapshot) ([]bool, int) {
	known := append([]bool(nil), snapshot.Facts...)
	inspected := 0
	for _, atom := range p.topology {
		if known[atom] {
			continue
		}
		for _, index := range p.producers[atom] {
			inspected++
			if !snapshot.Enabled[index] {
				continue
			}
			rule := p.rules[index]
			holds := true
			for _, premise := range rule.Premises {
				if !known[premise] {
					holds = false
					break
				}
			}
			if holds {
				known[atom] = true
				break
			}
		}
	}
	return known, inspected
}

func (p *HornProgram) VerifyProof(snapshot HornSnapshot, query int, proof HornProof) error {
	if err := p.validate(snapshot, query); err != nil {
		return err
	}
	if proof.Query != query || proof.Revision != snapshot.Revision || proof.Theory != p.digest ||
		len(proof.Facts) > p.atoms || len(proof.Rules) > len(p.rules) {
		return errors.New("Horn proof: wrong query, theory, revision or size")
	}
	known := make([]bool, p.atoms)
	for _, atom := range proof.Facts {
		if atom < 0 || atom >= p.atoms || !snapshot.Facts[atom] || known[atom] {
			return errors.New("Horn proof: invalid fact")
		}
		known[atom] = true
	}
	used := make([]bool, len(p.rules))
	for _, index := range proof.Rules {
		if index < 0 || index >= len(p.rules) || !snapshot.Enabled[index] || used[index] {
			return errors.New("Horn proof: invalid rule")
		}
		used[index] = true
		rule := p.rules[index]
		for _, premise := range rule.Premises {
			if !known[premise] {
				return errors.New("Horn proof: missing conjunctive premise")
			}
		}
		known[rule.Conclusion] = true
	}
	if !known[query] {
		return errors.New("Horn proof: query not established")
	}
	return nil
}

func (p *HornProgram) certificate(snapshot HornSnapshot, query int, parents []int, charge func(int) bool) (HornProof, error) {
	if !charge(p.atoms) {
		return HornProof{}, ErrHornBudget
	}
	proof := HornProof{Theory: p.digest, Revision: snapshot.Revision, Query: query}
	seen := make([]bool, p.atoms)
	var visit func(int) error
	visit = func(atom int) error {
		if !charge(1) {
			return ErrHornBudget
		}
		if seen[atom] {
			return nil
		}
		seen[atom] = true
		if snapshot.Facts[atom] {
			proof.Facts = append(proof.Facts, atom)
			return nil
		}
		index := parents[atom]
		if index < 0 || !snapshot.Enabled[index] {
			return errors.New("Horn certificate: missing support")
		}
		for _, premise := range p.rules[index].Premises {
			if err := visit(premise); err != nil {
				return err
			}
		}
		proof.Rules = append(proof.Rules, index)
		return nil
	}
	if err := visit(query); err != nil {
		return HornProof{}, err
	}
	return proof, nil
}

type HornChange struct {
	Index   int
	Present bool
}

type HornUpdate struct {
	Facts, Rules []HornChange
}

type HornState struct {
	program  *HornProgram
	snapshot HornSnapshot
	known    []bool
}

type HornUpdateStats struct {
	AffectedAtoms, InspectedRules    int
	DiscoveryVisits, RecomputedAtoms int
	Mode                             string
}

func NewHornState(p *HornProgram, snapshot HornSnapshot) (*HornState, error) {
	known, err := p.FullClosure(snapshot)
	if err != nil {
		return nil, err
	}
	return &HornState{p, cloneHornSnapshot(snapshot), known}, nil
}

func cloneHornSnapshot(s HornSnapshot) HornSnapshot {
	return HornSnapshot{s.Revision, append([]bool(nil), s.Facts...), append([]bool(nil), s.Enabled...)}
}

func (s *HornState) Snapshot() HornSnapshot { return cloneHornSnapshot(s.snapshot) }
func (s *HornState) Known() []bool          { return append([]bool(nil), s.known...) }

func (s *HornState) Update(update HornUpdate) (HornUpdateStats, error) {
	return s.update(update, "incremental", 0)
}

func (s *HornState) UpdateFull(update HornUpdate) (HornUpdateStats, error) {
	return s.update(update, "full", 0)
}

func (s *HornState) UpdateAdaptive(update HornUpdate, discoveryFraction float64) (HornUpdateStats, error) {
	if !finite(discoveryFraction) || discoveryFraction <= 0 || discoveryFraction > 1 {
		return HornUpdateStats{}, errors.New("Horn adaptive: discovery fraction must be in (0,1]")
	}
	return s.update(update, "adaptive", discoveryFraction)
}

func (s *HornState) validateUpdate(update HornUpdate) error {
	if s == nil || s.program == nil || s.snapshot.Revision == math.MaxUint64 ||
		len(update.Facts)+len(update.Rules) == 0 {
		return errors.New("Horn update: invalid state, empty update or revision overflow")
	}
	for _, batch := range []struct {
		changes []HornChange
		limit   int
	}{{update.Facts, s.program.atoms}, {update.Rules, len(s.program.rules)}} {
		seen := map[int]bool{}
		for _, change := range batch.changes {
			if change.Index < 0 || change.Index >= batch.limit || seen[change.Index] {
				return errors.New("Horn update: invalid or duplicate index")
			}
			seen[change.Index] = true
		}
	}
	return nil
}

func (s *HornState) update(update HornUpdate, mode string, fraction float64) (HornUpdateStats, error) {
	if err := s.validateUpdate(update); err != nil {
		return HornUpdateStats{}, err
	}
	p := s.program
	for _, change := range update.Facts {
		s.snapshot.Facts[change.Index] = change.Present
	}
	for _, change := range update.Rules {
		s.snapshot.Enabled[change.Index] = change.Present
	}
	s.snapshot.Revision++
	stats := HornUpdateStats{Mode: mode}
	full := func() (HornUpdateStats, error) {
		s.known, stats.InspectedRules = p.fullClosure(s.snapshot)
		stats.AffectedAtoms, stats.RecomputedAtoms = -1, p.atoms
		if mode == "adaptive" {
			stats.Mode = "adaptive-full"
		}
		return stats, nil
	}
	if mode == "full" {
		return full()
	}
	limit := math.MaxInt
	if mode == "adaptive" {
		limit = max(1, int(fraction*float64(p.atoms+p.edges)))
		stats.Mode = "adaptive-incremental"
	}
	affected := make([]bool, p.atoms)
	queue := []int{}
	add := func(atom int) bool {
		if !affected[atom] {
			if stats.DiscoveryVisits >= limit {
				return false
			}
			stats.DiscoveryVisits++
			affected[atom] = true
			queue = append(queue, atom)
		}
		return true
	}
	for _, change := range update.Facts {
		if !add(change.Index) {
			return full()
		}
	}
	for _, change := range update.Rules {
		if !add(p.rules[change.Index].Conclusion) {
			return full()
		}
	}
	for at := 0; at < len(queue); at++ {
		for _, index := range p.consumers[queue[at]] {
			if stats.DiscoveryVisits >= limit {
				return full()
			}
			stats.DiscoveryVisits++
			if !add(p.rules[index].Conclusion) {
				return full()
			}
		}
	}
	stats.AffectedAtoms, stats.RecomputedAtoms = len(queue), len(queue)
	for _, atom := range p.topology {
		if !affected[atom] {
			continue
		}
		s.known[atom] = s.snapshot.Facts[atom]
		if s.known[atom] {
			continue
		}
		for _, index := range p.producers[atom] {
			stats.InspectedRules++
			if !s.snapshot.Enabled[index] {
				continue
			}
			holds := true
			for _, premise := range p.rules[index].Premises {
				if !s.known[premise] {
					holds = false
					break
				}
			}
			if holds {
				s.known[atom] = true
				break
			}
		}
	}
	return stats, nil
}
