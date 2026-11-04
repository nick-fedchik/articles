package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestPublishedNormativeEvidence(t *testing.T) {
	expected := map[string]struct {
		rows   int
		digest string
	}{
		"semantics.csv":     {1448, "d36d6a4092534439b874a359a357291fc193fa0c9e8624d32d91d185497ae4d6"},
		"context-cache.csv": {2048, "0e3ef33bff6ff153664034b213fbe54743d3b97659e67834efbc15c9e645d602"},
	}
	for name, want := range expected {
		data, err := os.ReadFile(filepath.Join("..", "..", "testdata", "normative-study", name))
		if err != nil {
			t.Fatal(err)
		}
		if fmt.Sprintf("%x", sha256.Sum256(data)) != want.digest {
			t.Fatalf("%s changed without a corresponding normative protocol revision", name)
		}
		rows, err := csv.NewReader(bytes.NewReader(data)).ReadAll()
		if err != nil || len(rows) != want.rows+1 {
			t.Fatal("incorrect normative evidence count", err)
		}
		columns := map[string]int{}
		for i, key := range rows[0] {
			columns[key] = i
		}
		if name == "semantics.csv" {
			sat, priorityResolutions, stages := 0, 0, map[string]int{}
			for _, row := range rows[1:] {
				stages[row[columns["stage"]]]++
				if row[columns["sat_checked"]] == "true" {
					sat++
				}
				if row[columns["decision"]] == "ALLOW" && row[columns["applicable_evidence"]] == "3" &&
					row[columns["surviving_evidence"]] == "2" {
					priorityResolutions++
				}
			}
			if sat != 456 || priorityResolutions != 11 || stages["signed"] != 256 || stages["exceptions"] != 192 || stages["generated"] != 1000 {
				t.Fatal("published stages do not match the protocol")
			}
		} else {
			unsafe := 0
			for _, row := range rows[1:] {
				if row[columns["actual"]] != row[columns["expected"]] {
					t.Fatal("context cache returned an incorrect current decision")
				}
				if row[columns["repeat"]] == "0" && row[columns["unsafe_positive_only"]] == "ALLOW" && row[columns["expected"]] != "ALLOW" {
					unsafe++
				}
				if row[columns["repeat"]] == "1" && row[columns["reused"]] != "true" {
					t.Fatal("identical context failed to reuse its certificate")
				}
			}
			if unsafe != 768 {
				t.Fatal("naive-cache counterexample count changed")
			}
		}
	}
}
