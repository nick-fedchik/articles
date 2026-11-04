package evidence

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"math"
	"regexp"
	"sort"
	"strings"
	"unicode"
)

var ModalityLabels = [...]string{"OBLIGATION", "PROHIBITION", "RECOMMENDATION", "PERMISSION", "NO_ATOM"}

type CorpusItem struct {
	Document, Text, Label string
}

type CorpusAudit struct {
	InputRows, UniqueRows, DuplicateRows, Documents int
}

var rfcDocument = regexp.MustCompile(`(?i)\brfc[- ]?([0-9]+)\b`)

func CanonicalDocument(reference string) string {
	if match := rfcDocument.FindStringSubmatch(reference); len(match) > 0 {
		return "rfc:" + match[1]
	}
	return strings.ToLower(strings.TrimSpace(strings.SplitN(reference, "#", 2)[0]))
}

func normalizedText(text string) string {
	return strings.Join(strings.Fields(strings.ToLower(text)), " ")
}

func PrepareCorpus(items []CorpusItem) ([]CorpusItem, CorpusAudit, error) {
	audit := CorpusAudit{InputRows: len(items)}
	documents, texts := map[string]bool{}, map[string]CorpusItem{}
	clean := []CorpusItem{}
	for _, item := range items {
		item.Document = CanonicalDocument(item.Document)
		text := normalizedText(item.Text)
		validLabel := false
		for _, label := range ModalityLabels {
			validLabel = validLabel || label == item.Label
		}
		if item.Document == "" || text == "" || !validLabel {
			return nil, audit, errors.New("corpus: document, text and known label required")
		}
		if previous, ok := texts[text]; ok {
			if previous.Label != item.Label || previous.Document != item.Document {
				return nil, audit, errors.New("corpus: conflicting label or cross-document duplicate; resolve grouping before evaluation")
			}
			audit.DuplicateRows++
			continue
		}
		texts[text] = item
		documents[item.Document] = true
		clean = append(clean, item)
	}
	audit.UniqueRows, audit.Documents = len(clean), len(documents)
	if audit.Documents < 3 {
		return nil, audit, errors.New("corpus: at least three disjoint document groups required for train/calibration/test")
	}
	sort.Slice(clean, func(i, j int) bool {
		if clean[i].Document != clean[j].Document {
			return clean[i].Document < clean[j].Document
		}
		return normalizedText(clean[i].Text) < normalizedText(clean[j].Text)
	})
	return clean, audit, nil
}

func textFeatures(text string) map[string]int {
	features := map[string]int{}
	words := strings.FieldsFunc(strings.ToLower(text), func(r rune) bool { return !unicode.IsLetter(r) && !unicode.IsDigit(r) })
	for i, word := range words {
		features["w:"+word]++
		if i > 0 {
			features["b:"+words[i-1]+" "+word]++
		}
		runes := []rune("^" + word + "$")
		for k := 0; k+3 <= len(runes); k++ {
			features["c:"+string(runes[k:k+3])]++
		}
	}
	return features
}

type textClassifier struct {
	counts [5]map[string]int
	total  [5]int
	rows   [5]int
	vocab  map[string]bool
}

func trainClassifier(items []CorpusItem) *textClassifier {
	m := &textClassifier{vocab: map[string]bool{}}
	for i := range m.counts {
		m.counts[i] = map[string]int{}
	}
	for _, item := range items {
		for class, label := range ModalityLabels {
			if label != item.Label {
				continue
			}
			m.rows[class]++
			for feature, count := range textFeatures(item.Text) {
				m.vocab[feature] = true
				m.counts[class][feature] += count
				m.total[class] += count
			}
		}
	}
	return m
}

func (m *textClassifier) logits(text string) []float64 {
	out := make([]float64, len(ModalityLabels))
	features := textFeatures(text)
	// Sorted iteration makes floating-point accumulation reproducible across runs.
	keys := []string{}
	for feature := range features {
		if m.vocab[feature] {
			keys = append(keys, feature)
		}
	}
	sort.Strings(keys)
	for class := range out {
		out[class] = math.Log(float64(m.rows[class] + 1))
		for _, feature := range keys {
			out[class] += float64(features[feature]) * math.Log(float64(m.counts[class][feature]+1)/float64(m.total[class]+len(m.vocab)))
		}
	}
	return out
}

