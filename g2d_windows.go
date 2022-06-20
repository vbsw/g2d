/*
 *          Copyright 2022, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

package g2d

import (
	"github.com/vbsw/oglr"
	"github.com/vbsw/oglwnd"
)

// Parameters are the initialization parameters for Window.
type Parameters struct {
	oglwnd.Parameters
}

// Window is a window with OpenGL context.
type Window struct {
	data   oglwnd.Window
	params *Parameters
}

// tDummy is helper to load OpenGL 3.0 functions.
type tDummy struct {
	oglwnd.Window
}

// Init initializes g2d. Call this function before calling anyother function.
func Init() error {
	var err error
	if !initialized {
		var dummy tDummy
		err = dummy.Create()
		if err == nil {
			ctx := dummy.Context()
			err = ctx.MakeCurrent()
			if err == nil {
				err = oglr.Init()
				if err == nil {
					err = ctx.Release()
				} else {
					ctx.Release()
				}
			}
			err = dummy.Destroy(err)
			err = dummy.ReleaseMemory(err)
		}
		initialized = bool(err == nil)
	}
	return err
}

// ProcessWindowEvents retrieves messages from thread's message queue for all windows and calls window's
// handler to process it. This function blocks until further messages are available and returns only if
// all windows are destroyed.
func ProcessWindowEvents() {
	if initialized {
		oglwnd.ProcessEvents()
	} else {
		panic(notInitialized)
	}
}

// NewWindow creates a new instance of Window and returns it.
func NewWindow(params *Parameters) *Window {
	window := new(Window)
	window.params = params
	return window
}

// Init allocates window's ressources.
func (wnd *Window) Init() error {
	if initialized {
		return nil
	}
	panic(notInitialized)
}

// Destroy closes window and releases ressources associated with it.
func (wnd *Window) Destroy() error {
	return nil
}

// Create creates objects in win32.
func (dummy *tDummy) Create() error {
	err := dummy.Allocate()
	if err == nil {
		params := new(oglwnd.Parameters)
		params.Dummy = true
		err = dummy.Init(params)
		if err == nil {
			err = dummy.Window.Create()
		}
		if err != nil {
			dummy.Window.ReleaseMemory()
		}
	}
	return err
}

// Destroy releases win32 objects.
func (dummy *tDummy) Destroy(err error) error {
	if err == nil {
		err = dummy.Window.Destroy()
	} else {
		dummy.Window.Destroy()
	}
	return err
}

// ReleaseMemory releases struct memory allocated in C.
func (dummy *tDummy) ReleaseMemory(err error) error {
	if err == nil {
		err = dummy.Window.ReleaseMemory()
	} else {
		dummy.Window.ReleaseMemory()
	}
	return err
}
