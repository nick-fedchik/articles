package teoscompose

import (
	"crypto/sha256"
	"encoding/json"
	"math"
	"os"
	"reflect"
	"testing"
	"time"

	evidence "github.com/nick-fedchik/articles/Expert_System_Circuits/src"
)

func origin(t *testing.T, id, text string) Origin {
	t.Helper()
	o, err := studyOrigin(id, text)
	if err != nil {
		t.Fatal(err)
	}
	return o
}

func signal(t *testing.T, claim string, state evidence.State, refs ...Origin) Signal {
	t.Helper()
	s, err := AssumeSignal(claim, state, 0.6, refs...)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func fixture(t *testing.T, state, modal evidence.State) (FactPort, NormPort) {
	t.Helper()
	s := signal(t, "observed load", state, origin(t, "sensor", "load=930 K"))
	m, err := NewMeasurement(930, evidence.Interval{Low: 920, High: 940}, "K")
	if err != nil {
		t.Fatal(err)
	}
	roles := evidence.Roles{Actor: "controller", Action: "run", Target: "turbine"}
	f, err := NewFact(s, m, roles, 0)
	if err != nil {
		t.Fatal(err)
	}
	n, err := NewNorm(evidence.Rule{Pattern: roles, Sort: "K", Window: evidence.Interval{Low: 0, High: 950},
		MaxAge: time.Second, Modal: modal}, origin(t, "norm", "permit load<=950 K"))
	if err != nil {
		t.Fatal(err)
	}
	return f, n
}

func TestAllLogicalCombinations(t *testing.T) {
	// Independent tables, in order UNKNOWN, FALSE, TRUE, BOTH.
	andWant := [4][4]evidence.State{{0, 1, 0, 1}, {1, 1, 1, 1}, {0, 1, 2, 3}, {1, 1, 3, 3}}
	joinWant := [4][4]evidence.State{{0, 1, 2, 3}, {1, 1, 3, 3}, {2, 3, 2, 3}, {3, 3, 3, 3}}
	gateWant := [4][4]evidence.State{{0, 1, 0, 0}, {1, 1, 1, 1}, {0, 1, 2, 3}, {0, 1, 3, 3}}
	notWant := [4]evidence.State{0, 2, 1, 3}
	aOrigin, bOrigin := origin(t, "A", "A"), origin(t, "B", "B")
	for a := evidence.Unknown; a <= evidence.Both; a++ {
		x := signal(t, "p", a, aOrigin)
		neg, err := Not(x)
		if err != nil || neg.State() != notWant[a] || len(neg.Provenance()) != 1 {
			t.Fatalf("NOT(%s): %+v %v", a, neg, err)
		}
		twice, err := Not(neg)
		if err != nil || twice.State() != a {
			t.Fatalf("double negation: %+v %v", twice, err)
		}
		for b := evidence.Unknown; b <= evidence.Both; b++ {
			y := signal(t, "p", b, bOrigin)
			conj, err := And(x, y)
			if err != nil || conj.State() != andWant[a][b] || !conj.State().Valid() || len(conj.Provenance()) != 2 {
				t.Fatalf("AND(%s,%s): %+v %v", a, b, conj, err)
			}
			joined, err := Merge(x, y)
			if err != nil || joined.State() != joinWant[a][b] || len(joined.Provenance()) != 2 {
				t.Fatalf("JOIN(%s,%s): %+v %v", a, b, joined, err)
			}
			fact, norm := fixture(t, a, b)
			got, err := Gate(fact, &norm)
			if err != nil || got.State != gateWant[a][b] || got.Permitted != (a == evidence.True && b == evidence.True) ||
				len(got.Provenance) != 2 {
				t.Fatalf("GATE(%s,%s): %+v %v", a, b, got, err)
			}
		}
	}
}

func TestProvenanceFanOutAndDedupComposition(t *testing.T) {
	s := signal(t, "p", evidence.True, origin(t, "A", "A"))
	copies, err := FanOut(s, 4)
	if err != nil {
		t.Fatal(err)
	}
	merged := copies[0]
	total := 0.0
	for _, c := range copies {
		total += c.Rating()
		if c.Claim() != s.Claim() || !reflect.DeepEqual(c.Provenance(), s.Provenance()) {
			t.Fatal("fan-out changed meaning or lineage")
		}
		merged, err = Merge(merged, c)
		if err != nil {
			t.Fatal(err)
		}
	}
	if math.Abs(total-2.4) > 1e-14 || merged.Rating() != 0.6 || len(merged.Provenance()) != 1 {
		t.Fatalf("unexpected copy accounting: sum=%v merged=%+v", total, merged)
	}
	refs := copies[0].Provenance()
	refs[0].ID = "mutated by caller"
	if copies[1].Provenance()[0].ID != "A" || s.Provenance()[0].ID != "A" {
		t.Fatal("mutable provenance escaped")
	}
	if _, err := Merge(s, signal(t, "q", evidence.True)); err == nil {
		t.Fatal("merged distinct propositions")
	}
	// Associativity of supported truth composition and knowledge merge.
	for a := evidence.Unknown; a <= evidence.Both; a++ {
		for b := evidence.Unknown; b <= evidence.Both; b++ {
			for c := evidence.Unknown; c <= evidence.Both; c++ {
				x, y, z := signal(t, "p", a), signal(t, "p", b), signal(t, "p", c)
				for _, op := range []func(Signal, Signal) (Signal, error){And, Merge} {
					xy, _ := op(x, y)
					yz, _ := op(y, z)
					left, errL := op(xy, z)
					right, errR := op(x, yz)
					if errL != nil || errR != nil || left.State() != right.State() {
						t.Fatalf("associativity %s %s %s: %v %v", a, b, c, errL, errR)
					}
				}
			}
		}
	}
}

func TestIntervalAndExistingNormativeSemantics(t *testing.T) {
	for _, row := range []struct {
		low, value, high float64
		unit             string
		want             evidence.State
	}{
		{920, 930, 940, "K", evidence.True},
		{950, 950, 950, "K", evidence.True},
		{940, 950, 960, "K", evidence.Both},
		{960, 970, 980, "K", evidence.False},
		{920, 930, 940, "Pa", evidence.Unknown},
	} {
		m, err := NewMeasurement(row.value, evidence.Interval{Low: row.low, High: row.high}, row.unit)
		if err != nil {
			t.Fatal(err)
		}
		got, err := IntervalGuard(m, evidence.Interval{Low: 0, High: 950}, "K")
		if err != nil || got != row.want {
			t.Fatalf("guard %+v: %s %v", row, got, err)
		}
	}
	for _, name := range []string{"stale", "role inversion", "wrong unit", "boundary", "outside"} {
		t.Run(name, func(t *testing.T) {
			fact, norm := fixture(t, evidence.True, evidence.True)
			switch name {
			case "stale":
				fact.fact.Age = 2 * time.Second
			case "role inversion":
				fact.fact.Roles.Actor, fact.fact.Roles.Target = fact.fact.Roles.Target, fact.fact.Roles.Actor
			case "wrong unit":
				fact.fact.Sort = "Pa"
			case "boundary":
				fact.fact.Range.High = 960
			case "outside":
				fact.fact.Range = evidence.Interval{Low: 960, High: 980}
			}
			want, err := evidence.EvaluateRule(norm.rule, fact.fact)
			if err != nil {
				t.Fatal(err)
			}
			got, err := Gate(fact, &norm)
			if err != nil || got.State != want.State || got.Permitted {
				t.Fatalf("changed existing semantics: want=%+v got=%+v err=%v", want, got, err)
			}
		})
	}
}

func TestMissingProvenanceAndNormCannotAuthorize(t *testing.T) {
	for _, missing := range []string{"fact", "norm", "norm port"} {
		fact, norm := fixture(t, evidence.True, evidence.True)
		n := &norm
		switch missing {
		case "fact":
			fact.signal = signal(t, "unanchored", evidence.True)
		case "norm":
			norm.origins = nil
		case "norm port":
			n = nil
		}
		got, err := Gate(fact, n)
		if err != nil || got.Permitted || got.State != evidence.Unknown || len(got.Reasons) == 0 {
			t.Fatalf("missing %s: %+v %v", missing, got, err)
		}
	}
}

func TestSourceIntegrityAndIdentityCollisions(t *testing.T) {
	doc := []byte("load=930 K")
	source, err := evidence.NewSource("sensor", "1", doc, 0.6)
	if err != nil {
		t.Fatal(err)
	}
	digest := sha256.Sum256(doc)
	for _, row := range []struct {
		doc        []byte
		start, end int
		digest     [32]byte
	}{
		{[]byte("load=980 K"), 0, len(doc), digest},
		{doc, 0, len(doc), sha256.Sum256([]byte("other quote"))},
		{doc, -1, len(doc), digest},
		{doc, 2, 2, digest},
	} {
		if _, err := NewOrigin(source, row.doc, row.start, row.end, row.digest); err == nil {
			t.Fatal("invalid source/quote accepted")
		}
	}
	a := signal(t, "p", evidence.True, origin(t, "sensor", "first"))
	b := signal(t, "p", evidence.True, origin(t, "sensor", "different same version"))
	if _, err := Merge(a, b); err == nil {
		t.Fatal("conflicting source version accepted")
	}
	if _, err := And(a, b); err == nil {
		t.Fatal("conflicting source version accepted in conjunction")
	}
	if _, err := AssumeSignal("p", evidence.True, 0.6, Origin{}); err == nil {
		t.Fatal("zero origin accepted")
	}
}

func TestTypeAndNumericDomainClosure(t *testing.T) {
	for _, row := range []struct {
		state  evidence.State
		rating float64
	}{
		{evidence.State(4), 0.6}, {evidence.True, math.NaN()}, {evidence.True, math.Inf(1)},
		{evidence.True, -0.1}, {evidence.True, 1.1},
	} {
		if _, err := AssumeSignal("p", row.state, row.rating); err == nil {
			t.Fatalf("invalid signal %+v", row)
		}
	}
	s := signal(t, "p", evidence.True)
	for _, n := range []int{0, -1} {
		if _, err := FanOut(s, n); err == nil {
			t.Fatal("invalid fan-out accepted")
		}
	}
	if _, err := And(Signal{}, s); err == nil {
		t.Fatal("zero signal accepted")
	}
	if _, err := Not(Signal{}); err == nil {
		t.Fatal("zero signal accepted")
	}
	if _, err := FanOut(Signal{}, 1); err == nil {
		t.Fatal("zero signal accepted by fan-out")
	}
	for _, bounds := range []evidence.Interval{
		{Low: 2, High: 1}, {Low: math.NaN(), High: 3}, {Low: 0, High: math.Inf(1)}, {Low: 2, High: 3},
	} {
		if _, err := NewMeasurement(1, bounds, "K"); err == nil {
			t.Fatal("invalid measurement accepted")
		}
	}
	if _, err := NewMeasurement(1, evidence.Interval{Low: 0, High: 2}, ""); err == nil {
		t.Fatal("missing unit accepted")
	}
	if _, err := NewFact(s, Measurement{}, evidence.Roles{}, 0); err == nil {
		t.Fatal("zero measurement accepted")
	}
	if _, err := NewNorm(evidence.Rule{}); err == nil {
		t.Fatal("zero normative rule accepted")
	}
	if _, err := Gate(FactPort{}, nil); err == nil {
		t.Fatal("zero fact accepted")
	}
	if _, err := IntervalGuard(Measurement{}, evidence.Interval{}, "K"); err == nil {
		t.Fatal("zero measurement accepted by guard")
	}
}

func TestResourceAccountingNotTruth(t *testing.T) {
	for _, row := range []struct {
		ledger Ledger
		want   bool
	}{
		{Ledger{Input: 1, Output: 1}, true},
		{Ledger{Input: 1, Output: 4}, false},
		{Ledger{Input: 1, Created: 3, Output: 4}, true},
		{Ledger{Before: 1, Output: 1}, true},
		{Ledger{Input: 3, Output: 1, Rejected: 1, After: 1}, true},
	} {
		got, err := row.ledger.Balanced()
		if err != nil || got != row.want {
			t.Fatalf("ledger %+v: %v %v", row.ledger, got, err)
		}
	}
	if _, err := (Ledger{Input: math.MaxUint64, Created: 1}).Balanced(); err == nil {
		t.Fatal("counter overflow accepted")
	}
	report, err := Study()
	if err != nil {
		t.Fatal(err)
	}
	for _, model := range report.Countermodels {
		if !model.LedgerBalanced || !model.SourceBytesIntact || model.EmpiricallyCorrect || !model.GatePermitted {
			t.Fatalf("lost countermodel: %+v", model)
		}
	}
	if report.Countermodels[0].EntailedByQuote || !report.Countermodels[1].EntailedByQuote {
		t.Fatal("interpretation and false-source countermodels were conflated")
	}
}

func TestFlowScoreIsNotPermission(t *testing.T) {
	score, err := FlowScore(2, 0.85, 0.95, 0.45, 10, 20)
	if err != nil || math.Abs(score-0.425) > 1e-14 {
		t.Fatalf("score %v %v", score, err)
	}
	fact, norm := fixture(t, evidence.True, evidence.True)
	norm.rule.Window.High = 925
	result, err := Gate(fact, &norm)
	if err != nil || result.Permitted || result.State != evidence.Both || score <= 0 {
		t.Fatalf("positive score incorrectly implies permission: %+v %v", result, err)
	}
	for _, params := range [][6]float64{
		{2, 0.85, 0.95, 0.45, 10, 0}, {2, 1.1, 0.95, 0.45, 10, 20},
		{2, 0.85, 0.95, 0.45, -1, 20}, {math.Inf(1), 0.85, 0.95, 0.45, 10, 20},
	} {
		if _, err := FlowScore(params[0], params[1], params[2], params[3], params[4], params[5]); err == nil {
			t.Fatal("invalid score parameters accepted")
		}
	}
}

func TestPublishedSyntheticOutput(t *testing.T) {
	report, err := Study()
	if err != nil {
		t.Fatal(err)
	}
	got, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	want, err := os.ReadFile("../testdata/teos-compose/results.json")
	if err != nil {
		t.Fatal(err)
	}
	if string(append(got, '\n')) != string(want) {
		t.Fatal("published output differs: run go run ./cmd/teoscompose")
	}
}
