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

type tTestCreateWindow struct {
	tTestConfigWindow
	createCalled bool
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

func (window *tTestConfigWindow) OnConfig(config *Configuration) error {
	config.AutoUpdate = false
	window.configCalled = true
	return errors.New("test")
}

func (window *tTestConfigWindow) OnDestroy() {
	window.destroyCalled = true
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
