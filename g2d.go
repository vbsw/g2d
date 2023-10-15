/*
 *          Copyright 2023, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

// Package g2d is a framework to create 2D graphic applications.
package g2d

import "C"
import (
	"github.com/vbsw/golib/queue"
	"sync"
	"time"
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
	resizeType
	keyDownType
	keyUpType
	updateType
	quitReqType
	quitType
	leaveType
	refreshType
	swapIntervType
	imageType
	textureType
)

const (
	chunksCap = 1024 * 1024 - 32
	entitiesCap = 64
	entityValuesCount = 20
	layerEntitesCount = 1024 - 32
)

var fsm = [56]int{0, 1, 2, 0, 10, 2, 2, 1, 3, 1, 2, 0, 3, 4, 1, 0, 6, 5, 1, 2, 0, 4, 1, 0, 6, 7, 0, 2, 13, 8, 0, 1, 9, 7, 0, 2, 9, 5, 1, 2, 10, 11, 0, 1, 9, 12, 0, 2, 13, 11, 0, 1, 13, 2, 2, 1}

var (
	ErrConv     ErrorConvertor
	MaxTextureSize  int
	Err error
	quitting bool

	initialized bool
	initFailed  bool
	wndCbs     []*tWindow
	wndCbsNext []int
	wndsToStart []*tWindow
	msgs       queue.Queue
	wndsActive sync.WaitGroup 

	mutex       sync.Mutex
	startTime   time.Time
	running bool
)

type ErrorConvertor interface {
	ToError(err1, err2 int64, info string) error
}

type Window interface {
	OnConfig(config *Configuration) error
	OnCreate(widget *Widget) error
	OnShow() error
	OnResize() error
	OnKeyDown(keyCode int, repeated uint) error
	OnKeyUp(keyCode int) error
	OnTextureLoaded(textureId int) error
	OnUpdate() error
	OnClose() (bool, error)
	OnDestroy() error
}

type Widget struct {
	ClientX, ClientY          int
	ClientWidth, ClientHeight int
	MouseX, MouseY            int
	PrevUpdateNanos           int64
	CurrEventNanos            int64
	DeltaUpdateNanos          int64
	update                    bool
	Gfx                       Graphics
	msgs                      chan *tLMessage
	quitted                   chan bool
}

type Graphics struct {
	tEntities
	msgs           chan *tGMessage
	quitted                   chan bool
	rBuffer        *tBuffer
	wBuffer        *tBuffer
	buffers        [3]tBuffer
	mutex          sync.Mutex
	bufferState          int
	refresh        bool
	swapInterval          int
	running bool
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

type Entity struct {
	buffer    **tBuffer
	layer, offset int
	chunk, index int
}

type tEntities struct {
	chunks [][]Entity
	unused [][]int
	active []int
}

type tEngineProperties struct {
	errConv     ErrorConvertor
}

type tErrorConvertor struct {
}

type tWindow struct {
	state      int
	cbId       int
	cbIdStr string
	abst       Window
	wgt        *Widget
	dataC      unsafe.Pointer
	autoUpdate bool
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

type tBuffer struct {
	layers [][]C.float
	enabled     [][]C.char
	unused      [][]int
	active      []int
	enabledC []unsafe.Pointer
	layersC []unsafe.Pointer
	bgR, bgG, bgB C.float
}

func engineProperties() *tEngineProperties {
	props := new(tEngineProperties)
	if ErrConv != nil {
		props.errConv = ErrConv
	} else {
		props.errConv = new(tErrorConvertor)
	}
	return props
}

func register(wnd *tWindow) int {
	if len(wndCbsNext) == 0 {
		wndCbs = append(wndCbs, wnd)
		return len(wndCbs) - 1
	}
	indexLast := len(wndCbsNext) - 1
	cbId := wndCbsNext[indexLast]
	wndCbsNext = wndCbsNext[:indexLast]
	wndCbs[cbId] = wnd
	return cbId
}

func unregister(cbId int) {
	wndCbs[cbId] = nil
	wndCbsNext = append(wndCbsNext, cbId)
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

func (wgt *Widget) Update() {
	wgt.update = true
	wgt.msgs <- nil
}

func (wgt *Widget) Close() {
	wgt.msgs <- (&tLMessage{typeId: quitReqType, nanos: deltaNanos()})
}

func (gfx *Graphics) SetVSync(vsync bool) {
	if gfx.running {
		if vsync {
			gfx.swapInterval = 1
			gfx.msgs <- &tGMessage{typeId: swapIntervType, valA: 1}
		} else {
			gfx.swapInterval = 0
			gfx.msgs <- &tGMessage{typeId: swapIntervType, valA: 0}
		}
	} else {
		if vsync {
			gfx.swapInterval = 1
		} else {
			gfx.swapInterval = 0
		}
	}
}

func (gfx *Graphics) SetBGColor(r, g, b float32) {
	gfx.wBuffer.bgR, gfx.wBuffer.bgG, gfx.wBuffer.bgB = C.float(r), C.float(g), C.float(b)
}

func (gfx *Graphics) CreateLayers(count int, initCapacities ...int) {
	if count > 0 {
		layers := make([][]C.float, count)
		enabled := make([][]C.char, count)
		unused := make([][]int, count)
		active := make([]int, count)
		layersC := make([]unsafe.Pointer, count)
		enabledC := make([]unsafe.Pointer, count)
		for i := 0; i< count; i++ {
			var initCapacity int
			if i < len(initCapacities) && initCapacities[i] > 0 {
				initCapacity = initCapacities[i]
			} else {
				initCapacity = layerEntitesCount
			}
			layers[i] = make([]C.float, initCapacity * entityValuesCount)
			enabled[i] = make([]C.char, initCapacity)
			unused[i] = make([]int, 0, initCapacity)
			layersC[i] = unsafe.Pointer(&layers[i][0])
			enabledC[i] = unsafe.Pointer(&enabled[i][0])
		}
		if gfx.wBuffer.layers == nil {
			gfx.wBuffer.layers = layers
			gfx.wBuffer.enabled = enabled
			gfx.wBuffer.unused = unused
			gfx.wBuffer.active = active
			gfx.wBuffer.layersC = layersC
			gfx.wBuffer.enabledC = enabledC
		} else {
			gfx.wBuffer.layers = append(gfx.wBuffer.layers, layers...)
			gfx.wBuffer.enabled = append(gfx.wBuffer.enabled, enabled...)
			gfx.wBuffer.unused = append(gfx.wBuffer.unused, unused...)
			gfx.wBuffer.active = append(gfx.wBuffer.active, active...)
			gfx.wBuffer.layersC = append(gfx.wBuffer.layersC, layersC...)
			gfx.wBuffer.enabledC = append(gfx.wBuffer.enabledC, enabledC...)
		}
	}
}

func (gfx *Graphics) NewEntity(layer int) *Entity {
	entity, chunkIndex, entityIndex := gfx.newEntity()
	offset := gfx.wBuffer.newEntity(layer)
	entity.init(&gfx.wBuffer, layer, offset, chunkIndex, entityIndex)
	return entity
}

func (gfx *Graphics) Release(entity *Entity) {
	dataIndex := entity.offset / entityValuesCount
	(*entity.buffer).enabled[entity.layer][dataIndex] = 0
	(*entity.buffer).unused[entity.layer] = append((*entity.buffer).unused[entity.layer], dataIndex)
	(*entity.buffer).active[entity.layer]--
	entity.buffer = nil
	gfx.unused[entity.chunk] = append(gfx.unused[entity.chunk], entity.index)
}

func (gfx *Graphics) switchRBuffer() {
	gfx.mutex.Lock()
	indexCurr := gfx.bufferState * 4
	gfx.bufferState = fsm[indexCurr]
	indexNext := gfx.bufferState * 4
	gfx.rBuffer = &gfx.buffers[fsm[indexNext+2]]
	gfx.mutex.Unlock()
}

func (gfx *Graphics) switchWBuffer() {
	gfx.mutex.Lock()
	indexCurr := gfx.bufferState * 4
	gfx.bufferState = fsm[indexCurr+1]
	indexNext := gfx.bufferState * 4
	gfx.wBuffer = &gfx.buffers[fsm[indexNext+3]]
	gfx.mutex.Unlock()
	gfx.wBuffer.set(&gfx.buffers[fsm[indexCurr+3]])
}

func (buffer *tBuffer) set(other *tBuffer) {
	buffer.bgR, buffer.bgG, buffer.bgB = other.bgR, other.bgG, other.bgB
	for i, otherLayer := range other.layers {
		if len(buffer.layers) < i {
			if len(buffer.layers[i]) < len(otherLayer) {
				if cap(buffer.layers[i]) < len(otherLayer) {
					buffer.layers[i] = make([]C.float, len(otherLayer), cap(otherLayer))
					buffer.enabled[i] = make([]C.char, len(otherLayer), cap(otherLayer))
					buffer.unused[i] = make([]int, len(other.unused[i]), cap(otherLayer))
					buffer.layersC[i] = unsafe.Pointer(&buffer.layers[i][0])
					buffer.enabledC[i] = unsafe.Pointer(&buffer.enabled[i][0])
				} else {
					buffer.layers[i] = buffer.layers[i][:len(otherLayer)]
					buffer.enabled[i] = buffer.enabled[i][:len(otherLayer)]
					buffer.unused[i] = buffer.unused[i][:len(other.unused[i])]
				}
			}
		} else {
			buffer.layers = append(buffer.layers, make([]C.float, len(otherLayer), cap(otherLayer)))
			buffer.enabled = append(buffer.enabled, make([]C.char, len(otherLayer), cap(otherLayer)))
			buffer.unused = append(buffer.unused, make([]int, len(other.unused[i]), cap(otherLayer)))
			buffer.active = append(buffer.active, 0)
			buffer.layersC = append(buffer.layersC, unsafe.Pointer(&buffer.layers[i][0]))
			buffer.enabledC = append(buffer.enabledC, unsafe.Pointer(&buffer.enabled[i][0]))
		}
		copy(buffer.layers[i], otherLayer)
		copy(buffer.enabled[i], other.enabled[i])
		copy(buffer.unused[i], other.unused[i])
		buffer.active[i] = other.active[i]
	}
}

func (buffer *tBuffer) newEntity(layer int) int {
	var index int
	if len(buffer.layers[layer]) <= buffer.active[layer] {
		layers := make([]C.float, len(buffer.layers) * 2)
		enabled := make([]C.char, len(buffer.enabled) * 2)
		unused := make([]int, 0, len(buffer.enabled) * 2)
		copy(layers, buffer.layers[layer])
		copy(enabled, buffer.enabled[layer])
		buffer.layers[layer] = layers
		buffer.enabled[layer] = enabled
		buffer.unused[layer] = unused
		buffer.layersC[layer] = unsafe.Pointer(&layers[0])
		buffer.enabledC[layer] = unsafe.Pointer(&enabled[0])
		index = buffer.active[layer]
	} else if len(buffer.unused[layer]) > 0 {
		lastIndex := len(buffer.unused[layer])-1
		index = buffer.unused[layer][lastIndex]
		buffer.unused[layer] = buffer.unused[layer][:lastIndex]
	}
	buffer.enabled[layer][index] = 1
	buffer.active[layer]++
	return index * entityValuesCount
}

func (entity *Entity) RGBA() (float32, float32, float32, float32) {
	offset := entity.offset
	layer := (*entity.buffer).layers[entity.layer]
	return float32(layer[offset+4]), float32(layer[offset+5]), float32(layer[offset+6]), float32(layer[offset+7])
}

func (entity *Entity) SetRGBA(r, g, b, a float32) {
	offset := entity.offset
	layer := (*entity.buffer).layers[entity.layer]
	layer[offset+4] = C.float(r)
	layer[offset+5] = C.float(g)
	layer[offset+6] = C.float(b)
	layer[offset+7] = C.float(a)
}

func (entity *Entity) XYWH() (float32, float32, float32, float32) {
	offset := entity.offset
	layer := (*entity.buffer).layers[entity.layer]
	return float32(layer[offset+0]), float32(layer[offset+1]), float32(layer[offset+2]), float32(layer[offset+3])
}

func (entity *Entity) SetXY(x, y float32) {
	offset := entity.offset
	layer := (*entity.buffer).layers[entity.layer]
	layer[offset+0] = C.float(x)
	layer[offset+1] = C.float(y)
}

func (entity *Entity) SetWH(width, height float32) {
	offset := entity.offset
	layer := (*entity.buffer).layers[entity.layer]
	layer[offset+2] = C.float(width)
	layer[offset+3] = C.float(height)
}

func (entity *Entity) SetXYWH(x, y, width, height float32) {
	offset := entity.offset
	layer := (*entity.buffer).layers[entity.layer]
	layer[offset+0] = C.float(x)
	layer[offset+1] = C.float(y)
	layer[offset+2] = C.float(width)
	layer[offset+3] = C.float(height)
}

func (entity *Entity) init(buffer **tBuffer, layer, offset, chunkIndex, entityIndex int) {
	entity.buffer = buffer
	entity.layer = layer
	entity.offset = offset
	entity.chunk = chunkIndex
	entity.index = entityIndex
}

func (entities *tEntities) initEntities() {
	entities.chunks = make([][]Entity, 1, entitiesCap)
	entities.unused = make([][]int, 1, entitiesCap)
	entities.active = make([]int, 1, entitiesCap)
	entities.chunks[0] = make([]Entity, chunksCap)
	entities.unused[0] = make([]int, 0, chunksCap)
}

func (entities *tEntities) newEntity() (*Entity, int, int) {
	for i, active := range entities.active {
		if active < chunksCap {
			var index int
			if len(entities.unused) == 0 {
				index = active
			} else {
				unusedIndex := len(entities.unused) - 1
				index = entities.unused[i][unusedIndex]
				entities.unused[i] = entities.unused[i][:unusedIndex]
			}
			entities.active[i]++
			entity := &entities.chunks[i][index]
			return entity, i, index
		}
	}
	chunks := make([]Entity, chunksCap)
	unused := make([]int, 0, chunksCap)
	entities.chunks = append(entities.chunks, chunks)
	entities.unused = append(entities.unused, unused)
	entities.active = append(entities.active, 1)
	return &chunks[0], len(entities.active) - 1, 0
}

/*


import (
	"sync"
	"time"
	"unsafe"
)

type Graphics struct {
	MaxTextureSize int
}
*/

/*
import "C"
import (
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
	running    bool
	wndsUsed   []*tWindow
	wndsUnused []int
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
*/

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

func (_ *DefaultWindow) OnDestroy() error {
	return nil
}

func deltaNanos() int64 {
	timeNow := time.Now()
	d := timeNow.Sub(startTime)
	return d.Nanoseconds()
}

func anyAvailable(windows []Window) bool {
	for _, wnd := range windows {
		if wnd != nil {
			return true
		}
	}
	return false
}

func toCInt(b bool) C.int {
	if b {
		return C.int(1)
	}
	return C.int(0)
}
