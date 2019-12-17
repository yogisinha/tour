// Copyright 2010 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package distro

import (
	"tour/static"
)

// MainPage serves the main image distribution web page.
var MainPage = static.Serve("distro/distro.html")
