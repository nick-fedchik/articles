package evidence

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"sort"
	"strings"
)

type SignedLiteral struct {
	Atom     int
	Negative bool
}

func (l SignedLiteral) code() int {
	if l.Negative {
		return 2*l.Atom + 1
	}
	return 2 * l.Atom
}

func (l SignedLiteral) valid(atoms int) bool { return l.Atom >= 0 && l.Atom < atoms }

type SignedRule struct {
	ID, Version string
	Premises    []SignedLiteral
	Conclusion  SignedLiteral
}

type NormEffect string

const (
	Permit   NormEffect = "PERMIT"
	Prohibit NormEffect = "PROHIBIT"
)

type NormRule struct {
	ID, Version, Action string
	Effect              NormEffect
	Conditions          []SignedLiteral
	Exception           *SignedLiteral
}

type NormPriority struct{ Higher, Lower string }

type NormTheory struct {
	Atoms      int
	Evidence   []SignedRule
	Norms      []NormRule
	Priorities []NormPriority
}

type NormWorld struct {
	Revision        uint64
	Facts           []State
	EvidenceEnabled []bool
	NormEnabled     []bool
}

type NormDecision string

const (
	NormAllow    NormDecision = "ALLOW"
	NormDeny     NormDecision = "DENY"
	NormUnknown  NormDecision = "DEFER_UNKNOWN"
	NormConflict NormDecision = "DEFER_CONFLICT"
	NormNoBasis  NormDecision = "DEFER_NO_BASIS"
)

type NormObservation struct {
	ID                   string
	Condition, Exception State
	HasException         bool
	Status               string
	DominatedBy          []string
}

type NormResult struct {
	Action                        string
	ConditionEvidence, Applicable State
	Surviving                     State
	Decision                      NormDecision
	Rules                         []NormObservation
}

type NormCertificate struct {
	Theory, Context [32]byte
	Closure         []State
	Result          NormResult
}

type NormProgram struct {
	theory  NormTheory
	signed  *HornProgram
	higher  [][]bool
	byID    map[string]int
	actions map[string]bool
	digest  [32]byte
}

func CompileNormative(theory NormTheory) (*NormProgram, error) {
	if theory.Atoms < 1 || theory.Atoms > 128 || len(theory.Evidence) > 4096 ||
		len(theory.Norms) == 0 || len(theory.Norms) > 128 || len(theory.Priorities) > 4096 {
		return nil, errors.New("normative: finite model size limits exceeded")
	}
	p := &NormProgram{theory: NormTheory{Atoms: theory.Atoms}, byID: map[string]int{}, actions: map[string]bool{}}
	signedRules, copiedRules, err := signedHornRules(theory.Atoms, theory.Evidence)
	if err != nil {
		return nil, err
	}
	p.theory.Evidence = copiedRules
	p.signed, err = CompileHorn(2*theory.Atoms, signedRules)
	if err != nil {
		return nil, fmt.Errorf("normative signed fragment: %w", err)
	}
	for i, rule := range theory.Norms {
		if strings.TrimSpace(rule.ID) == "" || strings.TrimSpace(rule.Version) == "" || strings.TrimSpace(rule.Action) == "" ||
			(rule.Effect != Permit && rule.Effect != Prohibit) || len(rule.Conditions) > 64 {
			return nil, errors.New("normative: invalid rule identity, action, effect or arity")
		}
		if _, exists := p.byID[rule.ID]; exists {
			return nil, errors.New("normative: duplicate normative rule ID")
		}
		copyRule := rule
		copyRule.Conditions = slices.Clone(rule.Conditions)
		for _, literal := range copyRule.Conditions {
			if !literal.valid(theory.Atoms) {
				return nil, errors.New("normative: condition outside atom range")
			}
		}
		if rule.Exception != nil {
			if !rule.Exception.valid(theory.Atoms) {
				return nil, errors.New("normative: exception outside atom range")
			}
			copyLiteral := *rule.Exception
			copyRule.Exception = &copyLiteral
		}
		p.byID[rule.ID], p.actions[rule.Action] = i, true
		p.theory.Norms = append(p.theory.Norms, copyRule)
	}
	p.higher = make([][]bool, len(theory.Norms))
	for i := range p.higher {
		p.higher[i] = make([]bool, len(theory.Norms))
	}
	for _, edge := range theory.Priorities {
		h, hOK := p.byID[edge.Higher]
		l, lOK := p.byID[edge.Lower]
		if !hOK || !lOK || h == l || p.higher[h][l] || theory.Norms[h].Action != theory.Norms[l].Action {
			return nil, errors.New("normative: invalid, repeated or cross-action priority")
		}
		p.higher[h][l] = true
		p.theory.Priorities = append(p.theory.Priorities, edge)
	}
	for k := range p.higher {
		for i := range p.higher {
			for j := range p.higher {
				p.higher[i][j] = p.higher[i][j] || p.higher[i][k] && p.higher[k][j]
			}
		}
	}
	for i := range p.higher {
		if p.higher[i][i] {
			return nil, errors.New("normative: cyclic priority order")
		}
	}
	data, err := json.Marshal(p.theory)
	if err != nil {
		return nil, err
	}
	p.digest = sha256.Sum256(data)
	return p, nil
}

