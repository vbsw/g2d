/*
 *          Copyright 2025, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

package window

// #cgo CFLAGS: -DG2D_WINDOW_WIN32 -DUNICODE
// #cgo LDFLAGS: -luser32 -lgdi32
// #include "window.h"
import "C"
import (
	"errors"
	"fmt"
	"strconv"
	"unsafe"
)

var (
	mainLoopListener MainLoopListener
	running          bool
)

const (
	loadFunctionFailed = "g2d window load %s function failed"
)

type MainLoopListener interface {
	AtMainLoopStart()
	AtMainLoopEvent()
}

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

func MainLoop(listener MainLoopListener) {
	if listener != nil {
		mainLoopListener = listener
		C.g2d_window_mainloop_process()
	} else {
		panic("MainLoopListener must not be nil")
	}
}

func ProcessCustomMessage() (int64, int64) {
	var err1, err2 C.longlong
	C.g2d_window_post_custom_msg(&err1, &err2)
	return int64(err1), int64(err1)
}

func ProcessQuitMessage() (int64, int64) {
	var err1, err2 C.longlong
	C.g2d_window_post_quit_msg(&err1, &err2)
	return int64(err1), int64(err1)
}

func ClearMessageQueue() {
	C.g2d_window_mainloop_clean_up()
}

//export g2dWindowInit
func g2dWindowInit() {
	mainLoopListener.AtMainLoopStart()
}

//export g2dWindowProcessMessages
func g2dWindowProcessMessages() {
	mainLoopListener.AtMainLoopEvent()
}
