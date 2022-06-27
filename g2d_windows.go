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

const (
	msgClose = 0
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
			window.baseStruct().Time.NanosEvent = time()
			window.baseStruct().Time.NanosUpdate = window.baseStruct().Time.NanosEvent
			Err = window.Config(params)
			createWindow(window, params)
		}
	} else {
		panic(notInitialized)
	}
}

func createWindow(window AbstractWindow, params *Parameters) {
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
			var data unsafe.Pointer
			mgrId := cb.Register(nil)
			errC = C.g2d_window_create(&data, C.int(mgrId), x, y, w, h, wn, hn, wx, hx, b, d, r, f, l, c, t)
			if errC == nil {
				cb.mgrs[mgrId] = newManager(window, params, data)
				cb.mgrs[mgrId].onCreate(time())
				if Err == nil {
					errC = C.g2d_window_show(data)
				} else {
					cb.mgrs[mgrId].destroy()
				}
			} else {
				cb.Unregister(mgrId)
			}
		}
		C.g2d_string_free(t)
		setErr(toError(errC))
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
					Err = toError(errC)
					for i := 0; i < len(cb.mgrs) && cb.mgrs[i] != nil; i++ {
						cb.mgrs[i].destroy()
					}
				}
				processing = false
			} else {
				panic(alreadyProcessing)
			}
		}
	} else {
		panic(notInitialized)
	}
}

func (mgr *tManagerBase) updateProps() {
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
}

func (mgr *tManagerBase) applyProps(props Properties) {
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

func (mgr *tManagerBase) applyCmd(cmd Command) {
	if cmd.CloseUnc {
		mgr.destroy()
	} else if cmd.CloseReq {
		errC := C.g2d_message_post(mgr.data, msgClose)
		setErr(toError(errC))
	}
}

func (mgr *tManagerBase) createToWindow(nanos int64) error {
	mgr.wndBase.Time.NanosEvent = nanos
	err := mgr.wndAbst.Create()
	return err
}

func (mgr *tManagerBase) showToWindow(nanos int64) error {
	mgr.wndBase.Time.NanosEvent = nanos
	err := mgr.wndAbst.Show()
	return err
}

func (mgr *tManagerBase) keyDownToWindow(code int, repeated uint, nanos int64) error {
	mgr.wndBase.Time.NanosEvent = nanos
	err := mgr.wndAbst.KeyDown(int(code), uint(repeated))
	return err
}

func (mgr *tManagerBase) keyUpToWindow(code int, nanos int64) error {
	mgr.wndBase.Time.NanosEvent = nanos
	err := mgr.wndAbst.KeyUp(int(code))
	return err
}

func (mgr *tManagerBase) closeToWindow(nanos int64) (bool,error) {
	mgr.wndBase.Time.NanosEvent = nanos
	confirmed, err := mgr.wndAbst.Close()
	return confirmed, err
}

func (mgr *tManagerBase) destroyToWindow(nanos int64) {
	mgr.wndBase.Time.NanosEvent = nanos
	mgr.wndAbst.Destroy()
}

func (mgr *tManagerBase) destroy() {
	var errC unsafe.Pointer
	errC = C.g2d_window_destroy(mgr.data, &errC)
	setErr(toError(errC))
}

func (mgr *tManagerNoThreads) onCreate(nanos int64) {
	mgr.updateProps()
	mgr.wndAbst.updatePropsResetCmd(mgr.props)
	err := mgr.createToWindow(nanos)
	setErr(err)
	// only properties - no commands, because window isn't even shown, yet
	if Err == nil && mgr.props != mgr.wndBase.Props {
		// TODO test functionality
		mgr.applyProps(mgr.wndBase.Props)
	}
}

func (mgr *tManagerNoThreads) onShow(nanos int64) {
	mgr.updateProps()
	mgr.wndAbst.updatePropsResetCmd(mgr.props)
	err := mgr.showToWindow(nanos)
	setErr(err)
	mgr.applyPropsAndCmdFromWindow()
}

func (mgr *tManagerNoThreads) onKeyDown(code int, repeated uint, nanos int64) {
	mgr.updateProps()
	mgr.wndAbst.updatePropsResetCmd(mgr.props)
	err := mgr.keyDownToWindow(code, repeated, nanos)
	setErr(err)
	mgr.applyPropsAndCmdFromWindow()
}

func (mgr *tManagerNoThreads) onKeyUp(code int, nanos int64) {
	mgr.updateProps()
	mgr.wndAbst.updatePropsResetCmd(mgr.props)
	err := mgr.keyUpToWindow(code, nanos)
	setErr(err)
	mgr.applyPropsAndCmdFromWindow()
}

func (mgr *tManagerNoThreads) onClose(nanos int64) {
	mgr.updateProps()
	mgr.wndAbst.updatePropsResetCmd(mgr.props)
	confirmed, err := mgr.closeToWindow(nanos)
	setErr(err)
	if confirmed {
		mgr.destroy()
	}
	mgr.applyPropsAndCmdFromWindow()
}

func (mgr *tManagerNoThreads) onDestroy(nanos int64) {
	mgr.updateProps()
	mgr.wndAbst.updatePropsResetCmd(mgr.props)
	mgr.destroyToWindow(nanos)
}

func (mgr *tManagerNoThreads) applyPropsAndCmdFromWindow() {
	if Err == nil {
		if mgr.props != mgr.wndBase.Props {
			mgr.applyProps(mgr.wndBase.Props)
		}
		if mgr.cmd != mgr.wndBase.Cmd {
			mgr.applyCmd(mgr.wndBase.Cmd)
		}
	}
}

/*
func (mgr *tManagerLogicThread) onKeyDown(code int, repeated uint, nanos int64) {
	mgr.updateProps()
	mgr.wndAbst.updatePropsResetCmd(mgr.props)
	err := mgr.keyDownToWindow(code, repeated, nanos)
	setErr(err)
	mgr.applyPropsAndCmdFromWindow()
}
*/

//export g2dShow
func g2dShow(objIdC C.int) {
	cb.mgrs[int(objIdC)].onShow(time())
}

//export g2dKeyDown
func g2dKeyDown(objIdC, code C.int, repeated C.g2d_ui_t) {
	cb.mgrs[int(objIdC)].onKeyDown(int(code), uint(repeated), time())
}

//export g2dKeyUp
func g2dKeyUp(objIdC, code C.int) {
	cb.mgrs[int(objIdC)].onKeyUp(int(code), time())
}

//export g2dClose
func g2dClose(objIdC C.int) {
	cb.mgrs[int(objIdC)].onClose(time())
}

//export g2dDestroyBegin
func g2dDestroyBegin(objIdC C.int) {
	cb.mgrs[int(objIdC)].onDestroy(time())
}

//export g2dDestroyEnd
func g2dDestroyEnd(objIdC C.int) {
	cb.Unregister(int(objIdC))
}

//export goDebug
func goDebug(a, b C.int, c, d C.g2d_ul_t) {
	fmt.Println(a, b, c, d)
}

//export goDebugStr
func goDebugStr(code C.g2d_ul_t, strC C.g2d_lpcstr) {
	fmt.Println(C.GoString(strC), code)
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
		case 64:
			errStr = notInitialized
		case 65:
			// on show
			errStr = messageFailed
		case 66:
			errStr = messageFailed
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
