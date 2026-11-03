package baseline

import (
	"math"
	"reflect"
	"slices"
	"testing"
	"time"

	evidence "github.com/nick-fedchik/articles/Expert_System_Circuits/src"
	"github.com/nick-fedchik/articles/Expert_System_Circuits/src/kernel"
)

func near(a, b, tol float64) bool { return math.Abs(a-b) <= tol }

// Reference values: scipy.stats.beta.ppf(1-tail, x+1, n-x), rounded to 8 digits.
func TestClopperPearsonKnownValues(t *testing.T) {
	cases := []struct {
		x, n int
		tail float64
		want float64
	}{
		{0, 10, 0.025, 0.30849711},
		{1, 10, 0.025, 0.44501612},
		{5, 10, 0.025, 0.81291397},
		{0, 1000, 0.05, 0.00299125},
		{0, 1000, 0.025, 0.00368208},
		{3, 100, 0.025, 0.08517605},
	}
	for _, c := range cases {
		got, err := ClopperPearsonUpper(c.x, c.n, c.tail)
		if err != nil || !near(got, c.want, 6e-9) {
			t.Errorf("CP(%d,%d,%g) = %.8f, %v; want %.8f", c.x, c.n, c.tail, got, err, c.want)
		}
	}
	if got, _ := ClopperPearsonUpper(7, 7, 0.025); got != 1 {
		t.Fatalf("x=n must give 1, got %g", got)
	}
	for _, bad := range [][2]int{{0, 0}, {-1, 5}, {6, 5}} {
		if _, err := ClopperPearsonUpper(bad[0], bad[1], 0.025); err == nil {
			t.Fatalf("expected error for %v", bad)
		}
	}
}

// The bound p_U solves BinomialCDF(x; n, p_U) = tail and, equivalently,
// I_{p_U}(x+1, n-x) = 1 - tail (Beta quantile form).
func TestClopperPearsonMatchesBetaAndBinomial(t *testing.T) {
	for _, n := range []int{1, 2, 5, 17, 100, 1000, 12000} {
		for _, x := range []int{0, 1, 2, n / 3, n - 1} {
			if x < 0 || x >= n {
				continue
			}
			for _, tail := range []float64{0.025, 0.05} {
				p, err := ClopperPearsonUpper(x, n, tail)
				if err != nil {
					t.Fatal(err)
				}
				if cdf := BinomialCDF(x, n, p); !near(cdf, tail, 1e-9) {
					t.Errorf("n=%d x=%d: binomial CDF at bound = %g, want %g", n, x, cdf, tail)
				}
				if beta := RegIncBeta(p, float64(x+1), float64(n-x)); !near(beta, 1-tail, 1e-8) {
					t.Errorf("n=%d x=%d: I_p(x+1,n-x) = %.12f, want %g", n, x, beta, 1-tail)
				}
			}
		}
	}
	if RuleOfThree(1000) != 0.003 || !math.IsNaN(RuleOfThree(0)) {
		t.Fatal("rule of three")
	}
	if b := NewBound(0, 0, "none"); b.Defined {
		t.Fatal("bound with zero trials must be undefined")
	}
}

func TestRegIncBetaIdentities(t *testing.T) {
	// I_z(1, b) = 1 - (1-z)^b and I_z(a, 1) = z^a.
	for _, z := range []float64{0.001, 0.2, 0.5, 0.9, 0.999} {
		if got, want := RegIncBeta(z, 1, 7), 1-math.Pow(1-z, 7); !near(got, want, 1e-12) {
			t.Errorf("I_%g(1,7) = %g want %g", z, got, want)
		}
		if got, want := RegIncBeta(z, 3, 1), math.Pow(z, 3); !near(got, want, 1e-12) {
			t.Errorf("I_%g(3,1) = %g want %g", z, got, want)
		}
	}
}

