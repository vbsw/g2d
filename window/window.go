/*
 *          Copyright 2023, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

// Package window creates a window with OpenGL 3.0 context.
package window

import (
	"unsafe"
)

type Window struct {
	dataC unsafe.Pointer
}

type ErrorConvertor struct {
}
