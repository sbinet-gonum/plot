// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plotter

import (
	"errors"
	"fmt"
	"image/color"
	"math"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

// Histogram implements the Plotter interface,
// drawing a histogram of the data.
type Histogram struct {
	// Bins is the set of bins for this histogram.
	Bins []HistogramBin

	// Width is the width of each bin.
	Width float64

	// FillColor is the color used to fill each
	// bar of the histogram.  If the color is nil
	// then the bars are not filled.
	FillColor color.Color

	// LineStyle is the style of the outline of each
	// bar of the histogram.
	draw.LineStyle
}

// NewHistogram returns a new histogram
// that represents the distribution of values
// using the given number of bins.
//
// Each y value is assumed to be the frequency
// count for the corresponding x.
//
// If the number of bins is non-positive than
// a reasonable default is used.
func NewHistogram(xy XYer, n int) (*Histogram, error) {
	if n <= 0 {
		return nil, errors.New("Histogram with non-positive number of bins")
	}
	bins, width := binPoints(xy, n)
	return &Histogram{
		Bins:      bins,
		Width:     width,
		FillColor: color.Gray{128},
		LineStyle: DefaultLineStyle,
	}, nil
}

// NewHist returns a new histogram, as in
// NewHistogram, except that it accepts a Valuer
// instead of an XYer.
func NewHist(vs Valuer, n int) (*Histogram, error) {
	return NewHistogram(unitYs{vs}, n)
}

type unitYs struct {
	Valuer
}

func (u unitYs) XY(i int) (float64, float64) {
	return u.Value(i), 1.0
}

// Plot implements the Plotter interface, drawing a line
// that connects each point in the Line.
func (h *Histogram) Plot(c draw.Canvas, p *plot.Plot) {
	trX, trY := p.Transforms(&c)

	for _, bin := range h.Bins {
		pts := []vg.Point{
			{trX(bin.Min), trY(0)},
			{trX(bin.Max), trY(0)},
			{trX(bin.Max), trY(bin.Weight)},
			{trX(bin.Min), trY(bin.Weight)},
		}
		if h.FillColor != nil {
			c.FillPolygon(h.FillColor, c.ClipPolygonXY(pts))
		}
		pts = append(pts, vg.Point{X: trX(bin.Min), Y: trY(0)})
		c.StrokeLines(h.LineStyle, c.ClipLinesXY(pts)...)
	}
}

// DataRange returns the minimum and maximum X and Y values
func (h *Histogram) DataRange() (xmin, xmax, ymin, ymax float64) {
	xmin = math.Inf(1)
	xmax = math.Inf(-1)
	ymax = math.Inf(-1)
	for _, bin := range h.Bins {
		if bin.Max > xmax {
			xmax = bin.Max
		}
		if bin.Min < xmin {
			xmin = bin.Min
		}
		if bin.Weight > ymax {
			ymax = bin.Weight
		}
	}
	return
}

// Normalize normalizes the histogram so that the
// total area beneath it sums to a given value.
func (h *Histogram) Normalize(sum float64) {
	mass := 0.0
	for _, b := range h.Bins {
		mass += b.Weight
	}
	for i := range h.Bins {
		h.Bins[i].Weight *= sum / (h.Width * mass)
	}
}

// Thumbnail draws a rectangle in the given style of the histogram.
func (h *Histogram) Thumbnail(c *draw.Canvas) {
	ymin := c.Min.Y
	ymax := c.Max.Y
	xmin := c.Min.X
	xmax := c.Max.X

	pts := []vg.Point{
		{xmin, ymin},
		{xmax, ymin},
		{xmax, ymax},
		{xmin, ymax},
	}
	if h.FillColor != nil {
		c.FillPolygon(h.FillColor, c.ClipPolygonXY(pts))
	}
	pts = append(pts, vg.Point{X: xmin, Y: ymin})
	c.StrokeLines(h.LineStyle, c.ClipLinesXY(pts)...)
}

// binPoints returns a slice containing the
// given number of bins, and the width of
// each bin.
//
// If the given number of bins is not positive
// then a reasonable default is used.  The
// default is the square root of the sum of
// the y values.
func binPoints(xys XYer, n int) ([]HistogramBin, float64) {
	xmin, xmax := Range(XValues{xys})
	if n <= 0 {
		m := 0.0
		for i := 0; i < xys.Len(); i++ {
			_, y := xys.XY(i)
			m += math.Max(y, 1.0)
		}
		n = int(math.Ceil(math.Sqrt(m)))
	}
	if n < 1 || xmax <= xmin {
		n = 1
	}

	bins := make([]HistogramBin, n)

	w := (xmax - xmin) / float64(n)
	for i := range bins {
		bins[i].Min = xmin + float64(i)*w
		bins[i].Max = xmin + float64(i+1)*w
	}

	for i := 0; i < xys.Len(); i++ {
		x, y := xys.XY(i)
		bin := int((x - xmin) / w)
		if x == xmax {
			bin = n - 1
		}
		if bin < 0 || bin >= n {
			panic(fmt.Sprintf("%g, xmin=%g, xmax=%g, w=%g, bin=%d, n=%d\n",
				x, xmin, xmax, w, bin, n))
		}
		bins[bin].Weight += y
	}
	return bins, w
}

// A HistogramBin approximates the number of values
// within a range by a single number (the weight).
type HistogramBin struct {
	Min, Max float64
	Weight   float64
}

// LogHistogram draws an histogram of the data.
// LogHistogram is like Histogram, but tailored for the special case
// of a log-Y axis.
type LogHistogram struct {
	*Histogram
}

// Plot implements the Plotter interface, drawing a line
// that connects each point in the Line.
func (h LogHistogram) Plot(c draw.Canvas, p *plot.Plot) {
	trX, trY := p.Transforms(&c)

	_, _, ymin, _ := h.DataRange()
	for _, bin := range h.Bins {
		pts := []vg.Point{
			{trX(bin.Min), trY(ymin)},
			{trX(bin.Max), trY(ymin)},
			{trX(bin.Max), trY(bin.Weight)},
			{trX(bin.Min), trY(bin.Weight)},
		}
		if h.FillColor != nil {
			c.FillPolygon(h.FillColor, c.ClipPolygonXY(pts))
		}
		pts = append(pts, vg.Point{X: trX(bin.Min), Y: trY(ymin)})
		c.StrokeLines(h.LineStyle, c.ClipLinesXY(pts)...)
	}
}

// DataRange returns the minimum and maximum X and Y values
func (h LogHistogram) DataRange() (xmin, xmax, ymin, ymax float64) {
	xmin = math.Inf(1)
	xmax = math.Inf(-1)
	ymin = math.Inf(+1)
	ymax = math.Inf(-1)
	for _, bin := range h.Bins {
		if bin.Max > xmax {
			xmax = bin.Max
		}
		if bin.Min < xmin {
			xmin = bin.Min
		}
		if bin.Weight < ymin {
			ymin = bin.Weight
		}
		if bin.Weight > ymax {
			ymax = bin.Weight
		}
	}
	ymin /= 2.5
	return
}
