package main

import (
	"crypto/sha256"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	evidence "github.com/nick-fedchik/articles/Expert_System_Circuits/src"
)

func load(path string) ([]evidence.CorpusItem, string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, "", err
	}
	return decode(data)
}

func decode(data []byte) ([]evidence.CorpusItem, string, error) {
	var document struct {
		Clauses []struct {
			Source string `json:"source_ref"`
			Text   string `json:"text"`
			Label  struct {
				Modality string `json:"modality"`
			} `json:"arbitration"`
		} `json:"clauses"`
		Items []struct {
			Source, Text, Category string
		} `json:"items"`
	}
	if err := json.Unmarshal(data, &document); err != nil {
		return nil, "", err
	}
	if len(document.Clauses) == 0 && len(document.Items) == 0 {
		return nil, "", fmt.Errorf("corpus has neither clauses nor items")
	}
	out := []evidence.CorpusItem{}
	for _, row := range document.Clauses {
		out = append(out, evidence.CorpusItem{Document: row.Source, Text: row.Text, Label: row.Label.Modality})
	}
	for _, row := range document.Items {
		out = append(out, evidence.CorpusItem{Document: row.Source, Text: row.Text, Label: row.Category})
	}
	return out, fmt.Sprintf("%x", sha256.Sum256(data)), nil
}

func run(paths []string, target float64) error {
	if len(paths) == 0 {
		return fmt.Errorf("explicit local corpus files are required; no private paths are assumed")
	}
	items, digests := []evidence.CorpusItem{}, []string{}
	for _, path := range paths {
		rows, digest, err := load(path)
		if err != nil {
			return fmt.Errorf("corpus input: %w", err)
		}
		items, digests = append(items, rows...), append(digests, digest)
	}
	report, err := evidence.EvaluateCorpus(items, target)
	if err != nil {
		return err
	}
	out := struct {
		InputSHA256 []string
		Report      evidence.CorpusReport
	}{digests, report}
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(out)
}

func main() {
	target := flag.Float64("target-risk", .1, "calibration Wilson upper bound used to select the operating point")
	flag.Parse()
	if err := run(flag.Args(), *target); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
