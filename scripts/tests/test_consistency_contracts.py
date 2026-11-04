import itertools
import unittest


class ModelBoundaryTests(unittest.TestCase):
    def test_tara_example_matches_declared_profile_and_is_not_safety_closure(self):
        factor_scores = {
            "elapsed": 1,
            "expertise": 3,
            "knowledge": 3,
            "window": 1,
            "equipment_standard": 0,
        }
        attack_score = sum(factor_scores.values())
        category = "High" if attack_score <= 13 else "Other"

        def screening_gate(score, mitigation_verified):
            return score >= 20 or mitigation_verified

        def close_safety_case(gate_passed, safety_claim_verified, other_evidence_complete):
            return gate_passed and safety_claim_verified and other_evidence_complete

        self.assertEqual(attack_score, 8)
        self.assertEqual(category, "High")
        self.assertFalse(screening_gate(attack_score, False))
        self.assertTrue(screening_gate(20, False))
        self.assertFalse(close_safety_case(screening_gate(20, False), False, False))

    def test_modal_strength_does_not_resolve_unknown_applicability(self):
        def resolve(applicability, priority):
            if applicability != "both_applicable" or priority is None:
                return "REFUSAL_NORMATIVE_CONFLICT"
            return "higher_priority_rule"

        self.assertEqual(resolve("unknown", 1), "REFUSAL_NORMATIVE_CONFLICT")
        self.assertEqual(resolve("one_inapplicable", 1), "REFUSAL_NORMATIVE_CONFLICT")
        self.assertEqual(resolve("both_applicable", None), "REFUSAL_NORMATIVE_CONFLICT")
        self.assertEqual(resolve("both_applicable", 1), "higher_priority_rule")

    def test_fixed_frame_requires_a_write_set_bound(self):
        frame_bytes = 128
        illustrative_header_bytes = 32
        address_and_old_value_bytes = 16
        six_writes = illustrative_header_bytes + 6 * address_and_old_value_bytes
        seven_writes = illustrative_header_bytes + 7 * address_and_old_value_bytes
        self.assertLessEqual(six_writes, frame_bytes)
        self.assertGreater(seven_writes, frame_bytes)

    def test_grounding_is_not_semantic_or_execution_authority(self):
        cases = [(True, False, True), (True, True, False), (True, True, True)]
        results = [grounded and semantic_eligibility and authority
                   for grounded, semantic_eligibility, authority in cases]
        self.assertEqual(results, [False, False, True])
        self.assertTrue(all(grounded for grounded, _, _ in cases))

    def test_satisfiable_implication_does_not_establish_entailment(self):
        assignments = list(itertools.product((False, True), repeat=2))
        self.assertTrue(any(not hypothesis or observation for hypothesis, observation in assignments))
        self.assertTrue(any(hypothesis and not observation for hypothesis, observation in assignments))

    def test_consistency_is_required_to_avoid_vacuous_support(self):
        assignments = list(itertools.product((False, True), repeat=2))
        consistent = [(hypothesis, observation) for hypothesis, observation in assignments
                      if (not hypothesis or observation) and hypothesis]
        self.assertTrue(consistent)
        self.assertFalse(any(not observation for _, observation in consistent))
        contradictory = [(hypothesis, observation) for hypothesis, observation in assignments
                         if not hypothesis and hypothesis]
        self.assertFalse(contradictory)

    def test_conditioning_is_not_a_structural_intervention(self):
        assignments = list(itertools.product((False, True), repeat=3))
        conditioned = [(external, hypothesis, observation) for external, hypothesis, observation in assignments
                       if hypothesis == external and observation == external and not hypothesis and observation]
        intervened = [(external, hypothesis, observation) for external, hypothesis, observation in assignments
                      if not hypothesis and observation == external and observation]
        self.assertFalse(conditioned)
        self.assertIn((True, False, True), intervened)


if __name__ == "__main__":
    unittest.main()