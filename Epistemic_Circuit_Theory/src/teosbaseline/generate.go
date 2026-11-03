package teosbaseline

import (
	"fmt"
	"math/rand/v2"
	"time"

	evidence "github.com/nick-fedchik/articles/Epistemic_Circuit_Theory/src"
	"github.com/nick-fedchik/articles/Epistemic_Circuit_Theory/src/teoskernel"
)

// Class is a fault class derived from the unsafe control actions of the case.
type Class struct {
	Name, UCA, Description string
	// InModel is false for classes that violate a declared assumption of the
	// kernel (trusted sources honest, hard intervals contain the truth).
	InModel bool
	world   func(r *rand.Rand) World
	// build completes the world and returns its faulted observation and a
	// sub-variant label.
	build func(r *rand.Rand, w World) (World, teoskernel.Input, string)
	// control optionally changes the world of the matched clean control.
	control func(w World) World
}

// Scenario is one deterministic case: a hidden world and its observation.
type Scenario struct {
	Class   string
	Index   int
	Control bool
	Variant string
	World   World
	Input   teoskernel.Input
}

func uniform(r *rand.Rand, lo, hi float64) float64 { return lo + (hi-lo)*r.Float64() }

func durationIn(r *rand.Rand, lo, hi time.Duration) int64 {
	return int64(lo) + r.Int64N(int64(hi-lo)+1)
}

func bernoulli(r *rand.Rand, p float64) bool { return r.Float64() < p }

// baseWorld samples the nominal world distribution of the case.
func baseWorld(r *rand.Rand) World {
	w := World{
		Pressure:   uniform(r, 3, 10),
		Authorized: bernoulli(r, 0.9),
		Demand:     bernoulli(r, 0.9),
		Lock:       bernoulli(r, 0.1),
		HotWork:    bernoulli(r, 0.2),
		Gas:        bernoulli(r, 0.7),
	}
	w.Interlock = bernoulli(r, 0.05) || w.Lock
	return w
}

func fact(atom int, value bool, source string, t int64) teoskernel.Fact {
	return teoskernel.Fact{Literal: evidence.SignedLiteral{Atom: atom, Negative: !value},
		Sources: []string{source}, Time: t}
}

// measureAround builds an honest interval of half-width h containing p.
func measureAround(r *rand.Rand, p float64, t int64) teoskernel.Measurement {
	h := uniform(r, 0.05, 0.3)
	c := p + uniform(r, -0.9*h, 0.9*h)
	return teoskernel.Measurement{Quantity: Quantity, Unit: "bar", Low: c - h, High: c + h,
		Sources: []string{SourcePLC}, Time: t}
}

// Faithful is the fault-free observation of w: fresh, trusted, honest.
func Faithful(r *rand.Rand, w World) teoskernel.Input {
	recent := func() int64 { return Now - durationIn(r, 0, 1500*time.Millisecond) }
	return teoskernel.Input{
		Now:          Now,
		Measurements: []teoskernel.Measurement{measureAround(r, w.Pressure, Now-durationIn(r, 0, 500*time.Millisecond))},
		Facts: []teoskernel.Fact{
			fact(AtomAuthorized, w.Authorized, SourceBadge, recent()),
			fact(AtomLock, w.Lock, SourceMaint, recent()),
			fact(AtomInterlock, w.Interlock, SourcePLC, recent()),
			fact(AtomDemand, w.Demand, SourceDCS, recent()),
			fact(AtomHotWork, w.HotWork, SourcePTW, recent()),
			fact(AtomGasTest, w.Gas, SourceGas, recent()),
		},
	}
}

func findFact(in teoskernel.Input, atom int) int {
	for i, f := range in.Facts {
		if f.Literal.Atom == atom {
			return i
		}
	}
	return -1
}

func removeFact(in *teoskernel.Input, atom int) {
	out := in.Facts[:0]
	for _, f := range in.Facts {
		if f.Literal.Atom != atom {
			out = append(out, f)
		}
	}
	in.Facts = out
}

// aged models an old reading with a bad timestamp stamp(). For the
// measurement variant the reading was taken at a benign-looking past pressure
// while the true current pressure drifted (rising process); for the
// authorization variant an old positive authorization may since be revoked.
func aged(r *rand.Rand, w World, measurement bool, stamp func() int64, mv, fv string) (World, teoskernel.Input, string) {
	if measurement {
		past := uniform(r, 4, 7.8)
		w.Pressure = past + uniform(r, -0.5, 3)
		in := Faithful(r, w)
		in.Measurements[0] = measureAround(r, past, stamp())
		return w, in, mv
	}
	w.Authorized = bernoulli(r, 0.5)
	in := Faithful(r, w)
	in.Facts[findFact(in, AtomAuthorized)] = fact(AtomAuthorized, true, SourceBadge, stamp())
	return w, in, fv
}

// observe adapts an observation-only fault that keeps the world unchanged.
func observe(f func(r *rand.Rand, w World) (teoskernel.Input, string)) func(*rand.Rand, World) (World, teoskernel.Input, string) {
	return func(r *rand.Rand, w World) (World, teoskernel.Input, string) {
		in, v := f(r, w)
		return w, in, v
	}
}

