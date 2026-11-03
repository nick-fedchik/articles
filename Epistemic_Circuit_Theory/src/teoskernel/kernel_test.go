package teoskernel

import (
	"math"
	"slices"
	"testing"
	"time"

	evidence "github.com/nick-fedchik/articles/Epistemic_Circuit_Theory/src"
	"github.com/nick-fedchik/articles/Epistemic_Circuit_Theory/src/teoscompose"
)

func cloneCertificate(c Certificate) Certificate {
	c.Atoms = slices.Clone(c.Atoms)
	c.Guards = slices.Clone(c.Guards)
	c.Clean = slices.Clone(c.Clean)
	c.Violations = slices.Clone(c.Violations)
	c.Norm.Closure = slices.Clone(c.Norm.Closure)
	c.Norm.Result.Rules = slices.Clone(c.Norm.Result.Rules)
	for i := range c.Norm.Result.Rules {
		c.Norm.Result.Rules[i].DominatedBy = slices.Clone(c.Norm.Result.Rules[i].DominatedBy)
	}
	return c
}

func cloneInput(in Input) Input {
	in.Facts = slices.Clone(in.Facts)
	in.Measurements = slices.Clone(in.Measurements)
	for i := range in.Facts {
		in.Facts[i].Sources = slices.Clone(in.Facts[i].Sources)
	}
	for i := range in.Measurements {
		in.Measurements[i].Sources = slices.Clone(in.Measurements[i].Sources)
	}
	return in
}

func TestExampleAllowsAndSatisfiesP(t *testing.T) {
	spec, input := Example()
	c, err := Evaluate(spec, input)
	if err != nil {
		t.Fatal(err)
	}
	if c.Verdict != Allow || len(c.Violations) != 0 ||
		!slices.Equal(c.Norm.Closure, []evidence.State{evidence.True, evidence.True, evidence.False, evidence.False}) ||
		!slices.Equal(c.Clean, []evidence.State{evidence.True, evidence.True, evidence.False, evidence.False}) {
		t.Fatalf("unexpected certificate %+v", c)
	}
	if err := Verify(spec, input, c); err != nil {
		t.Fatal(err)
	}
	if violated := CheckSafety(spec, input); len(violated) != 0 {
		t.Fatal(violated)
	}
}

func TestUnreferencedNoiseFactDoesNotBlockAllow(t *testing.T) {
	spec, input := Example()
	// Add atom 4, which is not referenced by any guard, evidence rule, or norm.
	spec.Theory.Atoms = 5
	input.Facts = append(input.Facts,
		Fact{
			Literal: evidence.SignedLiteral{Atom: 4},
			Sources: []string{"untrusted-spam"},
			Time:    input.Now - int64(time.Hour),
		},
		Fact{
			Literal: evidence.SignedLiteral{Atom: 4, Negative: true},
			Sources: nil,
			Time:    input.Now + int64(time.Hour),
		},
	)
	c, err := Evaluate(spec, input)
	if err != nil {
		t.Fatal(err)
	}
	if c.Verdict != Allow || !c.Fresh || !c.Verified || len(c.Violations) != 0 {
		t.Fatalf("unreferenced noise fact blocked ALLOW: verdict=%s fresh=%v verified=%v violations=%v",
			c.Verdict, c.Fresh, c.Verified, c.Violations)
	}
	if err := Verify(spec, input, c); err != nil {
		t.Fatal(err)
	}
	if violated := CheckSafety(spec, input); len(violated) != 0 {
		t.Fatalf("expected P to hold with unreferenced noise fact, got %v", violated)
	}
}

