package evidence

import (
	"math"
	"testing"
)

func TestVocabularyUsesAlternativeSegmentations(t *testing.T) {
	tokens := []VocabularyToken{
		{0, "O", math.Log(.1)}, {1, "OB", math.Log(.2)}, {2, "NO", math.Log(.3)}, {3, "other", math.Log(.4)},
	}
	summary, _, err := AnalyzeVocabulary(tokens, []string{"OBLIGATION", "NO_ATOM"}, "NO_ATOM")
	if err != nil || summary.AllowedTokens != 3 || math.Abs(summary.RejectedMass-.4) > 1e-12 ||
		math.Abs(summary.NormMass-.3) > 1e-12 {
		t.Fatal(summary, err)
	}
	if err := ValidateGrammarToken(tokens, 1, []string{"OBLIGATION", "NO_ATOM"}); err != nil {
		t.Fatal(err)
	}
	if err := ValidateGrammarToken(tokens, 3, []string{"OBLIGATION", "NO_ATOM"}); err == nil {
		t.Fatal("ignored grammar was not detected")
	}
	if err := ValidateGrammarToken(tokens, 99, []string{"OBLIGATION"}); err == nil {
		t.Fatal("missing sampled token accepted")
	}
	tokens[3].ID = 0
	if _, _, err := AnalyzeVocabulary(tokens, []string{"OBLIGATION", "NO_ATOM"}, "NO_ATOM"); err == nil {
		t.Fatal("duplicate token IDs accepted")
	}
}

func TestTruncatedVocabularyAndOverlappingLabelsRejected(t *testing.T) {
	tokens := []VocabularyToken{{0, "O", math.Log(.1)}}
	if _, _, err := AnalyzeVocabulary(tokens, []string{"OBLIGATION", "NO_ATOM"}, "NO_ATOM"); err == nil {
		t.Fatal("top-k list mistaken for full vocabulary")
	}

	if _, _, err := AnalyzeVocabulary(tokens, []string{"NO", "NO_ATOM"}, "NO_ATOM"); err == nil {
		t.Fatal("overlapping full label strings accepted")
	}
}

func FuzzVocabularySummary(f *testing.F) {
	f.Add(math.Float64bits(-1), "NO")
	f.Fuzz(func(t *testing.T, bits uint64, text string) {
		if len(text) > 10000 {
			t.Skip("bounded token piece")
		}
		summary, _, err := AnalyzeVocabulary([]VocabularyToken{{0, text, math.Float64frombits(bits)}},
			[]string{"OBLIGATION", "NO_ATOM"}, "NO_ATOM")
		if err == nil && (math.IsNaN(summary.RejectedMass) || summary.RejectedMass < 0 || summary.RejectedMass > 1) {
			t.Fatal("invalid successful probability summary")
		}
	})
}
