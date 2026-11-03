package evidence

import (
	"errors"
	"math"
	"slices"
	"sort"
)

type CleanCacheMode string

const (
	CleanRecompute CleanCacheMode = "RECOMPUTE"
	CleanRevision  CleanCacheMode = "GLOBAL_REVISION"
	CleanCone      CleanCacheMode = "OPPOSITION_CONE"
)

type cleanEpoch struct {
	index int
	epoch uint64
}

type cleanWitnessEntry struct {
	witness      CleanWitnessCertificate
	facts, rules []cleanEpoch
	age          uint64
}

type CleanCacheResult struct {
	Answer                             CleanAnswer
	Reused, Evicted                    bool
	Invalidation                       string
	FactDependencies, RuleDependencies int
}

type CleanWitnessCache struct {
	program                *CleanProgram
	snapshot               HornSnapshot
	factEpochs, ruleEpochs []uint64
	mode                   CleanCacheMode
	entries                map[int]cleanWitnessEntry
	clock                  uint64
}

func NewCleanWitnessCache(p *CleanProgram, world SignedWorld, mode CleanCacheMode) (*CleanWitnessCache, error) {
	if mode != CleanRecompute && mode != CleanRevision && mode != CleanCone {
		return nil, errors.New("clean cache: unknown mode")
	}
	snapshot, err := p.snapshot(world, SignedLiteral{})
	if err != nil {
		return nil, err
	}
	cache := &CleanWitnessCache{program: p, snapshot: cloneHornSnapshot(snapshot), mode: mode, entries: map[int]cleanWitnessEntry{}}
	if mode == CleanCone {
		cache.factEpochs, cache.ruleEpochs = make([]uint64, p.core.atoms), make([]uint64, len(p.core.rules))
		for i := range cache.factEpochs {
			cache.factEpochs[i] = 1
		}
		for i := range cache.ruleEpochs {
			cache.ruleEpochs[i] = 1
		}
	}
	return cache, nil
}

func (c *CleanWitnessCache) World() (SignedWorld, error) {
	if c == nil || c.program == nil {
		return SignedWorld{}, errors.New("clean cache: uninitialized owner")
	}
	snapshot := c.snapshot
	world := SignedWorld{Revision: snapshot.Revision, Facts: make([]State, c.program.atoms), Enabled: slices.Clone(snapshot.Enabled)}
	for atom := range world.Facts {
		if snapshot.Facts[2*atom] {
			world.Facts[atom] |= True
		}
		if snapshot.Facts[2*atom+1] {
			world.Facts[atom] |= False
		}
	}
	return world, nil
}

func (c *CleanWitnessCache) Update(update HornUpdate) error {
	if c == nil || c.program == nil {
		return errors.New("clean cache: uninitialized owner")
	}
	if err := validateHornUpdate(c.program.core, c.snapshot, update); err != nil {
		return err
	}
	if c.mode == CleanCone {
		for _, change := range update.Facts {
			if c.factEpochs[change.Index] == math.MaxUint64 {
				return errors.New("clean cache: fact epoch overflow")
			}
		}
		for _, change := range update.Rules {
			if c.ruleEpochs[change.Index] == math.MaxUint64 {
				return errors.New("clean cache: rule epoch overflow")
			}
		}
	}
	for _, change := range update.Facts {
		c.snapshot.Facts[change.Index] = change.Present
		if c.mode == CleanCone {
			c.factEpochs[change.Index]++
		}
	}
	for _, change := range update.Rules {
		c.snapshot.Enabled[change.Index] = change.Present
		if c.mode == CleanCone {
			c.ruleEpochs[change.Index]++
		}
	}
	c.snapshot.Revision++
	return nil
}

func (c *CleanWitnessCache) dependencies(witness CleanWitnessCertificate) ([]cleanEpoch, []cleanEpoch) {
	facts, rules := map[int]bool{}, map[int]bool{}
	for _, code := range witness.Proof.Facts {
		facts[code] = true
	}
	for _, index := range witness.Proof.Rules {
		rules[index] = true
	}
	for _, value := range witness.Absence {
		facts[value.Code] = true
		for _, index := range c.program.core.producers[value.Code] {
			rules[index] = true
		}
	}
	collect := func(indices map[int]bool, epochs []uint64) []cleanEpoch {
		out := make([]cleanEpoch, 0, len(indices))
		for index := range indices {
			out = append(out, cleanEpoch{index, epochs[index]})
		}
		sort.Slice(out, func(i, j int) bool { return out[i].index < out[j].index })
		return out
	}
	return collect(facts, c.factEpochs), collect(rules, c.ruleEpochs)
}

func (c *CleanWitnessCache) invalidation(entry cleanWitnessEntry) string {
	if c.mode == CleanRevision {
		if entry.witness.Revision != c.snapshot.Revision {
			return "GLOBAL_REVISION"
		}
		return ""
	}
	for _, stamp := range entry.facts {
		if c.factEpochs[stamp.index] != stamp.epoch {
			return "FACT_EPOCH"
		}
	}
	for _, stamp := range entry.rules {
		if c.ruleEpochs[stamp.index] != stamp.epoch {
			return "RULE_EPOCH"
		}
	}
	return ""
}

func (c *CleanWitnessCache) Resolve(query SignedLiteral) (CleanCacheResult, error) {
	if c == nil || c.program == nil || !query.valid(c.program.atoms) {
		return CleanCacheResult{}, errors.New("clean cache: uninitialized owner or invalid query")
	}
	world, err := c.World()
	if err != nil {
		return CleanCacheResult{}, err
	}
	result := CleanCacheResult{}
	if entry, exists := c.entries[query.code()]; exists {
		result.Invalidation = c.invalidation(entry)
		if result.Invalidation == "" {
			witness := cloneCleanWitness(entry.witness)
			witness.Revision, witness.Proof.Revision = world.Revision, world.Revision
			answer := CleanAnswer{Status: CleanSupported, Witness: &witness}
			if err := c.program.VerifyCleanAnswer(world, query, answer); err != nil {
				return CleanCacheResult{}, err
			}
			return CleanCacheResult{Answer: answer, Reused: true,
				FactDependencies: len(entry.facts), RuleDependencies: len(entry.rules)}, nil
		}
		delete(c.entries, query.code())
	}
	result.Answer, err = c.program.FreshCleanAnswer(world, query)
	if err != nil {
		return CleanCacheResult{}, err
	}
	if c.mode == CleanRecompute || result.Answer.Status != CleanSupported {
		return result, nil
	}
	if c.clock == math.MaxUint64 {
		return CleanCacheResult{}, errors.New("clean cache: insertion counter overflow")
	}
	if len(c.entries) >= 128 {
		oldest, age := -1, uint64(math.MaxUint64)
		for code, entry := range c.entries {
			if entry.age < age {
				oldest, age = code, entry.age
			}
		}
		delete(c.entries, oldest)
		result.Evicted = true
	}
	c.clock++
	witness := cloneCleanWitness(*result.Answer.Witness)
	var facts, rules []cleanEpoch
	if c.mode == CleanCone {
		facts, rules = c.dependencies(witness)
	}
	c.entries[query.code()] = cleanWitnessEntry{witness, facts, rules, c.clock}
	result.FactDependencies, result.RuleDependencies = len(facts), len(rules)
	return result, nil
}
