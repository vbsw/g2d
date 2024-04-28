/*
 *          Copyright 2024, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

// Package g2d is a framework to create 2D graphic applications.
package g2d

import (
	"sync"
)

var (
	MaxTexSize, MaxTexUnits int
	initialized, initFailed bool
	quitting                bool
	mutex                   sync.Mutex
)
