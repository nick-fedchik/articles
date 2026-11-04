package baseline

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

//go:embed case.go arm_a.go arm_b.go units.go
var armSources embed.FS

// SpecRegion names the source file and marker of each arm's rule code.
var SpecRegion = map[string]string{"A": "arm_a.go", "B0": "arm_b.go", "B1": "arm_b.go", "C": "case.go"}

func countCodeLines(src string) int {
	count := 0
	for _, line := range strings.Split(src, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" && !strings.HasPrefix(trimmed, "//") {
			count++
		}
	}
	return count
}

// SpecLines counts the rule/specification lines of an arm deterministically:
// lines strictly between "// spec:<arm>:begin" and "// spec:<arm>:end" in the
// embedded source, excluding blank lines and lines that start with "//".
// Shared libraries (kernel, the Horn engine, typed units, projections) are
// reported separately by SharedEngineLines (F-06).
func SpecLines(arm string) (int, error) {
	file, ok := SpecRegion[arm]
	if !ok {
		return 0, fmt.Errorf("baseline: unknown arm %q", arm)
	}
	data, err := armSources.ReadFile(file)
	if err != nil {
		return 0, err
	}
	begin, end := "// spec:"+arm+":begin", "// spec:"+arm+":end"
	inside, found, count := false, false, 0
	for _, line := range strings.Split(string(data), "\n") {
		trimmed := strings.TrimSpace(line)
		switch {
		case trimmed == begin:
			inside, found = true, true
		case trimmed == end:
			inside = false
		case inside && trimmed != "" && !strings.HasPrefix(trimmed, "//"):
			count++
		}
	}
	if !found || inside {
		return 0, fmt.Errorf("baseline: markers for arm %q not found or unterminated", arm)
	}
	return count, nil
}

// SharedEngineLines reports non-blank, non-comment lines of the shared support
// code and the kernel TCB so that comparing B1 (imperative Go rules +
// local helpers) against C (declarative CaseSpec + kernel engine/verifier)
// accounts for the underlying engine size (F-06).
func SharedEngineLines() (map[string]int, error) {
	unitsData, err := armSources.ReadFile("units.go")
	if err != nil {
		return nil, err
	}
	armBData, err := armSources.ReadFile("arm_b.go")
	if err != nil {
		return nil, err
	}
	b0Spec, err := SpecLines("B0")
	if err != nil {
		return nil, err
	}
	b1Spec, err := SpecLines("B1")
	if err != nil {
		return nil, err
	}
	out := map[string]int{
		"conventional_units":   countCodeLines(string(unitsData)),
		"conventional_support": countCodeLines(string(armBData)) - b0Spec - b1Spec,
	}
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		return nil, fmt.Errorf("baseline: runtime.Caller unavailable")
	}
	kernelDir := filepath.Join(filepath.Dir(thisFile), "..", "kernel")
	kernelFiles := map[string]string{
		"kernel_evaluate": "kernel.go",
		"kernel_verify":   "verify.go",
		"kernel_safety":   "safety.go",
	}
	totalKernel := 0
	for key, name := range kernelFiles {
		data, err := os.ReadFile(filepath.Join(kernelDir, name))
		if err != nil {
			return nil, err
		}
		n := countCodeLines(string(data))
		out[key] = n
		totalKernel += n
	}
	out["kernel_total"] = totalKernel
	return out, nil
}
