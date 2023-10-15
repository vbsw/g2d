/*
 *          Copyright 2023, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

package g2d

// #cgo CFLAGS: -DVBSW_G2D_WIN32 -DUNICODE
// #cgo LDFLAGS: -luser32 -lgdi32 -lOpenGL32
// #include "g2d.h"
import "C"
import (
	"errors"
	"fmt"
	"strconv"
	"runtime"
	"time"
	"unsafe"
)

func Init() {
	mutex.Lock()
	defer mutex.Unlock()
	if !initialized {
		props := engineProperties()
		var maxTexSizeC C.int
		var err1, err2 C.longlong
		C.g2d_init(&maxTexSizeC, &err1, &err2)
		if err1 == 0 {
			ErrConv = props.errConv
			MaxTextureSize = int(maxTexSizeC)
			initialized = true
			initFailed = false
			quitting = false
			startTime = time.Now()
		} else {
			initFailed = true
			Err = props.errConv.ToError(int64(err1), int64(err2), "")
		}
	} else {
		panic("g2d engine is already initialized")
	}
}

func Show(window ...Window) {
	if anyAvailable(window) {
		mutex.Lock()
		if !initFailed {
			if initialized {
				if running {
					if !quitting {
						for _, abst := range window {
							if abst != nil {
								wnd := newWindow(abst)
								go wnd.logicThread()
								wnd.wgt.msgs <- (&tLMessage{typeId: configType, nanos: deltaNanos()})
							}
						}
					}
				} else {
					for _, abst := range window {
						if abst != nil {
							wnd := newWindow(abst)
							wndsToStart = append(wndsToStart, wnd)
						}
					}
					if len(wndsToStart) > 0 {
						running = true
						mutex.Unlock()
						C.g2d_process_messages()
						mutex.Lock()
						running = false
					}
					cleanUp()
				}
				mutex.Unlock()
			} else {
				mutex.Unlock()
				panic("g2d is not initialized")
			}
		} else {
			mutex.Unlock()
		}
	}
}

func postMessage(msg interface{}, errInfo string) {
	var err1, err2 C.longlong
	mutex.Lock()
	defer mutex.Unlock()
	C.g2d_post_window_msg(&err1, &err2)
	if err1 == 0 {
		msgs.Put(msg)
	} else {
		setError(err1, err2, nil, errInfo)
	}
}

func cleanUp() {
	for _, wnd := range wndCbs {
		if wnd != nil {
			var err1, err2 C.longlong
			if wnd.wgt != nil {
				wgt := wnd.wgt
				wgt.msgs <- (&tLMessage{typeId: quitType, nanos: deltaNanos()})
				<- wgt.quitted
			}
			unregister(wnd.cbId)
			C.g2d_window_destroy(wnd.dataC, &err1, &err2)
			if err1 != 0 {
				setError(err1, err2, nil, "window " + wnd.cbIdStr)
			}
		}
	}
	msgs.Reset(-1)
	C.g2d_clean_up_messages()
}

func newWindow(abst Window) *tWindow {
	wnd := new(tWindow)
	wnd.state = configState
	wnd.abst = abst
	wnd.wgt = new(Widget)
	wnd.wgt.msgs = make(chan *tLMessage, 1000)
	wnd.wgt.quitted = make(chan bool, 1)
	wnd.cbId = register(wnd)
	wnd.cbIdStr = strconv.FormatInt(int64(wnd.cbId), 10)
	wnd.wgt.Gfx.msgs = make(chan *tGMessage, 1000)
	wnd.wgt.Gfx.quitted = make(chan bool, 1)
	wnd.wgt.Gfx.rBuffer = &wnd.wgt.Gfx.buffers[0]
	wnd.wgt.Gfx.wBuffer = &wnd.wgt.Gfx.buffers[0]
	wnd.wgt.Gfx.initEntities()
	return wnd
}

func (wnd *tWindow) logicThread() {
	for wnd.state != quitState {
		msg := wnd.nextLMessage()
		if msg != nil {
			wnd.wgt.CurrEventNanos = msg.nanos
			switch msg.typeId {
			case configType:
				wnd.onConfig()
			case createType:
				wnd.onCreate()
			case showType:
				wnd.onShow()
			case resizeType:
				wnd.updateProps(msg)
				wnd.onResize()
			case keyDownType:
				wnd.onKeyDown(msg.valA, msg.repeated)
			case keyUpType:
				wnd.onKeyUp(msg.valA)
/*
			case textureType:
				wnd.onTextureLoaded(msg.valA)
*/
			case updateType:
				wnd.onUpdate()
			case quitReqType:
				wnd.onQuitReq()
			case quitType:
				wnd.onQuit()
			}
		}
	}
	wnd.wgt.Gfx.rBuffer = nil
	wnd.wgt.Gfx.wBuffer = nil
