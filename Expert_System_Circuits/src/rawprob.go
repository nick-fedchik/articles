package evidence

import (
	"errors"
	"math"
	"strings"
)

type VocabularyToken struct {
	ID      int
	Text    string
	LogProb float64
}

type VocabularySummary struct {
	RawMass, RejectedMass, NullMass, NormMass, ConditionalMax float64
	AllowedTokens, NullCompatibleTokens                       int
}

func ValidateGrammarToken(tokens []VocabularyToken, selectedID int, labels []string) error {
	for _, token := range tokens {
		if token.ID != selectedID {
			continue
		}
		if token.Text != "" {
			for _, label := range labels {
				if strings.HasPrefix(label, token.Text) {
					return nil
				}
			}
		}
		return errors.New("vocabulary: sampled token is outside the requested label grammar")
	}
	return errors.New("vocabulary: sampled token missing from full distribution")
}

// AnalyzeVocabulary includes every valid prefix segmentation, not only the tokenizer's canonical first token.
func AnalyzeVocabulary(tokens []VocabularyToken, labels []string, nullLabel string) (VocabularySummary, []float64, error) {
	if len(tokens) == 0 || len(labels) == 0 || nullLabel == "" {
		return VocabularySummary{}, nil, errors.New("vocabulary: nonempty tokens, labels and null label required")
	}
	seenLabels := map[string]bool{}
	for _, label := range labels {
		if label == "" || seenLabels[label] {
			return VocabularySummary{}, nil, errors.New("vocabulary: empty or duplicate label")
		}
		seenLabels[label] = true
	}
	for _, a := range labels {
		for _, b := range labels {
			if a != b && (strings.HasPrefix(a, b) || strings.HasPrefix(b, a)) {
				return VocabularySummary{}, nil, errors.New("vocabulary: labels must be prefix-free")
			}
		}
	}
	out := VocabularySummary{}
	probs := make([]float64, len(tokens))
	seen := make([]bool, len(tokens))
	maxAllowed := 0.0
	for _, token := range tokens {
		if token.ID < 0 || token.ID >= len(tokens) || seen[token.ID] ||
			math.IsNaN(token.LogProb) || math.IsInf(token.LogProb, 1) || token.LogProb > .001 {
			return VocabularySummary{}, nil, errors.New("vocabulary: incomplete IDs or invalid log probability")
		}
		seen[token.ID] = true
		prob := math.Exp(token.LogProb)
		probs[token.ID] = prob
		out.RawMass += prob
		allowed, normative, null := false, false, false
		if token.Text != "" {
			for _, label := range labels {
				if strings.HasPrefix(label, token.Text) {
					allowed = true
					normative = normative || label != nullLabel
					null = null || label == nullLabel
				}
			}
		}
		if allowed {
			out.AllowedTokens++
			maxAllowed = math.Max(maxAllowed, prob)
		} else {
			out.RejectedMass += prob
		}
		if normative {
			out.NormMass += prob
		}
		if null {
			out.NullCompatibleTokens++
			out.NullMass += prob
		}
	}
	if math.Abs(out.RawMass-1) > .001 {
		return VocabularySummary{}, nil, errors.New("vocabulary: probabilities do not cover the full normalized distribution")
	}
	allowedMass := out.RawMass - out.RejectedMass
	if allowedMass <= 0 {
		return VocabularySummary{}, nil, errors.New("vocabulary: allowed mass is numerically zero")
	}
	out.ConditionalMax = maxAllowed / allowedMass
	out.RejectedMass /= out.RawMass
	out.NormMass /= out.RawMass
	out.NullMass /= out.RawMass
	for i := range probs {
		probs[i] /= out.RawMass
	}
	return out, probs, nil
}
