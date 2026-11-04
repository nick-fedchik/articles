// Package baseline executes hypotheses H3 (comparison with baselines) and
// H4 (fault injection in a synthetic safety case) of the preregistered
// protocol PREREGISTRATION-2026-11-03.md.
//
// The synthetic case is the admission of process gas into a vessel through an
// admission valve (action "open_admission_valve"). The action is physically
// safe iff the true current vessel pressure lies in the hard window [0, 8] bar,
// the operator is authorized, a process demand exists, no maintenance lock and
// no interlock are active, and active hot work is covered by a passed gas test.
//
// Four arms decide on the same scenarios:
//
//   - A  naive: scalar confidences (noisy-OR of per-source reliabilities) and
//     point comparisons with a flat deny-first rule chain (ArmA);
//   - B0 conventional engineered stack with conventional information: typed
//     pressure units, point value plus a fixed instrument-uncertainty margin,
//     one boolean per fact (latest report wins), Horn closure, ordered
//     first-match rules, two-sided freshness and an append-only audit log that
//     is not consulted for decisions (ArmB0);
//   - B1 the same conventional stack given the same metadata as the kernel:
//     intervals, pro/contra evidence, source sets, timestamps, exceptions and
//     priorities, engineered from the derived STPA constraints (ArmB1);
//   - C  the kernel kernel.Evaluate over CaseSpec (ArmC).
//
// Ground truth is reported twice: the kernel safety property P
// (kernel.CheckSafety, circular for C because ALLOW <=> P) and an
// independent physical truth computed from the hidden world state (World.Safe).
package baseline

import (
	"time"

	evidence "github.com/nick-fedchik/articles/Expert_System_Circuits/src"
	"github.com/nick-fedchik/articles/Expert_System_Circuits/src/kernel"
)

// Atom indices of the case theory. AtomUnreferencedNoise is an unreferenced
// telemetry channel that does not appear in any guard, evidence rule, or norm
// of CaseSpec (used by F11_UNREFERENCED_SPAM to test F-01 relevantAtom scoping).
const (
	AtomPressureOK = iota
	AtomAuthorized
	AtomLock
	AtomInterlock
	AtomDemand
	AtomHotWork
	AtomGasTest
	AtomUnreferencedNoise
	atomCount
)

// Case constants shared by every arm.
const (
	Action        = "open_admission_valve"
	Quantity      = "pressure"
	LimitLowBar   = 0.0
	LimitHighBar  = 8.0
	MaxAge        = 2 * time.Second
	Now           = int64(1_793_700_000_000_000_000)
	SourcePLC     = "plc-a"
	SourceBadge   = "badge-reader"
	SourceMaint   = "maintenance-db"
	SourceDCS     = "dcs"
	SourcePTW     = "ptw-db"
	SourceGas     = "gas-detector"
	SourceLLM     = "llm-agent"
	SourceGateway = "field-gateway-x"
)

// AtomNames maps atom indices to names used by the conventional arms.
var AtomNames = [atomCount]string{"pressure_ok", "operator_authorized", "maintenance_lock",
	"interlock_engaged", "process_demand", "hot_work_active", "gas_test_passed", "unreferenced_telemetry"}

// TrustedSources is the set of identities verified upstream.
func TrustedSources() []string {
	return []string{SourcePLC, SourceBadge, SourceMaint, SourceDCS, SourcePTW, SourceGas}
}

func isTrusted(id string) bool {
	for _, t := range TrustedSources() {
		if t == id {
			return true
		}
	}
	return false
}

// CaseSpec is the arm C specification: the norms of the case written in the
// kernel's signed normative theory.
func CaseSpec() kernel.Spec {
	// spec:C:begin
	interlock := evidence.SignedLiteral{Atom: AtomInterlock}
	base := []evidence.SignedLiteral{{Atom: AtomPressureOK}, {Atom: AtomAuthorized}, {Atom: AtomDemand}}
	return kernel.Spec{
		Action: Action,
		Theory: evidence.NormTheory{
			Atoms: atomCount,
			Evidence: []evidence.SignedRule{{ID: "lock_engages_interlock", Version: "1",
				Premises: []evidence.SignedLiteral{{Atom: AtomLock}}, Conclusion: interlock}},
			Norms: []evidence.NormRule{
				{ID: "open_when_safe", Version: "1", Action: Action, Effect: evidence.Permit,
					Conditions: base, Exception: &interlock},
				{ID: "hot_work_cleared", Version: "1", Action: Action, Effect: evidence.Permit,
					Conditions: append(append([]evidence.SignedLiteral{}, base...),
						evidence.SignedLiteral{Atom: AtomHotWork}, evidence.SignedLiteral{Atom: AtomGasTest}),
					Exception: &interlock},
				{ID: "ban_hot_work", Version: "1", Action: Action, Effect: evidence.Prohibit,
					Conditions: []evidence.SignedLiteral{{Atom: AtomHotWork}}},
				{ID: "ban_under_lock", Version: "1", Action: Action, Effect: evidence.Prohibit,
					Conditions: []evidence.SignedLiteral{{Atom: AtomLock}}},
				{ID: "default_deny", Version: "1", Action: Action, Effect: evidence.Prohibit},
			},
			Priorities: []evidence.NormPriority{
				{Higher: "ban_under_lock", Lower: "open_when_safe"},
				{Higher: "ban_under_lock", Lower: "hot_work_cleared"},
				{Higher: "hot_work_cleared", Lower: "ban_hot_work"},
				{Higher: "ban_hot_work", Lower: "open_when_safe"},
				{Higher: "open_when_safe", Lower: "default_deny"},
				{Higher: "hot_work_cleared", Lower: "default_deny"},
			},
		},
		Guards: []kernel.Guard{{Atom: AtomPressureOK, Quantity: Quantity, Unit: "bar",
			Low: LimitLowBar, High: LimitHighBar}},
		TrustedSources: TrustedSources(),
		MaxAge:         MaxAge,
	}
	// spec:C:end
}

// World is the hidden physical state; Pressure is the true current value in
// bar. Interlock already includes the interlock engaged by a maintenance lock.
type World struct {
	Pressure                                          float64
	Authorized, Demand, Lock, Interlock, HotWork, Gas bool
}

// Safe is the independent physical ground truth of the case.
func (w World) Safe() bool {
	return w.Pressure >= LimitLowBar && w.Pressure <= LimitHighBar && w.Authorized && w.Demand &&
		!w.Lock && !w.Interlock && (!w.HotWork || w.Gas)
}
