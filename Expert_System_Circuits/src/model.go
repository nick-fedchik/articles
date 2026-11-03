package evidence

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"
)

type State uint8

const (
	Unknown State = iota
	False
	True
	Both
)

func (s State) Valid() bool { return s <= Both }

func (s State) String() string {
	switch s {
	case Unknown:
		return "UNKNOWN"
	case False:
		return "FALSE"
	case True:
		return "TRUE"
	case Both:
		return "BOTH"
	default:
		return fmt.Sprintf("INVALID(%d)", s)
	}
}

func TruthAnd(a, b State) (State, error) {
	if !a.Valid() || !b.Valid() {
		return Unknown, errors.New("truth conjunction: invalid state")
	}
	return (a & b & 2) | ((a | b) & 1), nil
}

func Priority(a, b State) (State, error) {
	if !a.Valid() || !b.Valid() {
		return Unknown, errors.New("priority conjunction: invalid state")
	}
	return priority(a, b), nil
}

func priority(a, b State) State {
	if a == False || b == False {
		return False
	}
	if a == Unknown || b == Unknown {
		return Unknown
	}
	if a == Both || b == Both {
		return Both
	}
	return True
}

func finite(x float64) bool { return !math.IsNaN(x) && !math.IsInf(x, 0) }

type Roles struct {
	Actor, Action, Target string
}

func (r Roles) validate(pattern bool) error {
	for _, slot := range []string{r.Actor, r.Action, r.Target} {
		if strings.TrimSpace(slot) == "" || (!pattern && slot == "*") {
			return errors.New("roles: empty or wildcard fact slot")
		}
	}
	return nil
}

type Interval struct {
	Low, High float64
}

func (i Interval) validate() error {
	if !finite(i.Low) || !finite(i.High) || i.Low > i.High {
		return errors.New("interval: finite ordered bounds required")
	}
	return nil
}

type Rule struct {
	Pattern Roles
	Sort    string
	Window  Interval
	MaxAge  time.Duration
	Modal   State
}

type Fact struct {
	Roles Roles
	Sort  string
	Range Interval
	Age   time.Duration
}

func (r Rule) validate() error {
	if err := r.Pattern.validate(true); err != nil {
		return err
	}
	if err := r.Window.validate(); err != nil {
		return err
	}
	if strings.TrimSpace(r.Sort) == "" || r.MaxAge <= 0 || !r.Modal.Valid() {
		return errors.New("rule: sort, positive max age and valid modality required")
	}
	return nil
}

func (f Fact) validate() error {
	if err := f.Roles.validate(false); err != nil {
		return err
	}
	if err := f.Range.validate(); err != nil {
		return err
	}
	if strings.TrimSpace(f.Sort) == "" || f.Age < 0 {
		return errors.New("fact: sort and nonnegative age required")
	}
	return nil
}

type Reason string

const (
	RoleMismatch  Reason = "ROLE_MISMATCH"
	RoleInversion Reason = "ROLE_INVERSION"
	SortMismatch  Reason = "SORT_MISMATCH"
	Boundary      Reason = "BOUNDARY"
	OutOfRange    Reason = "OUT_OF_RANGE"
	Stale         Reason = "STALE"
	ModalBlocked  Reason = "MODAL_BLOCKED"
)

type Decision struct {
	State   State
	Reasons []Reason
}

func EvaluateRule(rule Rule, fact Fact) (Decision, error) {
	if err := rule.validate(); err != nil {
		return Decision{}, err
	}
	if err := fact.validate(); err != nil {
		return Decision{}, err
	}
	out := Decision{State: True}
	apply := func(state State, reason Reason) {
		out.State = priority(out.State, state)
		if state != True {
			out.Reasons = append(out.Reasons, reason)
		}
	}
	bind := bindRoles(rule.Pattern, fact.Roles)
	reason := RoleMismatch
	if bind == False {
		reason = RoleInversion
	}
	apply(bind, reason)
	if rule.Sort != fact.Sort {
		apply(Unknown, SortMismatch)
	} else if fact.Range.High < rule.Window.Low || fact.Range.Low > rule.Window.High {
		apply(False, OutOfRange)
	} else if fact.Range.Low < rule.Window.Low || fact.Range.High > rule.Window.High {
		apply(Both, Boundary)
	}
	if fact.Age > rule.MaxAge {
		apply(False, Stale)
	}
	apply(rule.Modal, ModalBlocked)
	return out, nil
}

func bindRoles(pattern, fact Roles) State {
	match := func(p, f string) bool { return p == "*" || p == f }
	if match(pattern.Actor, fact.Actor) && match(pattern.Action, fact.Action) && match(pattern.Target, fact.Target) {
		return True
	}
	if pattern.Actor != "*" && pattern.Target != "*" && pattern.Actor != pattern.Target &&
		pattern.Actor == fact.Target && pattern.Target == fact.Actor && match(pattern.Action, fact.Action) {
		return False
	}
	return Unknown
}

type Source struct {
	ID, Version string
	Digest      [sha256.Size]byte
	Weight      float64
}

func NewSource(id, version string, content []byte, weight float64) (Source, error) {
	s := Source{ID: id, Version: version, Digest: sha256.Sum256(content), Weight: weight}
	return s, s.validate()
}

func (s Source) validate() error {
	if strings.TrimSpace(s.ID) == "" || strings.TrimSpace(s.Version) == "" ||
		strings.ContainsRune(s.ID+s.Version, '\x00') || !finite(s.Weight) ||
		s.Weight < 0 || s.Weight > 1 || (s.Weight > 0 && s.Weight < 1e-12) {
		return errors.New("source: identity and weight 0 or [1e-12,1] required")
	}
	return nil
}

func VerifyQuote(document []byte, start, end int, expected [sha256.Size]byte) error {
	if start < 0 || end <= start || end > len(document) {
		return errors.New("quote: invalid byte range")
	}
	if sha256.Sum256(document[start:end]) != expected {
		return errors.New("quote: digest mismatch")
	}
	return nil
}
