package evidence

import (
	"errors"
	"math"
	"sort"
	"strings"
)

type Candidate struct {
	ID, Version   string
	Score, Radius float64
}

type Route struct {
	Disposition string
	Candidates  []Candidate
}

func SelectCandidates(candidates []Candidate, minScore, margin float64) (Route, error) {
	if !finite(minScore) || !finite(margin) || minScore < 0 || minScore > 1 || margin < 0 || margin > 1 {
		return Route{}, errors.New("routing: score threshold and margin must be in [0,1]")
	}
	seen := map[string]bool{}
	for _, c := range candidates {
		key := c.ID + "\x00" + c.Version
		if strings.TrimSpace(c.ID) == "" || strings.TrimSpace(c.Version) == "" ||
			strings.ContainsRune(c.ID+c.Version, '\x00') || seen[key] || !finite(c.Score) || !finite(c.Radius) ||
			c.Score < 0 || c.Score > 1 || c.Radius < 0 || c.Radius > 1 {
			return Route{}, errors.New("routing: invalid or duplicate candidate")
		}
		seen[key] = true
	}
	return selectValidatedCandidates(candidates, minScore, margin), nil
}

func selectValidatedCandidates(candidates []Candidate, minScore, margin float64) Route {
	bestLower := 0.0
	for _, c := range candidates {
		bestLower = math.Max(bestLower, math.Max(0, c.Score-c.Radius))
	}
	out := Route{Disposition: "REFUSE"}
	for _, c := range candidates {
		upper := math.Min(1, c.Score+c.Radius)
		if upper >= minScore && upper >= bestLower-margin {
			out.Candidates = append(out.Candidates, c)
		}
	}
	sort.Slice(out.Candidates, func(i, j int) bool {
		a, b := out.Candidates[i], out.Candidates[j]
		if a.Score != b.Score {
			return a.Score > b.Score
		}
		if a.ID != b.ID {
			return a.ID < b.ID
		}
		return a.Version < b.Version
	})
	if len(out.Candidates) > 0 {
		out.Disposition = "CLARIFY"
		if len(out.Candidates) == 1 && out.Candidates[0].Score-out.Candidates[0].Radius >= minScore {
			out.Disposition = "DIRECT"
		}
	}
	return out
}

type LogitSummary struct {
	RejectedMass   float64 `json:"rejected_mass"`
	NormMass       float64 `json:"norm_mass"`
	ConditionalMax float64 `json:"conditional_max"`
	Entropy        float64 `json:"conditional_normalized_entropy"`
}

func SummarizeLogits(logits []float64, allowed, normative []bool, temperature float64) (LogitSummary, error) {
	if len(logits) == 0 || len(logits) != len(allowed) || len(logits) != len(normative) ||
		!finite(temperature) || temperature <= 0 {
		return LogitSummary{}, errors.New("logits: invalid lengths or temperature")
	}
	maxZ := math.Inf(-1)
	allowedCount, normCount := 0, 0
	for i, z := range logits {
		if !finite(z) || normative[i] && !allowed[i] {
			return LogitSummary{}, errors.New("logits: nonfinite value or normative class outside grammar")
		}
		maxZ = math.Max(maxZ, z)
		if allowed[i] {
			allowedCount++
		}
		if normative[i] {
			normCount++
		}
	}
	if allowedCount == 0 || normCount == 0 {
		return LogitSummary{}, errors.New("logits: grammar and normative class sets must be nonempty")
	}
	probs := make([]float64, len(logits))
	total := 0.0
	for i, z := range logits {
		probs[i] = math.Exp((z - maxZ) / temperature)
		total += probs[i]
	}
	out := LogitSummary{}
	accepted, maxAllowed := 0.0, 0.0
	for i := range probs {
		probs[i] /= total
		if allowed[i] {
			accepted += probs[i]
			maxAllowed = math.Max(maxAllowed, probs[i])
		} else {
			out.RejectedMass += probs[i]
		}
		if normative[i] {
			out.NormMass += probs[i]
		}
	}
	if accepted == 0 {
		return LogitSummary{}, errors.New("logits: allowed probability mass underflow")
	}
	out.ConditionalMax = maxAllowed / accepted
	if allowedCount > 1 {
		for i, prob := range probs {
			if allowed[i] && prob > 0 {
				p := prob / accepted
				out.Entropy -= p * math.Log(p) / math.Log(float64(allowedCount))
			}
		}
	}
	return out, nil
}

func ZeroFailureUpperBound(trials int, alpha float64) (float64, error) {
	if trials < 1 || !finite(alpha) || alpha <= 0 || alpha >= 1 {
		return 0, errors.New("binomial bound: positive trials and alpha in (0,1) required")
	}
	return -math.Expm1(math.Log(alpha) / float64(trials)), nil
}
