package annotation

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"reflect"
	"sort"
)

type Assignment struct {
	DocumentID, Family, Split string
}

type SplitPlan struct {
	Schema, CorpusDigest string
	Seed                 int64
	Assignments          []Assignment
}

func Prepare(c Corpus, seed int64) (SplitPlan, Annotation, Annotation, error) {
	digest, err := c.Digest()
	if err != nil {
		return SplitPlan{}, Annotation{}, Annotation{}, err
	}
	groups := map[string]bool{}
	for _, doc := range c.Documents {
		groups[doc.Family] = true
	}
	if len(groups) < 3 {
		return SplitPlan{}, Annotation{}, Annotation{}, errors.New("annotation split: at least three independent declared families required")
	}
	families := []string{}
	for family := range groups {
		families = append(families, family)
	}
	key := func(family string) string { return TextDigest(fmt.Sprintf("%d\x00%s", seed, family)) }
	sort.Slice(families, func(i, j int) bool {
		a, b := key(families[i]), key(families[j])
		if a != b {
			return a < b
		}
		return families[i] < families[j]
	})
	peripheral := max(1, len(families)/5)
	split := map[string]string{}
	for i, family := range families {
		value := "train"
		if i < peripheral {
			value = "test"
		} else if i < 2*peripheral {
			value = "calibration"
		}
		split[family] = value
	}
	plan := SplitPlan{Schema: Schema, CorpusDigest: digest, Seed: seed}
	for _, doc := range c.Documents {
		plan.Assignments = append(plan.Assignments, Assignment{doc.ID, doc.Family, split[doc.Family]})
	}
	sort.Slice(plan.Assignments, func(i, j int) bool { return plan.Assignments[i].DocumentID < plan.Assignments[j].DocumentID })
	a := Annotation{Schema: Schema, CorpusDigest: digest, Annotator: "FILL_ANNOTATOR_A", Status: "draft"}
	b := Annotation{Schema: Schema, CorpusDigest: digest, Annotator: "FILL_ANNOTATOR_B", Status: "draft"}
	items := append([]Item(nil), c.Items...)
	sort.Slice(items, func(i, j int) bool { return items[i].ID < items[j].ID })
	for _, item := range items {
		a.Records = append(a.Records, Record{ItemID: item.ID})
		b.Records = append(b.Records, Record{ItemID: item.ID})
	}
	return plan, a, b, nil
}

func ValidateAnnotation(c Corpus, annotation Annotation) (map[string]Label, error) {
	digest, err := c.Digest()
	if err != nil {
		return nil, err
	}
	if annotation.Schema != Schema || annotation.CorpusDigest != digest || !identity(annotation.Annotator) ||
		annotation.Status != "submitted" {
		return nil, errors.New("annotation: submitted status, annotator and matching corpus digest required; drafts are not data")
	}
	return recordsMap(c, annotation.Records, true)
}

type Count struct{ Matched, Total int }

type ExtractionAgreement struct {
	KindModality, Roles    Count
	Conditions, Exceptions Count
	Negations, Priorities  Count
	Parameters             Count
}

func (out *ExtractionAgreement) add(x, y Label) {
	count(&out.KindModality, x.Kind == y.Kind && x.Modality == y.Modality)
	if x.Kind != "NO_ATOM" || y.Kind != "NO_ATOM" {
		count(&out.Roles, x.Actor == y.Actor && x.Action == y.Action && x.Target == y.Target)
		count(&out.Conditions, reflect.DeepEqual(x.Conditions, y.Conditions))
		count(&out.Exceptions, reflect.DeepEqual(x.Exceptions, y.Exceptions))
		count(&out.Negations, reflect.DeepEqual(x.Negations, y.Negations))
		count(&out.Priorities, reflect.DeepEqual(x.PriorityOver, y.PriorityOver))
		count(&out.Parameters, reflect.DeepEqual(x.Parameters, y.Parameters))
	}
}

type Agreement struct {
	ExtractionAgreement
	Synthetic, IndependenceAttested bool
	Exact, Decisions                Count
	Kappa                           *float64
	KappaUndefined                  string
	Conflicts                       []string
}

func count(value *Count, agrees bool) {
	value.Total++
	if agrees {
		value.Matched++
	}
}

