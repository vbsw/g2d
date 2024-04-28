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
			Err = toError(int64(err1), int64(err2), 0, "")
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

func launchWindow(window Window) {
	/*
	   wnd := new(tWindow)
	   wnd.state = configState
	   wnd.abst = window
	   wnd.wgt = new(Widget)
	   wnd.wgt.msgs = make(chan *tLogicMessage, 1000)
	   wnd.wgt.quitted = make(chan bool, 1)
	   wnd.cbId = register(wnd)
	   wnd.cbIdStr = strconv.FormatInt(int64(wnd.cbId), 10)

	   wnd.wgt.Gfx.msgs = make(chan *tGMessage, 1000)
	   wnd.wgt.Gfx.quitted = make(chan bool, 1)
	   wnd.wgt.Gfx.rBuffer = &wnd.wgt.Gfx.buffers[0]
	   wnd.wgt.Gfx.wBuffer = &wnd.wgt.Gfx.buffers[0]
	   wnd.wgt.Gfx.initEntities()

	   go wnd.logicThread()
	   wnd.wgt.msgs <- (&tLogicMessage{typeId: configType, nanos: nanosDelta()})
	*/
}

func quitMainLoop() int64 {
	var err2 C.longlong
	C.g2d_mainloop_post_quit(&err2)
	return int64(err2)
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
