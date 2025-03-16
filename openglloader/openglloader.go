/*
 *        Copyright 2023, 2025 Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

// Package openglloader loads OpenGL functions.
package openglloader

import (
	"github.com/vbsw/golib/cdata"
)

type tInitializer struct {
}

func NewInitializer() cdata.CData {
	return new(tInitializer)
}