func probabilities(logits []float64, temperature float64) []float64 {
	top := logits[0]
	for _, x := range logits {
		top = math.Max(top, x)
	}
	out := make([]float64, len(logits))
	total := 0.0
	for i, x := range logits {
		out[i] = math.Exp((x - top) / temperature)
		total += out[i]
	}
	for i := range out {
		out[i] /= total
	}
	return out
}

type prediction struct {
	probs []float64
	label string
}

func calibrateTemperature(model *textClassifier, calibration []CorpusItem) float64 {
	best, loss := 1., math.Inf(1)
	for _, temperature := range []float64{.5, 1, 2, 4, 8, 16} {
		total := 0.0
		for _, item := range calibration {
			probs := probabilities(model.logits(item.Text), temperature)
			for i, label := range ModalityLabels {
				if label == item.Label {
					total -= math.Log(math.Max(probs[i], 1e-300))
				}
			}
		}
		if total < loss {
			best, loss = temperature, total
		}
	}
	return best
}

func routePrediction(p prediction, threshold, margin float64) Route {
	candidates := make([]Candidate, len(p.probs))
	for i, probability := range p.probs {
		candidates[i] = Candidate{ID: ModalityLabels[i], Version: "1", Score: probability}
	}
	return selectValidatedCandidates(candidates, threshold, margin)
}

func WilsonInterval(errorsCount, trials int) (float64, float64, error) {
	if trials < 1 || errorsCount < 0 || errorsCount > trials {
		return 0, 0, errors.New("Wilson: valid counts and positive denominator required")
	}
	lo, hi := wilsonInterval(errorsCount, trials)
	return lo, hi, nil
}

func wilsonInterval(errorsCount, trials int) (float64, float64) {
	const z = 1.959963984540054
	n, p := float64(trials), float64(errorsCount)/float64(trials)
	denom := 1 + z*z/n
	center := (p + z*z/(2*n)) / denom
	half := z * math.Sqrt(p*(1-p)/n+z*z/(4*n*n)) / denom
	return math.Max(0, center-half), math.Min(1, center+half)
}

type CorpusMetrics struct {
	Rows, Direct, Correct, Errors, Clarified, Refused, NormativeFalseAccepts int
	RiskLower, RiskUpper                                                     *float64
}

func scorePredictions(predictions []prediction, threshold, margin float64) CorpusMetrics {
	return scoreRoutes(predictions, func(p prediction) Route { return routePrediction(p, threshold, margin) })
}

func scoreArgmax(predictions []prediction) CorpusMetrics {
	return scoreRoutes(predictions, func(p prediction) Route {
		best := 0
		for i, probability := range p.probs {
			if probability > p.probs[best] {
				best = i
			}
		}
		return Route{Disposition: "DIRECT", Candidates: []Candidate{{ID: ModalityLabels[best]}}}
	})
}

func scoreRoutes(predictions []prediction, choose func(prediction) Route) CorpusMetrics {
	out := CorpusMetrics{Rows: len(predictions)}
	for _, p := range predictions {
		route := choose(p)
		switch route.Disposition {
		case "DIRECT":
			out.Direct++
			if route.Candidates[0].ID == p.label {
				out.Correct++
			} else {
				out.Errors++
			}
			if p.label == "NO_ATOM" && route.Candidates[0].ID != "NO_ATOM" {
				out.NormativeFalseAccepts++
			}
		case "CLARIFY":
			out.Clarified++
		case "REFUSE":
			out.Refused++
		}
	}
	if out.Direct > 0 {
		lo, hi := wilsonInterval(out.Errors, out.Direct)
		out.RiskLower, out.RiskUpper = &lo, &hi
	}
	return out
}

type CorpusFold struct {
	TestDocumentSHA256, CalibrationDocumentSHA256 string
	TrainRows, CalibrationRows, TestRows          int
	Temperature                                   float64
	Threshold, Margin                             float64
	CalibrationApproved                           bool
	Argmax, Selective                             CorpusMetrics
}

type RiskPoint struct {
	Threshold, Margin float64
	Metrics           CorpusMetrics
}

type CorpusReport struct {
	Audit      CorpusAudit
	Folds      []CorpusFold
	Argmax     CorpusMetrics
	Selective  CorpusMetrics
	TargetRisk float64
	Scope      string
	RiskCurve  []RiskPoint
}