func TestCleanProofBlocksTaintedDerivation(t *testing.T) {
	const now = int64(1_793_700_000_000_000_000)
	// Atom 0: guard, Atom 1: contested sensor (BOTH), Atom 2: derived safety flag (1+ -> 2+),
	// Atom 3: independent backup sensor.
	spec := Spec{
		Action: "open_valve",
		Theory: evidence.NormTheory{
			Atoms: 4,
			Evidence: []evidence.SignedRule{
				{ID: "primary_derives_safe", Version: "1",
					Premises: []evidence.SignedLiteral{{Atom: 1}}, Conclusion: evidence.SignedLiteral{Atom: 2}},
			},
			Norms: []evidence.NormRule{
				{ID: "permit_if_safe", Version: "1", Action: "open_valve", Effect: evidence.Permit,
					Conditions: []evidence.SignedLiteral{{Atom: 0}, {Atom: 2}}},
			},
		},
		Guards:         []Guard{{Atom: 0, Quantity: "pressure", Unit: "bar", Low: 0, High: 8}},
		TrustedSources: []string{"plc-a", "plc-b"},
		MaxAge:         2 * time.Second,
	}
	input := Input{
		Now: now,
		Facts: []Fact{
			{Literal: evidence.SignedLiteral{Atom: 1}, Sources: []string{"plc-a"}, Time: now},
			{Literal: evidence.SignedLiteral{Atom: 1, Negative: true}, Sources: []string{"plc-b"}, Time: now},
		},
		Measurements: []Measurement{
			{Quantity: "pressure", Unit: "bar", Low: 5, High: 6, Sources: []string{"plc-a"}, Time: now},
		},
	}
	c, err := Evaluate(spec, input)
	if err != nil {
		t.Fatal(err)
	}
	// Closure[1] is Both and Closure[2] is True, so Norm.Result.Decision is ALLOW,
	// but Clean[2] is Unknown, so K must downgrade Verdict to DEFER_CONFLICT!
	if c.Norm.Result.Decision != Allow {
		t.Fatalf("expected raw normative decision ALLOW, got %s", c.Norm.Result.Decision)
	}
	if c.Clean[2] != evidence.Unknown {
		t.Fatalf("expected Clean[2] == Unknown on tainted proof, got %v", c.Clean[2])
	}
	if c.Verdict != Conflict || !slices.Contains(c.Violations, codeConflictDependent) {
		t.Fatalf("expected DEFER_CONFLICT with %s, got %s (%v)", codeConflictDependent, c.Verdict, c.Violations)
	}
	if err := Verify(spec, input, c); err != nil {
		t.Fatal(err)
	}
	violated := CheckSafety(spec, input)
	if !slices.Equal(violated, []string{ClausePositive}) {
		t.Fatalf("expected CheckSafety violation [%s], got %v", ClausePositive, violated)
	}

	// Adding an independent clean derivation (3+ -> 2+ with Atom 3 = True) restores Clean[2] = True and ALLOW.
	spec.Theory.Evidence = append(spec.Theory.Evidence, evidence.SignedRule{
		ID: "backup_derives_safe", Version: "1",
		Premises: []evidence.SignedLiteral{{Atom: 3}}, Conclusion: evidence.SignedLiteral{Atom: 2},
	})
	input.Facts = append(input.Facts, Fact{
		Literal: evidence.SignedLiteral{Atom: 3}, Sources: []string{"plc-a"}, Time: now,
	})
	cClean, err := Evaluate(spec, input)
	if err != nil {
		t.Fatal(err)
	}
	if cClean.Clean[2] != evidence.True || cClean.Verdict != Allow || len(cClean.Violations) != 0 {
		t.Fatalf("expected clean alternative proof to restore ALLOW, got Clean[2]=%v Verdict=%s Violations=%v",
			cClean.Clean[2], cClean.Verdict, cClean.Violations)
	}
	if violated := CheckSafety(spec, input); len(violated) != 0 {
		t.Fatalf("expected P to hold with clean alternative proof, got %v", violated)
	}
}

