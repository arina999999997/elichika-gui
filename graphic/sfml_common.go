package graphic

import (
	"github.com/telroshan/go-sfml/v2/graphics"
)

func getVector2u(x, y int) graphics.SfVector2u {
	vec := graphics.NewSfVector2u()
	vec.SetX(uint(x))
	vec.SetY(uint(y))
	return vec
}

func getVector2f(x, y int) graphics.SfVector2f {
	vec := graphics.NewSfVector2f()
	vec.SetX(float32(x))
	vec.SetY(float32(y))
	return vec
}

func getVector2ff(x, y float32) graphics.SfVector2f {
	vec := graphics.NewSfVector2f()
	vec.SetX(float32(x))
	vec.SetY(float32(y))
	return vec
}
