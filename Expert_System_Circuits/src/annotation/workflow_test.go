package annotation

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func demoForTest(t *testing.T) DemoData {
	t.Helper()
	demo, err := SyntheticDemo()
	if err != nil {
		t.Fatal(err)
	}
	return demo
}

func TestBlindPreparationAndFamilySplit(t *testing.T) {
	demo := demoForTest(t)
	plan, a, b, err := Prepare(demo.Corpus, 42)
	if err != nil {
		t.Fatal(err)
	}
	again, _, _, err := Prepare(demo.Corpus, 42)
	if err != nil || !reflect.DeepEqual(plan, again) {
		t.Fatalf("nondeterministic plan: %v", err)
	}
	seen := map[string]string{}
	for _, entry := range plan.Assignments {
		if previous := seen[entry.Family]; previous != "" && previous != entry.Split {
			t.Fatal("document version leakage")
		}
		seen[entry.Family] = entry.Split
	}
	for _, template := range []Annotation{a, b} {
		if template.Status != "draft" || template.Human || template.Blind || template.ModelAssisted {
			t.Fatal("template invents an attestation")
		}
		for _, record := range template.Records {
			if !reflect.DeepEqual(record.Label, Label{}) {
				t.Fatal("blind packet contains labels")
			}
		}
		if _, err := ValidateAnnotation(demo.Corpus, template); err == nil {
			t.Fatal("draft counted as submitted human data")
		}
	}
}

func TestAgreementAdjudicationAndNoInventedIndependence(t *testing.T) {
	demo := demoForTest(t)
	agreement, err := Compare(demo.Corpus, demo.A, demo.B)
	if err != nil {
		t.Fatal(err)
	}
	if agreement.Exact.Matched != 14 || agreement.Exact.Total != 15 ||
		!reflect.DeepEqual(agreement.Conflicts, []string{"norm-0"}) || agreement.IndependenceAttested {
		t.Fatalf("wrong synthetic agreement: %+v", agreement)
	}
	gold, err := Adjudicate(demo.Corpus, demo.A, demo.B, demo.Adjudication)
	if err != nil || gold.IndependenceAttested || !gold.Synthetic || len(gold.Records) != 15 {
		t.Fatalf("gold status: %+v, %v", gold, err)
	}
	demo.Adjudication.Resolutions = nil
	if _, err := Adjudicate(demo.Corpus, demo.A, demo.B, demo.Adjudication); err == nil {
		t.Fatal("unresolved disagreement accepted")
	}
	demo = demoForTest(t)
	demo.Corpus.Synthetic = false
	digest, err := demo.Corpus.Digest()
	if err != nil {
		t.Fatal(err)
	}
	demo.A.CorpusDigest, demo.B.CorpusDigest, demo.Adjudication.CorpusDigest = digest, digest, digest
	if _, err := Adjudicate(demo.Corpus, demo.A, demo.B, demo.Adjudication); err == nil {
		t.Fatal("machine-generated annotations accepted as real independent gold")
	}
	demo.A.Human, demo.A.Blind, demo.B.Human, demo.B.Blind, demo.Adjudication.Human = true, true, true, true, true
	gold, err = Adjudicate(demo.Corpus, demo.A, demo.B, demo.Adjudication)
	if err != nil || !gold.IndependenceAttested {
		t.Fatalf("declared human workflow rejected: %v", err)
	}
	demo.B.Annotator = demo.A.Annotator
	if _, err := Compare(demo.Corpus, demo.A, demo.B); err == nil {
		t.Fatal("same annotator accepted twice")
	}
}

