// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plotter

import (
	"image/color"
	"math"

	"github.com/gonum/plot"
	"github.com/gonum/plot/vg"
	"github.com/gonum/plot/vg/vgdraw"
)

const (
	defaultFont = "Times-Roman"
)

type XAxis struct {
	plot.HorizontalAxis
}

func NewXAxis(label string) *XAxis {
	a := newAxis()
	a.Label.Text = label
	return &XAxis{plot.HorizontalAxis{a}}
}

func (x *XAxis) Axis() *plot.Axis {
	return &x.HorizontalAxis.Axis
}

func (*XAxis) IsXAxis() bool {
	return true
}

// Plot implements the plot.Plotter interface.
func (x *XAxis) Plot(c vgdraw.Canvas, plt *plot.Plot) {
	x.SanitizeRange()
	x.Draw(c)
}

type YAxis struct {
	plot.VerticalAxis
}

func NewYAxis(label string) *YAxis {
	a := newAxis()
	a.Label.Text = label
	return &YAxis{plot.VerticalAxis{a}}
}

func newAxis() plot.Axis {
	labelFont, err := vg.MakeFont(defaultFont, vg.Points(12))
	if err != nil {
		panic(err)
	}

	tickFont, err := vg.MakeFont(defaultFont, vg.Points(10))
	if err != nil {
		panic(err)
	}

	a := plot.Axis{
		Min: math.Inf(1),
		Max: math.Inf(-1),
		LineStyle: vgdraw.LineStyle{
			Color: color.Black,
			Width: vg.Points(0.5),
		},
		Padding: vg.Points(5),
		Scale:   plot.LinearScale{},
	}
	a.Label.TextStyle = vgdraw.TextStyle{
		Color: color.Black,
		Font:  labelFont,
	}
	a.Tick.Label = vgdraw.TextStyle{
		Color: color.Black,
		Font:  tickFont,
	}
	a.Tick.LineStyle = vgdraw.LineStyle{
		Color: color.Black,
		Width: vg.Points(0.5),
	}
	a.Tick.Length = vg.Points(8)
	a.Tick.Marker = plot.DefaultTicks{}

	return a
}

func (y *YAxis) Axis() *plot.Axis {
	return &y.VerticalAxis.Axis
}

func (*YAxis) IsYAxis() bool {
	return true
}

// Plot implements the plot.Plotter interface.
func (y *YAxis) Plot(c vgdraw.Canvas, plt *plot.Plot) {
	y.SanitizeRange()
	y.Draw(c)
}
