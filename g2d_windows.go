/*
 *          Copyright 2022, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

package g2d

// #cgo CFLAGS: -DVBSW_G2D_WIN32 -DUNICODE
// #cgo LDFLAGS: -luser32 -lgdi32 -lOpenGL32
// #include "g2d.h"
import "C"

func Init() {
	mutex.Lock()
	if !initialized {
		var n1, n2 C.int
		var err1, err2 C.longlong
		C.g2d_init(&n1, &n2, &err1, &err2)
		if err1 == 0 {
			MaxTexSize = int(n1)
			MaxTexUnits = int(n2)
			initialized = true
			initFailed = false
			quitting = false
			time.Reset()
		} else {
			initFailed = true
			Err = toError(int64(err1), int64(err2), 0, "")
		}
		mutex.Unlock()
	} else {
		mutex.Unlock()
		panic("g2d engine already initialized")
	}
}
