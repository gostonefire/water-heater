package pid

import "math"

// PID - The implementation of the Proportional Integral Derivative controller
type PID struct {
	DerState      float64
	IntegratState float64
	IntegratMax   float64
	IntegratMin   float64
	IntegratGain  float64
	PropGain      float64
	DerGain       float64
	inProgress    bool
	driveLow      float64
	driveHigh     float64
}

// NewPID - Returns a new PID regulator
//  - propGain is the proportional gain
//  - integratGain is the integral gain
//  - derGain is the derivative gain

func NewPID(propGain, integratGain, derGain float64) *PID {
	pid := PID{
		DerState:      0,
		IntegratState: 0,
		IntegratMax:   100,
		IntegratMin:   0,
		IntegratGain:  integratGain,
		PropGain:      propGain,
		DerGain:       derGain,
	}

	return &pid
}

// UpdatePID - Updates the drive signal to reflect to error and the current position.
// It returns the drive signal in the range 0.0 to 1.0, hence if looking at the code it does truncate very large
// internal drive values when the error, accumulated error and/or speed in system is very large. A truncation is
// necessary since the input signals are not expected to be normalized.
func (P *PID) UpdatePID(error, position float64) (drive float64) {
	var pTerm, dTerm, iTerm float64

	// Calculate the proportional term
	pTerm = P.PropGain * error

	// Calculate the integral state
	P.IntegratState += error

	// Limit the integrator state if necessary
	if P.IntegratState > P.IntegratMax {
		P.IntegratState = P.IntegratMax
	} else if P.IntegratState < P.IntegratMin {
		P.IntegratState = P.IntegratMin
	}

	// Calculate the integral term
	iTerm = P.IntegratGain * P.IntegratState

	// Calculate the derivative term, let it be zero if this is the first regulation iteration
	if !P.inProgress {
		P.inProgress = true
		P.DerState = position
	}
	dTerm = P.DerGain * (P.DerState - position)
	P.DerState = position

	// Sum up terms to drive output and adjusts it to output range 0.0 to and including 1.0
	drive = math.Min(1.0, math.Max(0.0, (pTerm+dTerm+iTerm)/100))
	return
}
