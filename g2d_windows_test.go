/*
 *          Copyright 2023, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

package g2d

import (
	"errors"
	"testing"
)

type tTestConfigWindow struct {
	DefaultWindow
	configCalled  bool
	destroyCalled bool
}

func (window *tTestConfigWindow) OnConfig(config *Configuration) error {
	config.AutoUpdate = false
	window.configCalled = true
	return errors.New("test")
}

func (window *tTestConfigWindow) OnDestroy() {
	window.destroyCalled = true
}

type tTestCreateWindow struct {
	tTestConfigWindow
	createCalled bool
}

func (window *tTestCreateWindow) OnConfig(config *Configuration) error {
	config.AutoUpdate = false
	window.configCalled = true
	return nil
}

func (window *tTestCreateWindow) OnCreate(widget *Widget) error {
	window.createCalled = true
	return errors.New("test")
}

func TestInit(t *testing.T) {
	if !initialized {
		Init()
		errs := Errors()
		if len(errs) > 0 {
			t.Error(errs[0].Error())
		}
	}
}

func TestConfigWindow(t *testing.T) {
	if !initialized {
		Init()
	}
	if !initFailed {
		clearErrors()
		window := new(tTestConfigWindow)
		Show(window)
		errs = Errors()
		if !window.configCalled {
			t.Error("OnConfig not called")
		} else if !window.destroyCalled {
			t.Error("OnDestroy not called")
		} else if len(errs) != 1 {
			t.Error(len(errs))
		} else if errs[0].Error() != "test" {
			t.Error(errs[0].Error())
		}
	} else {
		t.Error(errs[0].Error())
	}
}

func TestCreateWindow(t *testing.T) {
	if !initialized {
		Init()
	}
	if !initFailed {
		window := new(tTestCreateWindow)
		clearErrors()
		Show(window)
		errs = Errors()
		if !window.configCalled {
			t.Error("OnConfig not called")
		} else if !window.createCalled {
			t.Error("OnCreate not called")
		} else if !window.destroyCalled {
			t.Error("OnDestroy not called")
		} else if len(errs) != 1 {
			t.Error(len(errs))
		} else if errs[0].Error() != "test" {
			t.Error(errs[0].Error())
		}
	} else {
		t.Error(errs[0].Error())
	}
}

func TestLayerA(t *testing.T) {
	layer := newRectLayer(0, 2)
	rectA := layer.newRect()
	rectB := layer.newRect()
	if rectA.layer != 0 || rectB.layer != 0 {
		t.Error(rectA.layer, rectB.layer)
	} else if rectA.chunk != 0 || rectB.chunk != 0 {
		t.Error(rectA.chunk, rectB.chunk)
	} else if rectA.index != 0 || rectB.index != 1 {
		t.Error(rectA.index, rectB.index)
	} else {
		layer.release(rectA.chunk, rectA.index)
		rectA = layer.newRect()
		rectB = layer.newRect()
		if layer.totalActive != 3 || layer.size != 4 {
			t.Error(layer.totalActive, layer.size)
		} else if rectA.layer != 0 || rectB.layer != 0 {
			t.Error(rectA.layer, rectB.layer)
		} else if rectA.chunk != 0 || rectB.chunk != 1 {
			t.Error(rectA.chunk, rectB.chunk)
		} else if rectA.index != 0 || rectB.index != 0 {
			t.Error(rectA.index, rectB.index)
		}
	}
}

func TestLayerB(t *testing.T) {
	layer := newRectLayer(0, 2)
	layer.newRect()
	rectB := layer.newRect()
	layer.newRect()
	rectD := *rectB
	layer.release(rectB.chunk, rectB.index)
	rectB = rectB
	if rectD != *rectB {
		t.Error(rectD.layer, rectD.chunk, rectD.index, rectB.layer, rectB.chunk, rectB.index)
	}
}
