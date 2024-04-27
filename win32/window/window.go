/*
 *          Copyright 2024, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

// Package window creates a window with OpenGL 3.0 context.
package window

// #cgo CFLAGS: -DUNICODE
// #cgo LDFLAGS: -luser32 -lgdi32
// #include "window.h"
import "C"
import (
	"github.com/vbsw/g2d/win32/mainloop"
	"github.com/vbsw/golib/queue"
	"errors"
	"strconv"
	"sync"
	"time"
	"unsafe"
)

const (
	configState = iota
	runningState
	closingState
	quitState
)

const (
	configType = iota
	createType
	showType
	resizeType
	keyDownType
	keyUpType
	updateType
	quitReqType
	quitType
	leaveType
	refreshType
	swapIntervType
	imageType
	textureType
)

var (
	data tCData
	conv tErrorConv
	custErrConv cdata.ErrorConv
	mutex       sync.Mutex
	initialized bool
	running bool
	quitting bool
	Err error
	startTime   time.Time
	wndCbs     []*tWindow
	wndCbsNext []int
	msgs       queue.Queue
)

type Window interface {
	OnConfig(config *Configuration) error
	OnCreate(widget *Widget) error
	OnShow() error
/*
	OnResize() error
	OnKeyDown(keyCode int, repeated uint) error
	OnKeyUp(keyCode int) error
	OnTextureLoaded(textureId int) error
	OnUpdate() error
	OnClose() (bool, error)
	OnDestroy() error
*/
}

type Configuration struct {
	ClientX, ClientY                  int
	ClientWidth, ClientHeight         int
	ClientWidthMin, ClientHeightMin   int
	ClientWidthMax, ClientHeightMax   int
	MouseLocked, Borderless, Dragable bool
	Resizable, Fullscreen, Centered   bool
	AutoUpdate                        bool
	Title                             string
}

type Properties struct {
	MouseX, MouseY                    int
	ClientX, ClientY                  int
	ClientWidth, ClientHeight         int
	ClientWidthMin, ClientHeightMin   int
	ClientWidthMax, ClientHeightMax   int
	MouseLocked, Borderless, Dragable bool
	Resizable, Fullscreen             bool
	Title                             string
}

type Widget struct {
	ClientX, ClientY          int
	ClientWidth, ClientHeight int
	MouseX, MouseY            int
	NanosPrev           int64
	NanosCurr            int64
	NanosDelta          int64
	update                    bool
/*
	Gfx                       Graphics
*/
	msgs                      chan *tLogicMessage
	quitted                   chan bool
}

// tCData is the initialization for C.
type tCData struct {
}

// tErrorConv is the error converter.
type tErrorConv struct {
}

type tWindow struct {
	state      int
	cbId       int
	cbIdStr string
	abst       Window
	wgt        *Widget
	dataC      unsafe.Pointer
	autoUpdate bool
}

type tCreateWindowRequest struct {
	window *tWindow
	config *Configuration
}

type tShowWindowRequest struct {
	window *tWindow
}

type tDestroyWindowRequest struct {
	window *tWindow
}

// tLogicMessage communicates commands for logic handling.
type tLogicMessage struct {
	typeId   int
	valA     int
	repeated uint
	nanos    int64
	props    Properties
	obj      interface{}
}

// tGraphicsMessage communicates commands for graphics handling.
type tGraphicsMessage struct {
	typeId int
	valA   int
	valB   int
	valC   interface{}
	err    error
}

// CData returns an object to initialize this package.
// In cdata.Init first pass (pass = 0) initializes window configuration.
func CData() cdata.CData {
	return &data
}

// ErrorConv returns an instance of error convertor.
func ErrorConv() cdata.ErrorConv {
	if custErrConv == nil {
		return &conv
	}
	return custErrConv
}

// SetErrorConv overrides the default error convertor.
func SetErrorConv(errConv cdata.ErrorConv) {
	custErrConv = errConv
}

// MainLoop starts event processing. This function does not quit until an error occurs or all windows has been closed.
// (MainLoop needs at least one window to start.)
func MainLoop(windows []Window) error {
	mutex.Lock()
	if initialized {
		if !running {
			for i:=0;i< len(windows) &&!running; i++ {
				running = windows[i] != nil
			}
			if running {
				Err = nil
				startTime = time.Now()
				for _, window := range windows {
					if window != nil {
						launchWindow(window)
					}
				}
				mutex.Unlock()
				mainloop.Run()
				mutex.Lock()
				running = false
				cleanUp()
			}
			mutex.Unlock()
		} else {
			mutex.Unlock()
			panic("g2d.window MainLoop already running")
		}
	} else {
		mutex.Unlock()
		panic("g2d.window not initialized")
	}
	return nil
}

