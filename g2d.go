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
	configType     = 0
	createType     = 1
	showType       = 2
	wndMoveType    = 3
	wndResizeType  = 4
	keyDownType    = 5
	keyUpType      = 6
	msMoveType     = 7
	buttonDownType = 8
	buttonUpType   = 9
	wheelType      = 10
	updateType     = 11
	closeType      = 12
	destroyType    = 13
	leaveType      = 14
	swapIntervType = 15
	imageType      = 16
	textureType    = 17
	minimizeType   = 18
	restoreType    = 19
	focusType      = 20
	customType     = 21
	refreshType    = 22
)

var (
	MaxTexSize, MaxTexUnits int
	VSyncAvailable          bool
	AVSyncAvailable         bool
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
	OnDestroy(error) error
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
	BgR, BgG, BgB float32
	VSync, AVSync bool
	eventsChan    chan *tGraphicsEvent
	quittedChan   chan bool
	mutex         sync.Mutex
	w, h          int
	read          *tGraphics
	buffer        *tGraphics
	write         *tGraphics
	bufferReady   bool
	updating      bool
	running       bool
}

type RectanglesLayer struct {
	entities     []*Rectangle
	entityNextId []int
	count        int
	Enabled      bool
	buffer       []C.float
}

type Rectangle struct {
	id                  int
	X, Y, Width, Height float32
	R, G, B, A          float32
	Enabled             bool
}

type tWindow struct {
	eventsChan      chan *tLogicEvent
	quittedChan     chan bool
	abst            Window
	impl            *WindowImpl
	data            unsafe.Pointer
	title           string
	id, state, time int
	update          bool
}

type tGraphics struct {
	w, h, si C.int
	r, g, b  C.float
	layers   []tLayer
}

type tLayer interface {
	copyTo(tLayer)
	getProcessing([]tLayer) ([]tLayer, []C.float, int, unsafe.Pointer)
}

type tLogicEvent struct {
	typeId   int
	valA     int
	valB     float32
	repeated uint
	time     int
	props    Properties
	obj      interface{}
	err      error
}

