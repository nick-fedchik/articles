// Command gatevectors writes RTL test vectors whose expected outputs are
// computed by hardware.PrevalidatedGate. It defines only the bit packing of
// GateInput for the RTL port; it adds no gate policy.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"math"
	"math/big"
	"math/rand/v2"
	"os"
	"path/filepath"

	evidence "github.com/nick-fedchik/articles/Expert_System_Circuits/src"
	"github.com/nick-fedchik/articles/Expert_System_Circuits/src/hardware"
)

// WordBits is the packed GateInput width; ResultBits is allow + reason mask.
const (
	WordBits   = 462
	ResultBits = 10
	VectorBits = WordBits + ResultBits
	hexDigits  = (VectorBits + 3) / 4
)

// Reasons lists the PrevalidatedGate reason strings in check order; reason i
// is result bit i+1, result bit 0 is Allow.
var Reasons = []string{
	"custody_not_validated",
	"validation_binding_invalid",
	"si_mismatch",
	"numeric_invalid",
	"norm_blocked",
	"evidence_not_pro_only",
	"policy_revision_mismatch",
	"timestamp_invalid_or_stale",
	"guard_band_exceeds_limit",
}

type field struct {
	name  string
	lsb   int
	width int
}

// Layout is the RTL word layout (LSB first); it must match gate_core.v.
var Layout = []field{
	{"custody_validated", 0, 1}, {"binding_valid", 1, 1}, {"si_match", 2, 1},
	{"numeric_valid", 3, 1}, {"norm_permitted", 4, 1}, {"undercut", 5, 1},
	{"evidence_code", 6, 8},
	{"data_revision", 14, 32}, {"validated_revision", 46, 32},
	{"policy_revision", 78, 32}, {"expected_policy_revision", 110, 32},
	{"timestamp_ticks", 142, 64}, {"now_ticks", 206, 64}, {"max_age_ticks", 270, 64},
	{"prevalidated_upper_bound_raw", 334, 64}, {"limit_raw", 398, 64},
}

func b2u(v bool) uint64 {
	if v {
		return 1
	}
	return 0
}

func values(in hardware.GateInput) []uint64 {
	return []uint64{
		b2u(in.CustodyValidated), b2u(in.BindingValid), b2u(in.SIMatch),
		b2u(in.NumericValid), b2u(in.NormPermitted), b2u(in.Undercut),
		uint64(in.Evidence),
		uint64(in.DataRevision), uint64(in.ValidatedRevision),
		uint64(in.PolicyRevision), uint64(in.ExpectedPolicy),
		in.Timestamp, in.Now, in.MaxAge,
		uint64(in.UpperRaw), uint64(in.LimitRaw), // two's complement bits
	}
}

// Pack returns the 462-bit word; int64 fields keep their two's complement bits.
func Pack(in hardware.GateInput) *big.Int {
	word := new(big.Int)
	for i, v := range values(in) {
		part := new(big.Int).SetUint64(v)
		word.Or(word, part.Lsh(part, uint(Layout[i].lsb)))
	}
	return word
}

// Unpack is the inverse of Pack; it is used to derive inputs from bit flips.
func Unpack(word *big.Int) hardware.GateInput {
	get := func(i int) uint64 {
		f := Layout[i]
		v := new(big.Int).Rsh(word, uint(f.lsb))
		mask := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), uint(f.width)), big.NewInt(1))
		return v.And(v, mask).Uint64()
	}
	return hardware.GateInput{
		CustodyValidated: get(0) == 1, BindingValid: get(1) == 1, SIMatch: get(2) == 1,
		NumericValid: get(3) == 1, NormPermitted: get(4) == 1, Undercut: get(5) == 1,
		Evidence:     uint8(get(6)),
		DataRevision: uint32(get(7)), ValidatedRevision: uint32(get(8)),
		PolicyRevision: uint32(get(9)), ExpectedPolicy: uint32(get(10)),
		Timestamp: get(11), Now: get(12), MaxAge: get(13),
		UpperRaw: int64(get(14)), LimitRaw: int64(get(15)),
	}
}

// Expected maps the Go decision to the RTL result: bit 0 Allow, bit i+1 reason i.
func Expected(in hardware.GateInput) (uint16, error) {
	decision := hardware.PrevalidatedGate(in)
	var result uint16
	if decision.Allow {
		result = 1
	}
	for _, reason := range decision.Reasons {
		index := -1
		for i, known := range Reasons {
			if known == reason {
				index = i
			}
		}
		if index < 0 {
			return 0, fmt.Errorf("unknown gate reason %q", reason)
		}
		result |= 1 << (index + 1)
	}
	if decision.Allow != (len(decision.Reasons) == 0) {
		return 0, fmt.Errorf("inconsistent decision %+v", decision)
	}
	return result, nil
}

