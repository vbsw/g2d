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

// Handler
type Handler interface {
	OnCreate()
	OnShow()
	OnClose() bool
	OnCustom(interface{})
	OnDestroy()
}

// DefaultHandler
type DefaultHandler struct {
}

// Parameters are the initialization parameters for Window.
type Parameters struct {
	oglwnd.Parameters
	Handler Handler
}

// Window is a window with OpenGL context.
type Window struct {
	data   oglwnd.Window
	state int
}

// tDummy is helper to load OpenGL 3.0 functions.
type tDummy struct {
	oglwnd.Window
}

type tHandlerAdapter struct {
	handler Handler
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

// Init allocates window's ressources.
func (wnd *Window) Init(params *Parameters) error {
	if initialized {
		var err error
		if wnd.state == 0 {
			err = wnd.data.Allocate()
			if err == nil {
				paramsOglWnd := newOglWndParams(params)
				err = wnd.data.Init(paramsOglWnd)
				if err == nil {
					wglCPF, wglCCA := oglr.WGLFunctions()
					wnd.data.SetWGLFunctions(wglCPF, wglCCA)
					err = wnd.data.Create()
					if err == nil {
						wnd.state = 1
					} else {
						wnd.data.Destroy()
					}
				}
				if err != nil {
					wnd.data.ReleaseMemory()
				}
			}
		}
		return err
	}
	panic(notInitialized)
}

// Show makes window visible.
func (wnd *Window) Show() error {
	var err error
	if wnd.state == 1 {
		err = wnd.data.Show()
		if err == nil {
			wnd.state = 2
		}
	}
	return err
}

// Destroy closes window and releases ressources associated with it.
func (wnd *Window) Destroy() error {
	var err error
	if wnd.state > 0 {
		err = wnd.data.Destroy()
		if err == nil {
			err = wnd.data.ReleaseMemory()
		} else {
			wnd.data.ReleaseMemory()
		}
		wnd.state = 0
	}
	return err
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

func newOglWndParams(params *Parameters) *oglwnd.Parameters {
	if params != nil {
		paramsOglWnd := new(oglwnd.Parameters)
		paramsOglWnd.ClientX = params.ClientX
		paramsOglWnd.ClientY = params.ClientY
		paramsOglWnd.ClientWidth = params.ClientWidth
		paramsOglWnd.ClientHeight = params.ClientHeight
		paramsOglWnd.ClientMinWidth = params.ClientMinWidth
		paramsOglWnd.ClientMinHeight = params.ClientMinHeight
		paramsOglWnd.ClientMaxWidth = params.ClientMaxWidth
		paramsOglWnd.ClientMaxHeight = params.ClientMaxHeight
		paramsOglWnd.Centered = params.Centered
		paramsOglWnd.Borderless = params.Borderless
		paramsOglWnd.Dragable = params.Dragable
		paramsOglWnd.Resizable = params.Resizable
		paramsOglWnd.Fullscreen = params.Fullscreen
		paramsOglWnd.MouseLocked = params.MouseLocked
		paramsOglWnd.Handler = params.newOglWndHandler()
		return paramsOglWnd
	}
	return nil
}

func (params *Parameters) newOglWndHandler() oglwnd.Handler {
	adpr := new(tHandlerAdapter)
	adpr.handler = params.Handler
	return adpr
}

func (adpr *tHandlerAdapter) OnCreate(wnd *oglwnd.Window) {
	adpr.handler.OnCreate()
}

func (adpr *tHandlerAdapter) OnShow(wnd *oglwnd.Window) {
	adpr.handler.OnShow()
}

func (adpr *tHandlerAdapter) OnClose(wnd *oglwnd.Window) bool {
	return adpr.handler.OnClose()
}

func (adpr *tHandlerAdapter) OnCustom(wnd *oglwnd.Window, event interface{}) {
	adpr.handler.OnCustom(event)
}

func (adpr *tHandlerAdapter) OnDestroy(wnd *oglwnd.Window) {
	adpr.handler.OnDestroy()
}

func (adpr *DefaultHandler) OnCreate() {
}

func (adpr *DefaultHandler) OnShow() {
}

func (adpr *DefaultHandler) OnClose() bool {
	return true
}

func (adpr *DefaultHandler) OnCustom(event interface{}) {
}

func (adpr *DefaultHandler) OnDestroy() {
}
