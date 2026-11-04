package hardware

type Config struct {
	States  uint64 `json:"states"`
	BusBits uint64 `json:"bus_bits"`
	ClockHz uint64 `json:"clock_hz"`
	Leaves  uint64 `json:"provenance_leaves"`
}

func DefaultConfig() Config {
	return Config{States: 4096, BusBits: 64, ClockHz: 100_000_000, Leaves: 1024}
}

type LayoutCase struct {
	Layout             Layout   `json:"layout"`
	PackedMemory       Memory   `json:"bit_packed_memory"`
	AlignedMemory      Memory   `json:"aligned_memory"`
	BytePackedBytes    uint64   `json:"byte_packed_total_bytes"`
	PackedTransfer     Transfer `json:"bit_packed_record_transfer"`
	AlignedTransfer    Transfer `json:"aligned_record_transfer"`
	AlignedRefill      Transfer `json:"aligned_memory_refill"`
	RecordPaddedRefill Transfer `json:"record_padded_memory_refill"`
}

type GateExample struct {
	Name     string       `json:"name"`
	Input    GateInput    `json:"input"`
	Decision GateDecision `json:"decision"`
}

type Operation struct {
	Name          string `json:"name"`
	Scope         string `json:"scope"`
	Qualification string `json:"qualification"`
}

type Report struct {
	SchemaVersion                 uint64        `json:"schema_version"`
	Status                        string        `json:"status"`
	Config                        Config        `json:"config"`
	ClockPeriodNS                 float64       `json:"assumed_clock_period_ns"`
	Assumptions                   []string      `json:"assumptions"`
	Layouts                       []LayoutCase  `json:"layouts"`
	Provenance                    Provenance    `json:"provenance"`
	UnsharedSidecarAllStatesBytes uint64        `json:"unshared_provenance_sidecar_all_states_bytes"`
	HashModels                    []HashCost    `json:"hash_models"`
	GateExamples                  []GateExample `json:"gate_examples"`
	Operations                    []Operation   `json:"operations"`
	NotPerformed                  []string      `json:"not_performed"`
}

