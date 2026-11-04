package evidence

import (
	"testing"
	"time"
)

// TestChapter1ValveScenarios pins the five scenarios of table 1.6 in
// chapter 1: exact rule state, reason codes and the final permission.
func TestChapter1ValveScenarios(t *testing.T) {
	roles := Roles{Actor: "ctrl", Action: "open", Target: "valve_A"}
	rule := Rule{Pattern: roles, Sort: "bar", Window: Interval{Low: 7, High: 8},
		MaxAge: time.Second, Modal: True}
	cases := []struct {
		name      string
		sort      string
		rng       Interval
		age       time.Duration
		wantState State
		wantRsn   []Reason
		permitted bool
	}{
		{"A штатний вимір", "bar", Interval{7.7, 7.9}, 250 * time.Millisecond, True, nil, true},
		{"Б перетин межі", "bar", Interval{7.5, 8.1}, 250 * time.Millisecond, Both, []Reason{Boundary}, false},
		{"В застарілі дані", "bar", Interval{7.7, 7.9}, 1200 * time.Millisecond, False, []Reason{Stale}, false},
		{"Г інша одиниця", "kPa", Interval{770, 790}, 250 * time.Millisecond, Unknown, []Reason{SortMismatch}, false},
		{"Д аварійне перевищення", "bar", Interval{8.2, 8.5}, 250 * time.Millisecond, False, []Reason{OutOfRange}, false},
	}
	for _, c := range cases {
		d, err := EvaluateRule(rule, Fact{Roles: roles, Sort: c.sort, Range: c.rng, Age: c.age})
		if err != nil {
			t.Fatalf("%s: %v", c.name, err)
		}
		if d.State != c.wantState {
			t.Errorf("%s: state %s, want %s", c.name, d.State, c.wantState)
		}
		if len(d.Reasons) != len(c.wantRsn) {
			t.Errorf("%s: reasons %v, want %v", c.name, d.Reasons, c.wantRsn)
			continue
		}
		for i := range c.wantRsn {
			if d.Reasons[i] != c.wantRsn[i] {
				t.Errorf("%s: reasons %v, want %v", c.name, d.Reasons, c.wantRsn)
			}
		}
		// The sensor fact is TRUE in every scenario, so Priority returns the rule state.
		out, err := Priority(True, d.State)
		if err != nil {
			t.Fatalf("%s: %v", c.name, err)
		}
		if out != c.wantState || (out == True) != c.permitted {
			t.Errorf("%s: Priority=%s permitted=%t, want %s/%t", c.name, out, out == True, c.wantState, c.permitted)
		}
	}
}
