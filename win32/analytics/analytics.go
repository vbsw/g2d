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

// Analytics contains OpenGL information.
type Analytics struct {
	MaxTexSize  int
	MaxTexUnits int
}

// tCData is for use with cdata.
type tCData struct {
	analytics *Analytics
}

// tErrorConv is the error converter for use with cdata.
type tErrorConv struct {
}

// NewAnalytics returns a new instance of Analytics.
func NewAnalytics() *Analytics {
	return new(Analytics)
}

// NewCData returns a wrapper for Analytics.
// In cdata.Init first pass (pass = 0) initializes Analytics.
func NewCData(nltx *Analytics) cdata.CData {
	nltxData := new(tCData)
	nltxData.analytics = nltx
	return nltxData
}

// NewErrorConv returns a new instance of error convertor.
func NewErrorConv() cdata.ErrorConv {
	return new(tErrorConv)
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
