// Copyright ©2025 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vgtk_test

import (
	"bytes"
	"image/color"
	"image/png"
	"math"

	tk "modernc.org/tk9.0"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgtk"
)

// An example of making a Tk plot.
func Example() {
	p := plot.New()
	p.Title.Text = "My title"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"
	quad := plotter.NewFunction(func(x float64) float64 {
		return x * x
	})
	quad.Color = color.RGBA{B: 255, A: 255}

	exp := plotter.NewFunction(func(x float64) float64 {
		return math.Pow(2, x)
	})
	exp.Dashes = []vg.Length{vg.Points(2), vg.Points(2)}
	exp.Width = vg.Points(2)
	exp.Color = color.RGBA{G: 255, A: 255}

	sin := plotter.NewFunction(func(x float64) float64 {
		return 10*math.Sin(x) + 50
	})
	sin.Dashes = []vg.Length{vg.Points(4), vg.Points(5)}
	sin.Width = vg.Points(4)
	sin.Color = color.RGBA{R: 255, A: 255}

	p.Add(quad, exp, sin)
	p.Legend.Add("x^2", quad)
	p.Legend.Add("2^x", exp)
	p.Legend.Add("10*sin(x)+50", sin)
	p.Legend.ThumbnailWidth = 0.5 * vg.Inch

	p.X.Min = 0
	p.X.Max = 10
	p.Y.Min = 0
	p.Y.Max = 100

	p.Add(plotter.NewGrid())

	pic, err := png.Decode(bytes.NewReader(logo))
	if err != nil {
		panic(err)
	}

	cnv := vgtk.New(20*vg.Centimeter, 15*vg.Centimeter, tk.Background(tk.Red))
	cnv.DrawImage(vg.Rectangle{Max: vg.Point{220, 125}}.Add(vg.Point{50, 275}), pic)
	p.Draw(draw.New(cnv))
	cnv.Flush()

	tk.Pack(
		tk.TLabel(tk.Image(cnv)),
		tk.TExit(),
		tk.Padx("1m"), tk.Pady("2m"), tk.Ipadx("1m"), tk.Ipady("1m"),
	)
	tk.App.WmTitle("Gonum")
	//tk.ActivateTheme("azure light")
	tk.App.SetResizable(false, false)
	tk.App.Wait()
}