/*
	wnd.wgt.Gfx.buffers[0].layers = nil
	wnd.wgt.Gfx.buffers[1].layers = nil
	wnd.wgt.Gfx.buffers[2].layers = nil
	wnd.wgt.Gfx.entitiesLayers = nil
*/
	wnd.wgt = nil
}

func (wnd *tWindow) updateProps(msg *tLMessage) {
	wnd.wgt.ClientX = msg.props.ClientX
	wnd.wgt.ClientY = msg.props.ClientY
	wnd.wgt.ClientWidth = msg.props.ClientWidth
	wnd.wgt.ClientHeight = msg.props.ClientHeight
	wnd.wgt.MouseX = msg.props.MouseX
	wnd.wgt.MouseY = msg.props.MouseY
}

func (wnd *tWindow) nextLMessage() *tLMessage {
	var message *tLMessage
	if wnd.state > configState && wnd.state < closingState && (wnd.autoUpdate || wnd.wgt.update) {
		select {
		case msg := <-wnd.wgt.msgs:
			message = msg
		default:
			wnd.wgt.update = false
			message = &tLMessage{typeId: updateType, nanos: deltaNanos()}
		}
	} else {
		message = <-wnd.wgt.msgs
	}
	if wnd.state == closingState && message.typeId != quitType {
		message = nil
	}
	return message
}

func (wnd *tWindow) onConfig() {
	config := newConfiguration()
	err := wnd.abst.OnConfig(config)
	wnd.autoUpdate = config.AutoUpdate
	if err == nil {
		errInfo := "create-request, window " + wnd.cbIdStr
		postMessage(&tCreateWindowRequest{window: wnd, config: config}, errInfo)
	} else {
		wnd.onLError(err)
	}
}

func (wnd *tWindow) onCreate() {
	wnd.wgt.PrevUpdateNanos = wnd.wgt.CurrEventNanos
	err := wnd.abst.OnCreate(wnd.wgt)
	if err == nil {
		wnd.state = runningState
		wnd.wgt.Gfx.running = true
		wnd.wgt.Gfx.msgs <- &tGMessage{typeId: refreshType}
		go wnd.graphicsThread()
		wnd.wgt.Gfx.switchWBuffer()
		errInfo := "show-request, window " + wnd.cbIdStr
		postMessage(&tShowWindowRequest{window: wnd}, errInfo)
	} else {
		wnd.onLError(err)
	}
}

func (wnd *tWindow) onShow() {
	wnd.wgt.PrevUpdateNanos = wnd.wgt.CurrEventNanos
	err := wnd.abst.OnShow()
	if err == nil {
		wnd.wgt.Gfx.switchWBuffer()
		wnd.wgt.Gfx.msgs <- &tGMessage{typeId: refreshType}
	} else {
		wnd.onLError(err)
	}
}

func (wnd *tWindow) onResize() {
	err := wnd.abst.OnResize()
	if err != nil {
		wnd.onLError(err)
	}
}

func (wnd *tWindow) onKeyDown(keyCode int, repeated uint) {
	err := wnd.abst.OnKeyDown(keyCode, repeated)
	if err != nil {
		wnd.onLError(err)
	}
}

