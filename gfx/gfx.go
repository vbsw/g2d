/*
 *          Copyright 2023, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

// Package gfx provides graphics for the g2d engine.
package gfx

import (
	"unsafe"
)

type Rectangles struct {
	dataC unsafe.Pointer
}
