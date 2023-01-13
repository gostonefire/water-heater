package pid

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
}

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

	// Sum up terms to drive output
	drive = pTerm + dTerm + iTerm
	return
}