func (wnd *tWindow) onKeyUp(keyCode int) {
	err := wnd.abst.OnKeyUp(keyCode)
	if err != nil {
		wnd.onLError(err)
	}
}

func (wnd *tWindow) onUpdate() {
	wnd.wgt.DeltaUpdateNanos = wnd.wgt.CurrEventNanos - wnd.wgt.PrevUpdateNanos
	err := wnd.abst.OnUpdate()
	wnd.wgt.PrevUpdateNanos = wnd.wgt.CurrEventNanos
	if err == nil {
		wnd.wgt.Gfx.switchWBuffer()
		wnd.wgt.Gfx.msgs <- &tGMessage{typeId: refreshType}
	} else {
		wnd.onLError(err)
	}
}

func (wnd *tWindow) onQuitReq() {
	closeOk, err := wnd.abst.OnClose()
	if err == nil {
		if closeOk {
			wnd.wgt.CurrEventNanos = deltaNanos()
			wnd.onQuit()
			errInfo := "destroy-request, window " + wnd.cbIdStr
			postMessage(&tDestroyWindowRequest{window: wnd}, errInfo)
		}
	} else {
		wnd.onLError(err)
	}
}

func (wnd *tWindow) onQuit() {
	if wnd.wgt.Gfx.running {
		wnd.wgt.Gfx.msgs <- &tGMessage{typeId: quitType}
		<- wnd.wgt.Gfx.quitted
	}
	err := wnd.abst.OnDestroy()
	wnd.wgt.quitted <- true
	wnd.state = quitState
	if err != nil {
		setErrorSynced(err)
	}
}

func (wnd *tWindow) onLError(err error) {
	setErrorSynced(err)
	wnd.wgt.CurrEventNanos = deltaNanos()
	wnd.onQuit()
}

func (wnd *tWindow) graphicsThread() {
	var err1, err2 C.longlong
	var errStrC *C.char
	runtime.LockOSThread()
	C.g2d_context_make_current(wnd.dataC, &err1, &err2)
	if err1 == 0 {
		C.g2d_gfx_init(wnd.dataC, C.int(wnd.wgt.Gfx.swapInterval), &err1, &err2, &errStrC)
		if err1 == 0 {
			for wnd.wgt.Gfx.running {
				msg := wnd.nextGMessage()
				if msg != nil {
					switch msg.typeId {
					case refreshType:
						wnd.drawGraphics()
					case swapIntervType:
						C.g2d_gfx_set_swap_interval(C.int(msg.valA))
					case resizeType:
						C.g2d_gfx_set_view_size(wnd.dataC, C.int(msg.valA), C.int(msg.valB))
/*
					case imageType:
						texBytes, ok := msg.valC.([]byte)
						if ok {
							wnd.loadTexture(texBytes, msg.valA, msg.valB)
						} else {
							appendError(msg.err)
							processing = wnd.processGMessage(&tGMessage{typeId: quitType})
						}
*/
					case quitType:
						var err1, err2 C.longlong
						C.g2d_context_release(wnd.dataC, &err1, &err2)
						if (err1 == 0) {
							wnd.wgt.Gfx.quitted <- true
							wnd.wgt.Gfx.running = false
						}
					}
				}
			}
		}
	}
	if (err1 != 0) {
		wnd.onGError(err1, err2, errStrC)
	}
}

func (wnd *tWindow) nextGMessage() *tGMessage {
	var message *tGMessage
	if wnd.wgt.Gfx.refresh {
		select {
		case msg := <-wnd.wgt.Gfx.msgs:
			if msg.typeId != refreshType {
				message = msg
			}
		default:
			wnd.wgt.Gfx.refresh = false
			message = &tGMessage{typeId: refreshType}
		}
	} else {
		message = <-wnd.wgt.Gfx.msgs
		if message.typeId == refreshType {
			wnd.wgt.Gfx.refresh = true
			message = nil
		}
	}
	return message
}

