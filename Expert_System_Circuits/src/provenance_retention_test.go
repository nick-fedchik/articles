package evidence

import (
	"runtime"
	"testing"
)

func TestDenseRootDoesNotKeepMultiMegabyteArenaAlive(t *testing.T) {
	c, err := GenerateProvenanceCase("balanced-all", 4096, "resolved", 81001)
	if err != nil {
		t.Fatal(err)
	}
	runtime.GC()
	var before, after runtime.MemStats
	runtime.ReadMemStats(&before)
	snapshot, err := c.Program.BuildProvenance(c.Inputs, ProvenanceDense, 64<<20)
	if err != nil {
		t.Fatal(err)
	}
	runtime.GC()
	runtime.ReadMemStats(&after)
	// Keeping c alive removes program/source reclamation from the heap difference.
	runtime.KeepAlive(c)
	if after.HeapAlloc > before.HeapAlloc+(1<<20) {
		t.Fatalf("tiny root kept a multi-megabyte working arena: before=%d after=%d", before.HeapAlloc, after.HeapAlloc)
	}
	result, err := snapshot.Result()
	if err != nil || len(result.Sources) != 4096 || snapshot.Stats().RetainedBytes >= 1024 {
		t.Fatal("invalid retained root", err)
	}
}
