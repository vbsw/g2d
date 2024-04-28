/*
 *          Copyright 2024, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

package g2d

// #include "g2d.h"
import "C"

var (
	events   []tEvent
	eventsOn []tEvent
)

type tEvent interface {
	OnEvent()
}

type tLogicMessage struct {
	typeId   int
	valA     int
	repeated uint
	nanos    int64
	props    Properties
	obj      interface{}
}

type tGraphicsMessage struct {
	typeId int
	valA   int
	valB   int
	valC   interface{}
	err    error
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

func switchEvents() {
	eventsTmp := events
	events = eventsOn[:0]
	eventsOn = eventsTmp
}

func postEvent(event tEvent, err1Ev, wndId int64) {
	var err1, err2 C.longlong
	mutex.Lock()
	C.g2d_mainloop_post_custom(&err1, &err2)
	if err1 == 0 {
		events = append(events, event)
	} else if Err == nil {
		Err = toError(err1Ev, int64(err2), wndId, "", nil)
		quitMainLoop(err1Ev, wndId)
	}
	mutex.Unlock()
}

func (req *tCreateWindowRequest) OnEvent() {
	var err1, err2 C.longlong
	config := req.config
	wnd := req.window
	x := C.int(config.ClientX)
	y := C.int(config.ClientY)
	w := C.int(config.ClientWidth)
	h := C.int(config.ClientHeight)
	wn := C.int(config.ClientWidthMin)
	hn := C.int(config.ClientHeightMin)
	wx := C.int(config.ClientWidthMax)
	hx := C.int(config.ClientHeightMax)
	c := toIntC(config.Centered)
	l := toIntC(config.MouseLocked)
	b := toIntC(config.Borderless)
	d := toIntC(config.Dragable)
	r := toIntC(config.Resizable)
	f := toIntC(config.Fullscreen)
	t, err1 := toTString(config.Title)
	if err1 == 0 {
		C.g2d_window_create(&wnd.dataC, C.int(wnd.cbId), x, y, w, h, wn, hn, wx, hx, b, d, r, f, l, c, t, &err1, &err2)
		C.g2d_free(t)
		if err1 == 0 {
			msg := &tLogicMessage{typeId: createType, nanos: time.Nanos()}
			msg.props.update(wnd.dataC)
			wnd.wgt.msgs <- msg
		} else {
			setError(toError(int64(err1), int64(err2), int64(wnd.cbId), "", nil))
		}
	} else {
		setError(toError(int64(err1), int64(err2), int64(wnd.cbId), "", nil))
	}
}

func (req *tShowWindowRequest) OnEvent() {
	var err1, err2 C.longlong
	C.g2d_window_show(req.window.dataC, &err1, &err2)
	if err1 == 0 {
		req.window.wgt.msgs <- (&tLogicMessage{typeId: showType, nanos: time.Nanos()})
	} else {
		setError(toError(int64(err1), int64(err2), int64(req.window.cbId), "", nil))
	}
}

func (req *tDestroyWindowRequest) OnEvent() {
	var err1, err2 C.longlong
	C.g2d_window_destroy(req.window.dataC, &err1, &err2)
	unregister(req.window.cbId)
	req.window.cbId = -1
	if err1 != 0 {
		setError(toError(int64(err1), int64(err2), int64(req.window.cbId), "", nil))
	}
}
