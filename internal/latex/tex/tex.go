// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package tex provides a TeX-like box model.
//
// The following is based directly on the document 'woven' from the
// TeX82 source code.  This information is also available in printed
// form:
//
//    Knuth, Donald E.. 1986.  Computers and Typesetting, Volume B:
//    TeX: The Program.  Addison-Wesley Professional.
//
// An electronic version is also available from:
//
//    http://brokestream.com/tex.pdf
//
// The most relevant "chapters" are:
//    Data structures for boxes and their friends
//    Shipping pages out (Ship class)
//    Packaging (hpack and vpack)
//    Data structures for math mode
//    Subroutines for math mode
//    Typesetting math formulas
//
// Many of the docstrings below refer to a numbered "node" in that
// book, e.g., node123
//
// Note that (as TeX) y increases downward.
package tex

import (
	"fmt"
	"math"
)

const (
	// How much text shrinks when going to the next-smallest level.  GROW_FACTOR
	// must be the inverse of SHRINK_FACTOR.
	SHRINK_FACTOR = 0.7
	GROW_FACTOR   = 1.0 / SHRINK_FACTOR

	// The number of different sizes of chars to use, beyond which they will not
	// get any smaller
	NUM_SIZE_LEVELS = 6
)

// FontConstants is a set of magical values that control how certain things,
// such as sub- and superscripts are laid out.
// These are all metrics that can't be reliably retreived from the font metrics
// in the font itself.
type FontConstants struct {
	// Percentage of x-height of additional horiz. space after sub/superscripts
	ScriptSpace float64 // = 0.05

	// Percentage of x-height that sub/superscripts drop below the baseline
	SubDrop float64 // = 0.4

	// Percentage of x-height that superscripts are raised from the baseline
	Sup1 float64 // = 0.7

	// Percentage of x-height that subscripts drop below the baseline
	Sub1 float64 // = 0.3

	// Percentage of x-height that subscripts drop below the baseline when a
	// superscript is present
	Sub2 float64 // = 0.5

	// Percentage of x-height that sub/supercripts are offset relative to the
	// nucleus edge for non-slanted nuclei
	Delta float64 // = 0.025

	// Additional percentage of last character height above 2/3 of the
	// x-height that supercripts are offset relative to the subscript
	// for slanted nuclei
	DeltaSlanted float64 // = 0.2

	// Percentage of x-height that supercripts and subscripts are offset for
	// integrals
	DeltaIntegral float64 // = 0.1
}

var DefaultFontConstants = FontConstants{
	ScriptSpace:   0.05,
	SubDrop:       0.4,
	Sup1:          0.7,
	Sub1:          0.3,
	Sub2:          0.5,
	Delta:         0.025,
	DeltaSlanted:  0.2,
	DeltaIntegral: 0.1,
}

// Node represents a node in the TeX box model.
type Node interface {
	// Kerning returns the amount of kerning between this and the next node.
	Kerning(next Node) float64

	// Shrinks one level smaller.
	// There are only three levels of sizes, after which things
	// will no longer get smaller.
	Shrink()

	// Grows one level larger.
	// There is no limit to how big something can get.
	Grow()

	// Render renders the node at (x,y) on the canvas.
	Render(x, y float64)
}

// Box is a node with a physical location
type Box struct {
	size   int
	width  float64
	height float64
	depth  float64
}

func (*Box) Kerning(next Node) float64 { return 0 }

func (box *Box) Shrink() {
	box.size--
	if box.size >= NUM_SIZE_LEVELS {
		return
	}
	box.width *= SHRINK_FACTOR
	box.height *= SHRINK_FACTOR
	box.depth *= SHRINK_FACTOR
}

func (box *Box) Grow() {
	box.size++
	box.width *= GROW_FACTOR
	box.height *= GROW_FACTOR
	box.depth *= GROW_FACTOR
}

func (*Box) Render(x, y float64) {}

func (box *Box) hpackDims(width, height, depth *float64, stretch, shrink []float64) {
	*width += box.width
	if math.IsInf(box.height, 0) || math.IsInf(box.depth, 0) {
		return
	}
	*height = math.Max(*height, box.height)
	*depth = math.Max(*depth, box.depth)
}

func (box *Box) vpackDims(width, height, depth *float64, stretch, shrink []float64) {
	*height += *depth + box.height
	*depth = box.depth
	if math.IsInf(box.width, 0) {
		return
	}
	*width = math.Max(*width, box.width)
}

