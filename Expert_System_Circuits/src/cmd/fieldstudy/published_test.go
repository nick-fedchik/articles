package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

func TestPublishedEvidenceMatchesProtocol(t *testing.T) {
	expected := map[string]struct {
		rows   int
		digest string
	}{
		"search.csv":  {1200, "f3cb3d81df4b20217e3dc1e7d5ea69824d6d663d74563b1d9e23fcf754543be7"},
		"updates.csv": {480, "75785cab49b0c2f0ea0482d43f3d737e8b16d9e61e7c3880e3398ec30e9bc522"},
	}
	for name, want := range expected {
		data, err := os.ReadFile(filepath.Join("..", "..", "testdata", "field-study", name))
		if err != nil {
			t.Fatal(err)
		}
		if fmt.Sprintf("%x", sha256.Sum256(data)) != want.digest {
			t.Fatalf("%s changed without a corresponding protocol revision", name)
		}
		rows, err := csv.NewReader(bytes.NewReader(data)).ReadAll()
		if err != nil || len(rows) != want.rows+1 {
			t.Fatal("incorrect evidence row count", err)
		}
		index := map[string]int{}
		for i, column := range rows[0] {
			index[column] = i
		}
		for _, row := range rows[1:] {
			if name == "search.csv" {
				if row[index["proved"]] == "true" && row[index["query_supported"]] != "true" {
					t.Fatal("published false acceptance")
				}
				if row[index["budget_multiplier"]] == "128" &&
					(row[index["proved"]] != row[index["query_supported"]] || row[index["exhausted"]] != "false") {
					t.Fatal("largest-budget completeness claim does not match the evidence")
				}
			} else {
				changed, err := strconv.Atoi(row[index["changed_roots"]])
				if err != nil {
					t.Fatal(err)
				}
				total, err := strconv.Atoi(row[index["total_used_roots"]])
				if err != nil || changed < 1 || changed > total {
					t.Fatal("invalid update denominator", err)
				}
			}
		}
	}
}