func TestUnits(t *testing.T) {
	p, err := ParsePressure(0.8, "MPa")
	if err != nil || !near(float64(p.Bar()), 8, 1e-12) {
		t.Fatalf("0.8 MPa = %v bar, %v", p.Bar(), err)
	}
	if _, err := ParsePressure(1, "atm-ish"); err == nil {
		t.Fatal("unknown unit accepted")
	}
	v, _ := FromBar(8, "psi")
	if back, _ := ParsePressure(v, "psi"); !near(float64(back.Bar()), 8, 1e-12) {
		t.Fatal("psi round trip")
	}
	if Bar(2).Pascals() != 2e5 {
		t.Fatal("bar to pascal")
	}
}

func TestCaseSpecAndFaithfulAllow(t *testing.T) {
	w := World{Pressure: 5, Authorized: true, Demand: true}
	in := Faithful(classSeed(1, 0, 0), w)
	if !w.Safe() || len(kernel.CheckSafety(CaseSpec(), in)) != 0 {
		t.Fatalf("faithful safe world must satisfy P: %v", kernel.CheckSafety(CaseSpec(), in))
	}
	e := NewEvaluator()
	for _, arm := range Arms {
		o, err := e.Run(arm, in)
		if err != nil || o.Decision != Allow {
			t.Fatalf("arm %s on clean safe input: %v %v", arm, o, err)
		}
	}
	hot := World{Pressure: 5, Authorized: true, Demand: true, HotWork: true, Gas: true}
	in = Faithful(classSeed(1, 0, 1), hot)
	for _, arm := range Arms {
		if o, _ := e.Run(arm, in); o.Decision != Allow {
			t.Fatalf("arm %s must allow cleared hot work: %v", arm, o)
		}
	}
	for _, unsafe := range []World{
		{Pressure: 9, Authorized: true, Demand: true},
		{Pressure: 5, Authorized: true, Demand: true, Lock: true, Interlock: true},
		{Pressure: 5, Authorized: true, Demand: true, Interlock: true},
		{Pressure: 5, Authorized: true, Demand: true, HotWork: true},
		{Pressure: 5, Demand: true},
	} {
		in := Faithful(classSeed(1, 0, 2), unsafe)
		for _, arm := range Arms {
			if o, _ := e.Run(arm, in); o.Decision == Allow {
				t.Fatalf("arm %s allowed unsafe faithful world %+v", arm, unsafe)
			}
		}
	}
}

func TestArmDocumentedBehaviour(t *testing.T) {
	e := NewEvaluator()
	base := Faithful(classSeed(2, 0, 0), World{Pressure: 5, Authorized: true, Demand: true})
	// A: duplicated unverified copies inflate confidence over the 0.8 threshold.
	dup := base
	dup.Facts = slices.Clone(base.Facts)
	removeFact(&dup, AtomAuthorized)
	dup.Facts = append(dup.Facts, fact(AtomAuthorized, true, "relay-1", Now), fact(AtomAuthorized, true, "relay-2", Now))
	if c := ProjectNaive(dup).Confidence[AtomAuthorized]; !near(c, 0.5+0.5*(1-0.4*0.4), 1e-12) {
		t.Fatalf("noisy-OR confidence %g", c)
	}
	if ArmA(dup).Decision != Allow {
		t.Fatal("A must accept two unverified copies (documented naive behaviour)")
	}
	if o, _ := e.Run("C", dup); o.Decision != Defer {
		t.Fatalf("C must defer unverified support, got %v", o)
	}
	if o, _ := e.Run("B1", dup); o.Decision != Defer {
		t.Fatalf("B1 must defer unverified support, got %v", o)
	}
	// B0 rejects future timestamps (two-sided freshness) and converts units.
	future := base
	future.Measurements = slices.Clone(base.Measurements)
	future.Measurements[0].Time = Now + int64(time.Second)
	if o, _ := e.Run("B0", future); o.Decision != Deny {
		t.Fatalf("B0 must reject future timestamp, got %v", o)
	}
	mpa := base
	mpa.Measurements = slices.Clone(base.Measurements)
	mpa.Measurements[0].Low, mpa.Measurements[0].High, mpa.Measurements[0].Unit = 0.49, 0.51, "MPa"
	for arm, want := range map[string]Decision{"A": Allow, "B0": Allow, "B1": Allow, "B1S": Defer, "C": Defer} {
		if o, _ := e.Run(arm, mpa); o.Decision != want {
			t.Fatalf("arm %s on 5 bar in MPa: %v want %s", arm, o, want)
		}
	}
	mpa.Measurements[0].Low, mpa.Measurements[0].High = 0.89, 0.91 // 9 bar
	if ArmA(mpa).Decision != Allow {
		t.Fatal("A reads 0.9 as bar (documented unit blindness)")
	}
	if o, _ := e.Run("B0", mpa); o.Decision != Deny {
		t.Fatal("B0 must deny 9 bar given in MPa")
	}
}

