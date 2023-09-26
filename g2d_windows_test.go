/*
 *          Copyright 2023, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

package g2d

import (
	"testing"
)

func TestInit(t *testing.T) {
	var engine Engine
	err := engine.Init()
	if err != nil {
		t.Error(err.Error())
	}
}

/*
func BenchmarkInit(b *testing.B) {
	var collection cdata.Collection
	var loader ogfl.Loader
	var rects modules.Rectangles
	collection.Passes = 2
	collection.Init(&loader, &rects)
	var engine Engine
	params := new(EngineParams)
	params.Modules = append(params.Modules, new(modules.Rectangles))
	engine.Init(params)
}
*/

/*
type tTestConfigWindow struct {
	DefaultWindow
	configCalled  bool
	destroyCalled bool
}

func (window *tTestConfigWindow) OnConfig(config *Configuration) error {
	config.AutoUpdate = false
	window.configCalled = true
	return errors.New("test")
}

func (window *tTestConfigWindow) OnDestroy() {
	window.destroyCalled = true
}

type tTestCreateWindow struct {
	tTestConfigWindow
	createCalled bool
}

func (window *tTestCreateWindow) OnConfig(config *Configuration) error {
	config.AutoUpdate = false
	window.configCalled = true
	return nil
}

func (window *tTestCreateWindow) OnCreate(widget *Widget) error {
	window.createCalled = true
	return errors.New("test")
}

func TestInit(t *testing.T) {
	if !initialized {
		Init()
		errs := Errors()
		if len(errs) > 0 {
			t.Error(errs[0].Error())
		}
	}
}

func TestConfigWindow(t *testing.T) {
	if !initialized {
		Init()
	}
	if !initFailed {
		clearErrors()
		window := new(tTestConfigWindow)
		Show(window)
		errs = Errors()
		if !window.configCalled {
			t.Error("OnConfig not called")
		} else if !window.destroyCalled {
			t.Error("OnDestroy not called")
		} else if len(errs) != 1 {
			t.Error(len(errs))
		} else if errs[0].Error() != "test" {
			t.Error(errs[0].Error())
		}
	} else {
		t.Error(errs[0].Error())
	}
}

func TestCreateWindow(t *testing.T) {
	if !initialized {
		Init()
	}
	if !initFailed {
		window := new(tTestCreateWindow)
		clearErrors()
		Show(window)
		errs = Errors()
		if !window.configCalled {
			t.Error("OnConfig not called")
		} else if !window.createCalled {
			t.Error("OnCreate not called")
		} else if !window.destroyCalled {
			t.Error("OnDestroy not called")
		} else if len(errs) != 1 {
			t.Error(len(errs))
		} else if errs[0].Error() != "test" {
			t.Error(errs[0].Error())
		}
	} else {
		t.Error(errs[0].Error())
	}
}

func TestRectIndex(t *testing.T) {
	layer := newRectLayer(2)
	indexA := layer.newRectIndex()
	indexB := layer.newRectIndex()
	if indexA != 0 || indexB != 1 {
		t.Error(indexA, indexB)
	} else {
		layer.release(indexA)
		indexA = layer.newRectIndex()
		indexB = layer.newRectIndex()
		if layer.totalActive != 3 {
			t.Error(layer.totalActive)
		} else if indexA != 0 || indexB != 2 {
			t.Error(indexA, indexB)
		}
	}
}

func TestLayerA(t *testing.T) {
	layer := newRectEntitiesLayer(2)
	rectA := layer.newRectEntity(nil, 0, 0)
	rectB := layer.newRectEntity(nil, 0, 1)
	if rectA.chunk != 0 || rectB.chunk != 0 {
		t.Error(rectA.chunk, rectB.chunk)
	} else if rectA.index != 0 || rectB.index != 1 {
		t.Error(rectA.index, rectB.index)
	} else {
		layer.release(rectA.chunk, rectA.entityIndex)
		rectA = layer.newRectEntity(nil, 0, 0)
		rectB = layer.newRectEntity(nil, 0, 2)
		if layer.size != 2 {
			t.Error(layer.size)
		} else if rectA.layer != 0 || rectB.layer != 0 {
			t.Error(rectA.layer, rectB.layer)
		} else if rectA.chunk != 0 || rectB.chunk != 1 {
			t.Error(rectA.chunk, rectB.chunk)
		} else if rectA.index != 0 || rectB.index != 2 {
			t.Error(rectA.index, rectB.index)
		}
	}
}

func TestLayerB(t *testing.T) {
	layer := newRectEntitiesLayer(2)
	layer.newRectEntity(nil, 0, 0)
	rectB := layer.newRectEntity(nil, 0, 1)
	layer.newRectEntity(nil, 0, 2)
	rectD := rectB
	layer.release(rectB.chunk, rectB.entityIndex)
	rectB = layer.newRectEntity(nil, 0, 3)
	if rectD != rectB {
		t.Error(rectD.layer, rectD.chunk, rectD.entityIndex, rectB.layer, rectB.chunk, rectB.entityIndex)
	}
}

func TestSwitchBuffer(t *testing.T) {
	var gfx Graphics
	gfx.rBuffer = &gfx.buffers[0]
	gfx.wBuffer = &gfx.buffers[0]
	gfx.SetBGColor(1.0, 0.0, 0.0)
	gfx.NewRectLayer(100)
	if len(gfx.wBuffer.layers) != 1 || len(gfx.entitiesLayers) != 1 {
		t.Error(len(gfx.wBuffer.layers), len(gfx.entitiesLayers))
	} else if gfx.wBuffer.bgR != 1.0 || gfx.wBuffer.bgG != 0.0 {
		t.Error(gfx.wBuffer.bgR, gfx.wBuffer.bgG)
	} else {
		gfx.SetBGColor(0.0, 1.0, 0.0)
		gfx.switchWBuffer()
		if len(gfx.wBuffer.layers) != 1 || len(gfx.entitiesLayers) != 1 {
			t.Error(len(gfx.wBuffer.layers), len(gfx.entitiesLayers))
		} else if gfx.wBuffer.bgR != 0.0 || gfx.wBuffer.bgG != 1.0 {
			t.Error(gfx.wBuffer.bgR, gfx.wBuffer.bgG)
		}
		gfx.switchWBuffer()
		if len(gfx.wBuffer.layers) != 1 || len(gfx.entitiesLayers) != 1 {
			t.Error(len(gfx.wBuffer.layers), len(gfx.entitiesLayers))
		}
		gfx.NewRectLayer(100)
		if len(gfx.wBuffer.layers) != 2 || len(gfx.entitiesLayers) != 2 {
			t.Error(len(gfx.wBuffer.layers), len(gfx.entitiesLayers))
		} else if gfx.wBuffer.bgR != 0.0 || gfx.wBuffer.bgG != 1.0 {
			t.Error(gfx.wBuffer.bgR, gfx.wBuffer.bgG)
		}
		gfx.switchRBuffer()
		if gfx.rBuffer.bgR != 0.0 || gfx.rBuffer.bgG != 1.0 {
			t.Error(gfx.rBuffer.bgR, gfx.rBuffer.bgG)
		}
	}
}
*/
