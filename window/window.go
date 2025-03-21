/*
 *          Copyright 2025, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

// Package window creates a window with OpenGL 3.0 context.
package window

import (
	"github.com/vbsw/golib/cdata"
	"sync"
)

var (
	mutex          sync.Mutex
	bufferedEvents []Event
	bufferedDelta  []uint32
	currentEvents  []Event
	currentDelta   []uint32
	startMillis    uint64
)

type tInitializer struct {
}

func NewInitializer() cdata.CData {
	return new(tInitializer)
}

type Event interface {
	ProcessEvent(millis uint32)
}

type EventFactory interface {
	NewMainLoopStartedEvent() Event
	NewKeyPressedEvent(windowId, keyCode, keyRepeated int) Event
	NewKeyReleasedEvent(windowId, keyCode int) Event
}

type Window interface {
}

type Properties struct {
	ClientX, ClientY                  int
	ClientWidth, ClientHeight         int
	ClientWidthMin, ClientHeightMin   int
	ClientWidthMax, ClientHeightMax   int
	WindowId                          int
	MouseLocked, Borderless, Dragable bool
	Resizable, Fullscreen, Centered   bool
	AutoUpdate                        bool
	Title                             string
}
