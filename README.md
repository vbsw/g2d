# g2d

[![Go Reference](https://pkg.go.dev/badge/github.com/vbsw/g2d.svg)](https://pkg.go.dev/github.com/vbsw/g2d) [![Go Report Card](https://goreportcard.com/badge/github.com/vbsw/g2d)](https://goreportcard.com/report/github.com/vbsw/g2d) [![Stability: Experimental](https://masterminds.github.io/stability/experimental.svg)](https://masterminds.github.io/stability/experimental.html)

## About
g2d is a framework to create 2D graphic applications. It is published on <https://github.com/vbsw/g2d> and <https://gitlab.com/vbsw/g2d>.

Demo is available here <https://github.com/vbsw/g2d-demo>.

## Copyright
Copyright 2023, 2025, Vitali Baumtrok (vbsw@mailbox.org).

g2d is distributed under the Boost Software License, version 1.0. (See accompanying file LICENSE or copy at http://www.boost.org/LICENSE_1_0.txt)

g2d is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the Boost Software License for more details.

## Example

	package main

	import (
		"fmt"
		"github.com/vbsw/g2d"
		"runtime"
	)

	// Arrange that main.main runs on main thread.
	func init() {
		runtime.LockOSThread()
	}

	func Main() {
		g2d.Init()
		g2d.MainLoop(new(g2d.WindowImpl))
		if g2d.Err != nil {
			fmt.Println(g2d.Err.Error())
		}
	}

## References
- https://go.dev/doc/install
- https://jmeubank.github.io/tdm-gcc/
- https://git-scm.com/book/en/v2/Getting-Started-Installing-Git
- https://dave.cheney.net/2013/10/12/how-to-use-conditional-compilation-with-the-go-build-tool
- https://github.com/golang/go/wiki/cgo
- https://pkg.go.dev/cmd/go#hdr-Compile_packages_and_dependencies
- https://pkg.go.dev/cmd/link
