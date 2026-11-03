package evidence

import (
	"errors"
	"math"
	"slices"
)

type hornRuleKey struct{ id, version string }

type HornLedger struct {
	state       *HornState
	factEpochs  []uint64
	ruleEpochs  []uint64
	ruleIndices map[hornRuleKey]int
	history     map[hornRuleKey]HornRule
}

func copyHornRule(rule HornRule) HornRule {
	rule.Premises = append([]int(nil), rule.Premises...)
	return rule
}

func sameHornRule(a, b HornRule) bool {
	return a.ID == b.ID && a.Version == b.Version && a.Conclusion == b.Conclusion && slices.Equal(a.Premises, b.Premises)
}

func NewHornLedger(program *HornProgram, snapshot HornSnapshot) (*HornLedger, error) {
	state, err := NewHornState(program, snapshot)
	if err != nil {
		return nil, err
	}
	l := &HornLedger{
		state: state, factEpochs: make([]uint64, program.atoms), ruleEpochs: make([]uint64, len(program.rules)),
		ruleIndices: map[hornRuleKey]int{}, history: map[hornRuleKey]HornRule{},
	}
	for i := range l.factEpochs {
		l.factEpochs[i] = 1
	}
	for i, rule := range program.rules {
		l.ruleEpochs[i] = 1
		key := hornRuleKey{rule.ID, rule.Version}
		l.ruleIndices[key], l.history[key] = i, copyHornRule(rule)
	}
	return l, nil
}

func (l *HornLedger) Snapshot() HornSnapshot { return l.state.Snapshot() }
func (l *HornLedger) Known() []bool          { return l.state.Known() }

func (l *HornLedger) FullClosure() ([]bool, error) {
	if l == nil || l.state == nil {
		return nil, errors.New("Horn ledger: uninitialized")
	}
	return l.state.program.FullClosure(l.state.snapshot)
}

func (l *HornLedger) Rule(index int) (HornRule, error) {
	if l == nil || l.state == nil || index < 0 || index >= len(l.state.program.rules) {
		return HornRule{}, errors.New("Horn ledger: invalid rule index")
	}
	return copyHornRule(l.state.program.rules[index]), nil
}

func (l *HornLedger) Update(update HornUpdate) error {
	if l == nil || l.state == nil {
		return errors.New("Horn ledger: uninitialized")
	}
	if err := l.state.validateUpdate(update); err != nil {
		return err
	}
	for _, change := range update.Facts {
		if l.factEpochs[change.Index] == math.MaxUint64 {
			return errors.New("Horn ledger: fact epoch overflow")
		}
	}
	for _, change := range update.Rules {
		if l.ruleEpochs[change.Index] == math.MaxUint64 {
			return errors.New("Horn ledger: rule epoch overflow")
		}
	}
	if _, err := l.state.UpdateFull(update); err != nil {
		return err
	}
	for _, change := range update.Facts {
		l.factEpochs[change.Index]++
	}
	for _, change := range update.Rules {
		l.ruleEpochs[change.Index]++
	}
	return nil
}

func (l *HornLedger) ReplaceRule(index int, replacement HornRule) error {
	old, err := l.Rule(index)
	if err != nil {
		return err
	}
	if replacement.ID != old.ID || l.state.snapshot.Revision == math.MaxUint64 || l.ruleEpochs[index] == math.MaxUint64 {
		return errors.New("Horn ledger: replacement ID mismatch or revision overflow")
	}
	replacement = copyHornRule(replacement)
	slices.Sort(replacement.Premises)
	key := hornRuleKey{replacement.ID, replacement.Version}
	if historical, exists := l.history[key]; exists && !sameHornRule(historical, replacement) {
		return errors.New("Horn ledger: a rule version cannot be reused with a different structure")
	}
	if _, exists := l.history[key]; !exists && len(l.history) >= 100000 {
		return errors.New("Horn ledger: immutable version history limit exceeded")
	}
	rules := append([]HornRule(nil), l.state.program.rules...)
	rules[index] = replacement
	program, err := CompileHorn(l.state.program.atoms, rules)
	if err != nil {
		return err
	}
	snapshot := l.state.Snapshot()
	snapshot.Revision++
	state, err := NewHornState(program, snapshot)
	if err != nil {
		return err
	}
	indices := map[hornRuleKey]int{}
	for i, rule := range program.rules {
		indices[hornRuleKey{rule.ID, rule.Version}] = i
	}
	l.state, l.ruleIndices = state, indices
	l.ruleEpochs[index]++
	l.history[key] = copyHornRule(replacement)
	return nil
}

