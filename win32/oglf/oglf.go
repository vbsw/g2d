/*
 *          Copyright 2024, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

// Package oglf is an OpenGL function loader for win32.
package oglf

// #cgo CFLAGS: -DUNICODE
// #cgo LDFLAGS: -luser32 -lgdi32
// #include "oglf.h"
import "C"
import (
	"errors"
	"github.com/vbsw/golib/cdata"
	"strconv"
	"unsafe"
)

// tCData is for use with cdata.
type tCData struct {
}

// tErrorConv is the error converter for use with cdata.
type tErrorConv struct {
}

// CInitFunc returns a function to initialize C data.
func (data *tCData) CInitFunc() unsafe.Pointer {
	return C.vbsw_oglf_init
}

// SetCData sets initialized C data. (unused)
func (data *tCData) SetCData(unsafe.Pointer) {
}

// NewCData returns a new instance of OpenGL function loader.
// In cdata.Init first pass (pass = 0) initializes the loader,
// second pass (pass = 1) destroys it.
func NewCData() cdata.CData {
	return new(tCData)
}

// NewErrorConv returns a new instance of error convertor.
func NewErrorConv() cdata.ErrorConv {
	return new(tErrorConv)
}

// ToError returns error numbers/string as error.
func (errConv *tErrorConv) ToError(err1, err2 int64, info string) error {
	if err1 >= 1000000 && err1 < 1000100 {
		var errStr string
		if err1 == 1000000 {
			errStr = "vbsw.g2d.oglf GetModuleHandle failed"
		} else if err1 == 1000001 {
			errStr = "vbsw.g2d.oglf RegisterClassEx failed"
		} else if err1 == 1000002 {
			errStr = "vbsw.g2d.oglf CreateWindow failed"
		} else if err1 == 1000003 {
			errStr = "vbsw.g2d.oglf GetDC failed"
		} else if err1 == 1000004 {
			errStr = "vbsw.g2d.oglf ChoosePixelFormat failed"
		} else if err1 == 1000005 {
			errStr = "vbsw.g2d.oglf SetPixelFormat failed"
		} else if err1 == 1000006 {
			errStr = "vbsw.g2d.oglf wglCreateContext failed"
		} else if err1 == 1000007 {
			errStr = "vbsw.g2d.oglf wglMakeCurrent failed"
		} else if err1 == 1000008 {
			errStr = "vbsw.g2d.oglf get cdata failed"
		} else {
			errStr = "vbsw.g2d.oglf failed"
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
