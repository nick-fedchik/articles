package teoskernel

import (
	"time"

	evidence "github.com/nick-fedchik/articles/Epistemic_Circuit_Theory/src"
)

// Example returns the canonical valve-opening case used by tests, the command
// line tool and baseline comparisons. Atoms: 0 pressure_ok (guard on
// "pressure" in bar, window [0, 8]), 1 operator_authorized, 2
// maintenance_lock, 3 interlock_engaged (derived from maintenance_lock).
// The expected verdict is ALLOW.
func Example() (Spec, Input) {
	const now = int64(1_793_700_000_000_000_000)
	interlock := evidence.SignedLiteral{Atom: 3}
	spec := Spec{
		Action: "open_valve",
		Theory: evidence.NormTheory{
			Atoms: 4,
			Evidence: []evidence.SignedRule{{ID: "lock_engages_interlock", Version: "1",
				Premises: []evidence.SignedLiteral{{Atom: 2}}, Conclusion: interlock}},
			Norms: []evidence.NormRule{
				{ID: "open_when_safe", Version: "1", Action: "open_valve", Effect: evidence.Permit,
					Conditions: []evidence.SignedLiteral{{Atom: 0}, {Atom: 1}}, Exception: &interlock},
				{ID: "ban_under_lock", Version: "1", Action: "open_valve", Effect: evidence.Prohibit,
					Conditions: []evidence.SignedLiteral{{Atom: 2}}},
				{ID: "default_deny", Version: "1", Action: "open_valve", Effect: evidence.Prohibit},
			},
			Priorities: []evidence.NormPriority{
				{Higher: "ban_under_lock", Lower: "open_when_safe"},
				{Higher: "open_when_safe", Lower: "default_deny"},
			},
		},
		Guards:         []Guard{{Atom: 0, Quantity: "pressure", Unit: "bar", Low: 0, High: 8}},
		TrustedSources: []string{"plc-a", "badge-reader", "maintenance-db"},
		MaxAge:         2 * time.Second,
	}
	input := Input{
		Now: now,
		Facts: []Fact{
			{Literal: evidence.SignedLiteral{Atom: 1}, Sources: []string{"badge-reader"}, Time: now - int64(time.Second)},
			{Literal: evidence.SignedLiteral{Atom: 2, Negative: true}, Sources: []string{"maintenance-db"}, Time: now - int64(time.Second)},
			{Literal: evidence.SignedLiteral{Atom: 3, Negative: true}, Sources: []string{"plc-a"}, Time: now},
		},
		Measurements: []Measurement{{Quantity: "pressure", Unit: "bar", Low: 5.8, High: 6.2,
			Sources: []string{"plc-a"}, Time: now - int64(100*time.Millisecond)}},
	}
	return spec, input
}