func (wnd *tWindow) drawGraphics() {
	var err1, err2 C.longlong
	wnd.wgt.Gfx.switchRBuffer()
	buffer := wnd.wgt.Gfx.rBuffer
	C.g2d_gfx_clear_bg(buffer.bgR, buffer.bgG, buffer.bgB)
/*
	for _, layer := range wnd.wgt.Gfx.rBuffer.layers {
		err := layer.draw(wnd.dataC)
		if err != nil {
			appendError(err)
			wnd.wgt.Gfx.msgs <- &tGMessage{typeId: quitType}
		}
	}
*/
	C.g2d_gfx_swap_buffers(wnd.dataC, &err1, &err2)
	if err1 != 0 {
		wnd.onGError(err1, err2, nil)
	}
}

func (wnd *tWindow) onGError(err1, err2 C.longlong, errStrC *C.char) {
	mutex.Lock()
	wnd.wgt.Gfx.quitted <- true
	wnd.wgt.Gfx.running = false
	setError(err1, err2, errStrC, "window " + wnd.cbIdStr)
	C.g2d_context_release(wnd.dataC, &err1, &err2)
	mutex.Unlock()
}

//export g2dStartWindows
func g2dStartWindows() {
	mutex.Lock()
	defer mutex.Unlock()
	for _, wnd := range wndsToStart {
		go wnd.logicThread()
		wnd.wgt.msgs <- (&tLMessage{typeId: configType, nanos: deltaNanos()})
	}
	wndsToStart = wndsToStart[:0]
}

//export g2dProcessMessage
func g2dProcessMessage() {
	mutex.Lock()
	defer mutex.Unlock()
	if !quitting {
		message := msgs.First()
		if message != nil {
			switch msg := message.(type) {
			case *tCreateWindowRequest:
				createWindow(msg.window, msg.config)
			case *tShowWindowRequest:
				showWindow(msg.window)
			case *tDestroyWindowRequest:
				destroyWindow(msg.window)
			}
		}
	}
}

//export g2dResize
func g2dResize(cbIdC C.int) {
	wnd := wndCbs[int(cbIdC)]
	wgt := wnd.wgt
	if wgt != nil {
		msg := &tLMessage{typeId: resizeType, nanos: deltaNanos()}
		msg.props.update(wnd.dataC)
		wgt.msgs <- msg
		wnd.wgt.Gfx.msgs <- &tGMessage{typeId: resizeType, valA: msg.props.ClientWidth, valB: msg.props.ClientHeight}
	}
}

//export g2dKeyDown
func g2dKeyDown(cbIdC, code C.int, repeated C.int) {
	wnd := wndCbs[int(cbIdC)]
	wgt := wnd.wgt
	if wgt != nil {
		msg := &tLMessage{typeId: keyDownType, valA: int(code), repeated: uint(repeated), nanos: deltaNanos()}
		msg.props.update(wnd.dataC)
		wgt.msgs <- msg
	}
}

//export g2dKeyUp
func g2dKeyUp(cbIdC, code C.int) {
	wnd := wndCbs[int(cbIdC)]
	wgt := wnd.wgt
	if wgt != nil {
		msg := &tLMessage{typeId: keyUpType, valA: int(code), nanos: deltaNanos()}
		msg.props.update(wnd.dataC)
		wgt.msgs <- msg
	}
}

//export g2dClose
func g2dClose(cbIdC C.int) {
	wnd := wndCbs[int(cbIdC)]
	wgt := wnd.wgt
	if wgt != nil {
		wgt.msgs <- (&tLMessage{typeId: quitReqType, nanos: deltaNanos()})
	}
}

