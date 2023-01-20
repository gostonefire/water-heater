package plot

import (
	"fmt"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"image/color"
	"time"
)

// Plot - Plots a temperature curve to a PNG image file.
// The image file name will be built up using setTemp, pGain, iGain, dGain and a datetime string
//   - values is the X/Y coordinates given as time (a float64 representing seconds from start) and temperature
//   - setTemp is the temperature that the system is aiming for to reach and hold
//   - pGain is the used proportional gain
//   - iGain is the used integral gain
//   - dGain is the used derivative gain
//   - outDir is the directory in which to write the plot file
//
// It returns a standard Go error or nil if everything went well
func Plot(values [][2]float64, setTemp, pGain, iGain, dGain float64, outDir string) (err error) {
	p := plot.New()
	p.Title.Text = "Water heating"
	p.X.Label.Text = "Time"
	p.Y.Label.Text = "Celsius"
	p.Add(plotter.NewGrid())

	nValues := len(values)
	pts := make(plotter.XYs, nValues)

	for i, v := range values {
		pts[i].X = v[0]
		pts[i].Y = v[1]
	}

	// Add set temp line to plot
	ls, err := plotter.NewLine(plotter.XYs{{X: 0, Y: setTemp}, {X: values[nValues-1][0], Y: setTemp}})
	if err != nil {
		err = fmt.Errorf("error while creating plot: %s", err)
		return
	}
	ls.LineStyle.Width = vg.Points(1)
	ls.LineStyle.Color = color.RGBA{B: 200, A: 200}
	ls.LineStyle.Dashes = []vg.Length{vg.Points(5), vg.Points(5)}

	lt, err := plotter.NewLine(pts)
	if err != nil {
		err = fmt.Errorf("error while creating plot: %s", err)
		return
	}
	lt.LineStyle.Width = vg.Points(1)
	lt.LineStyle.Color = color.RGBA{R: 255, A: 255}

	p.Add(ls, lt)
	p.Legend.Add("SetTemp", ls)
	p.Legend.Add("Temp", lt)

	ts := time.Now().Format("20060102_150405")
	fName := fmt.Sprintf("%s/wh_%d_%.3f_%.3f_%.3f_%s.png", outDir, int(setTemp), pGain, iGain, dGain, ts)
	if err = p.Save(6*vg.Inch, 4*vg.Inch, fName); err != nil {
		err = fmt.Errorf("error while saving plot: %s", err)
		return
	}

	return
}
