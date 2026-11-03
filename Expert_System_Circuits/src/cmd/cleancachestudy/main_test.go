package main

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	e "github.com/nick-fedchik/articles/Expert_System_Circuits/src"
)

func TestCorrectnessStudy(t *testing.T) {
	s := &study{report: report{Reused: map[e.CleanCacheMode]int{}}}
	if err := s.correctness(); err != nil {
		t.Fatal(err)
	}
	if s.report.Events != 3072 || s.report.Queries != 9216 || s.report.CheckedAnswers != 27648 ||
		s.report.NaiveWrongClean != 930 || s.report.Reused[e.CleanCone] != 364 {
		t.Fatalf("changed deterministic counts: %+v", s.report)
	}
}

func TestPublishedHashesAndExclusiveOutput(t *testing.T) {
	for name, want := range map[string]string{
		"events.csv":  "320a8e6b08c097a9b24d297fe93982c39939c488e418df49e557a43e23be9116",
		"timings.csv": "aebafc3cc204f48f8f7a2940fc4064c01d5284218220074248967df5b018c103",
	} {
		data, err := os.ReadFile(filepath.Join("..", "..", "testdata", "clean-cache-study", name))
		if err != nil {
			t.Fatal(err)
		}
		if got := fmt.Sprintf("%x", sha256.Sum256(data)); got != want {
			t.Fatalf("%s hash %s", name, got)
		}
	}
	path := filepath.Join(t.TempDir(), "result.csv")
	if err := write(path, []string{"header"}, [][]string{{"row"}}); err != nil {
		t.Fatal(err)
	}
	if err := write(path, nil, nil); err == nil {
		t.Fatal("existing data overwritten")
	}
}
