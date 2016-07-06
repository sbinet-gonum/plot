// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plot

import (
	"math"
	"strconv"
	"time"

	"github.com/gonum/floats"
	"github.com/gonum/plot/style"
	"github.com/gonum/plot/vg"
	"github.com/gonum/plot/vg/draw"
)

// displayPrecision is a sane level of float precision for a plot.
const displayPrecision = 4

// Ticker creates Ticks in a specified range
type Ticker interface {
	// Ticks returns Ticks in a specified range
	Ticks(min, max float64) []Tick
}

// Normalizer rescales values from the data coordinate system to the
// normalized coordinate system.
type Normalizer interface {
	// Normalize transforms a value x in the data coordinate system to
	// the normalized coordinate system.
	Normalize(min, max, x float64) float64
}

// An Axis represents either a horizontal or vertical
// axis of a plot.
type Axis struct {
	// Min and Max are the minimum and maximum data
	// values represented by the axis.
	Min, Max float64

	Label struct {
		// Text is the axis label string.
		Text string
	}

	Tick struct {
		// Marker returns the tick marks.  Any tick marks
		// returned by the Marker function that are not in
		// range of the axis are not drawn.
		Marker Ticker
	}

	Style style.Axis

	// Scale transforms a value given in the data coordinate system
	// to the normalized coordinate system of the axis—its distance
	// along the axis as a fraction of the axis range.
	Scale Normalizer
}

// makeAxis returns a default Axis.
//
// The default range is (∞, ­∞), and thus any finite
// value is less than Min and greater than Max.
func makeAxis(sty style.Axis) (Axis, error) {
	a := Axis{
		Min:   math.Inf(1),
		Max:   math.Inf(-1),
		Scale: LinearScale{},
		Style: sty,
	}
	a.Tick.Marker = DefaultTicks{}
	return a, nil
}

// sanitizeRange ensures that the range of the
// axis makes sense.
func (a *Axis) sanitizeRange() {
	if math.IsInf(a.Min, 0) {
		a.Min = 0
	}
	if math.IsInf(a.Max, 0) {
		a.Max = 0
	}
	if a.Min > a.Max {
		a.Min, a.Max = a.Max, a.Min
	}
	if a.Min == a.Max {
		a.Min--
		a.Max++
	}
}

// LinearScale an be used as the value of an Axis.Scale function to
// set the axis to a standard linear scale.
type LinearScale struct{}

var _ Normalizer = LinearScale{}

// Normalize returns the fractional distance of x between min and max.
func (LinearScale) Normalize(min, max, x float64) float64 {
	return (x - min) / (max - min)
}

// LogScale can be used as the value of an Axis.Scale function to
// set the axis to a log scale.
type LogScale struct{}

var _ Normalizer = LogScale{}

// Normalize returns the fractional logarithmic distance of
// x between min and max.
func (LogScale) Normalize(min, max, x float64) float64 {
	logMin := log(min)
	return (log(x) - logMin) / (log(max) - logMin)
}

// Norm returns the value of x, given in the data coordinate
// system, normalized to its distance as a fraction of the
// range of this axis.  For example, if x is a.Min then the return
// value is 0, and if x is a.Max then the return value is 1.
func (a *Axis) Norm(x float64) float64 {
	return a.Scale.Normalize(a.Min, a.Max, x)
}

// drawTicks returns true if the tick marks should be drawn.
func (a *Axis) drawTicks() bool {
	return a.Style.Tick.Line.Width > 0 && a.Style.Tick.Length > 0
}

// A horizontalAxis draws horizontally across the bottom
// of a plot.
type horizontalAxis struct {
	Axis
}

// size returns the height of the axis.
func (a *horizontalAxis) size() (h vg.Length) {
	if a.Label.Text != "" {
		h -= a.Style.Label.Font.Extents().Descent
		h += a.Style.Label.Height(a.Label.Text)
	}
	if marks := a.Tick.Marker.Ticks(a.Min, a.Max); len(marks) > 0 {
		if a.drawTicks() {
			h += a.Style.Tick.Length
		}
		h += tickLabelHeight(a.Style.Tick.Label, marks)
	}
	h += a.Style.Line.Width / 2
	h += a.Style.Padding
	return
}

// draw draws the axis along the lower edge of a draw.Canvas.
func (a *horizontalAxis) draw(c draw.Canvas) {
	y := c.Min.Y
	if a.Label.Text != "" {
		y -= a.Style.Label.Font.Extents().Descent
		c.FillText(a.Style.Label, vg.Point{c.Center().X, y}, -0.5, 0, a.Label.Text)
		y += a.Style.Label.Height(a.Label.Text)
	}

	marks := a.Tick.Marker.Ticks(a.Min, a.Max)
	for _, t := range marks {
		x := c.X(a.Norm(t.Value))
		if !c.ContainsX(x) || t.IsMinor() {
			continue
		}
		c.FillText(a.Style.Tick.Label, vg.Point{x, y}, -0.5, 0, t.Label)
	}

	if len(marks) > 0 {
		y += tickLabelHeight(a.Style.Tick.Label, marks)
	} else {
		y += a.Style.Line.Width / 2
	}

	if len(marks) > 0 && a.drawTicks() {
		len := a.Style.Tick.Length
		for _, t := range marks {
			x := c.X(a.Norm(t.Value))
			if !c.ContainsX(x) {
				continue
			}
			start := t.lengthOffset(len)
			c.StrokeLine2(a.Style.Tick.Line, x, y+start, x, y+len)
		}
		y += len
	}

	c.StrokeLine2(a.Style.Line, c.Min.X, y, c.Max.X, y)
}

