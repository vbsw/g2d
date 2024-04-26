/*
 *          Copyright 2024, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

// Package analytics reads OpenGL info.
package analytics

// #cgo CFLAGS: -DUNICODE
// #cgo LDFLAGS: -luser32 -lOpenGL32
// #include "analytics.h"
import "C"
import (
	"errors"
	"github.com/vbsw/golib/cdata"
	"strconv"
	"unsafe"
)

const unknownError = "unknown error"

// Analytics contains OpenGL information.
type Analytics struct {
	MaxTexSize  int
	MaxTexUnits int
}

// ErrorConv converts error numbers/strings to error.
type ErrorConv struct {
}

// tCData implements CData interface (see github.com/vbsw/golib/cdata).
type tCData struct {
	analytics *Analytics
}

// NewCData returns initializer for Analytics.
func NewCData(nltx *Analytics) cdata.CData {
	nltxData := new(tCData)
	nltxData.analytics = nltx
	return nltxData
}

// CInitFunc returns a function to initialize C data.
func (nltxData *tCData) CInitFunc() unsafe.Pointer {
	return C.vbsw_nltx_init
}

// SetCData sets initialized C data.
func (nltxData *tCData) SetCData(data unsafe.Pointer) {
	if data != nil {
		var mts, mtu C.int
		C.vbsw_nltx_result_and_free(data, &mts, &mtu)
		nltxData.analytics.MaxTexSize = int(mts)
		nltxData.analytics.MaxTexUnits = int(mtu)
	} else {
		nltxData.analytics.MaxTexSize = 0
		nltxData.analytics.MaxTexUnits = 0
	}
}

// ToError returns error numbers/string as error.
func (errConv *ErrorConv) ToError(err1, err2 int64, info string) error {
	if err1 >= 1000100 && err1 < 1000200 {
		errStr := "g2d analytics failed"
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
