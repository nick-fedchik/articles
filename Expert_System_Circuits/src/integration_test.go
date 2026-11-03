package evidence

import (
	"context"
	"testing"
)

func TestTypedDecisionProvenanceAndDiagnosticRepair(t *testing.T) {
	rule, fact := exampleRule(), exampleFact()
	fact.Age = rule.MaxAge + 1
	decision, err := EvaluateRule(rule, fact)
	if err != nil || decision.State != False || len(decision.Reasons) != 1 || decision.Reasons[0] != Stale {
		t.Fatalf("stale decision: %+v %v", decision, err)
	}
	problem := RepairProblem{
		Actions: []Action{{"refresh-measurement", 10, true, false}, {"disable-timeout", 1, true, true}},
		Cuts:    [][]int{{0, 1}},
	}
	repair, err := FindRepair(context.Background(), problem, 100)
	if err != nil || repair.Mask != 1 {
		t.Fatalf("diagnostic repair: %+v %v", repair, err)
	}
	// This is an explicit test intervention, not autonomous execution of a repair.
	fact.Age = 0
	decision, err = EvaluateRule(rule, fact)
	if err != nil {
		t.Fatal(err)
	}
	program, err := Compile(&Expr{Op: All, Left: &Expr{Op: Input, Input: 0}, Right: &Expr{Op: Input, Input: 1}}, fixtureSources(t, 2))
	if err != nil {
		t.Fatal(err)
	}
	result, err := program.Evaluate([]State{decision.State, True})
	if err != nil || result.State != True || len(result.Sources) != 2 {
		t.Fatalf("re-evaluation: %+v %v", result, err)
	}
}

func TestRuleBankOwnsItsCompiledSnapshot(t *testing.T) {
	rules := []Rule{exampleRule()}
	bank, err := CompileRules(rules)
	if err != nil {
		t.Fatal(err)
	}
	rules[0].Modal = False
	out := make([]State, 1)
	if err := bank.Evaluate(exampleFact(), out); err != nil || out[0] != True {
		t.Fatalf("compiled snapshot mutated through caller data: %v %v", out, err)
	}
}

func TestUninitializedExecutorsReturnErrors(t *testing.T) {
	var program Program
	if _, err := program.Evaluate(nil); err == nil {
		t.Fatal("uninitialized program accepted")
	}

	var bank RuleBank
	if err := bank.Evaluate(exampleFact(), nil); err == nil {
		t.Fatal("uninitialized bank accepted")
	}
}

func TestSpecialNormRequiresResolvedCondition(t *testing.T) {
	general, special, fact := exampleRule(), exampleRule(), exampleFact()
	special.Window.High = 22
	fact.Range = Interval{20, 20}
	genDecision, err := EvaluateRule(general, fact)
	if err != nil {
		t.Fatal(err)
	}
	specDecision, err := EvaluateRule(special, fact)
	if err != nil {
		t.Fatal(err)
	}
	root := &Expr{Op: Select, Condition: &Expr{Op: Input, Input: 0},
		Left: &Expr{Op: Input, Input: 1}, Right: &Expr{Op: Input, Input: 2}}
	program, err := Compile(root, fixtureSources(t, 3))
	if err != nil {
		t.Fatal(err)
	}
	for _, condition := range []State{True, False, Unknown, Both} {
		got, err := program.Evaluate([]State{condition, genDecision.State, specDecision.State})
		if err != nil || (got.State == True) != (condition == True) {
			t.Fatalf("exception condition=%v, result=%+v, error=%v", condition, got, err)
		}
	}
}
