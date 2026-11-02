package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	evidence "github.com/nick-fedchik/articles/Epistemic_Circuit_Theory/src"
)

func main() {
	samples := flag.Int("samples", 20, "fixed number of timed samples per case")
	flag.Parse()
	rows, err := evidence.MeasureLayouts([]int{256, 10000, 100000}, *samples)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(rows); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
