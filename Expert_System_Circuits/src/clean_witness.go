package evidence

import (
	"errors"
	"slices"
	"sort"
)

type CleanSupportValue struct {
	Code    int
	Present bool
}

type CleanWitnessCertificate struct {
	Theory   [32]byte
	Revision uint64
	Query    SignedLiteral
	Proof    HornProof
	Absence  []CleanSupportValue
}

type CleanAnswer struct {
	Status  CleanStatus
	Witness *CleanWitnessCertificate
	Full    *CleanCertificate
}

func witnessNodes(proof HornProof, program *HornProgram) []int {
	nodes := map[int]bool{proof.Query: true}
	for _, code := range proof.Facts {
		nodes[code] = true
	}
	for _, index := range proof.Rules {
		rule := program.rules[index]
		nodes[rule.Conclusion] = true
		for _, premise := range rule.Premises {
			nodes[premise] = true
		}
	}
	out := make([]int, 0, len(nodes))
	for code := range nodes {
		out = append(out, code)
	}
	sort.Ints(out)
	return out
}

func cloneCleanWitness(witness CleanWitnessCertificate) CleanWitnessCertificate {
	witness.Absence = slices.Clone(witness.Absence)
	witness.Proof.Facts = slices.Clone(witness.Proof.Facts)
	witness.Proof.Rules = slices.Clone(witness.Proof.Rules)
	return witness
}

func (p *CleanProgram) extractWitness(full CleanCertificate) CleanWitnessCertificate {
	proof := *full.Proof
	witness := CleanWitnessCertificate{Theory: p.core.digest, Revision: proof.Revision, Query: full.Query, Proof: proof}
	pending := []int{}
	for _, code := range witnessNodes(proof, p.core) {
		pending = append(pending, code^1)
	}
	seen := map[int]bool{}
	for len(pending) > 0 {
		code := pending[len(pending)-1]
		pending = pending[:len(pending)-1]
		if seen[code] {
			continue
		}
		seen[code] = true
		witness.Absence = append(witness.Absence, CleanSupportValue{code, full.Closure[code]})
		for _, index := range p.core.producers[code] {
			pending = append(pending, p.core.rules[index].Premises...)
		}
	}
	sort.Slice(witness.Absence, func(i, j int) bool { return witness.Absence[i].Code < witness.Absence[j].Code })
	return witness
}

func (p *CleanProgram) VerifyCleanWitness(w SignedWorld, query SignedLiteral, witness CleanWitnessCertificate) error {
	snapshot, err := p.snapshot(w, query)
	if err != nil {
		return err
	}
	if witness.Theory != p.core.digest || witness.Revision != w.Revision || witness.Query != query ||
		len(witness.Absence) > p.core.atoms {
		return errors.New("clean witness: theory, revision, query or fragment size mismatch")
	}
	if err := p.core.VerifyProof(snapshot, query.code(), witness.Proof); err != nil {
		return err
	}
	values := map[int]bool{}
	for _, item := range witness.Absence {
		if item.Code < 0 || item.Code >= p.core.atoms {
			return errors.New("clean witness: literal outside range")
		}
		if _, exists := values[item.Code]; exists {
			return errors.New("clean witness: duplicate closure literal")
		}
		values[item.Code] = item.Present
	}
	for code, actual := range values {
		expected := snapshot.Facts[code]
		for _, index := range p.core.producers[code] {
			holds := snapshot.Enabled[index]
			for _, premise := range p.core.rules[index].Premises {
				value, exists := values[premise]
				if !exists {
					return errors.New("clean witness: closure fragment is not structurally closed")
				}
				holds = holds && value
			}
			expected = expected || holds
		}
		if actual != expected {
			return errors.New("clean witness: closure equation violated")
		}
	}
	for _, code := range witnessNodes(witness.Proof, p.core) {
		opposite, exists := values[code^1]
		if !exists || opposite {
			return errors.New("clean witness: opposition is missing or supported")
		}
	}
	return nil
}

func (p *CleanProgram) FreshCleanAnswer(w SignedWorld, query SignedLiteral) (CleanAnswer, error) {
	full, _, err := p.EvaluateClean(w, query)
	if err != nil {
		return CleanAnswer{}, err
	}
	answer := CleanAnswer{Status: full.Status}
	if full.Status == CleanSupported {
		witness := p.extractWitness(full)
		answer.Witness = &witness
	} else {
		answer.Full = &full
	}
	if err := p.VerifyCleanAnswer(w, query, answer); err != nil {
		return CleanAnswer{}, err
	}
	return answer, nil
}

func (p *CleanProgram) VerifyCleanAnswer(w SignedWorld, query SignedLiteral, answer CleanAnswer) error {
	if answer.Status == CleanSupported {
		if answer.Witness == nil || answer.Full != nil {
			return errors.New("clean answer: a positive answer needs exactly one sparse witness")
		}
		return p.VerifyCleanWitness(w, query, *answer.Witness)
	}
	if answer.Witness != nil || answer.Full == nil || answer.Full.Status != answer.Status ||
		(answer.Status != NotSupported && answer.Status != ConflictDependent) {
		return errors.New("clean answer: negative answer requires a complete certificate")
	}
	return p.VerifyClean(w, query, *answer.Full)
}
