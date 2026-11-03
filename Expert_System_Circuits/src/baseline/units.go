package baseline

import "fmt"

// Pascal is a pressure in pascals; distinct named types prevent mixing raw
// numbers of different units in the conventional arms.
type Pascal float64

// Bar is a pressure in bar.
type Bar float64

// Pascals converts bar to pascals.
func (b Bar) Pascals() Pascal { return Pascal(float64(b) * 1e5) }

// Bar converts pascals to bar.
func (p Pascal) Bar() Bar { return Bar(float64(p) / 1e5) }

// pascalsPerUnit lists the pressure units the conventional stack accepts.
var pascalsPerUnit = map[string]float64{
	"Pa":  1,
	"kPa": 1e3,
	"bar": 1e5,
	"MPa": 1e6,
	"psi": 6894.757293168361,
}

// ParsePressure attaches a unit to a raw reading; unknown units are errors.
func ParsePressure(value float64, unit string) (Pascal, error) {
	factor, ok := pascalsPerUnit[unit]
	if !ok {
		return 0, fmt.Errorf("baseline: unknown pressure unit %q", unit)
	}
	return Pascal(value * factor), nil
}

// FromBar expresses a value given in bar in the named unit (generator helper).
func FromBar(value float64, unit string) (float64, error) {
	factor, ok := pascalsPerUnit[unit]
	if !ok {
		return 0, fmt.Errorf("baseline: unknown pressure unit %q", unit)
	}
	return value * 1e5 / factor, nil
}
