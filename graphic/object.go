package graphic

import (
	"fmt"

	"github.com/telroshan/go-sfml/v2/graphics"
)

type childObject struct {
	x      int
	y      int
	width  int
	height int
	object *Object
}

type Object struct {
	// always available
	// these can be set upon creating the object
	// or the can be derived from children / texture
	// this is the native resolution, not the actual size this object will take up
	// this is defined to be the size of the texture if the texture exists
	Width  int
	Height int

	// position and size if there's a parent
	// children list OR texture
	// if texture != nil, then the children list are never considered
	texture   *Texture
	sprite    graphics.Struct_SS_sfSprite
	hasSprite bool
	children  []childObject
}

// Replace the existing texture with another texture:
// - The native resolution will get updated to the texture's resolution
// - If the texture's resolution is undefined (at least one dimension is <=0), then the resolution isn't updated
// - Generally, we don't have to care about the native resolution of a texture when loading from an image, so the resolution behaviour is for
// solid-filling or similar.
func (o *Object) LoadTexture(texture *Texture) {
	o.texture = texture
	width, height := texture.GetSize()
	if (width > 0) && (height > 0) {
		fmt.Printf("Updated size: (%d, %d) -> (%d, %d)\n", o.Width, o.Height, width, height)
		o.Width, o.Height = width, height
	}
}

func (o *Object) AnchorChild(child *Object, x, y, width, height int) {
	if o.texture != nil {
		panic("anchoring to a textured object")
	}
	o.children = append(o.children, childObject{
		x:      x,
		y:      y,
		width:  width,
		height: height,
		object: child,
	})
}

// anchor the child to a position, at native resolution
func (o *Object) AnchorChildNative(child *Object, x, y int) {
	// fmt.Println(o, child, x, y)
	o.AnchorChild(child, x, y, child.Width, child.Height)
}

// anchor the child to the center of this object, with the child's native resolution
// generally this is only used for composite object
func (o *Object) AnchorChildCenter(child *Object) {
	// x = o.Width - child.Width - x
	o.AnchorChildNative(child, (o.Width-child.Width)/2, (o.Height-child.Height)/2)
}

func (o *Object) draw(w *Window, x, y, width, height int) {
	// fmt.Println(x, y, width, height)
	if (width <= 0) || (height <= 0) {
		return
	}
	if o.texture != nil {
		// // has a texture, draw it
		// if o.hasSprite {
		// 	graphics.SfSprite_destroy(o.sprite)
		// } else {
		// 	o.hasSprite = true
		// }
		// o.sprite = graphics.SfSprite_create()
		// graphics.SfSprite_setTexture(o.sprite, o.texture.texture, 1)
		// graphics.SfSprite_setPosition(o.sprite, getVector2f(x, y))
		// graphics.SfRenderWindow_drawSprite(w.rw, o.sprite, (graphics.SfRenderStates)(graphics.SwigcptrSfRenderStates(0)))

		sprite := graphics.SfSprite_create()
		defer graphics.SfSprite_destroy(sprite)
		graphics.SfSprite_setTexture(sprite, o.texture.texture, 1)
		graphics.SfSprite_setPosition(sprite, getVector2f(x, y))
		graphics.SfSprite_setScale(sprite, getVector2ff(float32(width)/float32(o.Width), float32(height)/float32(o.Height)))
		graphics.SfRenderWindow_drawSprite(w.rw, sprite, (graphics.SfRenderStates)(graphics.SwigcptrSfRenderStates(0)))
		graphics.SfRenderTexture_drawSprite(w.rt, sprite, (graphics.SfRenderStates)(graphics.SwigcptrSfRenderStates(0)))
	} else {
		for _, child := range o.children {
			newX := x + child.x*width/o.Width
			newY := y + child.y*height/o.Height
			newWidth := child.width * width / o.Width
			newHeight := child.height * height / o.Height
			child.object.draw(w, newX, newY, newWidth, newHeight)
		}
	}
}
