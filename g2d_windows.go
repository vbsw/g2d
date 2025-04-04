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
	"fmt"
	"strconv"
	"unsafe"
)

const (
	functionFailedDummy = "g2d dummy window %s failed"
	loadFunctionFailed = "g2d load %s function failed"
)

func Init() {
	mutex.Lock()
	if !initialized {
		var n1, n2 C.int
		var err1, err2 C.longlong
		var errInfo *C.char
		C.g2d_init(&n1, &n2, &err1, &err2, &errInfo)
		if err1 == 0 {
			MaxTexSize, MaxTexUnits = int(n1), int(n2)
			initialized, initFailed, quitting = true, false, false
		} else {
			initFailed = true
			Err = toError(int64(err1), int64(err2), errInfo)
		}
		mutex.Unlock()
	} else {
		mutex.Unlock()
		panic("g2d engine is already initialized")
	}
}

func toError(err1, err2 int64, errInfo *C.char) error {
	var err error
	if err1 > 0 {
		var errStr, info string
		if err1 < 1000001 {
			errStr = "memory allocation failed"
		} else if err1 < 1000101 {
			switch err1 {
			case 1000001:
				errStr = "g2d GetModuleHandle failed"
			case 1000002:
				errStr = fmt.Sprintf(functionFailedDummy, "RegisterClassEx")
			case 1000003:
				errStr = fmt.Sprintf(functionFailedDummy, "CreateWindow")
			case 1000004:
				errStr = fmt.Sprintf(functionFailedDummy, "GetDC")
			case 1000005:
				errStr = fmt.Sprintf(functionFailedDummy, "ChoosePixelFormat")
			case 1000006:
				errStr = fmt.Sprintf(functionFailedDummy, "SetPixelFormat")
			case 1000007:
				errStr = fmt.Sprintf(functionFailedDummy, "wglCreateContext")
			case 1000008:
				errStr = fmt.Sprintf(functionFailedDummy, "wglMakeCurrent")
			case 1000009:
				errStr = fmt.Sprintf(functionFailedDummy, "wglMakeCurrent")
			case 1000010:
				errStr = fmt.Sprintf(functionFailedDummy, "wglDeleteContext")
			case 1000011:
				errStr = fmt.Sprintf(functionFailedDummy, "DestroyWindow")
			case 1000012:
				errStr = fmt.Sprintf(functionFailedDummy, "UnregisterClass")
			}
		} else if err1 < 1001001 {
			switch err1 {
			case 1000101:
				errStr = fmt.Sprintf(loadFunctionFailed, "WGL")
			case 1000102:
				errStr = fmt.Sprintf(loadFunctionFailed, "OpenGL")
			}
		}
		if len(errStr) == 0 {
			errStr = "unknown"
		}
		errStr = errStr + " (" + strconv.FormatInt(err1, 10)
		if err2 == 0 {
			errStr = errStr + ")"
		} else {
			errStr = errStr + ", " + strconv.FormatInt(err2, 10) + ")"
		}
		if errInfo != nil {
			info = C.GoString(errInfo)
			if err1 != 1000101 && err1 != 1000102 {
				C.g2d_free(unsafe.Pointer(errInfo))
			}
		}
		if len(info) > 0 {
			errStr = errStr + "; " + info
		}
		err = errors.New(errStr)
	}
	return err
}