func TestPostconditionsAndPolicyVerdicts(t *testing.T) {
	cases := []struct {
		name   string
		mutate func(*Spec, *Input)
		want   evidence.NormDecision
	}{
		{"stale fact", func(_ *Spec, in *Input) { in.Facts[0].Time = in.Now - int64(3*time.Second) }, Unknown},
		{"age exactly at bound", func(_ *Spec, in *Input) { in.Facts[0].Time = in.Now - int64(2*time.Second) }, Allow},
		{"future timestamp", func(_ *Spec, in *Input) { in.Facts[0].Time = in.Now + 1 }, Unknown},
		{"untrusted source", func(_ *Spec, in *Input) { in.Facts[0].Sources = []string{"chat-log"} }, Unknown},
		{"partly untrusted", func(_ *Spec, in *Input) { in.Measurements[0].Sources = []string{"plc-a", "x"} }, Unknown},
		{"no sources", func(_ *Spec, in *Input) { in.Facts[2].Sources = nil }, Unknown},
		{"unit mismatch", func(_ *Spec, in *Input) { in.Measurements[0].Unit = "kPa" }, Unknown},
		{"boundary overlap", func(_ *Spec, in *Input) { in.Measurements[0].Low, in.Measurements[0].High = 7.9, 8.1 }, Conflict},
		{"disjoint", func(_ *Spec, in *Input) { in.Measurements[0].Low, in.Measurements[0].High = 9, 9.5 }, Deny},
		{"closed bounds", func(_ *Spec, in *Input) { in.Measurements[0].Low, in.Measurements[0].High = 0, 8 }, Allow},
		{"missing measurement", func(_ *Spec, in *Input) { in.Measurements = nil }, Unknown},
		{"extra contained measurement", func(_ *Spec, in *Input) {
			in.Measurements = append(in.Measurements, Measurement{"pressure", "bar", 1, 2, []string{"plc-a"}, in.Now})
		}, Allow},
		{"second measurement disjoint", func(_ *Spec, in *Input) {
			in.Measurements = append(in.Measurements, Measurement{"pressure", "bar", 9, 10, []string{"plc-a"}, in.Now})
		}, Conflict},
		{"lock engaged", func(_ *Spec, in *Input) { in.Facts[1].Literal.Negative = false; in.Facts = in.Facts[:2] }, Deny},
		{"lock engaged, interlock denied", func(_ *Spec, in *Input) { in.Facts[1].Literal.Negative = false }, Conflict},
		{"interlock unknown", func(_ *Spec, in *Input) { in.Facts = in.Facts[:2] }, Unknown},
		{"authorization contested", func(_ *Spec, in *Input) {
			in.Facts = append(in.Facts, Fact{evidence.SignedLiteral{Atom: 1, Negative: true}, []string{"badge-reader"}, in.Now})
		}, Conflict},
		{"guard outside norm conditions, disjoint", func(s *Spec, in *Input) {
			s.Theory.Norms[0].Conditions = []evidence.SignedLiteral{{Atom: 1}}
			in.Measurements[0].Low, in.Measurements[0].High = 9, 9.5
		}, Deny},
		{"guard outside norm conditions, boundary", func(s *Spec, in *Input) {
			s.Theory.Norms[0].Conditions = []evidence.SignedLiteral{{Atom: 1}}
			in.Measurements[0].Low, in.Measurements[0].High = 7, 9
		}, Conflict},
		{"guard outside norm conditions, missing", func(s *Spec, in *Input) {
			s.Theory.Norms[0].Conditions = []evidence.SignedLiteral{{Atom: 1}}
			in.Measurements = nil
		}, Unknown},
		{"no applicable norm", func(s *Spec, in *Input) {
			s.Theory.Norms, s.Theory.Priorities = s.Theory.Norms[:1], nil
			in.Facts[0].Literal.Negative = true
		}, NoBasis},
	}
	for _, tc := range cases {
		spec, input := Example()
		tc.mutate(&spec, &input)
		c, err := Evaluate(spec, input)
		if err != nil {
			t.Fatalf("%s: %v", tc.name, err)
		}
		if c.Verdict != tc.want {
			t.Fatalf("%s: verdict %s, want %s (%v)", tc.name, c.Verdict, tc.want, c.Violations)
		}
		if err := Verify(spec, input, c); err != nil {
			t.Fatalf("%s: %v", tc.name, err)
		}
		if violated := CheckSafety(spec, input); c.Verdict == Allow && len(violated) != 0 {
			t.Fatalf("%s: ALLOW with violated P %v", tc.name, violated)
		}
	}
}

