package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nick-fedchik/articles/Epistemic_Circuit_Theory/src/annotation"
)

func TestCLIEndToEndAndExclusiveOutputs(t *testing.T) {
	demo, err := annotation.SyntheticDemo()
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	path := func(name string) string { return filepath.Join(dir, name) }
	for _, file := range []struct {
		name  string
		value any
	}{{"corpus.json", demo.Corpus}, {"a.json", demo.A}, {"b.json", demo.B}, {"arbiter.json", demo.Adjudication}} {
		if err := writeNew(path(file.name), file.value); err != nil {
			t.Fatal(err)
		}
	}
	check := func(args ...string) []byte {
		t.Helper()
		var buffer bytes.Buffer
		if err := run(args, &buffer); err != nil {
			t.Fatalf("%v: %v", args, err)
		}
		return buffer.Bytes()
	}
	check("prepare", "-corpus", path("corpus.json"), "-out", path("packets"))
	check("validate", "-corpus", path("corpus.json"), "-annotation", path("a.json"))
	var agreement annotation.Agreement
	if err := json.Unmarshal(check("compare", "-corpus", path("corpus.json"), "-a", path("a.json"), "-b", path("b.json")), &agreement); err != nil {
		t.Fatal(err)
	}
	if agreement.Exact.Matched != 14 || agreement.IndependenceAttested {
		t.Fatalf("wrong CLI agreement: %+v", agreement)
	}
	check("adjudicate", "-corpus", path("corpus.json"), "-a", path("a.json"), "-b", path("b.json"),
		"-arbiter", path("arbiter.json"), "-out", path("gold.json"))
	var gold annotation.Gold
	var plan annotation.SplitPlan
	if err := readJSON(path("gold.json"), &gold); err != nil {
		t.Fatal(err)
	}
	if err := readJSON(path("packets/split.json"), &plan); err != nil {
		t.Fatal(err)
	}
	documents, labels := map[string]string{}, map[string]annotation.Label{}
	for _, entry := range plan.Assignments {
		documents[entry.DocumentID] = entry.Split
	}
	for _, record := range gold.Records {
		labels[record.ItemID] = record.Label
	}
	predictions := annotation.Predictions{Schema: annotation.Schema, CorpusDigest: plan.CorpusDigest, System: "synthetic-only"}
	for _, item := range demo.Corpus.Items {
		if documents[item.DocumentID] != "test" {
			continue
		}
		label := labels[item.ID]
		if item.Task == "decision" {
			label.Decision = "ALLOW"
		}
		predictions.Records = append(predictions.Records, annotation.Prediction{ItemID: item.ID, Label: &label})
	}
	if err := writeNew(path("predictions.json"), predictions); err != nil {
		t.Fatal(err)
	}
	var report annotation.Evaluation
	if err := json.Unmarshal(check("evaluate", "-corpus", path("corpus.json"), "-plan", path("packets/split.json"),
		"-gold", path("gold.json"), "-predictions", path("predictions.json")), &report); err != nil {
		t.Fatal(err)
	}
	if !report.Synthetic || report.IndependenceAttested || report.FalseAllow.Trials == 0 ||
		report.FalseAllow.Errors != report.FalseAllow.Trials {
		t.Fatalf("wrong evaluation: %+v", report)
	}
	before, err := os.ReadFile(path("gold.json"))
	if err != nil {
		t.Fatal(err)
	}
	if err := writeNew(path("gold.json"), map[string]bool{"overwritten": true}); err == nil {
		t.Fatal("existing human work could be overwritten")
	}
	after, err := os.ReadFile(path("gold.json"))
	if err != nil || !bytes.Equal(before, after) {
		t.Fatal("exclusive output changed")
	}
	if err := run([]string{"prepare", "-corpus", path("corpus.json"), "-out", path("packets")}, new(bytes.Buffer)); err == nil {
		t.Fatal("existing packets overwritten")
	}
}

func TestCLIRejectsMalformedInputsAndMissingArguments(t *testing.T) {
	for _, args := range [][]string{nil, {"unknown"}, {"validate"}, {"demo", "extra"}} {
		if err := run(args, new(bytes.Buffer)); err == nil {
			t.Fatalf("invalid command accepted: %v", args)
		}
	}
	for i, data := range []string{
		`{"Unknown":true}`,
		`{} {}`,
		`{"Schema":"first","Schema":"second"}`,
		`{"Human":false,"human":true}`,
	} {
		path := filepath.Join(t.TempDir(), "bad.json")
		if err := os.WriteFile(path, []byte(data), 0600); err != nil {
			t.Fatal(err)
		}
		var corpus annotation.Corpus
		if err := readJSON(path, &corpus); err == nil {
			t.Fatalf("malformed input %d accepted", i)
		}
	}
}

func TestDemoMakesNoHumanClaim(t *testing.T) {
	var out bytes.Buffer
	if err := run([]string{"demo"}, &out); err != nil {
		t.Fatal(err)
	}
	var result struct {
		Scope string
		Gold  annotation.Gold
	}
	if err := json.Unmarshal(out.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(result.Scope, "synthetic") || !result.Gold.Synthetic || result.Gold.IndependenceAttested {
		t.Fatal("synthetic run presented as independent human evidence")
	}
}