func (l *HornLedger) Search(query, budget int) (HornSearchResult, error) {
	if l == nil || l.state == nil {
		return HornSearchResult{}, errors.New("Horn ledger: uninitialized")
	}
	return l.state.program.Search(l.state.snapshot, query, "backward", budget)
}

type hornFactStamp struct {
	atom  int
	epoch uint64
}

type hornRuleStamp struct {
	rule  HornRule
	epoch uint64
}

type hornCacheEntry struct {
	facts []hornFactStamp
	rules []hornRuleStamp
	age   uint64
}

type HornProofCache struct {
	ledger  *HornLedger
	entries map[int]hornCacheEntry
	clock   uint64
}

type HornCachedResult struct {
	Search       HornSearchResult
	Reused       bool
	Invalidation string
}

func NewHornProofCache(ledger *HornLedger) (*HornProofCache, error) {
	if ledger == nil || ledger.state == nil {
		return nil, errors.New("Horn cache: initialized owning ledger required")
	}
	return &HornProofCache{ledger: ledger, entries: map[int]hornCacheEntry{}}, nil
}

func (c *HornProofCache) bind(query int, entry hornCacheEntry) (HornProof, string) {
	l := c.ledger
	p, snapshot := l.state.program, l.state.snapshot
	proof := HornProof{Theory: p.digest, Revision: snapshot.Revision, Query: query}
	for _, fact := range entry.facts {
		if fact.atom < 0 || fact.atom >= len(l.factEpochs) || !snapshot.Facts[fact.atom] {
			return HornProof{}, "FACT_RETRACTED"
		}
		if fact.epoch != l.factEpochs[fact.atom] {
			return HornProof{}, "FACT_REVISION"
		}
		proof.Facts = append(proof.Facts, fact.atom)
	}
	for _, stamp := range entry.rules {
		index, exists := l.ruleIndices[hornRuleKey{stamp.rule.ID, stamp.rule.Version}]
		if !exists {
			return HornProof{}, "RULE_VERSION"
		}
		if !snapshot.Enabled[index] {
			return HornProof{}, "RULE_DISABLED"
		}
		if stamp.epoch != l.ruleEpochs[index] {
			return HornProof{}, "RULE_REVISION"
		}
		if !sameHornRule(stamp.rule, p.rules[index]) {
			return HornProof{}, "RULE_STRUCTURE"
		}
		proof.Rules = append(proof.Rules, index)
	}
	return proof, ""
}

func (c *HornProofCache) Resolve(query, budget int) (HornCachedResult, error) {
	if c == nil || c.ledger == nil || c.entries == nil || budget < 1 {
		return HornCachedResult{}, errors.New("Horn cache: initialized cache and positive miss-search budget required")
	}
	l := c.ledger
	if err := l.state.program.validate(l.state.snapshot, query); err != nil {
		return HornCachedResult{}, err
	}
	out := HornCachedResult{}
	if entry, exists := c.entries[query]; exists {
		proof, reason := c.bind(query, entry)
		if reason == "" {
			if err := l.state.program.VerifyProof(l.state.snapshot, query, proof); err != nil {
				return HornCachedResult{}, err
			}
			out.Reused = true
			out.Search = HornSearchResult{Method: "cache-rebind", Proved: true, Proof: proof}
			return out, nil
		}
		out.Invalidation = reason
		delete(c.entries, query)
	}
	result, err := l.Search(query, budget)
	if err != nil {
		return HornCachedResult{}, err
	}
	out.Search = result
	if !result.Proved {
		return out, nil
	}
	if c.clock == math.MaxUint64 {
		return HornCachedResult{}, errors.New("Horn cache: insertion counter overflow")
	}
	if len(c.entries) >= 1024 {
		oldestQuery, oldestAge := -1, uint64(math.MaxUint64)
		for q, entry := range c.entries {
			if entry.age < oldestAge {
				oldestQuery, oldestAge = q, entry.age
			}
		}
		delete(c.entries, oldestQuery)
	}
	c.clock++
	entry := hornCacheEntry{age: c.clock}
	for _, atom := range result.Proof.Facts {
		entry.facts = append(entry.facts, hornFactStamp{atom, l.factEpochs[atom]})
	}
	for _, index := range result.Proof.Rules {
		entry.rules = append(entry.rules, hornRuleStamp{copyHornRule(l.state.program.rules[index]), l.ruleEpochs[index]})
	}
	c.entries[query] = entry
	return out, nil
}
