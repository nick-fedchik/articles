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

func TestPublishedAdaptiveAndCacheEvidence(t *testing.T) {
	expected := map[string]struct {
		rows   int
		digest string
	}{
		"adaptive.csv":      {7200, "59843fab25643567d24a4ecd38e67072ecd55aeb443655e0a2b50ab03b031e13"},
		"cache-events.csv":  {3200, "9c0caae5c28b5b0dc32c09838234b5cab6a61ca8373feb9cf02990ef5fda3729"},
		"cache-queries.csv": {25600, "f9141823d2eaff5de68aaae8686ed9618a636adda8f5d1351218744d5aeb899d"},
	}
	for name, want := range expected {
		data, err := os.ReadFile(filepath.Join("..", "..", "testdata", "adaptive-study", name))
		if err != nil {
			t.Fatal(err)
		}
		if fmt.Sprintf("%x", sha256.Sum256(data)) != want.digest {
			t.Fatalf("%s changed without a corresponding protocol revision", name)
		}
		rows, err := csv.NewReader(bytes.NewReader(data)).ReadAll()
		if err != nil || len(rows) != want.rows+1 {
			t.Fatal("invalid primary measurement count", err)
		}
		columns := map[string]int{}
		for i, column := range rows[0] {
			columns[column] = i
		}
		if name == "adaptive.csv" {
			totals := map[string]int64{}
			for _, row := range rows[1:] {
				seed, err := strconv.Atoi(row[columns["seed"]])
				if err != nil {
					t.Fatal(err)
				}
				split := row[columns["split"]]
				if split == "train" && (seed < 1000 || seed > 1003) || split == "test" && (seed < 2000 || seed > 2007) {
					t.Fatal("training/test seed leakage")
				}
				if split == "train" && row[columns["mode"]] == "adaptive" {
					ns, err := strconv.ParseInt(row[columns["total_ns"]], 10, 64)
					if err != nil {
						t.Fatal(err)
					}
					totals[row[columns["discovery_fraction"]]] += ns
				}
				if split == "test" && row[columns["mode"]] == "adaptive" && row[columns["discovery_fraction"]] != "0.05" {
					t.Fatal("control rows do not use the frozen selected policy")
				}
			}
			for fraction, total := range totals {
				if total < totals["0.05"] {
					t.Fatalf("published policy was not the training-only minimum: %s", fraction)
				}
			}
		}
		if name == "cache-queries.csv" {
			hits, invalidations := 0, 0
			for _, row := range rows[1:] {
				if row[columns["expected"]] != row[columns["proved"]] {
					t.Fatal("published cache result disagrees with current closure")
				}
				if row[columns["reused"]] == "true" {
					hits++
					if row[columns["proved"]] != "true" || row[columns["invalidation"]] != "" {
						t.Fatal("invalid cached acceptance")
					}
				}
				if row[columns["invalidation"]] != "" {
					invalidations++
				}
			}
			if hits != 15099 || invalidations != 1654 {
				t.Fatal("published cache counts do not match the protocol")
			}
		}
	}
}
