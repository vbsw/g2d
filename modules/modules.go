/*
 *          Copyright 2023, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

// Package modules provides graphic modules for the g2d engine.
package modules

import (
	"unsafe"
)

type Rectangles struct {
	dataC unsafe.Pointer
}
