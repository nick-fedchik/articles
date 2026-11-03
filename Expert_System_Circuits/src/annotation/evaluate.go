package annotation

import (
	"errors"

	evidence "github.com/nick-fedchik/articles/Expert_System_Circuits/src"
)

type Prediction struct {
	ItemID  string
	Abstain bool
	Label   *Label
}

type Predictions struct {
	Schema, CorpusDigest, System string
	Records                      []Prediction
}

type Rate struct {
	Errors, Trials int
	Lower, Upper   *float64
	Undefined      string
}

func rate(errorsCount, trials int) (Rate, error) {
	result := Rate{Errors: errorsCount, Trials: trials}
	if trials == 0 {
		result.Undefined = "zero denominator; no measured zero risk"
		return result, nil
	}
	lo, hi, err := evidence.WilsonInterval(errorsCount, trials)
	if err != nil {
		return Rate{}, err
	}
	result.Lower, result.Upper = &lo, &hi
	return result, nil
}

type Evaluation struct {
	CorpusDigest, PlanDigest, GoldDigest string
	Split, System                        string
	Synthetic, IndependenceAttested      bool
	Items, Abstained                     int
	ExtractionItems, DecisionItems       int
	ExtractionExact                      Count
	ExtractionFields                     ExtractionAgreement
	DecisionExact                        Count
	AutomaticDecisions, Deferred         int
	FalseNormFromNoAtom                  int
	FalseAllow, FalseDeny                Rate
	AutomaticRisk                        Rate
	Scope                                string
}

func validatePlan(c Corpus, plan SplitPlan) (map[string]string, error) {
	digest, err := c.Digest()
	if err != nil {
		return nil, err
	}
	if plan.Schema != Schema || plan.CorpusDigest != digest || len(plan.Assignments) != len(c.Documents) {
		return nil, errors.New("evaluation: split plan does not match the corpus")
	}
	docs := map[string]Document{}
	for _, doc := range c.Documents {
		docs[doc.ID] = doc
	}
	out, families, counts := map[string]string{}, map[string]string{}, map[string]int{}
	for _, assignment := range plan.Assignments {
		doc, exists := docs[assignment.DocumentID]
		if !exists || assignment.Family != doc.Family || out[doc.ID] != "" ||
			assignment.Split != "train" && assignment.Split != "calibration" && assignment.Split != "test" {
			return nil, errors.New("evaluation: duplicate, unknown or invalid split assignment")
		}
		if prior := families[doc.Family]; prior != "" && prior != assignment.Split {
			return nil, errors.New("evaluation: document family leaks across splits")
		}
		out[doc.ID], families[doc.Family] = assignment.Split, assignment.Split
		counts[assignment.Split]++
	}
	if counts["train"] == 0 || counts["calibration"] == 0 || counts["test"] == 0 {
		return nil, errors.New("evaluation: nonempty train/calibration/test document groups required")
	}
	return out, nil
}

