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
	"runtime"
	"sort"
	"strconv"
	"time"

	evidence "github.com/nick-fedchik/articles/Epistemic_Circuit_Theory/src"
)

const managedBudget = 64 << 20

var formats = []evidence.ProvenanceFormat{
	evidence.ProvenanceDense, evidence.ProvenanceSparse, evidence.ProvenanceInterned, evidence.ProvenanceTrace,
}

type buildRow struct {
	Family, Regime         string
	Universe, Seed, Repeat int
	Format                 evidence.ProvenanceFormat
	NS, AllocBytes, Allocs uint64
	Stats                  evidence.ProvenanceStats
}

type weightRow struct {
	Family, Regime, Kind   string
	Universe, Seed, Repeat int
	Format                 evidence.ProvenanceFormat
	NS, AllocBytes, Allocs uint64
	Active                 bool
	Weight                 float64
}

type measurement struct {
	ns, bytes, allocs uint64
	snapshot          *evidence.ProvenanceSnapshot
	result            evidence.Result
}

func measure(operation func() (*evidence.ProvenanceSnapshot, error)) (measurement, error) {
	var before, after runtime.MemStats
	runtime.ReadMemStats(&before)
	start := time.Now()
	snapshot, err := operation()
	var result evidence.Result
	if err == nil {
		result, err = snapshot.Result()
	}
	ns := time.Since(start).Nanoseconds()
	runtime.ReadMemStats(&after)
	if err != nil {
		return measurement{}, err
	}
	return measurement{uint64(ns), after.TotalAlloc - before.TotalAlloc, after.Mallocs - before.Mallocs, snapshot, result}, nil
}

func sameResult(a, b evidence.Result) bool {
	if a.State != b.State || len(a.Sources) != len(b.Sources) || math.IsNaN(a.Weight) || math.IsInf(a.Weight, 0) ||
		math.IsNaN(b.Weight) || math.IsInf(b.Weight, 0) {
		return false
	}
	if b.Weight == 0 {
		if a.Weight != 0 {
			return false
		}
	} else if math.Abs(a.Weight-b.Weight) > 1e-12*math.Abs(b.Weight) {
		return false
	}
	counts := map[evidence.Source]int{}
	for _, source := range b.Sources {
		counts[source]++
	}
	for _, source := range a.Sources {
		if counts[source] == 0 {
			return false
		}
		counts[source]--
	}
	return true
}

type weightStep struct {
	kind   string
	source evidence.Source
	weight float64
	active bool
	want   evidence.Result
}

func weightPlan(c evidence.ProvenanceCase, original evidence.Result) ([]weightStep, bool, error) {
	active := original.Sources[0]
	plan := []weightStep{{kind: "active-weight", source: active, weight: .8, active: true}}
	type identity struct{ id, version string }
	present := map[identity]bool{}
	for _, source := range original.Sources {
		present[identity{source.ID, source.Version}] = true
	}
	inactive := false
	for _, source := range c.Sources {
		if !present[identity{source.ID, source.Version}] {
			plan = append(plan, weightStep{kind: "inactive-weight", source: source, weight: .5})
			inactive = true
			break
		}
	}
	plan = append(plan, weightStep{kind: "active-zero", source: active, weight: 0, active: true},
		weightStep{kind: "active-restore", source: active, weight: active.Weight, active: true})
	program := c.Program
	for i := range plan {
		var err error
		program, err = program.WithSourceWeight(plan[i].source.ID, plan[i].source.Version, plan[i].weight)
		if err != nil {
			return nil, false, err
		}
		plan[i].want, err = program.Evaluate(c.Inputs)
		if err != nil {
			return nil, false, err
		}
	}
	return plan, inactive, nil
}

