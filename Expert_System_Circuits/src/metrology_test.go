package evidence

import (
	"math"
	"testing"
)

func TestMemo03ProposedBoundCounterexample(t *testing.T) {
	// Analytic moments of max(0, X-threshold), X normally distributed.
	const mean, sigma, threshold, k = .55, .05, .50, 2.
	mu := mean - threshold
	alpha := mu / sigma
	density := math.Exp(-alpha*alpha/2) / math.Sqrt(2*math.Pi)
	cdf := .5 * (1 + math.Erf(alpha/math.Sqrt2))
	first := sigma*density + mu*cdf
	second := (mu*mu+sigma*sigma)*cdf + mu*sigma*density
	exact := math.Sqrt(second - first*first)
	proposed := math.Max(0, mean+k*sigma-threshold)/k - math.Max(0, mean-threshold)
	if math.Abs(proposed-.025) > 1e-12 || math.Abs(exact-.04333266) > 1e-8 || proposed >= exact {
		t.Fatalf("counterexample changed: proposed=%g, exact=%g", proposed, exact)
	}
	t.Logf("memo 03 proposed=%0.8f, exact standard deviation=%0.8f; not a conservative bound", proposed, exact)
}
