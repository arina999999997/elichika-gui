package graphic

import (
	// "fmt"

	"github.com/telroshan/go-sfml/v2/graphics"
)

// Canvas can be used to draw object on top
// Each object might also store a canvas which is what it looks like.

type sfRenderTexture = graphics.Struct_SS_sfRenderTexture

type Canvas struct {
	rendered      bool
	renderTexture sfRenderTexture
}

func (c *Canvas) IsRendered() bool {
	return (c != nil) && c.rendered
}

func (c *Canvas) InvalidateRenderCache() bool {
	if (c == nil) || (c.rendered == false) {
		return false
	}
	c.rendered = false
	graphics.SfRenderTexture_clear(c.renderTexture, graphics.GetSfTransparent())
	return true
}

func (c *Canvas) Finalize() {
	graphics.SfRenderTexture_display(c.renderTexture)
	c.rendered = true
}

var commonRenderState graphics.SfRenderStates
var loadedCommonState = false

func GetRenderState() graphics.SfRenderStates {
	return (graphics.SfRenderStates)(graphics.SwigcptrSfRenderStates(0))
	if loadedCommonState {
		return commonRenderState
	}
	loadedCommonState = true
	commonRenderState = graphics.NewSfRenderStates()

	blendMode := graphics.GetSfBlendAlpha()
	// defer graphics.DeleteSfBlendMode(blendMode)

	// blendMode.SetColorSrcFactor(6) // BlendMode::Factor::SrcAlpha
	// // blendMode.SetColorDstFactor(7) // BlendMode::Factor::OneMinusSrcAlpha
	// blendMode.SetColorDstFactor(1) // BlendMode::Factor::One
	// blendMode.SetColorEquation(0) // BlendMode::Equation::Add
	// blendMode.SetAlphaSrcFactor(1) // BlendMode::Factor::One
	// blendMode.SetAlphaDstFactor(1) // BlendMode::Factor::One
	// // blendMode.SetAlphaDstFactor(7) // BlendMode::Factor::OneMinusSrcAlpha
	// blendMode.SetAlphaEquation(0) // BlendMode::Equation::Add

	commonRenderState.SetBlendMode(blendMode)

	// shader := graphics.SfShader_createFromFile("empty", "empty", "empty")
	// renderState.SetShader(shader)

	return commonRenderState
}

// draw a texture in a rectangle [x, x + width) x [y, y + height) directly
// note that width < 0 or height < 0 is not
func (c *Canvas) DrawTexture(texture *Texture, x, y, width, height int) {
	if c.IsRendered() {
		panic("drawing on finalized canvas")
	}
	if (width <= 0) || (height <= 0) {
		return
	}
	tWidth, tHeight := texture.GetSize()
	wScale := float32(width) / float32(tWidth)
	hScale := float32(height) / float32(tHeight)

	// calculate the fit scale
	switch texture.StyleType {
	case StyleTypeNone:
		fallthrough
	case StyleTypeRepeat:
		wScale = 1
		hScale = 1
	case StyleTypeFitWidth:
		hScale = wScale
	case StyleTypeFitHeight:
		wScale = hScale
	case StyleTypeFitContainer:
		if wScale > hScale {
			wScale = hScale
		} else {
			hScale = wScale
		}
	case StyleTypeAutoCentered:
		if wScale > hScale {
			wScale = hScale
		} else {
			hScale = wScale
		}
		// fmt.Println(wScale, hScale)
		if wScale > 1 {
			wScale = 1
			hScale = 1
		}
	case StyleTypeIndependent:
	default:
		panic("Unknown style type")
	}
	// fmt.Println(width, height, tWidth, tHeight, wScale, hScale)
	graphics.SfTexture_setSmooth(texture.Texture, 1)
	sprite := graphics.SfSprite_create()
	defer graphics.SfSprite_destroy(sprite)

	graphics.SfSprite_setTexture(sprite, texture.Texture, 1)
	graphics.SfSprite_setScale(sprite, GetVector2ff(wScale, hScale))
	if texture.StyleType == StyleTypeAutoCentered {
		nX := float32(x) + (float32(width)-float32(tWidth)*wScale)/2
		nY := float32(y) + (float32(height)-float32(tHeight)*hScale)/2
		graphics.SfSprite_setPosition(sprite, GetVector2ff(nX, nY))
	} else {
		graphics.SfSprite_setPosition(sprite, GetVector2f(x, y))
	}

	if texture.StyleType == StyleTypeRepeat {
		repeated := graphics.SfTexture_isRepeated(texture.Texture)
		if repeated == 0 {
			graphics.SfTexture_setRepeated(texture.Texture, 1)
		}
		graphics.SfSprite_setTextureRect(sprite, GetIntRect(width, height))
		graphics.SfRenderTexture_drawSprite(c.renderTexture, sprite, GetRenderState())
		if repeated == 0 {
			graphics.SfTexture_setRepeated(texture.Texture, 0)
		}
	} else {
		graphics.SfRenderTexture_drawSprite(c.renderTexture, sprite, GetRenderState())
	}
}

