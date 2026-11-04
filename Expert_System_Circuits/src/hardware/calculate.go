// Package hardware dimensions finite representations of the СТЕС state.
// It does not simulate hardware, validate custody, or measure execution time.
package hardware

import (
	"fmt"
	"math"
)

const (
	SHA256DigestBytes = uint64(32)
	SHA256BlockBytes  = uint64(64)
	SHA256Rounds      = uint64(64)
)

type Field struct {
	Name string `json:"name"`
	Bits uint64 `json:"bits"`
}

type Layout struct {
	Name          string  `json:"name"`
	Fields        []Field `json:"fields"`
	AlignmentBits uint64  `json:"record_alignment_bits"`
	PackedBits    uint64  `json:"packed_bits"`
	PackedBytes   uint64  `json:"byte_packed_record_bytes"`
	AlignedBits   uint64  `json:"aligned_record_bits"`
	AlignedBytes  uint64  `json:"aligned_record_bytes"`
}

func add(a, b uint64) (uint64, error) {
	if b > math.MaxUint64-a {
		return 0, fmt.Errorf("addition overflows uint64")
	}
	return a + b, nil
}

func multiply(a, b uint64) (uint64, error) {
	if a != 0 && b > math.MaxUint64/a {
		return 0, fmt.Errorf("multiplication overflows uint64")
	}
	return a * b, nil
}

func ceilDiv(n, d uint64) uint64 {
	q := n / d
	if n%d != 0 {
		q++
	}
	return q
}

// NewLayout aligns the whole record, not individual fields and not a Go struct.
func NewLayout(name string, fields []Field, alignmentBits uint64) (Layout, error) {
	if name == "" || len(fields) == 0 || alignmentBits == 0 || alignmentBits%8 != 0 {
		return Layout{}, fmt.Errorf("name, fields and positive byte-multiple alignment required")
	}
	result := Layout{Name: name, Fields: append([]Field(nil), fields...), AlignmentBits: alignmentBits}
	seen := make(map[string]bool)
	for _, field := range fields {
		if field.Name == "" || field.Bits == 0 || seen[field.Name] {
			return Layout{}, fmt.Errorf("field names must be nonempty and unique; widths positive")
		}
		seen[field.Name] = true
		var err error
		result.PackedBits, err = add(result.PackedBits, field.Bits)
		if err != nil {
			return Layout{}, err
		}
	}
	var err error
	result.AlignedBits, err = multiply(ceilDiv(result.PackedBits, alignmentBits), alignmentBits)
	if err != nil {
		return Layout{}, err
	}
	result.PackedBytes = ceilDiv(result.PackedBits, 8)
	result.AlignedBytes = result.AlignedBits / 8
	return result, nil
}

// StudyLayouts are alternatives, not a canonical СТЕС ABI. Both beta scores
// remain present; evidence flags do not replace the graded Belnap coordinates.
func StudyLayouts() ([]Layout, error) {
	base := []Field{
		{"value", 32}, {"uncertainty", 32},
		{"beta_pro", 16}, {"beta_contra", 16}, {"evidence_flags", 8},
		{"si_exponents", 28}, {"type_descriptor", 12},
		{"data_revision", 32}, {"policy_revision", 32}, {"timestamp", 64},
	}
	variants := []struct {
		name string
		prov []Field
	}{
		{"reference64", []Field{{"provenance_reference", 64}}},
		{"digest256", []Field{{"sha256_digest", 256}}},
		{"leaf_custody", []Field{
			{"sha256_digest", 256}, {"document_id", 32},
			{"quote_start", 32}, {"quote_end_exclusive", 32},
		}},
	}
	var result []Layout
	for _, variant := range variants {
		fields := append(append([]Field(nil), base...), variant.prov...)
		layout, err := NewLayout(variant.name, fields, 128)
		if err != nil {
			return nil, err
		}
		result = append(result, layout)
	}
	wide := append([]Field(nil), base...)
	wide[0].Bits, wide[1].Bits = 64, 64
	wide[5].Bits, wide[6].Bits = 56, 16
	wide = append(wide, variants[2].prov...)
	layout, err := NewLayout("wide_leaf_custody", wide, 128)
	if err != nil {
		return nil, err
	}
	return append(result, layout), nil
}

type Geometry struct {
	Name     string `json:"name"`
	WordBits uint64 `json:"word_bits"`
	Depth    uint64 `json:"depth_words"`
}

// RAMB36SDP uses all 72 bits including parity lanes as data, with ECC disabled.
// This is an organization reference, not a selected FPGA or inferred netlist.
func RAMB36SDP() Geometry {
	return Geometry{"RAMB36E1_SDP_ECC_off", 72, 512}
}