func Compare(c Corpus, a, b Annotation) (Agreement, error) {
	left, err := ValidateAnnotation(c, a)
	if err != nil {
		return Agreement{}, err
	}
	right, err := ValidateAnnotation(c, b)
	if err != nil {
		return Agreement{}, err
	}
	if a.Annotator == b.Annotator {
		return Agreement{}, errors.New("annotation agreement: two distinct annotator identities required")
	}
	out := Agreement{Synthetic: c.Synthetic, IndependenceAttested: !c.Synthetic &&
		a.Human && a.Blind && !a.ModelAssisted && b.Human && b.Blind && !b.ModelAssisted}
	categoriesA, categoriesB := map[string]int{}, map[string]int{}
	for _, item := range c.Items {
		x, y := left[item.ID], right[item.ID]
		equal := LabelsEqual(x, y)
		count(&out.Exact, equal)
		if !equal {
			out.Conflicts = append(out.Conflicts, item.ID)
		}
		if item.Task == "decision" {
			count(&out.Decisions, x.Decision == y.Decision)
			continue
		}
		catA, catB := x.Kind+":"+x.Modality, y.Kind+":"+y.Modality
		categoriesA[catA]++
		categoriesB[catB]++
		out.ExtractionAgreement.add(x, y)
	}
	sort.Strings(out.Conflicts)
	n := out.KindModality.Total
	if n == 0 {
		out.KappaUndefined = "no extraction items"
	} else {
		expected := 0.0
		for category, total := range categoriesA {
			expected += float64(total) * float64(categoriesB[category]) / (float64(n) * float64(n))
		}
		if math.Abs(1-expected) < 1e-15 {
			out.KappaUndefined = "degenerate category marginals"
		} else {
			value := (float64(out.KindModality.Matched)/float64(n) - expected) / (1 - expected)
			out.Kappa = &value
		}
	}
	return out, nil
}

func objectDigest(value any) (string, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", sha256.Sum256(data)), nil
}

func Adjudicate(c Corpus, a, b Annotation, adjudication Adjudication) (Gold, error) {
	agreement, err := Compare(c, a, b)
	if err != nil {
		return Gold{}, err
	}
	digest, err := c.Digest()
	if err != nil {
		return Gold{}, err
	}
	if adjudication.Schema != Schema || adjudication.CorpusDigest != digest || !identity(adjudication.Arbiter) ||
		adjudication.Arbiter == a.Annotator || adjudication.Arbiter == b.Annotator {
		return Gold{}, errors.New("adjudication: distinct arbiter and matching corpus required")
	}
	if !c.Synthetic && (!agreement.IndependenceAttested || !adjudication.Human || adjudication.ModelAssisted) {
		return Gold{}, errors.New("adjudication: real independent gold requires declared blind human annotators and a human arbiter")
	}
	left, err := ValidateAnnotation(c, a)
	if err != nil {
		return Gold{}, err
	}
	resolutions, err := recordsMap(c, adjudication.Resolutions, false)
	if err != nil {
		return Gold{}, err
	}
	if len(resolutions) != len(agreement.Conflicts) {
		return Gold{}, errors.New("adjudication: resolve exactly the disagreements, not agreed items")
	}
	for _, id := range agreement.Conflicts {
		if _, exists := resolutions[id]; !exists {
			return Gold{}, errors.New("adjudication: missing a disagreement resolution")
		}
	}
	gold := Gold{Schema: Schema, CorpusDigest: digest, Annotators: []string{a.Annotator, b.Annotator},
		Arbiter: adjudication.Arbiter, Synthetic: c.Synthetic,
		IndependenceAttested: agreement.IndependenceAttested && adjudication.Human && !adjudication.ModelAssisted}
	for _, value := range []any{a, b, adjudication} {
		hash, err := objectDigest(value)
		if err != nil {
			return Gold{}, err
		}
		gold.InputDigests = append(gold.InputDigests, hash)
	}
	ids := []string{}
	for _, item := range c.Items {
		ids = append(ids, item.ID)
	}
	sort.Strings(ids)
	for _, id := range ids {
		label := left[id]
		if resolution, exists := resolutions[id]; exists {
			label = resolution
		}
		gold.Records = append(gold.Records, Record{id, label})
	}
	if _, err := recordsMap(c, gold.Records, true); err != nil {
		return Gold{}, err
	}
	return gold, nil
}
