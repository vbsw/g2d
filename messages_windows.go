/*
 *          Copyright 2024, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

package g2d

// #include "g2d.h"
import "C"
import (
	"sync"
)

var (
	toMainLoop tToMainLoop
)

type tMainLoopRequest interface {
	processRequest()
	err1(err1 C.longlong) int64
	wndId() int64
}

type tToMainLoop struct {
	errs       []*Error
	buffered   []tMainLoopRequest
	posted     []tMainLoopRequest
	processing []tMainLoopRequest
	nextAvail  chan bool
	quitted    chan bool
	mutex      sync.Mutex
	quitting   bool
	errAvail   bool
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

type tLaunchWindowRequest struct {
	windows []Window
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

func (toMainLoop *tToMainLoop) reset() {
	toMainLoop.errs = toMainLoop.errs[:0]
	toMainLoop.buffered = toMainLoop.buffered[:0]
	toMainLoop.posted = toMainLoop.posted[:0]
	toMainLoop.nextAvail = make(chan bool, 1000)
	toMainLoop.quitted = make(chan bool, 1)
	toMainLoop.quitting = false
	toMainLoop.errAvail = false
}

func (toMainLoop *tToMainLoop) messageThread() {
	for true {
		<-toMainLoop.nextAvail
		toMainLoop.mutex.Lock()
		if toMainLoop.errAvail {
			if toMainLoop.errs[0] != nil {
				setErrorSynced(toMainLoop.errs[0])
			} else {
				setErrorSynced(toError(4000, 0, 0, "toMainLoop.errs[0] == nil", nil))
			}
			toMainLoop.mutex.Unlock()
			break
		} else if toMainLoop.quitting {
			toMainLoop.mutex.Unlock()
			break
		} else if len(toMainLoop.buffered) > 0 {
			mutex.Lock()
			if running {
				var err1, err2 C.longlong
				C.g2d_mainloop_post_custom(&err1, &err2)
				if err1 == 0 {
					toMainLoop.posted = append(toMainLoop.posted, toMainLoop.buffered...)
					toMainLoop.buffered = toMainLoop.buffered[:0]
					mutex.Unlock()
					toMainLoop.mutex.Unlock()
				} else {
					setError(toError(int64(err1), int64(err2), 0, "", nil))
					mutex.Unlock()
					toMainLoop.mutex.Unlock()
					break
				}
			} else {
				mutex.Unlock()
				toMainLoop.mutex.Unlock()
				break
			}
		} else {
			toMainLoop.mutex.Unlock()
		}
	}
	toMainLoop.quitted <- true
}

func (toMainLoop *tToMainLoop) postMsg(msg tMainLoopRequest) {
	toMainLoop.mutex.Lock()
	toMainLoop.buffered = append(toMainLoop.buffered, msg)
	toMainLoop.nextAvail <- true
	toMainLoop.mutex.Unlock()
}

func (toMainLoop *tToMainLoop) postErr(err *Error) {
	toMainLoop.mutex.Lock()
	toMainLoop.errs = append(toMainLoop.errs, err)
	toMainLoop.errAvail = true
	toMainLoop.nextAvail <- true
	toMainLoop.mutex.Unlock()
}

func (toMainLoop *tToMainLoop) messages() []tMainLoopRequest {
	postedTmp := toMainLoop.posted
	toMainLoop.posted = toMainLoop.processing[:0]
	toMainLoop.processing = postedTmp
	return toMainLoop.processing
}

func (toMainLoop *tToMainLoop) quitMessageThread() {
	toMainLoop.mutex.Lock()
	toMainLoop.quitting = true
	toMainLoop.nextAvail <- true
	toMainLoop.mutex.Unlock()
	<-toMainLoop.quitted
}

func (req *tLaunchWindowRequest) processRequest() {
	for _, window := range req.windows {
		if window != nil {
			launchWindow(window)
		}
	}
}

func (req *tLaunchWindowRequest) err1(err1 C.longlong) int64 {
	return int64(err1)
}

func (req *tLaunchWindowRequest) wndId() int64 {
	return 0
}

func (req *tCreateWindowRequest) processRequest() {
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
			toMainLoop.postErr(toError(int64(err1), int64(err2), int64(wnd.cbId), "", nil))
		}
	} else {
		toMainLoop.postErr(toError(int64(err1), int64(err2), int64(wnd.cbId), "", nil))
	}
}

func (req *tCreateWindowRequest) err1(err1 C.longlong) int64 {
	return int64(err1)
}

func (req *tCreateWindowRequest) wndId() int64 {
	return int64(req.window.cbId)
}

func (req *tShowWindowRequest) processRequest() {
	var err1, err2 C.longlong
	C.g2d_window_show(req.window.dataC, &err1, &err2)
	if err1 == 0 {
		req.window.wgt.msgs <- (&tLogicMessage{typeId: showType, nanos: time.Nanos()})
	} else {
		toMainLoop.postErr(toError(int64(err1), int64(err2), int64(req.window.cbId), "", nil))
	}
}

func (req *tShowWindowRequest) err1(err1 C.longlong) int64 {
	return int64(err1)
}

func (req *tShowWindowRequest) wndId() int64 {
	return int64(req.window.cbId)
}

func (req *tDestroyWindowRequest) processRequest() {
	var err1, err2 C.longlong
	C.g2d_window_destroy(req.window.dataC, &err1, &err2)
	unregister(req.window.cbId)
	req.window.cbId = -1
	if err1 != 0 {
		toMainLoop.postErr(toError(int64(err1), int64(err2), int64(req.window.cbId), "", nil))
	}
}

func (req *tDestroyWindowRequest) err1(err1 C.longlong) int64 {
	return int64(err1)
}

func (req *tDestroyWindowRequest) wndId() int64 {
	return int64(req.window.cbId)
}
