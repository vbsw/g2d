/*
 *          Copyright 2024, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

package g2d

// #include "g2d.h"
import "C"
import (
	"sync"
)

var fsm = [56]int{0, 1, 2, 0, 10, 2, 2, 1, 3, 1, 2, 0, 3, 4, 1, 0, 6, 5, 1, 2, 0, 4, 1, 0, 6, 7, 0, 2, 13, 8, 0, 1, 9, 7, 0, 2, 9, 5, 1, 2, 10, 11, 0, 1, 9, 12, 0, 2, 13, 11, 0, 1, 13, 2, 2, 1}

type Graphics interface {
	NewRectLayer()
	NewGraphicsLayer()
}

type tGraphics struct {
	// tEntities
	msgs         chan *tGraphicsMessage
	quitted      chan bool
/*
	rBuffer      *tBuffer
	wBuffer      *tBuffer
	buffers      [3]tBuffer
*/
	mutex        sync.Mutex
	bufferState  int
	refresh      bool
	swapInterval int
	running      bool
}
