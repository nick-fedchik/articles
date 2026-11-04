package evidence

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
)

type CleanStatus string

const (
	CleanSupported    CleanStatus = "CLEAN_SUPPORTED"
	ConflictDependent CleanStatus = "CONFLICT_DEPENDENT"
	NotSupported      CleanStatus = "NOT_SUPPORTED"
)

type SignedWorld struct {
	Revision uint64
	Facts    []State
	Enabled  []bool
}

type CleanCertificate struct {
	Theory, Context [32]byte
	Query           SignedLiteral
	Status          CleanStatus
	Closure, Clean  []bool
	Proof           *HornProof
}

type CleanWork struct {
	Nodes, Rules, Premises int
}

type CleanProgram struct {
	core  *HornProgram
	rules []SignedRule
	atoms int
}

func CompileCleanSupport(atoms int, rules []SignedRule) (*CleanProgram, error) {
	if atoms < 1 || atoms > 10000 || len(rules) > 10000 {
		return nil, errors.New("clean support: atom/rule size limit exceeded")
	}
	encoded, copies, err := signedHornRules(atoms, rules)
	if err != nil {
		return nil, err
	}
	core, err := CompileHorn(2*atoms, encoded)
	if err != nil {
		return nil, fmt.Errorf("clean support: %w", err)
	}
	return &CleanProgram{core, copies, atoms}, nil
}

func (p *CleanProgram) snapshot(w SignedWorld, query SignedLiteral) (HornSnapshot, error) {
	if p == nil || p.core == nil || !query.valid(p.atoms) || len(w.Facts) != p.atoms {
		return HornSnapshot{}, errors.New("clean support: invalid program, world or query")
	}
	snapshot := HornSnapshot{Revision: w.Revision, Enabled: w.Enabled, Facts: make([]bool, 2*p.atoms)}
	for atom, state := range w.Facts {
		if !state.Valid() {
			return HornSnapshot{}, errors.New("clean support: invalid evidence state")
		}
		snapshot.Facts[2*atom], snapshot.Facts[2*atom+1] = state&2 != 0, state&1 != 0
	}
	if err := p.core.validate(snapshot, query.code()); err != nil {
		return HornSnapshot{}, err
	}
	return snapshot, nil
}

func (p *CleanProgram) context(w SignedWorld) ([32]byte, error) {
	data, err := json.Marshal(struct {
		Contract string
		Theory   [32]byte
		World    SignedWorld
	}{"clean-all-nodes-v1", p.core.digest, w})
	if err != nil {
		return [32]byte{}, err
	}
	return sha256.Sum256(data), nil
}

func cleanStatus(supported, clean bool) CleanStatus {
	if clean {
		return CleanSupported
	}
	if supported {
		return ConflictDependent
	}
	return NotSupported
}

func (p *CleanProgram) EvaluateClean(w SignedWorld, query SignedLiteral) (CleanCertificate, CleanWork, error) {
	snapshot, err := p.snapshot(w, query)
	if err != nil {
		return CleanCertificate{}, CleanWork{}, err
	}
	closure, err := p.core.FullClosure(snapshot)
	if err != nil {
		return CleanCertificate{}, CleanWork{}, err
	}
	known, parents := make([]bool, p.core.atoms), make([]int, p.core.atoms)
	work := CleanWork{}
	for _, code := range p.core.topology {
		work.Nodes++
		parents[code] = -1
		if !closure[code] || closure[code^1] {
			continue
		}
		if snapshot.Facts[code] {
			known[code] = true
			continue
		}
		for _, index := range p.core.producers[code] {
			work.Rules++
			if !snapshot.Enabled[index] {
				continue
			}
			holds := true
			for _, premise := range p.core.rules[index].Premises {
				work.Premises++
				if !known[premise] {
					holds = false
					break
				}
			}
			if holds {
				known[code], parents[code] = true, index
				break
			}
		}
	}
	context, err := p.context(w)
	if err != nil {
		return CleanCertificate{}, CleanWork{}, err
	}
	certificate := CleanCertificate{Theory: p.core.digest, Context: context, Query: query,
		Status: cleanStatus(closure[query.code()], known[query.code()]), Closure: closure, Clean: known}
	if known[query.code()] {
		proof, err := p.core.certificate(snapshot, query.code(), parents, func(int) bool { return true })
		if err != nil {
			return CleanCertificate{}, CleanWork{}, err
		}
		certificate.Proof = &proof
	}
	return certificate, work, nil
}

func (p *CleanProgram) VerifyClean(w SignedWorld, query SignedLiteral, certificate CleanCertificate) error {
	snapshot, err := p.snapshot(w, query)
	if err != nil {
		return err
	}
	context, err := p.context(w)
	if err != nil {
		return err
	}
	if certificate.Theory != p.core.digest || certificate.Context != context || certificate.Query != query ||
		len(certificate.Closure) != p.core.atoms || len(certificate.Clean) != p.core.atoms {
		return errors.New("clean certificate: theory, context, query or bitmap size mismatch")
	}
	for code, producers := range p.core.producers {
		support, derivation := snapshot.Facts[code], snapshot.Facts[code]
		for _, index := range producers {
			if !snapshot.Enabled[index] {
				continue
			}
			raw, clean := true, true
			for _, premise := range p.core.rules[index].Premises {
				raw = raw && certificate.Closure[premise]
				clean = clean && certificate.Clean[premise]
			}
			support, derivation = support || raw, derivation || clean
		}
		if certificate.Closure[code] != support ||
			certificate.Clean[code] != (certificate.Closure[code] && !certificate.Closure[code^1] && derivation) {
			return errors.New("clean certificate: complete closure/clean equation violated")
		}
	}
	want := cleanStatus(certificate.Closure[query.code()], certificate.Clean[query.code()])
	if certificate.Status != want || (certificate.Proof != nil) != (want == CleanSupported) {
		return errors.New("clean certificate: status or witness presence mismatch")
	}
	if certificate.Proof == nil {
		return nil
	}
	if err := p.core.VerifyProof(snapshot, query.code(), *certificate.Proof); err != nil {
		return err
	}
	for _, fact := range certificate.Proof.Facts {
		if !certificate.Clean[fact] {
			return errors.New("clean certificate: unclean witness fact")
		}
	}
	for _, index := range certificate.Proof.Rules {
		rule := p.core.rules[index]
		if !certificate.Clean[rule.Conclusion] {
			return errors.New("clean certificate: unclean witness conclusion")
		}
		for _, premise := range rule.Premises {
			if !certificate.Clean[premise] {
				return errors.New("clean certificate: unclean witness premise")
			}
		}
	}
	return nil
}