// Show will show windows, if MainLoop is running.
func Show(windows []Window) {
	mutex.Lock()
	if !quitting {
		if initialized {
			if running {
				for _, window := range windows {
					if window != nil {
						launchWindow(window)
					}
				}
				mutex.Unlock()
			} else {
				mutex.Unlock()
				panic("g2d.window MainLoop not running")
			}
		} else {
			mutex.Unlock()
			panic("g2d.window not initialized")
		}
	} else {
		mutex.Unlock()
	}
}

func launchWindow(window Window) {
	wnd := new(tWindow)
	wnd.state = configState
	wnd.abst = window
	wnd.wgt = new(Widget)
	wnd.wgt.msgs = make(chan *tLogicMessage, 1000)
	wnd.wgt.quitted = make(chan bool, 1)
	wnd.cbId = register(wnd)
	wnd.cbIdStr = strconv.FormatInt(int64(wnd.cbId), 10)
/*
	wnd.wgt.Gfx.msgs = make(chan *tGMessage, 1000)
	wnd.wgt.Gfx.quitted = make(chan bool, 1)
	wnd.wgt.Gfx.rBuffer = &wnd.wgt.Gfx.buffers[0]
	wnd.wgt.Gfx.wBuffer = &wnd.wgt.Gfx.buffers[0]
	wnd.wgt.Gfx.initEntities()
*/
	go wnd.logicThread()
	wnd.wgt.msgs <- (&tLogicMessage{typeId: configType, nanos: nanosDelta()})
}

