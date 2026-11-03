package evidence

import "errors"

var ErrCleanProofBudget = errors.New("clean reference: proof enumeration budget exhausted")

type CleanReference struct {
	Supported, ExistsClean, LeavesOnly, AllProofsClean bool
	Proofs, Work                                       int
	Closure                                            []bool
}

func (p *CleanProgram) ReferenceClean(w SignedWorld, query SignedLiteral, budget int) (CleanReference, error) {
	if _, err := p.snapshot(w, query); err != nil {
		return CleanReference{}, err
	}
	if budget < 1 || budget > 1000000 {
		return CleanReference{}, errors.New("clean reference: budget outside [1,1000000]")
	}
	work := 0
	tick := func() error {
		if work >= budget {
			return ErrCleanProofBudget
		}
		work++
		return nil
	}
	fact := func(code int) bool {
		state := w.Facts[code/2]
		if code%2 == 0 {
			return state == True || state == Both
		}
		return state == False || state == Both
	}
	closure := make([]bool, 2*p.atoms)
	for code := range closure {
		if err := tick(); err != nil {
			return CleanReference{}, err
		}
		closure[code] = fact(code)
	}
	changed := true
	for changed {
		changed = false
		for i := len(p.rules) - 1; i >= 0; i-- {
			if err := tick(); err != nil {
				return CleanReference{}, err
			}
			rule := p.rules[i]
			if !w.Enabled[i] || closure[rule.Conclusion.code()] {
				continue
			}
			holds := true
			for _, literal := range rule.Premises {
				if err := tick(); err != nil {
					return CleanReference{}, err
				}
				holds = holds && closure[literal.code()]
			}
			if holds {
				closure[rule.Conclusion.code()], changed = true, true
			}
		}
	}
	type quality struct{ clean, leaves bool }
	memo, active := map[int][]quality{}, map[int]bool{}
	var enumerate func(int) ([]quality, error)
	enumerate = func(code int) ([]quality, error) {
		if err := tick(); err != nil {
			return nil, err
		}
		if values, exists := memo[code]; exists {
			return values, nil
		}
		if active[code] {
			return nil, errors.New("clean reference: unexpected cyclic proof dependency")
		}
		active[code] = true
		values := []quality{}
		local := closure[code] && !closure[code^1]
		if fact(code) {
			values = append(values, quality{local, local})
		}
		for index, rule := range p.rules {
			if err := tick(); err != nil {
				return nil, err
			}
			if !w.Enabled[index] || rule.Conclusion.code() != code {
				continue
			}
			combinations := []quality{{true, true}}
			for _, premise := range rule.Premises {
				children, err := enumerate(premise.code())
				if err != nil {
					return nil, err
				}
				next := []quality{}
				for _, prefix := range combinations {
					for _, child := range children {
						if err := tick(); err != nil {
							return nil, err
						}
						next = append(next, quality{prefix.clean && child.clean, prefix.leaves && child.leaves})
					}
				}
				combinations = next
				if len(combinations) == 0 {
					break
				}
			}
			for _, combination := range combinations {
				if err := tick(); err != nil {
					return nil, err
				}
				values = append(values, quality{local && combination.clean, combination.leaves})
			}
		}
		active[code] = false
		memo[code] = values
		return values, nil
	}
	proofs, err := enumerate(query.code())
	if err != nil {
		return CleanReference{}, err
	}
	out := CleanReference{Supported: len(proofs) != 0, Proofs: len(proofs), Work: work, Closure: closure}
	out.AllProofsClean = out.Supported
	for _, proof := range proofs {
		out.ExistsClean = out.ExistsClean || proof.clean
		out.LeavesOnly = out.LeavesOnly || proof.leaves
		out.AllProofsClean = out.AllProofsClean && proof.clean
	}
	out.LeavesOnly = out.LeavesOnly && closure[query.code()] && !closure[query.code()^1]
	if out.Supported != closure[query.code()] {
		return CleanReference{}, errors.New("clean reference: proof enumeration contradicts independent closure")
	}
	return out, nil
}
