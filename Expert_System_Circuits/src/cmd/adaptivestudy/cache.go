package main

import (
	"fmt"
	"math/rand"
	"slices"
	"time"

	evidence "github.com/nick-fedchik/articles/Expert_System_Circuits/src"
)

type cacheEventRow struct {
	Family, Kind      string
	Rules, Seed, Step int
	NS                int64
}

type cacheQueryRow struct {
	Family, Invalidation     string
	Rules, Seed, Step, Query int
	ColdNS, CachedNS         int64
	Expected, Proved, Reused bool
}

type cacheSummary struct {
	Family                       string
	Rules, Queries, Proved       int
	Hits, Invalidations          int
	MedianColdNS, MedianCachedNS int64
	P95ColdNS, P95CachedNS       int64
	TotalColdNS, TotalCachedNS   int64
	EventNS                      int64
}

type cacheReport struct {
	Events, Queries, Hits, Invalidations int
	Summary                              []cacheSummary
}

type cacheCycle struct {
	root, rule int
	present    bool
}

func applyCacheEvent(ledger *evidence.HornLedger, c evidence.HornCase, step int, rng *rand.Rand, cycle cacheCycle) (string, error) {
	snapshot := ledger.Snapshot()
	root := c.Roots[rng.Intn(len(c.Roots))]
	ruleIndex := rng.Intn(len(snapshot.Enabled))
	switch step % 10 {
	case 0:
		return "fact-withdraw-or-add", ledger.Update(evidence.HornUpdate{Facts: []evidence.HornChange{{Index: cycle.root, Present: !cycle.present}}})
	case 1:
		root = c.Roots[0]
		return "fact-toggle", ledger.Update(evidence.HornUpdate{Facts: []evidence.HornChange{{Index: root, Present: !snapshot.Facts[root]}}})
	case 2:
		return "fact-restore", ledger.Update(evidence.HornUpdate{Facts: []evidence.HornChange{{Index: cycle.root, Present: cycle.present}}})
	case 3:
		return "rule-disable", ledger.Update(evidence.HornUpdate{Rules: []evidence.HornChange{{Index: cycle.rule, Present: false}}})
	case 4:
		return "rule-enable", ledger.Update(evidence.HornUpdate{Rules: []evidence.HornChange{{Index: cycle.rule, Present: true}}})
	case 5:
		return "fact-refresh", ledger.Update(evidence.HornUpdate{Facts: []evidence.HornChange{{Index: root, Present: snapshot.Facts[root]}}})
	case 6, 8:
		if step%10 == 6 {
			ruleIndex = cycle.rule
		} else if ruleIndex == cycle.rule {
			ruleIndex = (ruleIndex + 1) % len(snapshot.Enabled)
		}
		rule, err := ledger.Rule(ruleIndex)
		if err != nil {
			return "", err
		}
		rule.Version = fmt.Sprintf("event-%d", step)
		return "rule-version", ledger.ReplaceRule(ruleIndex, rule)
	case 7:
		ruleIndex = cycle.rule
		rule, err := ledger.Rule(ruleIndex)
		if err != nil {
			return "", err
		}
		rule.Premises = []int{root}
		rule.Version = fmt.Sprintf("event-%d", step)
		return "rule-structure", ledger.ReplaceRule(ruleIndex, rule)
	case 9:
		// Restore one original immutable version, not a relabelled modified body.
		ruleIndex = cycle.rule
		original, err := c.Program.Rule(ruleIndex)
		if err != nil {
			return "", err
		}
		return "rule-rollback", ledger.ReplaceRule(ruleIndex, original)
	}
	return "", fmt.Errorf("unreachable cache event")
}

