package main

import (
	"bytes"
	"encoding/json"
	"io"
	"testing"

	"github.com/nick-fedchik/articles/Expert_System_Circuits/src/hardware"
)

func TestCommandJSON(t *testing.T) {
	var out bytes.Buffer
	if err := run(&out, io.Discard, []string{"-states=1", "-bus-bits=512", "-hz=100000000", "-leaves=16"}); err != nil {
		t.Fatal(err)
	}
	var report hardware.Report
	if err := json.Unmarshal(out.Bytes(), &report); err != nil {
		t.Fatal(err)
	}
	if report.Config.States != 1 || report.Layouts[0].AlignedTransfer.BeatsMin != 1 ||
		!report.Provenance.FitsRetentionCap || report.Status != "calculated_not_hardware_measurement" {
		t.Fatalf("incorrect command result: %+v", report)
	}
}

func TestCommandErrors(t *testing.T) {
	for _, args := range [][]string{
		{"-hz=0"}, {"-bus-bits=0"}, {"-states=-1"}, {"-states=no"},
		{"-unknown"}, {"unexpected"}, {"-states=18446744073709551615"},
	} {
		var out bytes.Buffer
		if err := run(&out, io.Discard, args); err == nil || out.Len() != 0 {
			t.Errorf("args %v: err=%v, output bytes=%d", args, err, out.Len())
		}
	}
}

type brokenWriter struct{}

func (brokenWriter) Write([]byte) (int, error) { return 0, io.ErrClosedPipe }

func TestCommandOutputFailure(t *testing.T) {
	if err := run(brokenWriter{}, io.Discard, nil); err == nil {
		t.Fatal("output failure ignored")
	}
}
