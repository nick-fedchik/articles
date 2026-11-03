package teosbaseline

import (
	"embed"
	"fmt"
	"strings"
)

//go:embed case.go arm_a.go arm_b.go
var armSources embed.FS

// SpecRegion names the source file and marker of each arm's rule code.
var SpecRegion = map[string]string{"A": "arm_a.go", "B0": "arm_b.go", "B1": "arm_b.go", "C": "case.go"}

// SpecLines counts the rule/specification lines of an arm deterministically:
// lines strictly between "// spec:<arm>:begin" and "// spec:<arm>:end" in the
// embedded source, excluding blank lines and lines that start with "//".
// Shared libraries (teoskernel, the Horn engine, typed units, projections) are
// not counted for any arm.
func SpecLines(arm string) (int, error) {
	file, ok := SpecRegion[arm]
	if !ok {
		return 0, fmt.Errorf("teosbaseline: unknown arm %q", arm)
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
		return 0, fmt.Errorf("teosbaseline: markers for arm %q not found or unterminated", arm)
	}
	return count, nil
}
