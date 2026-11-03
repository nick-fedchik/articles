package evidence

import (
	"slices"
	"testing"
)

func TestCacheRebindsOnlyIndependentValidProofs(t *testing.T) {
	p, snapshot := hornFixture(t)
	ledger, err := NewHornLedger(p, snapshot)
	if err != nil {
		t.Fatal(err)
	}
	cache, err := NewHornProofCache(ledger)
	if err != nil {
		t.Fatal(err)
	}
	first, err := cache.Resolve(2, 100000)
	if err != nil || !first.Search.Proved || first.Reused {
		t.Fatal(first, err)
	}
	if err := ledger.Update(HornUpdate{Facts: []HornChange{{3, true}}}); err != nil {
		t.Fatal(err)
	}
	hit, err := cache.Resolve(2, 100000)
	if err != nil || !hit.Reused || hit.Search.Proof.Revision == first.Search.Proof.Revision {
		t.Fatal("unrelated update destroyed an independent proof", hit, err)
	}
	if err := ledger.state.program.VerifyProof(ledger.Snapshot(), 2, first.Search.Proof); err == nil {
		t.Fatal("old certificate bypassed global-revision validation")
	}
	if err := ledger.Update(HornUpdate{Facts: []HornChange{{0, false}}}); err != nil {
		t.Fatal(err)
	}
	alternative, err := cache.Resolve(2, 100000)
	if err != nil || !alternative.Search.Proved || alternative.Reused || alternative.Invalidation != "FACT_RETRACTED" {
		t.Fatal("missing premise was silently reused instead of finding an alternative", alternative, err)
	}
	if err := ledger.Update(HornUpdate{Facts: []HornChange{{3, false}}}); err != nil {
		t.Fatal(err)
	}
	missing, err := cache.Resolve(2, 100000)
	if err != nil || missing.Search.Proved || missing.Reused {
		t.Fatal("stale acceptance after both proofs disappeared", missing, err)
	}
}

func TestCacheRejectsFactAndRuleABA(t *testing.T) {
	for _, kind := range []string{"fact", "rule"} {
		p, snapshot := hornFixture(t)
		ledger, err := NewHornLedger(p, snapshot)
		if err != nil {
			t.Fatal(err)
		}
		cache, err := NewHornProofCache(ledger)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := cache.Resolve(2, 100000); err != nil {
			t.Fatal(err)
		}
		for _, present := range []bool{false, true} {
			event := HornUpdate{Facts: []HornChange{{0, present}}}
			if kind == "rule" {
				event = HornUpdate{Rules: []HornChange{{0, present}}}
			}
			if err := ledger.Update(event); err != nil {
				t.Fatal(err)
			}
		}
		got, err := cache.Resolve(2, 100000)
		if err != nil || got.Reused || !got.Search.Proved || got.Invalidation == "" {
			t.Fatal("withdraw-and-restore revived an obsolete dependency stamp", kind, got, err)
		}
	}
}

func TestVersionReplacementIsAtomicAndImmutable(t *testing.T) {
	p, snapshot := hornFixture(t)
	ledger, err := NewHornLedger(p, snapshot)
	if err != nil {
		t.Fatal(err)
	}
	cache, err := NewHornProofCache(ledger)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := cache.Resolve(2, 100000); err != nil {
		t.Fatal(err)
	}
	rule, err := ledger.Rule(0)
	if err != nil {
		t.Fatal(err)
	}
	bad := copyHornRule(rule)
	bad.Premises = []int{3}
	before := ledger.Snapshot()
	if err := ledger.ReplaceRule(0, bad); err == nil || ledger.Snapshot().Revision != before.Revision {
		t.Fatal("same version acquired another body")
	}
	rule.Version = "2"
	if err := ledger.ReplaceRule(0, rule); err != nil {
		t.Fatal(err)
	}
	changed, err := cache.Resolve(2, 100000)
	if err != nil || changed.Reused || changed.Invalidation != "RULE_VERSION" || !changed.Search.Proved {
		t.Fatal(changed, err)
	}
	rule.Version = "1"
	if err := ledger.ReplaceRule(0, rule); err != nil {
		t.Fatal(err)
	}
	reverted, err := cache.Resolve(2, 100000)
	if err != nil || reverted.Reused || !reverted.Search.Proved {
		t.Fatal("version rollback revived a cached certificate", reverted, err)
	}
	oldClosure, err := p.FullClosure(snapshot)
	if err != nil || !slices.Equal(oldClosure, []bool{true, true, true, false}) {
		t.Fatal("immutable old theory was mutated")
	}
}

func TestCacheDoesNotHideVerifierFailure(t *testing.T) {
	p, snapshot := hornFixture(t)
	ledger, err := NewHornLedger(p, snapshot)
	if err != nil {
		t.Fatal(err)
	}
	cache, err := NewHornProofCache(ledger)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := cache.Resolve(2, 100000); err != nil {
		t.Fatal(err)
	}
	entry := cache.entries[2]
	entry.rules = append(entry.rules, entry.rules[0])
	cache.entries[2] = entry
	if _, err := cache.Resolve(2, 100000); err == nil {
		t.Fatal("forged duplicate rule was accepted or silently replaced")
	}
}
