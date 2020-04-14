// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:generate stringer -type Kind

package token // import "gonum.org/v1/plot/internal/latex/token"

type Kind int

const (
	Invalid Kind = iota
	Macro
	EmptyLine
	Comment
	Space
	Word
	Number
	Dollar
	Lbrace
	Rbrace
	Lbrack
	Rbrack
	Equal
	Underscore
	Lparen
	Rparen
	Lt
	Gt
	Hat
	Div
	Mul
	Sub
	Add
	Not
	Colon
	Backslash
	Other
	Verbatim
	EOF
)

type Token struct {
	Kind Kind
	Pos  Pos
	Text string
}

func (t Token) String() string { return t.Text }

type Pos int

type Position struct {
	File   string
	Offset int
	Line   int
	Column int
}
