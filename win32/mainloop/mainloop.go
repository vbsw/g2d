/*
 *          Copyright 2024, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

// Package mainloop runs main message loop in win32.
package mainloop

// #cgo CFLAGS: -DUNICODE
// #cgo LDFLAGS: -luser32
// #include "mainloop.h"
import "C"
import (
	"github.com/vbsw/golib/queue"
	"sync"
)

var (
	mutex   sync.Mutex
	events  queue.Queue
	running bool
)

// Event is the interface to abstract event execution.
type Event interface {
	OnEvent()
}

// Run starts the main loop.
func Run() {
	C.g2d_mainloop_process_messages()
	mutex.Lock()
	running = false
	C.g2d_mainloop_clean_up()
	events.Reset(-1)
	mutex.Unlock()
}

// Post posts the event to the message queue. Returns an error code.
func Post(event Event) int64 {
	var err C.longlong
	mutex.Lock()
	if running {
		C.g2d_mainloop_post_custom(&err)
	}
	if err == 0 {
		events.Put(event)
	}
	mutex.Unlock()
	return int64(err)
}

// Quit stops the main loop. Unprocessed events are dropped. Returns an error code.
func Quit() int64 {
	var err C.longlong
	C.g2d_mainloop_post_quit(&err)
	return int64(err)
}

//export g2dMainLoopInit
func g2dMainLoopInit() {
	mutex.Lock()
	running = true
	evs := events.All()
	mutex.Unlock()
	for _, ev := range evs {
		switch event := ev.(type) {
		case Event:
			event.OnEvent()
		}
	}
}

//export g2dMainLoopProcessCustomEvents
func g2dMainLoopProcessCustomEvents(additional *C.int) {
	mutex.Lock()
	evs := events.All()
	mutex.Unlock()
	for _, ev := range evs {
		switch event := ev.(type) {
		case Event:
			event.OnEvent()
		}
	}
	if len(evs) > 1 {
		*additional = C.int(len(evs) - 1)
	}
}
