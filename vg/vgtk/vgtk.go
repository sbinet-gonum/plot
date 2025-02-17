// Copyright ©2025 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the GONUM-LICENSE file.

// Package vgtk provides tools to integrate plot/vg with a Tcl/Tk backend,
// provided by modernc.org/tk9.0
//
// More informations about this backend are available here: https://modernc.org/tk9.0
package vgtk // import "gonum.org/v1/plot/vg/vgtk"

import (
	"image"

	tk "modernc.org/tk9.0"
)

// Canvas returns a Tk image from the provided canvas image.
func Canvas(img image.Image) *tk.Img {
	return tk.NewPhoto(tk.Data(img))
}