// Classes returns the fault classes fixed before the experiment was run.
func Classes() []Class {
	return []Class{
		{Name: "CLEAN", UCA: "-", InModel: true,
			Description: "fault-free observation of the nominal world distribution",
			world:       baseWorld,
			build:       observe(func(r *rand.Rand, w World) (teoskernel.Input, string) { return Faithful(r, w), "clean" })},
		{Name: "F01_STALE", UCA: "UCA-1", InModel: true,
			Description: "measurement older than MaxAge while pressure drifted (75%), or stale authorization since revoked (25%)",
			world:       baseWorld,
			build: func(r *rand.Rand, w World) (World, teoskernel.Input, string) {
				return aged(r, w, r.IntN(4) < 3, func() int64 { return Now - durationIn(r, 2500*time.Millisecond, 120*time.Second) },
					"stale-measurement", "stale-authorization")
			}},
		{Name: "F02_UNIT", UCA: "UCA-2", InModel: true,
			Description: "honest interval reported in MPa, kPa or psi instead of the declared bar",
			world:       baseWorld,
			build: observe(func(r *rand.Rand, w World) (teoskernel.Input, string) {
				in := Faithful(r, w)
				unit := []string{"MPa", "kPa", "psi"}[r.IntN(3)]
				m := &in.Measurements[0]
				m.Low, _ = FromBar(m.Low, unit)
				m.High, _ = FromBar(m.High, unit)
				m.Unit = unit
				return in, unit
			})},
		{Name: "F03_STRADDLE", UCA: "UCA-3", InModel: true,
			Description: "honest interval straddles the 8 bar limit with its midpoint inside; truth uniform in the interval",
			world:       baseWorld,
			build:       straddle},
		{Name: "F04_CONTRARY", UCA: "UCA-4", InModel: true,
			Description: "two fresh trusted sources contradict each other on authorization or maintenance lock (state BOTH)",
			world:       baseWorld,
			build:       contrary},
		{Name: "F05_EXCEPTION_UNKNOWN", UCA: "UCA-5", InModel: true,
			Description: "interlock report missing; independent interlock causes active with probability 0.3",
			world: func(r *rand.Rand) World {
				w := baseWorld(r)
				w.Interlock = w.Lock || bernoulli(r, 0.3)
				return w
			},
			build: observe(func(r *rand.Rand, w World) (teoskernel.Input, string) {
				in := Faithful(r, w)
				removeFact(&in, AtomInterlock)
				return in, "interlock-missing"
			})},
		{Name: "F06_UNTRUSTED_PROVENANCE", UCA: "UCA-6", InModel: true,
			Description: "spoofed benign-looking pressure from an unverified gateway or without sources (2/3), or authorization from an unverified gateway (1/3)",
			world:       baseWorld,
			build:       untrusted},
		{Name: "F07_NORM_PRIORITY_CONFLICT", UCA: "UCA-7", InModel: true,
			Description: "hot work active (ban vs cleared permit by priority), gas test 50%, maintenance lock 30%; faithful observation",
			world: func(r *rand.Rand) World {
				w := baseWorld(r)
				w.HotWork, w.Gas, w.Lock = true, bernoulli(r, 0.5), bernoulli(r, 0.3)
				w.Interlock = w.Lock || bernoulli(r, 0.05)
				return w
			},
			build: observe(func(r *rand.Rand, w World) (teoskernel.Input, string) { return Faithful(r, w), "faithful" }),
			control: func(w World) World {
				w.HotWork, w.Lock = false, false
				w.Interlock = false
				return w
			}},
		{Name: "F08_LLM_UNSUPPORTED", UCA: "UCA-8", InModel: true,
			Description: "LLM proposes the action; process demand (truth 30%) absent or asserted only by the LLM agent",
			world: func(r *rand.Rand) World {
				w := baseWorld(r)
				w.Demand = bernoulli(r, 0.3)
				return w
			},
			build: observe(func(r *rand.Rand, w World) (teoskernel.Input, string) {
				in := Faithful(r, w)
				removeFact(&in, AtomDemand)
				if bernoulli(r, 0.5) {
					return in, "demand-absent"
				}
				copies := 1 + r.IntN(2)
				for range copies {
					in.Facts = append(in.Facts, fact(AtomDemand, true, SourceLLM, Now-durationIn(r, 0, time.Second)))
				}
				return in, fmt.Sprintf("llm-only-x%d", copies)
			})},
		{Name: "F09_DUPLICATED_SOURCES", UCA: "UCA-9", InModel: true,
			Description: "badge report replaced by k=2..5 relayed copies of one unverified positive authorization (truth 50%)",
			world: func(r *rand.Rand) World {
				w := baseWorld(r)
				w.Authorized = bernoulli(r, 0.5)
				return w
			},
			build: observe(func(r *rand.Rand, w World) (teoskernel.Input, string) {
				in := Faithful(r, w)
				removeFact(&in, AtomAuthorized)
				k := 2 + r.IntN(4)
				t := Now - durationIn(r, 0, time.Second)
				for j := 1; j <= k; j++ {
					in.Facts = append(in.Facts, fact(AtomAuthorized, true, fmt.Sprintf("relay-%d", j), t))
				}
				return in, fmt.Sprintf("copies-%d", k)
			})},
		{Name: "F10_CLOCK_FUTURE", UCA: "UCA-10", InModel: true,
			Description: "clock rollback/skew: old drifted measurement (2/3) or old authorization since revoked (1/3) stamped in the future",
			world:       baseWorld,
			build: func(r *rand.Rand, w World) (World, teoskernel.Input, string) {
				return aged(r, w, r.IntN(3) < 2, func() int64 { return Now + durationIn(r, 500*time.Millisecond, time.Hour) },
					"future-measurement", "future-authorization")
			}},
		{Name: "X11_OOM_TRUSTED_BIAS", UCA: "UCA-1/UCA-3", InModel: false,
			Description: "OUT OF MODEL: fresh trusted sensor reports a narrow interval that excludes the true pressure (bias below truth)",
			world: func(r *rand.Rand) World {
				w := baseWorld(r)
				w.Pressure = uniform(r, 6, 10)
				return w
			},
			build: observe(func(r *rand.Rand, w World) (teoskernel.Input, string) {
				in := Faithful(r, w)
				h := uniform(r, 0.05, 0.3)
				c := w.Pressure - h - uniform(r, 0.3, 2)
				in.Measurements[0].Low, in.Measurements[0].High = c-h, c+h
				return in, "biased-trusted"
			})},
	}
}

