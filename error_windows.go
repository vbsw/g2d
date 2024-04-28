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
	"strconv"
	"unsafe"
)

var (
	Err *Error
)

type Error struct {
	AllocErr, InitErr, RunErr int64
	UnkErr, SysErr, WndId     int64
	Str, SysInfo              string
}

func setError(err *Error) {
	if Err == nil {
		Err = err
	}
	quitMainLoop(0, 0)
}

func setErrorSynced(err *Error) {
	mutex.Lock()
	if Err == nil {
		Err = err
	}
	quitMainLoop(0, 0)
	mutex.Unlock()
}

func toError(err1, err2, wndId int64, sysInfo string, sysInfoC *C.char) *Error {
	if err1 > 0 {
		var err *Error
		if len(sysInfo) == 0 && sysInfoC != nil {
			sysInfo = C.GoString(sysInfoC)
			C.g2d_free(unsafe.Pointer(sysInfoC))
		}
		if err1 < 2000 {
			err = &Error{AllocErr: err1, SysErr: err2, WndId: wndId, SysInfo: sysInfo}
		} else if err1 < 3000 {
			err = &Error{InitErr: err1, SysErr: err2, WndId: wndId, SysInfo: sysInfo}
		} else if err1 < 4000 {
			err = &Error{RunErr: err1, SysErr: err2, WndId: wndId, SysInfo: sysInfo}
		} else {
			err = &Error{UnkErr: err1, SysErr: err2, WndId: wndId, SysInfo: sysInfo}
		}
		err.createInfo()
		return err
	}
	return nil
}

func (err *Error) createInfo() {
	if len(err.Str) == 0 {
		if err.AllocErr != 0 {
			err.Str = "memory allocation failed"
			err.Str += " (" + strconv.FormatInt(err.AllocErr, 10)
		} else if err.InitErr != 0 {
			err.Str = "g2d initialization failed"
			err.Str += " (" + strconv.FormatInt(err.InitErr, 10)
		} else if err.RunErr != 0 {
			err.Str = "g2d runtime failed"
			err.Str += " (" + strconv.FormatInt(err.RunErr, 10)
		} else {
			err.Str = "unknown"
			err.Str += " (" + strconv.FormatInt(err.UnkErr, 10)
		}
		if err.SysErr == 0 {
			err.Str = err.Str + ")"
		} else {
			err.Str = err.Str + ", " + strconv.FormatInt(err.SysErr, 10) + ")"
		}
		if len(err.SysInfo) != 0 {
			err.Str = err.Str + "; " + err.SysInfo
		}
	}
}
