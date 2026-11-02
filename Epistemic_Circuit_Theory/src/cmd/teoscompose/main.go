package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/nick-fedchik/articles/Epistemic_Circuit_Theory/src/teoscompose"
)

func main() {
	report, err := teoscompose.Study()
	if err == nil {
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		err = encoder.Encode(report)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
