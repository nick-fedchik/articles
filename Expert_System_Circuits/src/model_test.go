package evidence

import (
	"crypto/sha256"
	"math"
	"slices"
	"testing"
	"time"
)

func exampleRule() Rule {
	return Rule{Roles{"controller", "opens", "valve"}, "pressure:MPa", Interval{10, 15}, 20 * time.Millisecond, True}
}

func exampleFact() Fact {
	return Fact{Roles{"controller", "opens", "valve"}, "pressure:MPa", Interval{12, 13}, time.Millisecond}
}

func TestFourStatesAndDistinctConjunctions(t *testing.T) {
	for a := Unknown; a <= Both; a++ {
		for b := Unknown; b <= Both; b++ {
			got, err := TruthAnd(a, b)
			if err != nil {
				t.Fatal(err)
			}
			positive := (a == True || a == Both) && (b == True || b == Both)
			negative := a == False || a == Both || b == False || b == Both
			want := State(0)
			if positive {
				want += 2
			}
			if negative {
				want++
			}
			if got != want {
				t.Fatalf("truth(%v,%v)=%v, want %v", a, b, got, want)
			}
		}
	}
	truth, _ := TruthAnd(Unknown, Both)
	policy, _ := Priority(Unknown, Both)
	if truth != False || policy != Unknown {
		t.Fatal("standard truth conjunction must differ from priority policy")
	}
	if _, err := Priority(State(4), True); err == nil {
		t.Fatal("invalid logical state accepted")
	}
}

func TestRuleBoundariesAndReasons(t *testing.T) {
	cases := []struct {
		name   string
		change func(*Rule, *Fact)
		state  State
		reason Reason
	}{
		{"inside", func(*Rule, *Fact) {}, True, ""},
		{"exact upper", func(_ *Rule, f *Fact) { f.Range = Interval{15, 15} }, True, ""},
		{"touch upper", func(_ *Rule, f *Fact) { f.Range = Interval{15, 16} }, Both, Boundary},
		{"outside", func(_ *Rule, f *Fact) { f.Range = Interval{16, 17} }, False, OutOfRange},
		{"sort", func(_ *Rule, f *Fact) { f.Sort = "count:octet" }, Unknown, SortMismatch},
		{"inversion", func(_ *Rule, f *Fact) { f.Roles.Actor, f.Roles.Target = f.Roles.Target, f.Roles.Actor }, False, RoleInversion},
		{"self relation", func(r *Rule, f *Fact) {
			r.Pattern.Target = r.Pattern.Actor
			f.Roles.Target = f.Roles.Actor
		}, True, ""},
		{"wildcard", func(r *Rule, _ *Fact) { r.Pattern.Actor = "*" }, True, ""},
		{"age edge", func(r *Rule, f *Fact) { f.Age = r.MaxAge }, True, ""},
		{"stale", func(r *Rule, f *Fact) { f.Age = r.MaxAge + 1 }, False, Stale},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rule, fact := exampleRule(), exampleFact()
			tc.change(&rule, &fact)
			got, err := EvaluateRule(rule, fact)
			if err != nil || got.State != tc.state {
				t.Fatalf("%+v, %v; want %v", got, err, tc.state)
			}
			if tc.reason != "" && !slices.Contains(got.Reasons, tc.reason) {
				t.Fatalf("missing diagnostic %s: %+v", tc.reason, got)
			}
		})
	}
}

func TestRejectMalformedRulesAndFacts(t *testing.T) {
	for _, bad := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
		fact := exampleFact()
		fact.Range.Low = bad
		if _, err := EvaluateRule(exampleRule(), fact); err == nil {
			t.Fatal("nonfinite fact accepted")
		}
	}
	fact := exampleFact()
	fact.Age = -1
	if _, err := EvaluateRule(exampleRule(), fact); err == nil {
		t.Fatal("negative age accepted")
	}
	fact = exampleFact()
	fact.Range = Interval{2, 1}
	if _, err := EvaluateRule(exampleRule(), fact); err == nil {
		t.Fatal("inverted interval accepted")
	}
}

func TestQuoteAuthenticationIsOnlyByteAuthentication(t *testing.T) {
	document := []byte("public fixture: pressure <= 15 MPa")
	digest := sha256.Sum256(document[16:])
	if err := VerifyQuote(document, 16, len(document), digest); err != nil {
		t.Fatal(err)
	}
	if err := VerifyQuote(document, 15, len(document), digest); err == nil {
		t.Fatal("wrong byte range accepted")
	}
	if err := VerifyQuote(document, -1, 5, digest); err == nil {
		t.Fatal("invalid byte range accepted")
	}
}
