/*
 *          Copyright 2024, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

package g2d

var (
	abstCbs     []abstractWindow
	abstCbsNext []int
)

func register(abst abstractWindow) int {
	var cbIndex int
	if len(abstCbsNext) == 0 {
		abstCbs = append(abstCbs, abst)
		cbIndex = len(abstCbs) - 1
	} else {
		indexLast := len(abstCbsNext) - 1
		cbIndex = abstCbsNext[indexLast]
		abstCbsNext = abstCbsNext[:indexLast]
		abstCbs[cbIndex] = abst
	}
	return cbIndex
}

func unregister(cbIndex int) {
	abstCbs[cbIndex] = nil
	abstCbsNext = append(abstCbsNext, cbIndex)
}
