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
	"time"
	"unsafe"
)

const (
	mustNotBeNil   = "main window must not be nil"
	notInitialized = "g2d not initialized"
	alreadyRunning = "g2d main loop already running"
)

const (
	configState = iota
	showingState
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
	swapIntervType
	imageType
	textureType
	minimizeType
	restoreType
	focusType
	customType
	refreshType
)

var (
	MaxTexSize, MaxTexUnits int
	Err                     error
)

var (
	initialized, initFailed bool
	running, quitting       bool
	processingRequests      bool
	mutex                   sync.Mutex
	wnds                    []*tWindow
	wndNextId               []int
	requests                []tRequest
	appTime                 tAppTime
)

type Window interface {
	OnConfig(config *Configuration) error
	OnCreate() error
	OnShow() error
	OnResize() error
	OnMove() error
	OnKeyDown(keyCode int, repeated uint) error
	OnKeyUp(keyCode int) error
	OnMouseMove() error
	OnButtonDown(buttonCode int, doubleClicked bool) error
	OnButtonUp(buttonCode int, doubleClicked bool) error
	OnWheel(rotation float32) error
	OnCustom(obj interface{}) error
	OnTextureLoaded(textureId int) error
	OnUpdate() error
	OnClose() (bool, error)
	OnDestroy() error
	OnMinimize() error
	OnRestore() error
	OnFocus(focus bool) error
	Custom(obj interface{})
	Update()
	Close()
	Quit()
	Show(window Window)
	impl() *WindowImpl
}

// WindowImpl is the implementation of Window.
type WindowImpl struct {
	Props Properties
	Stats Stats
	Gfx   Graphics
	id    int
}

// Configuration for the starting window.
type Configuration struct {
	ClientX, ClientY                  int
	ClientWidth, ClientHeight         int
	ClientWidthMin, ClientHeightMin   int
	ClientWidthMax, ClientHeightMax   int
	MouseLocked, Borderless, Dragable bool
	Resizable, Fullscreen, Centered   bool
	Title                             string
}

// Properties of the current window.
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

// Stats holds time in milliseconds.
type Stats struct {
	AppTime, DeltaTime int
	lastUpdate         int
	UPS, FPS           int
	ups, lastUPSTime   int
	fps, lastFPSTime   int
}

// Graphics functions to draw in window.
type Graphics struct {
	eventsChan chan *tGraphicsEvent
	mutex      sync.Mutex
	state      int
	updating   bool
}

type tWindow struct {
	eventsChan      chan *tLogicEvent
	abst            Window
	impl            *WindowImpl
	data            unsafe.Pointer
	id, state, time int
	update          bool
}

type tLogicEvent struct {
	typeId   int
	valA     int
	valB     float32
	repeated uint
	time     int
	props    Properties
	obj      interface{}
}

type tGraphicsEvent struct {
	typeId int
	valA   int
	valB   int
	valC   interface{}
	err    error
}

type tRequest interface {
	process()
}

type tConfigWindowRequest struct {
	window Window
}

type tCreateWindowRequest struct {
	config *Configuration
	wndId  int
}

type tShowWindowRequest struct {
	wndId int
}

type tCloseWindowRequest struct {
	wndId int
}

type tDestroyWindowRequest struct {
	wndId int
}

type tCustomRequest struct {
	obj   interface{}
	wndId int
}

type tUpdateRequest struct {
	wndId int
}

type tSetPropertiesRequest struct {
	props                             Properties
	modPosSize, modStyle              bool
	modFullscreen, modMouse, modTitle bool
	wndId                             int
}

type tAppTime struct {
	start time.Time
}

func (stats *Stats) updateUPS() {
	diff := stats.AppTime - stats.lastUPSTime
	if diff < 1000 {
		stats.ups++
	} else if diff < 2000 {
		stats.UPS = stats.ups
		stats.ups = 1
		stats.lastUPSTime += 1000
	} else {
		stats.UPS = 0
		stats.ups = 1
		stats.lastUPSTime += diff % 1000 * 1000
	}
}

func (stats *Stats) updateFPS() {
	diff := stats.AppTime - stats.lastFPSTime
	if diff < 1000 {
		stats.fps++
	} else if diff < 2000 {
		stats.FPS = stats.fps
		stats.fps = 1
		stats.lastFPSTime += 1000
	} else {
		stats.FPS = 0
		stats.fps = 1
		stats.lastFPSTime += diff % 1000 * 1000
	}
}