func (wnd *tWindow) logicThread() {
	for wnd.state != quitState {
		msg := wnd.nextLogicMessage()
		if msg != nil {
			wnd.wgt.NanosCurr = msg.nanos
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
/*
	wnd.wgt.Gfx.rBuffer = nil
	wnd.wgt.Gfx.wBuffer = nil
	wnd.wgt.Gfx.buffers[0].layers = nil
	wnd.wgt.Gfx.buffers[1].layers = nil
	wnd.wgt.Gfx.buffers[2].layers = nil
	wnd.wgt.Gfx.entitiesLayers = nil
*/
	wnd.wgt = nil
}

func (wnd *tWindow) nextLogicMessage() *tLogicMessage {
	var message *tLogicMessage
	if wnd.state > configState && wnd.state < closingState && (wnd.autoUpdate || wnd.wgt.update) {
		select {
		case msg := <-wnd.wgt.msgs:
			message = msg
		default:
			wnd.wgt.update = false
			message = &tLogicMessage{typeId: updateType, nanos: nanosDelta()}
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
		wnd.onLogicError(err)
	}
}

func (wnd *tWindow) onCreate() {
	wnd.wgt.NanosPrev = wnd.wgt.NanosCurr
	err := wnd.abst.OnCreate(wnd.wgt)
	if err == nil {
		wnd.state = runningState
/*
		wnd.wgt.Gfx.running = true
		wnd.wgt.Gfx.msgs <- &tGMessage{typeId: refreshType}
		go wnd.graphicsThread()
		wnd.wgt.Gfx.switchWBuffer()
*/
		errInfo := "show-request, window " + wnd.cbIdStr
		postMessage(&tShowWindowRequest{window: wnd}, errInfo)
	} else {
		wnd.onLogicError(err)
	}
}

func (wnd *tWindow) onShow() {
	wnd.wgt.NanosPrev = wnd.wgt.NanosCurr
	err := wnd.abst.OnShow()
	if err == nil {
/*
		wnd.wgt.Gfx.switchWBuffer()
		wnd.wgt.Gfx.msgs <- &tGMessage{typeId: refreshType}
*/
	} else {
		wnd.onLogicError(err)
	}
}

func (wnd *tWindow) onResize() {
	err := wnd.abst.OnResize()
	if err != nil {
		wnd.onLogicError(err)
	}
}

func (wnd *tWindow) onKeyDown(keyCode int, repeated uint) {
	err := wnd.abst.OnKeyDown(keyCode, repeated)
	if err != nil {
		wnd.onLogicError(err)
	}
}

func (wnd *tWindow) onKeyUp(keyCode int) {
	err := wnd.abst.OnKeyUp(keyCode)
	if err != nil {
		wnd.onLogicError(err)
	}
}

func (wnd *tWindow) onUpdate() {
	wnd.wgt.DeltaUpdateNanos = wnd.wgt.NanosCurr - wnd.wgt.NanosPrev
	err := wnd.abst.OnUpdate()
	wnd.wgt.NanosPrev = wnd.wgt.NanosCurr
	if err == nil {
		wnd.wgt.Gfx.switchWBuffer()
		wnd.wgt.Gfx.msgs <- &tGMessage{typeId: refreshType}
	} else {
		wnd.onLogicError(err)
	}
}

func (wnd *tWindow) onQuitReq() {
	closeOk, err := wnd.abst.OnClose()
	if err == nil {
		if closeOk {
			wnd.wgt.NanosCurr = deltaNanos()
			wnd.onQuit()
			errInfo := "destroy-request, window " + wnd.cbIdStr
			postMessage(&tDestroyWindowRequest{window: wnd}, errInfo)
		}
	} else {
		wnd.onLogicError(err)
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

func (wnd *tWindow) onLogicError(err error) {
	setErrorSynced(err)
	wnd.wgt.NanosCurr = nanosDelta()
	wnd.onQuit()
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

func postMessage(event mainloop.Event, errInfo string) {
	err2 := mainloop.Post(event)
	if err2 != 0 {
		mutex.Lock()
		if Err == nil {
			Err = ErrorConv().ToError(12345678, err2, errInfo)
		}
		if running && !quitting {
			mainloop.Quit()
			quitting = true
		}
		mutex.Unlock()
	}
}

func setError(err1, err2 C.longlong, errStrC *C.char, info string) {
	if Err == nil {
		if errStrC != nil {
			info = C.GoString(errStrC)
			C.g2d_free(unsafe.Pointer(errStrC))
		}
		Err = ErrorConv().ToError(int64(err1), int64(err1), info)
	}
	if running && !quitting {
		var err1, err2 C.longlong
		C.g2d_post_quit_msg(&err1, &err2)
		quitting = true
	}
}

func setErrorSynced(err error) {
	mutex.Lock()
	if Err == nil {
		Err = err
	}
	if running && !quitting {
		var err1, err2 C.longlong
		C.g2d_post_quit_msg(&err1, &err2)
		quitting = true
	}
	mutex.Unlock()
}

func register(wnd *tWindow) int {
	var cbId int
	if len(wndCbsNext) == 0 {
		wndCbs = append(wndCbs, wnd)
		cbId = len(wndCbs) - 1
	} else {
		indexLast := len(wndCbsNext) - 1
		cbId = wndCbsNext[indexLast]
		wndCbsNext = wndCbsNext[:indexLast]
		wndCbs[cbId] = wnd
	}
	return cbId
}

func unregister(cbId int) {
	wndCbs[cbId] = nil
	wndCbsNext = append(wndCbsNext, cbId)
}

func cleanUp() {
	for _, wnd := range wndCbs {
		if wnd != nil {
			var err1, err2 C.longlong
			if wnd.wgt != nil {
				wgt := wnd.wgt
				wgt.msgs <- (&tLogicMessage{typeId: quitType, nanos: nanosDelta()})
				<- wgt.quitted
			}
			unregister(wnd.cbId)
			C.g2d_window_destroy(wnd.dataC, &err1, &err2)
			if err1 != 0 {
				setError(err1, err2, nil, "window " + wnd.cbIdStr)
			}
		}
	}
}

func newConfiguration() *Configuration {
	config := new(Configuration)
	config.ClientX = 50
	config.ClientY = 50
	config.ClientWidth = 640
	config.ClientHeight = 480
	config.ClientWidthMin = 0
	config.ClientHeightMin = 0
	config.ClientWidthMax = 99999
	config.ClientHeightMax = 99999
	config.MouseLocked = false
	config.Borderless = false
	config.Dragable = false
	config.Resizable = true
	config.Fullscreen = false
	config.Centered = true
	config.AutoUpdate = true
	config.Title = "g2d - 0.1.0"
	return config
}

// CInitFunc returns a function to initialize C data.
func (*tCData) CInitFunc() unsafe.Pointer {
	return C.g2d_window_init
}

// SetCData sets initialized C data.
func (*tCData) SetCData(data unsafe.Pointer) {
	mutex.Lock()
	initialized = true
	mutex.Unlock()
}

//export g2dProcessMessage
func g2dProcessMessage() {
	mutex.Lock()
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
	mutex.Unlock()
}

//export g2dResize
func g2dResize(cbIdC C.int) {
	wnd := wndCbs[int(cbIdC)]
	wgt := wnd.wgt
	if wgt != nil {
		msg := &tLogicMessage{typeId: resizeType, nanos: deltaNanos()}
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
		msg := &tLogicMessage{typeId: keyDownType, valA: int(code), repeated: uint(repeated), nanos: deltaNanos()}
		msg.props.update(wnd.dataC)
		wgt.msgs <- msg
	}
}

//export g2dKeyUp
func g2dKeyUp(cbIdC, code C.int) {
	wnd := wndCbs[int(cbIdC)]
	wgt := wnd.wgt
	if wgt != nil {
		msg := &tLogicMessage{typeId: keyUpType, valA: int(code), nanos: deltaNanos()}
		msg.props.update(wnd.dataC)
		wgt.msgs <- msg
	}
}

//export g2dClose
func g2dClose(cbIdC C.int) {
	wnd := wndCbs[int(cbIdC)]
	wgt := wnd.wgt
	if wgt != nil {
		wgt.msgs <- (&tLogicMessage{typeId: quitReqType, nanos: deltaNanos()})
	}
}

func (msg *tCreateWindowRequest) OnEvent() {
	var err1, err2 C.longlong
	config := msg.config
	wnd := msg.window
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
		C.g2d_window_free(t)
		if err1 == 0 {
			msg := &tLogicMessage{typeId: createType, nanos: nanosDelta()}
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
		wnd.wgt.msgs <- (&tLogicMessage{typeId: showType, nanos: deltaNanos()})
	} else {
		setError(err1, err2, nil, "window " + wnd.cbIdStr)
	}
}

func toTString(str string) (unsafe.Pointer, C.longlong) {
	var strT unsafe.Pointer
	var err1 C.longlong
	if len(str) > 0 {
		bytes := *(*[]byte)(unsafe.Pointer(&str))
		C.g2d_window_to_tstr(&strT, unsafe.Pointer(&bytes[0]), C.size_t(len(str)), &err1)
	} else {
		C.g2d_window_to_tstr(&strT, nil, C.size_t(len(str)), &err1)
	}
	return strT, err1
}

func nanosDelta() int64 {
	timeNow := time.Now()
	d := timeNow.Sub(startTime)
	return d.Nanoseconds()
}

func toCInt(b bool) C.int {
	if b {
		return C.int(1)
	}
	return C.int(0)
}

// ToError returns error numbers/string as error.
func (errConv *tErrorConv) ToError(err1, err2 int64, info string) error {
	if err1 >= 1000200 && err1 < 1000300 {
		var errStr string
		if err1 == 1000200 {
			errStr = "vbsw.g2d.window requires vbsw.g2d.loader"
		} else if err1 == 1000201 {
			errStr = "vbsw.g2d.window load wglChoosePixelFormatARB failed"
		} else if err1 == 1000202 {
			errStr = "vbsw.g2d.window load wglCreateContextAttribsARB failed"
		} else if err1 == 1000203 {
			errStr = "vbsw.g2d.window load wglSwapIntervalEXT failed"
		} else if err1 == 1000204 {
			errStr = "vbsw.g2d.window load wglGetSwapIntervalEXT failed"
		} else if err1 == 1000205 {
			errStr = "vbsw.g2d.window GetModuleHandle failed"
		} else if err1 == 1000206 {
			errStr = "vbsw.g2d.window RegisterClassEx failed"
		} else if err1 == 1000207 {
			errStr = "vbsw.g2d.window CreateWindow failed"
		} else if err1 == 1000208 {
			errStr = "vbsw.g2d.window GetDC failed"
		} else if err1 == 1000209 {
			errStr = "vbsw.g2d.window ChoosePixelFormat failed"
		} else if err1 == 1000210 {
			errStr = "vbsw.g2d.window SetPixelFormat failed"
		} else if err1 == 1000211 {
			errStr = "vbsw.g2d.window wglCreateContext failed"
		} else if err1 == 1000212 {
			errStr = "vbsw.g2d.window wglMakeCurrent failed"
		} else if err1 == 1000213 {
			errStr = "vbsw.g2d.window cdata.get failed"
		} else {
			errStr = "vbsw.g2d.window failed"
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
	return nil
}

type DefaultWindow struct {
}

func (_ *DefaultWindow) OnConfig(config *Configuration) error {
	return nil
}

func (_ *DefaultWindow) OnCreate(widget *Widget) error {
	return nil
}

func (_ *DefaultWindow) OnUpdate() error {
	return nil
}

func (_ *DefaultWindow) OnClose() (bool, error) {
	return true, nil
}

func (_ *DefaultWindow) OnShow() error {
	return nil
}

func (_ *DefaultWindow) OnResize() error {
	return nil
}

func (_ *DefaultWindow) OnKeyDown(keyCode int, repeated uint) error {
	return nil
}

func (_ *DefaultWindow) OnKeyUp(keyCode int) error {
	return nil
}

func (_ *DefaultWindow) OnTextureLoaded(textureId int) error {
	return nil
}

func (_ *DefaultWindow) OnDestroy() error {
	return nil
}
