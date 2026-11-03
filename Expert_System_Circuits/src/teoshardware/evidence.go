package teoshardware

import (
	"fmt"

	evidence "github.com/nick-fedchik/articles/Expert_System_Circuits/src"
)

// EncodeEvidence maps the software state to GateInput.Evidence: bit 0 is PRO,
// bit 1 is CONTRA. It only converts encoding; it does not validate evidence.
func EncodeEvidence(state evidence.State) (uint8, error) {
	switch state {
	case evidence.Unknown:
		return 0, nil
	case evidence.True:
		return 1, nil
	case evidence.False:
		return 2, nil
	case evidence.Both:
		return 3, nil
	default:
		return 0, fmt.Errorf("encode evidence: invalid state %d", state)
	}
}
