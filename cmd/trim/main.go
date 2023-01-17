package main

import (
	"fmt"
	"github.com/gostonefire/water-heater/internal/simulation"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"image/color"
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

	Vd := 1.0
	for i := 0; i < nPlot; i++ {
		Th = sim.SimIter(Vd, Ta)
		if Th > 100.0 {
			Th = 100.0
			Vd = 0
			sim = simulation.NewSim(Vd, Ta, Th)
		}

		pts[i].X = float64(i)
		pts[i].Y = Th
	}

	// Add set temp line to plot
	ls, err := plotter.NewLine(plotter.XYs{{X: 0, Y: Ts}, {X: float64(nPlot), Y: Ts}})
	if err != nil {
		panic(err)
	}
	ls.LineStyle.Width = vg.Points(1)
	ls.LineStyle.Color = color.RGBA{B: 200, A: 200}
	ls.LineStyle.Dashes = []vg.Length{vg.Points(5), vg.Points(5)}

	// Add 2 minute line to plot
	lm, err := plotter.NewLine(plotter.XYs{{X: 120, Y: 0}, {X: 120, Y: 100}})
	if err != nil {
		panic(err)
	}
	lm.LineStyle.Width = vg.Points(1)
	lm.LineStyle.Color = color.RGBA{G: 200, A: 200}
	lm.LineStyle.Dashes = []vg.Length{vg.Points(5), vg.Points(5)}

	lt, err := plotter.NewLine(pts)
	if err != nil {
		panic(err)
	}
	lt.LineStyle.Width = vg.Points(1)
	lt.LineStyle.Color = color.RGBA{R: 255, A: 255}

	p.Add(ls, lm, lt)
	p.Legend.Add("SetTemp", ls)
	p.Legend.Add("Temp", lt)

	fName := fmt.Sprintf("wh_trim.png")
	if err = p.Save(6*vg.Inch, 4*vg.Inch, fName); err != nil {
		panic(err)
	}

}
