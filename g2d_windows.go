/*
 *          Copyright 2023, Vitali Baumtrok.
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
		C.g2d_init(&errNumC, &errWin32C)
		if errNumC == 0 {
			startTime = time.Now()
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

func Show(window Window) {
	if !initFailed && window != nil {
		if initialized {
			wnd := newWindow(window)
			wnd.loopId = mainLoop.register(wnd)
			mainLoop.postMessage(&tConfigWindowRequest{window: wnd}, 1000)
			mainLoop.run()
		} else {
			panic("g2d is not initialized")
		}
	}
}

func (window *tWindow) logicThread() {
	for {
		msg := window.nextLMessage()
		window.wgt.CurrEventNanos = msg.nanos
		switch msg.typeId {
		case configType:
			window.onConfig()
		case createType:
			window.onCreate()
		case quitType:
			window.onQuit()
		case leaveType:
			mainLoop.postMessage(&tDestroyWindowRequest{window: window}, 1000)
			break
		}
	}
}

func (window *tWindow) onConfig() {
	config := newConfiguration()
	err := window.abst.OnConfig(config)
	window.autoUpdate = config.AutoUpdate
	if err == nil {
		mainLoop.postMessage(&tCreateWindowRequest{window: window, config: config}, 1000)
	} else {
		window.onError(err)
	}
}

func (window *tWindow) onCreate() {
	err := window.abst.OnCreate(window.wgt)
	if err == nil {
		mainLoop.postMessage(&tShowWindowRequest{window: window}, 1000)
	} else {
		window.onError(err)
	}
}

func (window *tWindow) onQuit() {
	window.abst.OnDestroy()
	if window.state == 0 {
		window.state = 10
		window.wgt.msgs <- &tLMessage{typeId: leaveType, nanos: deltaNanos()}
	}
}

func (window *tWindow) onError(err error) {
	appendError(err)
	window.onQuit()
}

func (window *tWindow) nextLMessage() *tLMessage {
	var message *tLMessage
	if window.autoUpdate {
		select {
		case msg := <-window.wgt.msgs:
			message = msg
		default:
			message = &tLMessage{typeId: updateType, nanos: deltaNanos()}
		}
	} else {
		message = <-window.wgt.msgs
	}
	if window.state == 10 && message.typeId != leaveType {
		message = nil
	}
	return message
}

func (loop *tMainLoop) register(window *tWindow) int {
	loop.mutex.Lock()
	defer loop.mutex.Unlock()
	if len(loop.wndsUnused) == 0 {
		loop.wndsUsed = append(loop.wndsUsed, window)
		return len(loop.wndsUsed) - 1
	}
	lastIndex := len(loop.wndsUnused) - 1
	index := loop.wndsUnused[lastIndex]
	loop.wndsUnused = loop.wndsUnused[:lastIndex]
	loop.wndsUsed[index] = window
	return index
}

func (loop *tMainLoop) unregister(loopId int) int {
	loop.mutex.Lock()
	defer loop.mutex.Unlock()
	loop.wndsUsed[loopId] = nil
	loop.wndsUnused = append(loop.wndsUnused, loopId)
	return len(loop.wndsUsed) - len(loop.wndsUnused)
}

func (loop *tMainLoop) postMessage(msg interface{}, errNumC C.int) {
	var errNumOrigC C.int
	var errWin32C C.g2d_ul_t
	loop.mutex.Lock()
	defer loop.mutex.Unlock()
	C.g2d_post_message(&errNumOrigC, &errWin32C)
	if errNumOrigC == 0 {
		loop.msgs.Put(msg)
	} else {
		C.g2d_quit_message_queue()
		appendError(toError(errNumC, errWin32C, nil))
	}
}

func (loop *tMainLoop) run() {
	loop.mutex.Lock()
	if !mainLoop.running {
		mainLoop.running = true
		mainLoop.mutex.Unlock()
		C.g2d_process_messages()
		loop.mutex.Lock()
		mainLoop.running = false
		mainLoop.mutex.Unlock()
	} else {
		mainLoop.mutex.Unlock()
	}
}

//export g2dProcessMessage
func g2dProcessMessage() {
	message := mainLoop.nextMessage()
	if message != nil {
		switch msg := message.(type) {
		case *tConfigWindowRequest:
			configWindow(msg.window)
		case *tCreateWindowRequest:
			createWindow(msg.window, msg.config)
		case *tShowWindowRequest:
			showWindow(msg.window)
		case *tDestroyWindowRequest:
			destroyWindow(msg.window)
		}
	}
}

func configWindow(window *tWindow) {
	window.wgt.msgs <- (&tLMessage{typeId: configType, nanos: deltaNanos()})
}

func createWindow(window *tWindow, config *Configuration) {
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
		var errWin32C C.g2d_ul_t
		window.cbId = cb.register(window)
		C.g2d_window_create(&window.dataC, C.int(window.cbId), x, y, w, h, wn, hn, wx, hx, b, d, r, f, l, c, t, &errNumC, &errWin32C)
		C.g2d_free(t)
		if errNumC == 0 {
			msg := &tLMessage{typeId: createType, nanos: deltaNanos()}
			msg.props.update(window.dataC)
			window.wgt.msgs <- msg
		} else {
			appendError(toError(errNumC, errWin32C, nil))
			window.wgt.msgs <- (&tLMessage{typeId: quitType, nanos: deltaNanos()})
		}
	} else {
		appendError(toError(errNumC+100, 0, nil))
		window.wgt.msgs <- (&tLMessage{typeId: quitType, nanos: deltaNanos()})
	}
}

func showWindow(window *tWindow) {
	// TODO show window
}

func destroyWindow(window *tWindow) {
	if window.cbId >= 0 {
		var errNumC C.int
		var errWin32C C.g2d_ul_t
		C.g2d_window_destroy(window.dataC, &errNumC, &errWin32C)
		cb.unregister(window.cbId)
		window.cbId = -1
		if errNumC != 0 {
			appendError(toError(errNumC, errWin32C, nil))
		}
	}
	registered := mainLoop.unregister(window.loopId)
	if registered <= 0 {
		C.g2d_quit_message_queue()
	}
}

func newWindow(abst Window) *tWindow {
	window := new(tWindow)
	window.cbId = -1
	window.abst = abst
	window.wgt = new(Widget)
	window.wgt.msgs = make(chan *tLMessage, 1024)
	go window.logicThread()
	return window
}

func initCustParams(params []interface{}) {
	for i, param := range params {
		var ok, used bool
		errGen, ok = param.(ErrorGenerator)
		used = used || ok
		errLog, ok = param.(ErrorLogger)
		used = used || ok
		if !used {
			panic(fmt.Sprintf("parameter %d is not used", i))
		}
	}
}

func initDefaultParams() {
	if errGen == nil {
		errGen = &errHandler
	}
	if errLog == nil {
		errLog = &errHandler
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

func toError(errNumC C.int, errWin32C C.g2d_ul_t, errStrC *C.char) error {
	var errStr string
	if errStrC != nil {
		errStr = C.GoString(errStrC)
		C.g2d_free(unsafe.Pointer(errStrC))
	}
	return errGen.ToError(uint64(errNumC), uint64(errWin32C), errStr)
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
		errStr = fmt.Sprintf(errStrGetProcAddr, "glUniformMatrix4fv")
	case 232:
		errStr = fmt.Sprintf(errStrGetProcAddr, "glUniformMatrix2x3fv")
	case 233:
		errStr = fmt.Sprintf(errStrGetProcAddr, "glGenerateMipmap")
	case 234:
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

func toTString(str string) (unsafe.Pointer, C.int) {
	var strT unsafe.Pointer
	var errNumC C.int
	strC := unsafe.Pointer(C.CString(str))
	if strC != nil {
		C.g2d_to_tstr(&strT, strC, &errNumC)
		C.g2d_free(strC)
	} else {
		errNumC = 100
	}
	return strT, errNumC
}