// Calculate returns deterministic synthetic design calculations: no clock
// sampling, Go object sizing, benchmarks, host timings, or FPGA tool execution.
func Calculate(config Config) (Report, error) {
	probe, err := Serialize(0, config.BusBits, config.ClockHz)
	if err != nil {
		return Report{}, err
	}
	result := Report{
		SchemaVersion: 1, Status: "calculated_not_hardware_measurement", Config: config,
		ClockPeriodNS: 1e9 / float64(probe.ClockHz),
		Assumptions: []string{
			"Layouts are finite design alternatives, not a fixed ABI or Go struct layout.",
			"Whole-record alignment is a chosen 128 bits; contiguous bit and byte packing are distinguished.",
			"One accepted bus beat per cycle; framing, stalls, arbitration, CDC and RAM read latency excluded.",
			"Aligned-memory refill is a dense-stream payload lower bound; record-padded refill starts every record on a new beat.",
			"RAMB36E1 512x72 SDP ECC-off is an organization reference only; parity lanes used as data.",
			"Rectangular tiling is calculated for one record per row, one read and one write; not mapped utilization.",
			"32-bit value: signed fixed point with 16 fraction bits; uncertainty unsigned with 16 fraction bits.",
			"beta scores use unsigned 16-bit integers scaled by 65535; flags retain only validated summaries.",
			"SI uses seven signed 4-bit exponents; type descriptor selects immutable scale/offset semantics.",
			"Wide layout: 64-bit value and uncertainty with 32 fraction bits; seven signed 8-bit SI exponents.",
			"Revision wrap, SI overflow, numeric overflow and unknown types must invalidate validation.",
			"Timestamp is a trusted monotonic counter; tick frequency, epoch and expiry policy are deployment choices.",
			"SHA-256: byte-aligned messages, 64-byte blocks, 32-byte digest, 64 rounds, 64-bit bit-length suffix.",
			"Iterative hash model assumes one round per cycle plus two cycles overhead per padded block.",
			"Unrolled one-block-per-cycle hash is hypothetical and requires timing closure, not a lower bound.",
			"64-stage hash pipeline has II=1 only for independent blocks; chained blocks depend on prior state.",
			"Hash core models exclude transfer, padding logic, trust-anchor/signature checks and custody traversal.",
			"Provenance leaves contain only 32-byte digest plus three 4-byte identifiers/offsets; not a full proof.",
			"Gate predicate consumes coherent trusted prevalidated summaries; it does not verify their authenticity.",
			"Gate summary/policy registers, arithmetic units, validation engine and bus staging storage are not included in state-memory capacity.",
			"One tick is conditional register-to-register gate response, not input-to-physical-actuation latency.",
		},
		Operations: []Operation{
			{"registered flags, revision/SI equality, evidence code, expiry and bound comparison",
				"conditional_prevalidated_gate", "Bounded combinational logic is plausible; frequency and route must be verified."},
			{"SHA-256, source-byte matching, proof/custody traversal and signature validation",
				"outside_prevalidated_gate", "Digest equality is not full custody validation; traversal and I/O must be bounded."},
			{"dual multiplication, uncertainty propagation, Hamacher division and interval construction",
				"separate_arithmetic_stage", "Latency depends on precision, rounding, overflow handling, mapping and clock."},
			{"input synchronization, RAM fetch, bus assembly and policy updates",
				"outside_gate_response", "Snapshot assembly and atomic commit precede the gate; partial records fail closed."},
			{"output driver, actuator and physical plant",
				"outside_gate_response", "No one-tick actuation assertion follows from a registered inhibit signal."},
		},
		NotPerformed: []string{
			"RTL simulation or synthesis for this calculator",
			"identified target device mapping or resource utilization",
			"post-route timing or power analysis",
			"HIL or physical actuation measurement",
			"WCET or certification assessment",
		},
	}
	layouts, err := StudyLayouts()
	if err != nil {
		return Report{}, err
	}
	for _, layout := range layouts {
		packed, err := DimensionMemory(layout.PackedBits, config.States, RAMB36SDP())
		if err != nil {
			return Report{}, err
		}
		aligned, err := DimensionMemory(layout.AlignedBits, config.States, RAMB36SDP())
		if err != nil {
			return Report{}, err
		}
		bytes, err := multiply(layout.PackedBytes, config.States)
		if err != nil {
			return Report{}, err
		}
		packedTransfer, _ := Serialize(layout.PackedBits, config.BusBits, config.ClockHz)
		alignedTransfer, _ := Serialize(layout.AlignedBits, config.BusBits, config.ClockHz)
		refill, _ := Serialize(aligned.PayloadBits, config.BusBits, config.ClockHz)
		paddedRefill, err := SerializeRecords(layout.AlignedBits, config.States, config.BusBits, config.ClockHz)
		if err != nil {
			return Report{}, err
		}
		result.Layouts = append(result.Layouts, LayoutCase{
			layout, packed, aligned, bytes, packedTransfer, alignedTransfer, refill, paddedRefill,
		})
	}
	result.Provenance, err = DimensionProvenance(config.Leaves, 44, 16)
	if err != nil {
		return Report{}, err
	}
	result.UnsharedSidecarAllStatesBytes, err = multiply(result.Provenance.LiteralChainBytes, config.States)
	if err != nil {
		return Report{}, err
	}
	for _, size := range []uint64{55, 56, 64, 1024, result.Provenance.LiteralChainBytes} {
		cost, err := DimensionSHA256(size, config.ClockHz, 2)
		if err != nil {
			return Report{}, err
		}
		result.HashModels = append(result.HashModels, cost)
	}
	input := GateInput{
		CustodyValidated: true, BindingValid: true, SIMatch: true,
		NumericValid: true, NormPermitted: true, Evidence: 1,
		DataRevision: 7, ValidatedRevision: 7, PolicyRevision: 2, ExpectedPolicy: 2,
		Timestamp: 100, Now: 105, MaxAge: 10, UpperRaw: 622_592, LimitRaw: 655_360,
	}
	result.GateExamples = append(result.GateExamples, GateExample{
		"pressure_9MPa_plus_0.5MPa_at_10MPa_limit_raw_unit_2^-16_MPa", input, PrevalidatedGate(input),
	})
	// ceil(10.4 * 65536), prepared upstream with outward rounding.
	input.UpperRaw = 681_575
	result.GateExamples = append(result.GateExamples, GateExample{
		"pressure_9.9MPa_plus_0.5MPa_crosses_10MPa_limit_raw_unit_2^-16_MPa", input, PrevalidatedGate(input),
	})
	return result, nil
}
