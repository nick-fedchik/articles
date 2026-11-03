package evidence

import (
	"math"
	"slices"
	"testing"
)

func TestRoutingDoesNotMergeVersionsOrInventLogicalStates(t *testing.T) {
	candidates := []Candidate{{"standard", "1", .9, .02}, {"standard", "2", .89, .02}}
	got, err := SelectCandidates(candidates, .5, .01)
	if err != nil || got.Disposition != "CLARIFY" || len(got.Candidates) != 2 {
		t.Fatalf("%+v %v", got, err)
	}
	slices.Reverse(candidates)
	reversed, err := SelectCandidates(candidates, .5, .01)
	if err != nil || !slices.Equal(got.Candidates, reversed.Candidates) {
		t.Fatal("candidate input order changed routing")
	}
	direct, _ := SelectCandidates([]Candidate{{"a", "1", .9, .01}, {"b", "1", .3, .01}}, .5, .01)
	refusal, _ := SelectCandidates([]Candidate{{"a", "1", .1, .01}}, .5, .01)
	if direct.Disposition != "DIRECT" || refusal.Disposition != "REFUSE" {
		t.Fatalf("direct=%+v refusal=%+v", direct, refusal)
	}
}

func TestMaskedLogitsShiftInvarianceAndFalseConfidenceWitness(t *testing.T) {
	allowed, normative := []bool{true, true, false}, []bool{true, false, false}
	a, err := SummarizeLogits([]float64{4, 1, 0}, allowed, normative, 1)
	if err != nil {
		t.Fatal(err)
	}
	b, err := SummarizeLogits([]float64{1, -2, -3}, allowed, normative, 1)
	if err != nil || math.Abs(a.NormMass-b.NormMass) > 1e-14 || math.Abs(a.Entropy-b.Entropy) > 1e-14 {
		t.Fatalf("logit shift changed diagnostics: %+v %+v %v", a, b, err)
	}
	illusion, err := SummarizeLogits(
		[]float64{math.Log(.0009), math.Log(.0001), math.Log(.999)},
		[]bool{true, true, false}, []bool{true, true, false}, 1)
	if err != nil || math.Abs(illusion.ConditionalMax-.9) > 1e-12 || math.Abs(illusion.NormMass-.001) > 1e-12 {
		t.Fatalf("mask illusion: %+v %v", illusion, err)
	}
	// The true label is NO_ATOM, but the model assigns it only .005.
	wrong, err := SummarizeLogits(
		[]float64{math.Log(.99), math.Log(.005), math.Log(.005)}, allowed, normative, 1)
	if err != nil || math.Abs((1-wrong.NormMass)-.01) > 1e-12 {
		t.Fatal("incorrect high-confidence witness")
	}
	if 1-wrong.NormMass >= .35 {
		t.Fatal("witness should evade the proposed scalar rejection threshold")
	}
}

func TestLogitsRejectMalformedInputs(t *testing.T) {
	for _, value := range []float64{math.NaN(), math.Inf(1)} {
		if _, err := SummarizeLogits([]float64{value, 1}, []bool{true, true}, []bool{true, false}, 1); err == nil {
			t.Fatal("nonfinite logits accepted")
		}
	}
	if _, err := SummarizeLogits([]float64{1}, []bool{false}, []bool{true}, 1); err == nil {
		t.Fatal("invalid class masks accepted")
	}
	if _, err := SelectCandidates([]Candidate{{"a", "1", math.NaN(), 0}}, .5, 0); err == nil {
		t.Fatal("NaN score accepted")
	}
}

func TestZeroFailuresAreNotZeroRisk(t *testing.T) {
	got, err := ZeroFailureUpperBound(100, .05)
	if err != nil || math.Abs(got-.0295130496) > 1e-9 {
		t.Fatalf("got %g, %v", got, err)
	}
}
