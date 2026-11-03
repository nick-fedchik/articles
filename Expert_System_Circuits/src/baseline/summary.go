package baseline

// Comparison lists unsafe permits of every arm on one in-model fault class.
type Comparison struct {
	Class                string
	N                    int
	UnsafeP, UnsafePhys  map[string]int
	FalseRefusalPhys     map[string]int
	CBetterThanB0P       bool
	CBetterThanB1P       bool
	CBetterThanB0Phys    bool
	CBetterThanB1Phys    bool
	B1MatchesCUnsafeP    bool
	B1MatchesCUnsafePhys bool
}

// H4Summary aggregates arm C on a defined scope of classes, tracking how many
// classes have Allow > 0 versus Allow == 0 (total refusal, conditional risk
// per ALLOW undefined 0/0).
type H4Summary struct {
	Scope                 string
	Classes               int
	ClassesWithAllow      int
	ClassesZeroAllow      int
	Scenarios, Allow      int
	UnsafeP, UnsafePhys   int
	UnsafePPerAllow       Bound
	UnsafePhysPerAllow    Bound
	UnsafePhysPerScenario Bound
}

// H4ClassSummary reports arm C's stratified H4 bounds on a single class (F-08).
// When Allow == 0 (TotalRefusal is true), UnsafePPerAllow.Defined and
// UnsafePhysPerAllow.Defined are false (Undefined 0/0 conditional risk per
// ALLOW), and the class is bounded by the unconditional scenario bound
// UnsafePhysPerScenario (denominator Scenarios).
type H4ClassSummary struct {
	Class                 string
	Group                 string
	InModel               bool
	TotalRefusal          bool
	Scenarios, Allow      int
	Deny, Defer           int
	UnsafeP, UnsafePhys   int
	UnsafePPerAllow       Bound
	UnsafePhysPerAllow    Bound
	UnsafePhysPerScenario Bound
}

// Summary holds the inputs to the preregistered H3 and H4 decision rules.
type Summary struct {
	H3        []Comparison
	H4        []H4Summary
	H4ByClass []H4ClassSummary
	// OutOfModel reports arm C on classes that violate declared assumptions.
	OutOfModel []H4Summary
}

func inModel(class string) bool {
	for _, c := range Classes() {
		if c.Name == class {
			return c.InModel
		}
	}
	return false
}

// Summarize derives the H3 comparison table and the stratified H4 bounds from rows.
func Summarize(rows []Row) Summary {
	var s Summary
	index := map[string]int{}
	for _, r := range rows {
		if r.Group != "fault" || !inModel(r.Class) {
			continue
		}
		i, ok := index[r.Class]
		if !ok {
			i = len(s.H3)
			index[r.Class] = i
			s.H3 = append(s.H3, Comparison{Class: r.Class, N: r.N, UnsafeP: map[string]int{},
				UnsafePhys: map[string]int{}, FalseRefusalPhys: map[string]int{}})
		}
		s.H3[i].UnsafeP[r.Arm] = r.UnsafeP
		s.H3[i].UnsafePhys[r.Arm] = r.UnsafePhys
		s.H3[i].FalseRefusalPhys[r.Arm] = r.FalseRefusalPhys
	}
	for i := range s.H3 {
		c := &s.H3[i]
		c.CBetterThanB0P = c.UnsafeP["C"] < c.UnsafeP["B0"]
		c.CBetterThanB1P = c.UnsafeP["C"] < c.UnsafeP["B1"]
		c.CBetterThanB0Phys = c.UnsafePhys["C"] < c.UnsafePhys["B0"]
		c.CBetterThanB1Phys = c.UnsafePhys["C"] < c.UnsafePhys["B1"]
		c.B1MatchesCUnsafeP = c.UnsafeP["C"] == c.UnsafeP["B1"]
		c.B1MatchesCUnsafePhys = c.UnsafePhys["C"] == c.UnsafePhys["B1"]
	}
	add := func(list []H4Summary, scope string, keep func(Row) bool) []H4Summary {
		h := H4Summary{Scope: scope}
		seenClass := map[string]bool{}
		allowByClass := map[string]int{}
		for _, r := range rows {
			if r.Arm == "C" && keep(r) {
				seenClass[r.Class] = true
				allowByClass[r.Class] += r.Allow
				h.Scenarios += r.N
				h.Allow += r.Allow
				h.UnsafeP += r.UnsafeP
				h.UnsafePhys += r.UnsafePhys
			}
		}
		h.Classes = len(seenClass)
		for class := range seenClass {
			if allowByClass[class] > 0 {
				h.ClassesWithAllow++
			} else {
				h.ClassesZeroAllow++
			}
		}
		h.UnsafePPerAllow = NewBound(h.UnsafeP, h.Allow, "P-violating among C ALLOW (circular: ALLOW<=>P)")
		h.UnsafePhysPerAllow = NewBound(h.UnsafePhys, h.Allow, "physically unsafe among C ALLOW")
		h.UnsafePhysPerScenario = NewBound(h.UnsafePhys, h.Scenarios, "unsafe C ALLOW among all scenarios")
		return append(list, h)
	}
	s.H4 = add(s.H4, "in-model, faulted scenarios only", func(r Row) bool { return inModel(r.Class) && r.Group == "fault" })
	s.H4 = add(s.H4, "in-model, all scenarios (faulted, controls, clean)", func(r Row) bool { return inModel(r.Class) })
	s.H4 = add(s.H4, "in-model, active-permit fault classes only (Allow > 0)", func(r Row) bool {
		return inModel(r.Class) && r.Group == "fault" && r.Allow > 0
	})
	s.H4 = add(s.H4, "in-model, total-refusal fault classes (Allow == 0, conditional risk Undefined 0/0)", func(r Row) bool {
		return inModel(r.Class) && r.Group == "fault" && r.Allow == 0
	})
	for _, r := range rows {
		if r.Arm != "C" || r.Group == "control" {
			continue
		}
		s.H4ByClass = append(s.H4ByClass, H4ClassSummary{
			Class:                 r.Class,
			Group:                 r.Group,
			InModel:               inModel(r.Class),
			TotalRefusal:          r.Allow == 0,
			Scenarios:             r.N,
			Allow:                 r.Allow,
			Deny:                  r.Deny,
			Defer:                 r.Defer,
			UnsafeP:               r.UnsafeP,
			UnsafePhys:            r.UnsafePhys,
			UnsafePPerAllow:       NewBound(r.UnsafeP, r.Allow, "P-violating among C ALLOW (circular: ALLOW<=>P)"),
			UnsafePhysPerAllow:    r.UnsafePhysPerAllow,
			UnsafePhysPerScenario: r.UnsafePhysPerScenario,
		})
	}
	for _, c := range Classes() {
		if c.InModel || c.Name == "CLEAN" {
			continue
		}
		name := c.Name
		s.OutOfModel = add(s.OutOfModel, name+" faulted", func(r Row) bool { return r.Class == name && r.Group == "fault" })
	}
	return s
}
