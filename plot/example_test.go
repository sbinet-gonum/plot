// Copyright 2012 The Plotinum Authors. All rights reserved.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

package plot

import (
	"code.google.com/p/plotinum/vg"
	"math/rand"
	"fmt"
)

// An example of making and saving a plot.
func Example() *Plot {
	// Get some data to plot.
	pts := make(XYLabels, 10)
	for i := range pts {
		if i == 0 {
			pts[i].X = rand.Float64()
		} else {
			pts[i].X = pts[i-1].X + rand.Float64()
		}
		pts[i].Y = rand.Float64()
		pts[i].Label = fmt.Sprintf("%05d", i)
	}

	// Make our plot and set some labels.
	p, err := New()
	if err != nil {
		panic(err)
	}
	p.Title.Text = "Plot Title"
	p.X.Label.Text = "X Values"
	p.Y.Label.Text = "Y Values"
	line := Line{pts, DefaultLineStyle}
	scatter := Scatter{pts, DefaultGlyphStyle}
	labels, err := MakeLabels(pts)
	if err != nil {
		panic(err)
	}
	labels.XOffs = scatter.Radius
	labels.YOffs = scatter.Radius
	p.Add(line, scatter, labels)
	p.Legend.Add("line", line, scatter)
	p.Legend.Top = true
	return p
}

// Draw the plotinum logo.
func Example_logo() *Plot {
	p, err := New()
	if err != nil {
		panic(err)
	}

	DefaultLineStyle.Width = vg.Points(1)
	DefaultGlyphStyle.Radius = vg.Points(3)

	p.Y.Tick.Marker = ConstantTicks([]Tick{{0, "0"}, {0.25, ""}, {0.5, "0.5"}, {0.75, ""}, {1, "1"}})
	p.X.Tick.Marker = ConstantTicks([]Tick{{0, "0"}, {0.25, ""}, {0.5, "0.5"}, {0.75, ""}, {1, "1"}})

	pts := XYs{{0,0}, {0, 1}, {0.5, 1}, {0.5, 0.6}, {0, 0.6}}
	line := Line{pts, DefaultLineStyle}
	scatter := Scatter{pts, DefaultGlyphStyle}
	p.Add(line, scatter)

	pts = XYs{{1, 0}, {0.75, 0}, {0.75, 0.75}}
	line = Line{pts, DefaultLineStyle}
	scatter = Scatter{pts, DefaultGlyphStyle}
	p.Add(line, scatter)

	pts = XYs{{0.5, 0.5}, {1, 0.5}}
	line = Line{pts, DefaultLineStyle}
	scatter = Scatter{pts, DefaultGlyphStyle}
	p.Add(line, scatter)

	return p
}

// An example of making a box plot.
func Example_boxPlot() *Plot {
	// Get some data to plot.
	n := 10
	uniform := make(Ys, n)
	normal := make(Ys, n)
	expon := make(Ys, n)
	for i := 0; i < n; i++ {
		uniform[i] = rand.Float64()
		normal[i] = rand.NormFloat64()
		expon[i] = rand.ExpFloat64()
	}

	// Make our plot and set some labels.
	p, err := New()
	if err != nil {
		panic(err)
	}
	p.Title.Text = "Plot Title"
	p.Y.Label.Text = "Values"

	// Make boxes for our data and add them to the plot.
	p.Add(NewBox(vg.Points(20), 0, uniform),
		NewBox(vg.Points(20), 1, normal),
		NewBox(vg.Points(20), 2, expon))

	// Set the X axis of the plot to nominal with
	// the given names for x=0, x=1 and x=2.
	p.NominalX("Uniform\nDistribution", "Normal\nDistribution",
		"Exponential\nDistribution")
	return p
}

// An example of making a horizontal box plot.
func Example_horizontalBoxes() *Plot {
	// Get some data to plot.
	n := 10
	uniform := make(Ys, n)
	normal := make(Ys, n)
	expon := make(Ys, n)
	for i := 0; i < n; i++ {
		uniform[i] = rand.Float64()
		normal[i] = rand.NormFloat64()
		expon[i] = rand.ExpFloat64()
	}

	// Make our plot and set some labels.
	p, err := New()
	if err != nil {
		panic(err)
	}
	p.Title.Text = "Plot Title"
	p.X.Label.Text = "Values"

	// Make horizontal boxes for our data and add
	// them to the plot.
	p.Add(MakeHorizBox(vg.Points(20), 0, uniform),
		MakeHorizBox(vg.Points(20), 1, normal),
		MakeHorizBox(vg.Points(20), 2, expon))

	// Set the Y axis of the plot to nominal with
	// the given names for y=0, y=1 and y=2.
	p.NominalY("Uniform\nDistribution", "Normal\nDistribution",
		"Exponential\nDistribution")
	return p
}