func (gfx *Graphics) Update() {
	gfx.mutex.Lock()
	if !gfx.updating {
		gfx.updating = true
		gfx.eventsChan <- &tGraphicsEvent{typeId: updateType}
	}
	gfx.mutex.Unlock()
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

func (props *Properties) boolsToCInt() (C.int, C.int, C.int, C.int, C.int) {
	var l, b, d, r, f C.int
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
	if props.Fullscreen {
		f = 1
	}
	return l, b, d, r, f
}

func (props *Properties) copyTo(target *Properties) {
	titleTmp := target.Title
	*target = *props
	target.Title = titleTmp
}

func (props *Properties) compare(target *Properties) *tSetPropertiesRequest {
	var req *tSetPropertiesRequest
	if *props != *target {
		req = new(tSetPropertiesRequest)
		req.props = *target
		req.modPosSize = bool(props.ClientX != target.ClientX || props.ClientY != target.ClientY)
		req.modPosSize = bool(req.modPosSize || props.ClientWidth != target.ClientWidth || props.ClientHeight != target.ClientHeight)
		req.modStyle = bool(props.ClientWidthMin != target.ClientWidthMin || props.ClientHeightMin != target.ClientHeightMin)
		req.modStyle = bool(req.modStyle || props.ClientWidthMax != target.ClientWidthMax || props.ClientHeightMax != target.ClientHeightMax)
		req.modStyle = bool(req.modStyle || props.MouseLocked != target.MouseLocked || props.Borderless != target.Borderless)
		req.modStyle = bool(req.modStyle || props.Dragable != target.Dragable || props.Resizable != target.Resizable)
		req.modFullscreen = bool(props.Fullscreen != target.Fullscreen)
		req.modMouse = bool(props.MouseX != target.MouseX || props.MouseY != target.MouseY)
		req.modTitle = bool(props.Title != target.Title)
	}
	return req
}

func newWindow(abst Window) *tWindow {
	wnd := new(tWindow)
	wnd.id = registerWnd(wnd)
	wnd.eventsChan = make(chan *tLogicEvent, 1000)
	wnd.abst = abst
	wnd.impl = abst.impl()
	wnd.state = configState
	wnd.impl.id = wnd.id
	wnd.impl.Gfx.eventsChan = make(chan *tGraphicsEvent, 1000)
	return wnd
}

func (wnd *tWindow) logicThread() {
	for wnd.state != quitState {
		event := <-wnd.eventsChan
		if event != nil {
			if event.typeId != destroyType {
				event.props.copyTo(&wnd.impl.Props)
			}
			wnd.impl.Stats.AppTime = event.time
			switch event.typeId {
			case configType:
				wnd.onConfig()
			case createType:
				wnd.onCreate()
			case showType:
				wnd.onShow()
			case wndMoveType:
				wnd.onMove()
			case wndResizeType:
				wnd.onResize()
			case keyDownType:
				wnd.onKeyDown(event.valA, event.repeated)
			case keyUpType:
				wnd.onKeyUp(event.valA)
			case msMoveType:
				wnd.onMouseMove()
			case buttonDownType:
				wnd.onButtonDown(event.valA, event.repeated != 0)
			case buttonUpType:
				wnd.onButtonUp(event.valA, event.repeated != 0)
			case wheelType:
				wnd.onWheel(event.valB)
			case updateType:
				wnd.onUpdate()
			case closeType:
				wnd.onClose()
			case destroyType:
				wnd.onDestroy()
			case minimizeType:
				wnd.onMinimize()
			case restoreType:
				wnd.onRestore()
			case focusType:
				wnd.onFocus(event.valA != 0)
			case customType:
				wnd.onCustom(event.obj)
			case refreshType:
				wnd.impl.Stats.updateFPS()
			}
		}
	}
}

func (wnd *tWindow) onConfig() {
	config := newConfiguration()
	err := wnd.abst.OnConfig(config)
	wnd.impl.Props.Title = config.Title
	if err == nil {
		postRequest(&tCreateWindowRequest{wndId: wnd.id, config: config})
	}
}

func (wnd *tWindow) onCreate() {
	err := wnd.abst.OnCreate()
	if err == nil {
		go wnd.graphicsThread()
		wnd.impl.Gfx.eventsChan <- &tGraphicsEvent{typeId: updateType}
		postRequest(&tShowWindowRequest{wndId: wnd.id})
	}
}

func (wnd *tWindow) onShow() {
	props := wnd.impl.Props
	wnd.state = showingState
	wnd.impl.Stats.lastUPSTime = wnd.impl.Stats.AppTime
	err := wnd.abst.OnShow()
	if err == nil {
		setPropsReq := props.compare(&wnd.impl.Props)
		if setPropsReq != nil {
			setPropsReq.wndId = wnd.id
			postRequest(setPropsReq)
		}
	}
}

func (wnd *tWindow) onMove() {
	props := wnd.impl.Props
	err := wnd.abst.OnMove()
	if err == nil {
		setPropsReq := props.compare(&wnd.impl.Props)
		if setPropsReq != nil {
			setPropsReq.wndId = wnd.id
			postRequest(setPropsReq)
		}
	}
}

func (wnd *tWindow) onResize() {
	props := wnd.impl.Props
	err := wnd.abst.OnResize()
	if err == nil {
		setPropsReq := props.compare(&wnd.impl.Props)
		if setPropsReq != nil {
			setPropsReq.wndId = wnd.id
			postRequest(setPropsReq)
		}
	}
}

func (wnd *tWindow) onKeyDown(keyCode int, repeated uint) {
	props := wnd.impl.Props
	err := wnd.abst.OnKeyDown(keyCode, repeated)
	if err == nil {
		setPropsReq := props.compare(&wnd.impl.Props)
		if setPropsReq != nil {
			setPropsReq.wndId = wnd.id
			postRequest(setPropsReq)
		}
	}
}

func (wnd *tWindow) onKeyUp(keyCode int) {
	props := wnd.impl.Props
	err := wnd.abst.OnKeyUp(keyCode)
	if err == nil {
		setPropsReq := props.compare(&wnd.impl.Props)
		if setPropsReq != nil {
			setPropsReq.wndId = wnd.id
			postRequest(setPropsReq)
		}
	}
}

func (wnd *tWindow) onMouseMove() {
	props := wnd.impl.Props
	err := wnd.abst.OnMouseMove()
	if err == nil {
		setPropsReq := props.compare(&wnd.impl.Props)
		if setPropsReq != nil {
			setPropsReq.wndId = wnd.id
			postRequest(setPropsReq)
		}
	}
}

func (wnd *tWindow) onButtonDown(buttonCode int, doubleClicked bool) {
	props := wnd.impl.Props
	err := wnd.abst.OnButtonDown(buttonCode, doubleClicked)
	if err == nil {
		setPropsReq := props.compare(&wnd.impl.Props)
		if setPropsReq != nil {
			setPropsReq.wndId = wnd.id
			postRequest(setPropsReq)
		}
	}
}

func (wnd *tWindow) onButtonUp(buttonCode int, doubleClicked bool) {
	props := wnd.impl.Props
	err := wnd.abst.OnButtonUp(buttonCode, doubleClicked)
	if err == nil {
		setPropsReq := props.compare(&wnd.impl.Props)
		if setPropsReq != nil {
			setPropsReq.wndId = wnd.id
			postRequest(setPropsReq)
		}
	}
}

func (wnd *tWindow) onWheel(rotation float32) {
	props := wnd.impl.Props
	err := wnd.abst.OnWheel(rotation)
	if err == nil {
		setPropsReq := props.compare(&wnd.impl.Props)
		if setPropsReq != nil {
			setPropsReq.wndId = wnd.id
			postRequest(setPropsReq)
		}
	}
}

func (wnd *tWindow) onUpdate() {
	wnd.update = false
	wnd.impl.Stats.DeltaTime = wnd.impl.Stats.AppTime - wnd.impl.Stats.lastUpdate
	wnd.impl.Stats.lastUpdate = wnd.impl.Stats.AppTime
	wnd.impl.Stats.updateUPS()
	props := wnd.impl.Props
	err := wnd.abst.OnUpdate()
	if err == nil {
		setPropsReq := props.compare(&wnd.impl.Props)
		if setPropsReq != nil {
			setPropsReq.wndId = wnd.id
			postRequest(setPropsReq)
		}
	}
}

func (wnd *tWindow) onClose() {
	props := wnd.impl.Props
	quit, err := wnd.abst.OnClose()
	if err == nil {
		if quit {
			postRequest(&tDestroyWindowRequest{wndId: wnd.id})
		} else {
			setPropsReq := props.compare(&wnd.impl.Props)
			if setPropsReq != nil {
				setPropsReq.wndId = wnd.id
				postRequest(setPropsReq)
			}
		}
	}
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

func (wnd *tWindow) onCustom(obj interface{}) {
	props := wnd.impl.Props
	err := wnd.abst.OnCustom(obj)
	if err == nil {
		setPropsReq := props.compare(&wnd.impl.Props)
		if setPropsReq != nil {
			setPropsReq.wndId = wnd.id
			postRequest(setPropsReq)
		}
	}
}

func (wnd *tWindow) onMinimize() {
	props := wnd.impl.Props
	err := wnd.abst.OnMinimize()
	if err == nil {
		setPropsReq := props.compare(&wnd.impl.Props)
		if setPropsReq != nil {
			setPropsReq.wndId = wnd.id
			postRequest(setPropsReq)
		}
	}
}

func (wnd *tWindow) onRestore() {
	props := wnd.impl.Props
	err := wnd.abst.OnRestore()
	if err == nil {
		setPropsReq := props.compare(&wnd.impl.Props)
		if setPropsReq != nil {
			setPropsReq.wndId = wnd.id
			postRequest(setPropsReq)
		}
	}
}

func (wnd *tWindow) onFocus(focus bool) {
	if wnd.state > configState {
		props := wnd.impl.Props
		err := wnd.abst.OnFocus(focus)
		if err == nil {
			setPropsReq := props.compare(&wnd.impl.Props)
			if setPropsReq != nil {
				setPropsReq.wndId = wnd.id
				postRequest(setPropsReq)
			}
		}
	}
}

func (wnd *WindowImpl) OnConfig(config *Configuration) error {
	return nil
}

func (wnd *WindowImpl) OnCreate() error {
	return nil
}

func (wnd *WindowImpl) OnShow() error {
	wnd.Update()
	return nil
}

func (wnd *WindowImpl) OnResize() error {
	return nil
}

func (wnd *WindowImpl) OnMove() error {
	return nil
}

func (wnd *WindowImpl) OnKeyDown(keyCode int, repeated uint) error {
	return nil
}

func (wnd *WindowImpl) OnKeyUp(keyCode int) error {
	return nil
}

func (wnd *WindowImpl) OnMouseMove() error {
	return nil
}

func (wnd *WindowImpl) OnButtonDown(buttonCode int, doubleClicked bool) error {
	return nil
}

func (wnd *WindowImpl) OnButtonUp(buttonCode int, doubleClicked bool) error {
	return nil
}

func (wnd *WindowImpl) OnWheel(rotation float32) error {
	return nil
}

func (wnd *WindowImpl) OnCustom(obj interface{}) error {
	return nil
}

func (wnd *WindowImpl) OnTextureLoaded(textureId int) error {
	return nil
}

func (wnd *WindowImpl) OnUpdate() error {
	return nil
}

func (wnd *WindowImpl) OnClose() (bool, error) {
	return true, nil
}

func (wnd *WindowImpl) OnDestroy() error {
	return nil
}

func (wnd *WindowImpl) OnMinimize() error {
	return nil
}

func (wnd *WindowImpl) OnRestore() error {
	return nil
}

func (wnd *WindowImpl) OnFocus(focus bool) error {
	return nil
}

func (wnd *WindowImpl) Close() {
	postRequest(&tCloseWindowRequest{wndId: wnd.id})
}

func (wnd *WindowImpl) Update() {
	mutex.Lock()
	wndWrapper := wnds[wnd.id]
	if !wndWrapper.update {
		wndWrapper.update = true
		event := &tLogicEvent{typeId: updateType, time: appTime.Millis()}
		event.props.update(wndWrapper.data)
		wndWrapper.eventsChan <- event
	}
	mutex.Unlock()
}

func (wnd *WindowImpl) Quit() {
	postRequest(&tDestroyWindowRequest{wndId: wnd.id})
}

func (wnd *WindowImpl) Custom(obj interface{}) {
	postRequest(&tCustomRequest{wndId: wnd.id, obj: obj})
}

func (wnd *WindowImpl) Show(window Window) {
	postRequest(&tConfigWindowRequest{window: window})
}

func (wnd *WindowImpl) impl() *WindowImpl {
	return wnd
}

func (t *tAppTime) Reset() {
	t.start = time.Now()
}

func (t *tAppTime) Nanos() int64 {
	now := time.Now()
	delta := now.Sub(t.start)
	return delta.Nanoseconds()
}

func (t *tAppTime) Millis() int {
	now := time.Now()
	delta := now.Sub(t.start)
	return int(delta.Milliseconds())
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
