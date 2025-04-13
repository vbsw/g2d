/*
 *          Copyright 2025, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

package g2d

// #cgo CFLAGS: -DG2D_WIN32 -DUNICODE
// #cgo LDFLAGS: -luser32 -lgdi32 -lOpenGL32
// #cgo noescape g2d_gfx_draw
// #cgo nocallback g2d_gfx_draw
// #include "g2d.h"
import "C"
import (
	"errors"
	"fmt"
	"runtime"
	"strconv"
	"unsafe"
)

const (
	functionFailedDummy    = "dummy window %s failed"
	loadFunctionFailed     = "load %s function failed"
	functionFailedWindow   = "window %s failed"
	functionFailedGraphics = "graphics %s failed"
	functionFailedG2D      = "%s failed"
)

// Init initialized the g2d framework.
func Init() {
	mutex.Lock()
	if !initialized {
		var numbers [4]C.int
		var err1, err2 C.longlong
		var errInfo *C.char
		C.g2d_init(&numbers[0], &err1, &err2, &errInfo)
		if err1 == 0 {
			Err = nil
			MaxTexSize, MaxTexUnits = int(numbers[0]), int(numbers[1])
			VSyncAvailable, AVSyncAvailable = (numbers[2] != 0), (numbers[3] != 0)
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

// MainLoop processes events and will initialize and show mainWindow.
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
					appTime.Reset()
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

func (props *Properties) update(data unsafe.Pointer, title string) {
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
	props.Title = title
}

func (layer *RectanglesLayer) getProcessing(layers []tLayer) ([]tLayer, []C.float, int, unsafe.Pointer) {
	var count, index int
	for len(layers) > 0 {
		if curr, ok := layers[0].(*RectanglesLayer); ok {
			if curr.Enabled && curr.count > 0 {
				layer.buffer = ensureCFloatLen(layer.buffer, (count+curr.count)*(2+2+4))
				for _, entity := range layer.entities {
					if entity.Enabled {
						layer.buffer[index+0] = C.float(entity.X)
						layer.buffer[index+1] = C.float(entity.Y)
						layer.buffer[index+2] = C.float(entity.Width)
						layer.buffer[index+3] = C.float(entity.Height)
						layer.buffer[index+4] = C.float(entity.R)
						layer.buffer[index+5] = C.float(entity.G)
						layer.buffer[index+6] = C.float(entity.B)
						layer.buffer[index+7] = C.float(entity.A)
						index += (2 + 2 + 4)
						count++
					}
				}
			}
			layers = layers[1:]
		} else {
			break
		}
	}
	return layers, layer.buffer, count, unsafe.Pointer(C.g2d_gfx_draw_rectangles)
}

func (wnd *tWindow) graphicsThread() {
	var err1, err2 C.longlong
	var errInfo *C.char
	runtime.LockOSThread()
	C.g2d_gfx_init(wnd.data, &err1, &err2, &errInfo)
	if err1 == 0 {
		wnd.impl.Gfx.running = true
		for wnd.impl.Gfx.running {
			event := <-wnd.impl.Gfx.eventsChan
			if event != nil {
				if event.err == nil {
					switch event.typeId {
					case refreshType:
						wnd.onGfxRefresh()
					case wndResizeType:
						wnd.impl.Gfx.w, wnd.impl.Gfx.h = event.valA, event.valB
					case leaveType:
						wnd.impl.Gfx.running = false
					case imageType:
						var texId C.int
						C.g2d_gfx_gen_tex(wnd.data, unsafe.Pointer(&event.valD[0]), C.int(event.valA), C.int(event.valB), &texId, &err1)
						if err1 == 0 {
							internalId := len(wnd.impl.Gfx.texturesIds)
							wnd.impl.Gfx.texturesIds = append(wnd.impl.Gfx.texturesIds, C.int(texId))
							wnd.eventsChan <- &tLogicEvent{typeId: textureType, valA: event.valC, valB: internalId, time: appTime.Millis()}
						} else {
							postRequest(&tErrorRequest{err: toError(err1, 0, nil)})
						}
					}
				} else {
					postRequest(&tErrorRequest{err: event.err})
				}
			}
		}
		C.g2d_gfx_release(wnd.data, &err1, &err2)
		if err1 != 0 {
			postRequest(&tErrorRequest{err: toError(err1, err2, nil)})
		}
	} else {
		postRequest(&tErrorRequest{err: toError(err1, err2, errInfo)})
	}
	wnd.impl.Gfx.quittedChan <- true
}

func (wnd *tWindow) onGfxRefresh() {
	var err1, err2 C.longlong
	wnd.impl.Gfx.mutex.Lock()
	wnd.impl.Gfx.updating = false
	read := wnd.impl.Gfx.getReadBuffer()
	wnd.impl.Gfx.mutex.Unlock()
	buffers, lengths, procs := read.getBatchProcessing()
	w, h, i, r, g, b := read.w, read.h, read.si, read.r, read.g, read.b
	if len(buffers) > 0 {
		// calling with &buffers[0] may cause "pointer to unpinned Go pointer" error
		// https://github.com/PowerDNS/lmdb-go/issues/28
		var pinner runtime.Pinner
		for _, buf := range buffers {
			pinner.Pin(buf)
		}
		C.g2d_gfx_draw(wnd.data, w, h, i, r, b, g, &buffers[0], &lengths[0], &procs[0], C.int(len(buffers)), &err1, &err2)
		pinner.Unpin()
	} else {
		// just draw background
		C.g2d_gfx_draw(wnd.data, w, h, i, r, b, g, nil, nil, nil, 0, &err1, &err2)
	}
	wnd.eventsChan <- &tLogicEvent{typeId: refreshType, time: appTime.Millis()}
}

func (gfx *tGraphics) getBatchProcessing() ([]*C.float, []C.int, []unsafe.Pointer) {
	var buffers []*C.float
	var lengths []C.int
	var procs []unsafe.Pointer
	layers := gfx.layers
	// multiple layers may be "batched" together, until no layer left
	for len(layers) > 0 {
		var buffer []C.float
		var length int
		var proc unsafe.Pointer
		layers, buffer, length, proc = layers[0].getProcessing(layers)
		if len(buffer) > 0 {
			buffers = append(buffers, &buffer[0])
			lengths = append(lengths, C.int(length))
			procs = append(procs, proc)
		}
	}
	return buffers, lengths, procs
}

func (request *tConfigWindowRequest) process() {
	wnd := newWindow(request.window)
	go wnd.logicThread()
	wnd.eventsChan <- &tLogicEvent{typeId: configType, time: appTime.Millis()}
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
		wnd.title = request.config.Title
		event := &tLogicEvent{typeId: createType, time: appTime.Millis()}
		event.props.update(data, wnd.title)
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
		event := &tLogicEvent{typeId: showType, time: appTime.Millis()}
		event.props.update(wnd.data, wnd.title)
		wnd.eventsChan <- event
	} else {
		Err = toError(err1, err2, nil)
	}
}

func (request *tCloseWindowRequest) process() {
	wnd := wnds[request.wndId]
	event := &tLogicEvent{typeId: closeType, time: appTime.Millis()}
	event.props.update(wnd.data, wnd.title)
	wnd.eventsChan <- event
}

func (request *tDestroyWindowRequest) process() {
	var err1, err2 C.longlong
	wnd := wnds[request.wndId]
	wnd.eventsChan <- &tLogicEvent{typeId: destroyType, time: appTime.Millis()}
	<-wnd.quittedChan
	C.g2d_window_destroy(wnd.data, &err1, &err2)
	wnd.data = nil
	unregisterWnd(wnd.id).impl = nil
	if err1 != 0 {
		(&tErrorRequest{err: toError(err1, err2, nil)}).process()
	}
}

func (request *tCustomRequest) process() {
	wnd := wnds[request.wndId]
	event := &tLogicEvent{typeId: customType, obj: request.obj}
	event.props.update(wnd.data, wnd.title)
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
		wnd.title = request.props.Title
		C.g2d_window_title_set(wnd.data, t, ts, &err1, &err2)
	}
	if request.modMouse {
		C.g2d_mouse_pos_set(wnd.data, C.int(request.props.MouseX), C.int(request.props.MouseY), &err1, &err2)
	}
}

func (request *tErrorRequest) process() {
	if Err == nil {
		Err = request.err
	}
	if running {
		var err1, err2 C.longlong
		C.g2d_post_quit(&err1, &err2)
		if err1 != 0 {
			if Err == nil {
				Err = toError(err1, err2, nil)
			}
		}
	}
}

func postRequest(request tRequest) {
	var err1, err2 C.longlong
	mutex.Lock()
	requests = append(requests, request)
	C.g2d_post_request(&err1, &err2)
	if err1 != 0 {
		(&tErrorRequest{err: toError(err1, err2, nil)}).process()
	}
	mutex.Unlock()
}

func cleanUp() {
	for _, wnd := range wnds {
		if wnd != nil {
			var err1, err2 C.longlong
			if wnd.data != nil {
				wnd.eventsChan <- &tLogicEvent{typeId: destroyType, err: Err, time: appTime.Millis()}
				<-wnd.quittedChan
				C.g2d_window_destroy(wnd.data, &err1, &err2)
				wnd.data = nil
				if err1 != 0 {
					(&tErrorRequest{err: toError(err1, err2, nil)}).process()
				}
			} else {
				wnd.eventsChan <- &tLogicEvent{typeId: leaveType, time: appTime.Millis()}
				<-wnd.quittedChan
			}
			unregisterWnd(wnd.id).impl = nil
		}
	}
	C.g2d_clean_up()
}

func postLogicEvent(id C.int, event *tLogicEvent) {
	if !processingRequests {
		mutex.Lock()
	}
	wnd := wnds[id]
	event.props.update(wnd.data, wnd.title)
	wnd.eventsChan <- event
	if !processingRequests {
		mutex.Unlock()
	}
}

//export g2dMainLoopStarted
func g2dMainLoopStarted() {
	wnds[0].eventsChan <- &tLogicEvent{typeId: configType, time: appTime.Millis()}
}

//export g2dProcessRequest
func g2dProcessRequest() {
	mutex.Lock()
	processingRequests = true
	for _, request := range requests {
		if Err == nil {
			request.process()
		}
	}
	requests = requests[:0]
	processingRequests = false
	mutex.Unlock()
}

//export g2dClose
func g2dClose(id C.int) {
	postLogicEvent(id, &tLogicEvent{typeId: closeType, time: appTime.Millis()})
}

//export g2dKeyDown
func g2dKeyDown(id, code C.int, repeated C.uint) {
	postLogicEvent(id, &tLogicEvent{typeId: keyDownType, valA: int(code), repeated: uint(repeated), time: appTime.Millis()})
}

//export g2dKeyUp
func g2dKeyUp(id, code C.int) {
	postLogicEvent(id, &tLogicEvent{typeId: keyUpType, valA: int(code), time: appTime.Millis()})
}

//export g2dMouseMove
func g2dMouseMove(id C.int) {
	postLogicEvent(id, &tLogicEvent{typeId: msMoveType, time: appTime.Millis()})
}

//export g2dWindowMove
func g2dWindowMove(id C.int) {
	postLogicEvent(id, &tLogicEvent{typeId: wndMoveType, time: appTime.Millis()})
}

//export g2dWindowResize
func g2dWindowResize(id C.int) {
	if !processingRequests {
		mutex.Lock()
	}
	wnd := wnds[id]
	event := &tLogicEvent{typeId: wndResizeType, time: appTime.Millis()}
	event.props.update(wnd.data, wnd.title)
	wnd.eventsChan <- event
	if !processingRequests {
		mutex.Unlock()
	}
}

//export g2dButtonDown
func g2dButtonDown(id, code, doubleClick C.int) {
	postLogicEvent(id, &tLogicEvent{typeId: buttonDownType, valA: int(code), repeated: uint(doubleClick), time: appTime.Millis()})
}

//export g2dButtonUp
func g2dButtonUp(id, code, doubleClick C.int) {
	postLogicEvent(id, &tLogicEvent{typeId: buttonUpType, valA: int(code), repeated: uint(doubleClick), time: appTime.Millis()})
}

//export g2dWheel
func g2dWheel(id C.int, wheel C.float) {
	postLogicEvent(id, &tLogicEvent{typeId: wheelType, valC: float32(wheel), time: appTime.Millis()})
}

//export g2dWindowMinimize
func g2dWindowMinimize(id C.int) {
	postLogicEvent(id, &tLogicEvent{typeId: minimizeType, time: appTime.Millis()})
}

//export g2dWindowRestore
func g2dWindowRestore(id C.int) {
	postLogicEvent(id, &tLogicEvent{typeId: restoreType, time: appTime.Millis()})
}

//export g2dOnFocus
func g2dOnFocus(id, focus C.int) {
	postLogicEvent(id, &tLogicEvent{typeId: focusType, valA: int(focus), time: appTime.Millis()})
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
		} else if err1 < 1002001 {
			if err1 < 1001007 {
				switch err1 {
				case 1001001:
					errStr = fmt.Sprintf(functionFailedWindow, "RegisterClassEx")
				case 1001002:
					errStr = fmt.Sprintf(functionFailedWindow, "CreateWindow")
				case 1001003:
					errStr = fmt.Sprintf(functionFailedWindow, "GetDC")
				case 1001004:
					errStr = fmt.Sprintf(functionFailedWindow, "wglChoosePixelFormatARB")
				case 1001005:
					errStr = fmt.Sprintf(functionFailedWindow, "SetPixelFormat")
				case 1001006:
					errStr = fmt.Sprintf(functionFailedWindow, "wglCreateContextAttribsARB")
				}
			} else if err1 < 1001015 {
				errStr = fmt.Sprintf(functionFailedWindow, "set fullscreen")
			} else {
				switch err1 {
				case 1001015:
					errStr = fmt.Sprintf(functionFailedWindow, "wglDeleteContext")
				case 1001016:
					errStr = fmt.Sprintf(functionFailedWindow, "DestroyWindow")
				case 1001017:
					errStr = fmt.Sprintf(functionFailedWindow, "UnregisterClass")
				case 1001018:
					errStr = fmt.Sprintf(functionFailedWindow, "set position")
				case 1001019:
					errStr = fmt.Sprintf(functionFailedWindow, "set position")
				case 1001020:
					errStr = fmt.Sprintf(functionFailedWindow, "move")
				case 1001021:
					errStr = fmt.Sprintf(functionFailedWindow, "set title")
				case 1001022:
					errStr = fmt.Sprintf(functionFailedG2D, "set mouse position")
				}
			}
		} else {
			switch err1 {
			case 1002001:
				errStr = fmt.Sprintf(functionFailedGraphics, "wglMakeCurrent")
			case 1002002:
				errStr = fmt.Sprintf(functionFailedG2D, "create vertex shader")
			case 1002003:
				errStr = fmt.Sprintf(functionFailedG2D, "create vertex shader")
			case 1002004:
				errStr = fmt.Sprintf(functionFailedG2D, "create fragment shader")
			case 1002005:
				errStr = fmt.Sprintf(functionFailedG2D, "create fragment shader")
			case 1002006:
				errStr = fmt.Sprintf(functionFailedG2D, "attach vertex shader")
			case 1002007:
				errStr = fmt.Sprintf(functionFailedG2D, "attach vertex shader")
			case 1002008:
				errStr = fmt.Sprintf(functionFailedG2D, "attach fragment shader")
			case 1002009:
				errStr = fmt.Sprintf(functionFailedG2D, "attach fragment shader")
			case 1002010:
				errStr = fmt.Sprintf(functionFailedG2D, "link shader program")
			case 1002011:
				errStr = fmt.Sprintf(functionFailedG2D, "create shader program")
			case 1002012:
				errStr = fmt.Sprintf(functionFailedG2D, "use shader program")
			case 1002013:
				errStr = fmt.Sprintf(functionFailedG2D, "use shader program")
			case 1002014:
				errStr = fmt.Sprintf(functionFailedG2D, "bind vertex array")
			case 1002015:
				errStr = fmt.Sprintf(functionFailedG2D, "bind buffer")
			case 1002016:
				errStr = fmt.Sprintf(functionFailedG2D, "bind buffer")
			case 1002017:
				errStr = fmt.Sprintf(functionFailedG2D, "draw rectangles")
			case 1002018:
				errStr = fmt.Sprintf(functionFailedG2D, "draw rectangles")
			case 1002019:
				errStr = fmt.Sprintf(functionFailedG2D, "draw rectangles")
			case 1002020:
				errStr = fmt.Sprintf(functionFailedG2D, "set buffer data")
			case 1002021:
				errStr = fmt.Sprintf(functionFailedG2D, "set buffer data")
			case 1002022:
				errStr = fmt.Sprintf(functionFailedG2D, "set buffer data")
			case 1002023:
				errStr = fmt.Sprintf(functionFailedG2D, "get attribute location")
			case 1002024:
				errStr = fmt.Sprintf(functionFailedG2D, "get attribute location")
			case 1002025:
				errStr = fmt.Sprintf(functionFailedG2D, "get uniform location")
			case 1002026:
				errStr = fmt.Sprintf(functionFailedG2D, "get uniform location")
			case 1002027:
				errStr = fmt.Sprintf(functionFailedG2D, "bind vertex array")
			case 1002028:
				errStr = fmt.Sprintf(functionFailedG2D, "enable attribute")
			case 1002029:
				errStr = fmt.Sprintf(functionFailedG2D, "enable attribute")
			case 1002030:
				errStr = fmt.Sprintf(functionFailedG2D, "enable attribute")
			case 1002031:
				errStr = fmt.Sprintf(functionFailedG2D, "enable attribute")
			case 1002032:
				errStr = fmt.Sprintf(functionFailedG2D, "bind buffer")
			case 1002033:
				errStr = fmt.Sprintf(functionFailedG2D, "bind buffer")
			case 1002034:
				errStr = fmt.Sprintf(functionFailedG2D, "set buffer data")
			case 1002035:
				errStr = fmt.Sprintf(functionFailedG2D, "set buffer data")
			case 1002036:
				errStr = fmt.Sprintf(functionFailedG2D, "set buffer data")
			case 1002037:
				errStr = fmt.Sprintf(functionFailedG2D, "set buffer data")
			case 1002038:
				errStr = fmt.Sprintf(functionFailedG2D, "set buffer data")
			case 1002039:
				errStr = fmt.Sprintf(functionFailedG2D, "set buffer data")
			case 1002040:
				errStr = fmt.Sprintf(functionFailedG2D, "set buffer data")
			case 1002041:
				errStr = fmt.Sprintf(functionFailedG2D, "set buffer data")
			case 1002042:
				errStr = fmt.Sprintf(functionFailedG2D, "set vertex data")
			case 1002043:
				errStr = fmt.Sprintf(functionFailedG2D, "set vertex data")
			case 1002044:
				errStr = fmt.Sprintf(functionFailedG2D, "set vertex data")
			case 1002045:
				errStr = fmt.Sprintf(functionFailedG2D, "set vertex data")
			case 1002046:
				errStr = fmt.Sprintf(functionFailedG2D, "set vertex data")
			case 1002047:
				errStr = fmt.Sprintf(functionFailedG2D, "set vertex data")
			case 1002048:
				errStr = fmt.Sprintf(functionFailedG2D, "bind buffer")
			case 1002049:
				errStr = fmt.Sprintf(functionFailedG2D, "bind buffer")
			case 1002050:
				errStr = fmt.Sprintf(functionFailedGraphics, "SwapBuffers")
			case 1002051:
				errStr = fmt.Sprintf(functionFailedGraphics, "wglMakeCurrent")
			case 1002052:
				errStr = fmt.Sprintf(functionFailedG2D, "load texture")
			case 1002053:
				errStr = fmt.Sprintf(functionFailedG2D, "load texture")
			case 1002054:
				errStr = fmt.Sprintf(functionFailedG2D, "load texture")
			case 1002055:
				errStr = fmt.Sprintf(functionFailedG2D, "load texture")
			case 1002056:
				errStr = fmt.Sprintf(functionFailedG2D, "load texture")
			case 1002057:
				errStr = fmt.Sprintf(functionFailedG2D, "load texture")
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
		for len(errStr) > 0 && (errStr[len(errStr)-1] == '\n' || errStr[len(errStr)-1] == '\r') {
			errStr = errStr[:len(errStr)-1]
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
