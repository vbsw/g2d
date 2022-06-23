/*
 *          Copyright 2022, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

// Package g2d is a framework to create 2D graphic applications.
package g2d

import "fmt"

type Engine struct {
	infoOnly bool
}

type AbstractEngine interface {
	baseStruct() *Engine
	ParseOSArgs() error
	Info()
	CreateWindow()
	Error(err error)
}

type WindowBuilder struct {
	ClientX, ClientY                  int
	ClientWidth, ClientHeight         int
	ClientMinWidth, ClientMinHeight   int
	ClientMaxWidth, ClientMaxHeight   int
	MouseLocked, Borderless, Dragable bool
	Resizable, Fullscreen, Centered   bool
	Handler                           EventHandler
}

type EventHandler interface {
}

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

func (engine *Engine) CreateWindow() {
}

func (engine *Engine) Error(err error) {
	fmt.Println("error:", err.Error())
}
