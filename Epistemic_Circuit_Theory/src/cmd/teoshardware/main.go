package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/nick-fedchik/articles/Epistemic_Circuit_Theory/src/teoshardware"
)

func run(out, diagnostics io.Writer, args []string) error {
	config := teoshardware.DefaultConfig()
	flags := flag.NewFlagSet("teoshardware", flag.ContinueOnError)
	flags.SetOutput(diagnostics)
	flags.Uint64Var(&config.States, "states", config.States, "number of stored states")
	flags.Uint64Var(&config.BusBits, "bus-bits", config.BusBits, "payload bits accepted per cycle")
	flags.Uint64Var(&config.ClockHz, "hz", config.ClockHz, "assumed bus and hash clock in Hz")
	flags.Uint64Var(&config.Leaves, "leaves", config.Leaves, "literal provenance leaf descriptors")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("unexpected positional arguments")
	}
	report, err := teoshardware.Calculate(config)
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(out)
	encoder.SetIndent("", "  ")
	return encoder.Encode(report)
}

func main() {
	if err := run(os.Stdout, os.Stderr, os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
