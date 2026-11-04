package hardware

import (
	"encoding/json"
	"math"
	"reflect"
	"testing"
)

func TestStudyLayouts(t *testing.T) {
	layouts, err := StudyLayouts()
	if err != nil {
		t.Fatal(err)
	}
	want := [][3]uint64{{336, 42, 48}, {528, 66, 80}, {624, 78, 80}, {720, 90, 96}}
	for i, layout := range layouts {
		if got := [3]uint64{layout.PackedBits, layout.PackedBytes, layout.AlignedBytes}; got != want[i] {
			t.Errorf("%s: got %v, want %v", layout.Name, got, want[i])
		}
		if layout.AlignmentBits != 128 || layout.AlignedBits != layout.AlignedBytes*8 {
			t.Errorf("inconsistent alignment: %+v", layout)
		}
	}
}

func TestLayoutBoundaries(t *testing.T) {
	tests := []struct {
		bits, alignment, bytes, aligned uint64
	}{{1, 8, 1, 8}, {8, 8, 1, 8}, {9, 8, 2, 16}, {64, 64, 8, 64}, {65, 64, 9, 128}}
	for _, tt := range tests {
		got, err := NewLayout("test", []Field{{"v", tt.bits}}, tt.alignment)
		if err != nil || got.PackedBytes != tt.bytes || got.AlignedBits != tt.aligned {
			t.Errorf("%+v: got %+v, err %v", tt, got, err)
		}
	}
	fields := []Field{{"v", 8}}
	layout, _ := NewLayout("copy", fields, 8)
	fields[0].Bits = 16
	if layout.Fields[0].Bits != 8 {
		t.Fatal("layout aliases input fields")
	}
}

func TestLayoutErrors(t *testing.T) {
	for _, tt := range []struct {
		name   string
		fields []Field
		align  uint64
	}{
		{"", []Field{{"a", 1}}, 8},
		{"empty", nil, 8},
		{"zero", []Field{{"a", 0}}, 8},
		{"no_name", []Field{{"", 1}}, 8},
		{"duplicate", []Field{{"a", 1}, {"a", 2}}, 8},
		{"align_zero", []Field{{"a", 1}}, 0},
		{"align_not_bytes", []Field{{"a", 1}}, 7},
		{"sum_overflow", []Field{{"a", math.MaxUint64}, {"b", 1}}, 8},
		{"rounding_overflow", []Field{{"a", math.MaxUint64}}, 8},
	} {
		if _, err := NewLayout(tt.name, tt.fields, tt.align); err == nil {
			t.Errorf("expected error for %s", tt.name)
		}
	}
}

func TestMemoryCapacityVersusTiling(t *testing.T) {
	tests := []struct {
		width, records, capacity, lanes, banks, tiles uint64
	}{
		{72, 512, 1, 1, 1, 1},
		{73, 512, 2, 2, 1, 2},
		{72, 513, 2, 1, 2, 2},
		{512, 4096, 57, 8, 8, 64},
		{436, 4096, 49, 7, 8, 56},
		{1, 0, 0, 1, 0, 0},
	}
	for _, tt := range tests {
		got, err := DimensionMemory(tt.width, tt.records, RAMB36SDP())
		if err != nil {
			t.Fatal(err)
		}
		if [4]uint64{got.CapacityBlocksMin, got.WidthLanes, got.DepthBanks, got.RectangularTiles} !=
			[4]uint64{tt.capacity, tt.lanes, tt.banks, tt.tiles} {
			t.Errorf("%+v: got %+v", tt, got)
		}
		if got.TileAllocationBits < got.PayloadBits {
			t.Fatal("allocated fewer bits than payload")
		}
	}
}

// One reference64 bank (4096 records x 384 aligned bits = 192 KiB) needs at
// least 384 iCE40 EBR blocks (12x the device), while iCE40HX8K has 32 (16 KiB); a ping-pong
// pair needs twice as much. The Xilinx RAMB36 geometry is a separate,
// conditional reference.
func TestReference64BankVersusICE40HX8K(t *testing.T) {
	got, err := DimensionMemory(384, 4096, ICE40EBR())
	if err != nil {
		t.Fatal(err)
	}
	if got.CapacityBytesMin != 196608 || got.CapacityBlocksMin != 384 || got.RectangularTiles != 384 {
		t.Fatalf("unexpected iCE40 dimensioning: %+v", got)
	}
	available := ICE40HX8KEBRBlocks * got.Geometry.WordBits * got.Geometry.Depth
	if available != 131072 || got.PayloadBits != 12*available {
		t.Fatalf("available=%d payload=%d", available, got.PayloadBits)
	}
	xilinx, err := DimensionMemory(384, 4096, RAMB36SDP())
	if err != nil {
		t.Fatal(err)
	}
	if xilinx.RectangularTiles != 48 || 2*xilinx.RectangularTiles != 96 {
		t.Fatalf("unexpected RAMB36 tiling: %+v", xilinx)
	}
}

