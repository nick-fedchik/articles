package evidence

import (
	"math"
	"testing"
	"time"
)

func TestResponseBudgetReservesCompleteReaction(t *testing.T) {
	budget := ResponseBudget{50 * time.Millisecond, 20 * time.Millisecond, 5 * time.Millisecond,
		time.Millisecond, 2 * time.Millisecond, 10 * time.Millisecond}
	remaining, err := budget.Remaining()
	if err != nil || remaining != 12*time.Millisecond {
		t.Fatalf("remaining %v, %v", remaining, err)
	}
	budget.MaxAge = budget.FTTI
	if _, err := budget.Remaining(); err == nil {
		t.Fatal("age equal to full FTTI consumed the actuator budget")
	}
	budget.FTTI, budget.MaxAge = time.Duration(math.MaxInt64), time.Duration(math.MaxInt64)
	if _, err := budget.Remaining(); err == nil {
		t.Fatal("overflowing budget accepted")
	}
}

func TestHysteresisBoundAndWorstCaseWitness(t *testing.T) {
	const A, delta, limit, correction = .38, .12, 15., .24
	for _, horizon := range []int{1, 3} {
		bound, err := MinimumHysteresis(A, delta, horizon)
		if err != nil {
			t.Fatal(err)
		}
		for _, offset := range []float64{0, -.01} {
			control, err := NewHysteresis(HysteresisConfig{Limit: limit, Correction: correction, Width: bound + offset})
			if err != nil {
				t.Fatal(err)
			}
			x := limit - correction - A + .005
			active, err := control.Step(0, x+A)
			if err != nil || !active {
				t.Fatal("witness did not trip")
			}
			released := false
			for j := 1; j <= horizon; j++ {
				active, err = control.Step(time.Duration(j), x-float64(j)*delta-A)
				if err != nil {
					t.Fatal(err)
				}
				released = released || !active
			}
			if offset == 0 && released {
				t.Fatalf("reversed within guaranteed horizon %d", horizon)
			}
			if offset < 0 && !released {
				t.Fatal("subcritical witness failed")
			}
		}
	}
}

func TestNoChatterDoesNotImplyRecovery(t *testing.T) {
	for _, tc := range []struct {
		width   float64
		recover bool
	}{{.88, true}, {1.12, false}} {
		result, err := SimulateMemoPlant(tc.width, 10000)
		if err != nil {
			t.Fatal(err)
		}
		if tc.recover && (result.Releases == 0 || result.RecoveryMargin <= 0) {
			t.Fatalf("expected reachable recovery: %+v", result)
		}
		if !tc.recover && (result.Releases != 0 || result.Trips != 1 || result.RecoveryMargin >= 0) {
			t.Fatalf("expected permanently latched protection: %+v", result)
		}
		t.Logf("%+v", result)
	}
}

func TestHysteresisHoldAndInvalidTime(t *testing.T) {
	h, _ := NewHysteresis(HysteresisConfig{Limit: 15, Width: 1, MinimumHold: 10})
	active, err := h.Step(0, 16)
	if err != nil || !active {
		t.Fatal(err)
	}
	active, err = h.Step(1, 13)
	if err != nil || !active {
		t.Fatal("minimum hold violated")
	}
	active, err = h.Step(10, 13)
	if err != nil || active {
		t.Fatal("failed to release after hold")
	}
	if _, err := h.Step(10, 13); err == nil {
		t.Fatal("nonmonotone clock accepted")
	}
	if _, err := h.Step(11, math.NaN()); err == nil {
		t.Fatal("NaN accepted")
	}
}
