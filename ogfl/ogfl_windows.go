/*
 *          Copyright 2023, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

package ogfl

// #cgo CFLAGS: -DG2D_OGFL_WIN32 -DUNICODE
// #cgo LDFLAGS: -luser32
// #include "ogfl.h"
import "C"
import (
	"errors"
	"strconv"
	"unsafe"
)

const unknownError = "unknown error"

func (loader *OGlFuncLoader) FuncCInit() unsafe.Pointer {
	return C.vbsw_ogfl_init
}

func (loader *OGlFuncLoader) SetCData(unsafe.Pointer) {
}

func (errConv *ErrorConvertor) ToError(err1, err2 int64, info string) error {
	var errStr string
	if err1 > 0 {
		if err1 < 1000 {
			errStr = "memory allocation failed"
		} else if err1 < 1010 {
			if err1 == 1000 {
				errStr = "ogfl GetModuleHandle failed"
			} else if err1 == 1001 {
				errStr = "ogfl RegisterClassEx failed"
			} else if err1 == 1002 {
				errStr = "ogfl CreateWindow failed"
			} else if err1 == 1003 {
				errStr = "ogfl GetDC failed"
			} else if err1 == 1004 {
				errStr = "ogfl ChoosePixelFormat failed"
			} else if err1 == 1005 {
				errStr = "ogfl SetPixelFormat failed"
			} else if err1 == 1006 {
				errStr = "ogfl wglCreateContext failed"
			} else if err1 == 1007 {
				errStr = "ogfl wglMakeCurrent failed"
			} else if err1 == 1008 {
				errStr = "ogfl get cdata failed"
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