package graphic

import (
	"github.com/telroshan/go-sfml/v2/graphics"
	"github.com/telroshan/go-sfml/v2/window"
)

var v2u graphics.SfVector2u

func init() {
	v2u = graphics.NewSfVector2u()
}
func getVector2u(x, y int) graphics.SfVector2u {
	v2u.SetX(uint(x))
	v2u.SetY(uint(y))
	return v2u
}

var v2f graphics.SfVector2f

func init() {
	v2f = graphics.NewSfVector2f()
}
func GetVector2f(x, y int) graphics.SfVector2f {
	v2f.SetX(float32(x))
	v2f.SetY(float32(y))
	return v2f
}

func GetVector2ff(x, y float32) graphics.SfVector2f {
	v2f.SetX(float32(x))
	v2f.SetY(float32(y))
	return v2f
}

var intRect graphics.SfIntRect

func init() {
	intRect = graphics.NewSfIntRect()
}
func GetIntRect(w, h int) graphics.SfIntRect {
	intRect.SetWidth(w)
	intRect.SetHeight(h)
	return intRect
}

func GetIntRectWithPosition(w, h, x, y int) graphics.SfIntRect {
	intRect.SetWidth(w)
	intRect.SetHeight(h)
	intRect.SetLeft(x)
	intRect.SetTop(y)
	return intRect
}

var contextSetting window.SfContextSettings

func init() {
	contextSetting = window.NewSfContextSettings()
	contextSetting.SetAntialiasingLevel(16)
}
func GetContextSetting() window.SfContextSettings {
	// contextSetting := window.NewSfContextSettings()
	// contextSetting.SetAntialiasingLevel(16)
	return contextSetting
}
