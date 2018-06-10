// Copyright ©2018 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vgpdf_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"testing"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/internal/cmpimg"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgpdf"
)

// ExampleEmbedFonts shows how one can embed (or not) fonts inside
// a PDF plot.
func ExampleEmbedFonts() {
	p, err := plot.New()
	if err != nil {
		log.Fatalf("could not create plot: %v", err)
	}

	pts := plotter.XYs{{0, 0}, {0, 1}, {1, 0}, {1, 1}}
	line, err := plotter.NewLine(pts)
	if err != nil {
		log.Fatalf("could not create line: %v", err)
	}
	p.Add(line)
	p.X.Label.Text = "X axis"
	p.Y.Label.Text = "Y axis"

	c := vgpdf.New(100, 100)

	// enable/disable embedding fonts
	c.EmbedFonts(true)
	p.Draw(draw.New(c))

	f, err := os.Create("testdata/enable-embedded-fonts.pdf")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	_, err = c.WriteTo(f)
	if err != nil {
		log.Fatalf("could not write canvas: %v", err)
	}

	err = f.Close()
	if err != nil {
		log.Fatalf("could not save canvas: %v", err)
	}
}

func TestEmbedFonts(t *testing.T) {
	for _, tc := range []struct {
		name  string
		embed bool
	}{
		{
			name:  "testdata/disable-embedded-fonts_golden.pdf",
			embed: false,
		},
		{
			name:  "testdata/enable-embedded-fonts_golden.pdf",
			embed: true,
		},
	} {
		t.Run(fmt.Sprintf("embed=%v", tc.embed), func(t *testing.T) {
			p, err := plot.New()
			if err != nil {
				t.Fatalf("could not create plot: %v", err)
			}

			pts := plotter.XYs{{0, 0}, {0, 1}, {1, 0}, {1, 1}}
			line, err := plotter.NewLine(pts)
			if err != nil {
				t.Fatalf("could not create line: %v", err)
			}
			p.Add(line)
			p.X.Label.Text = "X axis"
			p.Y.Label.Text = "Y axis"

			c := vgpdf.New(100, 100)

			// enable/disable embedding fonts
			c.EmbedFonts(tc.embed)
			p.Draw(draw.New(c))

			var buf bytes.Buffer
			_, err = c.WriteTo(&buf)
			if err != nil {
				t.Fatalf("could not write canvas: %v", err)
			}

			want, err := ioutil.ReadFile(tc.name)
			if err != nil {
				t.Fatalf("failed to read golden plot: %v", err)
			}

			ok, err := cmpimg.Equal("pdf", buf.Bytes(), want)
			if err != nil {
				t.Fatalf("failed to run cmpimg test: %v", err)
			}

			if !ok {
				t.Fatalf("plot mismatch: %v", tc.name)
			}
		})
	}
}

type MyCircleGlyph struct {
	Start float64
	Angle float64
}

func (gly MyCircleGlyph) DrawGlyph(c *draw.Canvas, sty draw.GlyphStyle, pt vg.Point) {
	var p vg.Path
	p.Move(vg.Point{X: pt.X + sty.Radius, Y: pt.Y})
	p.Arc(pt, sty.Radius, gly.Start, gly.Angle)
	p.Move(vg.Point{X: pt.X + sty.Radius, Y: pt.Y})
	p.Close()
	c.Fill(p)
}

type MyRingGlyph struct {
	Start float64
	Angle float64
}

func (gly MyRingGlyph) DrawGlyph(c *draw.Canvas, sty draw.GlyphStyle, pt vg.Point) {
	c.SetLineStyle(draw.LineStyle{Color: sty.Color, Width: vg.Points(0.5)})
	var p vg.Path
	p.Move(vg.Point{X: pt.X + sty.Radius, Y: pt.Y})
	p.Arc(pt, sty.Radius, gly.Start, gly.Angle)
	p.Move(vg.Point{X: pt.X + sty.Radius, Y: pt.Y})
	p.Close()
	c.Stroke(p)
}

func TestArc(t *testing.T) {
	p, err := plot.New()
	if err != nil {
		t.Fatal(err)
	}
	for ia, angle := range []float64{math.Pi / 2, math.Pi, 3 * math.Pi / 2} {
		for i, tc := range []struct {
			x, y float64
		}{
			{1, 1},
			{1, 2},
			{1, 3},
			{1, 4},
			{1, 5},
			{1, 6},
			{1, 7},
			{1, 8},
			{1, 9},
			{1, 10},
			{1, 11},
			{1, 12},
		} {
			pts1 := plotter.XYs{{tc.x + float64(ia), tc.y}}
			sca1, err := plotter.NewScatter(pts1)
			if err != nil {
				t.Fatalf("could not create scatter-1 for %v: %v", tc, err)
			}
			sca1.Shape = MyCircleGlyph{float64(i) * math.Pi / 6.0, angle}
			p.Add(sca1)

			pts2 := plotter.XYs{{tc.x + float64(ia) + 0.5, tc.y}}
			sca2, err := plotter.NewScatter(pts2)
			if err != nil {
				t.Fatalf("could not create scatter-2 for %v: %v", tc, err)
			}
			sca2.Shape = MyRingGlyph{float64(i) * math.Pi / 6.0, angle}
			p.Add(sca2)
		}
	}

	pts := plotter.XYs{{1, 0}, {2, 0}}
	scat, err := plotter.NewScatter(pts)
	if err != nil {
		t.Fatal(err)
	}
	p.Add(scat)

	c := vgpdf.New(100, 100)

	c.EmbedFonts(false)
	p.Draw(draw.New(c))
	p.Save(500, 500, "arc.png")
	p.Save(500, 500, "arc.pdf")

	f, err := os.Create("testdata/arc.pdf")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	_, err = c.WriteTo(f)
	if err != nil {
		t.Fatalf("could not write canvas: %v", err)
	}

	err = f.Close()
	if err != nil {
		t.Fatal(err)
	}

	want, err := ioutil.ReadFile("testdata/arc_golden.pdf")
	if err != nil {
		t.Fatalf("failed to read golden plot: %v", err)
	}

	got, err := ioutil.ReadFile("testdata/arc.pdf")
	if err != nil {
		t.Fatalf("failed to read plot: %v", err)
	}

	ok, err := cmpimg.Equal("pdf", got, want)
	if err != nil {
		t.Fatalf("failed to run cmpimg test: %v", err)
	}

	if !ok {
		t.Fatalf("plot mismatch")
	}
}
