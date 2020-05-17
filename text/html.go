// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package text

import (
	"bytes"
	"fmt"
	"image/color"
	"strings"
	"sync"

	"golang.org/x/net/html"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

// HTMLHandler parses, formats and renders HTML text.
type HTMLHandler struct {
	mu    sync.Mutex
	state struct {
		cnv        *draw.Canvas
		at         vg.Point
		sty        draw.TextStyle
		lineHeight vg.Length
		fntSize    vg.Length
	}

	// Font, BoldFont, and ItalicFont are the names of the fonts
	// to be used for regular, bold, italic, and bold-italic text, respectively.
	// The defaults are "Helvetica", "Helvetica-Bold", "Helvetica-Oblique",
	// and "Helvetica-BoldOblique", respectively.
	Font, BoldFont, ItalicFont, BoldItalicFont string

	// PMarginTop and PMarginBottom are the margins before and
	// after paragraphs. Defaults are 0 and 0.833 text height units, respectively.
	PMarginTop, PMarginBottom float64

	// H1 - H6 are HTML headings
	H1, H2, H3, H4, H5, H6 HTMLHeading

	// SuperscriptPosition, SubscriptPosition, and SuperSubScale
	// are the relative positions and sizes of superscripts and subscripts.
	// Defaults are +0.25, -1.25, and 0.583, respectively.
	SuperscriptPosition, SubscriptPosition, SuperSubScale float64

	// HRMarginTop and Bottom specify the spacing above and below horizontal
	// rules. Defaults are 0.0 and 1.833 text height units, respectively.
	HRMarginTop, HRMarginBottom float64

	// HRScale specifies the width of horizontal rules. The
	// default is 0.1 text height units.
	HRScale float64

	// HRColor specifies the color of horizontal rules. The
	// default is black.
	HRColor color.Color

	// WrapLines specifies if the text should be wrapped to the next line
	// if it is too long. The default is true.
	WrapLines bool
}

type HTMLHeading struct {
	Scale        float64 // Scale is a font size scaling factor for heading.
	MarginTop    float64 // MarginTop is the margin above the heading.
	MarginBottom float64 // MarginBottom is the margin below the heading.
	Bold         bool    // Bold specifies whether the heading should be bold-face.
}

// HTML returns a default HTML text handler.
func HTML() *HTMLHandler {
	h := &HTMLHandler{
		PMarginBottom:       0.833,
		SuperscriptPosition: 0.75,
		SubscriptPosition:   -0.25,
		SuperSubScale:       0.583,
		WrapLines:           true,

		H1: HTMLHeading{
			Scale:        2.0,
			MarginTop:    1,
			MarginBottom: 1,
			Bold:         true,
		},

		H2: HTMLHeading{
			Scale:        1.5,
			MarginTop:    0.833,
			MarginBottom: 0.833,
			Bold:         true,
		},

		H3: HTMLHeading{
			Scale:        1.25,
			MarginTop:    0.75,
			MarginBottom: 0.75,
			Bold:         true,
		},

		H4: HTMLHeading{
			Scale:        1,
			MarginTop:    0.5,
			MarginBottom: 0.5,
			Bold:         true,
		},

		H5: HTMLHeading{
			Scale:        1,
			MarginTop:    0.5,
			MarginBottom: 0.5,
			Bold:         true,
		},

		H6: HTMLHeading{
			Scale:        1,
			MarginTop:    0.5,
			MarginBottom: 0.5,
			Bold:         false,
		},

		Font:           "Helvetica",
		BoldFont:       "Helvetica-Bold",
		ItalicFont:     "Helvetica-Oblique",
		BoldItalicFont: "Helvetica-BoldOblique",

		HRMarginTop:    0,
		HRMarginBottom: 1.833,
		HRScale:        0.1,
		HRColor:        color.Black,
	}
	return h
}

// Box returns the bounding box of the given text where:
//  - width is the horizontal space from the origin.
//  - height is the vertical space above the baseline.
//  - depth is the vertical space below the baseline, a negative number.
func (hdlr *HTMLHandler) Box(txt string, fnt vg.Font) (width, height, depth vg.Length) {
	panic("not implemented")
}

// Draw renders the given text with the provided style and position on the canvas.
func (hdlr *HTMLHandler) Draw(c *draw.Canvas, txt string, sty draw.TextStyle, pt vg.Point) {
	fnt, err := vg.MakeFont(hdlr.Font, sty.Font.Size)
	if err != nil {
		panic(err)
	}

	hdlr.mu.Lock()
	defer hdlr.mu.Unlock()

	hdlr.state.cnv = c
	hdlr.state.at = pt
	hdlr.state.sty = sty
	hdlr.state.sty.Font = fnt
	hdlr.state.fntSize = sty.Font.Size

	doc, err := html.Parse(bytes.NewReader([]byte(txt)))
	if err != nil {
		panic(fmt.Errorf("text/html: could not parse html string: %+v", err))
	}

	_, err = hdlr.draw(doc)
	if err != nil {
		panic(err)
	}
}

func (hdlr *HTMLHandler) draw(n *html.Node) (vg.Point, error) {
	switch n.Type {
	case html.ErrorNode:
		return hdlr.state.at, fmt.Errorf("text/html: node error: %v", n)
	case html.TextNode:
		return hdlr.text(n)
	case html.DocumentNode, html.DoctypeNode, html.CommentNode:
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if at, err := hdlr.draw(c); err != nil {
				return at, err
			}
		}
	case html.ElementNode:
		if at, err := hdlr.element(n); err != nil {
			return at, err
		}
	default:
		panic(fmt.Errorf("invalid node type %v", n.Type))
	}
	return hdlr.state.at, nil
}

