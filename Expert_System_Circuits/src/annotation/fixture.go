package annotation

import (
	"fmt"
	"strings"
)

type DemoData struct {
	Corpus       Corpus
	A, B         Annotation
	Adjudication Adjudication
}

func SyntheticDemo() (DemoData, error) {
	corpus := Corpus{Schema: Schema, ID: "synthetic-workflow-only", Synthetic: true,
		Guidelines: "Generated workflow fixture only. Use literal actor/action/target spans, lowercase canonical role values; headings are NO_ATOM and explicit prohibitions are DENY."}
	expected := []Record{}
	for i := 0; i < 5; i++ {
		docID, family := fmt.Sprintf("doc-%d", i), fmt.Sprintf("family-%d", i)
		if i == 4 {
			family = "family-0"
		}
		normText := fmt.Sprintf("Client shall send message %d.", i)
		noteText := fmt.Sprintf("Fixture heading %d.", i)
		decisionText := fmt.Sprintf("Synthetic decision case %d: sending is prohibited.", i)
		docText := normText + "\n" + noteText + "\n" + decisionText
		corpus.Documents = append(corpus.Documents, Document{docID, family, "1", TextDigest(docText), "en", "public"})
		add := func(id, task, text, question string, label Label) {
			corpus.Items = append(corpus.Items, Item{id, docID, task, text, TextDigest(text), question})
			expected = append(expected, Record{id, label})
		}
		span := func(text string) Span {
			start := strings.Index(normText, text)
			return Span{start, start + len(text)}
		}
		norm := Label{Kind: "NORM", Modality: "OBLIGATION", Actor: Role{"client", span("Client")},
			Action: Role{"send", span("send")}, Target: Role{"message", span("message")},
			Rationale: "Generated synthetic fixture, not a human annotation."}
		add(fmt.Sprintf("norm-%d", i), "extract", normText, "", norm)
		add(fmt.Sprintf("note-%d", i), "extract", noteText, "", Label{Kind: "NO_ATOM", Rationale: "A generated heading without a proposition."})
		add(fmt.Sprintf("decision-%d", i), "decision", decisionText, "May the service send?",
			Label{Kind: "DECISION", Decision: "DENY", Evidence: []Span{{0, len(decisionText)}}, Rationale: "Synthetic answer, not independent gold."})
	}
	digest, err := corpus.Digest()
	if err != nil {
		return DemoData{}, err
	}
	a := Annotation{Schema: Schema, CorpusDigest: digest, Annotator: "synthetic-a", Status: "submitted", Records: expected}
	b := Annotation{Schema: Schema, CorpusDigest: digest, Annotator: "synthetic-b", Status: "submitted", Records: append([]Record(nil), expected...)}
	for i := range b.Records {
		b.Records[i].Label = canonicalLabel(b.Records[i].Label)
		if b.Records[i].ItemID == "norm-0" {
			b.Records[i].Label.Modality = "RECOMMENDATION"
		}
	}
	adjudication := Adjudication{Schema: Schema, CorpusDigest: digest, Arbiter: "synthetic-arbiter",
		Resolutions: []Record{expected[0]}}
	return DemoData{corpus, a, b, adjudication}, nil
}
