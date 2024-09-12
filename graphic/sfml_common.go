package graphic

import (
	"github.com/telroshan/go-sfml/v2/graphics"
	"github.com/telroshan/go-sfml/v2/window"
)

func getVector2u(x, y int) graphics.SfVector2u {
	vec := graphics.NewSfVector2u()
	vec.SetX(uint(x))
	vec.SetY(uint(y))
	return vec
}

func GetVector2f(x, y int) graphics.SfVector2f {
	vec := graphics.NewSfVector2f()
	vec.SetX(float32(x))
	vec.SetY(float32(y))
	return vec
}

func GetVector2ff(x, y float32) graphics.SfVector2f {
	vec := graphics.NewSfVector2f()
	vec.SetX(float32(x))
	vec.SetY(float32(y))
	return vec
}

func GetIntRect(w, h int) graphics.SfIntRect {
	rect := graphics.NewSfIntRect()
	rect.SetWidth(w)
	rect.SetHeight(h)
	return rect
}

func GetContextSetting() window.SfContextSettings {
	contextSetting := window.NewSfContextSettings()
	// contextSetting.SetAntialiasingLevel(16)
	contextSetting.SetAntialiasingLevel(16)
	return contextSetting
}