// CeilRaw converts a decimal quantity to ceil(x * 2^16), the outward rounding
// of the guard-band upper bound performed upstream of the gate (see report.go).
// The gate itself compares prepared raw integers and performs no rounding.
func CeilRaw(decimal string) (int64, error) {
	x, ok := new(big.Rat).SetString(decimal)
	if !ok {
		return 0, fmt.Errorf("bad decimal %q", decimal)
	}
	x.Mul(x, new(big.Rat).SetInt64(65536))
	q, r := new(big.Int).QuoRem(x.Num(), x.Denom(), new(big.Int))
	if r.Sign() > 0 {
		q.Add(q, big.NewInt(1))
	}
	if !q.IsInt64() {
		return 0, fmt.Errorf("%q overflows int64 raw", decimal)
	}
	return q.Int64(), nil
}

// Base is the passing pressure example from hardware/report.go.
func Base() hardware.GateInput {
	return hardware.GateInput{
		CustodyValidated: true, BindingValid: true, SIMatch: true,
		NumericValid: true, NormPermitted: true, Evidence: 1,
		DataRevision: 7, ValidatedRevision: 7, PolicyRevision: 2, ExpectedPolicy: 2,
		Timestamp: 100, Now: 105, MaxAge: 10, UpperRaw: 622_592, LimitRaw: 655_360,
	}
}

type vector struct {
	class string
	in    hardware.GateInput
}

func exhaustive() ([]vector, error) {
	var codes []uint8
	for _, state := range []evidence.State{evidence.Unknown, evidence.True, evidence.False, evidence.Both} {
		code, err := hardware.EncodeEvidence(state)
		if err != nil {
			return nil, err
		}
		codes = append(codes, code)
	}
	codes = append(codes, 4, 5, 9, 0x81, 0xff) // out-of-range 8-bit codes
	times := [][3]uint64{{100, 105, 10}, {106, 105, 10}, {100, 110, 10}, {100, 111, 10}}
	uppers := []int64{655_359, 655_360, 655_361}
	var out []vector
	for flags := 0; flags < 64; flags++ {
		for _, code := range codes {
			for binding := 0; binding < 2; binding++ {
				for policy := 0; policy < 2; policy++ {
					for _, t := range times {
						for _, upper := range uppers {
							in := Base()
							in.CustodyValidated = flags&1 != 0
							in.BindingValid = flags&2 != 0
							in.SIMatch = flags&4 != 0
							in.NumericValid = flags&8 != 0
							in.NormPermitted = flags&16 != 0
							in.Undercut = flags&32 != 0
							in.Evidence = code
							if binding == 1 {
								in.ValidatedRevision = 8
							}
							if policy == 1 {
								in.ExpectedPolicy = 3
							}
							in.Timestamp, in.Now, in.MaxAge = t[0], t[1], t[2]
							in.UpperRaw = upper
							out = append(out, vector{"exhaustive", in})
						}
					}
				}
			}
		}
	}
	for code := 0; code < 256; code++ {
		in := Base()
		in.Evidence = uint8(code)
		out = append(out, vector{"evidence_all_codes", in})
	}
	return out, nil
}

