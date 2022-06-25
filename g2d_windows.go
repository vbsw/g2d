/*
 *          Copyright 2022, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

package g2d

// #cgo CFLAGS: -DG2D_WIN32 -DUNICODE
// #cgo LDFLAGS: -luser32 -lgdi32 -lOpenGL32
// #include "g2d.h"
import "C"
import (
	"fmt"
	"errors"
	"strconv"
	"unsafe"
//	"reflect"
)

func Init(stub interface{}) {
	if !initialized {
		errC := C.g2d_init()
		initialized = bool(errC == nil)
		Err = toError(errC)
	}
}

func Show(window AbstractWindow) {
	if initialized {
		if Err == nil {
			params := newParameters()
			Err = window.Config(params)
			if Err == nil {
				x := C.int(params.ClientX)
				y := C.int(params.ClientY)
				w := C.int(params.ClientWidth)
				h := C.int(params.ClientHeight)
				wn := C.int(params.ClientMinWidth)
				hn := C.int(params.ClientMinHeight)
				wx := C.int(params.ClientMaxWidth)
				hx := C.int(params.ClientMaxHeight)
				c := toCInt(params.Centered)
				l := toCInt(params.MouseLocked)
				b := toCInt(params.Borderless)
				d := toCInt(params.Dragable)
				r := toCInt(params.Resizable)
				f := toCInt(params.Fullscreen)
				t, errC := toTString(params.Title)
				if errC == nil {
					mgr, mgrId := registerManager(window)
					errC = C.g2d_window_create(&mgr.data, mgrId, x, y, w, h, wn, hn, wx, hx, b, d, r, f, l, c, t)
					if errC != nil {
						cb.Unregister(int(mgrId))
					}
				}
				C.g2d_string_free(t)
				Err = toError(errC)
			}
		}
	} else {
		panic(notInitialized)
	}
}

func toTString(str string) (unsafe.Pointer, unsafe.Pointer) {
	var strT unsafe.Pointer
	strC := unsafe.Pointer(C.CString(str))
	errC := C.g2d_string_new(&strT, strC)
	C.g2d_string_free(strC)
	return strT, errC
}

func ProcessEvents() {
	if initialized {
		if Err == nil {
			if !processing {
				processing = true
				errC := C.g2d_process_events()
				if errC != nil {
					for i := 0; i < len(cb.mgrs) && cb.mgrs[i] != nil; i++ {
						errC = C.g2d_window_destroy(cb.mgrs[i].data, &errC)
					}
				}
				cb.UnregisterAll()
				Err = toError(errC)
				processing = false
			} else {
				panic(alreadyProcessing)
			}
		}
	} else {
		panic(notInitialized)
	}
}

func newParameters() *Parameters {
	params := new(Parameters)
	params.ClientX = 50
	params.ClientY = 50
	params.ClientWidth = 640
	params.ClientHeight = 480
	params.ClientMinWidth = 0
	params.ClientMinHeight = 0
	params.ClientMaxWidth = 99999
	params.ClientMaxHeight = 99999
	params.MouseLocked = false
	params.Borderless = false
	params.Dragable = false
	params.Resizable = true
	params.Fullscreen = false
	params.Centered = true
	params.Title = "g2d - 0.1.0 世界"
	return params
}

func registerManager(window AbstractWindow) (*tManager, C.int) {
	mgr := new(tManager)
	mgr.wndBase = window.baseStruct()
	mgr.wndAbst = window
	return mgr, C.int(cb.Register(mgr))
}

// toError converts C error to Go error.
func toError(errC unsafe.Pointer) error {
	if errC != nil {
		var errStr string
		var errNumC C.int
		var errWin32 C.g2d_ul_t
		var errStrC *C.char
		C.g2d_error(errC, &errNumC, &errWin32, &errStrC)
		switch errNumC {
		case 1:
			errStr = "memory allocation failed"
		case 2:
			errStr = "get module instance failed"
		case 3:
			errStr = "register dummy class failed"
		case 4:
			errStr = "create dummy window failed"
		case 5:
			errStr = "get dummy device context failed"
		case 6:
			errStr = "choose dummy pixel format failed"
		case 7:
			errStr = "set dummy pixel format failed"
		case 8:
			errStr = "create dummy render context failed"
		case 9:
			errStr = "make dummy context current failed"
		case 10:
			errStr = "release dummy context failed"
		case 11:
			errStr = "deleting dummy render context failed"
		case 12:
			errStr = "destroying dummy window failed"
		case 13:
			errStr = "unregister dummy class failed"
		case 14:
			errStr = "swap dummy buffer failed"
		case 15:
			errStr = "window functions not initialized"
		case 50:
			errStr = "register class failed"
		case 51:
			errStr = "create window failed"
		case 52:
			errStr = "get device context failed"
		case 53:
			errStr = "choose pixel format failed"
		case 54:
			errStr = "set pixel format failed"
		case 55:
			errStr = "create render context failed"
		case 56:
			errStr = "make context current failed"
		case 57:
			errStr = "release context failed"
		case 58:
			errStr = "deleting render context failed"
		case 59:
			errStr = "destroying window failed"
		case 60:
			errStr = "unregister class failed"
		case 61:
			errStr = "swap buffer failed"
		case 62:
			errStr = "set title failed"
		case 63:
			errStr = "wgl functions not initialized"
		case 100:
			C.g2d_error_free(errC)
			return Err
		default:
			errStr = "unknown error " + strconv.FormatUint(uint64(errNumC), 10)
		}
		if errWin32 != 0 {
			errStr = errStr + " (" + strconv.FormatUint(uint64(errWin32), 10) + ")"
		}
		if errStrC != nil {
			errStr = errStr + "; " + C.GoString(errStrC)
		}
		C.g2d_error_free(errC)
		return errors.New(errStr)
	}
	return nil
}

//export g2dOnClose
func g2dOnClose(objIdC C.int) {
	mgr := cb.mgrs[int(objIdC)]
	confirmed, err := mgr.wndAbst.Close()
	if confirmed {
		var errC unsafe.Pointer
		errC = C.g2d_window_destroy(mgr.data, &errC)
		if err == nil && errC != nil {
			err = toError(errC)
		}
	}
	Err = err
}

//export g2dOnDestroyBegin
func g2dOnDestroyBegin(objIdC C.int) {
	mgr := cb.mgrs[int(objIdC)]
	mgr.wndAbst.Destroy()
}

//export g2dOnDestroyEnd
func g2dOnDestroyEnd(objIdC C.int) {
	cb.Unregister(int(objIdC))
}

//export goDebug
func goDebug(a, b C.int, c, d C.g2d_ul_t) {
	fmt.Println(a, b, c, d)
}
