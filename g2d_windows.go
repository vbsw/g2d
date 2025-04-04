/*
 *          Copyright 2025, Vitali Baumtrok.
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
	loadFunctionFailed  = "g2d load %s function failed"
)

func Init() {
	mutex.Lock()
	if !initialized {
		var n1, n2 C.int
		var err1, err2 C.longlong
		var errInfo *C.char
		C.g2d_init(&n1, &n2, &err1, &err2, &errInfo)
		if err1 == 0 {
			Err = nil
			MaxTexSize, MaxTexUnits = int(n1), int(n2)
			initialized, initFailed, quitting = true, false, false
		} else {
			initFailed = true
			Err = toError(err1, err2, errInfo)
		}
		mutex.Unlock()
	} else {
		mutex.Unlock()
		panic("g2d engine is already initialized")
	}
}

func MainLoop(mainWindow Window) {
	if mainWindow != nil {
		mutex.Lock()
		if !initFailed {
			if initialized {
				if !running {
					running = true
					wnds = make([]*tWindow, 0, 2)
					wndNextId = make([]int, 0, 2)
					requests = make([]tRequest, 0, 2)
					wnd := newWindow(mainWindow)
					go wnd.logicThread()
					mutex.Unlock()
					C.g2d_main_loop()
					mutex.Lock()
					running = false
					mutex.Unlock()
					cleanUp()
				} else {
					mutex.Unlock()
					panic(alreadyRunning)
				}
			} else {
				mutex.Unlock()
				panic(notInitialized)
			}
		} else {
			mutex.Unlock()
		}
	} else {
		panic(mustNotBeNil)
	}
}

func (props *Properties) update(data unsafe.Pointer) {
	var mx, my, x, y, w, h, wn, hn, wx, hx, b, d, r, f, l C.int
	C.g2d_window_props(data, &mx, &my, &x, &y, &w, &h, &wn, &hn, &wx, &hx, &b, &d, &r, &f, &l)
	props.MouseX = int(mx)
	props.MouseY = int(my)
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
}

func (request *tCreateWindowRequest) process() {
	var err1, err2 C.longlong
	var data unsafe.Pointer
	var t unsafe.Pointer
	var ts C.size_t
	x := C.int(request.config.ClientX)
	y := C.int(request.config.ClientY)
	w := C.int(request.config.ClientWidth)
	h := C.int(request.config.ClientHeight)
	wn := C.int(request.config.ClientWidthMin)
	hn := C.int(request.config.ClientHeightMin)
	wx := C.int(request.config.ClientWidthMax)
	hx := C.int(request.config.ClientHeightMax)
	c, l, b, d, r, f := request.config.boolsToCInt()
	if len(request.config.Title) > 0 {
		bytes := *(*[]byte)(unsafe.Pointer(&(request.config.Title)))
		t, ts = unsafe.Pointer(&bytes[0]), C.size_t(len(request.config.Title))
	}
	C.g2d_window_create(&data, C.int(request.wndId), x, y, w, h, wn, hn, wx, hx, b, d, r, f, l, c, t, ts, &err1, &err2)
	if err1 == 0 {
		wnd := wnds[request.wndId]
		wnd.data = data
		event := &tLogicEvent{typeId: createType}
		event.props.update(data)
		wnd.eventsChan <- event
	} else {
		Err = toError(err1, err2, nil)
	}
}

func (request *tShowWindowRequest) process() {
	var err1, err2 C.longlong
	wnd := wnds[request.wndId]
	C.g2d_window_show(wnd.data, &err1, &err2)
	if err1 == 0 {
		event := &tLogicEvent{typeId: showType}
		event.props.update(wnd.data)
		wnd.eventsChan <- event
	} else {
		Err = toError(err1, err2, nil)
	}
}

func (request *tCloseWindowRequest) process() {
	wnd := wnds[request.wndId]
	event := &tLogicEvent{typeId: closeType}
	event.props.update(wnd.data)
	wnd.eventsChan <- event
}

func (request *tDestroyWindowRequest) process() {
	var err1, err2 C.longlong
	wnd := wnds[request.wndId]
	C.g2d_window_destroy(wnd.data, &err1, &err2)
	if err1 == 0 {
	} else {
		Err = toError(err1, err2, nil)
	}
}

func (request *tCustomRequest) process() {
	wnd := wnds[request.wndId]
	event := &tLogicEvent{typeId: customType, obj: request.obj}
	event.props.update(wnd.data)
	wnd.eventsChan <- event
}

func (request *tSetPropertiesRequest) process() {
	var err1, err2 C.longlong
	wnd := wnds[request.wndId]
	if request.modPosSize {
		C.g2d_window_pos_size_set(wnd.data, C.int(request.props.ClientX), C.int(request.props.ClientY), C.int(request.props.ClientWidth), C.int(request.props.ClientHeight))
	}
	if request.modStyle || request.modFullscreen {
		wn := C.int(request.props.ClientWidthMin)
		hn := C.int(request.props.ClientHeightMin)
		wx := C.int(request.props.ClientWidthMax)
		hx := C.int(request.props.ClientHeightMax)
		l, b, d, r, f := request.props.boolsToCInt()
		C.g2d_window_style_set(wnd.data, wn, hn, wx, hx, b, d, r, f, l)
	}
	if request.modFullscreen {
		C.g2d_window_fullscreen_set(wnd.data, &err1, &err2)
	} else if request.modStyle && !request.props.Fullscreen {
		C.g2d_window_pos_apply(wnd.data, &err1, &err2)
	} else if request.modPosSize {
		C.g2d_window_move(wnd.data, &err1, &err2)
	}
	if request.modTitle {
		var t unsafe.Pointer
		var ts C.size_t
		if len(request.props.Title) > 0 {
			bytes := *(*[]byte)(unsafe.Pointer(&(request.props.Title)))
			t, ts = unsafe.Pointer(&bytes[0]), C.size_t(len(request.props.Title))
		}
		C.g2d_window_title_set(wnd.data, t, ts, &err1, &err2)
	}
	if request.modMouse {
		C.g2d_mouse_pos_set(wnd.data, C.int(request.props.MouseX), C.int(request.props.MouseY), &err1, &err2)
	}
}

func (request *tUpdateRequest) process() {
	wnd := wnds[request.wndId]
	event := &tLogicEvent{typeId: updateType}
	event.props.update(wnd.data)
	wnd.eventsChan <- event
	wnd.update = false
}

func postRequest(request tRequest) {
	var err1, err2 C.longlong
	mutex.Lock()
	requests = append(requests, request)
	C.g2d_post_request(&err1, &err2)
	mutex.Unlock()
}

func postUpdateRequest(wndId int) {
	mutex.Lock()
	wnd := wnds[wndId]
	if !wnd.update {
		var err1, err2 C.longlong
		wnd.update = true
		requests = append(requests, &tUpdateRequest{wndId: wndId})
		C.g2d_post_request(&err1, &err2)
	}
	mutex.Unlock()
}

func cleanUp() {
	/*
		for _, abst := range abstCbs {
			if abst != nil {
				var err1, err2 C.longlong
				wnd := abst.impl()
				wnd.msgs <- (&tLogicMessage{typeId: quitType, nanos: time.Nanos()})
				<-wnd.quitted
				unregister(wnd.cbId)
				C.g2d_window_destroy(wnd.dataC, &err1, &err2)
				if err1 != 0 {
					setError(toError(int64(err1), 0, int64(wnd.cbId), "", nil))
				}
			}
		}
		toMainLoop.quitMessageThread()
	*/
	C.g2d_clean_up()
}

