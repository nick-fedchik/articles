package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	evidence "github.com/nick-fedchik/articles/Expert_System_Circuits/src"
)

func TestResultComparisonRejectsDuplicateSourcesAndNaN(t *testing.T) {
	a, err := evidence.NewSource("a", "1", []byte("public-a"), .9)
	if err != nil {
		t.Fatal(err)
	}
	b, err := evidence.NewSource("b", "1", []byte("public-b"), .9)
	if err != nil {
		t.Fatal(err)
	}
	want := evidence.Result{State: evidence.True, Weight: .8, Sources: []evidence.Source{a, b}}
	got := want
	got.Sources = []evidence.Source{a, a}
	if sameResult(got, want) {
		t.Fatal("duplicate source concealed an omitted source")
	}
	got = want
	got.Weight = math.NaN()
	if sameResult(got, want) {
		t.Fatal("NaN was treated as equivalent")
	}
}

func TestPublishedProvenanceEvidence(t *testing.T) {
	expected := map[string]struct {
		rows   int
		digest string
	}{
		"build.csv":   {2304, "ff035a15e85f13cb3698ce60d167afbb084b028a6866a22672c5fb28426e15f0"},
		"weights.csv": {7488, "852c5e4c65a42048046293c66f6e1c4d0891445cf29455fa6de1cf5bf2dbd3af"},
	}
	for name, want := range expected {
		data, err := os.ReadFile(filepath.Join("..", "..", "testdata", "provenance-study", name))
		if err != nil {
			t.Fatal(err)
		}
		if fmt.Sprintf("%x", sha256.Sum256(data)) != want.digest {
			t.Fatalf("%s changed without a corresponding provenance protocol revision", name)
		}
		rows, err := csv.NewReader(bytes.NewReader(data)).ReadAll()
		if err != nil || len(rows) != want.rows+1 {
			t.Fatal("incorrect provenance measurement count", err)
		}
		columns := map[string]int{}
		for i, column := range rows[0] {
			columns[column] = i
		}
		if name == "build.csv" {
			shapes := map[string][]string{}
			for _, row := range rows[1:] {
				key := row[columns["family"]] + "/" + row[columns["sources"]] + "/" + row[columns["regime"]] + "/" + row[columns["format"]]
				shape := []string{row[columns["nodes"]], row[columns["root_sources"]], row[columns["managed_working_bytes"]],
					row[columns["root_retained_bytes"]], row[columns["uint32_ids_allocated"]], row[columns["index_slots"]]}
				if previous, ok := shapes[key]; ok {
					for i := range shape {
						if shape[i] != previous[i] {
							t.Fatal("permutation-only repeats changed the reported shape")
						}
					}
				}
				shapes[key] = shape
				working, err := strconv.ParseInt(row[columns["managed_working_bytes"]], 10, 64)
				if err != nil {
					t.Fatal(err)
				}
				retained, err := strconv.ParseInt(row[columns["root_retained_bytes"]], 10, 64)
				if err != nil || retained >= working || working > managedBudget {
					t.Fatal("invalid root-retention or managed-budget claim", err)
				}
			}
			if len(shapes) != 96 {
				t.Fatal("missing family/source/regime/format group")
			}
		} else {
			for _, row := range rows[1:] {
				if row[columns["kind"]] == "active-zero" && row[columns["result_weight"]] != "0" {
					t.Fatal("active zero-weight source did not set the model rating to zero")
				}
			}
		}
	}
}