func TestVerifierRejectsForgedCertificates(t *testing.T) {
	spec, input := Example()
	honest, err := Evaluate(spec, input)
	if err != nil {
		t.Fatal(err)
	}
	var mutations []func(*Certificate)
	for i := range honest.Atoms {
		for s := evidence.Unknown; s <= evidence.Both; s++ {
			if s != honest.Atoms[i] {
				mutations = append(mutations, func(c *Certificate) { c.Atoms[i] = s })
			}
			if s != honest.Norm.Closure[i] {
				mutations = append(mutations, func(c *Certificate) { c.Norm.Closure[i] = s })
			}
			if s != honest.Clean[i] {
				mutations = append(mutations, func(c *Certificate) { c.Clean[i] = s })
			}
		}
	}
	for _, d := range []evidence.NormDecision{Deny, Unknown, Conflict, NoBasis} {
		mutations = append(mutations,
			func(c *Certificate) { c.Verdict = d },
			func(c *Certificate) { c.Norm.Result.Decision = d },
			func(c *Certificate) { c.Verdict, c.Norm.Result.Decision = d, d })
	}
	mutations = append(mutations,
		func(c *Certificate) { c.Guards[0] = evidence.Both },
		func(c *Certificate) { c.Fresh = false },
		func(c *Certificate) { c.Verified = false },
		func(c *Certificate) { c.Contained = false },
		func(c *Certificate) { c.Violations = append(c.Violations, "STALE:fact:0") },
		func(c *Certificate) { c.SpecDigest[0] ^= 1 },
		func(c *Certificate) { c.InputDigest[31] ^= 1 },
		func(c *Certificate) { c.Norm.Theory[0] ^= 1 },
		func(c *Certificate) { c.Norm.Context[0] ^= 1 },
		func(c *Certificate) { c.Norm.Result.Action = "close_valve" },
		func(c *Certificate) { c.Norm.Result.Rules[0].Status = "APPLICABLE" },
		func(c *Certificate) { c.Norm.Result.Rules[1].DominatedBy = nil },
		func(c *Certificate) { c.Norm.Result.Surviving = evidence.Both },
		func(c *Certificate) { c.Norm.Result.Applicable = evidence.True },
		func(c *Certificate) { c.Atoms = c.Atoms[:3] },
		func(c *Certificate) { c.Norm.Closure = append(c.Norm.Closure, evidence.True) },
	)
	for i, mutate := range mutations {
		forged := cloneCertificate(honest)
		mutate(&forged)
		if Verify(spec, input, forged) == nil {
			t.Fatalf("mutation %d accepted: %+v", i, forged)
		}
	}
	if err := Verify(spec, input, honest); err != nil {
		t.Fatal("mutations aliased the honest certificate:", err)
	}
	t.Logf("%d certificate mutations rejected", len(mutations))

	stale := cloneInput(input)
	stale.Facts[0].Time -= int64(time.Hour)
	if Verify(spec, stale, honest) == nil {
		t.Fatal("certificate replayed on a stale input")
	}
	looser := spec
	looser.MaxAge = time.Hour
	if Verify(looser, input, honest) == nil {
		t.Fatal("certificate replayed under another spec")
	}
}

func TestOnlyTheLeastClosureIsAccepted(t *testing.T) {
	spec, input := Example()
	honest, err := Evaluate(spec, input)
	if err != nil {
		t.Fatal(err)
	}
	accepted := 0
	for code := 0; code < 256; code++ {
		forged := cloneCertificate(honest)
		for i := range forged.Norm.Closure {
			forged.Norm.Closure[i] = evidence.State(code >> (2 * i) & 3)
		}
		if Verify(spec, input, forged) == nil {
			accepted++
			if !slices.Equal(forged.Norm.Closure, honest.Norm.Closure) {
				t.Fatalf("non-least closure accepted: %v", forged.Norm.Closure)
			}
		}
	}
	if accepted != 1 {
		t.Fatalf("accepted %d closures", accepted)
	}
}

func TestGuardRestatementMatchesIntervalGuard(t *testing.T) {
	values := []float64{-1, 0, math.Nextafter(0, 1), 2, 8, math.Nextafter(8, 9), 9}
	for _, unit := range []string{"bar", "kPa"} {
		for _, low := range values {
			for _, high := range values {
				if low > high {
					continue
				}
				m, err := teoscompose.NewMeasurement(low, evidence.Interval{Low: low, High: high}, unit)
				if err != nil {
					t.Fatal(err)
				}
				want, err := teoscompose.IntervalGuard(m, evidence.Interval{Low: 0, High: 8}, "bar")
				got := guardState(Guard{Unit: "bar", Low: 0, High: 8}, Measurement{Unit: unit, Low: low, High: high})
				if err != nil || got != want {
					t.Fatalf("[%v,%v] %s: %v != %v (%v)", low, high, unit, got, want, err)
				}
			}
		}
	}
}

func TestFreshnessArithmeticAtExtremes(t *testing.T) {
	values := []int64{math.MinInt64, math.MinInt64 + 1, -1, 0, 1, math.MaxInt64 - 1, math.MaxInt64}
	ages := []time.Duration{0, 1, time.Second, math.MaxInt64}
	for _, now := range values {
		for _, ts := range values {
			for _, age := range ages {
				if fresh(now, ts, age) != freshBig(now, ts, age) {
					t.Fatalf("now=%d t=%d age=%d", now, ts, age)
				}
			}
		}
	}
}

