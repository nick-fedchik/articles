package evidence

import (
	"errors"
	"math"
	"time"
)

type ResponseBudget struct {
	FTTI, MaxAge, Polling, Compute, Transport, Actuation time.Duration
}

func (b ResponseBudget) Remaining() (time.Duration, error) {
	if b.FTTI <= 0 || b.MaxAge <= 0 {
		return 0, errors.New("budget: positive FTTI and max age required")
	}
	total := time.Duration(0)
	for _, part := range []time.Duration{b.MaxAge, b.Polling, b.Compute, b.Transport, b.Actuation} {
		if part < 0 || total > time.Duration(math.MaxInt64)-part {
			return 0, errors.New("budget: negative duration or overflow")
		}
		total += part
	}
	if total > b.FTTI {
		return 0, errors.New("budget: complete response exceeds FTTI")
	}
	return b.FTTI - total, nil
}

func MinimumHysteresis(amplitude, maxStep float64, horizon int) (float64, error) {
	if !finite(amplitude) || !finite(maxStep) || amplitude < 0 || maxStep < 0 || horizon < 1 {
		return 0, errors.New("hysteresis: finite nonnegative bounds and positive horizon required")
	}
	width := 2*amplitude + float64(horizon)*maxStep
	if !finite(width) {
		return 0, errors.New("hysteresis: bound overflow")
	}
	return width, nil
}

type HysteresisConfig struct {
	Limit, Correction, Width float64
	MinimumHold              time.Duration
}

type Hysteresis struct {
	config      HysteresisConfig
	active      bool
	initialized bool
	last, setAt time.Duration
}

func NewHysteresis(config HysteresisConfig) (*Hysteresis, error) {
	if !finite(config.Limit) || !finite(config.Correction) || !finite(config.Width) ||
		config.Correction < 0 || config.Width < 0 || config.MinimumHold < 0 ||
		!finite(config.Limit-config.Width-config.Correction) {
		return nil, errors.New("hysteresis: invalid configuration")
	}
	return &Hysteresis{config: config}, nil
}

func (h *Hysteresis) Step(now time.Duration, measured float64) (bool, error) {
	if now < 0 || h.initialized && now <= h.last || !finite(measured) || !finite(measured+h.config.Correction) {
		return false, errors.New("hysteresis: finite measurement and strictly increasing nonnegative time required")
	}
	h.initialized, h.last = true, now
	upper := measured + h.config.Correction
	if !h.active && upper > h.config.Limit {
		h.active, h.setAt = true, now
	} else if h.active && upper < h.config.Limit-h.config.Width && now-h.setAt >= h.config.MinimumHold {
		h.active = false
	}
	return h.active, nil
}

func RecoveryMargin(config HysteresisConfig, openEquilibrium float64) (float64, error) {
	if _, err := NewHysteresis(config); err != nil {
		return 0, err
	}
	if !finite(openEquilibrium) {
		return 0, errors.New("recovery: nonfinite equilibrium")
	}
	margin := config.Limit - config.Correction - config.Width - openEquilibrium
	if !finite(margin) {
		return 0, errors.New("recovery: margin overflow")
	}
	return margin, nil
}

type PlantResult struct {
	Width           float64 `json:"width_mpa"`
	Trips           int     `json:"trips"`
	Releases        int     `json:"releases"`
	ProtectionTicks int     `json:"protection_ticks"`
	Peak            float64 `json:"peak_mpa"`
	Final           float64 `json:"final_mpa"`
	RecoveryMargin  float64 `json:"recovery_margin_mpa"`
}

// SimulateMemoPlant isolates the zero-ripple case; it is not a hardware experiment.
func SimulateMemoPlant(width float64, ticks int) (PlantResult, error) {
	if ticks < 1 || ticks > 1000000 {
		return PlantResult{}, errors.New("plant: tick count outside [1,1000000]")
	}
	config := HysteresisConfig{Limit: 15, Correction: .24, Width: width}
	control, err := NewHysteresis(config)
	if err != nil {
		return PlantResult{}, err
	}
	out := PlantResult{Width: width, Peak: 14.5}
	out.RecoveryMargin, err = RecoveryMargin(config, 13.8)
	if err != nil {
		return PlantResult{}, err
	}
	pressure, previous, physical := 14.5, false, false
	for tick := 0; tick < ticks; tick++ {
		active, err := control.Step(time.Duration(tick)*time.Millisecond, pressure)
		if err != nil {
			return PlantResult{}, err
		}
		if active && !previous {
			out.Trips++
		}
		if previous && !active {
			out.Releases++
		}
		if active {
			out.ProtectionTicks++
		}
		target := 15.6
		if physical {
			target = 13.8
		}
		step := math.Max(-.12, math.Min(.12, .065*(target-pressure)))
		pressure += step
		out.Peak = math.Max(out.Peak, pressure)
		previous, physical = active, active
	}
	out.Final = pressure
	return out, nil
}
