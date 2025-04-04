/*
 *          Copyright 2025, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

// Package g2d is a framework to create 2D graphic applications.
package g2d

import "C"
import (
	"sync"
)

const (
	configState = iota
	runningState
	closingState
	quitState
)

const (
	configType = iota
	createType
	showType
	wndMoveType
	wndResizeType
	keyDownType
	keyUpType
	msMoveType
	buttonDownType
	buttonUpType
	wheelType
	updateType
	quitReqType
	quitType
	leaveType
	refreshType
	swapIntervType
	imageType
	textureType
	minimizeType
	restoreType
)

var (
	MaxTexSize, MaxTexUnits int
	Err error
)

var (
	initialized, initFailed bool
	running, quitting       bool
	mutex                   sync.Mutex
)

type Window interface {
}

/*
const (
	notInitialized    = "g2d not initialized"
	alreadyProcessing = "already processing events"
	messageFailed     = "message post failed"
	postToInactive    = "can't post event to inactive window"
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
	LogicThread, GraphicThread        bool
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
	Title                             string
}

type Command struct {
	CloseReq bool
	CloseUnc bool
	Update   bool
}

type Time struct {
	NanosUpdate int64
	NanosEvent  int64
}

type Queue struct {
	events chan interface{}
}

type Window struct {
	Props      Properties
	Cmd        Command
	Time       Time
	Queue      Queue
	destroying bool
}

type AbstractWindow interface {
	Config(params *Parameters) error
	Create() error
	Show() error
	KeyDown(key int, repeated uint) error
	KeyUp(key int) error
	Update() error
	Close() (bool, error)
	Destroy()
	baseStruct() *Window
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

// tCallback holds objects identified by ids.
type tCallback struct {
	mgrs   []tManager
	unused []int
}

type tModification struct {
	mouse, style, title bool
	fsToggle, pos, size bool
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
	params.AutoUpdate = true
	params.Title = "g2d - 0.1.0"
	return params
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

func (window *Window) init() {
	window.Time.NanosEvent = time()
	window.baseStruct().Time.NanosUpdate = window.baseStruct().Time.NanosEvent
	window.Queue.events = make(chan interface{}, 1024 * 8)
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

func (window *Window) Update() error {
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
*/

/*
func newManagerLogicThread(data unsafe.Pointer, window AbstractWindow, params *Parameters) tManager {
	mgr := new(tManagerLogicThread)
	mgr.initBase(data, window, params)
	return mgr
}
*/

/*
func (mgr *tManagerBase) initBase(data unsafe.Pointer, window AbstractWindow, params *Parameters) {
	mgr.data = data
	mgr.wndBase = window.baseStruct()
	mgr.wndAbst = window
	mgr.props.Title = params.Title
	mgr.autoUpdate = params.AutoUpdate
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
*/