/*
 *          Copyright 2025, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

package window

// #cgo CFLAGS: -DG2D_WINDOW_WIN32 -DUNICODE
// #cgo LDFLAGS: -luser32 -lgdi32 -lOpenGL32
// #include "window.h"
import "C"
import (
	"errors"
	"fmt"
	"strconv"
	"unsafe"
)

const (
	loadFunctionFailed = "g2d window load %s function failed"
)

var (
	evFactory EventFactory
	running          bool
)

func (init tInitializer) CProcFunc() unsafe.Pointer {
	return C.g2d_window_init
}

func (props *Properties) boolsToCInt() (C.int, C.int, C.int, C.int, C.int, C.int) {
	var c, l, b, d, r, f C.int
	if props.Centered {
		c = 1
	}
	if props.MouseLocked {
		l = 1
	}
	if props.Borderless {
		b = 1
	}
	if props.Dragable {
		d = 1
	}
	if props.Resizable {
		r = 1
	}
	if props.Fullscreen {
		f = 1
	}
	return c, l, b, d, r, f
}

func (init tInitializer) SetCData(data unsafe.Pointer) {
}

func (init tInitializer) ToError(err1, err2 int64, info string) error {
	return ToError(err1, err2, info)
}

func ToError(err1, err2 int64, info string) error {
	var err error
	if err1 > 0 {
		var errStr string
		/* 201 - 300, 1000201 - 1000300 */
		if err1 < 1000000 {
			if err1 > 200 && err1 < 301 {
				errStr = "memory allocation failed"
			}
		} else {
			switch err1 {
			case 1000201:
				errStr = "g2d window GetModuleHandle failed"
			case 1000202:
				errStr = fmt.Sprintf(loadFunctionFailed, "wglChoosePixelFormatARB")
			case 1000203:
				errStr = fmt.Sprintf(loadFunctionFailed, "wglCreateContextAttribsARB")
			case 1000204:
				errStr = fmt.Sprintf(loadFunctionFailed, "wglSwapIntervalEXT")
			case 1000205:
				errStr = fmt.Sprintf(loadFunctionFailed, "wglGetSwapIntervalEXT")
			}
		}
		if len(errStr) > 0 {
			errStr = errStr + " (" + strconv.FormatInt(err1, 10)
			if err2 == 0 {
				errStr = errStr + ")"
			} else {
				errStr = errStr + ", " + strconv.FormatInt(err2, 10) + ")"
			}
			if len(info) > 0 {
				errStr = errStr + "; " + info
			}
			err = errors.New(errStr)
		}
	}
	return err
}

func MainLoop(eventFactory EventFactory) {
	if eventFactory != nil {
		if !running {
			running = true
			evFactory = eventFactory
			bufferedEvents = make([]Event, 0, 7)
			bufferedDelta = make([]uint32, 0, 7)
			currentEvents = make([]Event, 0, 7)
			currentDelta = make([]uint32, 0, 7)
			C.g2d_window_mainloop()
			running = false
		} else {
			panic("MainLoop already running")
		}
	} else {
		panic("EventFactory must not be nil")
	}
}

func PostEvent(event Event) (int64, int64) {
	var millis, err1, err2 C.longlong
	if event != nil {
		mutex.Lock()
		C.g2d_window_post_custom_msg(&millis, &err1, &err2)
		if err1 == 0 {
			bufferedEvents = append(bufferedEvents, event)
			bufferedDelta = append(bufferedDelta, uint32(uint64(millis) - startMillis))
		}
		mutex.Unlock()
	}
	return int64(err1), int64(err2)
}

func QuitMainLoop() (int64, int64) {
	var err1, err2 C.longlong
	C.g2d_window_post_quit_msg(&err1, &err2)
	return int64(err1), int64(err2)
}

func ClearMessageQueue() {
	C.g2d_window_mainloop_clean_up()
}

func Create(props *Properties) (unsafe.Pointer, int64, int64) {
	var err1, err2 C.longlong
	var data unsafe.Pointer
	var t unsafe.Pointer
	var ts C.size_t
	x := C.int(props.ClientX)
	y := C.int(props.ClientY)
	w := C.int(props.ClientWidth)
	h := C.int(props.ClientHeight)
	wn := C.int(props.ClientWidthMin)
	hn := C.int(props.ClientHeightMin)
	wx := C.int(props.ClientWidthMax)
	hx := C.int(props.ClientHeightMax)
	c, l, b, d, r, f := props.boolsToCInt()
	if len(props.Title) > 0 {
		bytes := *(*[]byte)(unsafe.Pointer(&(props.Title)))
		t, ts = unsafe.Pointer(&bytes[0]), C.size_t(len(props.Title))
	}
	C.g2d_window_create(&data, C.int(props.WindowId), x, y, w, h, wn, hn, wx, hx, b, d, r, f, l, c, t, ts, &err1, &err2)
/*
	if err1 == 0 {
		msg := &tLogicMessage{typeId: createType, nanos: time.Nanos()}
		msg.props.update(wnd.dataC)
		wnd.msgs <- msg
	} else {
		toMainLoop.postErr(toError(int64(err1), int64(err2), int64(wnd.WindowId), "", nil))
	}
*/
	return data, int64(err1), int64(err2)
}

func Show(data unsafe.Pointer) (int64, int64) {
	var err1, err2 C.longlong
	C.g2d_window_show(data, &err1, &err2)
	return int64(err1), int64(err2)
}

func processEvents(newEvent Event, newDelta uint32) {
	var tmpEvents []Event
	var tmpDelta []uint32
	mutex.Lock()
	if newEvent == nil {
		tmpEvents = bufferedEvents
		tmpDelta = bufferedDelta
	} else {
		tmpEvents = append(bufferedEvents, newEvent)
		tmpDelta = append(bufferedDelta, newDelta)
	}
	bufferedEvents = currentEvents[:0]
	bufferedDelta = currentDelta[:0]
	mutex.Unlock()
	currentEvents = tmpEvents
	currentDelta = tmpDelta
	for i, event := range currentEvents {
		event.ProcessEvent(currentDelta[i])
	}
}

//export g2dWindowMainLoopStart
func g2dWindowMainLoopStart(millis C.longlong) {
	startMillis = uint64(millis)
	evFactory.NewMainLoopStartedEvent().ProcessEvent(0)
}

//export g2dWindowMainLoopEvent
func g2dWindowMainLoopEvent() {
	processEvents(nil, 0)
}

//export g2dKeyDown
func g2dKeyDown(wndId, code, repeated C.int, millis C.longlong) {
	processEvents(evFactory.NewKeyPressedEvent(int(wndId), int(code), int(repeated)), uint32(uint64(millis) - startMillis))
}

//export g2dKeyUp
func g2dKeyUp(wndId, code C.int, millis C.longlong) {
	processEvents(evFactory.NewKeyReleasedEvent(int(wndId), int(code)), uint32(uint64(millis) - startMillis))
}
