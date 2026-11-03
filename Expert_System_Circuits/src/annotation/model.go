package annotation

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"regexp"
	"slices"
	"sort"
	"strings"
	"unicode/utf8"
)

const Schema = "ect.annotation.v1"

var exactNumber = regexp.MustCompile(`^[+-]?(?:0|[1-9][0-9]*)(?:\.[0-9]+|/[1-9][0-9]*)?$`)

type Document struct {
	ID, Family, Version, SHA256, Language, Access string
}

type Item struct {
	ID, DocumentID, Task, Text, TextSHA256, Question string
}

type Corpus struct {
	Schema, ID, Guidelines string
	Synthetic              bool
	Documents              []Document
	Items                  []Item
}

type Span struct{ Start, End int }

type Role struct {
	Value string
	Span  Span
}

type Parameter struct {
	Variable, Relation, Value, Unit string
	Span                            Span
}

type Label struct {
	Kind, Modality, Decision, Rationale string
	Actor, Action, Target               Role
	Conditions, Exceptions, Negations   []Span
	Evidence                            []Span
	PriorityOver                        []string
	Parameters                          []Parameter
}

type Record struct {
	ItemID string
	Label  Label
}

type Annotation struct {
	Schema, CorpusDigest, Annotator, Status string
	Human, Blind, ModelAssisted             bool
	Records                                 []Record
}

type Adjudication struct {
	Schema, CorpusDigest, Arbiter string
	Human, ModelAssisted          bool
	Resolutions                   []Record
}

type Gold struct {
	Schema, CorpusDigest string
	Annotators           []string
	Arbiter              string
	Synthetic            bool
	IndependenceAttested bool
	InputDigests         []string
	Records              []Record
}

func identity(value string) bool {
	return strings.TrimSpace(value) != "" && !strings.ContainsAny(value, "\x00\r\n")
}

func rationaleValid(value string) bool {
	return strings.TrimSpace(value) != "" && len(value) <= 1<<20 && utf8.ValidString(value) && !strings.ContainsRune(value, '\x00')
}

func digestValid(value string) bool {
	raw, err := hex.DecodeString(value)
	return err == nil && len(raw) == sha256.Size && strings.ToLower(value) == value
}

func TextDigest(text string) string { return fmt.Sprintf("%x", sha256.Sum256([]byte(text))) }

func (c Corpus) Validate() error {
	if c.Schema != Schema || !identity(c.ID) || !rationaleValid(c.Guidelines) || len(c.Documents) == 0 || len(c.Documents) > 10000 ||
		len(c.Items) == 0 || len(c.Items) > 100000 {
		return errors.New("annotation corpus: schema, identity or bounded size invalid")
	}
	docs := map[string]Document{}
	hashFamilies := map[string]string{}
	for _, doc := range c.Documents {
		if !identity(doc.ID) || !identity(doc.Family) || !identity(doc.Version) || !identity(doc.Language) ||
			!digestValid(doc.SHA256) || (doc.Access != "public" && doc.Access != "permission" && doc.Access != "private") {
			return errors.New("annotation corpus: incomplete document provenance or access declaration")
		}
		if _, exists := docs[doc.ID]; exists {
			return errors.New("annotation corpus: duplicate document ID")
		}
		if family, exists := hashFamilies[doc.SHA256]; exists && family != doc.Family {
			return errors.New("annotation corpus: identical document hashes assigned to different families")
		}
		hashFamilies[doc.SHA256] = doc.Family
		docs[doc.ID] = doc
	}
	items, texts := map[string]bool{}, map[string]string{}
	textBytes := int64(len(c.Guidelines))
	for _, item := range c.Items {
		if !identity(item.ID) || items[item.ID] || item.Text == "" || len(item.Text) > 1<<20 || !utf8.ValidString(item.Text) ||
			TextDigest(item.Text) != item.TextSHA256 {
			return errors.New("annotation corpus: duplicate/invalid item or text fingerprint mismatch")
		}
		if _, ok := docs[item.DocumentID]; !ok {
			return errors.New("annotation corpus: unknown source document")
		}
		textBytes += int64(len(item.Text)) + int64(len(item.Question))
		if textBytes > 64<<20 || !utf8.ValidString(item.Question) {
			return errors.New("annotation corpus: total text budget exceeded or invalid question encoding")
		}
		if item.Task != "extract" && item.Task != "decision" ||
			item.Task == "decision" && !identity(item.Question) || item.Task == "extract" && item.Question != "" {
			return errors.New("annotation corpus: extraction/decision task or question invalid")
		}
		key := strings.Join(strings.Fields(strings.ToLower(item.Text)), " ")
		if previous, exists := texts[key]; exists {
			return fmt.Errorf("annotation corpus: normalized duplicate %q and %q requires human resolution", previous, item.ID)
		}
		items[item.ID], texts[key] = true, item.ID
	}
	return nil
}

