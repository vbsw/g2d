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
	"unsafe"
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
	closeType
	destroyType
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
	wnds []*tWindow
	wndNextId []int
	requests []tRequest
)

// Window is the window event listener.
type Window interface {
	OnConfig(config *Configuration) error
	OnCreate() error
	OnShow() error
	OnResize() error
	OnKeyDown(keyCode int, repeated uint) error
	OnKeyUp(keyCode int) error
	OnTextureLoaded(textureId int) error
	OnUpdate() error
	OnClose() (bool, error)
	OnDestroy() error
}

// WindowImpl is the default implementation of Window.
type WindowImpl struct {
}

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

type tWindow struct {
	eventsChan       chan *tLogicEvent
	abst Window
	data unsafe.Pointer
	id, state int
	autoUpdate bool
}

type tLogicEvent struct {
	typeId   int
	valA     int
	valB     float32
	repeated uint
	time    int
	props    Properties
	obj      interface{}
}

type tRequest interface {
	process()
}

type tCreateWindowRequest struct {
	config *Configuration
	wndId int
}

type tShowWindowRequest struct {
	wndId int
}

type tDestroyWindowRequest struct {
	wndId int
}

type tSetPropertiesRequest struct {
	props Properties
	modPos, modSize, modStyle bool
	modFullscreen, modMouse, modTitle bool
	wndId int
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

func (config *Configuration) boolsToCInt() (C.int, C.int, C.int, C.int, C.int, C.int) {
	var c, l, b, d, r, f C.int
	if config.Centered {
		c = 1
	}
	if config.MouseLocked {
		l = 1
	}
	if config.Borderless {
		b = 1
	}
	if config.Dragable {
		d = 1
	}
	if config.Resizable {
		r = 1
	}
	if config.Fullscreen {
		f = 1
	}
	return c, l, b, d, r, f
}

func (props *Properties) boolsToCInt() (C.int, C.int, C.int, C.int) {
	var l, b, d, r C.int
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
	return l, b, d, r
}

func newWindow(abst Window) *tWindow {
	wnd := new(tWindow)
	wnd.eventsChan = make(chan *tLogicEvent, 1000)
	wnd.abst = abst
	wnd.id = registerWnd(wnd)
	wnd.state = configState
	return wnd
}

func (wnd *tWindow) nextEvent() *tLogicEvent {
	var event *tLogicEvent
	event = <-wnd.eventsChan
/*
	if wnd.state > configState && wnd.state < closingState && wnd.autoUpdate {
		select {
		case msg := <-wnd.msgs:
			event = msg
		default:
			wnd.update = false
			event = &tLogicEvent{typeId: updateType, nanos: time.Nanos()}
		}
	} else {
		event = <-wnd.msgs
		if event == nil && wnd.update {
			wnd.update = false
			event = &tLogicEvent{typeId: updateType, nanos: time.Nanos()}
		}
	}
	if wnd.state == closingState && event.typeId != quitType {
		event = nil
	}
*/
	return event
}

func (wnd *tWindow) logicThread() {
	for wnd.state != quitState {
		event := wnd.nextEvent()
		if event != nil {
			switch event.typeId {
			case configType:
				wnd.onConfig()
			case createType:
				wnd.onCreate()
			case closeType:
				wnd.onClose()
			case destroyType:
				wnd.onDestroy()
			}
		}
	}
/*
	wnd := abst.impl()
	for wnd.state != quitState {
		msg := wnd.nextEvent()
		if msg != nil {
			wnd.Time.Curr = msg.nanos
			switch msg.typeId {
			case configType:
				onConfig(abst, wnd)
			case createType:
				onCreate(abst, wnd)
			case showType:
				onShow(abst, wnd)
			case wndMoveType:
				wnd.updateProps(msg)
				onWindowMoved(abst, wnd)
			case wndResizeType:
				wnd.updateProps(msg)
				onWindowResize(abst, wnd)
			case keyDownType:
				onKeyDown(abst, wnd, msg.valA, msg.repeated)
			case keyUpType:
				onKeyUp(abst, wnd, msg.valA)
			case msMoveType:
				wnd.updateProps(msg)
				onMouseMove(abst, wnd)
			case buttonDownType:
				onButtonDown(abst, wnd, msg.valA, msg.repeated != 0)
			case buttonUpType:
				onButtonUp(abst, wnd, msg.valA, msg.repeated != 0)
			case wheelType:
				onWheel(abst, wnd, msg.valB)
			case minimizeType:
				onWindowMinimize(abst, wnd)
			case restoreType:
				onWindowRestore(abst, wnd)
			case textureType:
				onTextureLoaded(abst, wnd, msg.valA)
			case updateType:
				onUpdate(abst, wnd)
			case closeType:
				onQuitReq(abst, wnd)
			case quitType:
				onDestroy(abst, wnd)
			}
		}
	}
	wnd.gfxImpl.destroy()
*/
}

func (wnd *tWindow) onConfig() {
	config := newConfiguration()
	err := wnd.abst.OnConfig(config)
	wnd.autoUpdate = config.AutoUpdate
	if err == nil {
		postReqest(&tCreateWindowRequest{wndId: wnd.id, config: config})
	}
/*
	if err == nil {
		toMainLoop.postMsg(&tCreateWindowRequest{abst: abst, config: config})
	} else {
		onLogicError(abst, wnd, 4999, err)
	}
*/
}

func (wnd *tWindow) onCreate() {
	err := wnd.abst.OnCreate()
	if err == nil {
		postReqest(&tShowWindowRequest{wndId: wnd.id})
	}
/*
	if err == nil {
		toMainLoop.postMsg(&tCreateWindowRequest{abst: abst, config: config})
	} else {
		onLogicError(abst, wnd, 4999, err)
	}
*/
}

func (wnd *tWindow) onClose() {
	quit, err := wnd.abst.OnClose()
	if err == nil {
		if quit {
			postReqest(&tDestroyWindowRequest{wndId: wnd.id})
		}
	}
/*
	if err == nil {
		toMainLoop.postMsg(&tCreateWindowRequest{abst: abst, config: config})
	} else {
		onLogicError(abst, wnd, 4999, err)
	}
*/
}

func (wnd *tWindow) onDestroy() {
/*
	if wnd.wgt.Gfx.running {
		wnd.wgt.Gfx.msgs <- &tGMessage{typeId: quitType}
		<- wnd.wgt.Gfx.quitted
	}
*/
	wnd.abst.OnDestroy()
/*
	wnd.wgt.quitted <- true
	wnd.state = quitState
	if err != nil {
		setErrorSynced(err)
	}
*/
}

func (wndImpl *WindowImpl) OnConfig(config *Configuration) error {
	return nil
}

func (wndImpl *WindowImpl) OnCreate() error {
	return nil
}

func (wndImpl *WindowImpl) OnShow() error {
	return nil
}

func (wndImpl *WindowImpl) OnResize() error {
	return nil
}

func (wndImpl *WindowImpl) OnKeyDown(keyCode int, repeated uint) error {
	return nil
}

func (wndImpl *WindowImpl) OnKeyUp(keyCode int) error {
	return nil
}

func (wndImpl *WindowImpl) OnTextureLoaded(textureId int) error {
	return nil
}

func (wndImpl *WindowImpl) OnUpdate() error {
	return nil
}

func (wndImpl *WindowImpl) OnClose() (bool, error) {
	return true, nil
}

func (wndImpl *WindowImpl) OnDestroy() error {
	return nil
}

func registerWnd(wnd *tWindow) int {
	var id int
	if len(wndNextId) == 0 {
		wnds = append(wnds, wnd)
		id = len(wnds) - 1
	} else {
		idLast := len(wndNextId) - 1
		id = wndNextId[idLast]
		wndNextId = wndNextId[:idLast]
		wnds[id] = wnd
	}
	return id
}

func unregisterWnd(id int) *tWindow {
	wnd := wnds[id]
	wnds[id] = nil
	wndNextId = append(wndNextId, id)
	return wnd
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