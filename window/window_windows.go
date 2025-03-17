/*
 *          Copyright 2025, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

package window

// #cgo CFLAGS: -DG2D_WINDOW_WIN32 -DUNICODE
// #cgo LDFLAGS: -luser32
// #include "window.h"
import "C"
import (
	"errors"
	"fmt"
	"strconv"
	"unsafe"
)

const (
	loadFunctionFailed = "g2d window load %s function failed"
)

func (init tInitializer) CProcFunc() unsafe.Pointer {
	return C.g2d_window_init
}

func (init tInitializer) SetCData(data unsafe.Pointer) {
}

func (init tInitializer) ToError(err1, err2 int64, info string) error {
	var err error
	if err1 > 0 {
		var errStr string
		if err1 < 1000000 {
			/* 201 */
			if err1 > 200 && err1 < 301 {
				errStr = "memory allocation failed"
			}
		} else {
			switch err1 {
			case 1000201:
				errStr = "g2d window GetModuleHandle failed"
			case 1000202:
				errStr = fmt.Sprintf(loadFunctionFailed, "wglChoosePixelFormatARB")
			case 1000203:
				errStr = fmt.Sprintf(loadFunctionFailed, "wglCreateContextAttribsARB")
			case 1000204:
				errStr = fmt.Sprintf(loadFunctionFailed, "wglSwapIntervalEXT")
			case 1000205:
				errStr = fmt.Sprintf(loadFunctionFailed, "wglGetSwapIntervalEXT")
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