// GlyphBoxes returns the GlyphBoxes for the tick labels.
func (a *horizontalAxis) GlyphBoxes(*Plot) (boxes []GlyphBox) {
	for _, t := range a.Tick.Marker.Ticks(a.Min, a.Max) {
		if t.IsMinor() {
			continue
		}
		w := a.Style.Tick.Label.Width(t.Label)
		box := GlyphBox{
			X: a.Norm(t.Value),
			Rectangle: vg.Rectangle{
				Min: vg.Point{X: -w / 2},
				Max: vg.Point{X: w / 2},
			},
		}
		boxes = append(boxes, box)
	}
	return
}

// A verticalAxis is drawn vertically up the left side of a plot.
type verticalAxis struct {
	Axis
}

// size returns the width of the axis.
func (a *verticalAxis) size() (w vg.Length) {
	if a.Label.Text != "" {
		w -= a.Style.Label.Font.Extents().Descent
		w += a.Style.Label.Height(a.Label.Text)
	}
	if marks := a.Tick.Marker.Ticks(a.Min, a.Max); len(marks) > 0 {
		if lwidth := tickLabelWidth(a.Style.Tick.Label, marks); lwidth > 0 {
			w += lwidth
			w += a.Style.Label.Width(" ")
		}
		if a.drawTicks() {
			w += a.Style.Tick.Length
		}
	}
	w += a.Style.Line.Width / 2
	w += a.Style.Padding
	return
}

// draw draws the axis along the left side of a draw.Canvas.
func (a *verticalAxis) draw(c draw.Canvas) {
	x := c.Min.X
	if a.Label.Text != "" {
		x += a.Style.Label.Height(a.Label.Text)
		c.Push()
		c.Rotate(math.Pi / 2)
		c.FillText(a.Style.Label, vg.Point{c.Center().Y, -x}, -0.5, 0, a.Label.Text)
		c.Pop()
		x += -a.Style.Label.Font.Extents().Descent
	}
	marks := a.Tick.Marker.Ticks(a.Min, a.Max)
	if w := tickLabelWidth(a.Style.Tick.Label, marks); len(marks) > 0 && w > 0 {
		x += w
	}
	major := false
	for _, t := range marks {
		y := c.Y(a.Norm(t.Value))
		if !c.ContainsY(y) || t.IsMinor() {
			continue
		}
		c.FillText(a.Style.Tick.Label, vg.Point{x, y}, -1, -0.5, t.Label)
		major = true
	}
	if major {
		x += a.Style.Tick.Label.Width(" ")
	}
	if a.drawTicks() && len(marks) > 0 {
		len := a.Style.Tick.Length
		for _, t := range marks {
			y := c.Y(a.Norm(t.Value))
			if !c.ContainsY(y) {
				continue
			}
			start := t.lengthOffset(len)
			c.StrokeLine2(a.Style.Tick.Line, x+start, y, x+len, y)
		}
		x += len
	}
	c.StrokeLine2(a.Style.Line, x, c.Min.Y, x, c.Max.Y)
}

// GlyphBoxes returns the GlyphBoxes for the tick labels
func (a *verticalAxis) GlyphBoxes(*Plot) (boxes []GlyphBox) {
	for _, t := range a.Tick.Marker.Ticks(a.Min, a.Max) {
		if t.IsMinor() {
			continue
		}
		h := a.Style.Tick.Label.Height(t.Label)
		box := GlyphBox{
			Y: a.Norm(t.Value),
			Rectangle: vg.Rectangle{
				Min: vg.Point{Y: -h / 2},
				Max: vg.Point{Y: h / 2},
			},
		}
		boxes = append(boxes, box)
	}
	return
}

// DefaultTicks is suitable for the Tick.Marker field of an Axis,
// it returns a resonable default set of tick marks.
type DefaultTicks struct{}

var _ Ticker = DefaultTicks{}