type tGraphicsEvent struct {
	typeId int
	valA   int
	valB   int
	valC   interface{}
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

type tErrorRequest struct {
	err error
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

func (gfx *Graphics) NewRectanglesLayer() *RectanglesLayer {
	layer := new(RectanglesLayer)
	layer.Enabled = true
	gfx.write.layers = append(gfx.write.layers, layer)
	return layer
}

func (gfx *Graphics) getReadBuffer() *tGraphics {
	if gfx.bufferReady {
		read := gfx.buffer
		gfx.buffer = gfx.read
		gfx.read = read
		gfx.bufferReady = false
	}
	return gfx.read
}

func (gfx *Graphics) postRefresh() {
	var swapInt int
	gfx.mutex.Lock()
	if gfx.VSync {
		swapInt = 1
	} else if gfx.AVSync {
		swapInt = -1
	}
	gfx.write.copyTo(gfx.buffer, gfx.w, gfx.h, swapInt, gfx.BgR, gfx.BgG, gfx.BgB)
	gfx.bufferReady = true
	if !gfx.updating {
		gfx.updating = true
		gfx.eventsChan <- &tGraphicsEvent{typeId: refreshType}
	}
	gfx.mutex.Unlock()
}

func (layer *RectanglesLayer) NewEntity() *Rectangle {
	var entity *Rectangle
	if len(layer.entityNextId) == 0 {
		entity = new(Rectangle)
		entity.id = len(layer.entities)
		layer.entities = append(layer.entities, entity)
	} else {
		idLast := len(layer.entityNextId) - 1
		id := layer.entityNextId[idLast]
		layer.entityNextId = layer.entityNextId[:idLast]
		entity = layer.entities[id]
	}
	entity.Enabled = true
	layer.count++
	return entity
}

func (layer *RectanglesLayer) Release(r *Rectangle) *Rectangle {
	r.Enabled = false
	layer.entityNextId = append(layer.entityNextId, r.id)
	layer.count--
	return nil
}

func (layer *RectanglesLayer) copyTo(dest tLayer) {
	destLayer := dest.(*RectanglesLayer)
	destLayer.Enabled = layer.Enabled
	if layer.Enabled {
		destLayer.count = layer.count
		destLen := len(destLayer.entities)
		if cap(destLayer.entities) < len(layer.entities) {
			entitiesNew := make([]*Rectangle, len(layer.entities), cap(layer.entities))
			copy(entitiesNew, destLayer.entities)
			destLayer.entities = entitiesNew
		} else if destLen < len(layer.entities) {
			destLayer.entities = destLayer.entities[:len(layer.entities)]
		}
		for i := destLen; i < len(layer.entities); i++ {
			destLayer.entities[i] = new(Rectangle)
		}
		for i, entity := range layer.entities {
			*destLayer.entities[i] = *entity
		}
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
	wnd.quittedChan = make(chan bool, 1)
	wnd.abst = abst
	wnd.state = configState
	wnd.impl = abst.impl()
	wnd.impl.id = wnd.id
	wnd.impl.Gfx.VSync = VSyncAvailable
	wnd.impl.Gfx.eventsChan = make(chan *tGraphicsEvent, 1000)
	wnd.impl.Gfx.quittedChan = make(chan bool, 1)
	wnd.impl.Gfx.read = new(tGraphics)
	wnd.impl.Gfx.buffer = new(tGraphics)
	wnd.impl.Gfx.write = new(tGraphics)
	return wnd
}

func (wnd *tWindow) logicThread() {
	for wnd.state != quitState {
		event := <-wnd.eventsChan
		if event != nil {
			if event.typeId != destroyType && event.typeId != leaveType {
				wnd.impl.Props = event.props
			}
			wnd.impl.Stats.AppTime = event.time
			switch event.typeId {
			case configType:
				wnd.onConfig()
			case createType:
				wnd.onCreate()
			case showType:
				wnd.onShow()
			case destroyType:
				wnd.onDestroy(event.err)
			case leaveType:
				wnd.state = quitState
				wnd.impl.Gfx.quittedChan <- true
			default:
				if wnd.state == showingState {
					switch event.typeId {
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
	}
	<-wnd.impl.Gfx.quittedChan
	wnd.quittedChan <- true
}

func (wnd *tWindow) onConfig() {
	config := newConfiguration()
	err := wnd.abst.OnConfig(config)
	if err == nil {
		postRequest(&tCreateWindowRequest{wndId: wnd.id, config: config})
	} else {
		wnd.state = closingState
		postRequest(&tErrorRequest{err: err})
	}
}

func (wnd *tWindow) onCreate() {
	go wnd.graphicsThread()
	err := wnd.abst.OnCreate()
	if err == nil {
		wnd.impl.Gfx.w, wnd.impl.Gfx.h = wnd.impl.Props.ClientWidth, wnd.impl.Props.ClientHeight
		wnd.impl.Gfx.postRefresh()
		postRequest(&tShowWindowRequest{wndId: wnd.id})
	} else {
		wnd.state = closingState
		postRequest(&tErrorRequest{err: err})
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
	} else {
		wnd.state = closingState
		postRequest(&tErrorRequest{err: err})
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
	} else {
		wnd.state = closingState
		postRequest(&tErrorRequest{err: err})
	}
}

func (wnd *tWindow) onResize() {
	props := wnd.impl.Props
	wnd.impl.Gfx.w, wnd.impl.Gfx.h = props.ClientWidth, props.ClientHeight
	err := wnd.abst.OnResize()
	if err == nil {
		setPropsReq := props.compare(&wnd.impl.Props)
		if setPropsReq != nil {
			setPropsReq.wndId = wnd.id
			postRequest(setPropsReq)
		}
	} else {
		wnd.state = closingState
		postRequest(&tErrorRequest{err: err})
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
	} else {
		wnd.state = closingState
		postRequest(&tErrorRequest{err: err})
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
	} else {
		wnd.state = closingState
		postRequest(&tErrorRequest{err: err})
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
	} else {
		wnd.state = closingState
		postRequest(&tErrorRequest{err: err})
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
	} else {
		wnd.state = closingState
		postRequest(&tErrorRequest{err: err})
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
	} else {
		wnd.state = closingState
		postRequest(&tErrorRequest{err: err})
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
	} else {
		wnd.state = closingState
		postRequest(&tErrorRequest{err: err})
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
		wnd.impl.Gfx.postRefresh()
		setPropsReq := props.compare(&wnd.impl.Props)
		if setPropsReq != nil {
			setPropsReq.wndId = wnd.id
			postRequest(setPropsReq)
		}
	} else {
		wnd.state = closingState
		postRequest(&tErrorRequest{err: err})
	}
}

func (wnd *tWindow) onClose() {
	props := wnd.impl.Props
	quit, err := wnd.abst.OnClose()
	if err == nil {
		if quit {
			wnd.state = closingState
			postRequest(&tDestroyWindowRequest{wndId: wnd.id})
		} else {
			setPropsReq := props.compare(&wnd.impl.Props)
			if setPropsReq != nil {
				setPropsReq.wndId = wnd.id
				postRequest(setPropsReq)
			}
		}
	} else {
		wnd.state = closingState
		postRequest(&tErrorRequest{err: err})
	}
}

func (wnd *tWindow) onDestroy(err error) {
	wnd.state = quitState
	wnd.impl.Gfx.eventsChan <- &tGraphicsEvent{typeId: leaveType}
	err2 := wnd.abst.OnDestroy(err)
	if err2 != nil && err2 != err {
		postRequest(&tErrorRequest{err: err2})
	}
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
	} else {
		wnd.state = closingState
		postRequest(&tErrorRequest{err: err})
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
	} else {
		wnd.state = closingState
		postRequest(&tErrorRequest{err: err})
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
	} else {
		wnd.state = closingState
		postRequest(&tErrorRequest{err: err})
	}
}

func (wnd *tWindow) onFocus(focus bool) {
	props := wnd.impl.Props
	err := wnd.abst.OnFocus(focus)
	if err == nil {
		setPropsReq := props.compare(&wnd.impl.Props)
		if setPropsReq != nil {
			setPropsReq.wndId = wnd.id
			postRequest(setPropsReq)
		}
	} else {
		wnd.state = closingState
		postRequest(&tErrorRequest{err: err})
	}
}

func (gfx *tGraphics) copyTo(dest *tGraphics, w, h, si int, r, g, b float32) {
	dest.w, dest.h, dest.si = C.int(w), C.int(h), C.int(si)
	dest.r, dest.g, dest.b = C.float(r), C.float(g), C.float(b)
	for i, layer := range gfx.layers {
		// not enough layers in dest
		if len(dest.layers) == i {
			dest.layers = append(dest.layers, new(RectanglesLayer))
		}
		layer.copyTo(dest.layers[i])
	}
}

func (wnd *WindowImpl) OnConfig(config *Configuration) error {
	return nil
}

func (wnd *WindowImpl) OnCreate() error {
	return nil
}

func (wnd *WindowImpl) OnShow() error {
	return nil
}

func (wnd *WindowImpl) OnResize() error {
	wnd.Update()
	return nil
}

func (wnd *WindowImpl) OnMove() error {
	return nil
}

func (wnd *WindowImpl) OnKeyDown(keyCode int, repeated uint) error {
	if repeated == 0 {
		if keyCode == 41 { // ESC
			wnd.Close()
		}
	}
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

func (wnd *WindowImpl) OnDestroy(err error) error {
	return err
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
		event.props.update(wndWrapper.data, wndWrapper.title)
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

func ensureCFloatLen(arr []C.float, length int) []C.float {
	arrLen := len(arr)
	if arrLen < length {
		if cap(arr) < length {
			arrNew := make([]C.float, length)
			copy(arrNew, arr[:arrLen])
			return arrNew
		}
		return arr[:length]
	}
	return arr
}
