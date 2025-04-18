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
	"errors"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
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
	textureType    = 16
	texBufType     = 17
	minimizeType   = 18
	restoreType    = 19
	focusType      = 20
	customType     = 21
	refreshType    = 22
)

var (
	MaxTexSize, MaxTexUnits int
	MaxTextures             int
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

// Window is callback for window handling.
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
	OnTextureLoaded(texture Texture) error
	OnFramebufferCreated(buffer Framebuffer) error
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

// WindowImpl is the obligatory struct to embed, when using interface Window.
type WindowImpl struct {
	Props Properties
	Stats Stats
	Gfx   Graphics
	id    int
}

// Configuration is the initial setting of window.
type Configuration struct {
	ClientX, ClientY                  int
	ClientWidth, ClientHeight         int
	ClientWidthMin, ClientHeightMin   int
	ClientWidthMax, ClientHeightMax   int
	MouseLocked, Borderless, Dragable bool
	Resizable, Fullscreen, Centered   bool
	Title                             string
}

// Properties are the current window properties.
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

// Stats has useful data. Time is in milliseconds.
type Stats struct {
	AppTime, DeltaTime int
	lastUpdate         int
	UPS, FPS           int
	ups, lastUPSTime   int
	fps, lastFPSTime   int
}

// Graphics draws graphics.
type Graphics struct {
	BgR, BgG, BgB float32
	VSync, AVSync bool
	eventsChan    chan *tGraphicsEvent
	quittedChan   chan bool
	mutex         sync.Mutex
	w, h          int
	read          *tGfxBuffer
	buffer        *tGfxBuffer
	bufferReady   bool
	updating      bool
	running       bool
	glTexIds      []int
	texDims       []int
	Layers        []Layer
}

// RectanglesLayer is a layer holding rectangles.
type RectanglesLayer struct {
	entities     []*Rectangle
	entityNextId []int
	count        int
	Enabled      bool
	texMap       []int
}

// Rectangle is an entity from a RectanglesLayer.
type Rectangle struct {
	id                   int
	X, Y, Width, Height  float32
	R, G, B, A           float32
	RotX, RotY, RotAlpha float32
	TexRef, TexX, TexY   int
	TexWidth, TexHeight  int
	Enabled              bool
}

// Texture provides a texture.
type Texture interface {
	Id() int
	RGBABytes() ([]byte, error)
	Dimensions() (int, int)
	IsMipMap() bool
}

// Framebuffer provides a buffer to draw to.
type Framebuffer interface {
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

type tGfxBuffer struct {
	w, h, sw    C.int
	r, g, b     C.float
	batches     [][]C.float
	batchesPtrs []*C.float
	lengths     []C.int
	procs       []unsafe.Pointer
}

type Layer interface {
	getBatch([]Layer, []int, []C.float) ([]Layer, []C.float, C.int, unsafe.Pointer)
}

type tLogicEvent struct {
	typeId   int
	valA     int
	valB     int
	valC     float32
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
	valD   []byte
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

type tErrorRequest struct {
	err error
}

type tAppTime struct {
	start time.Time
}

// ImageFromFile reads image from file. (This is for test purposes.)
func ImageFromFile(path string) (image.Image, error) {
	var img image.Image
	file, err := os.Open(path)
	if err == nil {
		defer file.Close()
		ext := filepath.Ext(path)
		if ext == ".jpg" || ext == ".jpeg" {
			img, err = jpeg.Decode(file)
		} else if ext == ".png" || ext == ".apng" {
			img, err = png.Decode(file)
		} else if ext == ".gif" {
			img, err = gif.Decode(file)
		} else {
			img, _, err = image.Decode(file)
		}
	}
	return img, err
}

// BytesFromImage returns bytes from image. (This is for test purposes.)
func BytesFromImage(img image.Image) []byte {
	var bytes []byte
	switch imgStruct := img.(type) {
	case *image.RGBA:
		bytes = imgStruct.Pix
	case *image.RGBA64:
		bytes = imgStruct.Pix
	case *image.Alpha:
		bytes = imgStruct.Pix
	case *image.Alpha16:
		bytes = imgStruct.Pix
	case *image.CMYK:
		bytes = imgStruct.Pix
	case *image.Gray:
		bytes = imgStruct.Pix
	case *image.Gray16:
		bytes = imgStruct.Pix
	case *image.NRGBA:
		bytes = imgStruct.Pix
	case *image.NRGBA64:
		bytes = imgStruct.Pix
	case *image.Paletted:
		bytes = imgStruct.Pix
	}
	return bytes
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
	wnd.impl.Gfx.read = new(tGfxBuffer)
	wnd.impl.Gfx.buffer = new(tGfxBuffer)
	wnd.impl.Gfx.glTexIds = make([]int, MaxTextures, MaxTextures)
	wnd.impl.Gfx.texDims = make([]int, MaxTextures*2, MaxTextures*2)
	for i := range wnd.impl.Gfx.glTexIds {
		wnd.impl.Gfx.glTexIds[i] = -1
	}
	return wnd
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
		stats.lastUPSTime += diff - diff%1000
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
		stats.lastFPSTime += diff - diff%1000
	}
}

// LoadTexture loads a texture into video memory. This is asynchronous.
// After the texture has been loaded OnTextureLoaded is called.
// (Loading a new texture with same id as a previous one will release the
// previous texture.)
func (gfx *Graphics) LoadTexture(texture Texture) {
	go func() {
		var w, h int
		bytes, err := texture.RGBABytes()
		if err == nil {
			id := texture.Id()
			if id < 0 || id >= MaxTextures {
				err = errors.New(fmt.Sprintf("texture (%i) has invalid id", id))
			} else {
				w, h = texture.Dimensions()
				if w <= 0 || h <= 0 {
					err = errors.New(fmt.Sprintf("texture (%i) has invalid dimensions (%i, %i)", id, w, h))
				}
			}
		}
		if err == nil {
			gfx.eventsChan <- &tGraphicsEvent{typeId: textureType, valC: texture, valD: bytes, err: err}
		} else {
			gfx.eventsChan <- &tGraphicsEvent{typeId: textureType, err: err}
		}
	}()
}

// CreateFramebuffer creates a buffer to draw to.
func (gfx *Graphics) CreateFramebuffer(texture Texture) {
	go func() {
		var w, h int
		bytes, err := texture.RGBABytes()
		if err == nil {
			id := texture.Id()
			if texture.Id() < 0 || texture.Id() >= MaxTextures {
				err = errors.New(fmt.Sprintf("texture (%i) has invalid id", id))
			} else {
				w, h = texture.Dimensions()
				if w <= 0 || h <= 0 {
					err = errors.New(fmt.Sprintf("texture (%i) has invalid dimensions (%i, %i)", id, w, h))
				}
			}
		}
		if err == nil {
			gfx.eventsChan <- &tGraphicsEvent{typeId: texBufType, valA: w, valB: h, valC: texture, valD: bytes, err: err}
		} else {
			gfx.eventsChan <- &tGraphicsEvent{typeId: texBufType, err: err}
		}
	}()
}

func (gfx *Graphics) getReadBuffer() *tGfxBuffer {
	if gfx.bufferReady {
		tmp := gfx.buffer
		gfx.buffer = gfx.read
		gfx.read = tmp
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
	gfx.buffer.adopt(gfx.Layers, gfx.texDims, gfx.w, gfx.h, swapInt, gfx.BgR, gfx.BgG, gfx.BgB)
	gfx.bufferReady = true
	if !gfx.updating {
		gfx.updating = true
		gfx.eventsChan <- &tGraphicsEvent{typeId: refreshType}
	}
	gfx.mutex.Unlock()
}

// NewEntity returns a new instance of Rectangle.
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
	entity.TexRef = -1
	layer.count++
	return entity
}

// Release releases the entity. This entity may be reused when calling NewEntity.
func (layer *RectanglesLayer) Release(r *Rectangle) *Rectangle {
	r.Enabled = false
	layer.entityNextId = append(layer.entityNextId, r.id)
	layer.count--
	return nil
}

// UseTexture associates a reference with a texture. Reference must be in range of [0, 15].
func (layer *RectanglesLayer) UseTexture(ref, textureId int) {
	if ref >= 0 && ref <= 15 {
		if textureId >= 0 && textureId <= 79 {
			if len(layer.texMap) != 0 {
				layer.texMap[ref] = textureId
			} else {
				layer.texMap = make([]int, 16, 16)
				layer.texMap[ref] = textureId
			}
		} else {
			panic(fmt.Sprintf("invalid texture id (%i)", textureId))
		}
	} else {
		panic(fmt.Sprintf("invalid reference (%i) for texture (%i)", ref, textureId))
	}
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
				// error occurred while onConfig, when graphics thread has not been started
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
						wnd.onWheel(event.valC)
					case updateType:
						wnd.onUpdate()
					case closeType:
						wnd.onClose()
					case textureType:
						wnd.onTextureLoaded(event.obj.(Texture))
					case texBufType:
						wnd.onFramebufferCreated(event.obj.(Texture))
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

func (wnd *tWindow) onTextureLoaded(texture Texture) {
	props := wnd.impl.Props
	err := wnd.abst.OnTextureLoaded(texture)
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

func (wnd *tWindow) onFramebufferCreated(texture Framebuffer) {
	props := wnd.impl.Props
	err := wnd.abst.OnFramebufferCreated(texture)
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

func (buf *tGfxBuffer) adopt(layers []Layer, texDims []int, w, h, sw int, r, g, b float32) {
	var index int
	buf.w, buf.h, buf.sw = C.int(w), C.int(h), C.int(sw)
	buf.r, buf.g, buf.b = C.float(r), C.float(g), C.float(b)
	buf.batches = buf.batches[:0]
	buf.batchesPtrs = buf.batchesPtrs[:0]
	buf.lengths = buf.lengths[:0]
	buf.procs = buf.procs[:0]
	for len(layers) > 0 {
		if len(buf.batches) == index {
			if cap(buf.batches) == index {
				buf.batches = append(buf.batches, make([]C.float, 500))
				buf.lengths = append(buf.lengths, 0)
			} else {
				buf.batches = buf.batches[:index+1]
				buf.lengths = buf.lengths[:index+1]
			}
			if cap(buf.procs) == index {
				buf.procs = append(buf.procs, nil)
			} else {
				buf.procs = buf.procs[:index+1]
			}
		}
		layers, buf.batches[index], buf.lengths[index], buf.procs[index] = layers[index].getBatch(layers, texDims, buf.batches[index][:0])
		buf.batchesPtrs = append(buf.batchesPtrs, &buf.batches[index][0])
		index++
	}
}

// OnConfig is called before creating the window. Configuration
// is used to create a window.
func (wnd *WindowImpl) OnConfig(config *Configuration) error {
	return nil
}

// OnCreate is called after window has been created. This event
// could be used to allocate ressources and set graphics.
func (wnd *WindowImpl) OnCreate() error {
	return nil
}

// OnShow is called after window has been set visible. This event
// culd be used to start animations.
func (wnd *WindowImpl) OnShow() error {
	return nil
}

// OnResize is called when Window is resized. New properties are
// in wnd.Props.
func (wnd *WindowImpl) OnResize() error {
	wnd.Update()
	return nil
}

// OnMove is called when Window is moved. New properties are
// in wnd.Props.
func (wnd *WindowImpl) OnMove() error {
	return nil
}

// OnKeyDown is called when a key has been pressed. If key stayes pressed, parameter
// repeated is != 0.
func (wnd *WindowImpl) OnKeyDown(keyCode int, repeated uint) error {
	if repeated == 0 {
		if keyCode == 41 { // ESC
			wnd.Close()
		}
	}
	return nil
}

// OnKeyUp is called when a key has been released.
func (wnd *WindowImpl) OnKeyUp(keyCode int) error {
	return nil
}

// OnMouseMove is called when mouse is moved. New properties are
// in wnd.Props.
func (wnd *WindowImpl) OnMouseMove() error {
	return nil
}

// OnButtonDown is called when mouse button has been pressed.
func (wnd *WindowImpl) OnButtonDown(buttonCode int, doubleClicked bool) error {
	return nil
}

// OnButtonUp is called when mouse button has been released.
func (wnd *WindowImpl) OnButtonUp(buttonCode int, doubleClicked bool) error {
	return nil
}

// OnWheel is called when mouse wheel has been moved.
func (wnd *WindowImpl) OnWheel(rotation float32) error {
	return nil
}

// OnCustom is called after calling Custom().
func (wnd *WindowImpl) OnCustom(obj interface{}) error {
	return nil
}

// OnTextureLoaded is called after texture has been loaded to video memory.
func (wnd *WindowImpl) OnTextureLoaded(texture Texture) error {
	return nil
}

// OnFramebufferCreated is called after texture has been created.
func (wnd *WindowImpl) OnFramebufferCreated(texture Framebuffer) error {
	return nil
}

// OnCustom is called after calling Update(). After OnUpdate graphis
// is redrawn.
func (wnd *WindowImpl) OnUpdate() error {
	return nil
}

// OnClose is called after close button of window has been pressed. If
// function returns true, window will be destroyed.
func (wnd *WindowImpl) OnClose() (bool, error) {
	return true, nil
}

// OnDestroy is called before window gets destroyed.
func (wnd *WindowImpl) OnDestroy(err error) error {
	return err
}

// OnMinimize is called after window has been minimized.
func (wnd *WindowImpl) OnMinimize() error {
	return nil
}

// OnRestore is called after window has returned from minimized state.
func (wnd *WindowImpl) OnRestore() error {
	return nil
}

// OnFocus is called when window gets or loses focus.
func (wnd *WindowImpl) OnFocus(focus bool) error {
	return nil
}

// Close triggers OnClose event.
func (wnd *WindowImpl) Close() {
	postRequest(&tCloseWindowRequest{wndId: wnd.id})
}

// Update triggers OnUpdate event.
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

// Quit destroys window unconditionally.
func (wnd *WindowImpl) Quit() {
	postRequest(&tDestroyWindowRequest{wndId: wnd.id})
}

// Custom triggers OnCustom event.
func (wnd *WindowImpl) Custom(obj interface{}) {
	postRequest(&tCustomRequest{wndId: wnd.id, obj: obj})
}

// Show creates a new window.
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