// VBox is a box with a height but no width.
type VBox struct {
	size   int
	height float64
	depth  float64
}

func (*VBox) Kerning(next Node) float64 { return 0 }

func (box *VBox) Shrink() {
	box.size--
	if box.size >= NUM_SIZE_LEVELS {
		return
	}
	box.height *= SHRINK_FACTOR
	box.depth *= SHRINK_FACTOR
}

func (box *VBox) Grow() {
	box.size++
	box.height *= GROW_FACTOR
	box.depth *= GROW_FACTOR
}

func (*VBox) Render(x, y float64) {}

func (box *VBox) hpackDims(width, height, depth *float64, stretch, shrink []float64) {
	if math.IsInf(box.height, 0) || math.IsInf(box.depth, 0) {
		return
	}
	*height = math.Max(*height, box.height)
	*depth = math.Max(*depth, box.depth)
}

func (box *VBox) vpackDims(width, height, depth *float64, stretch, shrink []float64) {
	*height += *depth + box.height
	*depth = box.depth
	*width = math.Max(*width, 0)
}

// HBox is a box with a width but no height nor depth.
type HBox struct {
	size  int
	width float64
}

func (*HBox) Kerning(next Node) float64 { return 0 }

func (box *HBox) Shrink() {
	box.size--
	if box.size >= NUM_SIZE_LEVELS {
		return
	}
	box.width *= SHRINK_FACTOR
}

func (box *HBox) Grow() {
	box.size++
	box.width *= GROW_FACTOR
}

func (*HBox) Render(x, y float64) {}

func (box *HBox) hpackDims(width, height, depth *float64, stretch, shrink []float64) {
	*width += box.width
}

func (box *HBox) vpackDims(width, height, depth *float64, stretch, shrink []float64) {
	*height += *depth
	*depth = 0
	if math.IsInf(box.width, 0) {
		return
	}
	*width = math.Max(*width, box.width)
}

// Char is a single character.
//
// Unlike TeX, the font information and metrics are stored with each `Char`
// to make it easier to lookup the font metrics when needed.  Note that TeX
// boxes have a width, height, and depth, unlike Type1 and TrueType which use
// a full bounding box and an advance in the x-direction.  The metrics must
// be converted to the TeX model, and the advance (if different from width)
// must be converted into a `Kern` node when the `Char` is added to its parent
// `HList`.
type Char struct {
	c rune

	size   int
	width  float64
	height float64
	depth  float64

	font struct {
		size float64
	}
	dpi  float64
	math bool
}

func (c Char) String() string { return string(c.c) }

func (c *Char) Kerning(next Node) float64 { panic("not implemented") }

func (box *Char) Shrink() {
	box.size--
	if box.size >= NUM_SIZE_LEVELS {
		return
	}
	box.font.size *= SHRINK_FACTOR
	box.width *= SHRINK_FACTOR
	box.height *= SHRINK_FACTOR
	box.depth *= SHRINK_FACTOR
}

func (box *Char) Grow() {
	box.size++
	box.font.size *= GROW_FACTOR
	box.width *= GROW_FACTOR
	box.height *= GROW_FACTOR
	box.depth *= GROW_FACTOR
}

func (c Char) Render(x, y float64) { panic("not implemented") }

func (c Char) hpackDims(width, height, depth *float64, stretch, shrink []float64) {
	*width += c.width
	*height = math.Max(*height, c.height)
	*depth = math.Max(*depth, c.depth)
}

func (*Char) vpackDims(width, height, depth *float64, stretch, shrink []float64) {
	panic("Char node in VList")
}

// Accent is a character with an accent.
// Accents need to be dealt with separately as they are already offset
// from the baseline in TrueType fonts.
type Accent struct {
	char Char
}

func (box *Accent) String() string            { return box.char.String() }
func (box *Accent) Kerning(next Node) float64 { return box.char.Kerning(next) }

func (box *Accent) Shrink() {
	box.char.Shrink()
	box.updateMetrics()
}

func (box *Accent) Grow() {
	box.char.Grow()
	box.updateMetrics()
}

func (box *Accent) Render(x, y float64) { panic("not implemented") }

func (box *Accent) updateMetrics() { panic("not implemented") }

func (box *Accent) hpackDims(width, height, depth *float64, stretch, shrink []float64) {
	box.char.hpackDims(width, height, depth, stretch, shrink)
}