func EvaluateCorpus(items []CorpusItem, targetRisk float64) (CorpusReport, error) {
	if !finite(targetRisk) || targetRisk <= 0 || targetRisk >= 1 {
		return CorpusReport{}, errors.New("corpus: target risk must be in (0,1)")
	}
	clean, audit, err := PrepareCorpus(items)
	if err != nil {
		return CorpusReport{}, err
	}
	documents := []string{}
	for _, item := range clean {
		if len(documents) == 0 || documents[len(documents)-1] != item.Document {
			documents = append(documents, item.Document)
		}
	}
	report := CorpusReport{Audit: audit, TargetRisk: targetRisk, Scope: "document-disjoint cross-validation; existing local annotations, not an independent annotator; Wilson intervals assume independent accepted trials"}
	allArgmax, allSelective := []prediction{}, []CorpusMetrics{}
	for fold, testDocument := range documents {
		calDocument := documents[(fold+1)%len(documents)]
		train, calibration, test := []CorpusItem{}, []CorpusItem{}, []CorpusItem{}
		for _, item := range clean {
			switch item.Document {
			case testDocument:
				test = append(test, item)
			case calDocument:
				calibration = append(calibration, item)
			default:
				train = append(train, item)
			}
		}
		if len(train) == 0 || len(calibration) == 0 || len(test) == 0 {
			return CorpusReport{}, errors.New("corpus: empty document split")
		}
		model := trainClassifier(train)
		temperature := calibrateTemperature(model, calibration)
		predict := func(rows []CorpusItem) []prediction {
			out := make([]prediction, len(rows))
			for i, item := range rows {
				out[i] = prediction{probabilities(model.logits(item.Text), temperature), item.Label}
			}
			return out
		}
		calPredictions := predict(calibration)
		bestThreshold, bestMargin, bestCount := 0., 0., -1
		for _, threshold := range []float64{0, .3, .5, .6, .7, .8, .9, .95, .99} {
			for _, margin := range []float64{0, .05, .1, .2} {
				metrics := scorePredictions(calPredictions, threshold, margin)
				if metrics.Direct >= 10 && metrics.RiskUpper != nil && *metrics.RiskUpper <= targetRisk && metrics.Direct > bestCount {
					bestThreshold, bestMargin, bestCount = threshold, margin, metrics.Direct
				}
			}
		}
		testPredictions := predict(test)
		plain := scoreArgmax(testPredictions)
		selective := approvedMetrics(testPredictions, bestThreshold, bestMargin, bestCount >= 0)
		report.Folds = append(report.Folds, CorpusFold{
			TestDocumentSHA256:        fmt.Sprintf("%x", sha256.Sum256([]byte(testDocument))),
			CalibrationDocumentSHA256: fmt.Sprintf("%x", sha256.Sum256([]byte(calDocument))),
			TrainRows:                 len(train), CalibrationRows: len(calibration), TestRows: len(test),
			Temperature: temperature, Threshold: bestThreshold, Margin: bestMargin,
			CalibrationApproved: bestCount >= 0,
			Argmax:              plain, Selective: selective,
		})
		allArgmax = append(allArgmax, testPredictions...)
		allSelective = append(allSelective, selective)
	}
	report.Argmax = scoreArgmax(allArgmax)
	for _, threshold := range []float64{0, .5, .7, .9, .99} {
		for _, margin := range []float64{0, .1} {
			report.RiskCurve = append(report.RiskCurve, RiskPoint{threshold, margin, scorePredictions(allArgmax, threshold, margin)})
		}
	}
	for _, m := range allSelective {
		report.Selective.Rows += m.Rows
		report.Selective.Direct += m.Direct
		report.Selective.Correct += m.Correct
		report.Selective.Errors += m.Errors
		report.Selective.Clarified += m.Clarified
		report.Selective.Refused += m.Refused
		report.Selective.NormativeFalseAccepts += m.NormativeFalseAccepts
	}
	if report.Selective.Direct > 0 {
		lo, hi := wilsonInterval(report.Selective.Errors, report.Selective.Direct)
		report.Selective.RiskLower, report.Selective.RiskUpper = &lo, &hi
	}
	return report, nil
}

func approvedMetrics(predictions []prediction, threshold, margin float64, approved bool) CorpusMetrics {
	if !approved {
		return CorpusMetrics{Rows: len(predictions), Refused: len(predictions)}
	}
	return scorePredictions(predictions, threshold, margin)
}
