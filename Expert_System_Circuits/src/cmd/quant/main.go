// Command quant prints deterministic synthetic evidence for chapter 2.
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/nick-fedchik/articles/Expert_System_Circuits/src/quant"
)

func main() {
	e, err := quant.SyntheticEvidence()
	if err == nil {
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		err = encoder.Encode(e)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
