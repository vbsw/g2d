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
	poolStates  [56]int
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
	rPool *tPool
	wPool *tPool
	msgs  chan *tGMessage
	pools [3]tPool
	mutex sync.Mutex
	state int
	vsync bool
}

func (gfx *Graphics) Refresh() {
	gfx.msgs <- &tGMessage{typeId: refreshType}
}

func (gfx *Graphics) SetBGColor(r, g, b float32) {
	gfx.wPool.bgR, gfx.wPool.bgG, gfx.wPool.bgB = C.float(r), C.float(g), C.float(b)
}

func (gfx *Graphics) SetVSync(vsync bool) {
	gfx.vsync = vsync
	if vsync {
		gfx.msgs <- &tGMessage{typeId: vsyncType, val: 1}
	} else {
		gfx.msgs <- &tGMessage{typeId: vsyncType, val: 0}
	}
}

func (gfx *Graphics) updateRPool() {
	gfx.mutex.Lock()
	indexCurr := gfx.state * 4
	gfx.state = poolStates[indexCurr]
	indexNext := gfx.state * 4
	gfx.rPool = &gfx.pools[poolStates[indexNext+2]]
	gfx.mutex.Unlock()
}

func (gfx *Graphics) switchWPool() {
	gfx.mutex.Lock()
	indexCurr := gfx.state * 4
	gfx.state = poolStates[indexCurr+1]
	indexNext := gfx.state * 4
	poolPrev := &gfx.pools[poolStates[indexCurr+3]]
	gfx.wPool = &gfx.pools[poolStates[indexNext+3]]
	gfx.wPool.set(poolPrev)
	gfx.mutex.Unlock()
}

type tErrorGenerator interface {
	ToError(g2dErrNum, win32ErrNum uint64, info string) error
}

type tErrorLogger interface {
	LogError(err error)
}

type tErrorHandler struct {
}

type tPool struct {
	bgR, bgG, bgB C.float
}

func (pool *tPool) set(other *tPool) {
	pool.bgR, pool.bgG, pool.bgB = other.bgR, other.bgG, other.bgB
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

type tWindowError struct {
	window *tWindow
	err    error
}

type tStopMainLoop struct {
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
