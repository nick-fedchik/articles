package evidence

import (
	"math"
	"reflect"
	"testing"
)

func corpusFixture() []CorpusItem {
	rows := []CorpusItem{}
	for _, doc := range []string{"RFC-100#sec-1", "RFC 200, Section 2", "public-document-3", "public-document-4"} {
		for n := 0; n < 20; n++ {
			rows = append(rows, CorpusItem{doc, "The controller shall store event " + string(rune('A'+n)), "OBLIGATION"})
		}
	}
	// Distinct documents must not share identical test sentences.
	for i := range rows {
		rows[i].Text += " " + rows[i].Document
	}
	return rows
}

func TestCorpusDocumentSplitAndReproducibility(t *testing.T) {
	rows := corpusFixture()
	first, err := EvaluateCorpus(rows, .2)
	if err != nil {
		t.Fatal(err)
	}
	for i, j := 0, len(rows)-1; i < j; i, j = i+1, j-1 {
		rows[i], rows[j] = rows[j], rows[i]
	}
	second, err := EvaluateCorpus(rows, .2)
	if err != nil || !reflect.DeepEqual(first, second) {
		t.Fatal("corpus order changed the measured result", err)
	}
	for _, fold := range first.Folds {
		if fold.TestDocumentSHA256 == fold.CalibrationDocumentSHA256 || fold.TrainRows+fold.CalibrationRows+fold.TestRows != first.Audit.UniqueRows {
			t.Fatal("document leakage")
		}
	}
}

func TestCorpusRejectsLeakageAndConflictingAnnotations(t *testing.T) {
	rows := corpusFixture()
	rows = append(rows, rows[0])
	clean, audit, err := PrepareCorpus(rows)
	if err != nil || audit.DuplicateRows != 1 || len(clean) != len(rows)-1 {
		t.Fatal(audit, err)
	}
	rows[len(rows)-1].Document = "another-document"
	if _, _, err := PrepareCorpus(rows); err == nil {
		t.Fatal("cross-document duplicate accepted")
	}
	rows[len(rows)-1] = rows[0]
	rows[len(rows)-1].Label = "NO_ATOM"
	if _, _, err := PrepareCorpus(rows); err == nil {
		t.Fatal("conflicting duplicate label accepted")
	}
	if CanonicalDocument("RFC-9110#sec-10") != CanonicalDocument("RFC 9110 (rfc9110.txt)") {
		t.Fatal("RFC sections were not grouped by document")
	}
}

func TestWilsonAndZeroDenominator(t *testing.T) {
	lo, hi, err := WilsonInterval(0, 100)
	if err != nil || lo > 1e-15 || math.Abs(hi-.0369934982) > 1e-9 {
		t.Fatal(lo, hi, err)
	}
	metrics := scorePredictions([]prediction{{[]float64{.2, .2, .2, .2, .2}, "NO_ATOM"}}, 1, 1)
	if metrics.Direct != 0 || metrics.RiskUpper != nil {
		t.Fatal("zero accepted examples reported a measured zero risk")
	}
	unapproved := approvedMetrics([]prediction{{[]float64{1, 0, 0, 0, 0}, "NO_ATOM"}}, 1, 1, false)
	if unapproved.Direct != 0 || unapproved.Refused != 1 || unapproved.RiskUpper != nil {
		t.Fatal("unapproved calibration allowed a saturated floating-point probability")
	}
	tie := []prediction{{[]float64{.2, .2, .2, .2, .2}, "OBLIGATION"}}
	if scoreArgmax(tie).Direct != 1 || scorePredictions(tie, 0, 0).Clarified != 1 {
		t.Fatal("deterministic argmax baseline was conflated with candidate-set clarification")
	}
}
