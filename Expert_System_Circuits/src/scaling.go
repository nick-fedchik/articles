package evidence

import (
	"errors"
	"runtime"
	"sort"
	"time"
)

type ScalingSample struct {
	Rules           int
	Layout, Cache   string
	Samples         int
	MedianNS, P95NS int64
	EvictionBytes   int
}

func makeScalingRules(n int) []Rule {
	rules := make([]Rule, n)
	for i := range rules {
		rules[i] = Rule{
			Pattern: Roles{"controller", "opens", "valve"}, Sort: "pressure:MPa",
			Window: Interval{10, float64(11 + (i*37)%9)}, MaxAge: time.Duration(1+i%50) * time.Millisecond,
			Modal: State(i % 4),
		}
		if i%5 == 0 {
			rules[i].Pattern.Actor = "*"
		}
	}
	return rules
}

// MeasureLayouts uses explicit cache displacement, not a guarantee of cold hardware caches.
func MeasureLayouts(sizes []int, samples int) ([]ScalingSample, error) {
	if samples < 2 || samples > 100 {
		return nil, errors.New("scaling: samples outside [2,100]")
	}
	for _, size := range sizes {
		if size < 1 || size > 100000 {
			return nil, errors.New("scaling: rule count outside [1,100000]")
		}
	}
	eviction := make([]byte, 64<<20)
	for i := range eviction {
		eviction[i] = byte(i)
	}
	out := []ScalingSample{}
	fact := Fact{Roles{"controller", "opens", "valve"}, "pressure:MPa", Interval{12, 13}, 3 * time.Millisecond}
	for _, n := range sizes {
		rules := makeScalingRules(n)
		bank, err := CompileRules(rules)
		if err != nil {
			return nil, err
		}
		target, reference := make([]State, n), make([]State, n)
		for i, rule := range rules {
			reference[i], err = ReferenceRule(rule, fact)
			if err != nil {
				return nil, err
			}
		}
		run := func(mode int) error {
			if mode >= 2 {
				return bank.evaluate(fact, target, mode == 3)
			}
			if err := fact.validate(); err != nil {
				return err
			}
			for i, rule := range rules {
				target[i] = evaluateValidatedRuleUsing(rule, fact, mode == 1)
			}
			return nil
		}
		for mode, layout := range []string{"AoS-branch", "AoS-LUT", "SoA-branch", "SoA-LUT"} {
			if err := run(mode); err != nil {
				return nil, err
			}
			for i := range target {
				if target[i] != reference[i] {
					return nil, errors.New("scaling: optimized mode differs from independent reference")
				}
			}
			for _, cache := range []string{"warm", "displaced-64MiB"} {
				durations := make([]int64, samples)
				for sample := 0; sample < samples; sample++ {
					if cache != "warm" {
						for i := 0; i < len(eviction); i += 64 {
							eviction[i]++
						}
					}
					start := time.Now()
					if err := run(mode); err != nil {
						return nil, err
					}
					durations[sample] = time.Since(start).Nanoseconds()
				}
				sort.Slice(durations, func(i, j int) bool { return durations[i] < durations[j] })
				evicted := 0
				if cache != "warm" {
					evicted = len(eviction)
				}
				out = append(out, ScalingSample{n, layout, cache, samples, durations[samples/2], durations[(95*samples-1)/100], evicted})
			}
		}
		runtime.KeepAlive(target)
	}
	runtime.KeepAlive(eviction)
	return out, nil
}