func createWindow(wnd *tWindow, config *Configuration) {
	var err1, err2 C.longlong
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
	t, err1 := toTString(config.Title)
	if err1 == 0 {
		C.g2d_window_create(&wnd.dataC, C.int(wnd.cbId), x, y, w, h, wn, hn, wx, hx, b, d, r, f, l, c, t, &err1, &err2)
		C.g2d_free(t)
		if err1 == 0 {
			msg := &tLMessage{typeId: createType, nanos: deltaNanos()}
			msg.props.update(wnd.dataC)
			wnd.wgt.msgs <- msg
		} else {
			setError(err1, err2, nil, "window " + wnd.cbIdStr)
		}
	} else {
		setError(err1, err2, nil, "window " + wnd.cbIdStr)
	}
}

func showWindow(wnd *tWindow) {
	var err1, err2 C.longlong
	C.g2d_window_show(wnd.dataC, &err1, &err2)
	if err1 == 0 {
		wnd.wgt.msgs <- (&tLMessage{typeId: showType, nanos: deltaNanos()})
	} else {
		setError(err1, err2, nil, "window " + wnd.cbIdStr)
	}
}

func destroyWindow(wnd *tWindow) {
	var err1, err2 C.longlong
	C.g2d_window_destroy(wnd.dataC, &err1, &err2)
	unregister(wnd.cbId)
	wnd.cbId = -1
	if err1 != 0 {
		setError(err1, err2, nil, "window " + wnd.cbIdStr)
	}
}

func setError(err1, err2 C.longlong, errStrC *C.char, info string) {
	if Err == nil {
		if errStrC != nil {
			info = C.GoString(errStrC)
			C.g2d_free(unsafe.Pointer(errStrC))
		}
		Err = ErrConv.ToError(int64(err1), int64(err1), info)
	}
	if running && !quitting {
		var err1, err2 C.longlong
		C.g2d_post_quit_msg(&err1, &err2)
		quitting = true
	}
}

func setErrorSynced(err error) {
	mutex.Lock()
	defer mutex.Unlock()
	if Err == nil {
		Err = err
	}
	if running && !quitting {
		var err1, err2 C.longlong
		C.g2d_post_quit_msg(&err1, &err2)
		quitting = true
	}
}