// straddle completes the F03 world and observation.
func straddle(r *rand.Rand, w World) (World, teoskernel.Input, string) {
	h := uniform(r, 0.2, 1.0)
	mid := uniform(r, LimitHighBar-0.95*h, LimitHighBar)
	w.Pressure = uniform(r, mid-h, mid+h)
	in := Faithful(r, w)
	in.Measurements[0].Low, in.Measurements[0].High = mid-h, mid+h
	return w, in, "straddle"
}

// contrary builds the F04 world and its contradicting observation.
func contrary(r *rand.Rand, w World) (World, teoskernel.Input, string) {
	byLock := bernoulli(r, 0.5)
	if byLock {
		w.Lock = bernoulli(r, 0.5)
		w.Interlock = w.Lock || bernoulli(r, 0.05)
	} else {
		w.Authorized = bernoulli(r, 0.5)
	}
	in := Faithful(r, w)
	wrong := fact(AtomAuthorized, !w.Authorized, SourceDCS, Now-durationIn(r, 0, 1500*time.Millisecond))
	variant := "authorization"
	if byLock {
		wrong = fact(AtomLock, !w.Lock, SourcePLC, Now-durationIn(r, 0, 1500*time.Millisecond))
		variant = "lock"
	}
	if bernoulli(r, 0.5) {
		in.Facts = append([]teoskernel.Fact{wrong}, in.Facts...)
	} else {
		in.Facts = append(in.Facts, wrong)
	}
	return w, in, variant
}

// untrusted builds the F06 observation.
func untrusted(r *rand.Rand, w World) (World, teoskernel.Input, string) {
	if r.IntN(3) < 2 {
		in := Faithful(r, w)
		m := measureAround(r, uniform(r, 3, 7.5), in.Measurements[0].Time)
		variant := "untrusted-measurement"
		m.Sources = []string{SourceGateway}
		if bernoulli(r, 0.5) {
			m.Sources, variant = nil, "no-provenance-measurement"
		}
		in.Measurements[0] = m
		return w, in, variant
	}
	w.Authorized = bernoulli(r, 0.5)
	in := Faithful(r, w)
	i := findFact(in, AtomAuthorized)
	in.Facts[i] = fact(AtomAuthorized, true, SourceGateway, in.Facts[i].Time)
	return w, in, "untrusted-authorization"
}

func classSeed(seed uint64, class, index int) *rand.Rand {
	return rand.New(rand.NewPCG(seed, uint64(class)<<32|uint64(index)))
}

// Generate returns, for every class, faults faulted scenarios followed by
// controls matched clean controls (CLEAN has no separate controls). A control
// is the faithful observation of the world the faulted builder produces for
// the same random stream, optionally adjusted by Class.control. The output
// depends only on the arguments.
func Generate(seed uint64, faults, controls int) []Scenario {
	var out []Scenario
	for ci, class := range Classes() {
		n := faults + controls
		if class.Name == "CLEAN" {
			n = faults
		}
		for i := range n {
			r := classSeed(seed, ci, i)
			s := Scenario{Class: class.Name, Index: i, Control: i >= faults}
			s.World, s.Input, s.Variant = class.build(r, class.world(r))
			if s.Control {
				if class.control != nil {
					s.World = class.control(s.World)
				}
				s.Input, s.Variant = Faithful(r, s.World), "control"
			}
			out = append(out, s)
		}
	}
	return out
}
