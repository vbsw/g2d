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
	timepkg "time"
	"fmt"
	"strconv"
	"unsafe"
)

func Init(stub interface{}) {
	if !initialized {
		errC := C.g2d_init()
		timeStart = timepkg.Now()
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
				wn := C.int(params.ClientWidthMin)
				hn := C.int(params.ClientHeightMin)
				wx := C.int(params.ClientWidthMax)
				hx := C.int(params.ClientHeightMax)
				c := toCInt(params.Centered)
				l := toCInt(params.MouseLocked)
				b := toCInt(params.Borderless)
				d := toCInt(params.Dragable)
				r := toCInt(params.Resizable)
				f := toCInt(params.Fullscreen)
				t, errC := toTString(params.Title)
				if errC == nil {
					mgr, mgrId := registerNewManager(window, params.Title)
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
	params.ClientWidthMin = 0
	params.ClientHeightMin = 0
	params.ClientWidthMax = 99999
	params.ClientHeightMax = 99999
	params.MouseLocked = false
	params.Borderless = false
	params.Dragable = false
	params.Resizable = true
	params.Fullscreen = false
	params.Centered = true
	params.Title = "g2d - 0.1.0"
	return params
}

func registerNewManager(window AbstractWindow, title string) (*tManager, C.int) {
	mgr := new(tManager)
	mgr.wndBase = window.baseStruct()
	mgr.wndAbst = window
	mgr.props.Title = title
	return mgr, C.int(cb.Register(mgr))
}

func (mgr *tManager) updatePropsResetCmd() {
	var x, y, w, h, wn, hn, wx, hx, b, d, r, f, l C.int
	C.g2d_window_props(mgr.data, &x, &y, &w, &h, &wn, &hn, &wx, &hx, &b, &d, &r, &f, &l)
	mgr.props.ClientX = int(x)
	mgr.props.ClientY = int(y)
	mgr.props.ClientWidth = int(w)
	mgr.props.ClientHeight = int(h)
	mgr.props.ClientWidthMin = int(wn)
	mgr.props.ClientHeightMin = int(hn)
	mgr.props.ClientWidthMax = int(wx)
	mgr.props.ClientHeightMax = int(hx)
	mgr.props.Borderless = bool(b != 0)
	mgr.props.Dragable = bool(d != 0)
	mgr.props.Resizable = bool(r != 0)
	mgr.props.Fullscreen = bool(f != 0)
	mgr.props.MouseLocked = bool(l != 0)
	mgr.wndAbst.updatePropsResetCmd(mgr.props)
}

func (mgr *tManager) applyProps(props Properties) {
	x := C.int(mgr.props.ClientX)
	y := C.int(mgr.props.ClientY)
	w := C.int(mgr.props.ClientWidth)
	h := C.int(mgr.props.ClientHeight)
	wn := C.int(mgr.props.ClientWidthMin)
	hn := C.int(mgr.props.ClientHeightMin)
	wx := C.int(mgr.props.ClientWidthMax)
	hx := C.int(mgr.props.ClientHeightMax)
	b := toCInt(mgr.props.Borderless)
	d := toCInt(mgr.props.Dragable)
	r := toCInt(mgr.props.Resizable)
	f := toCInt(props.Fullscreen)
	l := toCInt(mgr.props.MouseLocked)
	C.g2d_window_props_apply(mgr.data, x, y, w, h, wn, hn, wx, hx, b, d, r, f, l)
}

func (mgr *tManager) applyCmd(cmd Command) {
	if cmd.CloseUnc {
		mgr.destroy()
	} else if cmd.CloseReq {
		C.g2d_message_close_post(mgr.data)
	}
}

func (mgr *tManager) destroy() {
	var errC unsafe.Pointer
	errC = C.g2d_window_destroy(mgr.data, &errC)
	setErr(toError(errC))
}

func (mgr *tManager) applyPropsAndCmd() {
	if Err == nil {
		props, cmd := mgr.wndAbst.propsAndCmd()
		if mgr.props != props {
			mgr.applyProps(props)
		}
		if mgr.cmd != cmd {
			mgr.applyCmd(cmd)
		}
	}
}

//export g2dKeyDown
func g2dKeyDown(objIdC, code C.int, repeated C.g2d_ui_t) {
	nanos := time()
	mgr := cb.mgrs[int(objIdC)]
	mgr.updatePropsResetCmd()
	mgr.wndBase.Time.NanosUpdateCurr = nanos
	err := mgr.wndAbst.KeyDown(int(code), uint(repeated), nanos)
	setErr(err)
	mgr.applyPropsAndCmd()
}

//export g2dKeyUp
func g2dKeyUp(objIdC, code C.int) {
	nanos := time()
	mgr := cb.mgrs[int(objIdC)]
	mgr.updatePropsResetCmd()
	mgr.wndBase.Time.NanosUpdateCurr = nanos
	err := mgr.wndAbst.KeyUp(int(code), nanos)
	setErr(err)
	mgr.applyPropsAndCmd()
}

//export g2dClose
func g2dClose(objIdC C.int) {
	nanos := time()
	mgr := cb.mgrs[int(objIdC)]
	mgr.updatePropsResetCmd()
	mgr.wndBase.Time.NanosUpdateCurr = nanos
	confirmed, err := mgr.wndAbst.Close(nanos)
	setErr(err)
	if confirmed {
		mgr.destroy()
	}
	mgr.applyPropsAndCmd()
}

//export g2dDestroyBegin
func g2dDestroyBegin(objIdC C.int) {
	nanos := time()
	mgr := cb.mgrs[int(objIdC)]
	mgr.updatePropsResetCmd()
	mgr.wndBase.Time.NanosUpdateCurr = nanos
	mgr.wndAbst.Destroy(nanos)
}

//export g2dDestroyEnd
func g2dDestroyEnd(objIdC C.int) {
	cb.Unregister(int(objIdC))
}

//export goDebug
func goDebug(a, b C.int, c, d C.g2d_ul_t) {
	fmt.Println(a, b, c, d)
}

func setErr(err error) {
	if err != nil && Err == nil {
		Err = err
		C.g2d_err_static_set(0)
	}
}

// toError converts C struct to Go error.
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
			errStr = "unknown error"
		}
		errStr = errStr + " (" + strconv.FormatUint(uint64(errNumC), 10)
		if errWin32 == 0 {
			errStr = errStr + ")"
		} else {
			errStr = errStr + ", " + strconv.FormatUint(uint64(errWin32), 10) + ")"
		}
		if errStrC != nil {
			errStr = errStr + "; " + C.GoString(errStrC)
		}
		C.g2d_error_free(errC)
		return errors.New(errStr)
	}
	return nil
}