// element renders an HTML element.
func (hdlr *HTMLHandler) element(e *html.Node) (vg.Point, error) {
	switch e.Data {
	case "p":
		return hdlr.paragraph(e)
	case "h1":
		return hdlr.heading(e, hdlr.H1)
	case "h2":
		return hdlr.heading(e, hdlr.H2)
	case "h3":
		return hdlr.heading(e, hdlr.H3)
	case "h4":
		return hdlr.heading(e, hdlr.H4)
	case "h5":
		return hdlr.heading(e, hdlr.H5)
	case "h6":
		return hdlr.heading(e, hdlr.H6)
	case "strong", "b":
		return hdlr.newFont(e, hdlr.BoldFont)
	case "em", "i":
		return hdlr.newFont(e, hdlr.ItalicFont)
	case "hr":
		return hdlr.hr()
	case "sup":
		return hdlr.subsuperscript(e, vg.Length(hdlr.SuperscriptPosition))
	case "sub":
		return hdlr.subsuperscript(e, vg.Length(hdlr.SubscriptPosition))
	case "html", "head", "body":
		for c := e.FirstChild; c != nil; c = c.NextSibling {
			if at, err := hdlr.draw(c); err != nil {
				return at, err
			}
		}
		return hdlr.state.at, nil
	default:
		return hdlr.state.at, fmt.Errorf("htmlvg: '%s' not implemented", e.Data)
	}
}

// paragraph renders an HTML p element.
func (hdlr *HTMLHandler) paragraph(p *html.Node) (vg.Point, error) {
	hdlr.state.at = vg.Point{
		X: hdlr.state.cnv.Min.X,
		Y: hdlr.state.at.Y - hdlr.state.fntSize*vg.Length(hdlr.PMarginTop),
	}
	hdlr.state.lineHeight = hdlr.state.sty.Font.Size
	for c := p.FirstChild; c != nil; c = c.NextSibling {
		if at, err := hdlr.draw(c); err != nil {
			return at, err
		}
	}
	hdlr.state.at = vg.Point{
		X: hdlr.state.cnv.Min.X,
		Y: hdlr.state.at.Y - hdlr.state.fntSize*(1+vg.Length(hdlr.PMarginBottom)),
	}
	return hdlr.state.at, nil
}

// text renders HTML normal text.
func (hdlr *HTMLHandler) text(t *html.Node) (vg.Point, error) {
	hdlr.writeLines(t.Data, hdlr.state.sty)
	return hdlr.state.at, nil
}

// subsuperscript renders superscript or subscript text.
func (hdlr *HTMLHandler) subsuperscript(s *html.Node, position vg.Length) (vg.Point, error) {
	hdlr.state.sty.Font.Size *= vg.Length(hdlr.SuperSubScale)
	hdlr.state.at.Y += hdlr.state.sty.Font.Size * position
	for c := s.FirstChild; c != nil; c = c.NextSibling {
		if at, err := hdlr.draw(c); err != nil {
			return at, err
		}
	}
	hdlr.state.at.Y -= hdlr.state.sty.Font.Size * position
	hdlr.state.sty.Font.Size /= vg.Length(hdlr.SuperSubScale)
	return hdlr.state.at, nil
}

