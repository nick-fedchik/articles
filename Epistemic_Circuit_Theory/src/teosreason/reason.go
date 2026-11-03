// Package teosreason implements research protocol 01: a minimal three-level
// model of one inference (content representation, content transformation and
// epistemic characteristics). It is a bounded teaching model for one case, not
// a general reasoner, natural-language interpreter or theory of diagnosis.
package teosreason

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	evidence "github.com/nick-fedchik/articles/Epistemic_Circuit_Theory/src"
	"github.com/nick-fedchik/articles/Epistemic_Circuit_Theory/src/teoscompose"
)

// Atom is a statement pred(args...). Two atoms are the same statement iff the
// predicate and the ordered argument list are equal. Arguments starting with
// '?' are pattern variables.
type Atom struct {
	Pred string
	Args []string
}

func A(pred string, args ...string) Atom { return Atom{pred, args} }

func (a Atom) Key() string { return a.Pred + "(" + strings.Join(a.Args, ",") + ")" }

func (a Atom) valid() bool {
	if strings.TrimSpace(a.Pred) == "" || strings.HasPrefix(a.Pred, "?") {
		return false
	}
	for _, x := range a.Args {
		if strings.TrimSpace(x) == "" {
			return false
		}
	}
	return true
}

func (a Atom) ground() bool {
	for _, x := range a.Args {
		if strings.HasPrefix(x, "?") {
			return false
		}
	}
	return true
}

type bindings map[string]string

func unify(pattern, ground Atom, b bindings) (bindings, bool) {
	if pattern.Pred != ground.Pred || len(pattern.Args) != len(ground.Args) {
		return nil, false
	}
	out := bindings{}
	for k, v := range b {
		out[k] = v
	}
	for i, p := range pattern.Args {
		g := ground.Args[i]
		if !strings.HasPrefix(p, "?") {
			if p != g {
				return nil, false
			}
			continue
		}
		if prev, ok := out[p]; ok && prev != g {
			return nil, false
		}
		out[p] = g
	}
	return out, true
}

func substitute(a Atom, b bindings) Atom {
	args := make([]string, len(a.Args))
	for i, x := range a.Args {
		if v, ok := b[x]; ok {
			x = v
		}
		args[i] = x
	}
	return Atom{a.Pred, args}
}

// Observation is one reported value. Root is the primary source (instrument or
// original measurement act); Relay only forwards a copy. Reliability is carried
// as a characteristic and deliberately not used by any content transition.
type Observation struct {
	Atom        Atom
	Value       evidence.State
	Relay       string
	Root        string
	Reliability float64
}

// ObserveInterval converts an interval reading into a Belnap observation with
// the unchanged teoscompose.IntervalGuard. Adaptation fixed by the protocol:
// interval overlap (BOTH in the guard convention) means insufficient precision
// in the reasoning layer and becomes UNKNOWN, not a contradiction.
func ObserveInterval(atom Atom, m teoscompose.Measurement, window evidence.Interval, unit, relay, root string) (Observation, evidence.State, error) {
	raw, err := teoscompose.IntervalGuard(m, window, unit)
	if err != nil {
		return Observation{}, raw, err
	}
	v := raw
	if raw == evidence.Both {
		v = evidence.Unknown
	}
	return Observation{Atom: atom, Value: v, Relay: relay, Root: root}, raw, nil
}

// Scheme is a modeler-supplied explanation scheme (abductive, non-deductive).
// If PerReportingRoot names a variable, one hypothesis is proposed per primary
// source that reported the effect as TRUE.
type Scheme struct {
	ID               string
	Hypothesis       Atom
	Effect           Atom
	PerReportingRoot string
}

type Strength uint8

const (
	Typical Strength = iota
	Strong
)

func (s Strength) String() string {
	if s == Strong {
		return "strong"
	}
	return "typical"
}

// Rule: if hypothesis If holds, the model expects atom Then to have value
// Expected (TRUE or FALSE), counting only primary sources other than the one
// bound to ExcludeRootVar.
type Rule struct {
	ID             string
	If             Atom
	Then           Atom
	Expected       evidence.State
	Strength       Strength
	Requires       []string
	ExcludeRootVar string
}

type Model struct {
	Schemes []Scheme
	Rules   []Rule
	Active  map[string]bool
}

type QueryKind uint8

const (
	Why QueryKind = iota
	Whether
)

type Query struct {
	Kind   QueryKind
	Target Atom
}

type Relation string

const (
	NoData     Relation = "NO_DATA"
	Confirms   Relation = "CONFIRMS"
	Fails      Relation = "FAILS"
	Conflicted Relation = "CONFLICTED"
)

