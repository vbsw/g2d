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
				var windowValid bool
				Err = nil
				time.Reset()
				for _, windw := range window {
					if windw != nil {
						windowValid = true
						launchWindow(windw)
					}
				}
				if windowValid {
					mutex.Unlock()
					C.g2d_mainloop_process_messages()
					mutex.Lock()
					running = false
					C.g2d_mainloop_clean_up()
					events = events[:0]
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

func Show(window ...Window) {
	mutex.Lock()
	if !initFailed {
		if initialized {
			if !quitting {
				for _, windw := range window {
					if windw != nil {
						launchWindow(windw)
					}
				}
			}
			mutex.Unlock()
		} else {
			mutex.Unlock()
			panic("g2d not initialized")
		}
	} else {
		mutex.Unlock()
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
	mutex.Lock()
	running = true
	switchEvents()
	mutex.Unlock()
	for _, event := range eventsOn {
		event.OnEvent()
	}
}

//export g2dMainLoopProcessCustomEvents
func g2dMainLoopProcessCustomEvents(additional *C.int) {
	mutex.Lock()
	switchEvents()
	mutex.Unlock()
	for _, event := range eventsOn {
		event.OnEvent()
	}
	if len(eventsOn) > 1 {
		*additional = C.int(len(eventsOn) - 1)
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
