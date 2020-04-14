// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ast // import "gonum.org/v1/plot/internal/latex/ast"

import (
	"fmt"
	"io"

	"gonum.org/v1/plot/internal/latex/token"
)

type Node interface {
	Pos() token.Pos
	End() token.Pos

	isNode()
}

type List []Node

func (x List) isNode() {}
func (x List) Pos() token.Pos {
	if len(x) == 0 {
		return -1
	}
	return x[0].Pos()
}
func (x List) End() token.Pos {
	if len(x) == 0 {
		return -1
	}
	return x[len(x)-1].End()
}

type Macro struct {
	Name Ident
	Args []Node
}

func (x *Macro) isNode()        {}
func (x *Macro) Pos() token.Pos { return x.Name.Pos() }
func (x *Macro) End() token.Pos {
	if len(x.Args) > 0 {
		return x.Args[len(x.Args)-1].End()
	}
	return x.Name.End()
}

type Arg struct {
	Lbrace token.Pos // position of '{'
	List   []Node    // or stmt?
	Rbrace token.Pos // position of '}'
}

func (x *Arg) Pos() token.Pos { return x.Lbrace }
func (x *Arg) End() token.Pos { return x.Rbrace }
func (x *Arg) isNode()        {}

type OptArg struct {
	Lbrack token.Pos // position of '['
	List   []Node
	Rbrack token.Pos // position of ']'
}

func (x *OptArg) Pos() token.Pos { return x.Lbrack }
func (x *OptArg) End() token.Pos { return x.Rbrack }
func (x *OptArg) isNode()        {}

type Ident struct {
	NamePos token.Pos // identifier position
	Name    string    // identifier name
}

func (x *Ident) Pos() token.Pos { return x.NamePos }
func (x *Ident) End() token.Pos { return token.Pos(int(x.NamePos) + len(x.Name)) }
func (x *Ident) isNode()        {}

type MathExpr struct {
	Left  token.Pos // position of opening '$'
	List  []Node
	Right token.Pos // position of closing '$'
}

func (x *MathExpr) isNode()        {}
func (x *MathExpr) Pos() token.Pos { return x.Left }
func (x *MathExpr) End() token.Pos { return x.Right }

type Var struct {
	VarPos token.Pos
	Name   string
}

func (x *Var) isNode()        {}
func (x *Var) Pos() token.Pos { return x.VarPos }
func (x *Var) End() token.Pos { return token.Pos(int(x.VarPos) + len(x.Name)) }

type Literal struct {
	LitPos token.Pos
	Text   string
}

func (x *Literal) isNode()        {}
func (x *Literal) Pos() token.Pos { return x.LitPos }
func (x *Literal) End() token.Pos { return token.Pos(int(x.LitPos) + len(x.Text)) }

type Op struct {
	OpPos token.Pos
	Text  string
}

func (x *Op) isNode()        {}
func (x *Op) Pos() token.Pos { return x.OpPos }
func (x *Op) End() token.Pos { return token.Pos(int(x.OpPos) + len(x.Text)) }

type Super struct {
	HatPos token.Pos
	Node   Node
}

func (x *Super) isNode()        {}
func (x *Super) Pos() token.Pos { return x.HatPos }
func (x *Super) End() token.Pos { return x.Node.End() }

func Print(o io.Writer, node Node) {
	switch node := node.(type) {
	case *Arg:
		fmt.Fprintf(o, "{")
		for i, n := range node.List {
			if i > 0 {
				fmt.Fprintf(o, ", ")
			}
			Print(o, n)
		}
		fmt.Fprintf(o, "}")

	case *Ident:
		fmt.Fprintf(o, "ast.Ident{%q}", node.Name)
	case *Macro:
		fmt.Fprintf(o, "ast.Macro{%q", node.Name.Name)
		switch len(node.Args) {
		case 0:
			// no-op
		default:
			fmt.Fprintf(o, ", Args:")
			for i, n := range node.Args {
				if i > 0 {
					fmt.Fprintf(o, ", ")
				}
				Print(o, n)
			}
		}
		fmt.Fprintf(o, "}")
	case *MathExpr:
		fmt.Fprintf(o, "ast.MathExpr{")
		switch len(node.List) {
		case 0:
			// no-op
		default:
			fmt.Fprintf(o, "List:")
			for i, n := range node.List {
				if i > 0 {
					fmt.Fprintf(o, ", ")
				}
				Print(o, n)
			}
		}
		fmt.Fprintf(o, "}")
	case *OptArg:
		fmt.Fprintf(o, "[")
		for i, n := range node.List {
			if i > 0 {
				fmt.Fprintf(o, ", ")
			}
			Print(o, n)
		}
		fmt.Fprintf(o, "]")
	case *Var:
		fmt.Fprintf(o, "ast.Var{%q}", node.Name)
	case *Literal:
		fmt.Fprintf(o, "ast.Lit{%q}", node.Text)
	case List:
		fmt.Fprintf(o, "ast.List{")
		for i, n := range node {
			if i > 0 {
				fmt.Fprintf(o, ", ")
			}
			Print(o, n)
		}
		fmt.Fprintf(o, "}")

	case *Super:
		fmt.Fprintf(o, "ast.Super{")
		Print(o, node.Node)
		fmt.Fprintf(o, "}")

	case *Op:
		fmt.Fprintf(o, "ast.Op{%q}", node.Text)

	default:
		panic(fmt.Errorf("unknown node %T", node))
	}
}

var (
	_ Node = (*List)(nil)
	_ Node = (*Arg)(nil)
	_ Node = (*Ident)(nil)
	_ Node = (*Macro)(nil)
	_ Node = (*MathExpr)(nil)
	_ Node = (*OptArg)(nil)
	_ Node = (*Var)(nil)
	_ Node = (*Literal)(nil)
	_ Node = (*Super)(nil)
	_ Node = (*Op)(nil)
)
