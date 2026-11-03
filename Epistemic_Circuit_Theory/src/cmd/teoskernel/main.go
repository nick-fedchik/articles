// Command teoskernel prints deterministic JSON for the TEOS admission kernel K:
//
//	teoskernel -mode example
//	teoskernel -mode exhaustive -norms 2 -conditions 1 -exceptions
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/nick-fedchik/articles/Epistemic_Circuit_Theory/src/teoskernel"
)

func main() {
	mode := flag.String("mode", "example", "example or exhaustive")
	norms := flag.Int("norms", 1, "maximum norms per theory (1 or 2)")
	conditions := flag.Int("conditions", 1, "maximum condition literals per norm (0..2)")
	exceptions := flag.Bool("exceptions", false, "enumerate optional exception literals")
	flag.Parse()
	var out any
	var err error
	switch *mode {
	case "example":
		spec, input := teoskernel.Example()
		var certificate teoskernel.Certificate
		if certificate, err = teoskernel.Evaluate(spec, input); err == nil {
			out = struct {
				Spec        teoskernel.Spec
				Input       teoskernel.Input
				Certificate teoskernel.Certificate
				Safety      []string
			}{spec, input, certificate, teoskernel.CheckSafety(spec, input)}
		}
	case "exhaustive":
		var report teoskernel.Report
		report, err = teoskernel.Exhaustive(teoskernel.Bounds{MaxNorms: *norms, MaxConditions: *conditions, Exceptions: *exceptions})
		if err == nil && len(report.Counterexamples) > 0 {
			err = fmt.Errorf("%d counterexamples", len(report.Counterexamples))
		}
		out = report
		fmt.Fprintf(os.Stderr, "elapsed %.1fs\n", report.ElapsedSeconds)
	default:
		err = fmt.Errorf("unknown mode %q", *mode)
	}
	if out != nil {
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		if encodeErr := encoder.Encode(out); err == nil {
			err = encodeErr
		}
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
