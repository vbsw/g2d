/*
 *          Copyright 2024, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

package g2d

// #include "g2d.h"
import "C"
import (
	"unsafe"
)

type Frame struct {
	X, Y, Width, Height int
}

type Time struct {
	Prev, Curr, Delta int64
}

type Mouse struct {
	X, Y int
}

type Window struct {
	Frame      Frame
	Time      Time
	Mouse      Mouse
	state      int
	cbId       int
	dataC      unsafe.Pointer
	autoUpdate bool
	update     bool
	msgs       chan *tLogicMessage
	quitted    chan bool
}

type abstractWindow interface {
	OnConfig(config *Configuration) error
	OnCreate() error
	OnShow() error
	OnWindowMoved() error
	OnWindowResize() error
	OnKeyDown(keyCode int, repeated uint) error
	OnKeyUp(keyCode int) error
	OnMouseMove() error
	OnButtonDown(buttonCode int, doubleClicked bool) error
	OnButtonUp(buttonCode int, doubleClicked bool) error
	OnWheel(rotation float32) error
	OnWindowMinimize() error
	OnWindowRestore() error
	OnTextureLoaded(textureId int) error
	OnUpdate() error
	OnClose() (bool, error)
	OnDestroy() error
	impl() *Window
}

func (wnd *Window) Update() {
	wnd.update = true
	wnd.msgs <- nil
}

func (wnd *Window) Close() {
	wnd.msgs <- (&tLogicMessage{typeId: quitReqType, nanos: time.Nanos()})
}

func (_ *Window) OnConfig(config *Configuration) error {
	return nil
}

func (_ *Window) OnCreate() error {
	return nil
}

func (_ *Window) OnShow() error {
	return nil
}

func (_ *Window) OnWindowMoved() error {
	return nil
}

func (_ *Window) OnWindowResize() error {
	return nil
}

func (_ *Window) OnKeyDown(keyCode int, repeated uint) error {
	return nil
}

func (_ *Window) OnKeyUp(keyCode int) error {
	return nil
}

func (_ *Window) OnMouseMove() error {
	return nil
}

func (_ *Window) OnButtonDown(buttonCode int, doubleClicked bool) error {
	return nil
}

func (_ *Window) OnButtonUp(buttonCode int, doubleClicked bool) error {
	return nil
}

func (_ *Window) OnWheel(rotation float32) error {
	return nil
}

func (_ *Window) OnWindowMinimize() error {
	return nil
}

func (_ *Window) OnWindowRestore() error {
	return nil
}

func (_ *Window) OnTextureLoaded(textureId int) error {
	return nil
}

func (_ *Window) OnUpdate() error {
	return nil
}

func (_ *Window) OnClose() (bool, error) {
	return true, nil
}

func (_ *Window) OnDestroy() error {
	return nil
}

func (wnd *Window) impl() *Window {
	return wnd
}

func logicThread(abst abstractWindow) {
	wnd := abst.impl()
	for wnd.state != quitState {
		msg := wnd.nextLogicMessage()
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
				/*
					case textureType:
						onTextureLoaded(abst, wnd, msg.valA)
				*/
			case updateType:
				onUpdate(abst, wnd)
			case quitReqType:
				onQuitReq(abst, wnd)
			case quitType:
				onQuit(abst, wnd)
			}
		}
	}
	/*
		wnd.Gfx.rBuffer = nil
		wnd.Gfx.wBuffer = nil
		wnd.Gfx.buffers[0].layers = nil
		wnd.Gfx.buffers[1].layers = nil
		wnd.Gfx.buffers[2].layers = nil
		wnd.Gfx.entitiesLayers = nil
	*/
	wnd = nil
}

func (wnd *Window) nextLogicMessage() *tLogicMessage {
	var message *tLogicMessage
	if wnd.state > configState && wnd.state < closingState && (wnd.autoUpdate || wnd.update) {
		select {
		case msg := <-wnd.msgs:
			message = msg
		default:
			wnd.update = false
			message = &tLogicMessage{typeId: updateType, nanos: time.Nanos()}
		}
	} else {
		message = <-wnd.msgs
	}
	if wnd.state == closingState && message.typeId != quitType {
		message = nil
	}
	return message
}

func (wnd *Window) updateProps(msg *tLogicMessage) {
	wnd.Frame.X = msg.props.ClientX
	wnd.Frame.Y = msg.props.ClientY
	wnd.Frame.Width = msg.props.ClientWidth
	wnd.Frame.Height = msg.props.ClientHeight
	wnd.Mouse.X = msg.props.MouseX
	wnd.Mouse.Y = msg.props.MouseY
}

func onConfig(abst abstractWindow, wnd *Window) {
	config := newConfiguration()
	err := abst.OnConfig(config)
	wnd.autoUpdate = config.AutoUpdate
	if err == nil {
		toMainLoop.postMsg(&tCreateWindowRequest{abst: abst, config: config})
	} else {
		onLogicError(abst, wnd, 4999, err)
	}
}

