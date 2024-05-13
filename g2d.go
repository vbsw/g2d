/*
 *          Copyright 2024, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

// Package g2d is a framework to create 2D graphic applications.
package g2d

import "C"
import (
	"sync"
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
	wndMoveType
	wndResizeType
	keyDownType
	keyUpType
	msMoveType
	buttonDownType
	buttonUpType
	wheelType
	updateType
	quitReqType
	quitType
	leaveType
	refreshType
	swapIntervType
	imageType
	textureType
)

var (
	MaxTexSize, MaxTexUnits int
	initialized, initFailed bool
	running, quitting       bool
	mutex                   sync.Mutex
)

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

func toIntC(b bool) C.int {
	if b {
		return C.int(1)
	}
	return C.int(0)
}
