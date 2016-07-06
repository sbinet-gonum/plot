// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plot

import (
	"fmt"
	"image/color"

	"github.com/gonum/plot/style"
	"github.com/gonum/plot/vg"
	"github.com/gonum/plot/vg/draw"
)

// DefaultTheme is the default style used by gonum/plot.
func DefaultStyle() *defaultStyle {
	return &defaultStyle{
		Style: makeDefaultStyle(DefaultFont, 12),
	}
}

type defaultStyle struct {
	Style style.Plot
}

func makeDefaultStyle(fontName string, fontSize vg.Length) style.Plot {
	var (
		titleFont     = mustMakeFont(fontName, fontSize)
		legendFont    = mustMakeFont(fontName, fontSize)
		axisLabelFont = mustMakeFont(fontName, fontSize)
		axisTickFont  = mustMakeFont(fontName, fontSize-2)
	)

	axis := style.Axis{
		Label: draw.TextStyle{
			Color: color.Black,
			Font:  axisLabelFont,
		},
		Line: draw.LineStyle{
			Color: color.Black,
			Width: vg.Points(0.5),
		},
		Padding: vg.Points(5),
		Tick: style.Tick{
			Label: draw.TextStyle{
				Color: color.Black,
				Font:  axisTickFont,
			},
			Line: draw.LineStyle{
				Color: color.Black,
				Width: vg.Points(0.5),
			},
			Length: vg.Points(8),
		},
	}

	return style.Plot{
		Font: fontName,
		Title: style.Title{
			Text: draw.TextStyle{
				Color: color.Black,
				Font:  titleFont,
			},
		},
		BackgroundColor: color.White,
		X:               axis,
		Y:               axis,
		Legend: style.Legend{
			Text: draw.TextStyle{
				Color: color.Black,
				Font:  legendFont,
			},
			Padding: 0,
		},
	}
}

func (sty *defaultStyle) DrawPlot(c draw.Canvas, p *Plot) {
	// TODO: actually draw the *Plot on the draw.Canvas,
	// following defaultTheme's theme.
}

func (sty *defaultStyle) Plot() style.Plot {
	return sty.Style
}

func mustMakeFont(name string, size vg.Length) vg.Font {
	font, err := vg.MakeFont(name, size)
	if err != nil {
		panic(fmt.Errorf("plot/theme: %v\n", err))
	}
	return font
}
