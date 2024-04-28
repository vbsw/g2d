/*
 *          Copyright 2024, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

package g2d

import (
	"errors"
	"testing"
)

type testWindow struct {
	WindowDummy
}

func (wnd *testWindow) OnConfig(config *Configuration) error {
	return errors.New("testA")
}

func TestInit(t *testing.T) {
	Init()
	if Err == nil {
		if MaxTexSize <= 0 {
			t.Error("MaxTexSize not initialized")
		} else if MaxTexUnits <= 0 {
			t.Error("MaxTexUnits not initialized")
		} else {
			MainLoop(new(testWindow))
			if Err == nil {
				t.Error("error missing")
			} else if Err.SysInfo != "testA" {
				t.Error("unexpected error:", Err.SysInfo)
			}
		}
	} else {
		t.Error(Err.Str)
	}
}