func (c *Canvas) DrawSprite(sprite *Sprite) {
	graphics.SfRenderTexture_drawSprite(c.renderTexture, sprite.Sprite, GetRenderState())
}

// this is used to draw a text directly
// it is different from drawing from an object that happen to be a text
// generally, this is used to draw to text's canvas, then that canvas would be used to draw on other things
func (c *Canvas) DrawText(text *Text) {
	if c.IsRendered() {
		panic("drawing on finalized canvas")
	}
	graphics.SfRenderTexture_drawText(c.renderTexture, text.text, (graphics.SfRenderStates)(graphics.SwigcptrSfRenderStates(0)))
}

func (c *Canvas) DrawObject(object Object, x, y, width, height int) {
	c.DrawTexture(object.ToTexture(), x, y, width, height)
}

func (c *Canvas) AsTexture() *Texture {
	if !c.IsRendered() {
		panic("Cannot create texture from not rendered canvas")
	}
	// note that there is no memory leak, because the texture in the C-side is a const reference
	// which also mean it's possible to keep a reference and reuse them
	texture := Texture{}
	texture.Texture = graphics.SfRenderTexture_getTexture(c.renderTexture)
	return &texture
}

func (c *Canvas) UpdateSize(object Object) {
	size := graphics.SfRenderTexture_getSize(c.renderTexture)
	if (size.GetX() == uint(object.GetWidth())) && (size.GetY() == uint(object.GetHeight())) {
		return
	}
	graphics.SfRenderTexture_destroy(c.renderTexture)
	c.renderTexture = graphics.SfRenderTexture_createWithSettings(uint(object.GetWidth()), uint(object.GetHeight()), GetContextSetting())
}

func (c *Canvas) UpdateSizeDirect(w, h int) {
	size := graphics.SfRenderTexture_getSize(c.renderTexture)
	if (size.GetX() == uint(w)) && (size.GetY() == uint(h)) {
		return
	}
	graphics.SfRenderTexture_destroy(c.renderTexture)
	c.renderTexture = graphics.SfRenderTexture_createWithSettings(uint(w), uint(h), GetContextSetting())
}

func (c *Canvas) UpdateSizeFromTexture(texture *Texture) {
	size := graphics.SfRenderTexture_getSize(c.renderTexture)
	newWidth, newHeight := texture.GetSize()
	if (size.GetX() == uint(newWidth)) && (size.GetY() == uint(newHeight)) {
		return
	}
	graphics.SfRenderTexture_destroy(c.renderTexture)
	c.renderTexture = graphics.SfRenderTexture_createWithSettings(uint(newWidth), uint(newHeight), GetContextSetting())
}

func NewCanvas(object Object) *Canvas {
	canvas := Canvas{}
	canvas.renderTexture = graphics.SfRenderTexture_createWithSettings(uint(object.GetWidth()), uint(object.GetHeight()), GetContextSetting())
	graphics.SfRenderTexture_setSmooth(canvas.renderTexture, 1)
	return &canvas
}

func NewCanvasNoSmooth(object Object) *Canvas {
	canvas := Canvas{}
	canvas.renderTexture = graphics.SfRenderTexture_createWithSettings(uint(object.GetWidth()), uint(object.GetHeight()), GetContextSetting())
	// graphics.SfRenderTexture_setSmooth(canvas.renderTexture, 1)
	return &canvas
}

func DeleteCanvas(canvas **Canvas) {
	if *canvas == nil {
		return
	}
	graphics.SfRenderTexture_destroy((*canvas).renderTexture)
	*canvas = nil
}
