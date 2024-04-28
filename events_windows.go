/*
 *          Copyright 2024, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

package g2d

// #include "g2d.h"
import "C"

var (
	events   []tEvent
	eventsOn []tEvent
)

type tEvent interface {
	OnEvent()
}

func switchEvents() {
	eventsTmp := events
	events = eventsOn[:0]
	eventsOn = eventsTmp
}