func Evaluate(c Corpus, plan SplitPlan, gold Gold, predictions Predictions, split string) (Evaluation, error) {
	assignments, err := validatePlan(c, plan)
	if err != nil {
		return Evaluation{}, err
	}
	digest, err := c.Digest()
	if err != nil {
		return Evaluation{}, err
	}
	if split != "train" && split != "calibration" && split != "test" ||
		gold.Schema != Schema || gold.CorpusDigest != digest || gold.Synthetic != c.Synthetic ||
		len(gold.Annotators) != 2 || !identity(gold.Annotators[0]) || !identity(gold.Annotators[1]) || gold.Annotators[0] == gold.Annotators[1] ||
		!identity(gold.Arbiter) || gold.Arbiter == gold.Annotators[0] || gold.Arbiter == gold.Annotators[1] ||
		len(gold.InputDigests) != 3 || c.Synthetic && gold.IndependenceAttested || !c.Synthetic && !gold.IndependenceAttested ||
		predictions.Schema != Schema || predictions.CorpusDigest != digest || !identity(predictions.System) {
		return Evaluation{}, errors.New("evaluation: malformed gold/prediction identity, lineage or selected split")
	}
	for _, hash := range gold.InputDigests {
		if !digestValid(hash) {
			return Evaluation{}, errors.New("evaluation: invalid annotation lineage fingerprint")
		}
	}
	labels, err := recordsMap(c, gold.Records, true)
	if err != nil {
		return Evaluation{}, err
	}
	selected := map[string]Item{}
	for _, item := range c.Items {
		if assignments[item.DocumentID] == split {
			selected[item.ID] = item
		}
	}
	if len(selected) == 0 || len(predictions.Records) != len(selected) {
		return Evaluation{}, errors.New("evaluation: provide an explicit prediction or abstention for every selected item")
	}
	predictedRecords := []Record{}
	seenPredictions := map[string]bool{}
	for _, prediction := range predictions.Records {
		if _, exists := selected[prediction.ItemID]; !exists || seenPredictions[prediction.ItemID] ||
			prediction.Abstain == (prediction.Label != nil) {
			return Evaluation{}, errors.New("evaluation: duplicate/out-of-split prediction or ambiguous abstention")
		}
		seenPredictions[prediction.ItemID] = true
		if prediction.Label != nil {
			predictedRecords = append(predictedRecords, Record{prediction.ItemID, *prediction.Label})
		}
	}
	validatedPredictions, err := recordsMap(c, predictedRecords, false)
	if err != nil {
		return Evaluation{}, err
	}
	out := Evaluation{Split: split, System: predictions.System, Synthetic: c.Synthetic,
		IndependenceAttested: gold.IndependenceAttested,
		Scope:                "exact structured labels; independence is metadata attestation, not verified human behavior; Wilson intervals assume independent accepted trials"}
	out.CorpusDigest = digest
	out.PlanDigest, err = objectDigest(plan)
	if err != nil {
		return Evaluation{}, err
	}
	out.GoldDigest, err = objectDigest(gold)
	if err != nil {
		return Evaluation{}, err
	}
	seen := map[string]bool{}
	falseAllows, allows, falseDenies, denies, automaticErrors := 0, 0, 0, 0, 0
	for _, prediction := range predictions.Records {
		item, exists := selected[prediction.ItemID]
		if !exists || seen[item.ID] || prediction.Abstain == (prediction.Label != nil) {
			return Evaluation{}, errors.New("evaluation: duplicate/out-of-split prediction or ambiguous abstention")
		}
		seen[item.ID] = true
		out.Items++
		if item.Task == "decision" {
			out.DecisionItems++
		} else {
			out.ExtractionItems++
		}
		if prediction.Abstain {
			out.Abstained++
			if item.Task == "decision" {
				out.Deferred++
			}
			continue
		}
		expected, actual := labels[item.ID], validatedPredictions[item.ID]
		if item.Task == "extract" {
			count(&out.ExtractionExact, LabelsEqual(expected, actual))
			out.ExtractionFields.add(expected, actual)
			if expected.Kind == "NO_ATOM" && actual.Kind == "NORM" {
				out.FalseNormFromNoAtom++
			}
			continue
		}
		count(&out.DecisionExact, expected.Decision == actual.Decision)
		switch actual.Decision {
		case "ALLOW":
			allows++
			out.AutomaticDecisions++
			if expected.Decision != "ALLOW" {
				falseAllows++
				automaticErrors++
			}
		case "DENY":
			denies++
			out.AutomaticDecisions++
			if expected.Decision != "DENY" {
				falseDenies++
				automaticErrors++
			}
		default:
			out.Deferred++
		}
	}
	out.FalseAllow, err = rate(falseAllows, allows)
	if err != nil {
		return Evaluation{}, err
	}
	out.FalseDeny, err = rate(falseDenies, denies)
	if err != nil {
		return Evaluation{}, err
	}
	out.AutomaticRisk, err = rate(automaticErrors, out.AutomaticDecisions)
	return out, err
}