func TestMalformedSpecsAndInputsRejected(t *testing.T) {
	cases := []func(*Spec, *Input){
		func(s *Spec, _ *Input) { s.Action = " " },
		func(s *Spec, _ *Input) { s.MaxAge = -1 },
		func(s *Spec, _ *Input) { s.Theory.Norms[2].Action = "close_valve" },
		func(s *Spec, _ *Input) { s.TrustedSources = append(s.TrustedSources, "") },
		func(s *Spec, _ *Input) { s.Guards[0].High = math.NaN() },
		func(s *Spec, _ *Input) { s.Guards[0].Low = 9 },
		func(s *Spec, _ *Input) { s.Guards = append(s.Guards, Guard{1, "pressure", "bar", 0, 1}) },
		func(s *Spec, _ *Input) { s.Guards[0].Atom = 4 },
		func(s *Spec, _ *Input) {
			s.Theory.Evidence = append(s.Theory.Evidence, evidence.SignedRule{ID: "g", Version: "1",
				Premises: []evidence.SignedLiteral{{Atom: 1}}, Conclusion: evidence.SignedLiteral{Atom: 0}})
		},
		func(s *Spec, _ *Input) {
			s.Theory.Evidence = append(s.Theory.Evidence, evidence.SignedRule{ID: "cycle", Version: "1",
				Premises: []evidence.SignedLiteral{{Atom: 3}}, Conclusion: evidence.SignedLiteral{Atom: 2}})
		},
		func(s *Spec, _ *Input) {
			s.Theory.Priorities = append(s.Theory.Priorities, evidence.NormPriority{Higher: "default_deny", Lower: "ban_under_lock"})
		},
		func(_ *Spec, in *Input) { in.Facts[0].Literal.Atom = 0 },
		func(_ *Spec, in *Input) { in.Facts[0].Literal.Atom = 9 },
		func(_ *Spec, in *Input) { in.Measurements[0].Quantity = "temperature" },
		func(_ *Spec, in *Input) { in.Measurements[0].High = math.Inf(1) },
		func(_ *Spec, in *Input) { in.Measurements[0].Low = 7 },
		func(_ *Spec, in *Input) { in.Measurements[0].Unit = "" },
	}
	for i, mutate := range cases {
		spec, input := Example()
		honest, err := Evaluate(spec, input)
		if err != nil {
			t.Fatal(err)
		}
		mutate(&spec, &input)
		if _, err := Evaluate(spec, input); err == nil {
			t.Fatalf("case %d evaluated", i)
		}
		if Verify(spec, input, honest) == nil {
			t.Fatalf("case %d verified", i)
		}
		if v := CheckSafety(spec, input); len(v) != 1 || v[0] != ClauseWellFormed {
			t.Fatalf("case %d safety %v", i, v)
		}
	}
}

func TestBoundedExhaustiveSoundness(t *testing.T) {
	bounds := Bounds{MaxNorms: 1, MaxConditions: 1, Exceptions: false}
	if testing.Short() {
		bounds.MaxConditions = 0
	}
	report, err := Exhaustive(bounds)
	if err != nil {
		t.Fatal(err)
	}
	if want := report.Theories * len(Inputs()); report.Cases != want {
		t.Fatalf("cases %d, want %d", report.Cases, want)
	}
	if len(report.Counterexamples) != 0 || report.ForgeriesAttempted != report.ForgeriesRejected ||
		report.Verdicts[Allow] != report.AllowWithSafetyOK || report.AllowWithSafetyOK == 0 ||
		report.Cases-report.PFalseCases != report.AllowWithSafetyOK {
		t.Fatalf("%+v", report)
	}
	t.Logf("%+v", report)
}

// The normative layer alone does not establish P: when no norm mentions the
// guard atom, it returns ALLOW for an out-of-bounds or missing measurement.
// Only the admission postconditions of K close this gap.
func TestNormativeAllowAloneIsInsufficient(t *testing.T) {
	for _, measurements := range [][]Measurement{nil, {{"pressure", "bar", 9, 9.5, []string{"plc-a"}, 0}}} {
		spec, input := Example()
		spec.Theory.Norms[0].Conditions = []evidence.SignedLiteral{{Atom: 1}}
		for i := range measurements {
			measurements[i].Time = input.Now
		}
		input.Measurements = measurements
		c, err := Evaluate(spec, input)
		if err != nil {
			t.Fatal(err)
		}
		if c.Norm.Result.Decision != Allow || c.Verdict == Allow || len(CheckSafety(spec, input)) == 0 {
			t.Fatalf("normative %s, kernel %s", c.Norm.Result.Decision, c.Verdict)
		}
	}
}
