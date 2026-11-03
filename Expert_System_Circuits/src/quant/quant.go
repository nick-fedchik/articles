// Package quant evaluates the declared chapter-2 rating model.
// Ratings are not calibrated probabilities; no inference or admission is performed.
package quant

import (
	"fmt"
	"math"
)

func nonnegative(name string, x float64) error {
	if math.IsNaN(x) || math.IsInf(x, 0) || x < 0 {
		return fmt.Errorf("%s must be finite and nonnegative", name)
	}
	return nil
}

func probability(name string, x float64) error {
	if err := nonnegative(name, x); err != nil {
		return err
	}
	if x > 1 {
		return fmt.Errorf("%s must not exceed 1", name)
	}
	return nil
}

// Generator is g(w)=1/w-1 on (0,1]. Zero has only an extended-real limit.
func Generator(w float64) (float64, error) {
	if err := probability("weight", w); err != nil {
		return 0, err
	}
	if w == 0 {
		return 0, fmt.Errorf("generator requires a positive weight")
	}
	d := (1 - w) / w
	return d, nonnegative("generator result", d)
}

// Weight is the chosen inverse generator, not a probability estimator.
func Weight(d float64) (float64, error) {
	if err := nonnegative("degradation", d); err != nil {
		return 0, err
	}
	return 1 / (1 + d), nil
}

// Hamacher composes ratings by the zero-parameter Hamacher t-norm.
func Hamacher(a, b float64) (float64, error) {
	if err := probability("a", a); err != nil {
		return 0, err
	}
	if err := probability("b", b); err != nil {
		return 0, err
	}
	if a == 0 || b == 0 {
		return 0, nil
	}
	if a > b {
		a, b = b, a
	}
	// Equivalent to ab/(a+b-ab), without underflow of the product ab.
	return a / (1 + a/b - a), nil
}

// SumVariance computes Var(X+Y) in squared input units for given standard
// uncertainties and correlation. It is exact for a linear measurement model.
func SumVariance(u1, u2, rho float64) (float64, error) {
	if err := nonnegative("u1", u1); err != nil {
		return 0, err
	}
	if err := nonnegative("u2", u2); err != nil {
		return 0, err
	}
	if math.IsNaN(rho) || math.IsInf(rho, 0) || rho < -1 || rho > 1 {
		return 0, fmt.Errorf("correlation must be finite in [-1,1]")
	}
	// This rearrangement avoids cancellation near rho=-1.
	v := (u1-u2)*(u1-u2) + 2*(1+rho)*u1*u2
	return v, nonnegative("variance", v)
}

// MetrologyRating normalizes expanded uncertainty U by a positive tolerance
// width in the same units. This rating transformation is a design choice.
func MetrologyRating(expanded, tolerance float64) (float64, error) {
	if err := nonnegative("expanded uncertainty", expanded); err != nil {
		return 0, err
	}
	if err := nonnegative("tolerance", tolerance); err != nil {
		return 0, err
	}
	if tolerance == 0 {
		return 0, fmt.Errorf("tolerance must be positive")
	}
	ratio := expanded / tolerance
	d := ratio * ratio
	return d, nonnegative("metrology rating", d)
}

type Metrics struct {
	Degradation    float64 `json:"degradation"`
	Sigma          float64 `json:"sigma"`
	LoadWeight     float64 `json:"load_weight"`
	PolicyRate     float64 `json:"policy_rate_per_tick"`
	WeightedRate   float64 `json:"weighted_rate_per_tick"`
	QuadraticScore float64 `json:"quadratic_score_per_tick_squared"`
}

// Evaluate applies formulas 2.2-2.4, including the OPTIONAL throttle
// Phi=Phi0*w_load. All degradation inputs are declared dimensionless ratings.
func Evaluate(w0, phi0, src, gum, sem float64) (Metrics, error) {
	if err := probability("initial weight", w0); err != nil {
		return Metrics{}, err
	}
	for _, x := range []struct {
		name string
		x    float64
	}{{"nominal rate", phi0}, {"source rating", src}, {"metrology rating", gum}, {"semantic rating", sem}} {
		if err := nonnegative(x.name, x.x); err != nil {
			return Metrics{}, err
		}
	}
	if phi0 == 0 {
		return Metrics{}, fmt.Errorf("nominal rate must be positive")
	}
	d := src + gum + sem
	sigma, err := Weight(d)
	if err != nil {
		return Metrics{}, err
	}
	w := w0 * sigma
	phi := phi0 * w
	m := Metrics{d, sigma, w, phi, w * phi, phi * phi * d}
	for _, x := range []float64{m.PolicyRate, m.WeightedRate, m.QuadraticScore} {
		if err := nonnegative("model output", x); err != nil {
			return Metrics{}, err
		}
	}
	return m, nil
}

// ExpectedErrors is E[K]=N*p for a declared fixed marginal error probability.
// Independence is not needed for this expectation. Float64 results need not
// distinguish adjacent transaction counts above 2^53.
func ExpectedErrors(transactions uint64, p float64) (float64, error) {
	if err := probability("error probability", p); err != nil {
		return 0, err
	}
	return float64(transactions) * p, nil
}
