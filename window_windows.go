/*
 *          Copyright 2024, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

package g2d

// #include "g2d.h"
import "C"

type Window interface {
	/*
	   OnConfig(config *Configuration) error
	   OnCreate(widget *Widget) error
	   OnShow() error
	   OnResize() error
	   OnKeyDown(keyCode int, repeated uint) error
	   OnKeyUp(keyCode int) error
	   OnTextureLoaded(textureId int) error
	   OnUpdate() error
	   OnClose() (bool, error)
	   OnDestroy() error
	*/
}
