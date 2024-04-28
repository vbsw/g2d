/*
 *          Copyright 2024, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

package g2d

import (
	"testing"
)

func TestInit(t *testing.T) {
	Init()
	if Err != nil {
		t.Error(Err.Str)
	}
}
