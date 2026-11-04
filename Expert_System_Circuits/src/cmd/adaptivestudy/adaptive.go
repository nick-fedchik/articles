package main

import (
	"fmt"
	"math/rand"
	"slices"
	"time"

	evidence "github.com/nick-fedchik/articles/Expert_System_Circuits/src"
)

type adaptiveRow struct {
	Split, Family, Mode, Path  string
	Rules, Seed, Step, Percent int
	Kind                       string
	Changed, Total             int
	Fraction                   float64
	NS                         int64
	Discovery, Recomputed      int
}

type updateMode struct {
	name     string
	fraction float64
	state    *evidence.HornState
}

func measureAdaptive(split string, families []string, sizes, seeds []int, fractions []float64) ([]adaptiveRow, error) {
	rows := []adaptiveRow{}
	for _, family := range families {
		for _, size := range sizes {
			for _, seed := range seeds {
				c, err := evidence.GenerateAdaptiveCase(family, size, int64(seed))
				if err != nil {
					return nil, err
				}
				plan, err := evidence.PlanAdaptiveUpdates(c, int64(seed+90000))
				if err != nil {
					return nil, err
				}
				modes := []updateMode{{name: "full"}, {name: "incremental"}}
				for _, fraction := range fractions {
					modes = append(modes, updateMode{name: "adaptive", fraction: fraction})
				}
				for i := range modes {
					modes[i].state, err = evidence.NewHornState(c.Program, c.Snapshot)
					if err != nil {
						return nil, err
					}
				}
				for step, event := range plan {
					order := rand.New(rand.NewSource(int64(seed*100 + step))).Perm(len(modes))
					for _, i := range order {
						mode := modes[i]
						start := time.Now()
						var stats evidence.HornUpdateStats
						var err error
						switch mode.name {
						case "full":
							stats, err = mode.state.UpdateFull(event.Update)
						case "incremental":
							stats, err = mode.state.Update(event.Update)
						case "adaptive":
							stats, err = mode.state.UpdateAdaptive(event.Update, mode.fraction)
						}
						ns := time.Since(start).Nanoseconds()
						if err != nil {
							return nil, err
						}
						total := len(c.Roots)
						if event.Kind == "rules" {
							total = len(c.Snapshot.Enabled)
						}
						rows = append(rows, adaptiveRow{
							split, family, mode.name, stats.Mode, size, seed, step, event.Percent, event.Kind,
							len(event.Update.Facts) + len(event.Update.Rules), total, mode.fraction, ns,
							stats.DiscoveryVisits, stats.RecomputedAtoms,
						})
					}
					want := modes[0].state.Known()
					for _, mode := range modes[1:] {
						if !slices.Equal(want, mode.state.Known()) {
							return nil, fmt.Errorf("adaptive mismatch %s/%d/%d/%d/%s", family, size, seed, step, mode.name)
						}
					}
				}
			}
		}
	}
	return rows, nil
}

type adaptiveSummary struct {
	Family, Mode           string
	Rules, Cases           int
	MedianNS, P95NS        int64
	FullFallbacks          int
	MaxObservedRatioToFull float64
}

type adaptiveReport struct {
	SelectedFraction float64
	TrainingTotalsNS map[string]int64
	Measurements     int
	Test             []adaptiveSummary
	PassPilot        bool
}

func adaptiveExperiment() ([]adaptiveRow, adaptiveReport, error) {
	fractions := []float64{.025, .05, .1, .2}
	train, err := measureAdaptive("train", []string{"clusters", "layered"}, []int{512, 2048}, []int{1000, 1001, 1002, 1003}, fractions)
	if err != nil {
		return nil, adaptiveReport{}, err
	}
	totals := map[string]int64{}
	for _, row := range train {
		if row.Mode == "adaptive" {
			totals[fmt.Sprint(row.Fraction)] += row.NS
		}
	}
	selected := fractions[0]
	for _, fraction := range fractions[1:] {
		if totals[fmt.Sprint(fraction)] < totals[fmt.Sprint(selected)] {
			selected = fraction
		}
	}
	test, err := measureAdaptive("test", evidence.AdaptiveFamilies, []int{1024, 4096},
		[]int{2000, 2001, 2002, 2003, 2004, 2005, 2006, 2007}, []float64{selected})
	if err != nil {
		return nil, adaptiveReport{}, err
	}
	report := adaptiveReport{SelectedFraction: selected, TrainingTotalsNS: totals, Measurements: len(train) + len(test), PassPilot: true}
	type pair struct {
		family           string
		size, seed, step int
	}
	full := map[pair]int64{}
	for _, row := range test {
		if row.Mode == "full" {
			full[pair{row.Family, row.Rules, row.Seed, row.Step}] = row.NS
		}
	}
	for _, family := range evidence.AdaptiveFamilies {
		familyFull, familyAdaptive := int64(0), int64(0)
		var clusterFullMedian, clusterAdaptiveMedian int64
		for _, size := range []int{1024, 4096} {
			for _, mode := range []string{"full", "incremental", "adaptive"} {
				summary := adaptiveSummary{Family: family, Mode: mode, Rules: size}
				times := []int64{}
				for _, row := range test {
					if row.Family != family || row.Rules != size || row.Mode != mode {
						continue
					}
					summary.Cases++
					times = append(times, row.NS)
					if row.Path == "adaptive-full" {
						summary.FullFallbacks++
					}
					ratio := float64(row.NS) / float64(max(1, full[pair{row.Family, row.Rules, row.Seed, row.Step}]))
					if ratio > summary.MaxObservedRatioToFull {
						summary.MaxObservedRatioToFull = ratio
					}
					if mode == "full" {
						familyFull += row.NS
					} else if mode == "adaptive" {
						familyAdaptive += row.NS
					}
				}
				summary.MedianNS, summary.P95NS = percentile(times, 50), percentile(times, 95)
				report.Test = append(report.Test, summary)
				if family == "clusters" && size == 4096 {
					if mode == "full" {
						clusterFullMedian = summary.MedianNS
					} else if mode == "adaptive" {
						clusterAdaptiveMedian = summary.MedianNS
					}
				}
			}
		}
		if float64(familyAdaptive) > 1.25*float64(familyFull) {
			report.PassPilot = false
		}
		if family == "clusters" && float64(clusterAdaptiveMedian) > .8*float64(clusterFullMedian) {
			report.PassPilot = false
		}
	}
	return append(train, test...), report, nil
}