func (p *NormProgram) NewWorld(facts []State) (NormWorld, error) {
	if p == nil {
		return NormWorld{}, errors.New("normative: uninitialized program")
	}
	w := NormWorld{Revision: 1, Facts: slices.Clone(facts),
		EvidenceEnabled: make([]bool, len(p.theory.Evidence)), NormEnabled: make([]bool, len(p.theory.Norms))}
	for i := range w.EvidenceEnabled {
		w.EvidenceEnabled[i] = true
	}
	for i := range w.NormEnabled {
		w.NormEnabled[i] = true
	}
	if err := p.validateWorld(w); err != nil {
		return NormWorld{}, err
	}
	return w, nil
}

func (p *NormProgram) validateWorld(w NormWorld) error {
	if p == nil || p.signed == nil || w.Revision == 0 || len(w.Facts) != p.theory.Atoms ||
		len(w.EvidenceEnabled) != len(p.theory.Evidence) || len(w.NormEnabled) != len(p.theory.Norms) {
		return errors.New("normative: malformed world or uninitialized program")
	}
	for _, state := range w.Facts {
		if !state.Valid() {
			return errors.New("normative: invalid evidence state")
		}
	}
	return nil
}

func (p *NormProgram) context(w NormWorld) ([32]byte, error) {
	if err := p.validateWorld(w); err != nil {
		return [32]byte{}, err
	}
	data, err := json.Marshal(struct {
		Theory [32]byte
		World  NormWorld
	}{p.digest, w})
	if err != nil {
		return [32]byte{}, err
	}
	return sha256.Sum256(data), nil
}

func literalState(states []State, literal SignedLiteral) State {
	state := states[literal.Atom]
	if literal.Negative {
		return ((state & 1) << 1) | ((state & 2) >> 1)
	}
	return state
}

func effectState(effect NormEffect) State {
	if effect == Permit {
		return True
	}
	return False
}

func (p *NormProgram) Evaluate(w NormWorld, action string) (NormCertificate, error) {
	context, err := p.context(w)
	if err != nil {
		return NormCertificate{}, err
	}
	if !p.actions[action] {
		return NormCertificate{}, errors.New("normative: unknown action")
	}
	facts := make([]bool, 2*p.theory.Atoms)
	for i, state := range w.Facts {
		facts[2*i], facts[2*i+1] = state&2 != 0, state&1 != 0
	}
	closure, err := p.signed.FullClosure(HornSnapshot{w.Revision, facts, w.EvidenceEnabled})
	if err != nil {
		return NormCertificate{}, err
	}
	states := make([]State, p.theory.Atoms)
	for i := range states {
		if closure[2*i] {
			states[i] |= True
		}
		if closure[2*i+1] {
			states[i] |= False
		}
	}
	result := p.evaluatePolicy(w, action, states)
	certificate := NormCertificate{p.digest, context, states, result}
	if err := p.Verify(w, action, certificate); err != nil {
		return NormCertificate{}, err
	}
	return certificate, nil
}

func (p *NormProgram) evaluatePolicy(w NormWorld, action string, states []State) NormResult {
	result := NormResult{Action: action}
	active := make([]bool, len(p.theory.Norms))
	unknown, conflict := false, false
	for i, rule := range p.theory.Norms {
		if rule.Action != action {
			continue
		}
		condition := True
		for _, literal := range rule.Conditions {
			next := literalState(states, literal)
			condition = (condition & next & 2) | ((condition | next) & 1)
		}
		exception := False
		if rule.Exception != nil {
			exception = literalState(states, *rule.Exception)
		}
		observation := NormObservation{ID: rule.ID, Condition: condition, Exception: exception, HasException: rule.Exception != nil}
		if w.NormEnabled[i] && condition == True {
			result.ConditionEvidence |= effectState(rule.Effect)
		}
		switch {
		case !w.NormEnabled[i]:
			observation.Status = "DISABLED"
		case condition == False:
			observation.Status = "INACTIVE"
		case exception == True:
			observation.Status = "EXCEPTED"
		case condition == Unknown:
			observation.Status, unknown = "UNKNOWN_CONDITION", true
		case condition == Both:
			observation.Status, conflict = "CONFLICTING_CONDITION", true
		case exception == Unknown:
			observation.Status, unknown = "UNKNOWN_EXCEPTION", true
		case exception == Both:
			observation.Status, conflict = "CONFLICTING_EXCEPTION", true
		default:
			observation.Status, active[i] = "APPLICABLE", true
			result.Applicable |= effectState(rule.Effect)
		}
		result.Rules = append(result.Rules, observation)
	}
	for i := range result.Rules {
		index := p.byID[result.Rules[i].ID]
		if !active[index] {
			continue
		}
		for higher, enabled := range active {
			if enabled && p.higher[higher][index] {
				result.Rules[i].DominatedBy = append(result.Rules[i].DominatedBy, p.theory.Norms[higher].ID)
			}
		}
		sort.Strings(result.Rules[i].DominatedBy)
		if len(result.Rules[i].DominatedBy) == 0 {
			result.Surviving |= effectState(p.theory.Norms[index].Effect)
		}
	}
	switch {
	case conflict:
		result.Decision = NormConflict
	case unknown:
		result.Decision = NormUnknown
	case result.Surviving == Both:
		result.Decision = NormConflict
	case result.Surviving == True:
		result.Decision = NormAllow
	case result.Surviving == False:
		result.Decision = NormDeny
	default:
		result.Decision = NormNoBasis
	}
	sort.Slice(result.Rules, func(i, j int) bool { return result.Rules[i].ID < result.Rules[j].ID })
	return result
}
