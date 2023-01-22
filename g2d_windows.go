/*
 *          Copyright 2023, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

package g2d

// #cgo CFLAGS: -DG2D_WIN32 -DUNICODE
// #cgo LDFLAGS: -luser32 -lgdi32 -lOpenGL32
// #include "g2d.h"
import "C"
import (
	"errors"
	"fmt"
	"strconv"
	"time"
	"unsafe"
)

const (
	postMessageErrStr  = "post message failed"
	mallocErrStr   = "memory allocation failed"
	getProcAddrErrStr = "load %s function failed"
)

type ErrorGenerator interface {
	ToError(g2dErrNum, win32ErrNum uint64, info string) error
}

type ErrorLogger interface {
	LogError(err error)
}

type tErrorHandler struct {
}

func Init(params ...interface{}) {
	if !initialized {
		var errNumC C.int
		var errWin32C C.g2d_ul_t
		initCustParams(params)
		initDefaultParams()
		C.g2d_init(&errNumC, &errWin32C)
		if errNumC == 0 {
			startTime = time.Now()
			initialized = true
		} else {
			appendError(toError(errNumC, errWin32C, nil))
		}
	} else {
		panic("g2d is already initialized")
	}
}

func initCustParams(params []interface{}) {
	for i, param := range params {
		var ok, used bool
		errGen, ok = param.(ErrorGenerator)
		used = used || ok
		errLog, ok = param.(ErrorLogger)
		used = used || ok
		if !used {
			panic(fmt.Sprintf("parameter %d is not used", i))
		}
	}
}

func initDefaultParams() {
	if errGen == nil {
		errGen = &errHandler
	}
	if errLog == nil {
		errLog = &errHandler
	}
}

func toError(errNumC C.int, errWin32C C.g2d_ul_t, errStrC *C.char) error {
	var errStr string
	if errStrC != nil {
		errStr = C.GoString(errStrC)
		C.g2d_free(unsafe.Pointer(errStrC))
	}
	return errGen.ToError(uint64(errNumC), uint64(errWin32C), errStr)
}

func (_ *tErrorHandler) ToError(g2dErrNum, win32ErrNum uint64, info string) error {
	var errStr string
	switch g2dErrNum {
	case 1:
		errStr = "get module instance failed"
	case 2:
		errStr = "register dummy class failed"
	case 3:
		errStr = "create dummy window failed"
	case 4:
		errStr = "get dummy device context failed"
	case 5:
		errStr = "choose dummy pixel format failed"
	case 6:
		errStr = "set dummy pixel format failed"
	case 7:
		errStr = "create dummy render context failed"
	case 8:
		errStr = "make dummy context current failed"
	case 9:
		errStr = "release dummy context failed"
	case 10:
		errStr = "deleting dummy render context failed"
	case 11:
		errStr = "destroying dummy window failed"
	case 12:
		errStr = "unregister dummy class failed"

	case 13:
		errStr = "register class failed"
	case 14:
		errStr = "create window failed"
	case 15:
		errStr = "get device context failed"
	case 16:
		errStr = "choose pixel format failed"
	case 17:
		errStr = "set pixel format failed"
	case 18:
		errStr = "create render context failed"
	case 19:
		errStr = "release context failed"
	case 20:
		errStr = "deleting render context failed"
	case 21:
		errStr = "destroying window failed"
	case 22:
		errStr = "unregister class failed"
	case 23:
		errStr = "show window failed; type Window is not embedded"

	case 56:
		errStr = "make context current failed"
	case 61:
		errStr = "swap buffer failed"
	case 62:
		errStr = "set title failed"
	case 63:
		errStr = "wgl functions not initialized"
	case 65:
		errStr = "set title failed"
	case 66:
		errStr = "set cursor position failed"
	case 67:
		errStr = "set fullscreen failed"
	case 68:
		errStr = "set window position failed"
	case 69:
		errStr = "move window failed"
	case 80:
		errStr = postMessageErrStr
	case 81:
		errStr = postMessageErrStr
	case 82:
		errStr = postMessageErrStr
	case 83:
		errStr = postMessageErrStr

	case 100:
		errStr = "not initialized"
	case 101:
		errStr = "not initialized"
	case 102:
		errStr = "not initialized"
	case 120:
		errStr = mallocErrStr
	case 121:
		errStr = mallocErrStr

	case 200:
		errStr = fmt.Sprintf(getProcAddrErrStr, "wglChoosePixelFormatARB")
	case 201:
		errStr = fmt.Sprintf(getProcAddrErrStr, "wglCreateContextAttribsARB")
	case 202:
		errStr = fmt.Sprintf(getProcAddrErrStr, "wglSwapIntervalEXT")
	case 203:
		errStr = fmt.Sprintf(getProcAddrErrStr, "wglGetSwapIntervalEXT")
	case 204:
		errStr = fmt.Sprintf(getProcAddrErrStr, "glCreateShader")
	case 205:
		errStr = fmt.Sprintf(getProcAddrErrStr, "glShaderSource")
	case 206:
		errStr = fmt.Sprintf(getProcAddrErrStr, "glCompileShader")
	case 207:
		errStr = fmt.Sprintf(getProcAddrErrStr, "glGetShaderiv")
	case 208:
		errStr = fmt.Sprintf(getProcAddrErrStr, "glGetShaderInfoLog")
	case 209:
		errStr = fmt.Sprintf(getProcAddrErrStr, "glCreateProgram")
	case 210:
		errStr = fmt.Sprintf(getProcAddrErrStr, "glAttachShader")
	case 211:
		errStr = fmt.Sprintf(getProcAddrErrStr, "glLinkProgram")
	case 212:
		errStr = fmt.Sprintf(getProcAddrErrStr, "glValidateProgram")
	case 213:
		errStr = fmt.Sprintf(getProcAddrErrStr, "glGetProgramiv")
	case 214:
		errStr = fmt.Sprintf(getProcAddrErrStr, "glGetProgramInfoLog")
	case 215:
		errStr = fmt.Sprintf(getProcAddrErrStr, "glGenBuffers")
	case 216:
		errStr = fmt.Sprintf(getProcAddrErrStr, "glGenVertexArrays")
	case 217:
		errStr = fmt.Sprintf(getProcAddrErrStr, "glGetAttribLocation")
	case 218:
		errStr = fmt.Sprintf(getProcAddrErrStr, "glBindVertexArray")
	case 219:
		errStr = fmt.Sprintf(getProcAddrErrStr, "glEnableVertexAttribArray")
	case 220:
		errStr = fmt.Sprintf(getProcAddrErrStr, "glVertexAttribPointer")
	case 221:
		errStr = fmt.Sprintf(getProcAddrErrStr, "glBindBuffer")
	case 222:
		errStr = fmt.Sprintf(getProcAddrErrStr, "glBufferData")
	case 223:
		errStr = fmt.Sprintf(getProcAddrErrStr, "glGetVertexAttribPointerv")
	case 224:
		errStr = fmt.Sprintf(getProcAddrErrStr, "glUseProgram")
	case 225:
		errStr = fmt.Sprintf(getProcAddrErrStr, "glDeleteVertexArrays")
	case 226:
		errStr = fmt.Sprintf(getProcAddrErrStr, "glDeleteBuffers")
	case 227:
		errStr = fmt.Sprintf(getProcAddrErrStr, "glDeleteProgram")
	case 228:
		errStr = fmt.Sprintf(getProcAddrErrStr, "glDeleteShader")
	case 229:
		errStr = fmt.Sprintf(getProcAddrErrStr, "glGetUniformLocation")
	case 230:
		errStr = fmt.Sprintf(getProcAddrErrStr, "glUniformMatrix3fv")
	case 231:
		errStr = fmt.Sprintf(getProcAddrErrStr, "glUniformMatrix4fv")
	case 232:
		errStr = fmt.Sprintf(getProcAddrErrStr, "glUniformMatrix2x3fv")
	case 233:
		errStr = fmt.Sprintf(getProcAddrErrStr, "glGenerateMipmap")
	case 234:
		errStr = fmt.Sprintf(getProcAddrErrStr, "glActiveTexture")
	default:
		errStr = "unknown error"
	}
	errStr = errStr + " (" + strconv.FormatUint(g2dErrNum, 10)
	if win32ErrNum == 0 {
		errStr = errStr + ")"
	} else {
		errStr = errStr + ", " + strconv.FormatUint(win32ErrNum, 10) + ")"
	}
	if len(info) > 0 {
		errStr = errStr + "; " + info
	}
	return errors.New(errStr)
}

func (_ *tErrorHandler) LogError(err error) {
}
