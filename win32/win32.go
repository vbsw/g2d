/*
 *          Copyright 2023, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

// Package win32 is the implementation of the g2d engine for windows OS.
package win32

// #cgo CFLAGS: -DG2D_IMPL_WIN32 -DUNICODE
// #cgo LDFLAGS: -luser32 -lgdi32 -lOpenGL32
// #include "win32.h"
import "C"
import (
	"unsafe"
)

type InitParams struct {
	Data   []*unsafe.Pointer
	Funcs  []unsafe.Pointer
	Engine *unsafe.Pointer
}

type InitResult struct {
	G2dErrNum, Win32ErrNum uint64
	ErrStr                 string
	MaxTextureSize         int
}

func Init(params *InitParams) *InitResult {
	var maxTexSizeC, errNumC C.int
	var errWin32C C.g2d_ul_t
	result := new(InitResult)
	dataLen := params.dataLength()
	if dataLen > 0 {
		// because &params.Data[0] is a Go pointer and can not be used in C
		data := make([]unsafe.Pointer, dataLen)
		dataC := (*unsafe.Pointer)(&data[0])
		funcsC := (*unsafe.Pointer)(&params.Funcs[0])
		C.g2d_win32_init(params.Engine, dataC, funcsC, C.int(dataLen), &maxTexSizeC, &errNumC, &errWin32C)
		if errNumC == 0 {
			for i, d := range data {
				*params.Data[i] = d
			}
		}
	} else {
		C.g2d_win32_init(params.Engine, nil, nil, C.int(dataLen), &maxTexSizeC, &errNumC, &errWin32C)
	}
	result.MaxTextureSize = int(maxTexSizeC)
	result.G2dErrNum = uint64(errNumC)
	result.Win32ErrNum = uint64(errWin32C)
	return result
}

func (params *InitParams) dataLength() int {
	if len(params.Data) <= len(params.Funcs) {
		return len(params.Data)
	}
	return len(params.Funcs)
}

/*
func (engine *Engine) toError(errNumC C.int, errWin32C C.g2d_ul_t, errStrC *C.char) error {
	var errStr string
	if errStrC != nil {
		errStr = C.GoString(errStrC)
		C.g2d_free(unsafe.Pointer(errStrC))
	}
	return engine.errConv.ToError(uint64(errNumC), uint64(errWin32C), errStr)
}
*/
