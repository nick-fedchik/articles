package evidence

import (
	"fmt"
	"math"
	"testing"
)

func TestVersionABACannotReviveOriginalTicket(t *testing.T) {
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
	rule.Version = "2"
	if err := ledger.ReplaceRule(0, rule); err != nil {
		t.Fatal(err)
	}
	rule.Version = "1"
	if err := ledger.ReplaceRule(0, rule); err != nil {
		t.Fatal(err)
	}
	result, err := cache.Resolve(2, 100000)
	if err != nil || result.Reused || !result.Search.Proved || result.Invalidation != "RULE_REVISION" {
		t.Fatal("version ABA revived a previously cached proof", result, err)
	}
}

func TestLedgerEpochOverflowDoesNotPartiallyUpdate(t *testing.T) {
	p, snapshot := hornFixture(t)
	ledger, err := NewHornLedger(p, snapshot)
	if err != nil {
		t.Fatal(err)
	}
	ledger.factEpochs[0] = math.MaxUint64
	if err := ledger.Update(HornUpdate{Facts: []HornChange{{0, false}}}); err == nil ||
		!ledger.Snapshot().Facts[0] || ledger.Snapshot().Revision != snapshot.Revision {
		t.Fatal("epoch overflow partially applied an update")
	}
}

func FuzzHornCacheSequence(f *testing.F) {
	f.Add([]byte{0, 1, 2, 3, 4, 5, 0, 0})
	f.Fuzz(func(t *testing.T, events []byte) {
		if len(events) > 128 {
			t.Skip("bounded versioned sequence")
		}
		p, snapshot := hornFixture(t)
		ledger, err := NewHornLedger(p, snapshot)
		if err != nil {
			t.Fatal(err)
		}
		cache, err := NewHornProofCache(ledger)
		if err != nil {
			t.Fatal(err)
		}
		for step, event := range events {
			current := ledger.Snapshot()
			switch event % 3 {
			case 0:
				atom := int(event>>2) % 4
				err = ledger.Update(HornUpdate{Facts: []HornChange{{atom, !current.Facts[atom]}}})
			case 1:
				index := int(event>>2) % 2
				err = ledger.Update(HornUpdate{Rules: []HornChange{{index, !current.Enabled[index]}}})
			case 2:
				rule, ruleErr := ledger.Rule(0)
				if ruleErr != nil {
					t.Fatal(ruleErr)
				}
				rule.Version = fmt.Sprintf("fuzz-%d", step)
				if event&4 == 0 {
					rule.Premises = []int{0, 1}
				} else {
					rule.Premises = []int{3}
				}
				err = ledger.ReplaceRule(0, rule)
			}
			if err != nil {
				t.Fatal(err)
			}
			closure, err := ledger.FullClosure()
			if err != nil {
				t.Fatal(err)
			}
			result, err := cache.Resolve(2, 100000)
			if err != nil || result.Search.Exhausted || result.Search.Proved != closure[2] {
				t.Fatalf("step %d cache differs from current closure: %+v %v", step, result, err)
			}
		}
	})
}