func TestUnreferencedSpamAndCleanProof(t *testing.T) {
	e := NewEvaluator()
	base := Faithful(classSeed(9, 0, 0), World{Pressure: 5, Authorized: true, Demand: true})

	// F-01: stale and unverified spam on AtomUnreferencedNoise must not block B1 or C,
	// whereas B0's unscoped flag freshness check denies stale unreferenced telemetry.
	spammed := base
	spammed.Facts = slices.Clone(base.Facts)
	spammed.Facts = append(spammed.Facts,
		fact(AtomUnreferencedNoise, true, SourceLLM, Now-int64(30*time.Second)),
		fact(AtomUnreferencedNoise, false, SourceGateway, Now+int64(5*time.Second)),
	)
	if violated := kernel.CheckSafety(CaseSpec(), spammed); len(violated) != 0 {
		t.Fatalf("unreferenced spam must not violate P: %v", violated)
	}
	for _, arm := range []string{"B1", "C"} {
		if o, err := e.Run(arm, spammed); err != nil || o.Decision != Allow {
			t.Fatalf("arm %s must ignore unreferenced spam, got %v (%v)", arm, o, err)
		}
	}
	if o, err := e.Run("B0", spammed); err != nil || o.Decision != Deny {
		t.Fatalf("B0 unscoped freshness expected to deny stale unreferenced fact, got %v (%v)", o, err)
	}

	// F-02: transitive conflict taint in CleanProof.
	// 1) In CaseSpec (as in F04_CONTRARY), contradicted AtomLock (Both) derives positive
	// support for AtomInterlock in closure, but Clean[AtomLock] and Clean[AtomInterlock]
	// have no clean positive support.
	lockContrary := base
	lockContrary.Facts = slices.Clone(base.Facts)
	lockContrary.Facts = append(lockContrary.Facts, fact(AtomLock, true, SourcePLC, Now))
	outC, cert, err := ArmC(CaseSpec(), lockContrary)
	if err != nil || outC.Decision != Defer {
		t.Fatalf("expected C to defer on contradicted lock, got %v (%v)", outC, err)
	}
	if cert.Clean[AtomLock] != evidence.Unknown || cert.Clean[AtomInterlock]&evidence.True != 0 {
		t.Fatalf("expected no clean positive support on contradicted lock/interlock: %+v", cert.Clean)
	}

	// 2) If an evidence rule derives a permit condition (AtomDemand) from a contradicted
	// premise (AtomUnreferencedNoise == Both) without direct clean support, raw closure
	// has Closure[AtomDemand] == True (Norm.Result.Decision == ALLOW), but Clean[AtomDemand]
	// is Unknown and ArmC downgrades to Defer (DEFER_CONFLICT).
	spec := CaseSpec()
	spec.Theory.Evidence = append(slices.Clone(spec.Theory.Evidence), evidence.SignedRule{
		ID: "telemetry_implies_demand", Version: "1",
		Premises:   []evidence.SignedLiteral{{Atom: AtomUnreferencedNoise}},
		Conclusion: evidence.SignedLiteral{Atom: AtomDemand},
	})
	tainted := base
	tainted.Facts = slices.Clone(base.Facts)
	removeFact(&tainted, AtomDemand)
	tainted.Facts = append(tainted.Facts,
		fact(AtomUnreferencedNoise, true, SourcePLC, Now),
		fact(AtomUnreferencedNoise, false, SourceDCS, Now),
	)
	taintedOut, taintedCert, err := ArmC(spec, tainted)
	if err != nil {
		t.Fatal(err)
	}
	if taintedCert.Norm.Result.Decision != kernel.Allow || taintedCert.Clean[AtomDemand] != evidence.Unknown ||
		taintedOut.Decision != Defer || taintedOut.Raw != string(kernel.Conflict) {
		t.Fatalf("expected CleanProof to downgrade conflict-derived permit to DEFER_CONFLICT: out=%+v cert=%+v",
			taintedOut, taintedCert)
	}
}

