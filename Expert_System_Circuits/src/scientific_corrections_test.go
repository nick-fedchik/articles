package evidence

import (
	"math"
	"testing"
)

func TestScientificCorrections(t *testing.T) {
	t.Run("clone normalization", func(t *testing.T) {
		for _, n := range []int{1, 4, 100} {
			weight := .6
			gamma := 1 / math.Sqrt(float64(n))
			result := math.Sqrt(float64(n) * math.Pow(gamma*weight, 2))
			if math.Abs(result-weight) > 1e-14 {
				t.Fatalf("%d copies changed the weight: %.17g", n, result)
			}
		}
	})
	t.Run("cascade losses are attenuated downstream", func(t *testing.T) {
		weight := 1.0
		expected := []float64{.98, .96, .94, .732, .712}
		for i, coefficient := range []float64{1, 1, 1, .8, 1} {
			weight = math.Max(0, coefficient*weight-.02)
			if math.Abs(weight-expected[i]) > 1e-14 {
				t.Fatalf("stage %d: %.17g", i+1, weight)
			}
		}
	})
	t.Run("finite product error includes cross term", func(t *testing.T) {
		nominal := 1500.0 * 19.6
		linear := 1500*.5 + 19.6*30
		hard := linear + 30*.5
		low, high := 1470*19.1, 1530*20.1
		if linear != 1338 || hard != 1353 || math.Abs(low-28077) > 1e-9 ||
			math.Abs(high-30753) > 1e-9 || nominal+linear >= high || nominal+hard < high {
			t.Fatalf("product bounds: nominal=%g linear=%g hard=%g [%g,%g]", nominal, linear, hard, low, high)
		}
	})
	t.Run("comparator endpoints", func(t *testing.T) {
		classify := func(x, halfwidth, limit float64) State {
			if x+halfwidth <= limit {
				return True
			}
			if x-halfwidth > limit {
				return False
			}
			return Unknown
		}
		for _, test := range []struct {
			x, halfwidth, limit float64
			want                State
		}{{10, 0, 10, True}, {9, 1, 10, True}, {11, 1, 10, Unknown}, {12, 1, 10, False}, {29400, 1353, 30000, Unknown}} {
			if got := classify(test.x, test.halfwidth, test.limit); got != test.want {
				t.Fatalf("%+v: %v", test, got)
			}
		}
	})
	t.Run("threshold bits do not imply a differential margin", func(t *testing.T) {
		for _, test := range []struct {
			pro, contra float64
			want        State
		}{{.21, .19, True}, {.8, .3, Both}} {
			got := Unknown
			if test.pro >= .2 {
				got |= True
			}
			if test.contra >= .2 {
				got |= False
			}
			if got != test.want {
				t.Fatalf("%+v: %v", test, got)
			}
		}
	})
	t.Run("small curvature index is not sufficient", func(t *testing.T) {
		coefficients := []float64{0, 1, 0, 1}
		moments := []float64{1, 0, 1, 0, 3, 0, 15}
		first, second := coefficients[1], 2*coefficients[2]
		index := .5 * math.Abs(second/first)
		linearVariance := first * first
		mean, secondMoment := 0.0, 0.0
		for i, a := range coefficients {
			mean += a * moments[i]
			for j, b := range coefficients {
				secondMoment += a * b * moments[i+j]
			}
		}
		exactVariance := secondMoment - mean*mean
		if index != 0 || linearVariance != 1 || exactVariance != 22 {
			t.Fatal("normal moments counterexample failed")
		}
	})
	t.Run("withdrawn threshold uncertainty bound underestimates", func(t *testing.T) {
		mean, sigma, threshold := .55, .05, .50
		delta := mean - threshold
		z := delta / sigma
		phi := math.Exp(-z*z/2) / math.Sqrt(2*math.Pi)
		cdf := .5 * math.Erfc(-z/math.Sqrt2)
		first := sigma*phi + delta*cdf
		second := (delta*delta+sigma*sigma)*cdf + delta*sigma*phi
		exactSD := math.Sqrt(second - first*first)
		withdrawn := math.Max(0, mean+2*sigma-threshold)/2 - math.Max(0, mean-threshold)
		if math.Abs(withdrawn-.025) > 1e-14 || exactSD < .0433 || exactSD > .0434 || withdrawn >= exactSD {
			t.Fatalf("withdrawn=%g exact=%g", withdrawn, exactSD)
		}
		t.Logf("withdrawn bound %.8f < exact SD %.8f", withdrawn, exactSD)
	})
	t.Run("tree repair output may be exponential", func(t *testing.T) {
		const pairs = 20
		count := 0
		for choices := uint64(0); choices < 1<<pairs; choices++ {
			var repair, decoded uint64
			for pair := range pairs {
				choice := (choices >> pair) & 1
				repair |= 1 << (2*pair + int(choice))
			}
			for pair := range pairs {
				selected := (repair >> (2 * pair)) & 3
				if selected != 1 && selected != 2 {
					t.Fatal("repair is not minimal for an OR pair")
				}
				if selected == 2 {
					decoded |= 1 << pair
				}
			}
			if decoded != choices {
				t.Fatal("repair encoding is not injective")
			}
			count++
		}
		if count != 1048576 {
			t.Fatalf("wrong number of minimal repairs: %d", count)
		}
	})
}