func boundaries() ([]vector, error) {
	const maxU = math.MaxUint64
	var out []vector
	add := func(class string, edit func(*hardware.GateInput)) {
		in := Base()
		edit(&in)
		out = append(out, vector{class, in})
	}
	timeCases := [][3]uint64{
		{0, 0, 0}, {5, 5, 0}, {5, 6, 0}, {5, 6, 1}, {6, 5, maxU}, {maxU, 0, maxU},
		{1, 0, maxU}, {0, maxU, maxU}, {0, maxU, maxU - 1}, {maxU, maxU, 0},
		{maxU - 1, maxU, 0}, {maxU - 1, maxU, 1}, {1 << 63, (1 << 63) - 1, maxU},
		{(1 << 63) - 1, 1 << 63, 0}, {(1 << 63) - 1, 1 << 63, 1}, {1 << 32, (1 << 33) - 1, (1 << 32) - 1},
		{1 << 32, (1 << 33) - 1, (1 << 32) - 2}, {0, 1 << 63, 1 << 63}, {0, 1 << 63, (1 << 63) - 1},
	}
	for _, t := range timeCases {
		t := t
		add("boundary_timestamp", func(in *hardware.GateInput) { in.Timestamp, in.Now, in.MaxAge = t[0], t[1], t[2] })
	}
	signed := []int64{math.MinInt64, math.MinInt64 + 1, -65536, -2, -1, 0, 1, 65535, 65536, 655_360, math.MaxInt64 - 1, math.MaxInt64}
	for _, u := range signed {
		for _, l := range signed {
			u, l := u, l
			add("boundary_signed_compare", func(in *hardware.GateInput) { in.UpperRaw, in.LimitRaw = u, l })
		}
	}
	for bit := 0; bit < 32; bit++ {
		bit := bit
		add("boundary_revision_bit", func(in *hardware.GateInput) { in.ValidatedRevision ^= 1 << bit })
		add("boundary_revision_bit", func(in *hardware.GateInput) { in.ExpectedPolicy ^= 1 << bit })
	}
	for _, r := range []uint32{0, 1, 1 << 31, math.MaxUint32} {
		r := r
		add("boundary_revision_extreme", func(in *hardware.GateInput) {
			in.DataRevision, in.ValidatedRevision, in.PolicyRevision, in.ExpectedPolicy = r, r, r, r
		})
	}
	// Outward-rounded upper bounds (value + half-width) against a 10 MPa limit.
	rounding := []string{
		"9.5", "10.4", "10", "9.9999847412109375", "9.99999", "10.0000000001",
		"10.0000076293945", "10.0000152587890625", "10.00001525878907",
		"-0.5", "-0.0000001", "0", "0.0000001", "32767.9999847412109375", "-32768",
	}
	limit, err := CeilRaw("10")
	if err != nil {
		return nil, err
	}
	for _, decimal := range rounding {
		upper, err := CeilRaw(decimal)
		if err != nil {
			return nil, err
		}
		add("rounding_edge", func(in *hardware.GateInput) { in.UpperRaw, in.LimitRaw = upper, limit })
		add("rounding_edge_negative_limit", func(in *hardware.GateInput) { in.UpperRaw, in.LimitRaw = upper, -limit })
	}
	base := Pack(Base())
	for bit := 0; bit < WordBits; bit++ {
		word := new(big.Int).Set(base)
		word.SetBit(word, bit, word.Bit(bit)^1)
		out = append(out, vector{"single_bit_flip", Unpack(word)})
	}
	return out, nil
}

func randomVectors(seed uint64, count int) []vector {
	rng := rand.New(rand.NewPCG(seed, 0x7465_6f73_6761_7465))
	var out []vector
	for i := 0; i < count; i++ {
		in := hardware.GateInput{
			CustodyValidated: rng.IntN(2) == 1, BindingValid: rng.IntN(2) == 1,
			SIMatch: rng.IntN(2) == 1, NumericValid: rng.IntN(2) == 1,
			NormPermitted: rng.IntN(2) == 1, Undercut: rng.IntN(2) == 1,
			Evidence: uint8(rng.Uint32()), DataRevision: rng.Uint32(), ValidatedRevision: rng.Uint32(),
			PolicyRevision: rng.Uint32(), ExpectedPolicy: rng.Uint32(),
			Timestamp: rng.Uint64(), Now: rng.Uint64(), MaxAge: rng.Uint64(),
			UpperRaw: int64(rng.Uint64()), LimitRaw: int64(rng.Uint64()),
		}
		out = append(out, vector{"random_uniform", in})
		// Near-pass: satisfy every check with random magnitudes, then perturb.
		near := in
		near.CustodyValidated, near.BindingValid, near.SIMatch = true, true, true
		near.NumericValid, near.NormPermitted, near.Undercut = true, true, false
		near.Evidence, near.ValidatedRevision, near.ExpectedPolicy = 1, near.DataRevision, near.PolicyRevision
		if near.Timestamp > near.Now {
			near.Timestamp, near.Now = near.Now, near.Timestamp
		}
		age := near.Now - near.Timestamp
		near.MaxAge = age + uint64(rng.IntN(3)) - 1 // age-1, age or age+1, wrapping at extremes
		if near.UpperRaw > near.LimitRaw {
			near.UpperRaw, near.LimitRaw = near.LimitRaw, near.UpperRaw
		}
		switch rng.IntN(4) {
		case 0:
			near.UpperRaw = near.LimitRaw
		case 1:
			if near.LimitRaw != math.MaxInt64 {
				near.UpperRaw = near.LimitRaw + 1
			}
		case 2:
			near.Timestamp = near.Now
			near.MaxAge = uint64(rng.IntN(2))
		}
		if rng.IntN(3) == 0 {
			near.ExpectedPolicy ^= 1 << rng.IntN(32)
		}
		out = append(out, vector{"random_near_pass", near})
		// Small magnitudes exercise carries and sign handling near equality.
		narrow := near
		narrow.Timestamp = uint64(rng.IntN(16))
		narrow.Now = uint64(rng.IntN(16))
		narrow.MaxAge = uint64(rng.IntN(16))
		narrow.UpperRaw = int64(rng.IntN(16)) - 8
		narrow.LimitRaw = int64(rng.IntN(16)) - 8
		out = append(out, vector{"random_narrow", narrow})
	}
	return out
}