func (c Corpus) Digest() (string, error) {
	if err := c.Validate(); err != nil {
		return "", err
	}
	c.Documents, c.Items = slices.Clone(c.Documents), slices.Clone(c.Items)
	sort.Slice(c.Documents, func(i, j int) bool { return c.Documents[i].ID < c.Documents[j].ID })
	sort.Slice(c.Items, func(i, j int) bool { return c.Items[i].ID < c.Items[j].ID })
	data, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", sha256.Sum256(data)), nil
}

func spanValid(text string, span Span) bool {
	return span.Start >= 0 && span.End > span.Start && span.End <= len(text) &&
		(span.Start == 0 || utf8.RuneStart(text[span.Start])) &&
		(span.End == len(text) || utf8.RuneStart(text[span.End]))
}

func decisionValid(value string) bool {
	switch value {
	case "ALLOW", "DENY", "DEFER_UNKNOWN", "DEFER_CONFLICT", "DEFER_NO_BASIS":
		return true
	default:
		return false
	}
}

func (c Corpus) validateLabel(item Item, label Label) error {
	if len(label.Conditions)+len(label.Exceptions)+len(label.Negations)+len(label.Evidence)+len(label.Parameters) > 1024 ||
		len(label.PriorityOver) > 1024 {
		return errors.New("annotation label: annotation size limit exceeded")
	}
	for _, spans := range [][]Span{label.Conditions, label.Exceptions, label.Negations, label.Evidence} {
		seen := map[Span]bool{}
		for _, span := range spans {
			if !spanValid(item.Text, span) || seen[span] {
				return errors.New("annotation label: invalid UTF-8 byte span or duplicate span")
			}
			seen[span] = true
		}
	}
	if item.Task == "decision" {
		if label.Kind != "DECISION" || !decisionValid(label.Decision) || !rationaleValid(label.Rationale) ||
			label.Modality != "" || label.Actor != (Role{}) || label.Action != (Role{}) || label.Target != (Role{}) ||
			len(label.Conditions)+len(label.Exceptions)+len(label.Negations)+len(label.PriorityOver)+len(label.Parameters) != 0 {
			return errors.New("annotation label: decision task requires only a grounded decision/rationale")
		}
		if label.Decision != "DEFER_NO_BASIS" && len(label.Evidence) == 0 {
			return errors.New("annotation label: decision evidence spans required")
		}
		return nil
	}
	if label.Decision != "" || label.Kind != "NORM" && label.Kind != "FACT" && label.Kind != "NO_ATOM" {
		return errors.New("annotation label: invalid extraction kind or decision on extraction task")
	}
	if label.Kind == "NO_ATOM" {
		if label.Modality != "" || label.Actor != (Role{}) || label.Action != (Role{}) || label.Target != (Role{}) ||
			len(label.Conditions)+len(label.Exceptions)+len(label.Negations)+len(label.PriorityOver)+len(label.Parameters) != 0 ||
			!rationaleValid(label.Rationale) {
			return errors.New("annotation label: NO_ATOM requires a rationale without extracted roles")
		}
		return nil
	}
	for _, role := range []Role{label.Actor, label.Action, label.Target} {
		if !identity(role.Value) || !spanValid(item.Text, role.Span) {
			return errors.New("annotation label: each role needs a canonical value and a source span")
		}
	}
	if label.Kind == "FACT" {
		if label.Modality != "" || len(label.PriorityOver) != 0 {
			return errors.New("annotation label: factual extraction has no normative modality/priority")
		}
	} else {
		switch label.Modality {
		case "OBLIGATION", "PROHIBITION", "PERMISSION", "RECOMMENDATION":
		default:
			return errors.New("annotation label: unknown normative modality")
		}
	}
	for _, parameter := range label.Parameters {
		if !identity(parameter.Variable) || !identity(parameter.Unit) || !spanValid(item.Text, parameter.Span) {
			return errors.New("annotation label: numeric parameter needs variable, unit and span")
		}
		if len(parameter.Value) > 128 || !exactNumber.MatchString(parameter.Value) {
			return errors.New("annotation label: numeric value must be a bounded decimal/fraction without an exponent")
		}
		if _, ok := new(big.Rat).SetString(parameter.Value); !ok {
			return errors.New("annotation label: numeric value must be an exact finite rational/decimal")
		}
		switch parameter.Relation {
		case "<", "<=", ">", ">=", "==":
		default:
			return errors.New("annotation label: invalid numeric relation")
		}
	}
	return nil
}

