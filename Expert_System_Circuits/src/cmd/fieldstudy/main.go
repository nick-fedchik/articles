package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strconv"
	"time"

	evidence "github.com/nick-fedchik/articles/Expert_System_Circuits/src"
)

type searchSample struct {
	Family, Method                string
	Rules, Seed, Multiplier       int
	NS, Operations, Applied, Rank int64
	Expected, Proved, Exhausted   bool
}

type updateSample struct {
	Family                          string
	Rules, Seed, Percent, Round     int
	Changed, Roots, Atoms, Affected int
	FullNS, IncrementalNS           int64
}

type searchSummary struct {
	Family, Method                    string
	Rules, Multiplier, Cases          int
	Supported, Proved, Exhausted      int
	MedianNS, P95NS, MedianOperations int64
}

type updateSummary struct {
	Family                            string
	Rules, Percent, Cases             int
	MedianFullNS, MedianIncrementalNS int64
	MedianAffectedPercent             float64
}

type report struct {
	Scope                   string
	SearchCases, Updates    int
	Search                  []searchSummary
	Incremental             []updateSummary
	HeatWinningFamilies     int
	PageRankWinningFamilies int
	ContinueHeatPilot       bool
}

func quantile(values []int64, percentile int) int64 {
	copyValues := append([]int64(nil), values...)
	sort.Slice(copyValues, func(i, j int) bool { return copyValues[i] < copyValues[j] })
	return copyValues[(len(copyValues)*percentile-1)/100]
}

func measure(seeds int) ([]searchSample, []updateSample, error) {
	if seeds < 1 || seeds > 20 {
		return nil, nil, fmt.Errorf("seeds outside [1,20]")
	}
	searches, updates := []searchSample{}, []updateSample{}
	for _, family := range evidence.HornFamilies {
		for _, size := range []int{256, 2048} {
			for seed := 0; seed < seeds; seed++ {
				c, err := evidence.GenerateHornCase(family, size, int64(seed))
				if err != nil {
					return nil, nil, err
				}
				expected, err := c.Expected()
				if err != nil {
					return nil, nil, err
				}
				// Full-budget regression checks are separate from timed samples.
				for _, method := range evidence.HornMethods {
					got, err := c.Program.Search(c.Snapshot, c.Query, method, 100000000)
					if err != nil || got.Exhausted || got.Proved != expected {
						return nil, nil, fmt.Errorf("unbounded %s/%d/%d/%s: %+v %v", family, size, seed, method, got, err)
					}
				}
				for _, multiplier := range []int{4, 16, 128} {
					order := rand.New(rand.NewSource(int64(1000*seed + multiplier))).Perm(len(evidence.HornMethods))
					for _, index := range order {
						method := evidence.HornMethods[index]
						start := time.Now()
						got, err := c.Program.Search(c.Snapshot, c.Query, method, multiplier*c.Program.BaseVolume())
						ns := time.Since(start).Nanoseconds()
						if err != nil || got.Proved && !expected {
							return nil, nil, fmt.Errorf("search %s: false acceptance or error: %v", method, err)
						}
						searches = append(searches, searchSample{
							family, method, size, seed, multiplier, ns, int64(got.Operations), int64(got.Applied),
							int64(got.RankVisits), expected, got.Proved, got.Exhausted,
						})
					}
				}
				for _, percent := range []int{1, 5, 10} {
					state, err := evidence.NewHornState(c.Program, c.Snapshot)
					if err != nil {
						return nil, nil, err
					}
					for round := 0; round < 2; round++ {
						snapshot := state.Snapshot()
						count := int(math.Ceil(float64(len(c.Roots)*percent) / 100))
						order := rand.New(rand.NewSource(int64(seed*100 + percent*2 + round))).Perm(len(c.Roots))
						update := evidence.HornUpdate{}
						for _, i := range order[:count] {
							atom := c.Roots[i]
							update.Facts = append(update.Facts, evidence.HornChange{Index: atom, Present: !snapshot.Facts[atom]})
							snapshot.Facts[atom] = !snapshot.Facts[atom]
						}
						snapshot.Revision++
						start := time.Now()
						full, err := c.Program.FullClosure(snapshot)
						fullNS := time.Since(start).Nanoseconds()
						if err != nil {
							return nil, nil, err
						}
						start = time.Now()
						stats, err := state.Update(update)
						incrementalNS := time.Since(start).Nanoseconds()
						if err != nil || !slices.Equal(full, state.Known()) {
							return nil, nil, fmt.Errorf("incremental mismatch: %s/%d/%d/%d: %v", family, size, seed, percent, err)
						}
						updates = append(updates, updateSample{
							family, size, seed, percent, round, count, len(c.Roots), c.AtomCount(), stats.AffectedAtoms, fullNS, incrementalNS,
						})
					}
				}
			}
		}
	}
	return searches, updates, nil
}

