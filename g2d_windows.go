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

func MainLoop(window ...Window) {
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

func hasAny(windows []Window) bool {
	for _, window := range windows {
		if window != nil {
			return true
		}
	}
	return false
}

func cleanUp() {
	for _, wnd := range wndCbs {
		if wnd != nil {
			var err1, err2 C.longlong
			if wnd.wgt != nil {
				wgt := wnd.wgt
				wgt.msgs <- (&tLogicMessage{typeId: quitType, nanos: time.Nanos()})
				<-wgt.quitted
			}
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

func Show(window ...Window) {
	if hasAny(window) {
		toMainLoop.postMsg(&tLaunchWindowRequest{windows: window})
	}
}

func launchWindow(window Window) {
	wnd := newWindowWrapper(window)
	go wnd.logicThread()
	wnd.wgt.msgs <- (&tLogicMessage{typeId: configType, nanos: time.Nanos()})
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

//export g2dResize
func g2dResize(cbIdC C.int) {
	wnd := wndCbs[int(cbIdC)]
	wgt := wnd.wgt
	if wgt != nil {
		msg := &tLogicMessage{typeId: resizeType, nanos: time.Nanos()}
		msg.props.update(wnd.dataC)
		wgt.msgs <- msg
		//wnd.wgt.Gfx.msgs <- &tGMessage{typeId: resizeType, valA: msg.props.ClientWidth, valB: msg.props.ClientHeight}
	}
}

//export g2dKeyDown
func g2dKeyDown(cbIdC, code C.int, repeated C.int) {
	wnd := wndCbs[int(cbIdC)]
	wgt := wnd.wgt
	if wgt != nil {
		msg := &tLogicMessage{typeId: keyDownType, valA: int(code), repeated: uint(repeated), nanos: time.Nanos()}
		msg.props.update(wnd.dataC)
		wgt.msgs <- msg
	}
}

//export g2dKeyUp
func g2dKeyUp(cbIdC, code C.int) {
	wnd := wndCbs[int(cbIdC)]
	wgt := wnd.wgt
	if wgt != nil {
		msg := &tLogicMessage{typeId: keyUpType, valA: int(code), nanos: time.Nanos()}
		msg.props.update(wnd.dataC)
		wgt.msgs <- msg
	}
}

//export g2dClose
func g2dClose(cbIdC C.int) {
	wnd := wndCbs[int(cbIdC)]
	wgt := wnd.wgt
	if wgt != nil {
		wgt.msgs <- (&tLogicMessage{typeId: quitReqType, nanos: time.Nanos()})
	}
}
