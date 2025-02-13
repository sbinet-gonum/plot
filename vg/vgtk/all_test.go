// Copyright ©2025 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vgtk // import "gonum.org/v1/plot/vg/vgtk"

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	rc := m.Run()
	os.Exit(rc)
}
