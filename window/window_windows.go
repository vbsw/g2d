/*
 *          Copyright 2023, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

package window

// #cgo CFLAGS: -DG2D_WINDOW_WIN32 -DUNICODE
// #cgo LDFLAGS: -luser32 -lgdi32 -lOpenGL32
// #include "window.h"
import "C"
import (
	"unsafe"
)

const unknownError = "unknown error"

func (wnd *Window) FuncCInit() unsafe.Pointer {
	return C.g2d_window_init
}

func (wnd *Window) SetCData(dataC unsafe.Pointer) {
	wnd.dataC = dataC
}

func (errConv *ErrorConvertor) ToError(err1, err2 int64, info string) error {
	var errStr string
	if err1 > 0 {
		if err1 < 1000 {
			errStr = "memory allocation failed"
		} else if err1 < 1100 {
			errStr = unknownError
		} else if err1 < 1200 {
			if err1 == 1200 {
				errStr = "g2d window GetModuleHandle failed"
			} else if err1 == 1201 {
				errStr = "g2d window RegisterClassEx failed"
			} else if err1 == 1202 {
				errStr = "g2d window CreateWindow failed"
			} else if err1 == 1203 {
				errStr = "g2d window GetDC failed"
			} else if err1 == 1204 {
				errStr = "g2d window ChoosePixelFormat failed"
			} else if err1 == 1205 {
				errStr = "g2d window SetPixelFormat failed"
			} else if err1 == 1206 {
				errStr = "g2d window wglCreateContext failed"
			} else if err1 == 1207 {
				errStr = "g2d window wglMakeCurrent failed"
			} else if err1 == 1208 {
				errStr = "g2d window get cdata failed"
			} else {
				errStr = unknownError
			}
		} else {
			errStr = unknownError
		}
	} else {
		errStr = unknownError
	}
	errStr = errStr + " (" + strconv.FormatInt(err1, 10)
	if err2 == 0 {
		errStr = errStr + ")"
	} else {
		errStr = errStr + ", " + strconv.FormatInt(err2, 10) + ")"
	}
	if len(info) > 0 {
		errStr = errStr + "; " + info
	}
	return errors.New(errStr)
}
