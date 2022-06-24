/*
 *          Copyright 2022, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

// Package g2d is a framework to create 2D graphic applications.
package g2d

import "C"
import (
	"fmt"
	"unsafe"
)

type Engine struct {
	infoOnly bool
}

type AbstractEngine interface {
	baseStruct() *Engine
	ParseOSArgs() error
	Info()
	CreateWindow() error
	Error(err error)
}

type WindowBuilder struct {
	ClientX, ClientY                  int
	ClientWidth, ClientHeight         int
	ClientMinWidth, ClientMinHeight   int
	ClientMaxWidth, ClientMaxHeight   int
	MouseLocked, Borderless, Dragable bool
	Resizable, Fullscreen, Centered   bool
	Handler                           AbstractEventHandler
	Title string
}

type EventHandler struct {
}

type AbstractEventHandler interface {
	OnClose() (bool, error)
	OnDestroy()
}

type tManager struct {
	handler AbstractEventHandler
	data unsafe.Pointer
	err error
}

// tCallback holds objects identified by ids.
type tCallback struct {
	mgrs  []*tManager
	unused []int
}

var (
	running, initialized bool
	cb tCallback
)

func (engine *Engine) baseStruct() *Engine {
	return engine
}

func (engine *Engine) ParseOSArgs() error {
	return nil
}

func (engine *Engine) SetInfoOnly(infoOnly bool) {
	engine.infoOnly = infoOnly
}

func (engine *Engine) Info() {
}

func (engine *Engine) CreateWindow() error {
	return nil
}

func (engine *Engine) Error(err error) {
	fmt.Println("error:", err.Error())
}

func (handler *EventHandler) OnClose() (bool, error) {
	return true, nil
}

func (handler *EventHandler) OnDestroy() {
}

// Register returns a new id number for mgr. It will not be garbage collected until
// Unregister is called with this id.
func (cb *tCallback) Register(mgr *tManager) int {
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

// toCInt converts bool value to C int value.
func toCInt(b bool) C.int {
	if b {
		return C.int(1)
	}
	return C.int(0)
}