func (*Accent) vpackDims(width, height, depth *float64, stretch, shrink []float64) {
	panic("Accent node in VList")
}

// List is a list of vertical or horizontal nodes.
type List struct {
	box      Box
	shift    float64 // shift is an arbitrary offset.
	children []Node  // children nodes of this list.

	glue struct {
		set   float64 // glue setting of this list
		sign  int     // 0: normal, -1: shrinking, 1: stretching
		order int     // the order of infinity (0 - 3) for the glue.
		ratio float64
	}
}

func ListOf(elements []Node) *List {
	list := &List{children: make([]Node, len(elements))}
	copy(list.children, elements)
	return list
}

func (lst *List) setGlue(x float64, sign int, totals []float64, errMsg string) {
	panic("FIXME")
}

func (lst *List) Kerning(next Node) float64 {
	return lst.box.Kerning(next)
}

func (lst *List) Shrink() {
	for _, node := range lst.children {
		node.Shrink()
	}
	lst.box.Shrink()
	if lst.box.size < NUM_SIZE_LEVELS {
		lst.shift *= SHRINK_FACTOR
		lst.glue.set *= SHRINK_FACTOR
	}
}

func (lst *List) Grow() {
	for _, node := range lst.children {
		node.Grow()
	}
	lst.box.Grow()
	lst.shift *= GROW_FACTOR
	lst.glue.set *= GROW_FACTOR
}

func (lst *List) Render(x, y float64) {
	lst.box.Render(x, y)
}

func (lst *List) hpackDims(width, height, depth *float64, stretch, shrink []float64) {
	*width += lst.box.width
	if math.IsInf(lst.box.height, 0) || math.IsInf(lst.box.depth, 0) {
		return
	}
	*height = math.Max(*height, lst.box.height-lst.shift)
	*depth = math.Max(*depth, lst.box.depth+lst.shift)
}

func (lst *List) vpackDims(width, height, depth *float64, stretch, shrink []float64) {
	*height += *depth + lst.box.height
	*depth = lst.box.depth
	if math.IsInf(lst.box.width, 0) {
		return
	}
	*width = math.Max(*width, lst.box.width)
}

// HList is a horizontal list of boxes.
type HList struct {
	lst List
}

func HListOf(elements []Node, doKern bool) *HList {
	lst := &HList{
		lst: *ListOf(elements),
	}
	if doKern {
		lst.kern()
	}
	const (
		width      = 0
		additional = true
	)
	lst.hpack(width, additional)
	return lst
}

// kern inserts Kern nodes between Char nodes to set kerning.
//
// The Char nodes themselves determine the amount of kerning they need.
// This method just creates the correct list.
func (lst *HList) kern() {
	if len(lst.lst.children) == 0 {
		return
	}
	var (
		n        = len(lst.lst.children)
		children = make([]Node, 0, n)
	)
	for i := range lst.lst.children {
		var (
			elem = lst.lst.children[i]
			next Node
			dist float64
		)
		if i < n-1 {
			next = lst.lst.children[i+1]
			dist = elem.Kerning(next)
		}
		children = append(children, elem)
		if dist != 0 {
			children = append(children, NewKern(dist))
		}
	}
	lst.lst.children = children
}

// hpack computes the dimensions of the resulting boxes, and adjusts the glue
// if one of those dimensions is pre-specified.
//
// The computed sizes normally enclose all of the material inside the new box;
// but some items may stick out if negative glue is used, if the box is
// overfull, or if a `\vbox` includes other boxes that have been shifted left.
//
// If additional is false, hpack will produce a box whose width is exactly as
// wide as the given 'width'.
// Otherwise, hpack will produce a box with the natural width of the contents,
// plus the given 'width'.
func (lst *HList) hpack(width float64, additional bool) {
	var (
		h float64
		d float64
		x float64

		totStretch = make([]float64, 4)
		totShrink  = make([]float64, 4)
	)

	for _, node := range lst.lst.children {
		switch node := node.(type) {
		case hpacker:
			node.hpackDims(&x, &h, &d, totStretch, totShrink)
		default:
			panic(fmt.Errorf("unknown node type %T", node))
		}
	}
	lst.lst.box.height = h
	lst.lst.box.depth = d

	if additional {
		width += x
	}
	lst.lst.box.width = width
	x = width - x
	switch {
	case x == 0:
		lst.lst.glue.sign = 0
		lst.lst.glue.order = 0
		lst.lst.glue.ratio = 0
	case x > 0:
		lst.lst.setGlue(x, 1, totStretch, "overfull")
	default:
		lst.lst.setGlue(x, -1, totShrink, "underfull")
	}
}