func TestCorpusAndLabelRejections(t *testing.T) {
	for _, test := range []struct {
		name string
		edit func(*DemoData)
	}{
		{"duplicate text", func(d *DemoData) {
			d.Corpus.Items[1].Text = "  " + strings.ToUpper(d.Corpus.Items[0].Text) + "  "
			d.Corpus.Items[1].TextSHA256 = TextDigest(d.Corpus.Items[1].Text)
		}},
		{"text tampering", func(d *DemoData) { d.Corpus.Items[0].Text += "tampered" }},
		{"duplicate document across families", func(d *DemoData) { d.Corpus.Documents[1].SHA256 = d.Corpus.Documents[0].SHA256 }},
		{"missing annotation", func(d *DemoData) { d.A.Records = d.A.Records[1:] }},
		{"duplicate annotation", func(d *DemoData) { d.A.Records[1] = d.A.Records[0] }},
		{"bad source span", func(d *DemoData) { d.A.Records[0].Label.Actor.Span = Span{0, 100000} }},
		{"unknown priority", func(d *DemoData) { d.A.Records[0].Label.PriorityOver = []string{"missing"} }},
		{"cross-family priority", func(d *DemoData) { d.A.Records[0].Label.PriorityOver = []string{"norm-1"} }},
		{"numeric exponent", func(d *DemoData) {
			d.A.Records[0].Label.Parameters = []Parameter{{"x", "<=", "1e999999999", "s", Span{0, 1}}}
		}},
		{"nonfinite number", func(d *DemoData) {
			d.A.Records[0].Label.Parameters = []Parameter{{"x", "<=", "NaN", "s", Span{0, 1}}}
		}},
	} {
		t.Run(test.name, func(t *testing.T) {
			demo := demoForTest(t)
			test.edit(&demo)
			if _, err := ValidateAnnotation(demo.Corpus, demo.A); err == nil {
				t.Fatal("invalid annotation accepted")
			}
		})
	}
	if spanValid("АБ", Span{1, 2}) || !spanValid("АБ", Span{0, 2}) {
		t.Fatal("UTF-8 byte boundaries are not checked")
	}
}

func TestExactAgreementIgnoresOnlyFreeTextAndSetOrder(t *testing.T) {
	x := Label{Kind: "NO_ATOM", Rationale: "one", Evidence: []Span{{2, 4}, {0, 2}}}
	y := Label{Kind: "NO_ATOM", Rationale: "another", Evidence: []Span{{0, 2}, {2, 4}}}
	if !LabelsEqual(x, y) {
		t.Fatal("free-text rationale or span order creates a false conflict")
	}
	y.Kind = "FACT"
	if LabelsEqual(x, y) {
		t.Fatal("semantic difference ignored")
	}
}

func TestKappaDegenerateMarginalsAreExplicitlyUndefined(t *testing.T) {
	demo := demoForTest(t)
	items, records := []Item{}, []Record{}
	for i, item := range demo.Corpus.Items {
		if demo.A.Records[i].Label.Kind == "NO_ATOM" {
			items, records = append(items, item), append(records, demo.A.Records[i])
		}
	}
	demo.Corpus.Items = items
	digest, err := demo.Corpus.Digest()
	if err != nil {
		t.Fatal(err)
	}
	demo.A.CorpusDigest, demo.B.CorpusDigest = digest, digest
	demo.A.Records, demo.B.Records = records, records
	report, err := Compare(demo.Corpus, demo.A, demo.B)
	if err != nil || report.Kappa != nil || report.KappaUndefined == "" || report.Exact.Matched != len(items) {
		t.Fatalf("degenerate kappa reported as a number: %+v %v", report, err)
	}
}

func testEvaluationInputs(t *testing.T) (DemoData, SplitPlan, Gold, Predictions) {
	t.Helper()
	demo := demoForTest(t)
	plan, _, _, err := Prepare(demo.Corpus, 42)
	if err != nil {
		t.Fatal(err)
	}
	gold, err := Adjudicate(demo.Corpus, demo.A, demo.B, demo.Adjudication)
	if err != nil {
		t.Fatal(err)
	}
	docs := map[string]bool{}
	for _, assignment := range plan.Assignments {
		docs[assignment.DocumentID] = assignment.Split == "test"
	}
	labels := map[string]Label{}
	for _, record := range gold.Records {
		labels[record.ItemID] = record.Label
	}
	predictions := Predictions{Schema: Schema, CorpusDigest: plan.CorpusDigest, System: "synthetic-predictor"}
	for _, item := range demo.Corpus.Items {
		if docs[item.DocumentID] {
			label := labels[item.ID]
			if item.Task == "decision" {
				label.Decision = "ALLOW"
			}
			predictions.Records = append(predictions.Records, Prediction{ItemID: item.ID, Label: &label})
		}
	}
	return demo, plan, gold, predictions
}

