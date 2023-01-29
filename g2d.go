/*
 *          Copyright 2023, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

// Package g2d is a framework to create 2D graphic applications.
package g2d

// typedef struct { float x, y, w, h, r, g, b, a; } g2d_rect_t;
import "C"
import (
	"github.com/vbsw/golib/queue"
	"sync"
	"time"
	"unsafe"
)

const (
	wrongLayerType = "cast to wrong layer type"
)

const (
	configType = iota
	createType
	showType
	updateType
	quitReqType
	quitType
	leaveType
	refreshType
	vsyncType
)

var (
	errs        []error
	mutex       sync.Mutex
	errGen      tErrorGenerator
	errLog      tErrorLogger
	errHandler  tErrorHandler
	mainLoop    tMainLoop
	initialized bool
	initFailed  bool
	startTime   time.Time
	cb          tCallback
	fsm         [56]int
)

func Errors() []error {
	mutex.Lock()
	defer mutex.Unlock()
	return errs
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
	Title                             string
}

type Window interface {
	OnConfig(config *Configuration) error
	OnCreate(widget *Widget) error
	OnShow() error
	OnClose() (bool, error)
	OnDestroy()
}

type DefaultWindow struct {
}

func (_ *DefaultWindow) OnConfig(config *Configuration) error {
	return nil
}

func (_ *DefaultWindow) OnCreate(widget *Widget) error {
	return nil
}

func (_ *DefaultWindow) OnClose() (bool, error) {
	return true, nil
}

func (_ *DefaultWindow) OnShow() error {
	return nil
}

func (_ *DefaultWindow) OnDestroy() {
}

type Widget struct {
	ClientX, ClientY          int
	ClientWidth, ClientHeight int
	MouseX, MouseY            int
	PrevUpdateNanos           int64
	CurrEventNanos            int64
	Gfx                       Graphics
	msgs                      chan *tLMessage
}

type Graphics struct {
	rBuffer        *tBuffer
	wBuffer        *tBuffer
	msgs           chan *tGMessage
	buffers        [3]tBuffer
	entitiesLayers []tEntitiesLayer
	mutex          sync.Mutex
	state          int
	vsync          bool
}

func (gfx *Graphics) Refresh() {
	gfx.msgs <- &tGMessage{typeId: refreshType}
}

func (gfx *Graphics) SetBGColor(r, g, b float32) {
	gfx.wBuffer.bgR, gfx.wBuffer.bgG, gfx.wBuffer.bgB = C.float(r), C.float(g), C.float(b)
}

func (gfx *Graphics) SetVSync(vsync bool) {
	gfx.vsync = vsync
	if vsync {
		gfx.msgs <- &tGMessage{typeId: vsyncType, val: 1}
	} else {
		gfx.msgs <- &tGMessage{typeId: vsyncType, val: 0}
	}
}

func (gfx *Graphics) NewRectLayer(size int) int {
	layerId := len(gfx.wBuffer.layers)
	gfx.wBuffer.layers = append(gfx.wBuffer.layers, newRectLayer(size))
	gfx.entitiesLayers = append(gfx.entitiesLayers, newRectEntitiesLayer(size))
	return layerId
}

func (gfx *Graphics) NewRect(layer int) *Rect {
	index := gfx.wBuffer.layers[layer].newRectIndex()
	rect := gfx.entitiesLayers[layer].newRectEntity(&gfx.wBuffer, layer, index)
	return rect
}

func (gfx *Graphics) switchRBuffer() {
	gfx.mutex.Lock()
	indexCurr := gfx.state * 4
	gfx.state = fsm[indexCurr]
	indexNext := gfx.state * 4
	gfx.rBuffer = &gfx.buffers[fsm[indexNext+2]]
	gfx.mutex.Unlock()
}

func (gfx *Graphics) switchWBuffer() {
	gfx.mutex.Lock()
	indexCurr := gfx.state * 4
	gfx.state = fsm[indexCurr+1]
	indexNext := gfx.state * 4
	gfx.wBuffer = &gfx.buffers[fsm[indexNext+3]]
	gfx.mutex.Unlock()
	gfx.wBuffer.set(&gfx.buffers[fsm[indexCurr+3]])
}

type Rect struct {
	buffer             **tBuffer
	entitylayer        tEntitiesLayer
	chunk, entityIndex int
	layer, index       int
}

func (rect *Rect) init(buffer **tBuffer, layer, index, chunk, entityIndex int, entitylayer tEntitiesLayer) *Rect {
	rect.buffer = buffer
	rect.layer, rect.index = layer, index
	rect.chunk, rect.entityIndex = chunk, entityIndex
	rect.entitylayer = entitylayer
	return rect
}

func (rect *Rect) XY() (float32, float32) {
	rectC := (*rect.buffer).layers[rect.layer].rect(rect.index)
	return float32(rectC.x), float32(rectC.y)
}

func (rect *Rect) WH() (float32, float32) {
	rectC := (*rect.buffer).layers[rect.layer].rect(rect.index)
	return float32(rectC.w), float32(rectC.h)
}

func (rect *Rect) XYWH() (float32, float32, float32, float32) {
	rectC := (*rect.buffer).layers[rect.layer].rect(rect.index)
	return float32(rectC.x), float32(rectC.y), float32(rectC.w), float32(rectC.h)
}

func (rect *Rect) RGBA() (float32, float32, float32, float32) {
	rectC := (*rect.buffer).layers[rect.layer].rect(rect.index)
	return float32(rectC.r), float32(rectC.g), float32(rectC.b), float32(rectC.a)
}

func (rect *Rect) SetXY(x, y float32) {
	rectC := (*rect.buffer).layers[rect.layer].rect(rect.index)
	rectC.x, rectC.y = C.float(x), C.float(y)
}

func (rect *Rect) SetWH(w, h float32) {
	rectC := (*rect.buffer).layers[rect.layer].rect(rect.index)
	rectC.w, rectC.h = C.float(w), C.float(h)
}

func (rect *Rect) SetXYWH(x, y, w, h float32) {
	rectC := (*rect.buffer).layers[rect.layer].rect(rect.index)
	rectC.x, rectC.y, rectC.w, rectC.h = C.float(x), C.float(y), C.float(w), C.float(h)
}

func (rect *Rect) SetRGBA(r, g, b, a float32) {
	rectC := (*rect.buffer).layers[rect.layer].rect(rect.index)
	rectC.r, rectC.g, rectC.b, rectC.a = C.float(r), C.float(g), C.float(b), C.float(a)
}

func (rect *Rect) SetEnabled(enabled bool) {
	if enabled {
		(*rect.buffer).layers[rect.layer].enable(rect.index)
	} else {
		(*rect.buffer).layers[rect.layer].disable(rect.index)
	}
}

func (rect *Rect) Release() {
	(*rect.buffer).layers[rect.layer].release(rect.index)
	rect.entitylayer.release(rect.chunk, rect.entityIndex)
	rect.buffer = nil
	rect.entitylayer = nil
}

type tBuffer struct {
	bgR, bgG, bgB C.float
	layers        []tLayer
}

func (buffer *tBuffer) set(other *tBuffer) {
	buffer.bgR, buffer.bgG, buffer.bgB = other.bgR, other.bgG, other.bgB
	for i, layer := range buffer.layers {
		layer.set(other.layers[i])
	}
	for _, otherLayer := range buffer.layers[len(buffer.layers):] {
		buffer.layers = append(buffer.layers, otherLayer.clone())
	}
}

type tErrorGenerator interface {
	ToError(g2dErrNum, win32ErrNum uint64, info string) error
}

type tErrorLogger interface {
	LogError(err error)
}

type tErrorHandler struct {
}

type tMainLoop struct {
	mutex      sync.Mutex
	msgs       queue.Queue
	running    bool
	wndsUsed   []*tWindow
	wndsUnused []int
}

func (loop *tMainLoop) nextMessage() interface{} {
	loop.mutex.Lock()
	defer loop.mutex.Unlock()
	return loop.msgs.First()
}

type tWindow struct {
	abst       Window
	wgt        *Widget
	dataC      unsafe.Pointer
	autoUpdate bool
	loopId     int
	cbId       int
	state      int
}

type tCallback struct {
	wnds   []*tWindow
	unused []int
}

// register returns a new id number for wnd. It will not be garbage collected until
// unregister is called with this id.
func (cb *tCallback) register(wnd *tWindow) int {
	var index int
	if len(cb.unused) == 0 {
		cb.wnds = append(cb.wnds, wnd)
		index = len(cb.wnds) - 1
	} else {
		lastIndex := len(cb.unused) - 1
		index = cb.unused[lastIndex]
		cb.unused = cb.unused[:lastIndex]
		cb.wnds[index] = wnd
	}
	return index
}

// unregister makes wnd no more identified by id.
// This object may be garbage collected, now.
func (cb *tCallback) unregister(id int) {
	cb.wnds[id] = nil
	cb.unused = append(cb.unused, id)
}

// unregisterAll makes all regiestered wnds no more identified by id.
// These objects may be garbage collected, now.
func (cb *tCallback) unregisterAll() {
	for i := 0; i < len(cb.wnds) && cb.wnds[i] != nil; i++ {
		cb.unregister(i)
	}
}

type tLMessage struct {
	typeId int
	nanos  int64
	props  Properties
	obj    interface{}
}

type tGMessage struct {
	typeId int
	val    int
}

type tConfigWindowRequest struct {
	window *tWindow
}

type tCreateWindowRequest struct {
	window *tWindow
	config *Configuration
}

type tShowWindowRequest struct {
	window *tWindow
}

type tDestroyWindowRequest struct {
	window *tWindow
}

type tLayer interface {
	newRectIndex() int
	rect(index int) *C.g2d_rect_t
	enable(index int)
	disable(index int)
	clone() tLayer
	set(other tLayer)
	release(index int)
}

type tBaseLayer struct {
	enabled     []C.char
	unused      []int
	totalActive int
}

func (layer *tBaseLayer) initBase(size int) {
	layer.enabled = make([]C.char, 0, size)
	layer.unused = make([]int, 0, size)
}

func (layer *tBaseLayer) usedIndex() int {
	lengthNew := len(layer.unused) - 1
	indexNew := layer.unused[lengthNew]
	layer.enabled[indexNew] = 1
	layer.unused = layer.unused[:lengthNew]
	return indexNew
}

func (layer *tBaseLayer) release(index int) {
	layer.enabled[index] = 0
	layer.unused = append(layer.unused, index)
	layer.totalActive--
}

func (layer *tBaseLayer) enable(index int) {
	layer.enabled[index] = 1
}

func (layer *tBaseLayer) disable(index int) {
	layer.enabled[index] = 0
}

func (layer *tBaseLayer) set(enabled []C.char, unused []int, totalActive int) {
	if cap(layer.enabled) <= len(enabled) {
		layer.enabled = layer.enabled[:len(enabled)]
	} else {
		layer.enabled = make([]C.char, len(enabled), cap(enabled))
	}
	if cap(layer.unused) <= len(unused) {
		layer.unused = layer.unused[:len(unused)]
	} else {
		layer.unused = make([]int, len(unused), cap(unused))
	}
	copy(layer.enabled, enabled)
	copy(layer.unused, unused)
	layer.totalActive = totalActive
}

type tRectLayer struct {
	tBaseLayer
	rects []C.g2d_rect_t
}

func newRectLayer(size int) *tRectLayer {
	layer := new(tRectLayer)
	layer.rects = make([]C.g2d_rect_t, 0, size)
	layer.initBase(size)
	return layer
}

func (layer *tRectLayer) rect(index int) *C.g2d_rect_t {
	return &layer.rects[index]
}

func (layer *tRectLayer) newRectIndex() int {
	layer.totalActive++
	if len(layer.unused) == 0 {
		layer.enabled = append(layer.enabled, 1)
		layer.rects = append(layer.rects, C.g2d_rect_t{})
		return len(layer.enabled) - 1
	}
	return layer.usedIndex()
}

func (layer *tRectLayer) clone() tLayer {
	other := new(tRectLayer)
	other.rects = make([]C.g2d_rect_t, len(layer.rects), cap(layer.rects))
	copy(other.rects, layer.rects)
	other.tBaseLayer.set(layer.enabled, layer.unused, layer.totalActive)
	return other
}

func (layer *tRectLayer) set(other tLayer) {
	otherLayer, ok := other.(*tRectLayer)
	if ok {
		if cap(layer.rects) <= len(otherLayer.rects) {
			layer.rects = layer.rects[:len(otherLayer.rects)]
		} else {
			layer.rects = make([]C.g2d_rect_t, len(otherLayer.rects), cap(otherLayer.rects))
		}
		copy(layer.rects, otherLayer.rects)
		layer.tBaseLayer.set(otherLayer.enabled, otherLayer.unused, otherLayer.totalActive)
	} else {
		panic(wrongLayerType)
	}
}

type tEntitiesLayer interface {
	newRectEntity(buffer **tBuffer, layer, index int) *Rect
	release(chunk, index int)
}

type tBaseEntitiesLayer struct {
	unused [][]int
	size   int
}

func (layer *tBaseEntitiesLayer) initBase(size int) {
	layer.unused = make([][]int, 1)
	layer.unused[0] = make([]int, 0, size)
}

func (layer *tBaseEntitiesLayer) appendBase(size int) {
	if layer.size == 0 {
		layer.size = size
	} else {
		layer.size *= 2
	}
	layer.unused = append(layer.unused, make([]int, 0, layer.size))
}

func (layer *tBaseEntitiesLayer) release(chunk, index int) {
	layer.unused[chunk] = append(layer.unused[chunk], index)
}

type tRectEntitiesLayer struct {
	tBaseEntitiesLayer
	rects [][]Rect
}

func newRectEntitiesLayer(size int) *tRectEntitiesLayer {
	layer := new(tRectEntitiesLayer)
	layer.rects = make([][]Rect, 1)
	layer.rects[0] = make([]Rect, 0, size)
	layer.initBase(size)
	return layer
}

func (layer *tRectEntitiesLayer) newRectEntity(buffer **tBuffer, bufferLayer, index int) *Rect {
	for chunk, rects := range layer.rects {
		unused := layer.unused[chunk]
		lengthUnused := len(unused)
		if lengthUnused == 0 {
			if len(rects) < cap(rects) {
				entityIndex := len(rects)
				layer.rects[chunk] = append(rects, Rect{})
				rect := &layer.rects[chunk][entityIndex]
				return rect.init(buffer, bufferLayer, index, chunk, entityIndex, layer)
			}
		} else if lengthUnused > 0 {
			lengthUnusedNew := lengthUnused - 1
			entityIndex := unused[lengthUnusedNew]
			rect := &layer.rects[chunk][entityIndex]
			layer.unused[chunk] = unused[:lengthUnusedNew]
			return rect.init(buffer, bufferLayer, index, chunk, entityIndex, layer)
		}
	}
	layer.appendBase(len(layer.rects[0]))
	layer.rects = append(layer.rects, make([]Rect, 1, layer.size))
	return (&layer.rects[len(layer.rects)-1][0]).init(buffer, bufferLayer, index, len(layer.unused)-1, 0, layer)
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

func appendError(err error) {
	mutex.Lock()
	errs = append(errs, err)
	errLog.LogError(err)
	mutex.Unlock()
}

func clearErrors() {
	mutex.Lock()
	errs = errs[:0]
	mutex.Unlock()
}

func deltaNanos() int64 {
	timeNow := time.Now()
	d := timeNow.Sub(startTime)
	return d.Nanoseconds()
}

// toCInt converts bool value to C int value.
func toCInt(b bool) C.int {
	if b {
		return C.int(1)
	}
	return C.int(0)
}
