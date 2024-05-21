/*
 *          Copyright 2024, Vitali Baumtrok.
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
	"unsafe"
)

func Init() {
	mutex.Lock()
	if !initialized {
		var n1, n2 C.int
		var err1, err2 C.longlong
		C.g2d_init(&n1, &n2, &err1, &err2)
		if err1 == 0 {
			MaxTexSize, MaxTexUnits = int(n1), int(n2)
			initialized, initFailed = true, false
			quitting = false
		} else {
			initFailed = true
			Err = toError(int64(err1), int64(err2), 0, "", nil)
		}
		mutex.Unlock()
	} else {
		mutex.Unlock()
		panic("g2d engine already initialized")
	}
}

func MainLoop(window ...interface{}) {
	mutex.Lock()
	if !initFailed {
		if initialized {
			if !running {
				Err = nil
				time.Reset()
				if hasAny(window) {
					running = true
					toMainLoop.reset()
					Show(window...)
					mutex.Unlock()
					C.g2d_mainloop_process_messages()
					mutex.Lock()
					running = false
					cleanUp()
				}
				mutex.Unlock()
			} else {
				mutex.Unlock()
				panic("g2d MainLoop already running")
			}
		} else {
			mutex.Unlock()
			panic("g2d not initialized")
		}
	} else {
		mutex.Unlock()
	}
}

func hasAny(windows []interface{}) bool {
	for _, window := range windows {
		if window != nil {
			return true
		}
	}
	return false
}

func cleanUp() {
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
	C.g2d_mainloop_clean_up()
}

func Show(window ...interface{}) {
	for _, wnd := range window {
		if wnd != nil {
			switch abst := wnd.(type) {
			case abstractWindow:
				toMainLoop.postMsg(&tLaunchWindowRequest{abst: abst})
			default:
				panic("g2d Window not embedded")
			}
		}
	}
}

func quitMainLoop(err1Ev, wndId int64) {
	if running && !quitting {
		var err1, err2 C.longlong
		C.g2d_mainloop_post_quit(&err1, &err2)
		quitting = true
		if err1 != 0 && Err != nil {
			Err = toError(err1Ev, int64(err2), wndId, "", nil)
		}
	}
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

//export g2dMainLoopInit
func g2dMainLoopInit() {
	go toMainLoop.messageThread()
}

//export g2dProcessToMainLoopMessages
func g2dProcessToMainLoopMessages() {
	messages := toMainLoop.messages()
	for _, message := range messages {
		message.processRequest()
	}
}

//export g2dWindowMove
func g2dWindowMove(cbIdC C.int) {
	wnd := abstCbs[int(cbIdC)].impl()
	msg := &tLogicMessage{typeId: wndMoveType, nanos: time.Nanos()}
	msg.props.update(wnd.dataC)
	wnd.msgs <- msg
}

//export g2dWindowResize
func g2dWindowResize(cbIdC C.int) {
	wnd := abstCbs[int(cbIdC)].impl()
	msg := &tLogicMessage{typeId: wndResizeType, nanos: time.Nanos()}
	msg.props.update(wnd.dataC)
	wnd.msgs <- msg
	wnd.gfxImpl.msgs <- &tGraphicsMessage{typeId: wndResizeType, valA: msg.props.ClientWidth, valB: msg.props.ClientHeight}
}

//export g2dClose
func g2dClose(cbIdC C.int) {
	wnd := abstCbs[int(cbIdC)].impl()
	wnd.msgs <- (&tLogicMessage{typeId: quitReqType, nanos: time.Nanos()})
}

//export g2dKeyDown
func g2dKeyDown(cbIdC, code C.int, repeated C.int) {
	wnd := abstCbs[int(cbIdC)].impl()
	msg := &tLogicMessage{typeId: keyDownType, valA: int(code), repeated: uint(repeated), nanos: time.Nanos()}
	msg.props.update(wnd.dataC)
	wnd.msgs <- msg
}

//export g2dKeyUp
func g2dKeyUp(cbIdC, code C.int) {
	wnd := abstCbs[int(cbIdC)].impl()
	msg := &tLogicMessage{typeId: keyUpType, valA: int(code), nanos: time.Nanos()}
	msg.props.update(wnd.dataC)
	wnd.msgs <- msg
}

//export g2dMouseMove
func g2dMouseMove(cbIdC C.int) {
	wnd := abstCbs[int(cbIdC)].impl()
	msg := &tLogicMessage{typeId: msMoveType, nanos: time.Nanos()}
	msg.props.update(wnd.dataC)
	wnd.msgs <- msg
}

//export g2dButtonDown
func g2dButtonDown(cbIdC, code, doubleClick C.int) {
	wnd := abstCbs[int(cbIdC)].impl()
	msg := &tLogicMessage{typeId: buttonDownType, valA: int(code), repeated: uint(doubleClick), nanos: time.Nanos()}
	msg.props.update(wnd.dataC)
	wnd.msgs <- msg
}

//export g2dButtonUp
func g2dButtonUp(cbIdC, code, doubleClick C.int) {
	wnd := abstCbs[int(cbIdC)].impl()
	msg := &tLogicMessage{typeId: buttonUpType, valA: int(code), repeated: uint(doubleClick), nanos: time.Nanos()}
	msg.props.update(wnd.dataC)
	wnd.msgs <- msg
}

//export g2dWheel
func g2dWheel(cbIdC C.int, wheel C.float) {
	wnd := abstCbs[int(cbIdC)].impl()
	msg := &tLogicMessage{typeId: wheelType, valB: float32(wheel), nanos: time.Nanos()}
	msg.props.update(wnd.dataC)
	wnd.msgs <- msg
}

//export g2dWindowMinimize
func g2dWindowMinimize(cbIdC C.int) {
	wnd := abstCbs[int(cbIdC)].impl()
	msg := &tLogicMessage{typeId: minimizeType, nanos: time.Nanos()}
	wnd.msgs <- msg
}

//export g2dWindowRestore
func g2dWindowRestore(cbIdC C.int) {
	wnd := abstCbs[int(cbIdC)].impl()
	msg := &tLogicMessage{typeId: restoreType, nanos: time.Nanos()}
	wnd.msgs <- msg
}