func canonicalLabel(label Label) Label {
	sortSpans := func(spans []Span) []Span {
		out := slices.Clone(spans)
		sort.Slice(out, func(i, j int) bool {
			if out[i].Start != out[j].Start {
				return out[i].Start < out[j].Start
			}
			return out[i].End < out[j].End
		})
		if len(out) == 0 {
			return nil
		}
		return out
	}
	label.Conditions, label.Exceptions = sortSpans(label.Conditions), sortSpans(label.Exceptions)
	label.Negations, label.Evidence = sortSpans(label.Negations), sortSpans(label.Evidence)
	label.PriorityOver = slices.Clone(label.PriorityOver)
	sort.Strings(label.PriorityOver)
	if len(label.PriorityOver) == 0 {
		label.PriorityOver = nil
	}
	label.Parameters = slices.Clone(label.Parameters)
	sort.Slice(label.Parameters, func(i, j int) bool {
		key := func(p Parameter) string {
			return fmt.Sprintf("%s\x00%s\x00%s\x00%s\x00%d\x00%d", p.Variable, p.Relation, p.Value, p.Unit, p.Span.Start, p.Span.End)
		}
		return key(label.Parameters[i]) < key(label.Parameters[j])
	})
	if len(label.Parameters) == 0 {
		label.Parameters = nil
	}
	return label
}

func recordsMap(c Corpus, records []Record, complete bool) (map[string]Label, error) {
	items := map[string]Item{}
	for _, item := range c.Items {
		items[item.ID] = item
	}
	families := map[string]string{}
	for _, doc := range c.Documents {
		families[doc.ID] = doc.Family
	}
	out := map[string]Label{}
	for _, record := range records {
		item, exists := items[record.ItemID]
		if !exists {
			return nil, errors.New("annotation records: unknown item ID")
		}
		if _, exists := out[record.ItemID]; exists {
			return nil, errors.New("annotation records: duplicate item ID")
		}
		if err := c.validateLabel(item, record.Label); err != nil {
			return nil, fmt.Errorf("item %q: %w", record.ItemID, err)
		}
		seenTargets := map[string]bool{}
		for _, target := range record.Label.PriorityOver {
			other, exists := items[target]
			if !exists || target == item.ID || seenTargets[target] || other.Task != "extract" ||
				families[other.DocumentID] != families[item.DocumentID] {
				return nil, errors.New("annotation label: priority reference must be distinct and grouped with its source")
			}
			seenTargets[target] = true
		}
		out[record.ItemID] = canonicalLabel(record.Label)
	}
	if complete && len(out) != len(items) {
		return nil, errors.New("annotation records: every item must have an explicit record")
	}
	for _, label := range out {
		for _, target := range label.PriorityOver {
			if other, exists := out[target]; exists && other.Kind != "NORM" {
				return nil, errors.New("annotation records: priority points to a non-normative item")
			}
		}
	}
	return out, nil
}

func LabelsEqual(a, b Label) bool {
	a, b = canonicalLabel(a), canonicalLabel(b)
	a.Rationale, b.Rationale = "", ""
	return reflect.DeepEqual(a, b)
}
