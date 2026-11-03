package main

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"testing"
)

func TestCSVRecordsDeclaredBudgetAndUpdateDenominator(t *testing.T) {
	dir := t.TempDir()
	searches := []searchSample{{Family: "public", Method: "fifo", Rules: 2, Seed: 1, Multiplier: 4, NS: 100}}
	updates := []updateSample{{Family: "public", Rules: 2, Percent: 1, Changed: 1, Roots: 32, Atoms: 34}}
	if err := writeCSV(dir, searches, updates); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"search.csv", "updates.csv"} {
		file, err := os.Open(filepath.Join(dir, name))
		if err != nil {
			t.Fatal(err)
		}
		rows, err := csv.NewReader(file).ReadAll()
		if closeErr := file.Close(); closeErr != nil {
			t.Fatal(closeErr)
		}
		if err != nil || len(rows) != 2 || len(rows[0]) != len(rows[1]) {
			t.Fatal("inconsistent evidence table", err)
		}
	}
}
