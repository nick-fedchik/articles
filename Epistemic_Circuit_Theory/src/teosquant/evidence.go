package teosquant

// Evidence is a deterministic set of synthetic examples, not observations.
type Evidence struct {
	Status      string              `json:"status"`
	Assumptions []string            `json:"assumptions"`
	Composition CompositionEvidence `json:"composition"`
	Measurement MeasurementEvidence `json:"measurement"`
	Errors      []ErrorEvidence     `json:"expected_errors"`
	Valve       []Metrics           `json:"valve_scenarios_A_B"`
	Throttle    []Metrics           `json:"throttle_D_1_9"`
}

type CompositionEvidence struct {
	A                      float64 `json:"a"`
	B                      float64 `json:"b"`
	GeneratorA             float64 `json:"generator_a"`
	GeneratorB             float64 `json:"generator_b"`
	Hamacher               float64 `json:"hamacher"`
	InverseGeneratorSum    float64 `json:"inverse_generator_sum"`
	IndependentConjunction float64 `json:"independent_conjunction"`
	EquivalentSequentialD  float64 `json:"equivalent_sequential_D"`
}

type MeasurementEvidence struct {
	Standard1KPa         float64 `json:"standard_1_kPa"`
	Standard2KPa         float64 `json:"standard_2_kPa"`
	Correlation          float64 `json:"correlation"`
	UncorrelatedVariance float64 `json:"uncorrelated_variance_kPa_squared"`
	CorrelatedVariance   float64 `json:"correlated_variance_kPa_squared"`
}

type ErrorEvidence struct {
	Transactions uint64  `json:"transactions"`
	Probability  float64 `json:"assumed_error_probability"`
	Expected     float64 `json:"expected_error_count"`
}

// SyntheticEvidence uses fixed, explicitly assumed inputs and no random draws.
func SyntheticEvidence() (Evidence, error) {
	g, err := Generator(0.8)
	if err != nil {
		return Evidence{}, err
	}
	h, err := Hamacher(0.8, 0.8)
	if err != nil {
		return Evidence{}, err
	}
	sumWeight, err := Weight(g + g)
	if err != nil {
		return Evidence{}, err
	}
	product := 0.8 * 0.8
	equivalentD, err := Generator(product)
	if err != nil {
		return Evidence{}, err
	}
	uncorrelated, err := SumVariance(1, 1, 0)
	if err != nil {
		return Evidence{}, err
	}
	correlated, err := SumVariance(1, 1, 0.8)
	if err != nil {
		return Evidence{}, err
	}
	e := Evidence{
		Status: "synthetic_conditional_calculations_not_empirical_calibration",
		Assumptions: []string{
			"w is a dimensionless model rating, not a calibrated probability",
			"independent_conjunction treats 0.8 as hypothetical event probabilities only for that comparison",
			"measurement is X+Y with standard uncertainties 1 kPa each and assumed correlation 0.8",
			"error expectation uses fixed marginal probability 0.01; independence is not required",
			"valve uses assumed U=2 or 6 kPa, tolerance width=10 kPa, w0=1 and phi0=1/tick",
			"policy rate is an optional throttle, not measured verifier throughput",
			"quadratic score has units tick^-2, not entropy/time or an error count",
		},
		Composition: CompositionEvidence{0.8, 0.8, g, g, h, sumWeight, product, equivalentD},
		Measurement: MeasurementEvidence{1, 1, 0.8, uncorrelated, correlated},
	}
	for _, n := range []uint64{100, 200, 400} {
		expected, err := ExpectedErrors(n, 0.01)
		if err != nil {
			return Evidence{}, err
		}
		e.Errors = append(e.Errors, ErrorEvidence{n, 0.01, expected})
	}
	for _, u := range []float64{2, 6} {
		d, err := MetrologyRating(u, 10)
		if err != nil {
			return Evidence{}, err
		}
		m, err := Evaluate(1, 1, 0, d, 0)
		if err != nil {
			return Evidence{}, err
		}
		e.Valve = append(e.Valve, m)
	}
	for _, d := range []float64{1, 9} {
		m, err := Evaluate(1, 1, 0, d, 0)
		if err != nil {
			return Evidence{}, err
		}
		e.Throttle = append(e.Throttle, m)
	}
	return e, nil
}