func TestGenerateDeterministicAndEvaluate(t *testing.T) {
	a, b := Generate(7, 40, 10), Generate(7, 40, 10)
	if !reflect.DeepEqual(a, b) {
		t.Fatal("generator not deterministic")
	}
	if want := 40 + 12*50; len(a) != want {
		t.Fatalf("scenario count %d want %d", len(a), want)
	}
	records, err := Evaluate(a)
	if err != nil {
		t.Fatal(err)
	}
	rows := Aggregate(records)
	if len(rows) != (1+2*12)*len(Arms) {
		t.Fatalf("rows %d", len(rows))
	}
	for _, r := range rows {
		if r.Allow+r.Deny+r.Defer != r.N || r.Automatic != r.Allow+r.Deny {
			t.Fatalf("row counts inconsistent: %+v", r)
		}
		if r.Arm == "C" && r.UnsafeP != 0 {
			t.Fatalf("C must have zero P-unsafe permits: %+v", r)
		}
		if (r.Arm == "A" || r.Arm == "B0") && r.Defer != 0 {
			t.Fatalf("binary arm deferred: %+v", r)
		}
	}
	// Every faulted in-model scenario except F07 and F11 violates P by construction.
	for _, r := range records {
		if r.Group() == "fault" && r.Class != "F07_NORM_PRIORITY_CONFLICT" &&
			r.Class != "F11_UNREFERENCED_SPAM" && r.Class != "X11_OOM_TRUSTED_BIAS" && r.PHolds {
			t.Fatalf("faulted scenario satisfies P: %+v", r)
		}
	}
	s := Summarize(rows)
	if len(s.H3) != 11 || len(s.H4) != 4 || len(s.H4ByClass) != 13 || len(s.OutOfModel) != 1 {
		t.Fatalf("summary shape %d %d %d %d", len(s.H3), len(s.H4), len(s.H4ByClass), len(s.OutOfModel))
	}
	if s.H4[0].Classes != 11 || s.H4[0].ClassesWithAllow != 2 || s.H4[0].ClassesZeroAllow != 9 {
		t.Fatalf("unexpected H4 class stratification: %+v", s.H4[0])
	}
	if !s.H4[2].UnsafePhysPerAllow.Defined || s.H4[3].UnsafePhysPerAllow.Defined {
		t.Fatalf("expected active-permit scope defined and total-refusal scope undefined: %+v %+v", s.H4[2], s.H4[3])
	}
}

func TestSpecLinesAndAuditLog(t *testing.T) {
	for _, arm := range []string{"A", "B0", "B1", "C"} {
		n, err := SpecLines(arm)
		if err != nil || n <= 0 {
			t.Fatalf("spec lines %s: %d %v", arm, n, err)
		}
	}
	if _, err := SpecLines("Z"); err == nil {
		t.Fatal("unknown arm")
	}
	shared, err := SharedEngineLines()
	if err != nil || shared["kernel_total"] <= 0 || shared["conventional_units"] <= 0 || shared["conventional_support"] <= 0 {
		t.Fatalf("shared engine lines: %v %v", shared, err)
	}
	var log AuditLog
	first := log.Append(map[string]any{"x": 1})
	second := log.Append(map[string]any{"x": 1})
	if log.Len() != 2 || first <= 0 || second <= first {
		t.Fatalf("audit log append-only chain: %d %d %d", log.Len(), first, second)
	}
}

func TestEvaluateConcurrentEvaluators(t *testing.T) {
	scenarios := Generate(3, 5, 2)
	done := make(chan []Record, 2)
	for range 2 {
		go func() {
			r, err := Evaluate(scenarios)
			if err != nil {
				t.Error(err)
			}
			done <- r
		}()
	}
	if a, b := <-done, <-done; !reflect.DeepEqual(a, b) {
		t.Fatal("concurrent evaluations differ")
	}
}
