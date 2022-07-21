/*
 *          Copyright 2022, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

// Package g2d is a framework to create 2D graphic applications.
package g2d

import "C"
import (
	"time"
	"unsafe"
)

const (
	notInitialized    = "g2d not initialized"
	incativeWindow    = "inactive window"
	alreadyProcessing = "already processing events"
	messageFailed     = "post message failed"
	memoryAllocation  = "memory allocation failed"
	notEmbedded       = "type Window is not embedded"
)

var (
	Err error
	initialized bool
	processing  bool
	cb          tCallback
	startTime   time.Time
)

type Configuration struct {
	ClientX, ClientY                  int
	ClientWidth, ClientHeight         int
	ClientWidthMin, ClientHeightMin   int
	ClientWidthMax, ClientHeightMax   int
	MouseLocked, Borderless, Dragable bool
	Resizable, Fullscreen, Centered   bool
	AutoUpdate                        bool
	Title                             string
}

type Properties struct {
	MouseX, MouseY                    int
	ClientX, ClientY                  int
	ClientWidth, ClientHeight         int
	ClientWidthMin, ClientHeightMin   int
	ClientWidthMax, ClientHeightMax   int
	MouseLocked, Borderless, Dragable bool
	Resizable, Fullscreen             bool
	EventTime, UpdateTime             int64
	Title                             string
}

type Window struct {
	Props Properties
	abst tWindow
	dataC unsafe.Pointer
	active bool
	events chan interface{}
}

type tWindow interface {
	OnConfig(config *Configuration) error
	OnCreate() error
	OnShow() error
	OnKeyDown(key int, repeated uint) error
	OnKeyUp(key int) error
	OnUpdate() error
	OnClose() (bool, error)
	OnDestroy()
	config(window tWindow) *Configuration
	create(*Configuration)
	show()
	destroy()
	onCreate()
	onShow()
}

// Time returns nanoseconds.
func Time() int64 {
	if initialized {
		timeNow := time.Now()
		d := timeNow.Sub(startTime)
		return d.Nanoseconds()
	}
	panic(notInitialized)
}

type tCallback struct {
	wnds []tWindow
	unused []int
}

func (window *Window) OnConfig(config *Configuration) error {
	return nil
}

func (window *Window) OnCreate() error {
	return nil
}

func (window *Window) OnShow() error {
	return nil
}

func (window *Window) OnKeyDown(key int, repeated uint) error {
	return nil
}

func (window *Window) OnKeyUp(key int) error {
	return nil
}

func (window *Window) OnUpdate() error {
	return nil
}

func (window *Window) OnClose() (bool, error) {
	return true, nil
}

func (window *Window) OnDestroy() {
}

func (window *Window) Close() {
}

func (window *Window) Destroy() {
}

func (window *Window) PostEvent(event interface{}) {
}

func (window *Window) destroy() {
}

func (window *Window) config(abst tWindow) *Configuration {
	config := newConfiguration()
	window.Props.EventTime = Time()
	window.Props.UpdateTime = window.Props.EventTime
	window.events = make(chan interface{}, 1024 * 8)
	Err = abst.OnConfig(config)
	if Err == nil {
		window.abst = abst
	}
	return config
}

// Register returns a new id number for wnd. It will not be garbage collected until
// Unregister is called with this id.
func (cb *tCallback) Register(wnd tWindow) int {
	if len(cb.unused) == 0 {
		cb.wnds = append(cb.wnds, wnd)
		return len(cb.wnds) - 1
	}
	indexLast := len(cb.unused) - 1
	indexObj := cb.unused[indexLast]
	cb.unused = cb.unused[:indexLast]
	cb.wnds[indexObj] = wnd
	return indexObj
}

// Unregister makes wnd no more identified by id.
// This object may be garbage collected, now.
func (cb *tCallback) Unregister(id int) {
	cb.wnds[id] = nil
	cb.unused = append(cb.unused, id)
}

// UnregisterAll makes all regiestered wnds no more identified by id.
// These objects may be garbage collected, now.
func (cb *tCallback) UnregisterAll() {
	for i := 0; i < len(cb.wnds) && cb.wnds[i] != nil; i++ {
		cb.Unregister(i)
	}
}

func newConfiguration() *Configuration {
	config := new(Configuration)
	config.ClientX = 50
	config.ClientY = 50
	config.ClientWidth = 640
	config.ClientHeight = 480
	config.ClientWidthMin = 0
	config.ClientHeightMin = 0
	config.ClientWidthMax = 99999
	config.ClientHeightMax = 99999
	config.MouseLocked = false
	config.Borderless = false
	config.Dragable = false
	config.Resizable = true
	config.Fullscreen = false
	config.Centered = true
	config.AutoUpdate = true
	config.Title = "g2d - 0.1.0"
	return config
}

// toCInt converts bool value to C int value.
func toCInt(b bool) C.int {
	if b {
		return C.int(1)
	}
	return C.int(0)
}

/*
type Command struct {
	CloseReq bool
	CloseUnc bool
	Update   bool
}

type tManager interface {
	onCreate(nanos int64)
	onShow(nanos int64)
	onKeyDown(key int, repeated uint, nanos int64)
	onKeyUp(key int, nanos int64)
	onUpdate(nanos int64)
	onClose(nanos int64)
	onDestroy(nanos int64)
	onProps(nanos int64)
	onError(nanos int64)
	destroy()
}

type tManagerBase struct {
	data       unsafe.Pointer
	wndAbst    AbstractWindow
	wndBase    *Window
	props      Properties
	autoUpdate bool
}

type tManagerNoThreads struct {
	tManagerBase
	cmd Command
}

type tManagerLogicThread struct {
	tManagerBase
}

type tManagerGraphicThread struct {
}

type tManagerLogicGraphicThread struct {
}

type tModification struct {
	mouse, style, title bool
	fsToggle, pos, size bool
}

func (t *Time) NanosNow() int64 {
	return time()
}

func (que *Queue) init() {
}

func (que *Queue) Post(event interface{}) {
	if cap(que.events) > 0 {
		que.events <- event
	} else {
		panic(postToInactive)
	}
}

func (window *Window) resetPropsAndCmd(props Properties) {
	window.Props = props
	window.Cmd.CloseReq = false
	window.Cmd.CloseUnc = false
	window.Cmd.Update = false
}

func (window *Window) modified(props Properties) tModification {
	var mod tModification
	if window.Props != props {
		mod.fsToggle = bool(window.Props.Fullscreen != props.Fullscreen)
		mod.mouse = bool(window.Props.MouseX != props.MouseX || window.Props.MouseY != props.MouseY)
		mod.title = bool(window.Props.Title != props.Title)
		mod.pos = bool(window.Props.ClientX != props.ClientX || window.Props.ClientY != props.ClientY)
		mod.size = bool(window.Props.ClientWidth != props.ClientWidth || window.Props.ClientHeight != props.ClientHeight)
		mod.style = bool(window.Props.Borderless != props.Borderless || window.Props.Resizable != props.Resizable)
	}
	return mod
}

func newManager(data unsafe.Pointer, window AbstractWindow, params *Parameters) tManager {
	var mgr tManager
	if params.LogicThread {
		if params.GraphicThread {
			mgr = newManagerNoThreads(data, window, params)
		} else {
			mgr = newManagerNoThreads(data, window, params)
			//mgr = newManagerLogicThread(data, window, params)
		}
	} else {
		if params.GraphicThread {
			mgr = newManagerNoThreads(data, window, params)
		} else {
			mgr = newManagerNoThreads(data, window, params)
		}
	}
	return mgr
}

func newManagerNoThreads(data unsafe.Pointer, window AbstractWindow, params *Parameters) tManager {
	mgr := new(tManagerNoThreads)
	mgr.initBase(data, window, params)
	return mgr
}

func newManagerLogicThread(data unsafe.Pointer, window AbstractWindow, params *Parameters) tManager {
	mgr := new(tManagerLogicThread)
	mgr.initBase(data, window, params)
	return mgr
}

func (mgr *tManagerBase) initBase(data unsafe.Pointer, window AbstractWindow, params *Parameters) {
	mgr.data = data
	mgr.wndBase = window.baseStruct()
	mgr.wndAbst = window
	mgr.props.Title = params.Title
	mgr.autoUpdate = params.AutoUpdate
}
*/