// ICE40EBR is one iCE40 SB_RAM40_4K embedded block in its 256 x 16 mode.
// ICE40HX8KEBRBlocks is the number of such blocks on iCE40HX8K
// (32 x 4 Kbit = 131072 bits = 16 KiB), used to check whether a bank fits.
func ICE40EBR() Geometry {
	return Geometry{"iCE40_SB_RAM40_4K_256x16", 16, 256}
}

const ICE40HX8KEBRBlocks = uint64(32)

type Memory struct {
	RecordBits         uint64   `json:"record_bits"`
	Records            uint64   `json:"records"`
	PayloadBits        uint64   `json:"payload_bits"`
	CapacityBytesMin   uint64   `json:"dense_capacity_bytes_min"`
	Geometry           Geometry `json:"reference_geometry"`
	CapacityBlocksMin  uint64   `json:"capacity_blocks_min"`
	WidthLanes         uint64   `json:"width_lanes"`
	DepthBanks         uint64   `json:"depth_banks"`
	RectangularTiles   uint64   `json:"rectangular_tiles"`
	TileAllocationBits uint64   `json:"tile_allocation_bits"`
}

// DimensionMemory assumes one entire record per addressed row, one read and
// one write port. Packing multiple records per row or extra ports changes tiling.
func DimensionMemory(recordBits, records uint64, geometry Geometry) (Memory, error) {
	if recordBits == 0 || geometry.WordBits == 0 || geometry.Depth == 0 {
		return Memory{}, fmt.Errorf("record and geometry dimensions must be positive")
	}
	total, err := multiply(recordBits, records)
	if err != nil {
		return Memory{}, err
	}
	tileBits, err := multiply(geometry.WordBits, geometry.Depth)
	if err != nil {
		return Memory{}, err
	}
	lanes, banks := ceilDiv(recordBits, geometry.WordBits), ceilDiv(records, geometry.Depth)
	tiles, err := multiply(lanes, banks)
	if err != nil {
		return Memory{}, err
	}
	allocation, err := multiply(tiles, tileBits)
	if err != nil {
		return Memory{}, err
	}
	return Memory{recordBits, records, total, ceilDiv(total, 8), geometry,
		ceilDiv(total, tileBits), lanes, banks, tiles, allocation}, nil
}

type Transfer struct {
	Bits          uint64  `json:"payload_bits"`
	BusBits       uint64  `json:"bus_bits"`
	ClockHz       uint64  `json:"clock_hz"`
	BeatsMin      uint64  `json:"beats_min"`
	DurationNSMin float64 `json:"duration_ns_min"`
}

// Serialize assumes one accepted beat per cycle, no framing, stalls, memory
// read latency, arbitration or CDC. This lower bound is not a WCET.
func Serialize(bits, busBits, clockHz uint64) (Transfer, error) {
	if busBits == 0 || clockHz == 0 {
		return Transfer{}, fmt.Errorf("bus width and clock must be positive")
	}
	beats := ceilDiv(bits, busBits)
	return Transfer{bits, busBits, clockHz, beats, float64(beats) * 1e9 / float64(clockHz)}, nil
}

// SerializeRecords starts every record on a new beat. Unlike Serialize over
// the concatenated payload, unused lanes cannot carry the following record.
func SerializeRecords(recordBits, records, busBits, clockHz uint64) (Transfer, error) {
	record, err := Serialize(recordBits, busBits, clockHz)
	if err != nil {
		return Transfer{}, err
	}
	totalBits, err := multiply(recordBits, records)
	if err != nil {
		return Transfer{}, err
	}
	beats, err := multiply(record.BeatsMin, records)
	if err != nil {
		return Transfer{}, err
	}
	return Transfer{totalBits, busBits, clockHz, beats, float64(beats) * 1e9 / float64(clockHz)}, nil
}

type Provenance struct {
	Leaves                 uint64 `json:"leaf_records"`
	LeafBytes              uint64 `json:"leaf_descriptor_bytes"`
	LiteralChainBytes      uint64 `json:"literal_chain_bytes"`
	DigestBytes            uint64 `json:"fixed_digest_bytes"`
	ReferenceBytes         uint64 `json:"fixed_reference_bytes"`
	UnsharedSidecarPlusRef uint64 `json:"unshared_sidecar_plus_reference_bytes"`
	RetentionCapLeaves     uint64 `json:"retention_cap_leaves"`
	BoundedChainBytesMax   uint64 `json:"bounded_chain_bytes_max"`
	FitsRetentionCap       bool   `json:"fits_retention_cap"`
}

// DimensionProvenance counts literal leaf descriptors only, excluding internal
// proof nodes, signatures and certificates. A cap cannot silently truncate proof.
func DimensionProvenance(leaves, leafBytes, capLeaves uint64) (Provenance, error) {
	if leafBytes == 0 || capLeaves == 0 {
		return Provenance{}, fmt.Errorf("leaf size and retention cap must be positive")
	}
	literal, err := multiply(leaves, leafBytes)
	if err != nil {
		return Provenance{}, err
	}
	bounded, err := multiply(capLeaves, leafBytes)
	if err != nil {
		return Provenance{}, err
	}
	sidecar, err := add(literal, 8)
	if err != nil {
		return Provenance{}, err
	}
	return Provenance{leaves, leafBytes, literal, SHA256DigestBytes, 8,
		sidecar, capLeaves, bounded, leaves <= capLeaves}, nil
}

