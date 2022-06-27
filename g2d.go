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
	timepkg "time"
	"unsafe"
)

const (
	notInitialized    = "g2d not initialized"
	alreadyProcessing = "already processing events"
	messageFailed = "message post failed"
)

var (
	Err         error
	initialized bool
	processing  bool
	cb          tCallback
	timeStart   timepkg.Time
)

type Parameters struct {
	ClientX, ClientY                  int
	ClientWidth, ClientHeight         int
	ClientWidthMin, ClientHeightMin   int
	ClientWidthMax, ClientHeightMax   int
	MouseLocked, Borderless, Dragable bool
	Resizable, Fullscreen, Centered   bool
	LogicThread, GraphicThread bool
	Title                             string
}

type Properties struct {
	ClientX, ClientY                  int
	ClientWidth, ClientHeight         int
	ClientWidthMin, ClientHeightMin   int
	ClientWidthMax, ClientHeightMax   int
	MouseLocked, Borderless, Dragable bool
	Resizable, Fullscreen   bool
	Title   string
}

type Command struct {
	CloseReq bool
	CloseUnc bool
}

type Time struct {
	NanosUpdate int64
	NanosEvent int64
}

type Window struct {
	Props Properties
	Cmd Command
	Time Time
}

type AbstractWindow interface {
	Config(params *Parameters) error
	Create() error
	Show() error
	KeyDown(key int, repeated uint) error
	KeyUp(key int) error
	Close() (bool, error)
	Destroy()
	baseStruct() *Window
	updatePropsResetCmd(props Properties)
	propsAndCmd() (Properties, Command)
}

type tManager interface {
	onCreate(nanos int64)
	onShow(nanos int64)
	onKeyDown(key int, repeated uint, nanos int64)
	onKeyUp(key int, nanos int64)
	onClose(nanos int64)
	onDestroy(nanos int64)
	destroy()
}

type tManagerBase struct {
	data    unsafe.Pointer
	wndAbst AbstractWindow
	wndBase *Window
	props Properties
}

type tManagerNoThreads struct {
	tManagerBase
	cmd  Command
}

type tManagerLogicThread struct {
	tManagerBase
}

type tManagerGraphicThread struct {
}

type tManagerLogicGraphicThread struct {
}

// tCallback holds objects identified by ids.
type tCallback struct {
	mgrs   []tManager
	unused []int
}

func newParameters() *Parameters {
	params := new(Parameters)
	params.ClientX = 50
	params.ClientY = 50
	params.ClientWidth = 640
	params.ClientHeight = 480
	params.ClientWidthMin = 0
	params.ClientHeightMin = 0
	params.ClientWidthMax = 99999
	params.ClientHeightMax = 99999
	params.MouseLocked = false
	params.Borderless = false
	params.Dragable = false
	params.Resizable = true
	params.Fullscreen = false
	params.Centered = true
	params.LogicThread = true
	params.Title = "g2d - 0.1.0"
	return params
}

func (t *Time) NanosNow() int64 {
	return time()
}

func (window *Window) Config(params *Parameters) error {
	return nil
}

func (window *Window) Create() error {
	return nil
}

func (window *Window) Show() error {
	return nil
}

func (window *Window) KeyDown(key int, repeated uint) error {
	return nil
}

func (window *Window) KeyUp(key int) error {
	return nil
}

func (window *Window) Close() (bool, error) {
	return true, nil
}

func (window *Window) Destroy() {
}

func (window *Window) baseStruct() *Window {
	return window
}

func (window *Window) updatePropsResetCmd(props Properties) {
	window.Props = props
	window.Cmd.CloseReq = false
	window.Cmd.CloseUnc = false
}

func (window *Window) propsAndCmd() (Properties, Command) {
	return window.Props, window.Cmd
}

// Register returns a new id number for mgr. It will not be garbage collected until
// Unregister is called with this id.
func (cb *tCallback) Register(mgr tManager) int {
	if len(cb.unused) == 0 {
		cb.mgrs = append(cb.mgrs, mgr)
		return len(cb.mgrs) - 1
	}
	indexLast := len(cb.unused) - 1
	indexObj := cb.unused[indexLast]
	cb.unused = cb.unused[:indexLast]
	cb.mgrs[indexObj] = mgr
	return indexObj
}

// Unregister makes mgr no more identified by id.
// This object may be garbage collected, now.
func (cb *tCallback) Unregister(id int) {
	cb.mgrs[id] = nil
	cb.unused = append(cb.unused, id)
}

// UnregisterAll makes all regiestered mgrs no more identified by id.
// These objects may be garbage collected, now.
func (cb *tCallback) UnregisterAll() {
	for i := 0; i < len(cb.mgrs) && cb.mgrs[i] != nil; i++ {
		cb.Unregister(i)
	}
}

func newManager(window AbstractWindow, params *Parameters, data unsafe.Pointer) tManager {
	var mgr tManager
	if params.LogicThread {
		if params.GraphicThread {
			mgr = newManagerNoThreads(window, params, data)
		} else {
			mgr = newManagerNoThreads(window, params, data)
			//mgr = newManagerLogicThread(window, params, data)
		}
	} else {
		if params.GraphicThread {
			mgr = newManagerNoThreads(window, params, data)
		} else {
			mgr = newManagerNoThreads(window, params, data)
		}
	}
	return mgr
}

func newManagerNoThreads(window AbstractWindow, params *Parameters, data unsafe.Pointer) tManager {
	mgr := new(tManagerNoThreads)
	mgr.initBase(window, params, data)
	return mgr
}

/*
func newManagerLogicThread(window AbstractWindow, params *Parameters, data unsafe.Pointer) tManager {
	mgr := new(tManagerLogicThread)
	mgr.initBase(window, params, data)
	return mgr
}
*/

func (mgr *tManagerBase) initBase(window AbstractWindow, params *Parameters, data unsafe.Pointer) {
	mgr.wndBase = window.baseStruct()
	mgr.wndAbst = window
	mgr.props.Title = params.Title
	mgr.data = data
}

func time() int64 {
	timeNow := timepkg.Now()
	d := timeNow.Sub(timeStart)
	return d.Nanoseconds()
}

// toCInt converts bool value to C int value.
func toCInt(b bool) C.int {
	if b {
		return C.int(1)
	}
	return C.int(0)
}
