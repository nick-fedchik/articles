package main

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestExperimentReportHasAllSixStages(t *testing.T) {
	var output bytes.Buffer
	if err := run(&output); err != nil {
		t.Fatal(err)
	}
	var report map[string]json.RawMessage
	if err := json.Unmarshal(output.Bytes(), &report); err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"01_semantics", "02_provenance", "03_runtime", "04_repairs", "05_zero_ripple_plant", "06_routing"} {
		if len(report[key]) == 0 {
			t.Errorf("missing stage %q", key)
		}
	}
}
