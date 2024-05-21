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
	"errors"
	"runtime"
	"sync"
)

var fsm = [56]int{0, 1, 2, 0, 10, 2, 2, 1, 3, 1, 2, 0, 3, 4, 1, 0, 6, 5, 1, 2, 0, 4, 1, 0, 6, 7, 0, 2, 13, 8, 0, 1, 9, 7, 0, 2, 9, 5, 1, 2, 10, 11, 0, 1, 9, 12, 0, 2, 13, 11, 0, 1, 13, 2, 2, 1}

type Graphics interface {
	NewRectLayer() RectLayer
	NewGraphicsLayer() Graphics
	SetVSync(vsync bool)
	SetBGColor(r, g, b float32)
}

type RectLayer interface {
}

type tGraphics struct {
	// tEntities
	msgs    chan *tGraphicsMessage
	quitted chan bool
	/*
		rBuffer      *tBuffer
		wBuffer      *tBuffer
		buffers      [3]tBuffer
	*/
	mutex         sync.Mutex
	bufferState   int
	refresh       bool
	swapInterval  int
	running       bool
	bgR, bgG, bgB C.float
}

func (gfx *tGraphics) NewRectLayer() RectLayer {
	return nil
}

func (gfx *tGraphics) NewGraphicsLayer() Graphics {
	graphics := new(tGraphics)
	return graphics
}

func (gfx *tGraphics) SetVSync(vsync bool) {
	if gfx.running {
		if vsync {
			gfx.swapInterval = 1
			gfx.msgs <- &tGraphicsMessage{typeId: swapIntervType, valA: 1}
		} else {
			gfx.swapInterval = 0
			gfx.msgs <- &tGraphicsMessage{typeId: swapIntervType, valA: 0}
		}
	} else {
		if vsync {
			gfx.swapInterval = 1
		} else {
			gfx.swapInterval = 0
		}
	}
}

func (gfx *tGraphics) SetBGColor(r, g, b float32) {
	gfx.bgR, gfx.bgG, gfx.bgB = C.float(r), C.float(g), C.float(b)
	if gfx.running {
		gfx.msgs <- &tGraphicsMessage{typeId: refreshType}
	}
}

func graphicsThread(abst abstractWindow, wnd *Window) {
	var err1, err2 C.longlong
	var errStrC *C.char
	runtime.LockOSThread()
	C.g2d_gfx_make_current(wnd.dataC, &err1, &err2)
	if err1 == 0 {
		C.g2d_gfx_init(wnd.dataC, C.int(wnd.gfxImpl.swapInterval), &err1, &err2, &errStrC)
		if err1 == 0 {
			for wnd.gfxImpl.running {
				msg := wnd.gfxImpl.nextGMessage()
				if msg != nil {
					switch msg.typeId {
					case refreshType:
						draw(wnd)
					case swapIntervType:
						C.g2d_gfx_set_swap_interval(C.int(msg.valA))
					case wndResizeType:
						C.g2d_gfx_set_view_size(wnd.dataC, C.int(msg.valA), C.int(msg.valB))
						draw(wnd)
						/*
							case imageType:
								texBytes, ok := msg.valC.([]byte)
								if ok {
									wnd.loadTexture(texBytes, msg.valA, msg.valB)
								} else {
									appendError(msg.err)
									processing = wnd.processGMessage(&tGraphicsMessage{typeId: quitType})
								}
						*/
					case quitType:
						C.g2d_gfx_release(wnd.dataC, &err1, &err2)
						if err1 == 0 {
							wnd.gfxImpl.quitted <- true
							wnd.gfxImpl.running = false
						}
					}
				}
			}
		}
	}
	if err1 != 0 {
		onLogicError(abst, wnd, 4999, errors.New("g2d graphics error"))
	}
}

func draw(wnd *Window) {
	var err1, err2 C.longlong
	/*
		wnd.wgt.Gfx.switchRBuffer()
		buffer := wnd.wgt.Gfx.rBuffer
		C.g2d_gfx_clear_bg(buffer.bgR, buffer.bgG, buffer.bgB)
	*/
	C.g2d_gfx_clear_bg(wnd.gfxImpl.bgR, wnd.gfxImpl.bgG, wnd.gfxImpl.bgB)
	/*
		for _, layer := range wnd.wgt.Gfx.rBuffer.layers {
			err := layer.draw(wnd.dataC)
			if err != nil {
				appendError(err)
				wnd.wgt.Gfx.msgs <- &tGMessage{typeId: quitType}
			}
		}
	*/
	C.g2d_gfx_swap_buffers(wnd.dataC, &err1, &err2)
	if err1 != 0 {
		//wnd.onGError(err1, err2, nil)
	}
}

func (gfx *tGraphics) switchWBuffer() {
}

func (gfx *tGraphics) destroy() {
	/*
		wnd.Gfx.rBuffer = nil
		wnd.Gfx.wBuffer = nil
		wnd.Gfx.buffers[0].layers = nil
		wnd.Gfx.buffers[1].layers = nil
		wnd.Gfx.buffers[2].layers = nil
		wnd.Gfx.entitiesLayers = nil
	*/
}

func (gfx *tGraphics) nextGMessage() *tGraphicsMessage {
	var message *tGraphicsMessage
	if gfx.refresh {
		select {
		case msg := <-gfx.msgs:
			if msg.typeId != refreshType {
				message = msg
			}
		default:
			gfx.refresh = false
			message = &tGraphicsMessage{typeId: refreshType}
		}
	} else {
		message = <-gfx.msgs
		if message.typeId == refreshType {
			gfx.refresh = true
			message = nil
		}
	}
	return message
}