/*
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
		params := initWindow(window)
		mgrId, data := createWindow(window, params)
		if Err == nil {
			cb.mgrs[mgrId] = newManager(data, window, params)
			cb.mgrs[mgrId].onCreate(time())
			if Err == nil {
				Err = toError(C.g2d_window_show(data))
			}
			if Err != nil {
				// calls cb.Unregister
				cb.mgrs[mgrId].destroy()
			}
		}
	} else {
		panic(notInitialized)
	}
}

func ProcessEvents() {
	if initialized {
		if Err == nil {
			if !processing {
				processing = true
				pollEvents()
				if Err != nil {
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

func initWindow(window AbstractWindow) *Parameters {
	if Err == nil {
		params := newParameters()
		window.baseStruct().init()
		Err = window.Config(params)
		return params
	}
	return nil
}

func createWindow(window AbstractWindow, params *Parameters) (int, unsafe.Pointer) {
	var mgrId int
	var data unsafe.Pointer
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
			mgrId = cb.Register(nil)
			errC = C.g2d_window_create(&data, C.int(mgrId), x, y, w, h, wn, hn, wx, hx, b, d, r, f, l, c, t)
			C.g2d_string_free(t)
			if errC != nil {
				cb.Unregister(mgrId)
			}
		}
		Err = toError(errC)
	}
	return mgrId, data
}

func toTString(str string) (unsafe.Pointer, unsafe.Pointer) {
	var strT unsafe.Pointer
	strC := unsafe.Pointer(C.CString(str))
	errC := C.g2d_string_new(&strT, strC)
	C.g2d_string_free(strC)
	return strT, errC
}

func pollEvents() {
	errC := C.g2d_process_events()
	Err = toError(errC)
}

func (mgr *tManagerBase) newProps() Properties {
	var x, y, w, h, wn, hn, wx, hx, b, d, r, f, l C.int
	var props Properties
	C.g2d_window_props(mgr.data, &x, &y, &w, &h, &wn, &hn, &wx, &hx, &b, &d, &r, &f, &l)
	props.ClientX = int(x)
	props.ClientY = int(y)
	props.ClientWidth = int(w)
	props.ClientHeight = int(h)
	props.ClientWidthMin = int(wn)
	props.ClientHeightMin = int(hn)
	props.ClientWidthMax = int(wx)
	props.ClientHeightMax = int(hx)
	props.Borderless = bool(b != 0)
	props.Dragable = bool(d != 0)
	props.Resizable = bool(r != 0)
	props.Fullscreen = bool(f != 0)
	props.MouseLocked = bool(l != 0)
	return props
}

func (mgr *tManagerBase) setClientPos(x, y int) {
	C.g2d_client_pos_set(mgr.data, C.int(x), C.int(y))
}

func (mgr *tManagerBase) setClientSize(width, height int) {
	C.g2d_client_size_set(mgr.data, C.int(width), C.int(height))
}

func (mgr *tManagerBase) setWindowStyle(props Properties) {
	wn := C.int(props.ClientWidthMin)
	hn := C.int(props.ClientHeightMin)
	wx := C.int(props.ClientWidthMax)
	hx := C.int(props.ClientHeightMax)
	b := toCInt(props.Borderless)
	d := toCInt(props.Dragable)
	r := toCInt(props.Resizable)
	l := toCInt(props.MouseLocked)
	C.g2d_window_style_set(mgr.data, wn, hn, wx, hx, b, d, r, l)
}

func (mgr *tManagerBase) applyFullscreen() {
	if Err == nil {
		errC := C.g2d_window_fullscreen_set(mgr.data)
		setErr(toError(errC))
	}
}

func (mgr *tManagerBase) applyClientPos() {
	if Err == nil {
		errC := C.g2d_client_pos_apply(mgr.data)
		setErr(toError(errC))
	}
}

func (mgr *tManagerBase) applyClientMove() {
	if Err == nil {
		errC := C.g2d_client_move(mgr.data)
		setErr(toError(errC))
	}
}

func (mgr *tManagerBase) applyCmd(cmd Command) {
	if Err == nil {
		if cmd.CloseUnc {
			mgr.destroy()
		} else if cmd.CloseReq {
			errC := C.g2d_post_close(mgr.data)
			setErr(toError(errC))
		}
	}
}

func (mgr *tManagerBase) applyMouse(x, y int) {
	if Err == nil {
		errC := C.g2d_mouse_pos_set(mgr.data, C.int(x), C.int(y))
		setErr(toError(errC))
	}
}

func (mgr *tManagerBase) applyTitle(title string) {
	if Err == nil {
		t, errC := toTString(title)
		if errC == nil {
			errC = C.g2d_window_title_set(mgr.data, t)
			C.g2d_string_free(t)
		}
		setErr(toError(errC))
	}
}

func (mgr *tManagerBase) applyProps(props Properties, mod tModification) {
	if Err == nil && !mgr.wndBase.destroying {
		if mod.pos {
			mgr.setClientPos(props.ClientX, props.ClientY)
		}
		if mod.size {
			mgr.setClientSize(props.ClientWidth, props.ClientHeight)
		}
		if mod.style {
			mgr.setWindowStyle(props)
		}
		if mod.fsToggle && props.Fullscreen {
			mgr.applyFullscreen()
		} else if mod.fsToggle {
			C.g2d_client_restore_bak(mgr.data)
			mgr.applyClientPos()
		} else if mod.style {
			mgr.applyClientPos()
		} else if mod.pos || mod.size {
			mgr.applyClientMove()
		}
		if mod.mouse {
			mgr.applyMouse(props.MouseX, props.MouseY)
		}
		if mod.title {
			mgr.applyTitle(props.Title)
		}
	}
}

func (mgr *tManagerBase) createToWindow(nanos int64) error {
	mgr.wndBase.Time.NanosEvent = nanos
	err := mgr.wndAbst.Create()
	return err
}

// showToWindow is same as updateToWindow
func (mgr *tManagerBase) showToWindow(nanos int64) error {
	mgr.wndBase.Time.NanosEvent = nanos
	err := mgr.wndAbst.Show()
	mgr.wndBase.Time.NanosUpdate = nanos
	if err == nil && mgr.autoUpdate {
		errC := C.g2d_post_update(mgr.data)
		err = toError(errC)
	}
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

func (mgr *tManagerBase) updateToWindow(nanos int64) error {
	mgr.wndBase.Time.NanosEvent = nanos
	err := mgr.wndAbst.Update()
	mgr.wndBase.Time.NanosUpdate = nanos
	if err == nil && mgr.autoUpdate {
		errC := C.g2d_post_update(mgr.data)
		err = toError(errC)
	}
	return err
}

func (mgr *tManagerBase) closeToWindow(nanos int64) (bool, error) {
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
	mgr.wndBase.destroying = true
	errC = C.g2d_window_destroy(mgr.data, &errC)
	setErr(toError(errC))
}

func (mgr *tManagerNoThreads) applyCmdUpdate() {
	if Err == nil && !mgr.cmd.CloseUnc && !mgr.autoUpdate && mgr.wndBase.Cmd.Update {
		errC := C.g2d_post_update(mgr.data)
		setErr(toError(errC))
	}
}

func (mgr *tManagerNoThreads) onCreate(nanos int64) {
	props := mgr.newProps()
	mgr.wndBase.resetPropsAndCmd(props)
	err := mgr.createToWindow(nanos)
	setErr(err)
	// ignore commands, when window not shown
	mgr.applyProps(mgr.wndBase.Props, mgr.wndBase.modified(props))
}

func (mgr *tManagerNoThreads) onShow(nanos int64) {
	props := mgr.newProps()
	mgr.wndBase.resetPropsAndCmd(props)
	err := mgr.showToWindow(nanos)
	setErr(err)
	mgr.applyCmd(mgr.wndBase.Cmd)
	mgr.applyProps(mgr.wndBase.Props, mgr.wndBase.modified(props))
	mgr.applyCmdUpdate()
}

func (mgr *tManagerNoThreads) onKeyDown(code int, repeated uint, nanos int64) {
	props := mgr.newProps()
	mgr.wndBase.resetPropsAndCmd(props)
	err := mgr.keyDownToWindow(code, repeated, nanos)
	setErr(err)
	mgr.applyCmd(mgr.wndBase.Cmd)
	mgr.applyProps(mgr.wndBase.Props, mgr.wndBase.modified(props))
	mgr.applyCmdUpdate()
}

func (mgr *tManagerNoThreads) onKeyUp(code int, nanos int64) {
	props := mgr.newProps()
	mgr.wndBase.resetPropsAndCmd(props)
	err := mgr.keyUpToWindow(code, nanos)
	setErr(err)
	mgr.applyCmd(mgr.wndBase.Cmd)
	mgr.applyProps(mgr.wndBase.Props, mgr.wndBase.modified(props))
	mgr.applyCmdUpdate()
}

func (mgr *tManagerNoThreads) onUpdate(nanos int64) {
	props := mgr.newProps()
	mgr.wndBase.resetPropsAndCmd(props)
	err := mgr.updateToWindow(nanos)
	setErr(err)
	mgr.applyCmd(mgr.wndBase.Cmd)
	mgr.applyProps(mgr.wndBase.Props, mgr.wndBase.modified(props))
	mgr.applyCmdUpdate()
}

func (mgr *tManagerNoThreads) onClose(nanos int64) {
	props := mgr.newProps()
	mgr.wndBase.resetPropsAndCmd(props)
	confirmed, err := mgr.closeToWindow(nanos)
	setErr(err)
	if confirmed {
		mgr.destroy()
	} else {
		mgr.applyCmd(mgr.wndBase.Cmd)
		mgr.applyProps(mgr.wndBase.Props, mgr.wndBase.modified(props))
		mgr.applyCmdUpdate()
	}
}

func (mgr *tManagerNoThreads) onDestroy(nanos int64) {
	props := mgr.newProps()
	mgr.wndBase.resetPropsAndCmd(props)
	mgr.destroyToWindow(nanos)
}

func (mgr *tManagerNoThreads) onProps(nanos int64) {
}

func (mgr *tManagerNoThreads) onError(nanos int64) {
}

func (mgr *tManagerLogicThread) onKeyDown(code int, repeated uint, nanos int64) {
	props := mgr.newProps()<
	mgr.wndBase.resetPropsAndCmd(props)
	err := mgr.keyDownToWindow(code, repeated, nanos)
	setErr(err)
	mgr.applyProps(mgr.wndBase.Props, mgr.wndBase.modified(props))
}

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

//export g2dUpdate
func g2dUpdate(objIdC C.int) {
	cb.mgrs[int(objIdC)].onUpdate(time())
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

//export g2dProps
func g2dProps(objIdC C.int) {
	cb.mgrs[int(objIdC)].onProps(time())
}

//export g2dError
func g2dError(objIdC C.int) {
	cb.mgrs[int(objIdC)].onError(time())
}

//export goDebug
func goDebug(a, b C.int, c, d C.g2d_ul_t) {
	fmt.Println(a, b, c, d)
}

//export goDebugMessage
func goDebugMessage(code C.g2d_ul_t, strC C.g2d_lpcstr) {
	fmt.Println("Msg:", C.GoString(strC), code)
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
		case 67:
			errStr = messageFailed
		case 68:
			errStr = "set title failed"
		case 69:
			errStr = "set cursor position failed"
		case 70:
			errStr = "set fullscreen failed"
		case 71:
			errStr = "set window position failed"
		case 72:
			errStr = "move window failed"
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
*/
