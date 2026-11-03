package teosbaseline

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

// H4Summary aggregates arm C on all in-model classes.
type H4Summary struct {
	Scope                 string
	Scenarios, Allow      int
	UnsafeP, UnsafePhys   int
	UnsafePPerAllow       Bound
	UnsafePhysPerAllow    Bound
	UnsafePhysPerScenario Bound
}

// Summary holds the inputs to the preregistered H3 and H4 decision rules.
type Summary struct {
	H3 []Comparison
	H4 []H4Summary
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

// Summarize derives the H3 comparison table and the H4 bounds from rows.
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
		for _, r := range rows {
			if r.Arm == "C" && keep(r) {
				h.Scenarios += r.N
				h.Allow += r.Allow
				h.UnsafeP += r.UnsafeP
				h.UnsafePhys += r.UnsafePhys
			}
		}
		h.UnsafePPerAllow = NewBound(h.UnsafeP, h.Allow, "P-violating among C ALLOW (circular: ALLOW<=>P)")
		h.UnsafePhysPerAllow = NewBound(h.UnsafePhys, h.Allow, "physically unsafe among C ALLOW")
		h.UnsafePhysPerScenario = NewBound(h.UnsafePhys, h.Scenarios, "unsafe C ALLOW among all scenarios")
		return append(list, h)
	}
	s.H4 = add(s.H4, "in-model, faulted scenarios only", func(r Row) bool { return inModel(r.Class) && r.Group == "fault" })
	s.H4 = add(s.H4, "in-model, all scenarios (faulted, controls, clean)", func(r Row) bool { return inModel(r.Class) })
	for _, c := range Classes() {
		if c.InModel || c.Name == "CLEAN" {
			continue
		}
		name := c.Name
		s.OutOfModel = add(s.OutOfModel, name+" faulted", func(r Row) bool { return r.Class == name && r.Group == "fault" })
	}
	return s
}
