// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package style provides theming capabilities to plots.
package style

import (
	"image/color"

	"github.com/gonum/plot/vg"
	"github.com/gonum/plot/vg/draw"
)

type Plot struct {
	Font  string
	Title Title
	// BackgroundColor is the background color of the plot.
	// The default is White.
	BackgroundColor color.Color
	// X and Y are the styles of the horizontal and vertical axes
	// of the plot respectively.
	X, Y Axis
	// Legend is the style of the plot's legend.
	Legend Legend
}

type Title struct {
	// Padding is the amount of padding
	// between the bottom of the title and
	// the top of the plot.
	Padding vg.Length
	Text    draw.TextStyle
}

type Axis struct {
	// Label is the style of the axis label.
	Label draw.TextStyle

	// Line is the style of the axis line.
	Line draw.LineStyle

	// Padding between the axis line and the data.  Having
	// non-zero padding ensures that the data is never drawn
	// on the axis, thus making it easier to see.
	Padding vg.Length

	// Tick is the style of the axis' ticks.
	Tick Tick
}

type Tick struct {
	// Label is the TextStyle on the tick labels.
	Label draw.TextStyle

	// Line is the LineStyle of the tick lines.
	Line draw.LineStyle

	// Length is the length of a major tick mark.
	// Minor tick marks are half of the length of major
	// tick marks.
	Length vg.Length
}

type Legend struct {
	// Text is the style given to the legend
	// entry texts.
	Text draw.TextStyle

	// Padding is the amount of padding to add
	// between each entry of the legend.  If Padding
	// is zero then entries are spaced based on the
	// font size.
	Padding vg.Length
}