func (props *Properties) update(dataC unsafe.Pointer) {
	var x, y, w, h, wn, hn, wx, hx, b, d, r, f, l C.int
	C.g2d_window_props(dataC, &x, &y, &w, &h, &wn, &hn, &wx, &hx, &b, &d, &r, &f, &l)
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

func (errConv *tErrorConvertor) ToError(err1, err2 int64, info string) error {
	var errStr string
	if err1 > 0 && err1 < 1000 {
		errStr = "memory allocation failed"
	} else {
		errStr = "unknown error"
	}
	errStr = errStr + " (" + strconv.FormatInt(err1, 10)
	if err2 == 0 {
		errStr = errStr + ")"
	} else {
		errStr = errStr + ", " + strconv.FormatInt(err2, 10) + ")"
	}
	if len(info) > 0 {
		errStr = errStr + "; " + info
	}
	return errors.New(errStr)
}

func toTString(str string) (unsafe.Pointer, C.longlong) {
	var strT unsafe.Pointer
	var err1 C.longlong
	if len(str) > 0 {
		bytes := *(*[]byte)(unsafe.Pointer(&str))
		C.g2d_to_tstr(&strT, unsafe.Pointer(&bytes[0]), C.size_t(len(str)), &err1)
	} else {
		C.g2d_to_tstr(&strT, nil, C.size_t(len(str)), &err1)
	}
	return strT, err1
}

//export goDebug
func goDebug(a, b C.int, c, d C.longlong) {
	fmt.Println(a, b, c, d)
}

/*
//export goDebugMessage
func goDebugMessage(code C.longlong, strC C.g2d_lpcstr) {
	fmt.Println("Msg:", C.GoString(strC), code)
}

import (
	"github.com/vbsw/golib/cdata"
	"github.com/vbsw/g2d/ogfl"
	"github.com/vbsw/g2d/modules"
	"unsafe"
)

func initA() {
	var collection cdata.Collection
	var loader ogfl.Loader
	var rects modules.Rectangles
	collection.Passes = 2
	collection.Init(&loader, &rects)
}

import (
	"github.com/vbsw/g2d/win32"
)

type ErrorConvertor interface {
	ToError(g2dErrNum, win32ErrNum uint64, info string) error
}

type ErrorLogger interface {
	LogError(err error)
}

type EngineParams struct {
	ErrConv   ErrorConvertor
	ErrLogger ErrorLogger
	Modules   []Module
}

type Module interface {
	DataFunc() (*unsafe.Pointer, unsafe.Pointer)
}

func (engine *Engine) Errors() []error {
	engine.mutex.Lock()
	defer engine.mutex.Unlock()
	return engine.errs
}

func (engine *Engine) setErrConv(errConv ErrorConvertor) {
	if errConv == nil {
		engine.errConv = new(tDefaultErrorConvertor)
	} else {
		engine.errConv = errConv
	}
}

func (engine *Engine) setErrLogger(errLogger ErrorLogger) {
	if errLogger == nil {
		engine.errLogger = new(tDefaultErrorLogger)
	} else {
		engine.errLogger = errLogger
	}
}

func (engine *Engine) appendError(err error) {
	engine.mutex.Lock()
	engine.errs = append(engine.errs, err)
	engine.errLogger.LogError(err)
	engine.mutex.Unlock()
}

type tDefaultErrorConvertor struct {
}

type tDefaultErrorLogger struct {
}

func (errLoger *tDefaultErrorLogger) LogError(err error) {
}

func newInitParams(engineDataC *unsafe.Pointer, engineParams *EngineParams) *win32.InitParams {
	initParams := new(win32.InitParams)
	length := len(engineParams.Modules)
	if length > 0 {
		initParams.Data = make([]*unsafe.Pointer, length)
		initParams.Funcs = make([]unsafe.Pointer, length)
		for i, m := range engineParams.Modules {
			initParams.Data[i], initParams.Funcs[i] = m.DataFunc()
		}
	}
	initParams.Engine = engineDataC
	return initParams
}
*/

/*
const (
	errStrPostMessage = "post message failed"
	errStrMalloc      = "memory allocation failed"
	errStrGetProcAddr = "load %s function failed"
)

func Init(params ...interface{}) {
	if !initialized {
		var errNumC C.int
		var errWin32C C.g2d_ul_t
		clearErrors()
		initCustParams(params)
		initDefaultParams()
		C.g2d_init(&maxTexSize, &errNumC, &errWin32C)
		if errNumC == 0 {
			startTime = time.Now()
			fsm = [56]int{0, 1, 2, 0, 10, 2, 2, 1, 3, 1, 2, 0, 3, 4, 1, 0, 6, 5, 1, 2, 0, 4, 1, 0, 6, 7, 0, 2, 13, 8, 0, 1, 9, 7, 0, 2, 9, 5, 1, 2, 10, 11, 0, 1, 9, 12, 0, 2, 13, 11, 0, 1, 13, 2, 2, 1}
			initialized = true
			initFailed = false
		} else {
			initFailed = true
			appendError(toError(errNumC, errWin32C, nil))
		}
	} else {
		panic("g2d is already initialized")
	}
}

func (window *tWindow) onTextureLoaded(textureId int) {
	err := window.abst.OnTextureLoaded(textureId)
	if err != nil {
		window.onLError(err)
	}
}

func (window *tWindow) loadTexture(bytes []byte, w, h int) {
	var errNumC, texIdC C.int
	C.g2d_gfx_gen_tex(window.dataC, unsafe.Pointer(&bytes[0]), C.int(w), C.int(h), &texIdC, &errNumC)
	if errNumC == 0 {
		window.wgt.msgs <- &tLMessage{typeId: textureType, valA: int(texIdC), nanos: deltaNanos()}
	} else {
		window.state = 2
		appendError(toError(errNumC, 0, nil))
		window.wgt.Gfx.msgs <- &tGMessage{typeId: quitType}
	}
}

func (layer *tRectLayer) draw(dataC unsafe.Pointer) error {
	var errNumC C.int
	var errStrC *C.char
	length := len(layer.enabled)
	if length > 0 {
		C.g2d_gfx_draw_rect(dataC, &layer.enabled[0], &layer.rects[0], C.int(length), C.int(layer.totalActive), &errNumC, &errStrC)
		if errNumC != 0 {
			return toError(errNumC, 0, errStrC)
		}
	}
	return nil
}

func (layer *tImageLayer) draw(dataC unsafe.Pointer) error {
	var errNumC C.int
	var errStrC *C.char
	length := len(layer.enabled)
	if length > 0 {
		C.g2d_gfx_draw_image(dataC, &layer.enabled[0], &layer.rects[0], C.int(length), C.int(layer.totalActive), C.int(layer.textureId), &errNumC, &errStrC)
		if errNumC != 0 {
			return toError(errNumC, 0, errStrC)
		}
	}
	return nil
}

func initDefaultParams() {
	if errGen == nil {
		errGen = &errHandler
	}
	if errLog == nil {
		errLog = &errHandler
	}
}

func (_ *tErrorHandler) ToError(g2dErrNum, win32ErrNum uint64, info string) error {
	var errStr string
	switch g2dErrNum {
	case 1:
		errStr = "get module instance failed"
	case 2:
		errStr = "register dummy class failed"
	case 3:
		errStr = "create dummy window failed"
	case 4:
		errStr = "get dummy device context failed"
	case 5:
		errStr = "choose dummy pixel format failed"
	case 6:
		errStr = "set dummy pixel format failed"
	case 7:
		errStr = "create dummy render context failed"
	case 8:
		errStr = "make dummy context current failed"
	case 9:
		errStr = "release dummy context failed"
	case 10:
		errStr = "deleting dummy render context failed"
	case 11:
		errStr = "destroying dummy window failed"
	case 12:
		errStr = "unregister dummy class failed"

	case 13:
		errStr = "register class failed"
	case 14:
		errStr = "create window failed"
	case 15:
		errStr = "get device context failed"
	case 16:
		errStr = "choose pixel format failed"
	case 17:
		errStr = "set pixel format failed"
	case 18:
		errStr = "create render context failed"
	case 19:
		errStr = "release context failed"
	case 20:
		errStr = "deleting render context failed"
	case 21:
		errStr = "destroying window failed"
	case 22:
		errStr = "unregister class failed"
	case 23:
		errStr = "show window failed; type Window is not embedded"

	case 56:
		errStr = "make context current failed"
	case 61:
		errStr = "swap buffer failed"
	case 62:
		errStr = "set title failed"
	case 63:
		errStr = "wgl functions not initialized"
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
		errStr = errStrPostMessage
	case 81:
		errStr = errStrPostMessage
	case 82:
		errStr = errStrPostMessage
	case 83:
		errStr = errStrPostMessage

	case 100:
		errStr = "not initialized"
	case 101:
		errStr = "not initialized"
	case 102:
		errStr = "not initialized"
	case 120:
		errStr = errStrMalloc
	case 121:
		errStr = errStrMalloc

	case 200:
		errStr = fmt.Sprintf(errStrGetProcAddr, "wglChoosePixelFormatARB")
	case 201:
		errStr = fmt.Sprintf(errStrGetProcAddr, "wglCreateContextAttribsARB")
	case 202:
		errStr = fmt.Sprintf(errStrGetProcAddr, "wglSwapIntervalEXT")
	case 203:
		errStr = fmt.Sprintf(errStrGetProcAddr, "wglGetSwapIntervalEXT")
	case 204:
		errStr = fmt.Sprintf(errStrGetProcAddr, "glCreateShader")
	case 205:
		errStr = fmt.Sprintf(errStrGetProcAddr, "glShaderSource")
	case 206:
		errStr = fmt.Sprintf(errStrGetProcAddr, "glCompileShader")
	case 207:
		errStr = fmt.Sprintf(errStrGetProcAddr, "glGetShaderiv")
	case 208:
		errStr = fmt.Sprintf(errStrGetProcAddr, "glGetShaderInfoLog")
	case 209:
		errStr = fmt.Sprintf(errStrGetProcAddr, "glCreateProgram")
	case 210:
		errStr = fmt.Sprintf(errStrGetProcAddr, "glAttachShader")
	case 211:
		errStr = fmt.Sprintf(errStrGetProcAddr, "glLinkProgram")
	case 212:
		errStr = fmt.Sprintf(errStrGetProcAddr, "glValidateProgram")
	case 213:
		errStr = fmt.Sprintf(errStrGetProcAddr, "glGetProgramiv")
	case 214:
		errStr = fmt.Sprintf(errStrGetProcAddr, "glGetProgramInfoLog")
	case 215:
		errStr = fmt.Sprintf(errStrGetProcAddr, "glGenBuffers")
	case 216:
		errStr = fmt.Sprintf(errStrGetProcAddr, "glGenVertexArrays")
	case 217:
		errStr = fmt.Sprintf(errStrGetProcAddr, "glGetAttribLocation")
	case 218:
		errStr = fmt.Sprintf(errStrGetProcAddr, "glBindVertexArray")
	case 219:
		errStr = fmt.Sprintf(errStrGetProcAddr, "glEnableVertexAttribArray")
	case 220:
		errStr = fmt.Sprintf(errStrGetProcAddr, "glVertexAttribPointer")
	case 221:
		errStr = fmt.Sprintf(errStrGetProcAddr, "glBindBuffer")
	case 222:
		errStr = fmt.Sprintf(errStrGetProcAddr, "glBufferData")
	case 223:
		errStr = fmt.Sprintf(errStrGetProcAddr, "glGetVertexAttribPointerv")
	case 224:
		errStr = fmt.Sprintf(errStrGetProcAddr, "glUseProgram")
	case 225:
		errStr = fmt.Sprintf(errStrGetProcAddr, "glDeleteVertexArrays")
	case 226:
		errStr = fmt.Sprintf(errStrGetProcAddr, "glDeleteBuffers")
	case 227:
		errStr = fmt.Sprintf(errStrGetProcAddr, "glDeleteProgram")
	case 228:
		errStr = fmt.Sprintf(errStrGetProcAddr, "glDeleteShader")
	case 229:
		errStr = fmt.Sprintf(errStrGetProcAddr, "glGetUniformLocation")
	case 230:
		errStr = fmt.Sprintf(errStrGetProcAddr, "glUniformMatrix3fv")
	case 231:
		errStr = fmt.Sprintf(errStrGetProcAddr, "glUniform1fv")
	case 232:
		errStr = fmt.Sprintf(errStrGetProcAddr, "glUniformMatrix4fv")
	case 233:
		errStr = fmt.Sprintf(errStrGetProcAddr, "glUniformMatrix2x3fv")
	case 234:
		errStr = fmt.Sprintf(errStrGetProcAddr, "glGenerateMipmap")
	case 235:
		errStr = fmt.Sprintf(errStrGetProcAddr, "glActiveTexture")
	default:
		errStr = "unknown error"
	}
	errStr = errStr + " (" + strconv.FormatUint(g2dErrNum, 10)
	if win32ErrNum == 0 {
		errStr = errStr + ")"
	} else {
		errStr = errStr + ", " + strconv.FormatUint(win32ErrNum, 10) + ")"
	}
	if len(info) > 0 {
		errStr = errStr + "; " + info
	}
	return errors.New(errStr)
}

func (_ *tErrorHandler) LogError(err error) {
}

func getType(myvar interface{}) string {
	if t := reflect.TypeOf(myvar); t.Kind() == reflect.Ptr {
		return "*" + t.Elem().Name()
	} else {
		return t.Name()
	}
}
*/
