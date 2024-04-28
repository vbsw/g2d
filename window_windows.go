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
	NanosPrev                 int64
	NanosCurr                 int64
	NanosDelta                int64
	update                    bool
	// Gfx                       Graphics
	msgs    chan *tLogicMessage
	quitted chan bool
}

type tWindow struct {
	state      int
	cbId       int
	abst       Window
	wgt        *Widget
	dataC      unsafe.Pointer
	autoUpdate bool
}

func newWindowWrapper(window Window) *tWindow {
	wnd := new(tWindow)
	wnd.state = configState
	wnd.abst = window
	wnd.wgt = new(Widget)
	wnd.wgt.msgs = make(chan *tLogicMessage, 1000)
	wnd.wgt.quitted = make(chan bool, 1)
	wnd.cbId = register(wnd)
	/*
	   wnd.wgt.Gfx.msgs = make(chan *tGMessage, 1000)
	   wnd.wgt.Gfx.quitted = make(chan bool, 1)
	   wnd.wgt.Gfx.rBuffer = &wnd.wgt.Gfx.buffers[0]
	   wnd.wgt.Gfx.wBuffer = &wnd.wgt.Gfx.buffers[0]
	   wnd.wgt.Gfx.initEntities()
	*/
	return wnd
}

func (wnd *tWindow) logicThread() {
	for wnd.state != quitState {
		msg := wnd.nextLogicMessage()
		if msg != nil {
			wnd.wgt.NanosCurr = msg.nanos
			switch msg.typeId {
			case configType:
				wnd.onConfig()
				/*
					case createType:
						wnd.onCreate()
					case showType:
						wnd.onShow()
					case resizeType:
						wnd.updateProps(msg)
						wnd.onResize()
					case keyDownType:
						wnd.onKeyDown(msg.valA, msg.repeated)
					case keyUpType:
						wnd.onKeyUp(msg.valA)
					case textureType:
						wnd.onTextureLoaded(msg.valA)
					case updateType:
						wnd.onUpdate()
					case quitReqType:
						wnd.onQuitReq()
					case quitType:
						wnd.onQuit()
				*/
			}
		}
	}
	/*
		wnd.wgt.Gfx.rBuffer = nil
		wnd.wgt.Gfx.wBuffer = nil
		wnd.wgt.Gfx.buffers[0].layers = nil
		wnd.wgt.Gfx.buffers[1].layers = nil
		wnd.wgt.Gfx.buffers[2].layers = nil
		wnd.wgt.Gfx.entitiesLayers = nil
	*/
	wnd.wgt = nil
}

func (wnd *tWindow) nextLogicMessage() *tLogicMessage {
	var message *tLogicMessage
	if wnd.state > configState && wnd.state < closingState && (wnd.autoUpdate || wnd.wgt.update) {
		select {
		case msg := <-wnd.wgt.msgs:
			message = msg
		default:
			wnd.wgt.update = false
			message = &tLogicMessage{typeId: updateType, nanos: time.Nanos()}
		}
	} else {
		message = <-wnd.wgt.msgs
	}
	if wnd.state == closingState && message.typeId != quitType {
		message = nil
	}
	return message
}

func (wnd *tWindow) onConfig() {
	config := newConfiguration()
	err := wnd.abst.OnConfig(config)
	wnd.autoUpdate = config.AutoUpdate
	if err == nil {
		postEvent(&tCreateWindowRequest{window: wnd, config: config}, 3998, int64(wnd.cbId))
	} else {
		wnd.onLogicError(err)
	}
}

func (wnd *tWindow) onQuit() {
	/*
		if wnd.wgt.Gfx.running {
			wnd.wgt.Gfx.msgs <- &tGMessage{typeId: quitType}
			<- wnd.wgt.Gfx.quitted
		}
	*/
	err := wnd.abst.OnDestroy()
	wnd.wgt.quitted <- true
	wnd.state = quitState
	if err != nil {
		setErrorSynced(toError(4000, 0, int64(wnd.cbId), err.Error(), nil))
	}
}

func (wnd *tWindow) onLogicError(err error) {
	setErrorSynced(toError(4000, 0, int64(wnd.cbId), err.Error(), nil))
	wnd.wgt.NanosCurr = time.Nanos()
	wnd.onQuit()
}
