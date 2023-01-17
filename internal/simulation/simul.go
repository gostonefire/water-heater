package simulation

import (
	"math"
)

const Tau1 float64 = 0.1
const Tau2 float64 = 0.3

// const kh float64 = 300
// const stepTime float64 = 0.00167
const kh float64 = 800
const stepTime float64 = 0.0008

type Sim struct {
	Aj []float64
	Bj []float64
	Yj []float64
	Tj []float64
	v  float64
	u  float64
	j  int
}

// NewSim - Initializes start position at time zero
func NewSim(Vd0, Ta0, Y0 float64) *Sim {
	v := math.Sqrt(math.Pow(1/Tau1+1/Tau2, 2)/4 - 1/(Tau1*Tau2))
	u := (1/Tau1 + 1/Tau2) / 2

	X0 := kh*Vd0 + Ta0
	A0 := 1 / (2 * v) * (v + u) * (Y0 - X0)
	B0 := Y0 - A0 - X0

	return &Sim{
		Aj: []float64{A0},
		Bj: []float64{B0},
		Yj: []float64{Y0},
		Tj: []float64{0},
		v:  v,
		u:  u,
		j:  0,
	}
}

// SimIter - Runs one iteration over time t, and t is defined as one second
func (S *Sim) SimIter(Vd, Ta float64) float64 {
	S.j++

	Xj := kh*Vd + Ta
	a1 := (S.v - S.u) * S.Aj[S.j-1] * math.Exp((S.v-S.u)*S.Tj[S.j-1])
	a2 := (S.v + S.u) * (S.Yj[S.j-1] - S.Bj[S.j-1]*math.Exp(-(S.v+S.u)*S.Tj[S.j-1]) - Xj)
	Aj := 1 / (2 * S.v) * (a1 + a2)
	Bj := S.Yj[S.j-1] - Aj - Xj
	Yj := Aj*math.Exp((S.v-S.u)*stepTime) + Bj*math.Exp(-(S.v+S.u)*stepTime) + Xj

	S.Aj = append(S.Aj, Aj)
	S.Bj = append(S.Bj, Bj)
	S.Yj = append(S.Yj, Yj)
	S.Tj = append(S.Tj, stepTime)

	return Yj
}
