package evidence

import (
	"fmt"
	"math"
	"sync"
	"testing"
)

func balancedExpression(begin, end int) *Expr {
	if end-begin == 1 {
		return &Expr{Op: Input, Input: begin}
	}
	middle := begin + (end-begin)/2
	return &Expr{Op: All, Left: balancedExpression(begin, middle), Right: balancedExpression(middle, end)}
}

func TestScratchBudgetAndImmutableUpdates(t *testing.T) {
	sources := fixtureSources(t, 130)
	root := balancedExpression(0, len(sources))
	program, err := Compile(root, sources)
	if err != nil {
		t.Fatal(err)
	}
	inputs := make([]State, len(sources))
	for i := range inputs {
		inputs[i] = True
	}
	required := int64(program.Nodes()) * (3*8 + 1)
	if program.ScratchBytes() != required {
		t.Fatal("incorrect exact scratch estimate")
	}
	if _, err := program.EvaluateWithBudget(inputs, required-1); err == nil {
		t.Fatal("scratch limit was ignored")
	}
	old, err := program.EvaluateWithBudget(inputs, required)
	if err != nil {
		t.Fatal(err)
	}
	updated, err := program.WithSourceWeight("s0", "1", .5)
	if err != nil {
		t.Fatal(err)
	}
	newResult, err := updated.Evaluate(inputs)
	if err != nil || newResult.Weight >= old.Weight {
		t.Fatal("updated weight was not propagated")
	}
	again, err := program.Evaluate(inputs)
	if err != nil || math.Abs(again.Weight-old.Weight) > 1e-14 {
		t.Fatal("old immutable snapshot changed")
	}
}

func TestConcurrentProgramAndRuleBankSnapshots(t *testing.T) {
	program, err := Compile(balancedExpression(0, 8), fixtureSources(t, 8))
	if err != nil {
		t.Fatal(err)
	}
	bank, err := CompileRules(makeScalingRules(256))
	if err != nil {
		t.Fatal(err)
	}
	var workers sync.WaitGroup
	for worker := 0; worker < 8; worker++ {
		workers.Add(1)
		go func() {
			defer workers.Done()
			inputs := []State{True, True, True, True, True, True, True, True}
			dst := make([]State, 256)
			for i := 0; i < 30; i++ {
				snapshot, err := program.WithSourceWeight("s0", "1", .8)
				if err != nil {
					t.Error(err)
					return
				}
				if _, err := snapshot.Evaluate(inputs); err != nil {
					t.Error(err)
				}
				if err := bank.Evaluate(exampleFact(), dst); err != nil {
					t.Error(err)
				}
			}
		}()
	}
	workers.Wait()
}

func TestAllLayoutAblationsAgree(t *testing.T) {
	rules := makeScalingRules(1000)
	bank, err := CompileRules(rules)
	if err != nil {
		t.Fatal(err)
	}
	dst := make([]State, len(rules))
	for _, lut := range []bool{false, true} {
		if err := bank.evaluate(exampleFact(), dst, lut); err != nil {
			t.Fatal(err)
		}
		for i, rule := range rules {
			want, err := ReferenceRule(rule, exampleFact())
			if err != nil || dst[i] != want || evaluateValidatedRuleUsing(rule, exampleFact(), lut) != want {
				t.Fatal("layout or switching ablation changed semantics")
			}
		}
	}
}

func BenchmarkGraphScaling(b *testing.B) {
	for _, count := range []int{16, 128, 512, 2048} {
		sources := fixtureSources(b, count)
		root := balancedExpression(0, count)
		program, err := Compile(root, sources)
		if err != nil {
			b.Fatal(err)
		}
		inputs := make([]State, count)
		for i := range inputs {
			inputs[i] = True
		}
		b.Run(fmt.Sprintf("Compile/%d", count), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				if _, err := Compile(root, sources); err != nil {
					b.Fatal(err)
				}
			}
		})
		b.Run(fmt.Sprintf("Evaluate/%d", count), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				out, err := program.Evaluate(inputs)
				if err != nil {
					b.Fatal(err)
				}
				benchmarkState = out.State
			}
		})
		b.Run(fmt.Sprintf("Reweight/%d", count), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				if _, err := program.WithSourceWeight("s0", "1", .8); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
