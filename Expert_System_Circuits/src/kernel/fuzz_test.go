package kernel

import (
	"slices"
	"testing"
	"time"

	evidence "github.com/nick-fedchik/articles/Expert_System_Circuits/src"
)

type byteReader struct {
	data []byte
	pos  int
}

func (r *byteReader) next() int {
	if r.pos >= len(r.data) {
		return 0
	}
	r.pos++
	return int(r.data[r.pos-1])
}

func (r *byteReader) literal(atoms int) evidence.SignedLiteral {
	b := r.next()
	return evidence.SignedLiteral{Atom: (b >> 1) % atoms, Negative: b&1 == 1}
}

func (r *byteReader) sources() []string {
	b := r.next()
	var out []string
	for i, id := range []string{"s1", "s2", "rogue"} {
		if b>>i&1 == 1 {
			out = append(out, id)
		}
	}
	return out
}

func fuzzCase(data []byte) (Spec, Input) {
	r := &byteReader{data: data}
	atoms := 1 + r.next()%4
	spec := Spec{Action: "act", Theory: evidence.NormTheory{Atoms: atoms},
		TrustedSources: []string{"s1", "s2"}, MaxAge: time.Duration(r.next()%4) * time.Second}
	if r.next()&1 == 1 {
		spec.Guards = []Guard{{Atom: 0, Quantity: "q", Unit: "bar", Low: -2, High: 2}}
	}
	for i := r.next() % 3; i > 0; i-- {
		var premises []evidence.SignedLiteral
		for j := r.next() % 3; j > 0; j-- {
			premises = append(premises, r.literal(atoms))
		}
		spec.Theory.Evidence = append(spec.Theory.Evidence, evidence.SignedRule{ID: "e", Version: "1", Premises: premises, Conclusion: r.literal(atoms)})
	}
	ids := []string{"a", "b", "c"}
	norms := 1 + r.next()%3
	for i := 0; i < norms; i++ {
		rule := evidence.NormRule{ID: ids[i], Version: "1", Action: "act", Effect: evidence.Permit}
		if r.next()&1 == 1 {
			rule.Effect = evidence.Prohibit
		}
		for j := r.next() % 3; j > 0; j-- {
			rule.Conditions = append(rule.Conditions, r.literal(atoms))
		}
		if r.next()&1 == 1 {
			exception := r.literal(atoms)
			rule.Exception = &exception
		}
		spec.Theory.Norms = append(spec.Theory.Norms, rule)
	}
	for i := r.next() % 3; i > 0; i-- {
		spec.Theory.Priorities = append(spec.Theory.Priorities, evidence.NormPriority{Higher: ids[r.next()%norms], Lower: ids[r.next()%norms]})
	}
	const now = int64(1 << 40)
	in := Input{Now: now}
	for i := r.next() % 4; i > 0; i-- {
		in.Facts = append(in.Facts, Fact{Literal: r.literal(atoms), Sources: r.sources(), Time: now + int64(int8(r.next()))*int64(time.Second/2)})
	}
	for i := r.next() % 3; i > 0; i-- {
		unit := "bar"
		if r.next()&1 == 1 {
			unit = "kPa"
		}
		low, high := float64(int8(r.next()))/16, float64(int8(r.next()))/16
		if low > high {
			low, high = high, low
		}
		in.Measurements = append(in.Measurements, Measurement{Quantity: "q", Unit: unit, Low: low, High: high, Sources: r.sources(), Time: now - int64(r.next()%8)*int64(time.Second/2)})
	}
	return spec, in
}

func FuzzKernelSoundness(f *testing.F) {
	f.Add([]byte{})
	f.Add([]byte{3, 2, 1, 1, 0, 0, 0, 2, 0, 1, 2, 0, 1, 2, 1, 2, 0, 0, 0, 3, 3, 1, 3, 1, 0, 0, 16, 32, 1, 0})
	f.Add([]byte("\x02\x01\x01\x00\x01\x00\x02\x00\x01\x00\x00\x01\x01\x02\x01\x00\x03\x01\x00\x10\x20\x01"))
	f.Fuzz(func(t *testing.T, data []byte) {
		spec, in := fuzzCase(data)
		c, err := Evaluate(spec, in)
		if err != nil {
			if Verify(spec, in, c) == nil {
				t.Fatal("zero certificate accepted for rejected case")
			}
			return
		}
		if err := Verify(spec, in, c); err != nil {
			t.Fatal(err)
		}
		if violated := CheckSafety(spec, in); c.Verdict == Allow && len(violated) != 0 {
			t.Fatalf("ALLOW violates P: %v", violated)
		}
		if c.Verdict != Allow {
			forged := cloneCertificate(c)
			forged.Verdict, forged.Norm.Result.Decision = Allow, Allow
			forged.Fresh, forged.Verified, forged.Contained, forged.Violations = true, true, true, nil
			forged.Clean = slices.Clone(forged.Norm.Closure)
			if Verify(spec, in, forged) == nil {
				t.Fatal("forged ALLOW accepted")
			}
		}
		if len(data) > 0 && len(c.Norm.Closure) > 0 {
			forged := cloneCertificate(c)
			i := int(data[len(data)-1]) % len(forged.Norm.Closure)
			forged.Norm.Closure[i] ^= evidence.State(1 + data[0]%3)
			if Verify(spec, in, forged) == nil {
				t.Fatal("altered closure accepted")
			}
		}
	})
}
