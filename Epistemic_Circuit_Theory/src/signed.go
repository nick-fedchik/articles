package evidence

import (
	"errors"
	"slices"
)

func signedHornRules(atoms int, rules []SignedRule) ([]HornRule, []SignedRule, error) {
	encoded := []HornRule{}
	var copies []SignedRule
	for _, rule := range rules {
		if !rule.Conclusion.valid(atoms) {
			return nil, nil, errors.New("normative: evidence conclusion outside atom range")
		}
		copyRule := rule
		copyRule.Premises = slices.Clone(rule.Premises)
		indices := make([]int, len(rule.Premises))
		for i, literal := range rule.Premises {
			if !literal.valid(atoms) {
				return nil, nil, errors.New("normative: evidence premise outside atom range")
			}
			indices[i] = literal.code()
		}
		copies = append(copies, copyRule)
		encoded = append(encoded, HornRule{rule.ID, rule.Version, indices, rule.Conclusion.code()})
	}
	return encoded, copies, nil
}
