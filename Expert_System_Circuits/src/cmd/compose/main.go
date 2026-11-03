package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/nick-fedchik/articles/Expert_System_Circuits/src/compose"
)

func main() {
	report, err := compose.Study()
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
