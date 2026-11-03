package teosbaseline

import (
	"errors"
	"math"
)

// RuleOfThree is the approximate one-sided 95% upper bound 3/n for zero
// events in n trials.
func RuleOfThree(n int) float64 {
	if n <= 0 {
		return math.NaN()
	}
	return 3 / float64(n)
}

// BinomialCDF returns P(X <= x) for X ~ Binomial(n, p), summed in log space.
func BinomialCDF(x, n int, p float64) float64 {
	switch {
	case x < 0:
		return 0
	case x >= n:
		return 1
	case p <= 0:
		return 1
	case p >= 1:
		return 0
	}
	lgn, _ := math.Lgamma(float64(n + 1))
	sum := 0.0
	for k := 0; k <= x; k++ {
		lk, _ := math.Lgamma(float64(k + 1))
		lnk, _ := math.Lgamma(float64(n - k + 1))
		sum += math.Exp(lgn - lk - lnk + float64(k)*math.Log(p) + float64(n-k)*math.Log1p(-p))
	}
	return math.Min(sum, 1)
}

// ClopperPearsonUpper returns the exact Clopper-Pearson upper bound for x
// events in n trials: the p with P(X <= x; n, p) = tail, found by bisection
// (tail = 0.025 for the two-sided 95% interval, 0.05 for one-sided 95%). It
// equals the Beta(x+1, n-x) quantile at 1-tail.
func ClopperPearsonUpper(x, n int, tail float64) (float64, error) {
	if n <= 0 || x < 0 || x > n || !(tail > 0 && tail < 1) {
		return math.NaN(), errors.New("teosbaseline: invalid Clopper-Pearson arguments")
	}
	if x == n {
		return 1, nil
	}
	if x == 0 {
		return -math.Expm1(math.Log(tail) / float64(n)), nil
	}
	lo, hi := float64(x)/float64(n), 1.0
	for range 200 {
		mid := (lo + hi) / 2
		if BinomialCDF(x, n, mid) > tail {
			lo = mid
		} else {
			hi = mid
		}
	}
	return (lo + hi) / 2, nil
}

// RegIncBeta is the regularized incomplete beta function I_z(a, b), by the
// Lentz continued fraction; it is used to cross-check the binomial route.
func RegIncBeta(z, a, b float64) float64 {
	if z <= 0 {
		return 0
	}
	if z >= 1 {
		return 1
	}
	la, _ := math.Lgamma(a)
	lb, _ := math.Lgamma(b)
	lab, _ := math.Lgamma(a + b)
	front := math.Exp(lab - la - lb + a*math.Log(z) + b*math.Log1p(-z))
	if z > (a+1)/(a+b+2) {
		return 1 - RegIncBeta(1-z, b, a)
	}
	const tiny = 1e-300
	f, c, d := 1.0, 1.0, 0.0
	for i := 0; i <= 10000; i++ {
		m := float64(i / 2)
		var num float64
		switch {
		case i == 0:
			num = 1
		case i%2 == 0:
			num = m * (b - m) * z / ((a + 2*m - 1) * (a + 2*m))
		default:
			num = -(a + m) * (a + b + m) * z / ((a + 2*m) * (a + 2*m + 1))
		}
		d = 1 + num*d
		if math.Abs(d) < tiny {
			d = tiny
		}
		d = 1 / d
		c = 1 + num/c
		if math.Abs(c) < tiny {
			c = tiny
		}
		cd := c * d
		f *= cd
		if math.Abs(1-cd) < 1e-15 {
			break
		}
	}
	return front * (f - 1) / a
}
