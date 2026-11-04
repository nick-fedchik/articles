// Command baseline runs the preregistered H3/H4 baseline comparison and
// fault-injection study and writes deterministic results:
//
//	baseline -out testdata/baseline            # results.json, results.csv, scenarios.csv
//	baseline -out testdata/baseline -timing    # additionally timing.csv (host microbenchmark)
package main

import (
	"crypto/sha256"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/nick-fedchik/articles/Expert_System_Circuits/src/baseline"
)

type classInfo struct {
	Name, UCA, Description string
	InModel                bool
}

type output struct {
	Protocol          string
	Seed              uint64
	Faults            int
	Controls          int
	Arms              map[string]string
	SpecLines         map[string]int
	SharedEngineLines map[string]int
	Classes           []classInfo
	Rows              []baseline.Row
	Summary           baseline.Summary
	Notes             []string
}

var armDescriptions = map[string]string{
	"A":   "naive: noisy-OR scalar confidences (trusted 0.95, other 0.6, duplicates counted), point midpoint of latest measurement read as bar, flat deny-first chain, no freshness or provenance checks",
	"B0":  "conventional stack, conventional information: typed units with conversion, point value +/- fixed 0.3 bar instrument margin, one boolean per fact (latest report wins), Horn closure, ordered first-match rules, two-sided freshness, hash-chained audit log not consulted",
	"B1":  "conventional stack, same metadata as C: typed units with conversion, interval check, pro/contra Horn closure, source sets checked against trusted list, two-sided freshness of every item, fail-safe priority-ordered rules with exceptions, audit log",
	"C":   "kernel.Evaluate over CaseSpec (DEFER_* mapped to DEFER)",
	"B1S": "EXPLORATORY, not preregistered: B1 without unit conversion (non-bar unit -> DEFER)",
}

func write(path string, data []byte) error {
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "%x  %s\n", sha256.Sum256(data), path)
	return nil
}

func writeCSV(path string, rows [][]string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	w := csv.NewWriter(f)
	if err := w.WriteAll(rows); err != nil {
		f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "%x  %s\n", sha256.Sum256(data), path)
	return nil
}

func itoa(v int) string { return strconv.Itoa(v) }

func ftoa(v float64) string { return strconv.FormatFloat(v, 'g', 8, 64) }