func TestSelectiveRiskAndExplicitAbstention(t *testing.T) {
	demo, plan, gold, predictions := testEvaluationInputs(t)
	report, err := Evaluate(demo.Corpus, plan, gold, predictions, "test")
	if err != nil {
		t.Fatal(err)
	}
	if report.FalseAllow.Errors != report.FalseAllow.Trials || report.FalseAllow.Trials == 0 ||
		report.AutomaticRisk.Lower == nil || report.FalseDeny.Upper != nil || report.IndependenceAttested {
		t.Fatalf("wrong selective-risk denominators: %+v", report)
	}
	for i := range predictions.Records {
		predictions.Records[i].Abstain, predictions.Records[i].Label = true, nil
	}
	report, err = Evaluate(demo.Corpus, plan, gold, predictions, "test")
	if err != nil || report.AutomaticDecisions != 0 || report.AutomaticRisk.Upper != nil ||
		report.AutomaticRisk.Undefined == "" || report.Abstained != report.Items {
		t.Fatalf("abstention misreported as zero risk: %+v %v", report, err)
	}
	predictions.Records = predictions.Records[1:]
	if _, err := Evaluate(demo.Corpus, plan, gold, predictions, "test"); err == nil {
		t.Fatal("missing prediction silently treated as abstention")
	}
}

func TestEvaluationRejectsLeakageAndInvalidPredictions(t *testing.T) {
	for _, test := range []struct {
		name string
		edit func(*SplitPlan, *Predictions)
	}{
		{"wrong fingerprint", func(_ *SplitPlan, p *Predictions) { p.CorpusDigest = "wrong" }},
		{"duplicate", func(_ *SplitPlan, p *Predictions) { p.Records[1] = p.Records[0] }},
		{"ambiguous abstention", func(_ *SplitPlan, p *Predictions) { p.Records[0].Abstain = true }},
		{"unknown priority", func(_ *SplitPlan, p *Predictions) {
			for i := range p.Records {
				if p.Records[i].Label.Kind == "NORM" {
					p.Records[i].Label.PriorityOver = []string{"unknown"}
					break
				}
			}
		}},
		{"family leakage", func(plan *SplitPlan, _ *Predictions) {
			plan.Assignments[4].Split = "test"
			if plan.Assignments[0].Split == "test" {
				plan.Assignments[4].Split = "train"
			}
		}},
	} {
		t.Run(test.name, func(t *testing.T) {
			demo, plan, gold, predictions := testEvaluationInputs(t)
			test.edit(&plan, &predictions)
			if _, err := Evaluate(demo.Corpus, plan, gold, predictions, "test"); err == nil {
				t.Fatal("invalid evaluation accepted")
			}
		})
	}
}

func FuzzCorpusValidation(f *testing.F) {
	demo, err := SyntheticDemo()
	if err != nil {
		f.Fatal(err)
	}
	data, err := json.Marshal(demo.Corpus)
	if err != nil {
		f.Fatal(err)
	}
	f.Add(data)
	f.Add([]byte(`{"Schema":"ect.annotation.v1"}`))
	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) > 1<<20 {
			t.Skip()
		}
		var corpus Corpus
		if err := json.Unmarshal(data, &corpus); err != nil {
			return
		}
		if err := corpus.Validate(); err != nil {
			return
		}
		if _, err := corpus.Digest(); err != nil {
			t.Fatal(err)
		}
	})
}
