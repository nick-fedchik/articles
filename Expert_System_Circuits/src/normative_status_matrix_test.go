package evidence

import "testing"

// TestRuleStatusMatrixAllConditionExceptionPairs fixes the documented
// precedence of formula (5.1) in chapter 5 for all 4 x 4 pairs of condition
// and exception states: exception TRUE excepts a rule whose condition is not
// FALSE; otherwise the condition is examined before the exception.
func TestRuleStatusMatrixAllConditionExceptionPairs(t *testing.T) {
	exception := SignedLiteral{Atom: 1}
	theory := NormTheory{Atoms: 2, Norms: []NormRule{{
		ID: "permit", Version: "1", Action: "run", Effect: Permit,
		Conditions: []SignedLiteral{{Atom: 0}}, Exception: &exception,
	}}}
	p, err := CompileNormative(theory)
	if err != nil {
		t.Fatal(err)
	}
	status := func(c, e State) (string, NormDecision) {
		switch {
		case c == False:
			return "INACTIVE", NormNoBasis
		case e == True:
			return "EXCEPTED", NormNoBasis
		case c == Unknown:
			return "UNKNOWN_CONDITION", NormUnknown
		case c == Both:
			return "CONFLICTING_CONDITION", NormConflict
		case e == False:
			return "APPLICABLE", NormAllow
		case e == Unknown:
			return "UNKNOWN_EXCEPTION", NormUnknown
		default:
			return "CONFLICTING_EXCEPTION", NormConflict
		}
	}
	pairs := 0
	for _, c := range []State{Unknown, False, True, Both} {
		for _, e := range []State{Unknown, False, True, Both} {
			world, err := p.NewWorld([]State{c, e})
			if err != nil {
				t.Fatal(err)
			}
			cert, err := p.Evaluate(world, "run")
			if err != nil {
				t.Fatal(err)
			}
			wantStatus, wantDecision := status(c, e)
			got := cert.Result.Rules[0]
			if got.Status != wantStatus || cert.Result.Decision != wantDecision {
				t.Errorf("C=%v E=%v: status %s decision %s, want %s %s",
					c, e, got.Status, cert.Result.Decision, wantStatus, wantDecision)
			}
			pairs++
		}
	}
	if pairs != 16 {
		t.Fatalf("checked %d pairs, want 16", pairs)
	}
}
