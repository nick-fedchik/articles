package teoshardware

import (
	"testing"

	evidence "github.com/nick-fedchik/articles/Epistemic_Circuit_Theory/src"
)

func TestEncodeEvidenceAllStates(t *testing.T) {
	for _, tt := range []struct {
		state evidence.State
		code  uint8
		allow bool
	}{
		{evidence.Unknown, 0, false},
		{evidence.False, 2, false},
		{evidence.True, 1, true},
		{evidence.Both, 3, false},
	} {
		t.Run(tt.state.String(), func(t *testing.T) {
			got, err := EncodeEvidence(tt.state)
			if err != nil || got != tt.code {
				t.Fatalf("EncodeEvidence(%s) = %d, %v; want %d, nil", tt.state, got, err, tt.code)
			}
			in := safeGate()
			in.Evidence = got
			if decision := PrevalidatedGate(in); decision.Allow != tt.allow {
				t.Fatalf("encoded %s: gate = %+v; want allow=%t", tt.state, decision, tt.allow)
			}
		})
	}
}

func TestEncodeEvidenceRejectsInvalidStates(t *testing.T) {
	for raw := 0; raw <= 255; raw++ {
		state := evidence.State(raw)
		if state.Valid() {
			continue
		}
		got, err := EncodeEvidence(state)
		if err == nil || got != 0 {
			t.Fatalf("EncodeEvidence(%d) = %d, %v; want 0 and error", raw, got, err)
		}
		in := safeGate()
		in.Evidence = got
		if PrevalidatedGate(in).Allow {
			t.Fatalf("invalid state %d produced an allowing code", raw)
		}
	}
}
