/*
 *          Copyright 2025, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

package g2d

import (
	"testing"
)

func TestInit(t *testing.T) {
	if Err != nil {
		t.Error("Err is not nil")
	} else if MaxTexSize != 0 {
		t.Error("MaxTexSize is not 0")
	} else if MaxTexUnits != 0 {
		t.Error("MaxTexUnits is not 0")
	} else {
		Init()
		if Err != nil {
			t.Error(Err.Error())
		}
		if MaxTexSize == 0 {
			t.Error("MaxTexSize is 0")
		}
		if MaxTexUnits == 0 {
			t.Error("MaxTexUnits is 0")
		}
	}
}
