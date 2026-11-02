package evidence

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
	"time"
)

func TestPriorityLUTExhaustive256(t *testing.T) {
	for a := Unknown; a <= Both; a++ {
		for b := Unknown; b <= Both; b++ {
			for c := Unknown; c <= Both; c++ {
				for d := Unknown; d <= Both; d++ {
					got := lookupPriority(lookupPriority(a, b), lookupPriority(c, d))
					want, err := ReferencePriority(a, b, c, d)
					if err != nil || got != want {
						t.Fatalf("%v %v %v %v: got %v, want %v", a, b, c, d, got, want)
					}
				}
			}
		}
	}
}

func TestRuleBankAgreement(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	rules := make([]Rule, 1000)
	for i := range rules {
		rules[i] = exampleRule()
		rules[i].Window = Interval{float64(rng.Intn(10)), float64(10 + rng.Intn(10))}
		rules[i].MaxAge = time.Duration(1+rng.Intn(50)) * time.Millisecond
		rules[i].Modal = State(rng.Intn(4))
		if i%3 == 0 {
			rules[i].Pattern.Actor = "*"
		}
	}
	bank, err := CompileRules(rules)
	if err != nil {
		t.Fatal(err)
	}
	dst := make([]State, len(rules))
	for n := 0; n < 100; n++ {
		fact := exampleFact()
		fact.Range = Interval{float64(rng.Intn(20)), 20}
		fact.Age = time.Duration(rng.Intn(60)) * time.Millisecond
		if n%5 == 0 {
			fact.Roles.Actor, fact.Roles.Target = fact.Roles.Target, fact.Roles.Actor
		}
		if err := bank.Evaluate(fact, dst); err != nil {
			t.Fatal(err)
		}
		for i, rule := range rules {
			want, err := ReferenceRule(rule, fact)
			if err != nil || dst[i] != want {
				t.Fatalf("case %d rule %d: got %v, want %v (%v)", n, i, dst[i], want, err)
			}
		}
	}
}

func TestRuleBankRejectsInvalidInput(t *testing.T) {
	bank, _ := CompileRules([]Rule{exampleRule()})
	fact := exampleFact()
	fact.Range.High = math.Inf(1)
	if err := bank.Evaluate(fact, make([]State, 1)); err == nil {
		t.Fatal("infinite bound accepted")
	}
	if err := bank.Evaluate(exampleFact(), nil); err == nil {
		t.Fatal("wrong output length accepted")
	}
}

func FuzzRuleAgreement(f *testing.F) {
	f.Add(uint64(0), uint64(0), uint64(0), uint8(0), uint8(0))
	f.Add(math.Float64bits(15), math.Float64bits(1), uint64(20), uint8(2), uint8(1))
	f.Fuzz(func(t *testing.T, xb, rb, age uint64, modal, roles uint8) {
		rule, fact := exampleRule(), exampleFact()
		x, radius := math.Float64frombits(xb), math.Abs(math.Float64frombits(rb))
		fact.Range = Interval{x - radius, x + radius}
		fact.Age = time.Duration(age)
		rule.Modal = State(modal % 5)
		if roles%2 != 0 {
			fact.Roles.Target = fact.Roles.Actor
		}
		want, refErr := ReferenceRule(rule, fact)
		got, err := EvaluateRule(rule, fact)
		if (err == nil) != (refErr == nil) {
			t.Fatalf("validation mismatch: %v / %v", err, refErr)
		}
		if err != nil {
			return
		}
		bank, err := CompileRules([]Rule{rule})
		if err != nil {
			t.Fatal(err)
		}
		dst := make([]State, 1)
		if err := bank.Evaluate(fact, dst); err != nil || got.State != want || dst[0] != want {
			t.Fatalf("got %v/%v, want %v: %v", got.State, dst[0], want, err)
		}
	})
}

var benchmarkState State

func BenchmarkRuleBank(b *testing.B) {
	for _, n := range []int{1, 256, 10000} {
		rules := make([]Rule, n)
		for i := range rules {
			rules[i] = exampleRule()
			rules[i].Window.High = float64(11 + i%7)
			rules[i].Modal = State(i % 4)
		}
		bank, err := CompileRules(rules)
		if err != nil {
			b.Fatal(err)
		}
		fact := exampleFact()
		output := make([]State, n)
		b.Run(fmt.Sprintf("AoS-branching/%d", n), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for k := 0; k < b.N; k++ {
				if err := fact.validate(); err != nil {
					b.Fatal(err)
				}
				for i, rule := range rules {
					output[i] = evaluateValidatedRule(rule, fact)
				}
				benchmarkState = output[k%n]
			}
		})
		b.Run(fmt.Sprintf("SoA-LUT/%d", n), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for k := 0; k < b.N; k++ {
				if err := bank.Evaluate(fact, output); err != nil {
					b.Fatal(err)
				}
				benchmarkState = output[k%n]
			}
		})
	}
}
