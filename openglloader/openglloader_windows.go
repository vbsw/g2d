/*
 *       Copyright 2023, 2025, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

package openglloader

// #cgo CFLAGS: -DG2D_OPENGLLOADER_WIN32 -DUNICODE
// #cgo LDFLAGS: -luser32 -lOpenGL32
// #include "openglloader.h"
import "C"
import (
	"errors"
	"strconv"
	"unsafe"
)

func (ini tInitializer) CProcFunc() unsafe.Pointer {
	return C.g2d_openglloader_init
}

func (ini tInitializer) SetCData(data unsafe.Pointer) {
}

func (ini tInitializer) ToError(err1, err2 int64, info string) error {
	var err error
	if err1 > 0 {
		var errStr string
		if err1 < 1000000 {
			/* 101 */
			if err1 > 100 && err1 < 201 {
				errStr = "memory allocation failed"
			}
		} else {
			switch err1 {
			case 1000101:
				errStr = "g2d OpenGL Loader GetModuleHandle failed"
			case 1000102:
				errStr = "g2d OpenGL Loader RegisterClassEx failed"
			case 1000103:
				errStr = "g2d OpenGL Loader CreateWindow failed"
			case 1000104:
				errStr = "g2d OpenGL Loader GetDC failed"
			case 1000105:
				errStr = "g2d OpenGL Loader ChoosePixelFormat failed"
			case 1000106:
				errStr = "g2d OpenGL Loader SetPixelFormat failed"
			case 1000107:
				errStr = "g2d OpenGL Loader wglCreateContext failed"
			case 1000108:
				errStr = "g2d OpenGL Loader wglMakeCurrent failed"
			case 1000109:
				errStr = "g2d OpenGL Loader get cdata failed"
			}
		}
		if len(errStr) > 0 {
			errStr = errStr + " (" + strconv.FormatInt(err1, 10)
			if err2 == 0 {
				errStr = errStr + ")"
			} else {
				errStr = errStr + ", " + strconv.FormatInt(err2, 10) + ")"
			}
			if len(info) > 0 {
				errStr = errStr + "; " + info
			}
			err = errors.New(errStr)
		}
	}
	return err
}
