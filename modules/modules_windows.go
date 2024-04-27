/*
 *          Copyright 2023, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

// Package modules provides graphic modules for the g2d engine.
package modules

// #cgo CFLAGS: -DG2D_MODULES_WIN32 -DUNICODE
// #cgo LDFLAGS: -luser32 -lgdi32 -lOpenGL32
// #include "modules.h"
import "C"
import (
	"unsafe"
)

func (rects *Rectangles) FuncCInit() unsafe.Pointer {
	return C.g2d_mods_rects_init
}

func (rects *Rectangles) SetCData(dataC unsafe.Pointer) {
	rects.dataC = dataC
}