func (lst *HList) Kerning(next Node) float64 { return lst.lst.Kerning(next) }
func (lst *HList) Shrink()                   { lst.lst.Shrink() }
func (lst *HList) Grow()                     { lst.lst.Grow() }
func (lst *HList) Render(x, y float64)       { lst.lst.Render(x, x) }

func (lst *HList) hpackDims(width, height, depth *float64, stretch, shrink []float64) {
	lst.lst.hpackDims(width, height, depth, stretch, shrink)
}

func (lst *HList) vpackDims(width, height, depth *float64, stretch, shrink []float64) {
	lst.lst.vpackDims(width, height, depth, stretch, shrink)
}

// VList is a vertical list of boxes.
type VList struct {
	lst List
}

func VListOf(elements []Node) *VList {
	lst := &VList{lst: *ListOf(elements)}
	var (
		height     float64
		additional = true
		max        = math.Inf(+1)
	)
	lst.vpack(height, additional, max)
	return lst
}

// vpack computes the dimensions of the resulting boxes, and adjusts the
// glue if one of those dimensions is pre-specified.
//
// If additional is false, vpack will produce a box whose height is exactly as
// tall as the given 'height'.
// Otherwise, vpack will produce a box with the natural height of the contents,
// plus the given 'height'.
func (lst *VList) vpack(height float64, additional bool, l float64) {
	var (
		w float64
		d float64
		x float64

		totStretch = make([]float64, 4)
		totShrink  = make([]float64, 4)
	)

	for _, node := range lst.lst.children {
		switch node := node.(type) {
		case vpacker:
			node.vpackDims(&w, &x, &d, totStretch, totShrink)
		}
	}

	lst.lst.box.width = w
	switch {
	case d > l:
		x += d - l
		lst.lst.box.depth = l
	default:
		lst.lst.box.depth = d
	}

	if additional {
		height += x
	}
	lst.lst.box.height = height
	x = height - x

	switch {
	case x == 0:
		lst.lst.glue.sign = 0
		lst.lst.glue.order = 0
		lst.lst.glue.ratio = 0
	case x > 0:
		lst.lst.setGlue(x, +1, totStretch, "overfull")
	default:
		lst.lst.setGlue(x, -1, totShrink, "underfull")
	}
}

func (lst *VList) Kerning(next Node) float64 { return lst.lst.Kerning(next) }
func (lst *VList) Shrink()                   { lst.lst.Shrink() }
func (lst *VList) Grow()                     { lst.lst.Grow() }
func (lst *VList) Render(x, y float64)       { lst.lst.Render(x, y) }

func (lst *VList) hpackDims(width, height, depth *float64, stretch, shrink []float64) {
	lst.lst.hpackDims(width, height, depth, stretch, shrink)
}

func (lst *VList) vpackDims(width, height, depth *float64, stretch, shrink []float64) {
	lst.lst.vpackDims(width, height, depth, stretch, shrink)
}

// Rule is a solid black rectangle.
//
// Like a HList, Rule has a width, a depth and a height.
// However, if any of these dimensions is ∞, the actual value will be
// determined by running the rule up to the boundary of the innermost
// enclosing box.
// This is called a "running dimension".
// The width is never running in an HList; the height and depth are never
// running in a VList.
type Rule struct {
	box Box
	out backend
}

func NewRule(w, h, d float64, state State) *Rule {
	return &Rule{
		box: Box{},
		out: state.Backend(),
	}
}

func (rule *Rule) render(x, y, w, h float64) {
	rule.out.RenderRectFilled(x, y, x+w, y+h)
}

func (rule *Rule) Kerning(next Node) float64 { return rule.box.Kerning(next) }
func (rule *Rule) Shrink()                   { rule.box.Shrink() }
func (rule *Rule) Grow()                     { rule.box.Grow() }
func (rule *Rule) Render(x, y float64)       { rule.box.Render(x, y) }

func (rule *Rule) hpackDims(width, height, depth *float64, stretch, shrink []float64) {
	rule.box.hpackDims(width, height, depth, stretch, shrink)
}

