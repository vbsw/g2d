/*
 *          Copyright 2024, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

// Package analytics reads OpenGL info.
package analytics

// #cgo LDFLAGS: -lOpenGL32
// #include "analytics.h"
import "C"
import (
	"errors"
	"github.com/vbsw/golib/cdata"
	"strconv"
	"unsafe"
)

var (
	MaxTexSize  int
	MaxTexUnits int
	data tCData
	conv tErrorConv
)

// tCData is for use with cdata.
type tCData struct {
}

// tErrorConv is the error converter for use with cdata.
type tErrorConv struct {
}

// CData returns an instance of OpenGL analytics.
// In cdata.Init first pass (pass = 0) initializes analytics.
func CData() cdata.CData {
	return &data
}

// ErrorConv returns a new instance of error convertor.
func ErrorConv() cdata.ErrorConv {
	return &conv
}

// CInitFunc returns a function to initialize C data.
func (*tCData) CInitFunc() unsafe.Pointer {
	return C.vbsw_nltx_init
}

// SetCData sets initialized C data.
func (*tCData) SetCData(data unsafe.Pointer) {
	if data != nil {
		var mts, mtu C.int
		C.vbsw_nltx_result_and_free(data, &mts, &mtu)
		MaxTexSize = int(mts)
		MaxTexUnits = int(mtu)
	} else {
		MaxTexSize = 0
		MaxTexUnits = 0
	}
}

// ToError returns error numbers/string as error.
func (errConv *tErrorConv) ToError(err1, err2 int64, info string) error {
	if err1 >= 1000100 && err1 < 1000200 {
		var errStr string
		if err1 == 1000100 {
			errStr = "vbsw.g2d.analytics requires vbsw.g2d.oglf"
		} else {
			errStr = "vbsw.g2d.analytics failed"
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
