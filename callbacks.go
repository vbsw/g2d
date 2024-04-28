/*
 *          Copyright 2024, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

package g2d

var (
	wndCbs     []*tWindow
	wndCbsNext []int
)

func register(wnd *tWindow) int {
	var cbId int
	if len(wndCbsNext) == 0 {
		wndCbs = append(wndCbs, wnd)
		cbId = len(wndCbs) - 1
	} else {
		indexLast := len(wndCbsNext) - 1
		cbId = wndCbsNext[indexLast]
		wndCbsNext = wndCbsNext[:indexLast]
		wndCbs[cbId] = wnd
	}
	return cbId
}

func unregister(cbId int) {
	wndCbs[cbId] = nil
	wndCbsNext = append(wndCbsNext, cbId)
}