func (rule *Rule) vpackDims(width, height, depth *float64, stretch, shrink []float64) {
	rule.box.vpackDims(width, height, depth, stretch, shrink)
}

// HRule is a horizontal rule.
type HRule struct {
	rule Rule
}

func NewHRule(state State, thickness float64) *HRule {
	if thickness < 0 {
		thickness = state.UnderlineThickness()
	}
	var (
		height = 0.5 * thickness
		depth  = 0.5 * thickness
	)
	return &HRule{
		rule: *NewRule(math.Inf(+1), height, depth, state),
	}
}

func (rule *HRule) Kerning(next Node) float64 { return rule.rule.Kerning(next) }
func (rule *HRule) Shrink()                   { rule.rule.Shrink() }
func (rule *HRule) Grow()                     { rule.rule.Grow() }
func (rule *HRule) Render(x, y float64)       { rule.rule.Render(x, y) }

func (rule *HRule) hpackDims(width, height, depth *float64, stretch, shrink []float64) {
	rule.rule.hpackDims(width, height, depth, stretch, shrink)
}

func (rule *HRule) vpackDims(width, height, depth *float64, stretch, shrink []float64) {
	rule.rule.vpackDims(width, height, depth, stretch, shrink)
}

// VRule is a vertical rule.
type VRule struct {
	rule Rule
}

func NewVRule(state State) *VRule {
	thickness := state.UnderlineThickness()
	return &VRule{
		rule: *NewRule(thickness, math.Inf(+1), math.Inf(+1), state),
	}
}

func (rule *VRule) Kerning(next Node) float64 { return rule.rule.Kerning(next) }
func (rule *VRule) Shrink()                   { rule.rule.Shrink() }
func (rule *VRule) Grow()                     { rule.rule.Grow() }
func (rule *VRule) Render(x, y float64)       { rule.rule.Render(x, y) }

func (rule *VRule) hpackDims(width, height, depth *float64, stretch, shrink []float64) {
	rule.rule.hpackDims(width, height, depth, stretch, shrink)
}

func (rule *VRule) vpackDims(width, height, depth *float64, stretch, shrink []float64) {
	rule.rule.vpackDims(width, height, depth, stretch, shrink)
}

type Glue struct {
	size         int
	width        float64
	stretch      float64
	stretchOrder int
	shrink       float64
	shrinkOrder  int
}

func NewGlue(typ string) *Glue {
	switch typ {
	case "fil":
		return newGlue(0, 1, 1, 0, 0)
	case "fill":
		return newGlue(0, 1, 2, 0, 0)
	case "filll":
		return newGlue(0, 1, 3, 0, 0)
	case "neg_fil":
		return newGlue(0, 0, 0, 1, 1)
	case "neg_fill":
		return newGlue(0, 0, 0, 1, 2)
	case "neg_filll":
		return newGlue(0, 0, 0, 1, 3)
	case "empty":
		return &Glue{}
	case "ss":
		return newGlue(0, 1, 1, -1, 1)
	default:
		panic(fmt.Errorf("tex: unknown Glue spec %q", typ))
	}
}

func newGlue(w, st float64, sto int, sh float64, sho int) *Glue {
	return &Glue{
		size:         0,
		width:        w,
		stretch:      st,
		stretchOrder: sto,
		shrink:       sh,
		shrinkOrder:  sho,
	}
}

func (g *Glue) Kerning(next Node) float64 { return 0 }

func (g *Glue) Shrink() {
	g.size--
	if g.size >= NUM_SIZE_LEVELS {
		return
	}
	g.width *= SHRINK_FACTOR
}

func (g *Glue) Grow() {
	g.size++
	g.width *= GROW_FACTOR
}

func (g *Glue) Render(x, y float64) {}

func (g *Glue) hpackDims(width, height, depth *float64, stretch, shrink []float64) {
	*width += g.width
	stretch[g.stretchOrder] += g.stretch
	shrink[g.shrinkOrder] += g.shrink
}

func (g *Glue) vpackDims(width, height, depth *float64, stretch, shrink []float64) {
	*height += *depth
	*depth = 0
	*height += g.width
	stretch[g.stretchOrder] += g.stretch
	shrink[g.shrinkOrder] += g.shrink
}