type HashCost struct {
	MessageBytes                    uint64  `json:"message_bytes"`
	PaddedBlocks                    uint64  `json:"padded_blocks"`
	CompressionRounds               uint64  `json:"compression_rounds"`
	IterativeOverheadCyclesPerBlock uint64  `json:"iterative_overhead_cycles_per_block"`
	IterativeCycles                 uint64  `json:"iterative_one_round_per_cycle_cycles"`
	IterativeNS                     float64 `json:"iterative_ns"`
	UnrolledCombinationalCycles     uint64  `json:"hypothetical_fully_combinational_block_cycles"`
	Pipeline64LatencyCycles         uint64  `json:"hypothetical_64_stage_same_message_latency_cycles"`
	PipelineIndependentBlockII      uint64  `json:"hypothetical_pipeline_independent_block_ii_cycles"`
}

// DimensionSHA256 models three declared architectures, not universal hardware
// lower bounds. Full unrolling requires timing closure; a 64-stage pipeline's
// II=1 applies to independent blocks, not dependent blocks of one message.
func DimensionSHA256(messageBytes, clockHz, iterativeOverhead uint64) (HashCost, error) {
	if clockHz == 0 || messageBytes > (uint64(1)<<61)-1 {
		return HashCost{}, fmt.Errorf("positive clock and SHA-256 bit length below 2^64 required")
	}
	blocks := messageBytes/SHA256BlockBytes + 1
	if messageBytes%SHA256BlockBytes >= 56 {
		blocks++
	}
	perBlock, err := add(SHA256Rounds, iterativeOverhead)
	if err != nil {
		return HashCost{}, err
	}
	cycles, err := multiply(blocks, perBlock)
	if err != nil {
		return HashCost{}, err
	}
	rounds, err := multiply(blocks, SHA256Rounds)
	if err != nil {
		return HashCost{}, err
	}
	return HashCost{messageBytes, blocks, rounds, iterativeOverhead, cycles,
		float64(cycles) * 1e9 / float64(clockHz), blocks, rounds, 1}, nil
}

type GateInput struct {
	CustodyValidated  bool   `json:"custody_validated"`
	BindingValid      bool   `json:"binding_valid"`
	SIMatch           bool   `json:"si_match"`
	NumericValid      bool   `json:"numeric_valid"`
	NormPermitted     bool   `json:"norm_permitted"`
	Undercut          bool   `json:"undercut"`
	Evidence          uint8  `json:"evidence_code"`
	DataRevision      uint32 `json:"data_revision"`
	ValidatedRevision uint32 `json:"validated_revision"`
	PolicyRevision    uint32 `json:"policy_revision"`
	ExpectedPolicy    uint32 `json:"expected_policy_revision"`
	Timestamp         uint64 `json:"timestamp_ticks"`
	Now               uint64 `json:"now_ticks"`
	MaxAge            uint64 `json:"max_age_ticks"`
	UpperRaw          int64  `json:"prevalidated_upper_bound_raw"`
	LimitRaw          int64  `json:"limit_raw"`
}

type GateDecision struct {
	Allow   bool     `json:"allow"`
	Reasons []string `json:"reasons"`
}

// PrevalidatedGate is a functional predicate only, NOT a custody validator or
// timed RTL. Evidence codes: 0=neither, 1=pro only, 2=contra only, 3=both.
// All inputs must be coherently registered and bound by trusted validation;
// an attacker-controlled "validated" flag establishes no authenticity.
func PrevalidatedGate(in GateInput) GateDecision {
	result := GateDecision{Reasons: []string{}}
	checks := []struct {
		ok     bool
		reason string
	}{
		{in.CustodyValidated, "custody_not_validated"},
		{in.BindingValid && in.DataRevision == in.ValidatedRevision, "validation_binding_invalid"},
		{in.SIMatch, "si_mismatch"},
		{in.NumericValid, "numeric_invalid"},
		{in.NormPermitted && !in.Undercut, "norm_blocked"},
		{in.Evidence == 1, "evidence_not_pro_only"},
		{in.PolicyRevision == in.ExpectedPolicy, "policy_revision_mismatch"},
		{in.Timestamp <= in.Now && in.Now-in.Timestamp <= in.MaxAge, "timestamp_invalid_or_stale"},
		{in.UpperRaw <= in.LimitRaw, "guard_band_exceeds_limit"},
	}
	for _, check := range checks {
		if !check.ok {
			result.Reasons = append(result.Reasons, check.reason)
		}
	}
	result.Allow = len(result.Reasons) == 0
	return result
}