func runMeasurements() ([]buildRow, []weightRow, int, error) {
	builds, updates := []buildRow{}, []weightRow{}
	skippedInactive := 0
	for _, family := range evidence.ProvenanceFamilies {
		for _, universe := range []int{256, 1024, 4096} {
			for _, regime := range []string{"resolved", "uncertain"} {
				for seed := 81000; seed <= 81002; seed++ {
					c, err := evidence.GenerateProvenanceCase(family, universe, regime, int64(seed))
					if err != nil {
						return nil, nil, 0, err
					}
					want, err := evidence.ReferenceSupport(c.Root, c.Sources, c.Inputs)
					if err != nil {
						return nil, nil, 0, err
					}
					old, err := c.Program.Evaluate(c.Inputs)
					if err != nil || !sameResult(old, want) {
						return nil, nil, 0, fmt.Errorf("old executor/oracle mismatch: %v", err)
					}
					bases := make([]*evidence.ProvenanceSnapshot, len(formats))
					for repeat := 0; repeat < 8; repeat++ {
						order := rand.New(rand.NewSource(int64(seed*100 + repeat))).Perm(len(formats))
						for _, i := range order {
							format := formats[i]
							got, err := measure(func() (*evidence.ProvenanceSnapshot, error) {
								return c.Program.BuildProvenance(c.Inputs, format, managedBudget)
							})
							if err != nil || !sameResult(got.result, want) {
								return nil, nil, 0, fmt.Errorf("build mismatch %s/%d/%s/%s: %v", family, universe, regime, format, err)
							}
							bases[i] = got.snapshot
							builds = append(builds, buildRow{family, regime, universe, seed, repeat, format,
								got.ns, got.bytes, got.allocs, got.snapshot.Stats()})
						}
					}
					plan, inactive, err := weightPlan(c, want)
					if err != nil {
						return nil, nil, 0, err
					}
					if !inactive {
						skippedInactive++
					}
					for repeat := 0; repeat < 8; repeat++ {
						current := append([]*evidence.ProvenanceSnapshot(nil), bases...)
						for step, event := range plan {
							order := rand.New(rand.NewSource(int64(seed*1000 + repeat*10 + step))).Perm(len(formats))
							for _, i := range order {
								got, err := measure(func() (*evidence.ProvenanceSnapshot, error) {
									return current[i].WithSourceWeight(event.source.ID, event.source.Version, event.weight)
								})
								if err != nil || !sameResult(got.result, event.want) {
									return nil, nil, 0, fmt.Errorf("weight update mismatch %s/%d/%s: %v", family, universe, event.kind, err)
								}
								current[i] = got.snapshot
								updates = append(updates, weightRow{family, regime, event.kind, universe, seed, repeat,
									formats[i], got.ns, got.bytes, got.allocs, event.active, got.result.Weight})
							}
						}
					}
				}
			}
		}
	}
	return builds, updates, skippedInactive, nil
}

func quantile(values []uint64, percent int) uint64 {
	copyValues := append([]uint64(nil), values...)
	sort.Slice(copyValues, func(i, j int) bool { return copyValues[i] < copyValues[j] })
	return copyValues[(len(copyValues)*percent-1)/100]
}

type buildSummary struct {
	Family, Regime                    string
	Universe, Samples                 int
	Format                            evidence.ProvenanceFormat
	MedianNS, P95NS, MedianAllocBytes uint64
	MedianAllocs                      uint64
	Stats                             evidence.ProvenanceStats
}

type weightSummary struct {
	Family, Regime, Kind       string
	Universe, Samples          int
	Format                     evidence.ProvenanceFormat
	MedianNS, MedianAllocBytes uint64
}

type report struct {
	BuildMeasurements, WeightMeasurements int
	CasesWithoutInactiveSource            int
	Builds                                []buildSummary
	Weights                               []weightSummary
	SparseMemoryCriterion                 bool
	InternMemoryCriterion                 bool
}

