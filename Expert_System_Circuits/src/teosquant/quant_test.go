package teosquant

import (
	"bytes"
	"encoding/json"
	"math"
	"math/big"
	"testing"
)

func closeTo(t *testing.T, got, want float64) {
	t.Helper()
	if math.IsNaN(got) || math.IsInf(got, 0) || math.Abs(got-want) > 1e-12*math.Max(1, math.Abs(want)) {
		t.Fatalf("got %.17g, want %.17g", got, want)
	}
}

func TestExactGeneratorIdentityAndCounterexample(t *testing.T) {
	a, b := big.NewRat(4, 5), big.NewRat(4, 5)
	one := big.NewRat(1, 1)
	generator := func(w *big.Rat) *big.Rat {
		return new(big.Rat).Sub(new(big.Rat).Inv(w), one)
	}
	product := new(big.Rat).Mul(a, b)
	denominator := new(big.Rat).Sub(new(big.Rat).Add(a, b), product)
	h := new(big.Rat).Quo(product, denominator)
	sum := new(big.Rat).Add(generator(a), generator(b))
	inverse := new(big.Rat).Inv(new(big.Rat).Add(one, sum))
	if h.Cmp(inverse) != 0 || h.Cmp(big.NewRat(2, 3)) != 0 {
		t.Fatalf("generator identity failed: %s vs %s", h, inverse)
	}
	if product.Cmp(big.NewRat(16, 25)) != 0 || h.Cmp(product) == 0 {
		t.Fatal("Hamacher must differ from independent conjunction")
	}
	if generator(product).Cmp(big.NewRat(9, 16)) != 0 {
		t.Fatal("sequential equivalent D must include the cross term")
	}
	t.Log("exact H(4/5,4/5)=2/3; independent conjunction=16/25; sequential D=9/16")
}

func TestCompositionGrid(t *testing.T) {
	for i := 1; i <= 32; i++ {
		a := float64(i) / 32
		ga, err := Generator(a)
		if err != nil {
			t.Fatal(err)
		}
		for j := 1; j <= 32; j++ {
			b := float64(j) / 32
			gb, _ := Generator(b)
			h, err := Hamacher(a, b)
			if err != nil {
				t.Fatal(err)
			}
			want, _ := Weight(ga + gb)
			closeTo(t, h, want)
			reversed, _ := Hamacher(b, a)
			closeTo(t, h, reversed)
			if h > math.Min(a, b)+1e-15 {
				t.Fatal("composition increased a rating")
			}
			hc, _ := Hamacher(h, 0.7)
			bc, _ := Hamacher(b, 0.7)
			ah, _ := Hamacher(a, bc)
			closeTo(t, hc, ah)
		}
	}
	for _, a := range []float64{0, 0.2, 1} {
		h, _ := Hamacher(a, 0)
		closeTo(t, h, 0)
		h, _ = Hamacher(a, 1)
		closeTo(t, h, a)
	}
	h, err := Hamacher(1e-250, 1e-250)
	if err != nil || math.Abs(h/1e-250-0.5) > 1e-15 {
		t.Fatalf("small inputs underflowed: %g, %v", h, err)
	}
}

func TestCorrelatedMeasurement(t *testing.T) {
	for _, tc := range []struct{ rho, variance float64 }{{0, 2}, {0.8, 3.6}, {1, 4}, {-1, 0}} {
		v, err := SumVariance(1, 1, tc.rho)
		if err != nil {
			t.Fatal(err)
		}
		closeTo(t, v, tc.variance)
	}
	// Common-mode uncertainty cancels only for equal sensitivities/uncertainties.
	v, _ := SumVariance(2, 1, -1)
	closeTo(t, v, 1)
	d, err := MetrologyRating(2, 10)
	if err != nil {
		t.Fatal(err)
	}
	closeTo(t, d, 0.04)
	t.Log("Var(X+Y)=2 kPa^2 for rho=0; 3.6 kPa^2 for rho=0.8")
}

func TestLinearExpectedErrors(t *testing.T) {
	a, _ := ExpectedErrors(100, 0.01)
	b, _ := ExpectedErrors(200, 0.01)
	c, _ := ExpectedErrors(400, 0.01)
	closeTo(t, a, 1)
	closeTo(t, b, 2)
	closeTo(t, c, 4)
	if b == 4*a {
		t.Fatal("expected errors incorrectly scale quadratically")
	}
	t.Log("N=100,200,400; assumed p=0.01; expected errors=1,2,4 (not observed counts)")
}