type Link struct {
	Hypothesis    string
	Rule          string
	Strength      Strength
	Atom          Atom
	Expected      evidence.State
	ExcludedRoots []string
	Assumptions   []string
	Observed      evidence.State
	Roots         []string
	Relation      Relation
}

type HypothesisResult struct {
	Hypothesis   Atom
	Scheme       string
	Status       evidence.State
	SupportRoots []string
	RefuteRoots  []string
	RefutedBy    []string
	Conflicts    []string
	Unmet        []string
	Links        []Link
}

type Check struct {
	Hypothesis    string
	Rule          string
	Atom          Atom
	Expected      evidence.State
	ExcludedRoots []string
}

type Answer struct {
	Kind                QueryKind
	Target              Atom
	Presupposition      evidence.State
	PresuppositionRoots []string
	Hypotheses          []HypothesisResult
	Compatible          []string
	Excluded            []string
	Contested           []string
	Direction           string
	NextChecks          []Check
	UnusedObservations  []string
	Assumptions         []string
	Label               string
}

const (
	LabelWhy          = "напрям перевірки серед розглянутих гіпотез; причину не доведено"
	LabelNoDirection  = "сумісні пояснення серед розглянутих гіпотез; напрям не визначено"
	LabelUnknown      = "передумову не встановлено"
	LabelFalse        = "передумова хибна"
	LabelContested    = "передумову оскаржено суперечливими підставами"
	LabelWhetherState = "стан твердження з першоджерелами"
)

// StepKind separates the nature of each transition in the trace.
type StepKind string

const (
	HumanInput     StepKind = "людський ввід"
	Characteristic StepKind = "характеристика"
	Refinement     StepKind = "уточнення цілі"
	Abduction      StepKind = "пропонування гіпотези (недедуктивно)"
	Deduction      StepKind = "дедукція в моделі"
	Matching       StepKind = "зіставлення"
	Selection      StepKind = "політика ранжування"
	Formation      StepKind = "формування відповіді"
)

type Step struct {
	N        int
	Operator string
	Kind     StepKind
	Basis    string
	Inputs   []string
	Outputs  []string
}

type store struct {
	byAtom map[string][]Observation
	used   map[string]bool
}

func (s store) join(atom Atom, excluded map[string]bool) (evidence.State, []string) {
	v := evidence.Unknown
	roots := map[string]bool{}
	for _, o := range s.byAtom[atom.Key()] {
		if excluded[o.Root] {
			continue
		}
		v |= o.Value
		if o.Value != evidence.Unknown {
			roots[o.Root] = true
		}
	}
	return v, sortedKeys(roots)
}

func (s store) rootValues(atom Atom) map[string]evidence.State {
	out := map[string]evidence.State{}
	for _, o := range s.byAtom[atom.Key()] {
		out[o.Root] |= o.Value
	}
	return out
}

func sortedKeys(m map[string]bool) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

func relation(observed, expected evidence.State) Relation {
	switch observed {
	case evidence.Unknown:
		return NoData
	case evidence.Both:
		return Conflicted
	case expected:
		return Confirms
	default:
		return Fails
	}
}

func validate(q Query, m Model, obs []Observation) error {
	if !q.Target.valid() || !q.Target.ground() || (q.Kind != Why && q.Kind != Whether) {
		return errors.New("reason: ground target and known query kind required")
	}
	for _, o := range obs {
		if !o.Atom.valid() || !o.Atom.ground() || !o.Value.Valid() || strings.TrimSpace(o.Root) == "" ||
			o.Reliability < 0 || o.Reliability > 1 {
			return fmt.Errorf("reason: invalid observation %s", o.Atom.Key())
		}
	}
	for _, s := range m.Schemes {
		if s.ID == "" || !s.Hypothesis.valid() || !s.Effect.valid() {
			return fmt.Errorf("reason: invalid scheme %q", s.ID)
		}
	}
	for _, r := range m.Rules {
		if r.ID == "" || !r.If.valid() || !r.Then.valid() || (r.Expected != evidence.True && r.Expected != evidence.False) {
			return fmt.Errorf("reason: invalid rule %q", r.ID)
		}
		bound := map[string]bool{}
		for _, x := range r.If.Args {
			bound[x] = true
		}
		for _, x := range r.Then.Args {
			if strings.HasPrefix(x, "?") && !bound[x] {
				return fmt.Errorf("reason: rule %q has unbound variable %s", r.ID, x)
			}
		}
		if r.ExcludeRootVar != "" && !bound[r.ExcludeRootVar] {
			return fmt.Errorf("reason: rule %q excludes unbound %s", r.ID, r.ExcludeRootVar)
		}
	}
	return nil
}