func TestMemoryErrors(t *testing.T) {
	for _, tt := range []struct {
		width, records uint64
		geometry       Geometry
	}{
		{0, 1, RAMB36SDP()},
		{1, 1, Geometry{}},
		{1, 1, Geometry{"zero_depth", 72, 0}},
		{math.MaxUint64, 2, RAMB36SDP()},
		{1, 1, Geometry{"tile_overflow", math.MaxUint64, 2}},
		{1, math.MaxUint64, Geometry{"allocation_overflow", 8, 1}},
		{math.MaxUint64, 1, Geometry{"lane_overflow", 1, 2}},
	} {
		if _, err := DimensionMemory(tt.width, tt.records, tt.geometry); err == nil {
			t.Errorf("expected error for %+v", tt)
		}
	}
}

func TestSerialize(t *testing.T) {
	for _, tt := range []struct {
		bits, bus, beats uint64
		ns               float64
	}{
		{0, 64, 0, 0}, {64, 64, 1, 10}, {65, 64, 2, 20},
		{512, 64, 8, 80}, {512, 256, 2, 20}, {512, 512, 1, 10},
		{768, 512, 2, 20}, {math.MaxUint64, math.MaxUint64, 1, 10},
	} {
		got, err := Serialize(tt.bits, tt.bus, 100_000_000)
		if err != nil || got.BeatsMin != tt.beats || math.Abs(got.DurationNSMin-tt.ns) > 1e-8 {
			t.Errorf("%+v: got %+v, %v", tt, got, err)
		}
	}
	if _, err := Serialize(1, 0, 1); err == nil {
		t.Fatal("zero bus width accepted")
	}
	if _, err := Serialize(1, 1, 0); err == nil {
		t.Fatal("zero clock accepted")
	}
}

func TestSerializeRecords(t *testing.T) {
	for _, tt := range []struct{ bus, beats uint64 }{
		{64, 24576}, {256, 8192}, {512, 4096},
	} {
		got, err := SerializeRecords(384, 4096, tt.bus, 100_000_000)
		if err != nil || got.Bits != 1572864 || got.BeatsMin != tt.beats {
			t.Errorf("%+v: got %+v, %v", tt, got, err)
		}
	}
	for _, args := range [][4]uint64{
		{1, 1, 0, 1}, {math.MaxUint64, 2, 1, 1},
	} {
		if _, err := SerializeRecords(args[0], args[1], args[2], args[3]); err == nil {
			t.Errorf("accepted invalid sequence %v", args)
		}
	}
}

func TestProvenance(t *testing.T) {
	got, err := DimensionProvenance(1024, 44, 16)
	if err != nil || got.LiteralChainBytes != 45056 || got.UnsharedSidecarPlusRef != 45064 ||
		got.BoundedChainBytesMax != 704 || got.FitsRetentionCap ||
		got.DigestBytes != 32 || got.ReferenceBytes != 8 {
		t.Fatalf("unexpected provenance: %+v, %v", got, err)
	}
	for _, leaves := range []uint64{0, 16, 17} {
		got, err := DimensionProvenance(leaves, 44, 16)
		if err != nil || got.FitsRetentionCap != (leaves <= 16) {
			t.Fatalf("unexpected cap decision: %+v, %v", got, err)
		}
	}
	for _, args := range [][3]uint64{
		{1, 0, 1}, {1, 1, 0}, {math.MaxUint64, 2, 1},
		{1, 2, math.MaxUint64}, {math.MaxUint64, 1, 1},
	} {
		if _, err := DimensionProvenance(args[0], args[1], args[2]); err == nil {
			t.Errorf("accepted invalid provenance arguments %v", args)
		}
	}
}

func TestSHA256PaddingAndArchitectures(t *testing.T) {
	for _, tt := range []struct{ bytes, blocks uint64 }{
		{0, 1}, {55, 1}, {56, 2}, {63, 2}, {64, 2},
		{119, 2}, {120, 3}, {1024, 17}, {45056, 705},
		{(uint64(1) << 61) - 1, (uint64(1) << 55) + 1},
	} {
		got, err := DimensionSHA256(tt.bytes, 100_000_000, 2)
		if err != nil || got.PaddedBlocks != tt.blocks || got.CompressionRounds != tt.blocks*64 ||
			got.IterativeCycles != tt.blocks*66 || got.UnrolledCombinationalCycles != tt.blocks ||
			got.Pipeline64LatencyCycles != tt.blocks*64 || got.PipelineIndependentBlockII != 1 {
			t.Errorf("%+v: got %+v, %v", tt, got, err)
		}
	}
	got, _ := DimensionSHA256(1024, 100_000_000, 2)
	if math.Abs(got.IterativeNS-11220) > 1e-8 {
		t.Fatalf("1 KiB iterative duration: %g ns", got.IterativeNS)
	}
	for _, args := range [][3]uint64{
		{0, 0, 0}, {uint64(1) << 61, 1, 0}, {0, 1, math.MaxUint64},
		{56, 1, math.MaxUint64 / 2},
	} {
		if _, err := DimensionSHA256(args[0], args[1], args[2]); err == nil {
			t.Errorf("accepted invalid hash arguments %v", args)
		}
	}
}