func (hdlr *HTMLHandler) heading(h *html.Node, heading HTMLHeading) (vg.Point, error) {
	var (
		scale        = heading.Scale
		marginTop    = heading.MarginTop
		marginBottom = heading.MarginBottom
		bold         = heading.Bold
	)
	if bold {
		f := hdlr.state.sty.Font
		if err := hdlr.state.sty.Font.SetName(hdlr.BoldFont); err != nil {
			return hdlr.state.at, err
		}
		defer func() {
			hdlr.state.sty.Font = f
		}()
	}
	hdlr.state.at.X = hdlr.state.cnv.Min.X
	hdlr.state.at.Y -= hdlr.state.fntSize * vg.Length(marginTop)
	hdlr.state.sty.Font.Size *= vg.Length(scale)
	hdlr.state.lineHeight = hdlr.state.sty.Font.Size
	for c := h.FirstChild; c != nil; c = c.NextSibling {
		if at, err := hdlr.draw(c); err != nil {
			return at, err
		}
	}
	hdlr.state.at.Y -= hdlr.state.sty.Font.Size * vg.Length(marginBottom)
	hdlr.state.sty.Font.Size /= vg.Length(scale)
	hdlr.state.at.X = hdlr.state.cnv.Min.X
	hdlr.state.at.Y -= hdlr.state.fntSize * vg.Length(marginBottom)
	return hdlr.state.at, nil
}

func (hdlr *HTMLHandler) hr() (vg.Point, error) {
	hdlr.state.at.Y -= hdlr.state.fntSize * vg.Length(hdlr.HRMarginTop)
	hdlr.state.cnv.StrokeLine2(draw.LineStyle{
		Color: hdlr.HRColor,
		Width: hdlr.state.fntSize * vg.Length(hdlr.HRScale),
	}, hdlr.state.cnv.Min.X, hdlr.state.at.Y, hdlr.state.cnv.Max.X, hdlr.state.at.Y)
	hdlr.state.at.Y -= hdlr.state.fntSize * vg.Length(hdlr.HRMarginBottom)
	return hdlr.state.at, nil
}

// newFont temporarily changes the font.
func (hdlr *HTMLHandler) newFont(n *html.Node, font string) (vg.Point, error) {
	f := hdlr.state.sty.Font
	if err := hdlr.state.sty.Font.SetName(font); err != nil {
		return hdlr.state.at, err
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if at, err := hdlr.draw(c); err != nil {
			return at, err
		}
	}
	hdlr.state.sty.Font = f
	return hdlr.state.at, nil
}

// writeLines writes the given text to the canvas, inserting line breaks
// as necessary.
func (hdlr *HTMLHandler) writeLines(text string, sty draw.TextStyle) {
	splitFunc := func(r rune) bool {
		return r == ' ' || r == '-' // Function for choosing possible line breaks.
	}

	str := strings.Replace(text, " \n ", " ", -1)
	str = strings.Replace(str, " \n", " ", -1)
	str = strings.Replace(str, "\n ", " ", -1)
	str = strings.Replace(str, "\n", " ", -1)

	var lineStart int
	var line string
	for {
		nextBreak := -1
		if len(str) > 1 {
			nextBreak = strings.IndexFunc(str[lineStart+len(line)+1:], splitFunc)
		}
		if nextBreak != -1 && str[lineStart+len(line)+1+nextBreak] == '-' {
			// Break line after dash, not before.
			nextBreak++
		}
		var lineEnd int
		if nextBreak == -1 {
			lineEnd = len(str)
		} else {
			lineEnd = lineStart + len(line) + 1 + nextBreak
		}

		if hdlr.WrapLines && sty.Font.Width(str[lineStart:lineEnd]) > hdlr.state.cnv.Max.X-hdlr.state.at.X {
			// If we go to the next break, will the line be too long? If so,
			// insert a line break.
			lineStart += len(line)
			if hdlr.state.at.X == hdlr.state.cnv.Min.X { // Remove any trailing space at the beginning of a line.
				line = strings.TrimLeft(line, " ")
			}
			hdlr.state.cnv.Canvas.FillString(sty.Font, hdlr.state.at, line)
			hdlr.newLine()
			line = ""
		} else {
			line = str[lineStart:lineEnd]
		}
		if nextBreak == -1 {
			out := str[lineStart:]
			if hdlr.state.at.X == hdlr.state.cnv.Min.X { // Remove any trailing space at the beginning of a line.
				out = strings.TrimLeft(out, " ")
			}
			hdlr.state.cnv.Canvas.FillString(sty.Font, hdlr.state.at, out)
			hdlr.state.at.X += sty.Width(out)
			break
		}
	}
}

func (hdlr *HTMLHandler) newLine() {
	hdlr.state.at.X = hdlr.state.cnv.Min.X
	hdlr.state.at.Y -= hdlr.state.lineHeight
}

var (
	_ draw.TextHandler = (*HTMLHandler)(nil)
)