// Ticks returns Ticks in a specified range
func (DefaultTicks) Ticks(min, max float64) (ticks []Tick) {
	const SuggestedTicks = 3
	if max < min {
		panic("illegal range")
	}
	tens := math.Pow10(int(math.Floor(math.Log10(max - min))))
	n := (max - min) / tens
	for n < SuggestedTicks {
		tens /= 10
		n = (max - min) / tens
	}

	majorMult := int(n / SuggestedTicks)
	switch majorMult {
	case 7:
		majorMult = 6
	case 9:
		majorMult = 8
	}
	majorDelta := float64(majorMult) * tens
	val := math.Floor(min/majorDelta) * majorDelta
	prec := maxInt(precisionOf(min), precisionOf(max))
	for val <= max {
		if val >= min && val <= max {
			ticks = append(ticks, Tick{Value: val, Label: formatFloatTick(val, prec)})
		}
		if math.Nextafter(val, val+majorDelta) == val {
			break
		}
		val += majorDelta
	}

	minorDelta := majorDelta / 2
	switch majorMult {
	case 3, 6:
		minorDelta = majorDelta / 3
	case 5:
		minorDelta = majorDelta / 5
	}

	val = math.Floor(min/minorDelta) * minorDelta
	for val <= max {
		found := false
		for _, t := range ticks {
			if t.Value == val {
				found = true
			}
		}
		if val >= min && val <= max && !found {
			ticks = append(ticks, Tick{Value: val})
		}
		if math.Nextafter(val, val+minorDelta) == val {
			break
		}
		val += minorDelta
	}
	return
}

// LogTicks is suitable for the Tick.Marker field of an Axis,
// it returns tick marks suitable for a log-scale axis.
type LogTicks struct{}

var _ Ticker = LogTicks{}

// Ticks returns Ticks in a specified range
func (LogTicks) Ticks(min, max float64) []Tick {
	var ticks []Tick
	val := math.Pow10(int(math.Floor(math.Log10(min))))
	if min <= 0 {
		panic("Values must be greater than 0 for a log scale.")
	}
	prec := precisionOf(max)
	for val < max*10 {
		for i := 1; i < 10; i++ {
			tick := Tick{Value: val * float64(i)}
			if i == 1 {
				tick.Label = formatFloatTick(val*float64(i), prec)
			}
			ticks = append(ticks, tick)
		}
		val *= 10
	}
	tick := Tick{Value: val, Label: formatFloatTick(val, prec)}
	ticks = append(ticks, tick)
	return ticks
}

// ConstantTicks is suitable for the Tick.Marker field of an Axis.
// This function returns the given set of ticks.
type ConstantTicks []Tick

var _ Ticker = ConstantTicks{}

// Ticks returns Ticks in a specified range
func (ts ConstantTicks) Ticks(float64, float64) []Tick {
	return ts
}

// UnixTimeTicks is suitable for axes representing time values.
// UnixTimeTicks expects values in Unix time seconds.
type UnixTimeTicks struct {
	// Ticker is used to generate a set of ticks.
	// If nil, DefaultTicks will be used.
	Ticker Ticker

	// Format is the textual representation of the time value.
	// If empty, time.RFC3339 will be used
	Format string
}

var _ Ticker = UnixTimeTicks{}

// Ticks implements plot.Ticker.
func (utt UnixTimeTicks) Ticks(min, max float64) []Tick {
	if utt.Ticker == nil {
		utt.Ticker = DefaultTicks{}
	}
	if utt.Format == "" {
		utt.Format = time.RFC3339
	}

	ticks := utt.Ticker.Ticks(min, max)
	for i := range ticks {
		tick := &ticks[i]
		if tick.Label == "" {
			continue
		}
		t := time.Unix(int64(tick.Value), 0)
		tick.Label = t.Format(utt.Format)
	}
	return ticks
}

// A Tick is a single tick mark on an axis.
type Tick struct {
	// Value is the data value marked by this Tick.
	Value float64

	// Label is the text to display at the tick mark.
	// If Label is an empty string then this is a minor
	// tick mark.
	Label string
}

// IsMinor returns true if this is a minor tick mark.
func (t Tick) IsMinor() bool {
	return t.Label == ""
}

// lengthOffset returns an offset that should be added to the
// tick mark's line to accout for its length.  I.e., the start of
// the line for a minor tick mark must be shifted by half of
// the length.
func (t Tick) lengthOffset(len vg.Length) vg.Length {
	if t.IsMinor() {
		return len / 2
	}
	return 0
}

// tickLabelHeight returns height of the tick mark labels.
func tickLabelHeight(sty draw.TextStyle, ticks []Tick) vg.Length {
	maxHeight := vg.Length(0)
	for _, t := range ticks {
		if t.IsMinor() {
			continue
		}
		h := sty.Height(t.Label)
		if h > maxHeight {
			maxHeight = h
		}
	}
	return maxHeight
}

// tickLabelWidth returns the width of the widest tick mark label.
func tickLabelWidth(sty draw.TextStyle, ticks []Tick) vg.Length {
	maxWidth := vg.Length(0)
	for _, t := range ticks {
		if t.IsMinor() {
			continue
		}
		w := sty.Width(t.Label)
		if w > maxWidth {
			maxWidth = w
		}
	}
	return maxWidth
}

func log(x float64) float64 {
	if x <= 0 {
		panic("Values must be greater than 0 for a log scale.")
	}
	return math.Log(x)
}

// formatFloatTick returns a g-formated string representation of v
// to the specified precision.
func formatFloatTick(v float64, prec int) string {
	return strconv.FormatFloat(floats.Round(v, prec), 'g', displayPrecision, 64)
}

// precisionOf returns the precision needed to display x without e notation.
func precisionOf(x float64) int {
	return int(math.Max(math.Ceil(-math.Log10(math.Abs(x))), displayPrecision))
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
