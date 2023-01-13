package main

import (
	"fmt"
	"github.com/gostonefire/water-heater/internal/pid"
	"github.com/gostonefire/water-heater/internal/simulation"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"image/color"
	"math"
)

func main() {
	Ts := 95.0
	Th := 15.0
	Ta := 0.0
	nPlot := 480

	// Set up plot
	p := plot.New()
	p.Title.Text = "Water heating"
	p.X.Label.Text = "Time"
	p.Y.Label.Text = "Celsius"
	p.Add(plotter.NewGrid())

	pts := make(plotter.XYs, nPlot)

	sim := simulation.NewSim(0, Ta, Th)

	control := pid.PID{
		DerState:      0,
		IntegratState: 0,
		IntegratMax:   100,
		IntegratMin:   0,
		IntegratGain:  0.3,
		PropGain:      20,
		DerGain:       500,
	}

	for i := 0; i < nPlot; i++ {
		Vd := control.UpdatePID(Ts-Th, Th)
		Vd = math.Min(1.0, math.Max(0.0, Vd/100))

		Th = sim.SimIter(Vd, Ta)
		if Th > 100.0 {
			Th = 100.0
			sim = simulation.NewSim(0, Ta, Th)
		}

		pts[i].X = float64(i)
		pts[i].Y = Th

		// fmt.Printf("Time: %d, Temp: %0.2f\n", i, Th)
		//fmt.Println(Th)
	}

	// Add set temp line to plot
	ls, err := plotter.NewLine(plotter.XYs{{X: 0, Y: Ts}, {X: float64(nPlot), Y: Ts}})
	if err != nil {
		panic(err)
	}
	ls.LineStyle.Width = vg.Points(1)
	ls.LineStyle.Color = color.RGBA{B: 200, A: 200}
	ls.LineStyle.Dashes = []vg.Length{vg.Points(5), vg.Points(5)}

	lt, err := plotter.NewLine(pts)
	if err != nil {
		panic(err)
	}
	lt.LineStyle.Width = vg.Points(1)
	lt.LineStyle.Color = color.RGBA{R: 255, A: 255}

	p.Add(ls, lt)
	p.Legend.Add("SetTemp", ls)
	p.Legend.Add("Temp", lt)

	fName := fmt.Sprintf("wh_%d_%.2f_%d_%d.png", int(Ts), control.IntegratGain, int(control.PropGain), int(control.DerGain))
	if err = p.Save(6*vg.Inch, 4*vg.Inch, fName); err != nil {
		panic(err)
	}
}
