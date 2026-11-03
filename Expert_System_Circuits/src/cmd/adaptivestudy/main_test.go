package main

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"testing"
)

func TestAdaptiveEvidenceTablesHaveExplicitPolicyAndCosts(t *testing.T) {
	dir := t.TempDir()
	rows := []adaptiveRow{{Split: "test", Family: "public", Mode: "adaptive", Fraction: .05, NS: 100}}
	events := []cacheEventRow{{Family: "public", NS: 10}}
	queries := []cacheQueryRow{{Family: "public", ColdNS: 100, CachedNS: 20}}
	if err := writeEvidence(dir, rows, events, queries); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"adaptive.csv", "cache-events.csv", "cache-queries.csv"} {
		file, err := os.Open(filepath.Join(dir, name))
		if err != nil {
			t.Fatal(err)
		}
		data, err := csv.NewReader(file).ReadAll()
		if closeErr := file.Close(); closeErr != nil {
			t.Fatal(closeErr)
		}
		if err != nil || len(data) != 2 || len(data[0]) != len(data[1]) {
			t.Fatal("malformed evidence table", err)
		}
	}
}