func onCreate(abst abstractWindow, wnd *Window) {
	wnd.Time.Prev = wnd.Time.Curr
	err := abst.OnCreate()
	if err == nil {
		wnd.state = runningState
		/*
			wnd.Gfx.running = true
			wnd.Gfx.msgs <- &tGraphicsMessage{typeId: refreshType}
			go wnd.graphicsThread()
			wnd.Gfx.switchWBuffer()
		*/
		toMainLoop.postMsg(&tShowWindowRequest{abst: abst})
	} else {
		onLogicError(abst, wnd, 4999, err)
	}
}

func onShow(abst abstractWindow, wnd *Window) {
	wnd.Time.Prev = wnd.Time.Curr
	err := abst.OnShow()
	if err == nil {
		/*
			wnd.Gfx.switchWBuffer()
			wnd.Gfx.msgs <- &tGraphicsMessage{typeId: refreshType}
		*/
	} else {
		onLogicError(abst, wnd, 4999, err)
	}
}

func onWindowMoved(abst abstractWindow, wnd *Window) {
	err := abst.OnWindowMoved()
	if err != nil {
		onLogicError(abst, wnd, 4999, err)
	}
}

func onWindowResize(abst abstractWindow, wnd *Window) {
	err := abst.OnWindowResize()
	if err != nil {
		onLogicError(abst, wnd, 4999, err)
	}
}

func onKeyDown(abst abstractWindow, wnd *Window, keyCode int, repeated uint) {
	err := abst.OnKeyDown(keyCode, repeated)
	if err != nil {
		onLogicError(abst, wnd, 4999, err)
	}
}

func onKeyUp(abst abstractWindow, wnd *Window, keyCode int) {
	err := abst.OnKeyUp(keyCode)
	if err != nil {
		onLogicError(abst, wnd, 4999, err)
	}
}

func onMouseMove(abst abstractWindow, wnd *Window) {
	err := abst.OnMouseMove()
	if err != nil {
		onLogicError(abst, wnd, 4999, err)
	}
}

func onButtonDown(abst abstractWindow, wnd *Window, buttonCode int, doubleClicked bool) {
	err := abst.OnButtonDown(buttonCode, doubleClicked)
	if err != nil {
		onLogicError(abst, wnd, 4999, err)
	}
}

func onButtonUp(abst abstractWindow, wnd *Window, buttonCode int, doubleClicked bool) {
	err := abst.OnButtonUp(buttonCode, doubleClicked)
	if err != nil {
		onLogicError(abst, wnd, 4999, err)
	}
}

func onWheel(abst abstractWindow, wnd *Window, rotation float32) {
	err := abst.OnWheel(rotation)
	if err != nil {
		onLogicError(abst, wnd, 4999, err)
	}
}

func onWindowMinimize(abst abstractWindow, wnd *Window) {
	err := abst.OnWindowMinimize()
	if err != nil {
		onLogicError(abst, wnd, 4999, err)
	}
}

func onWindowRestore(abst abstractWindow, wnd *Window) {
	err := abst.OnWindowRestore()
	if err != nil {
		onLogicError(abst, wnd, 4999, err)
	}
}

func onUpdate(abst abstractWindow, wnd *Window) {
	wnd.Time.Delta = wnd.Time.Curr - wnd.Time.Prev
	err := abst.OnUpdate()
	wnd.Time.Prev = wnd.Time.Curr
	if err == nil {
		/*
			wnd.Gfx.switchWBuffer()
			wnd.Gfx.msgs <- &tGraphicsMessage{typeId: refreshType}
		*/
	} else {
		onLogicError(abst, wnd, 4999, err)
	}
}

func onQuitReq(abst abstractWindow, wnd *Window) {
	closeOk, err := abst.OnClose()
	if err == nil {
		if closeOk {
			wnd.Time.Curr = time.Nanos()
			onQuit(abst, wnd)
			toMainLoop.postMsg(&tDestroyWindowRequest{abst: abst})
		}
	} else {
		onLogicError(abst, wnd, 4999, err)
	}
}

func onQuit(abst abstractWindow, wnd *Window) {
	/*
		if wnd.Gfx.running {
			wnd.Gfx.msgs <- &tGraphicsMessage{typeId: quitType}
			<- wnd.Gfx.quitted
		}
	*/
	err := abst.OnDestroy()
	wnd.quitted <- true
	wnd.state = quitState
	if err != nil {
		setErrorSynced(toError(4000, 0, int64(wnd.cbId), err.Error(), nil))
	}
}

func onLogicError(abst abstractWindow, wnd *Window, err1 int64, err error) {
	toMainLoop.postErr(toError(err1, 0, int64(wnd.cbId), err.Error(), nil))
	wnd.Time.Curr = time.Nanos()
	onQuit(abst, wnd)
}
