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
	"errors"
	"strconv"
	"unsafe"
)

func Start(engine AbstractEngine) {
	if !running {
		running = true
		err := engine.ParseOSArgs()
		if err == nil {
			base := engine.baseStruct()
			if base.infoOnly {
				engine.Info()
			} else {
				err = initC()
				if err == nil {
					err = engine.CreateWindow()
					err = processEvents(err)
				}
			}
		}
		if err != nil {
			engine.Error(err)
		}
		cb.UnregisterAll()
		running = false
	} else {
		panic("g2d engine already running")
	}
}

func processEvents(err error) error {
	if err == nil {
		var errC unsafe.Pointer
		C.g2d_process_events(&errC)
		if errC != nil {
			for i := 0; i < len(cb.mgrs) && cb.mgrs[i] != nil; i++ {
				C.g2d_window_destroy(cb.mgrs[i].data, &errC)
			}
		}
		err = toError(errC)
	}
	return err
}

func (bulder *WindowBuilder) CreateWindow() error {
	if !initialized {
		var errC unsafe.Pointer
		bulder.ensureParams()
		x := C.int(bulder.ClientX)
		y := C.int(bulder.ClientY)
		w := C.int(bulder.ClientWidth)
		h := C.int(bulder.ClientHeight)
		wn := C.int(bulder.ClientMinWidth)
		hn := C.int(bulder.ClientMinHeight)
		wx := C.int(bulder.ClientMaxWidth)
		hx := C.int(bulder.ClientMaxHeight)
		c := toCInt(bulder.Centered)
		l := toCInt(bulder.MouseLocked)
		b := toCInt(bulder.Borderless)
		d := toCInt(bulder.Dragable)
		r := toCInt(bulder.Resizable)
		f := toCInt(bulder.Fullscreen)
		t := C.g2d_string_new(unsafe.Pointer(&bulder.Title), C.int(len(bulder.Title)), &errC)
		if errC != nil {
			mgr, mgrId := registerManager(bulder.Handler)
			C.g2d_window_create(&mgr.data, mgrId, x, y, w, h, wn, hn, wx, hx, b, d, r, f, l, c, t, &errC)
			if errC != nil {
				cb.Unregister(int(mgrId))
				mgr.handler = nil
			}
		}
		return toError(errC)
	}
	panic("g2d engine not initialized")
}

func registerManager(handler AbstractEventHandler) (*tManager, C.int) {
	mgr := new(tManager)
	mgr.handler = handler
	return mgr, C.int(cb.Register(mgr))
}

func (bulder *WindowBuilder) ensureParams() {
	if bulder.ClientMinWidth < 0 {
		bulder.ClientMinWidth = 0
	}
	if bulder.ClientMinHeight < 0 {
		bulder.ClientMinHeight = 0
	}
	if bulder.ClientMaxWidth <= 0 {
		bulder.ClientMaxWidth = 99999
	}
	if bulder.ClientMaxHeight <= 0 {
		bulder.ClientMaxHeight = 99999
	}
	if len(bulder.Title) == 0 {
		bulder.Title = "OpenGL"
	}
	if bulder.Handler == nil {
		bulder.Handler = new(EventHandler)
	}
}

func initC() error {
	if !initialized {
		var errC unsafe.Pointer
		C.g2d_init(&errC)
		initialized = bool(errC != nil)
		return toError(errC)
	}
	return nil
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
			return cb.mgrs[int(errWin32)].err
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
	confirmed, err := mgr.handler.OnClose()
	if confirmed {
		var errC unsafe.Pointer
		C.g2d_window_destroy(mgr.data, &errC)
		if err == nil && errC != nil {
			err = toError(errC)
		}
	}
	mgr.setError(err, objIdC)
}

//export g2dOnDestroy
func g2dOnDestroy(objIdC C.int) {
	mgr := cb.mgrs[int(objIdC)]
	mgr.handler.OnDestroy()
}

// setError sets err_static to err_num = 100.
func (mgr *tManager) setError(err error, objIdC C.int) {
	if err != nil {
		mgr.err = err
		C.g2d_set_static_err(objIdC)
	}
}
