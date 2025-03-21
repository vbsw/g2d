/*
 *       Copyright 2023, 2025, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

// Package dummycontext creates a window with an OpenGL context.
package dummycontext

// #cgo CFLAGS: -DUNICODE
// #cgo LDFLAGS: -luser32 -lgdi32
// #include "dummycontext.h"
import "C"
import (
	"errors"
	"fmt"
	"github.com/vbsw/golib/cdata"
	"strconv"
	"unsafe"
)

const functionFailed = "g2d dummy context %s failed"

type tInitializer struct {
}

func NewInitializer() cdata.CData {
	return new(tInitializer)
}

func (ini tInitializer) CProcFunc() unsafe.Pointer {
	return C.g2d_dummycontext_init
}

func (ini tInitializer) SetCData(data unsafe.Pointer) {
}

func (ini tInitializer) ToError(err1, err2 int64, info string) error {
	var err error
	if err1 > 0 {
		var errStr string
		/* 101 - 200, 1000101 - 1000200 */
		if err1 < 1000000 {
			if err1 > 100 && err1 < 201 {
				errStr = "memory allocation failed"
			}
		} else {
			switch err1 {
			case 1000101:
				errStr = fmt.Sprintf(functionFailed, "GetModuleHandle")
			case 1000102:
				errStr = fmt.Sprintf(functionFailed, "RegisterClassEx")
			case 1000103:
				errStr = fmt.Sprintf(functionFailed, "CreateWindow")
			case 1000104:
				errStr = fmt.Sprintf(functionFailed, "GetDC")
			case 1000105:
				errStr = fmt.Sprintf(functionFailed, "ChoosePixelFormat")
			case 1000106:
				errStr = fmt.Sprintf(functionFailed, "SetPixelFormat")
			case 1000107:
				errStr = fmt.Sprintf(functionFailed, "wglCreateContext")
			case 1000108:
				errStr = fmt.Sprintf(functionFailed, "wglMakeCurrent")
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