func toEventsChan(id C.int, event *tLogicEvent) {
	if !processingRequests {
		mutex.Lock()
	}
	wnd := wnds[id]
	event.props.update(wnd.data)
	wnd.eventsChan <- event
	if !processingRequests {
		mutex.Unlock()
	}
}

//export g2dMainLoopStarted
func g2dMainLoopStarted() {
	wnds[0].eventsChan <- &tLogicEvent{typeId: configType}
}

//export g2dProcessRequest
func g2dProcessRequest() {
	mutex.Lock()
	processingRequests = true
	for _, request := range requests {
		request.process()
	}
	requests = requests[:0]
	processingRequests = false
	mutex.Unlock()
}

//export g2dClose
func g2dClose(id C.int) {
	toEventsChan(id, &tLogicEvent{typeId: closeType})
}

//export g2dKeyDown
func g2dKeyDown(id, code C.int, repeated C.uint) {
	toEventsChan(id, &tLogicEvent{typeId: keyDownType, valA: int(code), repeated: uint(repeated)})
}

//export g2dKeyUp
func g2dKeyUp(id, code C.int) {
	toEventsChan(id, &tLogicEvent{typeId: keyUpType, valA: int(code)})
}

//export g2dMouseMove
func g2dMouseMove(id C.int) {
	toEventsChan(id, &tLogicEvent{typeId: msMoveType})
}

//export g2dWindowMove
func g2dWindowMove(id C.int) {
	toEventsChan(id, &tLogicEvent{typeId: wndMoveType})
}

//export g2dWindowResize
func g2dWindowResize(id C.int) {
	if !processingRequests {
		mutex.Lock()
	}
	wnd := wnds[id]
	event := &tLogicEvent{typeId: wndResizeType}
	event.props.update(wnd.data)
	wnd.eventsChan <- event
	// wnd.gfx.eventsChan <- &tGraphicsEvent{typeId: wndResizeType, valA: event.props.ClientWidth, valB: event.props.ClientHeight}
	if !processingRequests {
		mutex.Unlock()
	}
}

//export g2dButtonDown
func g2dButtonDown(id, code, doubleClick C.int) {
	toEventsChan(id, &tLogicEvent{typeId: buttonDownType, valA: int(code), repeated: uint(doubleClick)})
}

//export g2dButtonUp
func g2dButtonUp(id, code, doubleClick C.int) {
	toEventsChan(id, &tLogicEvent{typeId: buttonUpType, valA: int(code), repeated: uint(doubleClick)})
}

//export g2dWheel
func g2dWheel(id C.int, wheel C.float) {
	toEventsChan(id, &tLogicEvent{typeId: wheelType, valB: float32(wheel)})
}

//export g2dWindowMinimize
func g2dWindowMinimize(id C.int) {
	toEventsChan(id, &tLogicEvent{typeId: minimizeType})
}

//export g2dWindowRestore
func g2dWindowRestore(id C.int) {
	toEventsChan(id, &tLogicEvent{typeId: restoreType})
}

func toError(err1, err2 C.longlong, errInfo *C.char) error {
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
		errStr = errStr + " (" + strconv.FormatInt(int64(err1), 10)
		if err2 == 0 {
			errStr = errStr + ")"
		} else {
			errStr = errStr + ", " + strconv.FormatInt(int64(err2), 10) + ")"
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

//export goDebug
func goDebug(a, b C.int, c, d C.longlong) {
	fmt.Println(a, b, c, d)
}

//export goDebugMessage
func goDebugMessage(code C.ulong, strC *C.char) {
	fmt.Println("Msg:", C.GoString(strC), code)
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