func safeGate() GateInput {
	return GateInput{
		CustodyValidated: true, BindingValid: true, SIMatch: true,
		NumericValid: true, NormPermitted: true, Evidence: 1,
		DataRevision: 7, ValidatedRevision: 7, PolicyRevision: 2, ExpectedPolicy: 2,
		Timestamp: 100, Now: 105, MaxAge: 10, UpperRaw: 622592, LimitRaw: 655360,
	}
}

func TestPrevalidatedGateFailClosed(t *testing.T) {
	tests := []struct {
		name   string
		change func(*GateInput)
	}{
		{"custody", func(in *GateInput) { in.CustodyValidated = false }},
		{"binding", func(in *GateInput) { in.BindingValid = false }},
		{"SI", func(in *GateInput) { in.SIMatch = false }},
		{"numeric", func(in *GateInput) { in.NumericValid = false }},
		{"norm", func(in *GateInput) { in.NormPermitted = false }},
		{"undercut", func(in *GateInput) { in.Undercut = true }},
		{"revision", func(in *GateInput) { in.DataRevision++ }},
		{"policy", func(in *GateInput) { in.PolicyRevision++ }},
		{"future", func(in *GateInput) { in.Timestamp = in.Now + 1 }},
		{"expired", func(in *GateInput) { in.Now = in.Timestamp + in.MaxAge + 1 }},
		{"guard_band", func(in *GateInput) { in.UpperRaw = in.LimitRaw + 1 }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := safeGate()
			tt.change(&in)
			if decision := PrevalidatedGate(in); decision.Allow || len(decision.Reasons) == 0 {
				t.Fatalf("failure passed: %+v", decision)
			}
		})
	}
	for evidence := 0; evidence <= 255; evidence++ {
		in := safeGate()
		in.Evidence = uint8(evidence)
		if PrevalidatedGate(in).Allow != (evidence == 1) {
			t.Errorf("incorrect evidence verdict for %d", evidence)
		}
	}
	in := safeGate()
	in.UpperRaw = in.LimitRaw
	in.Now = in.Timestamp + in.MaxAge
	if !PrevalidatedGate(in).Allow {
		t.Fatal("inclusive bound/expiry equality rejected")
	}
	if PrevalidatedGate(GateInput{}).Allow {
		t.Fatal("uninitialized snapshot accepted")
	}
}

func TestDefaultReportAndDeterminism(t *testing.T) {
	config := DefaultConfig()
	report, err := Calculate(config)
	if err != nil {
		t.Fatal(err)
	}
	ref := report.Layouts[0]
	if ref.AlignedMemory.CapacityBytesMin != 196608 ||
		ref.AlignedMemory.CapacityBlocksMin != 43 || ref.AlignedMemory.RectangularTiles != 48 ||
		ref.AlignedRefill.BeatsMin != 24576 || math.Abs(ref.AlignedRefill.DurationNSMin-245760) > 1e-8 {
		t.Fatalf("default reference calculation: %+v", ref)
	}
	if !report.GateExamples[0].Decision.Allow || report.GateExamples[1].Decision.Allow {
		t.Fatal("incorrect pressure examples")
	}
	if report.UnsharedSidecarAllStatesBytes != 184549376 {
		t.Fatal("unshared provenance sidecar capacity not accounted for")
	}
	again, err := Calculate(config)
	if err != nil {
		t.Fatal(err)
	}
	a, _ := json.Marshal(report)
	b, _ := json.Marshal(again)
	if !reflect.DeepEqual(a, b) {
		t.Fatal("JSON is not deterministic")
	}
	for _, c := range []Config{
		{1, 0, 1, 1}, {1, 1, 0, 1},
		{math.MaxUint64, 64, 1, 1}, {1, 64, 1, math.MaxUint64},
		{1_000_000_000, 64, 1, 1_000_000_000},
	} {
		if _, err := Calculate(c); err == nil {
			t.Errorf("accepted invalid config %+v", c)
		}
	}
	config.States, config.Leaves = 0, 0
	empty, err := Calculate(config)
	if err != nil || empty.Layouts[0].AlignedRefill.BeatsMin != 0 ||
		empty.Provenance.LiteralChainBytes != 0 {
		t.Fatalf("empty design: %+v, %v", empty, err)
	}
}
