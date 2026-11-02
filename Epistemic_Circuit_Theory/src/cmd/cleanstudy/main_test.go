package main

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestExhaustiveSemanticStudy(t *testing.T) {
	s := newStudy()
	if err := s.semantics(); err != nil {
		t.Fatal(err)
	}
	if s.report.Cases != 38912 || s.report.Clean != 9014 || s.report.ConflictDependent != 12682 ||
		s.report.Unsupported != 17216 || s.report.LeavesOnlyDisagreements != 16 || s.report.AllProofsOverblocks != 500 {
		t.Fatalf("semantic counts changed: %+v", s.report)
	}
}

func TestExclusiveCSVOutput(t *testing.T) {
	path := filepath.Join(t.TempDir(), "rows.csv")
	if err := writeCSV(path, []string{"a", "b"}, [][]string{{"one", "two"}}); err != nil {
		t.Fatal(err)
	}
	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := writeCSV(path, nil, nil); err == nil {
		t.Fatal("existing research data overwritten")
	}
	after, err := os.ReadFile(path)
	if err != nil || string(before) != string(after) {
		t.Fatal("published data changed")
	}
}

func TestPublishedCleanStudyHashes(t *testing.T) {
	for name, want := range map[string]string{
		"semantics.csv": "c03426dea848b84a27debf4983d6f1546ddb1cec1d8f7e87c8bd97e26b3a0498",
		"scaling.csv":   "e0ab7cfff1fd451856678484c6f6a04e82e40f7f9f03b9e67c530d3fa770d208",
	} {
		data, err := os.ReadFile(filepath.Join("..", "..", "testdata", "clean-study", name))
		if err != nil {
			t.Fatal(err)
		}
		if got := fmt.Sprintf("%x", sha256.Sum256(data)); got != want {
			t.Fatalf("%s hash %s != %s", name, got, want)
		}
	}
}
