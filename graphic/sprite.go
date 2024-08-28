package graphic

import (
	"github.com/telroshan/go-sfml/v2/graphics"
)

// A sprite is an instruction of how to draw a texture
// for now it just wrap the smfl sprite type
type Sprite struct {
	Sprite graphics.Struct_SS_sfSprite
}
