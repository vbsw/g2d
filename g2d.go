/*
 *          Copyright 2023, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

// Package g2d is a framework to create 2D graphic applications.
package g2d

import (
	"sync"
	"time"
	"unsafe"
)

var (
	mutex       sync.Mutex
	engines     []*Engine
	enginesNext []int
)

type ErrorConvertor interface {
	ToError(err1, err2 int64, info string) error
}

type Window interface {
/*
	OnConfig(config *Configuration) error
	OnCreate(widget *Widget) error
	OnShow() error
	OnResize() error
	OnKeyDown(keyCode int, repeated uint) error
	OnKeyUp(keyCode int) error
	OnTextureLoaded(textureId int) error
	OnUpdate() error
	OnClose() (bool, error)
*/
	OnDestroy()
}

type Engine struct {
	ErrConv     ErrorConvertor
	MaxTexSize  int
	dataC       unsafe.Pointer
	initialized bool
	initFailed  bool
	startTime   time.Time
	mutex       sync.Mutex
	running bool
/*
	errLogger   ErrorLogger
	errs        []error
	Gfx         Graphics
*/
}

type engineProperties struct {
	errConv     ErrorConvertor
}

type defaultErrorConvertor struct {
}

func (engine *Engine) getProperties() *engineProperties {
	props := new(engineProperties)
	if engine.ErrConv != nil {
		props.errConv = engine.ErrConv
	} else {
		props.errConv = new(defaultErrorConvertor)
	}
	return props
}

/*
func unregister(id int) int {
	mutex.Lock()
	defer mutex.Unlock()
	wndsUsed[id] = nil
	wndsUnused = append(enginesNext, id)
	return len(engines) - len(enginesNext)
}


import (
	"sync"
	"time"
	"unsafe"
)

var fsm = [56]int{0, 1, 2, 0, 10, 2, 2, 1, 3, 1, 2, 0, 3, 4, 1, 0, 6, 5, 1, 2, 0, 4, 1, 0, 6, 7, 0, 2, 13, 8, 0, 1, 9, 7, 0, 2, 9, 5, 1, 2, 10, 11, 0, 1, 9, 12, 0, 2, 13, 11, 0, 1, 13, 2, 2, 1}

type Graphics struct {
	MaxTextureSize int
}
*/

