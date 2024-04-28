/*
 *          Copyright 2024, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

package g2d

import (
	stdtime "time"
)

var (
	time tTime
)

type tTime struct {
	start stdtime.Time
}

func (t *tTime) Reset() {
	t.start = stdtime.Now()
}

func (t *tTime) Nanos() int64 {
	now := stdtime.Now()
	delta := now.Sub(t.start)
	return delta.Nanoseconds()
}
