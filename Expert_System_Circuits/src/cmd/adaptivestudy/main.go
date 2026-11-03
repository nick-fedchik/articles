package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
)

func percentile(values []int64, percent int) int64 {
	copyValues := append([]int64(nil), values...)
	sort.Slice(copyValues, func(i, j int) bool { return copyValues[i] < copyValues[j] })
	return copyValues[(len(copyValues)*percent-1)/100]
}

func writeRows(dir, name string, header []string, rows [][]string) error {
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
	writer := csv.NewWriter(file)
	if err := writer.Write(header); err != nil {
		file.Close()
		return err
	}
	for _, row := range rows {
		if err := writer.Write(row); err != nil {
			file.Close()
			return err
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		file.Close()
		return err
	}
	return file.Close()
}

func values(fields ...any) []string {
	out := make([]string, len(fields))
	for i, field := range fields {
		out[i] = fmt.Sprint(field)
	}
	return out
}

func writeEvidence(dir string, adaptive []adaptiveRow, events []cacheEventRow, queries []cacheQueryRow) error {
	rows := [][]string{}
	for _, row := range adaptive {
		rows = append(rows, values(row.Split, row.Family, row.Rules, row.Seed, row.Step, row.Kind, row.Percent,
			row.Changed, row.Total, row.Mode, row.Fraction, row.Path, row.NS, row.Discovery, row.Recomputed))
	}
	if err := writeRows(dir, "adaptive.csv", []string{"split", "family", "rules", "seed", "step", "kind", "target_percent",
		"changed", "total_inputs", "mode", "discovery_fraction", "actual_path", "total_ns", "discovery_visits", "recomputed_atoms"}, rows); err != nil {
		return err
	}
	rows = nil
	for _, row := range events {
		rows = append(rows, values(row.Family, row.Rules, row.Seed, row.Step, row.Kind, row.NS))
	}
	if err := writeRows(dir, "cache-events.csv", []string{"family", "rules", "seed", "step", "kind", "event_ns"}, rows); err != nil {
		return err
	}
	rows = nil
	for _, row := range queries {
		rows = append(rows, values(row.Family, row.Rules, row.Seed, row.Step, row.Query, row.Expected, row.Proved,
			row.Reused, row.Invalidation, row.ColdNS, row.CachedNS))
	}
	return writeRows(dir, "cache-queries.csv", []string{"family", "rules", "seed", "step", "query", "expected", "proved",
		"reused", "invalidation", "cold_ns", "cached_ns"}, rows)
}

func run(dir string, output io.Writer) error {
	adaptive, adaptiveResult, err := adaptiveExperiment()
	if err != nil {
		return err
	}
	events, queries, cacheResult, err := cacheExperiment()
	if err != nil {
		return err
	}
	if err := writeEvidence(dir, adaptive, events, queries); err != nil {
		return err
	}
	encoder := json.NewEncoder(output)
	encoder.SetIndent("", "  ")
	return encoder.Encode(struct {
		Adaptive adaptiveReport
		Cache    cacheReport
	}{adaptiveResult, cacheResult})
}

func main() {
	dir := flag.String("csv-dir", "", "explicit directory for new primary CSV measurements; three files are replaced")
	flag.Parse()
	if err := run(*dir, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
