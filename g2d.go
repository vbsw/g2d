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
	"unsafe"
)

const (
	notInitialized    = "g2d not initialized"
	alreadyProcessing = "already processing events"
)

var (
	Err         error
	initialized bool
	processing  bool
	cb          tCallback
)

type Parameters struct {
	ClientX, ClientY                  int
	ClientWidth, ClientHeight         int
	ClientMinWidth, ClientMinHeight   int
	ClientMaxWidth, ClientMaxHeight   int
	MouseLocked, Borderless, Dragable bool
	Resizable, Fullscreen, Centered   bool
	Title                             string
}

type Properties struct {
	Fullscreen bool
	Title   string
}

type Command struct {
	CloseReq bool
	CloseUnc bool
}

type Window struct {
	Props Properties
	Cmd Command
}

type AbstractWindow interface {
	Config(params *Parameters) error
	KeyDown(key int, repeated uint) error
	KeyUp(key int) error
	Close() (bool, error)
	Destroy()
	baseStruct() *Window
	updatePropsResetCmd(props Properties)
	propsAndCmd() (Properties, Command)
}

type tManager struct {
	data    unsafe.Pointer
	wndBase *Window
	wndAbst AbstractWindow
	props Properties
	cmd  Command
}

// tCallback holds objects identified by ids.
type tCallback struct {
	mgrs   []*tManager
	unused []int
}

func (window *Window) Config(params *Parameters) error {
	return nil
}

func (window *Window) KeyDown(key int, repeated uint) error {
	return nil
}

func (window *Window) KeyUp(key int) error {
	return nil
}

func (window *Window) Close() (bool, error) {
	return true, nil
}

func (window *Window) Destroy() {
}

func (window *Window) updatePropsResetCmd(props Properties) {
	window.Props = props
	window.Cmd.CloseReq = false
	window.Cmd.CloseUnc = false
}

func (window *Window) propsAndCmd() (Properties, Command) {
	return window.Props, window.Cmd
}

func (window *Window) baseStruct() *Window {
	return window
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
