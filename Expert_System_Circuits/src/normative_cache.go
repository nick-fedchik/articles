package evidence

import (
	"errors"
	"slices"
)

type NormCache struct {
	entries map[string]NormCertificate
}

type NormCachedResult struct {
	Certificate NormCertificate
	Reused      bool
	Invalidated bool
}

func NewNormCache() *NormCache {
	return &NormCache{entries: map[string]NormCertificate{}}
}

func cloneNormCertificate(certificate NormCertificate) NormCertificate {
	certificate.Closure = slices.Clone(certificate.Closure)
	certificate.Result.Rules = slices.Clone(certificate.Result.Rules)
	for i := range certificate.Result.Rules {
		certificate.Result.Rules[i].DominatedBy = slices.Clone(certificate.Result.Rules[i].DominatedBy)
	}
	return certificate
}

func (c *NormCache) Resolve(p *NormProgram, w NormWorld, action string) (NormCachedResult, error) {
	if c == nil || c.entries == nil {
		return NormCachedResult{}, errors.New("normative cache: uninitialized")
	}
	context, err := p.context(w)
	if err != nil {
		return NormCachedResult{}, err
	}
	if !p.actions[action] {
		return NormCachedResult{}, errors.New("normative cache: unknown action")
	}
	out := NormCachedResult{}
	if previous, exists := c.entries[action]; exists {
		if previous.Context == context && previous.Theory == p.digest {
			if err := p.Verify(w, action, previous); err != nil {
				return NormCachedResult{}, err
			}
			return NormCachedResult{Certificate: cloneNormCertificate(previous), Reused: true}, nil
		}
		out.Invalidated = true
		delete(c.entries, action)
	}
	if len(c.entries) >= 128 {
		return NormCachedResult{}, errors.New("normative cache: action capacity exceeded")
	}
	certificate, err := p.Evaluate(w, action)
	if err != nil {
		return NormCachedResult{}, err
	}
	c.entries[action] = cloneNormCertificate(certificate)
	out.Certificate = certificate
	return out, nil
}