// Vectors returns all classes in deterministic order.
func Vectors(seed uint64, randomCount int) ([]vector, error) {
	ex, err := exhaustive()
	if err != nil {
		return nil, err
	}
	bd, err := boundaries()
	if err != nil {
		return nil, err
	}
	return append(append(ex, bd...), randomVectors(seed, randomCount)...), nil
}

// Line formats one vector as {expected[9:0], word[461:0]} in hex for $readmemh.
func Line(in hardware.GateInput) (string, uint16, error) {
	expected, err := Expected(in)
	if err != nil {
		return "", 0, err
	}
	word := Pack(in)
	if word.BitLen() > WordBits {
		return "", 0, fmt.Errorf("word exceeds %d bits", WordBits)
	}
	full := new(big.Int).Lsh(new(big.Int).SetUint64(uint64(expected)), WordBits)
	full.Or(full, word)
	return fmt.Sprintf("%0*x", hexDigits, full), expected, nil
}

// render writes one hex line per vector and a JSON summary; it returns the count.
func render(vectors []vector, hexOut, summary io.Writer) (int, error) {
	buffered := bufio.NewWriter(hexOut)
	total, allow := 0, 0
	classes := map[string]int{}
	var reasonCounts [9]int
	for _, v := range vectors {
		line, expected, err := Line(v.in)
		if err != nil {
			return 0, err
		}
		fmt.Fprintln(buffered, line)
		total++
		classes[v.class]++
		if expected&1 == 1 {
			allow++
		}
		for i := range Reasons {
			if expected>>(i+1)&1 == 1 {
				reasonCounts[i]++
			}
		}
	}
	if err := buffered.Flush(); err != nil {
		return 0, err
	}
	fmt.Fprintf(summary, "{\n  \"vectors\": %d,\n  \"allow\": %d,\n  \"refuse\": %d,\n  \"classes\": {", total, allow, total-allow)
	names := []string{"exhaustive", "evidence_all_codes", "boundary_timestamp", "boundary_signed_compare",
		"boundary_revision_bit", "boundary_revision_extreme", "rounding_edge", "rounding_edge_negative_limit",
		"single_bit_flip", "random_uniform", "random_near_pass", "random_narrow"}
	listed := 0
	for i, name := range names {
		sep := ","
		if i == len(names)-1 {
			sep = ""
		}
		fmt.Fprintf(summary, "\n    %q: %d%s", name, classes[name], sep)
		listed += classes[name]
	}
	if listed != total {
		return 0, fmt.Errorf("class accounting mismatch: %d != %d", listed, total)
	}
	fmt.Fprintf(summary, "\n  },\n  \"reason_counts\": {")
	for i, r := range Reasons {
		sep := ","
		if i == len(Reasons)-1 {
			sep = ""
		}
		fmt.Fprintf(summary, "\n    %q: %d%s", r, reasonCounts[i], sep)
	}
	fmt.Fprintf(summary, "\n  }\n}\n")
	return total, nil
}

func write(dir string, vectors []vector, summary io.Writer) error {
	file, err := os.Create(filepath.Join(dir, "vectors.hex"))
	if err != nil {
		return err
	}
	total, err := render(vectors, file, summary)
	if closeErr := file.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, "vectors.count"), []byte(fmt.Sprintf("%d\n", total)), 0o644)
}

func run(out, diagnostics io.Writer, args []string) error {
	flags := flag.NewFlagSet("gatevectors", flag.ContinueOnError)
	flags.SetOutput(diagnostics)
	dir := flags.String("out", ".", "output directory for vectors.hex and vectors.count")
	seed := flags.Uint64("seed", 20261103, "PCG seed for random vectors")
	count := flags.Int("random", 10000, "random iterations (3 vectors each)")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 0 || *count < 0 {
		return fmt.Errorf("unexpected arguments or negative -random")
	}
	vectors, err := Vectors(*seed, *count)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(*dir, 0o755); err != nil {
		return err
	}
	return write(*dir, vectors, out)
}

func main() {
	if err := run(os.Stdout, os.Stderr, os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
