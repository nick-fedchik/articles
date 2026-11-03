package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"math"
	"os"
	"time"

	evidence "github.com/nick-fedchik/articles/Expert_System_Circuits/src"
	"github.com/nick-fedchik/articles/Expert_System_Circuits/src/internal/localdecoder"
)

type row struct {
	Fixture, Expected, Prefix, NullLabel string
	NullAllowed                          bool
	Summary                              evidence.VocabularySummary
	GrammarMaxProbabilityDifference      float64
	CanonicalNullTokens                  int
}

type sequence struct {
	Fixture, Prefix, Label string
	Tokens                 int
	LogProbability         float64
}

func vocabulary(response localdecoder.Completion) ([]evidence.VocabularyToken, error) {
	out := make([]evidence.VocabularyToken, len(response.Positions[0].Top))
	for i, token := range response.Positions[0].Top {
		if token.LogProb == nil {
			return nil, fmt.Errorf("native decoder returned null log probability")
		}
		out[i] = evidence.VocabularyToken{ID: token.ID, Text: string(token.Bytes), LogProb: *token.LogProb}
	}
	return out, nil
}

func run(ctx context.Context, endpoint string, vocab int) error {
	client, err := localdecoder.New(endpoint)
	if err != nil {
		return err
	}
	fixtures := []struct{ id, text, label string }{
		{"obligation", "The controller shall record each alarm.", "OBLIGATION"},
		{"prohibition", "The client shall not send an empty password.", "PROHIBITION"},
		{"description-en", "The controller recorded an alarm yesterday.", "NO_ATOM"},
		{"description-uk", "Учора сервер надіслав повідомлення клієнтові.", "NO_ATOM"},
	}
	normLabels := []string{"OBLIGATION", "PROHIBITION", "RECOMMENDATION", "PERMISSION"}
	rows, sequences := []row{}, []sequence{}
	for _, fixture := range fixtures {
		for _, prefix := range []string{"plain", "json"} {
			prompt := "Classify explicit normative modality as OBLIGATION, PROHIBITION, RECOMMENDATION, PERMISSION, or NO_ATOM for descriptions. Sentence: " + fixture.text
			if prefix == "plain" {
				prompt += "\nAnswer:"
			} else {
				prompt += "\n{\"modality\":\""
			}
			base, err := client.Next(ctx, prompt, "", vocab)
			if err != nil {
				return err
			}
			baseTokens, err := vocabulary(base)
			if err != nil {
				return err
			}
			labels := append(append([]string{}, normLabels...), "NO_ATOM")
			_, baseProbs, err := evidence.AnalyzeVocabulary(baseTokens, labels, "NO_ATOM")
			if err != nil {
				return err
			}
			nullTokens, err := client.Tokenize(ctx, "NO_ATOM", false)
			if err != nil {
				return err
			}
			for _, withNull := range []bool{false, true} {
				allowed := normLabels
				if withNull {
					allowed = labels
				}
				masked, err := client.Next(ctx, prompt, localdecoder.Grammar(allowed), vocab)
				if err != nil {
					return err
				}
				tokens, err := vocabulary(masked)
				if err != nil {
					return err
				}
				summary, probs, err := evidence.AnalyzeVocabulary(tokens, allowed, "NO_ATOM")
				if err != nil {
					return err
				}
				if err := evidence.ValidateGrammarToken(tokens, masked.Positions[0].ID, allowed); err != nil {
					return err
				}
				maxDiff := 0.0
				for i, prob := range probs {
					maxDiff = math.Max(maxDiff, math.Abs(prob-baseProbs[i]))
				}
				if maxDiff > .001 {
					return fmt.Errorf("grammar changed reported pre-mask distribution: max difference %g", maxDiff)
				}
				rows = append(rows, row{fixture.id, fixture.label, prefix, "NO_ATOM", withNull, summary, maxDiff, len(nullTokens)})
			}
			// A second refusal spelling explicitly checks tokenizer dependence.
			alias := append(append([]string{}, normLabels...), "NONE")
			aliasSummary, _, err := evidence.AnalyzeVocabulary(baseTokens, alias, "NONE")
			if err != nil {
				return err
			}
			aliasTokens, err := client.Tokenize(ctx, "NONE", false)
			if err != nil {
				return err
			}
			rows = append(rows, row{fixture.id, fixture.label, prefix, "NONE", true, aliasSummary, 0, len(aliasTokens)})
			if fixture.id != "obligation" && fixture.id != "description-en" {
				continue
			}
			promptIDs, err := client.Tokenize(ctx, prompt, true)
			if err != nil {
				return err
			}
			for _, label := range labels {
				joinedIDs, err := client.Tokenize(ctx, prompt+label, true)
				if err != nil {
					return err
				}
				if len(joinedIDs) <= len(promptIDs) {
					return fmt.Errorf("label tokenization has no continuation")
				}
				for i, id := range promptIDs {
					if joinedIDs[i] != id {
						return fmt.Errorf("tokenizer changes the prompt boundary; refusing to score a different prefix")
					}
				}
				labelIDs := joinedIDs[len(promptIDs):]
				contextIDs := append([]int{}, promptIDs...)
				logp := 0.0
				for _, tokenID := range labelIDs {
					next, err := client.Next(ctx, contextIDs, "", vocab)
					if err != nil {
						return err
					}
					tokens, err := vocabulary(next)
					if err != nil {
						return err
					}
					_, probs, err := evidence.AnalyzeVocabulary(tokens, labels, "NO_ATOM")
					if err != nil {
						return err
					}
					if tokenID < 0 || tokenID >= len(probs) || probs[tokenID] <= 0 {
						return fmt.Errorf("teacher-forced label probability underflow or invalid token")
					}
					logp += math.Log(probs[tokenID])
					contextIDs = append(contextIDs, tokenID)
				}
				sequences = append(sequences, sequence{fixture.id, prefix, label, len(labelIDs), logp})
			}
		}
	}
	result := struct {
		Vocabulary int
		Scope      string
		Rows       []row
		Sequences  []sequence
	}{vocab, "full pre-sampling normalized log probabilities, not raw unnormalized logits; synthetic probes; canonical sequence paths do not sum all segmentations", rows, sequences}
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(result)
}

func main() {
	endpoint := flag.String("endpoint", "http://127.0.0.1:18085", "temporary local native llama.cpp server")
	vocabulary := flag.Int("vocab", 32064, "required complete vocabulary size of the selected model")
	flag.Parse()
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Minute)
	defer cancel()
	if err := run(ctx, *endpoint, *vocabulary); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
