/*
 *          Copyright 2024, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

// Package loader is an OpenGL function loader for win32.
package loader

// #cgo CFLAGS: -DUNICODE
// #cgo LDFLAGS: -luser32 -lgdi32
// #include "loader.h"
import "C"
import (
	"errors"
	"github.com/vbsw/golib/cdata"
	"strconv"
	"unsafe"
)

var (
	data tCData
	conv tErrorConv
)

// tCData is the initialization for C.
type tCData struct {
}

// tErrorConv is the error converter.
type tErrorConv struct {
}

// CInitFunc returns a function to initialize C data.
func (data *tCData) CInitFunc() unsafe.Pointer {
	return C.g2d_loader_init
}

// SetCData sets initialized C data. (unused)
func (data *tCData) SetCData(unsafe.Pointer) {
}

// CData returns an instance of OpenGL function loader.
// In cdata.Init first pass (pass = 0) initializes the loader,
// second pass (pass = 1) destroys it.
func CData() cdata.CData {
	return &data
}

// ErrorConv returns an instance of error convertor.
func ErrorConv() cdata.ErrorConv {
	return &conv
}

// ToError returns error numbers/string as error.
func (errConv *tErrorConv) ToError(err1, err2 int64, info string) error {
	if err1 >= 1000000 && err1 < 1000100 {
		var errStr string
		if err1 == 1000000 {
			errStr = "vbsw.g2d.loader GetModuleHandle failed"
		} else if err1 == 1000001 {
			errStr = "vbsw.g2d.loader RegisterClassEx failed"
		} else if err1 == 1000002 {
			errStr = "vbsw.g2d.loader CreateWindow failed"
		} else if err1 == 1000003 {
			errStr = "vbsw.g2d.loader GetDC failed"
		} else if err1 == 1000004 {
			errStr = "vbsw.g2d.loader ChoosePixelFormat failed"
		} else if err1 == 1000005 {
			errStr = "vbsw.g2d.loader SetPixelFormat failed"
		} else if err1 == 1000006 {
			errStr = "vbsw.g2d.loader wglCreateContext failed"
		} else if err1 == 1000007 {
			errStr = "vbsw.g2d.loader wglMakeCurrent failed"
		} else if err1 == 1000008 {
			errStr = "vbsw.g2d.loader cdata.get failed"
		} else {
			errStr = "vbsw.g2d.loader failed"
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