// Reason composes the operators as a fixed finite DAG:
// input -> join -> refine -> propose -> derive -> match -> status -> select -> answer.
func Reason(q Query, m Model, obs []Observation) (Answer, []Step, error) {
	if err := validate(q, m, obs); err != nil {
		return Answer{}, nil, err
	}
	var trace []Step
	emit := func(op string, kind StepKind, basis string, in, out []string) {
		trace = append(trace, Step{len(trace) + 1, op, kind, basis, in, out})
	}

	st := store{byAtom: map[string][]Observation{}, used: map[string]bool{}}
	var obsLines []string
	for _, o := range obs {
		st.byAtom[o.Atom.Key()] = append(st.byAtom[o.Atom.Key()], o)
		obsLines = append(obsLines, fmt.Sprintf("%s=%s від %s через %s", o.Atom.Key(), o.Value, o.Root, o.Relay))
	}
	sort.Strings(obsLines)
	emit("INPUT", HumanInput, "формалізація питання та спостережень людиною",
		[]string{"питання природною мовою"}, append([]string{"запит " + kindName(q.Kind) + "(" + q.Target.Key() + ")"}, obsLines...))

	var joined []string
	keys := make([]string, 0, len(st.byAtom))
	for k := range st.byAtom {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		v, roots := st.join(st.byAtom[k][0].Atom, nil)
		joined = append(joined, fmt.Sprintf("%s=%s за %v", k, v, roots))
	}
	emit("JOIN", Characteristic, "L1: об'єднання за порядком знання; першоджерела як множина", obsLines, joined)

	ans := Answer{Kind: q.Kind, Target: q.Target}
	ans.Presupposition, ans.PresuppositionRoots = st.join(q.Target, nil)
	st.used[q.Target.Key()] = true

	if q.Kind == Whether {
		ans.Label = LabelWhetherState
		emit("STATE", Formation, "тип відповіді Whether: стан одного твердження",
			[]string{q.Target.Key()}, []string{fmt.Sprintf("%s=%s за %v", q.Target.Key(), ans.Presupposition, ans.PresuppositionRoots)})
		ans.UnusedObservations = unused(st)
		return ans, trace, nil
	}

	goal := "explain(" + q.Target.Key() + ")"
	emit("REFINE", Refinement, "шаблон Why(e) → ціль explain(e) і передумова e=TRUE",
		[]string{"Why(" + q.Target.Key() + ")"}, []string{goal, "передумова " + q.Target.Key() + "=TRUE"})

	switch ans.Presupposition {
	case evidence.Unknown:
		ans.Label = LabelUnknown
	case evidence.False:
		ans.Label = LabelFalse
	case evidence.Both:
		ans.Label = LabelContested
	}
	if ans.Presupposition != evidence.True {
		emit("PRESUPPOSITION", Formation, "L4: гіпотези не висуваються",
			[]string{q.Target.Key()}, []string{ans.Label})
		ans.UnusedObservations = unused(st)
		return ans, trace, nil
	}

	type hyp struct {
		atom   Atom
		scheme string
	}
	hyps := map[string]hyp{}
	reporting := st.rootValues(q.Target)
	for _, s := range m.Schemes {
		b, ok := unify(s.Effect, q.Target, bindings{})
		if !ok {
			continue
		}
		var made []string
		if s.PerReportingRoot != "" {
			var roots []string
			for r, v := range reporting {
				if v == evidence.True {
					roots = append(roots, r)
				}
			}
			sort.Strings(roots)
			for _, r := range roots {
				b2 := bindings{s.PerReportingRoot: r}
				for k, v := range b {
					b2[k] = v
				}
				if h := substitute(s.Hypothesis, b2); h.ground() {
					hyps[h.Key()] = hyp{h, s.ID}
					made = append(made, h.Key())
				}
			}
		} else if h := substitute(s.Hypothesis, b); h.ground() {
			hyps[h.Key()] = hyp{h, s.ID}
			made = append(made, h.Key())
		}
		emit("PROPOSE", Abduction, "схема "+s.ID+" з каталогу модельєра", []string{goal}, made)
	}

	hkeys := make([]string, 0, len(hyps))
	for k := range hyps {
		hkeys = append(hkeys, k)
	}
	sort.Strings(hkeys)

	usedAssumptions := map[string]bool{}
	for _, hk := range hkeys {
		h := hyps[hk]
		res := HypothesisResult{Hypothesis: h.atom, Scheme: h.scheme, Status: evidence.Unknown}
		support, refute := map[string]bool{}, map[string]bool{}
		for _, r := range m.Rules {
			b, ok := unify(r.If, h.atom, bindings{})
			if !ok {
				continue
			}
			active := true
			for _, a := range r.Requires {
				if !m.Active[a] {
					active = false
				}
			}
			if !active {
				emit("DERIVE", Deduction, "правило "+r.ID+" неактивне: бракує припущень "+strings.Join(r.Requires, ","),
					[]string{hk}, nil)
				continue
			}
			for _, a := range r.Requires {
				usedAssumptions[a] = true
			}
			then := substitute(r.Then, b)
			excluded := map[string]bool{}
			if r.ExcludeRootVar != "" {
				excluded[b[r.ExcludeRootVar]] = true
			}
			ex := sortedKeys(excluded)
			pred := fmt.Sprintf("прогноз %s=%s без першоджерел %v", then.Key(), r.Expected, ex)
			emit("DERIVE", Deduction, fmt.Sprintf("правило %s (%s) за припущень %v", r.ID, r.Strength, r.Requires),
				[]string{hk}, []string{pred})

			observed, roots := st.join(then, excluded)
			st.used[then.Key()] = true
			link := Link{hk, r.ID, r.Strength, then, r.Expected, ex, append([]string(nil), r.Requires...), observed, roots, relation(observed, r.Expected)}
			res.Links = append(res.Links, link)
			emit("MATCH", Matching, "L2", []string{pred}, []string{fmt.Sprintf("%s: спостережено %s за %v", link.Relation, observed, roots)})

			switch {
			case link.Relation == Confirms:
				res.Status |= evidence.True
				for _, x := range roots {
					support[x] = true
				}
			case link.Relation == Fails && r.Strength == Strong:
				res.Status |= evidence.False
				res.RefutedBy = append(res.RefutedBy, r.ID)
				for _, x := range roots {
					refute[x] = true
				}
			case link.Relation == Fails:
				res.Unmet = append(res.Unmet, r.ID)
			case link.Relation == Conflicted:
				res.Conflicts = append(res.Conflicts, r.ID)
			}
		}
		res.SupportRoots, res.RefuteRoots = sortedKeys(support), sortedKeys(refute)
		emit("STATUS", Characteristic, "L3: біти «за/проти» лише з чистих зв'язків",
			[]string{hk}, []string{fmt.Sprintf("статус %s; за %v; проти %v; конфлікти %v; невиконані типові %v",
				res.Status, res.SupportRoots, res.RefuteRoots, res.Conflicts, res.Unmet)})
		ans.Hypotheses = append(ans.Hypotheses, res)
	}
	ans.Assumptions = sortedKeys(usedAssumptions)

	best, bestN, tie := "", -1, false
	for _, h := range ans.Hypotheses {
		k := h.Hypothesis.Key()
		if h.Status == evidence.False {
			ans.Excluded = append(ans.Excluded, k)
			continue
		}
		ans.Compatible = append(ans.Compatible, k)
		for _, l := range h.Links {
			if l.Relation == NoData {
				ans.NextChecks = append(ans.NextChecks, Check{k, l.Rule, l.Atom, l.Expected, l.ExcludedRoots})
			}
		}
		if h.Status == evidence.Both || len(h.Conflicts) > 0 {
			ans.Contested = append(ans.Contested, k)
			continue
		}
		if h.Status == evidence.True {
			switch n := len(h.SupportRoots); {
			case n > bestN:
				best, bestN, tie = k, n, false
			case n == bestN:
				tie = true
			}
		}
	}
	if best != "" && !tie {
		ans.Direction = best
		ans.Label = LabelWhy
	} else {
		ans.Label = LabelNoDirection
	}
	emit("SELECT", Selection, "C1–C3: виключення, оскарження, єдиний максимум першоджерел підтвердження",
		hkeys, []string{fmt.Sprintf("сумісні %v; виключені %v; оскаржені %v; напрям %q", ans.Compatible, ans.Excluded, ans.Contested, ans.Direction)})

	var checks []string
	for _, c := range ans.NextChecks {
		checks = append(checks, fmt.Sprintf("перевірити %s (очікується %s для %s за %s, без %v)", c.Atom.Key(), c.Expected, c.Hypothesis, c.Rule, c.ExcludedRoots))
	}
	ans.UnusedObservations = unused(st)
	emit("ANSWER", Formation, "відповідь типу «напрям перевірки»", []string{goal}, append([]string{ans.Label}, checks...))
	return ans, trace, nil
}

func unused(st store) []string {
	var out []string
	for k := range st.byAtom {
		if !st.used[k] {
			out = append(out, k)
		}
	}
	sort.Strings(out)
	return out
}

func kindName(k QueryKind) string {
	if k == Whether {
		return "Whether"
	}
	return "Why"
}