// HCentered creates an HList whose contents are centered within
// its enclosing box.
func HCentered(elements []Node) *HList {
	const doKern = false
	nodes := make([]Node, 0, len(elements)+2)
	nodes = append(nodes, NewGlue("ss"))
	nodes = append(nodes, elements...)
	nodes = append(nodes, NewGlue("ss"))
	return HListOf(nodes, doKern)
}

// VCentered creates a VList whose contents are centered within
// its enclosing box.
func VCentered(elements []Node) *VList {
	nodes := make([]Node, 0, len(elements)+2)
	nodes = append(nodes, NewGlue("ss"))
	nodes = append(nodes, elements...)
	nodes = append(nodes, NewGlue("ss"))
	return VListOf(nodes)
}

// Kern is a node with a width to specify a (normally negative) amount of spacing.
//
// This spacing correction appears in horizontal lists between letters
// like A and V, when the font designer decided it looks better to move them
// closer together or further apart.
// A Kern node can also appear in a vertical list, when its width denotes
// spacing in the vertical direction.
type Kern struct {
	size  int
	width float64
}

func NewKern(width float64) *Kern {
	return &Kern{width: width}
}

func (k *Kern) String() string { return fmt.Sprintf("k%.02f", k.width) }

func (k *Kern) Kerning(next Node) float64 { return 0 }

func (k *Kern) Shrink() {
	k.size--
	if k.size >= NUM_SIZE_LEVELS {
		return
	}
	k.width *= SHRINK_FACTOR
}

func (k *Kern) Grow() {
	k.size++
	k.width *= GROW_FACTOR
}

func (k *Kern) Render(x, y float64) {}

func (k *Kern) hpackDims(width, height, depth *float64, stretch, shrink []float64) {
	*width += k.width
}

func (k *Kern) vpackDims(width, height, depth *float64, stretch, shrink []float64) {
	*height += *depth + k.width
	*depth = 0
}

type SubSuperCluster struct {
	*HList
	nucleus interface{} // FIXME
	sub     interface{} // FIXME
	super   interface{} // FIXME
}

type hpacker interface {
	hpackDims(width, height, depth *float64, stretch, shrink []float64)
}

type vpacker interface {
	vpackDims(width, height, depth *float64, stretch, shrink []float64)
}

type State struct{}

func (State) Backend() backend            { panic("not implemented") }
func (State) UnderlineThickness() float64 { panic("not implemented") }

type backend interface {
	RenderGlyph()
	RenderRectFilled(x1, y1, x2, y2 float64)
}

var (
	_ Node = (*Box)(nil)
	_ Node = (*VBox)(nil)
	_ Node = (*HBox)(nil)
	_ Node = (*Char)(nil)
	_ Node = (*Accent)(nil)
	_ Node = (*List)(nil)
	_ Node = (*HList)(nil)
	_ Node = (*VList)(nil)
	_ Node = (*Rule)(nil)
	_ Node = (*HRule)(nil)
	_ Node = (*VRule)(nil)
	_ Node = (*Glue)(nil)
	_ Node = (*Kern)(nil)
	_ Node = (*SubSuperCluster)(nil)

	_ hpacker = (*Box)(nil)
	_ hpacker = (*VBox)(nil)
	_ hpacker = (*HBox)(nil)
	_ hpacker = (*Char)(nil)
	_ hpacker = (*Accent)(nil)
	_ hpacker = (*List)(nil)
	_ hpacker = (*HList)(nil)
	_ hpacker = (*VList)(nil)
	_ hpacker = (*Rule)(nil)
	_ hpacker = (*HRule)(nil)
	_ hpacker = (*VRule)(nil)
	_ hpacker = (*Glue)(nil)
	_ hpacker = (*Kern)(nil)
	_ hpacker = (*SubSuperCluster)(nil)

	_ vpacker = (*Box)(nil)
	_ vpacker = (*VBox)(nil)
	_ vpacker = (*HBox)(nil)
	_ vpacker = (*Char)(nil)
	_ vpacker = (*Accent)(nil)
	_ vpacker = (*List)(nil)
	_ vpacker = (*HList)(nil)
	_ vpacker = (*VList)(nil)
	_ vpacker = (*Rule)(nil)
	_ vpacker = (*HRule)(nil)
	_ vpacker = (*VRule)(nil)
	_ vpacker = (*Glue)(nil)
	_ vpacker = (*Kern)(nil)
	_ vpacker = (*SubSuperCluster)(nil)
)