func summarize(builds []buildRow, updates []weightRow, skipped int) (report, error) {
	out := report{BuildMeasurements: len(builds), WeightMeasurements: len(updates), CasesWithoutInactiveSource: skipped,
		SparseMemoryCriterion: true, InternMemoryCriterion: true}
	for _, family := range evidence.ProvenanceFamilies {
		for _, universe := range []int{256, 1024, 4096} {
			for _, regime := range []string{"resolved", "uncertain"} {
				budgets := map[evidence.ProvenanceFormat]int64{}
				for _, format := range formats {
					summary := buildSummary{Family: family, Regime: regime, Universe: universe, Format: format}
					times, allocations, counts := []uint64{}, []uint64{}, []uint64{}
					for _, row := range builds {
						if row.Family != family || row.Universe != universe || row.Regime != regime || row.Format != format {
							continue
						}
						if summary.Samples > 0 && summary.Stats != row.Stats {
							return report{}, fmt.Errorf("provenance shape differs across the declared permutation-only replicates")
						}
						times, allocations, counts = append(times, row.NS), append(allocations, row.AllocBytes), append(counts, row.Allocs)
						summary.Stats = row.Stats
						summary.Samples++
					}
					if len(times) == 0 {
						return report{}, fmt.Errorf("missing provenance measurement group")
					}
					summary.MedianNS, summary.P95NS = quantile(times, 50), quantile(times, 95)
					summary.MedianAllocBytes, summary.MedianAllocs = quantile(allocations, 50), quantile(counts, 50)
					out.Builds = append(out.Builds, summary)
					budgets[format] = summary.Stats.ManagedWorkingBytes
				}
				if universe == 4096 && family == "local-select" &&
					float64(budgets[evidence.ProvenanceSparse]) > .8*float64(budgets[evidence.ProvenanceDense]) {
					out.SparseMemoryCriterion = false
				}
				if universe == 4096 && family == "repeated-support" &&
					float64(budgets[evidence.ProvenanceInterned]) > .8*float64(budgets[evidence.ProvenanceSparse]) {
					out.InternMemoryCriterion = false
				}
				for _, kind := range []string{"active-weight", "inactive-weight", "active-zero", "active-restore"} {
					for _, format := range formats {
						times, allocations := []uint64{}, []uint64{}
						for _, row := range updates {
							if row.Family == family && row.Universe == universe && row.Regime == regime && row.Format == format && row.Kind == kind {
								times, allocations = append(times, row.NS), append(allocations, row.AllocBytes)
							}
						}
						if len(times) > 0 {
							out.Weights = append(out.Weights, weightSummary{family, regime, kind, universe, len(times), format,
								quantile(times, 50), quantile(allocations, 50)})
						}
					}
				}
			}
		}
	}
	return out, nil
}

func writeCSV(dir, name string, header []string, rows [][]string) error {
	if dir == "" {
		return nil
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	file, err := os.Create(filepath.Join(dir, name))
	if err != nil {
		return err
	}
	w := csv.NewWriter(file)
	if err := w.Write(header); err != nil {
		file.Close()
		return err
	}
	for _, row := range rows {
		if err := w.Write(row); err != nil {
			file.Close()
			return err
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		file.Close()
		return err
	}
	return file.Close()
}

func fields(values ...any) []string {
	out := make([]string, len(values))
	for i, value := range values {
		out[i] = fmt.Sprint(value)
	}
	return out
}

func writeEvidence(dir string, builds []buildRow, updates []weightRow) error {
	rows := [][]string{}
	for _, row := range builds {
		s := row.Stats
		rows = append(rows, fields(row.Family, row.Universe, row.Regime, row.Seed, row.Repeat, row.Format,
			row.NS, row.AllocBytes, row.Allocs, s.Nodes, s.Sources, s.RootSources, s.ManagedWorkingBytes, s.RetainedBytes,
			s.VectorIDsAllocated, s.MaskWords, s.InternedVectors, s.IndexSlots))
	}
	if err := writeCSV(dir, "build.csv", []string{"family", "sources", "regime", "seed", "repeat", "format",
		"total_ns", "allocated_bytes", "allocations", "nodes", "source_universe", "root_sources", "managed_working_bytes",
		"root_retained_bytes", "uint32_ids_allocated", "mask_words", "interned_vectors", "index_slots"}, rows); err != nil {
		return err
	}
	rows = nil
	for _, row := range updates {
		rows = append(rows, fields(row.Family, row.Universe, row.Regime, row.Seed, row.Repeat, row.Format,
			row.Kind, row.Active, row.NS, row.AllocBytes, row.Allocs, strconv.FormatFloat(row.Weight, 'g', 17, 64)))
	}
	return writeCSV(dir, "weights.csv", []string{"family", "sources", "regime", "seed", "repeat", "format",
		"kind", "changed_source_active", "total_ns", "allocated_bytes", "allocations", "result_weight"}, rows)
}

func run(dir string, output io.Writer) error {
	builds, updates, skipped, err := runMeasurements()
	if err != nil {
		return err
	}
	summary, err := summarize(builds, updates, skipped)
	if err != nil {
		return err
	}
	if err := writeEvidence(dir, builds, updates); err != nil {
		return err
	}
	encoder := json.NewEncoder(output)
	encoder.SetIndent("", "  ")
	return encoder.Encode(summary)
}

func main() {
	dir := flag.String("csv-dir", "", "explicit directory for two new provenance measurement tables; existing files are replaced")
	flag.Parse()
	if err := run(*dir, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
