// Copyright 2010 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stack

import (
	"fmt"
	"runtime"
)

// Stack returns a textual description of the call stack.
func Stack() string {
	var s string
	for n := 2;; n++ {
		pc, file, line, ok := runtime.Caller(n)
		if !ok {
			break
		}
		var funcname string
		f := runtime.FuncForPC(pc)
		if f != nil {
			funcname = f.Name()
		}
		s += fmt.Sprintf("%-30s %s:%d\n", funcname, file, line)
	}
	return s
}
