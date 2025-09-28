/*
 *          Copyright 2025, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

package api

import (
	"errors"
	"runtime/debug"
	"strings"
	"strconv"
)

type MainLoop interface {
	InitMainLoop(params *g2d.InitParams)
	MainWindow() Window
	ErrMainLoop(err error)
}

type Application interface {
	InitApp() error
	FailApp(err error)
	CloseApp() error
}

func NewError(message string, id int64) error {
	var builder strings.Builder
	stack := debug.Stack()
	stack = skip5Lines(stack)
	builder.Grow(len(message) + len(stack))
	builder.WriteString(message)
	builder.WriteString(" (#")
	builder.WriteString(strconv.FormatInt(id, 10))
	builder.WriteString(")")
	for len(stack) > 0 {
		lineAEnd := seekByte('\n', stack)
		lineBEnd := seekByte('\n', stack[lineAEnd+1:]) + 1 + lineAEnd
		funcNameEnd := seekByteRev('(', stack[:lineAEnd])
		funcNameStart := seekByteRev('.', stack[:funcNameEnd]) + 1
		fileNameEnd := seekByteRev(':', stack[:lineBEnd])
		fileNameEndNum := seekNumEnd(stack[fileNameEnd+1:]) + 1 + fileNameEnd
		fileNameStart := seekByteRev('/', stack[:fileNameEnd]) + 1
		builder.WriteString("\n    ")
		builder.Write(stack[fileNameStart:fileNameEndNum])
		builder.WriteByte(' ')
		builder.Write(stack[funcNameStart:funcNameEnd])
		stack = stack[lineBEnd+1:]
	}
	return errors.New(builder.String())
}

func skip5Lines(bytes []byte) []byte {
	var counter int
	for i, b := range bytes {
		if b == '\n' {
			counter++
			if counter == 5 {
				return bytes[i+1:]
			}
		}
	}
	return bytes
}

func seekByte(target byte, bytes []byte) int {
	for i, b := range bytes {
		if b == target {
			return i
		}
	}
	return len(bytes)
}

func seekByteRev(target byte, bytes []byte) int {
	for i := len(bytes) - 1; i >= 0; i-- {
		if bytes[i] == target {
			return i
		}
	}
	return 0
}

func seekNumEnd(bytes []byte) int {
	for i, b := range bytes {
		if b < '0' || b > '9' {
			return i
		}
	}
	return len(bytes)
}

/*

ERROR: unknown state, hallo failed (#54211)
  main.go:14/bar
  main.go:20/func1

  main.go:14, bar
  main.go:20, func1

err.Append(api.NewError("hasdfasf asdf awef w", 10330))
err.AllString()

string(debug.Stack())

params := getParams()

if params.App {
	app := newApplication(params)
	audio.Init(app)
	wndapp.Start(app)
} else {
	info := getInfo(params)
	fmt.Println(info)
}

app api.Application
wnd api.Window

app abst.Application
wnd abst.Window
abst.NewError()

err := api.NewError()


batch0 := newBatch()

batch0.Enabled = false
x, y := batch0.XY(i)

batch0.SetShader(shader)
batch0.SetXY(i, x, y)
batch0.SetTex(0, tex)
batch0.SetSample(

texture.Bind(0)
shader.UseTexture(0, unit0)
shader.MaxSamples()

tex := newTexture()

SetLayer(squares)
batch0 := NewSquares()
batch0.EnsureCap(batch0, 10)
batch0 = batch0.Release()
length := beatch.Length()
for i := 0; i < length; i++ {
	
}
for i, j, length := 0, 0, len(batch0.Actives); i < length; i++ {
	if batch0.Active[i] {
		x, y, w := batch0.Values[j], batch0.Values[j+1], batch0.Values[j+2]
	}
	j += 3
}

Entities

cmdargs

g2d.Application
g2d

vb0.Application
vb1.Application
vb2.Application

p0.Application

exp001.Application
exp001.DefaultApplication

api.Application
api.DefaultApplication

win32ogl

awnd.Application
awnd.DefaultApplication
awa.Application
awa.DefaultApplication

api.Application
	ConfigApp
	InitApp
	FailApp
	CloseApp
api.ConfigAppParams

stack.NewError()

func ConfigApp(config *api.ConfigAppParams) {
	config.Window = true
	config.Audio = true
}

func InitApp(config *api.ConfigAppParams, err error) bool {
	if err == nil {
		return true
	}
	return false
}

func (ooo *App) ConfigApp(config *api.ConfigAppParams) {
	config.Window = true
	config.Audio = true
}

func (ooo *App) ConfigMainWindow(config *api.ConfigMainWindowParams) {
	config.Fullscreen = true
	config.Display = 2
}

api.NewError("hallo")

ERROR: unknown state, hallo failed (#54211)
  main.go/bar:14
  main.go/func1:20

  main.go:14, bar
  main.go:20, func1

err.Append(api.NewError("hasdfasf asdf awef w", 10330))
err.AllString()

string(debug.Stack())
*/