func TestChapterExamplesAndUnits(t *testing.T) {
	e, err := SyntheticEvidence()
	if err != nil {
		t.Fatal(err)
	}
	a, b := e.Valve[0], e.Valve[1]
	closeTo(t, a.LoadWeight, 25.0/26)
	closeTo(t, a.WeightedRate, 625.0/676)
	closeTo(t, a.QuadraticScore, 25.0/676)
	closeTo(t, b.LoadWeight, 25.0/34)
	closeTo(t, b.QuadraticScore, 225.0/1156)
	if a.LoadWeight < 0.8 || b.LoadWeight >= 0.8 {
		t.Fatal("chapter's synthetic threshold comparison changed")
	}
	zero, err := Evaluate(0, 1, 0, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	closeTo(t, zero.LoadWeight, 0)
	closeTo(t, zero.Sigma, 1) // w_load/w0 is undefined here, not zero or one.
	// Doubling nominal rate doubles weighted rate but quadruples the chosen score.
	doubled, _ := Evaluate(1, 2, 0, 0.04, 0)
	closeTo(t, doubled.WeightedRate, 2*a.WeightedRate)
	closeTo(t, doubled.QuadraticScore, 4*a.QuadraticScore)
	// Under the optional throttle, increasing D need not increase the score.
	closeTo(t, e.Throttle[0].QuadraticScore, 0.25)
	closeTo(t, e.Throttle[1].QuadraticScore, 0.09)
}

func TestInvalidInputs(t *testing.T) {
	for _, x := range []float64{-1, math.NaN(), math.Inf(1)} {
		if _, err := Weight(x); err == nil {
			t.Fatalf("accepted degradation %v", x)
		}
		if _, err := MetrologyRating(x, 10); err == nil {
			t.Fatalf("accepted uncertainty %v", x)
		}
		if _, err := Evaluate(1, 1, 0, x, 0); err == nil {
			t.Fatalf("accepted model input %v", x)
		}
	}
	for _, x := range []float64{0, -1, 1.1, math.NaN(), math.Inf(1), math.SmallestNonzeroFloat64} {
		if _, err := Generator(x); err == nil {
			t.Fatalf("accepted nonrepresentable generator %v", x)
		}
	}
	for _, x := range []float64{-1, 1.1, math.NaN(), math.Inf(1)} {
		if _, err := Hamacher(x, 0.5); err == nil {
			t.Fatalf("accepted rating %v", x)
		}
		if _, err := ExpectedErrors(100, x); err == nil {
			t.Fatalf("accepted probability %v", x)
		}
	}
	for _, x := range []float64{-2, 2, math.NaN(), math.Inf(1)} {
		if _, err := SumVariance(1, 1, x); err == nil {
			t.Fatalf("accepted correlation %v", x)
		}
	}
	for _, tolerance := range []float64{0, -1, math.NaN(), math.Inf(1)} {
		if _, err := MetrologyRating(2, tolerance); err == nil {
			t.Fatalf("accepted tolerance %v", tolerance)
		}
	}
	for _, call := range []func() error{
		func() error { _, err := Evaluate(2, 1, 0, 0, 0); return err },
		func() error { _, err := Evaluate(1, 0, 0, 0, 0); return err },
		func() error { _, err := Evaluate(1, 1, math.MaxFloat64, math.MaxFloat64, 0); return err },
		func() error { _, err := Evaluate(1, math.MaxFloat64, 0, 1, 0); return err },
		func() error { _, err := MetrologyRating(math.MaxFloat64, 1); return err },
		func() error { _, err := SumVariance(math.MaxFloat64, 1, 0); return err },
		func() error { _, err := SumVariance(-1, 1, 0); return err },
	} {
		if call() == nil {
			t.Fatal("accepted invalid or overflowing calculation")
		}
	}
}

func TestEvidenceDeterministicJSON(t *testing.T) {
	a, err := SyntheticEvidence()
	if err != nil {
		t.Fatal(err)
	}
	b, err := SyntheticEvidence()
	if err != nil {
		t.Fatal(err)
	}
	ja, err := json.Marshal(a)
	if err != nil {
		t.Fatal(err)
	}
	jb, err := json.Marshal(b)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(ja, jb) {
		t.Fatal("evidence is not deterministic")
	}
	if a.Status != "synthetic_conditional_calculations_not_empirical_calibration" {
		t.Fatal("missing synthetic status")
	}
}
