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
	"time"
	"unsafe"
)

func Init(stub interface{}) {
	if !initialized {
		var errNumC C.int
		var errWin32C C.g2d_ul_t
		var errStrC *C.char
		C.g2d_init(&errNumC, &errWin32C, &errStrC)
		Err = toError(errNumC, errWin32C, errStrC)
		if Err == nil {
			startTime = time.Now()
			initialized = true
		}
	}
}

func Show(window interface{}) {
	if Err == nil {
		if initialized {
			if window != nil {
				wnd, ok := window.(tWindow)
				if ok {
					config := wnd.config(wnd)
					wnd.create(config)
					wnd.show()
				} else {
					panic(notEmbedded)
				}
			}
		} else {
			panic(notInitialized)
		}
	}
}

func (window *Window) create(config *Configuration) {
	if Err == nil {
		x := C.int(config.ClientX)
		y := C.int(config.ClientY)
		w := C.int(config.ClientWidth)
		h := C.int(config.ClientHeight)
		wn := C.int(config.ClientWidthMin)
		hn := C.int(config.ClientHeightMin)
		wx := C.int(config.ClientWidthMax)
		hx := C.int(config.ClientHeightMax)
		c := toCInt(config.Centered)
		l := toCInt(config.MouseLocked)
		b := toCInt(config.Borderless)
		d := toCInt(config.Dragable)
		r := toCInt(config.Resizable)
		f := toCInt(config.Fullscreen)
		t, errNumC := toTString(config.Title)
		if errNumC == 0 {
			var dataC unsafe.Pointer
			var errWin32C C.g2d_ul_t
			cbId := cb.Register(window.abst)
			C.g2d_window_create(&dataC, C.int(cbId), x, y, w, h, wn, hn, wx, hx, b, d, r, f, l, c, t, &errNumC, &errWin32C)
			C.g2d_free(t)
			if errNumC == 0 {
				window.dataC = dataC
				window.onCreate()
			} else {
				Err = toError(errNumC, 0, nil)
			}
		} else {
			Err = toError(-18, 0, nil)
		}
	}
}

func (window *Window) show() {
	if Err == nil {
		var errNumC C.int
		var errWin32C C.g2d_ul_t
		C.g2d_window_show(window.dataC, &errNumC, &errWin32C)
		if errNumC == 0 {
			window.onShow()
		} else {
			Err = toError(errNumC, errWin32C, nil)
		}
	}
}

func toTString(str string) (unsafe.Pointer, C.int) {
	var strT unsafe.Pointer
	var errNumC C.int
	strC := unsafe.Pointer(C.CString(str))
	if strC != nil {
		C.g2d_to_tstr(&strT, strC, &errNumC)
		C.g2d_free(strC)
	} else {
		errNumC = 2
	}
	return strT, errNumC
}

func (window *Window) onCreate() {
}

func (window *Window) onShow() {
}

/*
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

*/

/*
func (mgr *tManagerLogicThread) onKeyDown(code int, repeated uint, nanos int64) {
	props := mgr.newProps()<
	mgr.wndBase.resetPropsAndCmd(props)
	err := mgr.keyDownToWindow(code, repeated, nanos)
	setErr(err)
	mgr.applyProps(mgr.wndBase.Props, mgr.wndBase.modified(props))
}
*/

//export g2dShow
func g2dShow(objIdC C.int) {
	//cb.wnds[int(objIdC)].onShow()
}

//export g2dKeyDown
func g2dKeyDown(objIdC, code C.int, repeated C.g2d_ui_t) {
	//cb.wnds[int(objIdC)].onKeyDown(int(code), uint(repeated))
}

//export g2dKeyUp
func g2dKeyUp(objIdC, code C.int) {
	//cb.wnds[int(objIdC)].onKeyUp(int(code), time())
}

//export g2dUpdate
func g2dUpdate(objIdC C.int) {
	//cb.wnds[int(objIdC)].onUpdate()
}

//export g2dClose
func g2dClose(objIdC C.int) {
	//cb.wnds[int(objIdC)].onClose()
}

//export g2dDestroyBegin
func g2dDestroyBegin(objIdC C.int) {
	//cb.wnds[int(objIdC)].onDestroy()
}

//export g2dDestroyEnd
func g2dDestroyEnd(objIdC C.int) {
	cb.Unregister(int(objIdC))
}

//export g2dProps
func g2dProps(objIdC C.int) {
	//cb.wnds[int(objIdC)].onProps(time())
}

//export g2dError
func g2dError(objIdC C.int) {
	//cb.wnds[int(objIdC)].onError(time())
}

func setErr(err error) {
	if err != nil && Err == nil {
		Err = err
	}
}

func toError(errNumC C.int, errWin32C C.g2d_ul_t, errStrC *C.char) error {
	if errNumC != 0 {
		var errStr string
		if errNumC < 0 {
			errStr = memoryAllocation
			errNumC = -1 * errNumC
		} else {
			switch errNumC {
			case 1:
				errStr = "get module instance failed"
			case 2:
				errStr = memoryAllocation
			case 10:
				errStr = "register dummy class failed"
			case 11:
				errStr = "create dummy window failed"
			case 12:
				errStr = "get dummy device context failed"
			case 13:
				errStr = "choose dummy pixel format failed"
			case 14:
				errStr = "set dummy pixel format failed"
			case 15:
				errStr = "create dummy render context failed"
			case 16:
				errStr = "make dummy context current failed"
			case 17:
				errStr = "load OpenGL function failed"
			case 18:
				errStr = "release dummy context failed"
			case 19:
				errStr = "deleting dummy render context failed"
			case 20:
				errStr = "destroying dummy window failed"
			case 21:
				errStr = "unregister dummy class failed"
			case 22:
				errStr = "swap dummy buffer failed"
			case 30:
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
				errStr = "set title failed"
			case 66:
				errStr = "set cursor position failed"
			case 67:
				errStr = "set fullscreen failed"
			case 68:
				errStr = "set window position failed"
			case 69:
				errStr = "move window failed"
			case 80:
				errStr = messageFailed
			case 81:
				errStr = messageFailed
			case 82:
				errStr = messageFailed
			case 83:
				errStr = messageFailed
			default:
				errStr = "unknown error"
			}
		}
		errStr = errStr + " (" + strconv.FormatUint(uint64(errNumC), 10)
		if errWin32C == 0 {
			errStr = errStr + ")"
		} else {
			errStr = errStr + ", " + strconv.FormatUint(uint64(errWin32C), 10) + ")"
		}
		if errStrC != nil {
			errStr = errStr + "; " + C.GoString(errStrC)
			C.g2d_free(unsafe.Pointer(errStrC))
		}
		return errors.New(errStr)
	}
	return nil
}

//export goDebug
func goDebug(a, b C.int, c, d C.g2d_ul_t) {
	fmt.Println(a, b, c, d)
}

//export goDebugMessage
func goDebugMessage(code C.g2d_ul_t, strC C.g2d_lpcstr) {
	fmt.Println("Msg:", C.GoString(strC), code)
}