func cacheExperiment() ([]cacheEventRow, []cacheQueryRow, cacheReport, error) {
	events, queries := []cacheEventRow{}, []cacheQueryRow{}
	for _, family := range evidence.AdaptiveFamilies {
		for _, size := range []int{256, 2048} {
			for seed := 5000; seed <= 5003; seed++ {
				c, err := evidence.GenerateAdaptiveCase(family, size, int64(seed))
				if err != nil {
					return nil, nil, cacheReport{}, err
				}
				ledger, err := evidence.NewHornLedger(c.Program, c.Snapshot)
				if err != nil {
					return nil, nil, cacheReport{}, err
				}
				cache, err := evidence.NewHornProofCache(ledger)
				if err != nil {
					return nil, nil, cacheReport{}, err
				}
				targets, err := c.DerivedQueries(8, int64(seed+1))
				if err != nil {
					return nil, nil, cacheReport{}, err
				}
				rng := rand.New(rand.NewSource(int64(seed)))
				cycle := cacheCycle{}
				for step := 0; step < 100; step++ {
					if step%10 == 0 {
						snapshot := ledger.Snapshot()
						root := c.Roots[rng.Intn(len(c.Roots))]
						cycle = cacheCycle{root, rng.Intn(len(snapshot.Enabled)), snapshot.Facts[root]}
					}
					start := time.Now()
					kind, err := applyCacheEvent(ledger, c, step, rng, cycle)
					eventNS := time.Since(start).Nanoseconds()
					if err != nil {
						return nil, nil, cacheReport{}, err
					}
					events = append(events, cacheEventRow{family, kind, size, seed, step, eventNS})
					expected := ledger.Known()
					// The ledger result is independently checked against full recomputation.
					full, err := ledger.FullClosure()
					if err != nil || !slices.Equal(expected, full) {
						return nil, nil, cacheReport{}, fmt.Errorf("versioned ledger mismatch: %v", err)
					}
					for _, query := range targets {
						var cold evidence.HornSearchResult
						var cached evidence.HornCachedResult
						var coldNS, cachedNS int64
						for _, mode := range rng.Perm(2) {
							start := time.Now()
							if mode == 0 {
								cold, err = ledger.Search(query, 100000000)
								coldNS = time.Since(start).Nanoseconds()
							} else {
								cached, err = cache.Resolve(query, 100000000)
								cachedNS = time.Since(start).Nanoseconds()
							}
							if err != nil {
								return nil, nil, cacheReport{}, err
							}
						}
						if cold.Exhausted || cached.Search.Exhausted || cold.Proved != expected[query] || cached.Search.Proved != expected[query] {
							return nil, nil, cacheReport{}, fmt.Errorf("cache acceptance differs from current closure: %s/%d/%d/%d", family, size, step, query)
						}
						queries = append(queries, cacheQueryRow{
							family, cached.Invalidation, size, seed, step, query, coldNS, cachedNS,
							expected[query], cached.Search.Proved, cached.Reused,
						})
					}
				}
			}
		}
	}
	report := cacheReport{Events: len(events), Queries: len(queries)}
	for _, family := range evidence.AdaptiveFamilies {
		for _, size := range []int{256, 2048} {
			summary := cacheSummary{Family: family, Rules: size}
			cold, cached := []int64{}, []int64{}
			for _, row := range queries {
				if row.Family != family || row.Rules != size {
					continue
				}
				summary.Queries++
				if row.Proved {
					summary.Proved++
				}
				if row.Reused {
					summary.Hits++
					report.Hits++
				}
				if row.Invalidation != "" {
					summary.Invalidations++
					report.Invalidations++
				}
				cold, cached = append(cold, row.ColdNS), append(cached, row.CachedNS)
				summary.TotalColdNS += row.ColdNS
				summary.TotalCachedNS += row.CachedNS
			}
			for _, row := range events {
				if row.Family == family && row.Rules == size {
					summary.EventNS += row.NS
				}
			}
			summary.MedianColdNS, summary.MedianCachedNS = percentile(cold, 50), percentile(cached, 50)
			summary.P95ColdNS, summary.P95CachedNS = percentile(cold, 95), percentile(cached, 95)
			report.Summary = append(report.Summary, summary)
		}
	}
	return events, queries, report, nil
}