func run() error {
	seed := flag.Uint64("seed", 20261103, "PCG seed")
	faults := flag.Int("faults", 1000, "faulted scenarios per class")
	controls := flag.Int("controls", 250, "matched clean controls per fault class")
	dir := flag.String("out", "testdata/baseline", "output directory")
	timing := flag.Bool("timing", false, "also write timing.csv (non-deterministic)")
	perClass := flag.Int("timing-scenarios", 200, "scenarios per class for timing")
	reps := flag.Int("timing-reps", 20, "repetitions per scenario for timing")
	flag.Parse()

	scenarios := baseline.Generate(*seed, *faults, *controls)
	records, err := baseline.Evaluate(scenarios)
	if err != nil {
		return err
	}
	rows := baseline.Aggregate(records)
	out := output{Protocol: "PREREGISTRATION-2026-11-03.md H3/H4", Seed: *seed, Faults: *faults, Controls: *controls,
		Arms: armDescriptions, SpecLines: map[string]int{}, Rows: rows, Summary: baseline.Summarize(rows),
		Notes: []string{
			"Ground truth P = kernel.CheckSafety; arm C is correct w.r.t. P by construction (ALLOW<=>P).",
			"Physical ground truth = World.Safe on the hidden synthetic world; in-model classes keep trusted sources honest and honest intervals containing the truth.",
			"F-01/F-02: kernel scopes freshness/provenance to rule-referenced atoms (tested in F11_UNREFERENCED_SPAM) and requires clean non-conflicting derivation support Clean (Theorem B.7).",
			"F-06: SpecLines counts per-case specification lines between markers; SharedEngineLines reports the shared kernel TCB and conventional support lines.",
			"F-08: H4 stratifies total-refusal fault classes (Allow == 0, conditional risk per ALLOW Undefined 0/0, bounded per scenario n=1000) from active-permit fault classes (Allow > 0: F07_NORM_PRIORITY_CONFLICT, F11_UNREFERENCED_SPAM).",
			"Scenarios within a class are generated from one distribution and are not independent real-world trials; bounds describe the synthetic generator only.",
			"STPA-lite, classes and arms were written with AI assistance and without a domain expert; this is not industry validation.",
		}}
	for _, arm := range []string{"A", "B0", "B1", "C"} {
		if out.SpecLines[arm], err = baseline.SpecLines(arm); err != nil {
			return err
		}
	}
	if out.SharedEngineLines, err = baseline.SharedEngineLines(); err != nil {
		return err
	}
	for _, c := range baseline.Classes() {
		out.Classes = append(out.Classes, classInfo{c.Name, c.UCA, c.Description, c.InModel})
	}
	if err := os.MkdirAll(*dir, 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return err
	}
	if err := write(filepath.Join(*dir, "results.json"), append(data, '\n')); err != nil {
		return err
	}
	table := [][]string{{"class", "group", "arm", "n", "allow", "deny", "defer", "p_safe", "phys_safe",
		"unsafe_permit_p", "unsafe_permit_phys", "false_refusal_p", "false_refusal_phys", "automatic",
		"audit_median_bytes", "audit_max_bytes", "cp95_upper_unsafe_phys_per_allow", "cp95_upper_unsafe_phys_per_scenario"}}
	for _, r := range rows {
		perAllow := ""
		if r.UnsafePhysPerAllow.Defined {
			perAllow = ftoa(r.UnsafePhysPerAllow.CPUpper95TwoSided)
		}
		table = append(table, []string{r.Class, r.Group, r.Arm, itoa(r.N), itoa(r.Allow), itoa(r.Deny), itoa(r.Defer),
			itoa(r.PSafe), itoa(r.PhysSafe), itoa(r.UnsafeP), itoa(r.UnsafePhys), itoa(r.FalseRefusalP),
			itoa(r.FalseRefusalPhys), itoa(r.Automatic), itoa(r.AuditMedianBytes), itoa(r.AuditMaxBytes),
			perAllow, ftoa(r.UnsafePhysPerScenario.CPUpper95TwoSided)})
	}
	if err := writeCSV(filepath.Join(*dir, "results.csv"), table); err != nil {
		return err
	}
	header := []string{"class", "index", "group", "variant", "phys_safe", "p_holds", "violations"}
	for _, arm := range baseline.Arms {
		header = append(header, "decision_"+arm, "raw_"+arm)
	}
	perScenario := [][]string{header}
	for _, r := range records {
		line := []string{r.Class, itoa(r.Index), r.Group(), r.Variant, strconv.FormatBool(r.PhysSafe),
			strconv.FormatBool(r.PHolds), r.Violations}
		for i := range baseline.Arms {
			line = append(line, string(r.Decision[i]), r.Raw[i])
		}
		perScenario = append(perScenario, line)
	}
	if err := writeCSV(filepath.Join(*dir, "scenarios.csv"), perScenario); err != nil {
		return err
	}
	if !*timing {
		return nil
	}
	timings, err := baseline.Timing(scenarios, *perClass, *reps)
	if err != nil {
		return err
	}
	timingTable := [][]string{{"class", "arm", "scenarios", "reps", "median_ns_per_eval", "label"}}
	for _, t := range timings {
		timingTable = append(timingTable, []string{t.Class, t.Arm, itoa(t.Scenarios), itoa(t.Reps),
			strconv.FormatInt(t.MedianNs, 10), "host microbenchmark, non-deterministic"})
	}
	return writeCSV(filepath.Join(*dir, "timing.csv"), timingTable)
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