/*
import "C"
import (
	"github.com/vbsw/golib/queue"
	"sync"
	"time"
	"unsafe"
	"image/png"
	"os"
	"image"
	"path/filepath"
	"image/jpeg"
	"image/gif"
)

const (
	wrongLayerType = "cast to wrong layer type"
	notImplemented = "function not implemented"
)

const (
	configType = iota
	createType
	showType
	resizeType
	keyDownType
	keyUpType
	updateType
	quitReqType
	quitType
	leaveType
	refreshType
	vsyncType
	imageType
	textureType
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
	maxTexSize  C.int
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

type TextureProvider interface {
	RGBABytes() ([]byte, int, int, error)
}

type ImageLoader struct {
}

func (loader *ImageLoader) RGBABytes() ([]byte, int, int, error) {
	var img image.Image
	var err error
	var bytes []byte
	var width, height int
	path := "./test.png"
	if len(path) > 0 {
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
			bounds := img.Bounds()
			xMin := bounds.Min.X
			xMax := bounds.Max.X
			yMin := bounds.Min.Y
			yMax := bounds.Max.Y
			width = xMax - xMin
			height = yMax - yMin

			switch imgStruct := img.(type) {
			case *image.RGBA:
				bytes = imgStruct.Pix
				println("RGBA stride", imgStruct.Stride, "width", width, "height", height, "bytes", len(bytes))
			case *image.RGBA64:
				bytes = imgStruct.Pix
				println("RGBA64 stride", imgStruct.Stride, "width", width, "height", height, "bytes", len(bytes))
			case *image.Alpha:
				bytes = imgStruct.Pix
				println("Alpha stride", imgStruct.Stride, "width", width, "height", height, "bytes", len(bytes))
			case *image.Alpha16:
				bytes = imgStruct.Pix
				println("Alpha16 stride", imgStruct.Stride, "width", width, "height", height, "bytes", len(bytes))
			case *image.CMYK:
				bytes = imgStruct.Pix
				println("CMYK stride", imgStruct.Stride, "width", width, "height", height, "bytes", len(bytes))
			case *image.Gray:
				bytes = imgStruct.Pix
				println("Gray stride", imgStruct.Stride, "width", width, "height", height, "bytes", len(bytes))
			case *image.Gray16:
				bytes = imgStruct.Pix
				println("Gray16 stride", imgStruct.Stride, "width", width, "height", height, "bytes", len(bytes))
			case *image.NRGBA:
				bytes = imgStruct.Pix
				println("NRGBA stride", imgStruct.Stride, "width", width, "height", height, "bytes", len(bytes))
			case *image.NRGBA64:
				bytes = imgStruct.Pix
				println("NRGBA64 stride", imgStruct.Stride, "width", width, "height", height, "bytes", len(bytes))
			case *image.Paletted:
				bytes = imgStruct.Pix
				println("Paletted stride", imgStruct.Stride, "width", width, "height", height, "bytes", len(bytes))
			default:
				panic("image format not supported")
			}
		}
	}
	return bytes, width, height, err
}

type DefaultWindow struct {
}

func (_ *DefaultWindow) OnConfig(config *Configuration) error {
	return nil
}

func (_ *DefaultWindow) OnCreate(widget *Widget) error {
	return nil
}

func (_ *DefaultWindow) OnUpdate() error {
	return nil
}

func (_ *DefaultWindow) OnClose() (bool, error) {
	return true, nil
}

func (_ *DefaultWindow) OnShow() error {
	return nil
}

func (_ *DefaultWindow) OnResize() error {
	return nil
}

func (_ *DefaultWindow) OnKeyDown(keyCode int, repeated uint) error {
	return nil
}

func (_ *DefaultWindow) OnKeyUp(keyCode int) error {
	return nil
}

func (_ *DefaultWindow) OnTextureLoaded(textureId int) error {
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
	update                    bool
	Gfx                       Graphics
	msgs                      chan *tLMessage
}

func (wgt *Widget) Update() {
	wgt.update = true
	wgt.msgs <- nil
}

func (wgt *Widget) RequestClose() {
	msg := &tLMessage{typeId: quitReqType, nanos: deltaNanos()}
	wgt.msgs <- msg
}

func (wgt *Widget) Close() {
	wgt.msgs <- (&tLMessage{typeId: quitType, nanos: deltaNanos()})
}

type Graphics struct {
	MaxTextureSize int
	rBuffer        *tBuffer
	wBuffer        *tBuffer
	msgs           chan *tGMessage
	buffers        [3]tBuffer
	entitiesLayers []tEntitiesLayer
	mutex          sync.Mutex
	state          int
	refresh        bool
	vsync          bool
}

func (gfx *Graphics) SetBGColor(r, g, b float32) {
	gfx.wBuffer.bgR, gfx.wBuffer.bgG, gfx.wBuffer.bgB = C.float(r), C.float(g), C.float(b)
}

func (gfx *Graphics) SetVSync(vsync bool) {
	gfx.vsync = vsync
	if vsync {
		gfx.msgs <- &tGMessage{typeId: vsyncType, valA: 1}
	} else {
		gfx.msgs <- &tGMessage{typeId: vsyncType, valA: 0}
	}
}

func (gfx *Graphics) LoadTexture(texture TextureProvider) {
	go func() {
		bytes, w, h, err := texture.RGBABytes()
		gfx.msgs <- &tGMessage{typeId: imageType, valA: w, valB: h, valC: bytes, err: err}
	}()
}

func (gfx *Graphics) NewRectLayer(size int) int {
	layerId := len(gfx.wBuffer.layers)
	gfx.wBuffer.layers = append(gfx.wBuffer.layers, newRectLayer(size))
	gfx.entitiesLayers = append(gfx.entitiesLayers, newRectEntitiesLayer(size))
	return layerId
}

func (gfx *Graphics) NewImageLayer(textureId, size int) int {
	layerId := len(gfx.wBuffer.layers)
	gfx.wBuffer.layers = append(gfx.wBuffer.layers, newImageLayer(textureId, size))
	gfx.entitiesLayers = append(gfx.entitiesLayers, newImageEntitiesLayer(size))
	return layerId
}

func (gfx *Graphics) NewRect(layer int) *Rect {
	index := gfx.wBuffer.layers[layer].newRectIndex()
	rect := gfx.entitiesLayers[layer].newRectEntity(&gfx.wBuffer, layer, index)
	return rect
}

func (gfx *Graphics) NewImage(layer int) *Image {
	index := gfx.wBuffer.layers[layer].newImageIndex()
	rect := gfx.entitiesLayers[layer].newImageEntity(&gfx.wBuffer, layer, index)
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
	entityLayer        tEntitiesLayer
	chunk, entityIndex int
	layer, index       int
}

func (rect *Rect) init(buffer **tBuffer, layer, index, chunk, entityIndex int, entityLayer tEntitiesLayer) *Rect {
	rect.buffer = buffer
	rect.layer, rect.index = layer, index
	rect.chunk, rect.entityIndex = chunk, entityIndex
	rect.entityLayer = entityLayer
	return rect
}

func (rect *Rect) XY() (float32, float32) {
	offset := rect.index * 20
	data := (*rect.buffer).layers[rect.layer].data()
	return float32(data[offset]), float32(data[offset+1])
}

func (rect *Rect) WH() (float32, float32) {
	offset := rect.index * 20
	data := (*rect.buffer).layers[rect.layer].data()
	return float32(data[offset+2]), float32(data[offset+3])
}

func (rect *Rect) XYWH() (float32, float32, float32, float32) {
	offset := rect.index * 20
	data := (*rect.buffer).layers[rect.layer].data()
	return float32(data[offset]), float32(data[offset+1]), float32(data[offset+2]), float32(data[offset+3])
}

func (rect *Rect) RGBA() (float32, float32, float32, float32) {
	offset := rect.index * 20
	data := (*rect.buffer).layers[rect.layer].data()
	return float32(data[offset+4]), float32(data[offset+5]), float32(data[offset+6]), float32(data[offset+7])
}

func (rect *Rect) SetXY(x, y float32) {
	(*rect.buffer).layers[rect.layer].setData2(rect.index*20, C.float(x), C.float(y))
}

func (rect *Rect) SetWH(w, h float32) {
	(*rect.buffer).layers[rect.layer].setData2(rect.index*20+2, C.float(w), C.float(h))
}

func (rect *Rect) SetXYWH(x, y, w, h float32) {
	(*rect.buffer).layers[rect.layer].setData4(rect.index*20, C.float(x), C.float(y), C.float(w), C.float(h))
}

func (rect *Rect) SetRGBA(r, g, b, a float32) {
	offset := rect.index * 20
	(*rect.buffer).layers[rect.layer].setData4(offset+4, C.float(r), C.float(g), C.float(b), C.float(a))
	(*rect.buffer).layers[rect.layer].setData4(offset+8, C.float(r), C.float(g), C.float(b), C.float(a))
	(*rect.buffer).layers[rect.layer].setData4(offset+12, C.float(r), C.float(g), C.float(b), C.float(a))
	(*rect.buffer).layers[rect.layer].setData4(offset+16, C.float(r), C.float(g), C.float(b), C.float(a))
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
	rect.entityLayer.release(rect.chunk, rect.entityIndex)
	rect.buffer = nil
	rect.entityLayer = nil
}





type Image struct {
	buffer             **tBuffer
	entityLayer        tEntitiesLayer
	chunk, entityIndex int
	layer, index       int
}

func (rect *Image) init(buffer **tBuffer, layer, index, chunk, entityIndex int, entityLayer tEntitiesLayer) *Image {
	rect.buffer = buffer
	rect.layer, rect.index = layer, index
	rect.chunk, rect.entityIndex = chunk, entityIndex
	rect.entityLayer = entityLayer
	return rect
}

func (rect *Image) XY() (float32, float32) {
	offset := rect.index * 26
	data := (*rect.buffer).layers[rect.layer].data()
	return float32(data[offset]), float32(data[offset+1])
}

func (rect *Image) WH() (float32, float32) {
	offset := rect.index * 26
	data := (*rect.buffer).layers[rect.layer].data()
	return float32(data[offset+2]), float32(data[offset+3])
}

func (rect *Image) XYWH() (float32, float32, float32, float32) {
	offset := rect.index * 26
	data := (*rect.buffer).layers[rect.layer].data()
	return float32(data[offset]), float32(data[offset+1]), float32(data[offset+2]), float32(data[offset+3])
}

func (rect *Image) TexXYWH() (float32, float32, float32, float32) {
	offset := rect.index * 26
	data := (*rect.buffer).layers[rect.layer].data()
	return float32(data[offset+4]), float32(data[offset+5]), float32(data[offset+6]), float32(data[offset+7])
}

func (rect *Image) RGBA() (float32, float32, float32, float32) {
	offset := rect.index * 26
	data := (*rect.buffer).layers[rect.layer].data()
	return float32(data[offset+8]), float32(data[offset+9]), float32(data[offset+10]), float32(data[offset+11])
}

func (rect *Image) SetXY(x, y float32) {
	(*rect.buffer).layers[rect.layer].setData2(rect.index*26, C.float(x), C.float(y))
}

func (rect *Image) SetWH(w, h float32) {
	(*rect.buffer).layers[rect.layer].setData2(rect.index*26+2, C.float(w), C.float(h))
}

func (rect *Image) SetXYWH(x, y, w, h float32) {
	(*rect.buffer).layers[rect.layer].setData4(rect.index*26, C.float(x), C.float(y), C.float(w), C.float(h))
}

func (rect *Image) SetTexXYWH(x, y, w, h float32) {
	(*rect.buffer).layers[rect.layer].setData4(rect.index*26+4, C.float(x), C.float(y), C.float(w), C.float(h))
}

func (rect *Image) SetRGBA(r, g, b, a float32) {
	offset := rect.index * 26
	(*rect.buffer).layers[rect.layer].setData4(offset+8, C.float(r), C.float(g), C.float(b), C.float(a))
	(*rect.buffer).layers[rect.layer].setData4(offset+12, C.float(r), C.float(g), C.float(b), C.float(a))
	(*rect.buffer).layers[rect.layer].setData4(offset+16, C.float(r), C.float(g), C.float(b), C.float(a))
	(*rect.buffer).layers[rect.layer].setData4(offset+20, C.float(r), C.float(g), C.float(b), C.float(a))
}

func (rect *Image) SetEnabled(enabled bool) {
	if enabled {
		(*rect.buffer).layers[rect.layer].enable(rect.index)
	} else {
		(*rect.buffer).layers[rect.layer].disable(rect.index)
	}
}

func (rect *Image) Release() {
	(*rect.buffer).layers[rect.layer].release(rect.index)
	rect.entityLayer.release(rect.chunk, rect.entityIndex)
	rect.buffer = nil
	rect.entityLayer = nil
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
	for _, otherLayer := range other.layers[len(buffer.layers):] {
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
	typeId   int
	valA     int
	repeated uint
	nanos    int64
	props    Properties
	obj      interface{}
}

type tGMessage struct {
	typeId int
	valA   int
	valB   int
	valC   interface{}
	err    error
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
	newImageIndex() int
	data() []C.float
	setData2(offset int, a, b C.float)
	setData4(offset int, a, b, c, d C.float)
	enable(index int)
	disable(index int)
	clone() tLayer
	draw(dataC unsafe.Pointer) error
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
	if cap(layer.enabled) >= len(enabled) {
		layer.enabled = layer.enabled[:len(enabled)]
	} else {
		layer.enabled = make([]C.char, len(enabled), cap(enabled))
	}
	if cap(layer.unused) >= len(unused) {
		layer.unused = layer.unused[:len(unused)]
	} else {
		layer.unused = make([]int, len(unused), cap(unused))
	}
	copy(layer.enabled, enabled)
	copy(layer.unused, unused)
	layer.totalActive = totalActive
}

func (layer *tBaseLayer) clone(enabled []C.char, unused []int, totalActive int) {
	layer.enabled = make([]C.char, len(enabled), cap(enabled))
	layer.unused = make([]int, len(unused), cap(unused))
	copy(layer.enabled, enabled)
	copy(layer.unused, unused)
	layer.totalActive = totalActive
}

type tRectLayer struct {
	tBaseLayer
	rects []C.float
}

func newRectLayer(size int) *tRectLayer {
	layer := new(tRectLayer)
	layer.rects = make([]C.float, 0, size*20)
	layer.initBase(size)
	return layer
}

func (layer *tRectLayer) newRectIndex() int {
	layer.totalActive++
	if len(layer.unused) == 0 {
		layer.enabled = append(layer.enabled, 1)
		if cap(layer.rects)-len(layer.rects) >= 20 {
			layer.rects = layer.rects[:len(layer.rects)+20]
		} else {
			rectsNew := make([]C.float, len(layer.rects)+20, cap(layer.rects)*2)
			copy(rectsNew, layer.rects)
			layer.rects = rectsNew
		}
		return len(layer.enabled) - 1
	}
	return layer.usedIndex()
}

func (layer *tRectLayer) newImageIndex() int {
	panic(notImplemented)
}

func (layer *tRectLayer) data() []C.float {
	return layer.rects
}

func (layer *tRectLayer) setData2(offset int, a, b C.float) {
	layer.rects[offset] = a
	layer.rects[offset+1] = b
}

func (layer *tRectLayer) setData4(offset int, a, b, c, d C.float) {
	layer.rects[offset] = a
	layer.rects[offset+1] = b
	layer.rects[offset+2] = c
	layer.rects[offset+3] = d
}

func (layer *tRectLayer) clone() tLayer {
	other := new(tRectLayer)
	other.rects = make([]C.float, len(layer.rects), cap(layer.rects))
	copy(other.rects, layer.rects)
	other.tBaseLayer.clone(layer.enabled, layer.unused, layer.totalActive)
	return other
}

func (layer *tRectLayer) set(other tLayer) {
	otherLayer, ok := other.(*tRectLayer)
	if ok {
		if cap(layer.rects) >= len(otherLayer.rects) {
			layer.rects = layer.rects[:len(otherLayer.rects)]
		} else {
			layer.rects = make([]C.float, len(otherLayer.rects), cap(otherLayer.rects))
		}
		copy(layer.rects, otherLayer.rects)
		layer.tBaseLayer.set(otherLayer.enabled, otherLayer.unused, otherLayer.totalActive)
	} else {
		panic(wrongLayerType)
	}
}




type tImageLayer struct {
	tBaseLayer
	rects []C.float
	textureId int
}

func newImageLayer(textureId, size int) *tImageLayer {
	layer := new(tImageLayer)
	layer.rects = make([]C.float, 0, size*26)
	layer.textureId = textureId
	layer.initBase(size)
	return layer
}

func (layer *tImageLayer) newRectIndex() int {
	panic(notImplemented)
}

func (layer *tImageLayer) newImageIndex() int {
	layer.totalActive++
	if len(layer.unused) == 0 {
		layer.enabled = append(layer.enabled, 1)
		if cap(layer.rects)-len(layer.rects) >= 26 {
			layer.rects = layer.rects[:len(layer.rects)+26]
		} else {
			rectsNew := make([]C.float, len(layer.rects)+26, cap(layer.rects)*2)
			copy(rectsNew, layer.rects)
			layer.rects = rectsNew
		}
		return len(layer.enabled) - 1
	}
	return layer.usedIndex()
}

func (layer *tImageLayer) data() []C.float {
	return layer.rects
}

func (layer *tImageLayer) setData2(offset int, a, b C.float) {
	layer.rects[offset] = a
	layer.rects[offset+1] = b
}

func (layer *tImageLayer) setData4(offset int, a, b, c, d C.float) {
	layer.rects[offset] = a
	layer.rects[offset+1] = b
	layer.rects[offset+2] = c
	layer.rects[offset+3] = d
}

func (layer *tImageLayer) clone() tLayer {
	other := new(tImageLayer)
	other.rects = make([]C.float, len(layer.rects), cap(layer.rects))
	other.textureId = layer.textureId
	copy(other.rects, layer.rects)
	other.tBaseLayer.clone(layer.enabled, layer.unused, layer.totalActive)
	return other
}

func (layer *tImageLayer) set(other tLayer) {
	otherLayer, ok := other.(*tImageLayer)
	if ok {
		otherLayer.textureId = layer.textureId
		if cap(layer.rects) >= len(otherLayer.rects) {
			layer.rects = layer.rects[:len(otherLayer.rects)]
		} else {
			layer.rects = make([]C.float, len(otherLayer.rects), cap(otherLayer.rects))
		}
		copy(layer.rects, otherLayer.rects)
		layer.tBaseLayer.set(otherLayer.enabled, otherLayer.unused, otherLayer.totalActive)
	} else {
		panic(wrongLayerType)
	}
}




type tEntitiesLayer interface {
	newRectEntity(buffer **tBuffer, layer, index int) *Rect
	newImageEntity(buffer **tBuffer, layer, index int) *Image
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

func (layer *tRectEntitiesLayer) newImageEntity(buffer **tBuffer, bufferLayer, index int) *Image {
	panic(notImplemented)
}




type tImageEntitiesLayer struct {
	tBaseEntitiesLayer
	rects [][]Image
}

func newImageEntitiesLayer(size int) *tImageEntitiesLayer {
	layer := new(tImageEntitiesLayer)
	layer.rects = make([][]Image, 1)
	layer.rects[0] = make([]Image, 0, size)
	layer.initBase(size)
	return layer
}

func (layer *tImageEntitiesLayer) newRectEntity(buffer **tBuffer, bufferLayer, index int) *Rect {
	panic(notImplemented)
}

func (layer *tImageEntitiesLayer) newImageEntity(buffer **tBuffer, bufferLayer, index int) *Image {
	for chunk, rects := range layer.rects {
		unused := layer.unused[chunk]
		lengthUnused := len(unused)
		if lengthUnused == 0 {
			if len(rects) < cap(rects) {
				entityIndex := len(rects)
				layer.rects[chunk] = append(rects, Image{})
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
	layer.rects = append(layer.rects, make([]Image, 1, layer.size))
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
*/
