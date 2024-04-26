/*
 *          Copyright 2024, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

// Package oglf loads OpenGL functions.
package oglf

// #cgo CFLAGS: -DUNICODE
// #cgo LDFLAGS: -luser32
// #include "oglf.h"
import "C"
import (
	"errors"
	"strconv"
	"unsafe"
)

// CData is for use with github.com/vbsw/golib/cdata.
type CData struct {
}

// ErrorConv converts error numbers/strings to error.
type ErrorConv struct {
}

// CInitFunc returns a function to initialize C data.
func (data *CData) CInitFunc() unsafe.Pointer {
	return C.vbsw_oglf_init
}

// SetCData sets initialized C data. (unused)
func (data *CData) SetCData(unsafe.Pointer) {
}

// ToError returns error numbers/string as error.
func (errConv *ErrorConv) ToError(err1, err2 int64, info string) error {
	if err1 >= 1000000 && err1 < 1000100 {
		var errStr string
		if err1 == 1000000 {
			errStr = "g2d oglf GetModuleHandle failed"
		} else if err1 == 1000001 {
			errStr = "g2d oglf RegisterClassEx failed"
		} else if err1 == 1000002 {
			errStr = "g2d oglf CreateWindow failed"
		} else if err1 == 1000003 {
			errStr = "g2d oglf GetDC failed"
		} else if err1 == 1000004 {
			errStr = "g2d oglf ChoosePixelFormat failed"
		} else if err1 == 1000005 {
			errStr = "g2d oglf SetPixelFormat failed"
		} else if err1 == 1000006 {
			errStr = "g2d oglf wglCreateContext failed"
		} else if err1 == 1000007 {
			errStr = "g2d oglf wglMakeCurrent failed"
		} else if err1 == 1000008 {
			errStr = "g2d oglf get cdata failed"
		} else {
			errStr = "g2d oglf failed"
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
	return nil
}