func summarize(searches []searchSample, updates []updateSample) report {
	out := report{
		Scope:       "synthetic positive acyclic Horn theories; full per-query ranking and proof-check time; graph-operation budget is not CPU instructions",
		SearchCases: len(searches), Updates: len(updates),
	}
	for _, family := range evidence.HornFamilies {
		for _, size := range []int{256, 2048} {
			for _, multiplier := range []int{4, 16, 128} {
				bestCheap, heat, pageRank := int64(math.MaxInt64), int64(0), int64(0)
				heatComplete, pageComplete := false, false
				for _, method := range evidence.HornMethods {
					summary := searchSummary{Family: family, Method: method, Rules: size, Multiplier: multiplier}
					times, operations := []int64{}, []int64{}
					for _, sample := range searches {
						if sample.Family != family || sample.Rules != size || sample.Multiplier != multiplier || sample.Method != method {
							continue
						}
						summary.Cases++
						if sample.Expected {
							summary.Supported++
						}
						if sample.Proved {
							summary.Proved++
						}
						if sample.Exhausted {
							summary.Exhausted++
						}
						times = append(times, sample.NS)
						operations = append(operations, sample.Operations)
					}
					summary.MedianNS, summary.P95NS = quantile(times, 50), quantile(times, 95)
					summary.MedianOperations = quantile(operations, 50)
					out.Search = append(out.Search, summary)
					complete := summary.Proved == summary.Supported && summary.Exhausted == 0
					switch method {
					case "fifo", "backward", "greedy":
						if complete && summary.MedianNS < bestCheap {
							bestCheap = summary.MedianNS
						}
					case "heat":
						heat, heatComplete = summary.MedianNS, complete
					case "pagerank":
						pageRank, pageComplete = summary.MedianNS, complete
					}
				}
				if size == 2048 && multiplier == 128 && bestCheap != math.MaxInt64 {
					if heatComplete && float64(heat) <= .8*float64(bestCheap) {
						out.HeatWinningFamilies++
					}
					if pageComplete && float64(pageRank) <= .8*float64(bestCheap) {
						out.PageRankWinningFamilies++
					}
				}
			}
			for _, percent := range []int{1, 5, 10} {
				summary := updateSummary{Family: family, Rules: size, Percent: percent}
				full, incremental, affected := []int64{}, []int64{}, []int64{}
				for _, sample := range updates {
					if sample.Family != family || sample.Rules != size || sample.Percent != percent {
						continue
					}
					summary.Cases++
					full = append(full, sample.FullNS)
					incremental = append(incremental, sample.IncrementalNS)
					affected = append(affected, int64(1000000*sample.Affected/sample.Atoms))
				}
				summary.MedianFullNS, summary.MedianIncrementalNS = quantile(full, 50), quantile(incremental, 50)
				summary.MedianAffectedPercent = float64(quantile(affected, 50)) / 10000
				out.Incremental = append(out.Incremental, summary)
			}
		}
	}
	out.ContinueHeatPilot = out.HeatWinningFamilies >= 2
	return out
}

func writeCSV(dir string, searches []searchSample, updates []updateSample) error {
	if dir == "" {
		return nil
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	write := func(name string, emit func(*csv.Writer) error) error {
		file, err := os.Create(filepath.Join(dir, name))
		if err != nil {
			return err
		}
		writer := csv.NewWriter(file)
		if err := emit(writer); err != nil {
			file.Close()
			return err
		}
		writer.Flush()
		if err := writer.Error(); err != nil {
			file.Close()
			return err
		}
		return file.Close()
	}
	if err := write("search.csv", func(w *csv.Writer) error {
		if err := w.Write([]string{"family", "rules", "seed", "method", "budget_multiplier", "total_ns", "graph_operations", "applied_rules", "rank_edge_visits", "query_supported", "proved", "exhausted"}); err != nil {
			return err
		}
		for _, x := range searches {
			if err := w.Write([]string{x.Family, strconv.Itoa(x.Rules), strconv.Itoa(x.Seed), x.Method,
				strconv.Itoa(x.Multiplier), strconv.FormatInt(x.NS, 10), strconv.FormatInt(x.Operations, 10),
				strconv.FormatInt(x.Applied, 10), strconv.FormatInt(x.Rank, 10),
				strconv.FormatBool(x.Expected), strconv.FormatBool(x.Proved), strconv.FormatBool(x.Exhausted)}); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}
	return write("updates.csv", func(w *csv.Writer) error {
		if err := w.Write([]string{"family", "rules", "seed", "percent", "round", "changed_roots", "total_used_roots", "atoms", "affected_atoms", "full_ns", "incremental_ns"}); err != nil {
			return err
		}
		for _, x := range updates {
			if err := w.Write([]string{x.Family, strconv.Itoa(x.Rules), strconv.Itoa(x.Seed), strconv.Itoa(x.Percent), strconv.Itoa(x.Round),
				strconv.Itoa(x.Changed), strconv.Itoa(x.Roots), strconv.Itoa(x.Atoms), strconv.Itoa(x.Affected), strconv.FormatInt(x.FullNS, 10), strconv.FormatInt(x.IncrementalNS, 10)}); err != nil {
				return err
			}
		}
		return nil
	})
}

func run(seeds int, csvDir string, writer io.Writer) error {
	searches, updates, err := measure(seeds)
	if err != nil {
		return err
	}
	if err := writeCSV(csvDir, searches, updates); err != nil {
		return err
	}
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(summarize(searches, updates))
}

func main() {
	seeds := flag.Int("seeds", 10, "fixed graph seeds, from zero; maximum 20")
	csvDir := flag.String("csv-dir", "", "explicit directory for reproducible per-case CSV evidence; files are replaced")
	flag.Parse()
	if err := run(*seeds, *csvDir, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
