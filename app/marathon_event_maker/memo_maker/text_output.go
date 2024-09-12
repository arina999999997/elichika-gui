package memo_maker

import (
	"elichika/gui/graphic"

	"fmt"

	"github.com/telroshan/go-sfml/v2/graphics"
)

type TextOutput struct {
	Maker  *MemoMaker
	Canvas *graphic.Canvas

	Background *graphic.Texture
	Text       *MultilineCenteredText

	Height int
	Width  int
}

func (t *TextOutput) SetSize(newWidth, newHeight int) {
	fmt.Println("New text output size: ", newWidth, newHeight)
	t.Width = newWidth
	t.Height = newHeight
}

// graphic.Object
func (t *TextOutput) GetWidth() int {
	return t.Width
}
func (t *TextOutput) GetHeight() int {
	return t.Height
}

func (t *TextOutput) InvalidateRenderCache() bool {
	return t.Canvas.InvalidateRenderCache()
}

func (t *TextOutput) ToTexture() *graphic.Texture {
	t.Draw()
	return t.Canvas.AsTexture()
}

func (t *TextOutput) Draw() {
	if t.Canvas.IsRendered() {
		return
	}
	t.SetSize(MemoWidth*t.Maker.TextScalingFactor, MemoHeight*t.Maker.TextScalingFactor)

	if t.Canvas == nil {
		t.Canvas = graphic.NewCanvas(t)
	} else {
		t.Canvas.UpdateSize(t)
	}

	if t.Maker.TextScaleBackground {
		t.Canvas.DrawTexture(t.Maker.BackgroundTexture, 0, 0, t.GetWidth(), t.GetHeight())
	}
	if t.Text == nil {
		t.Text = &MultilineCenteredText{}
	}

	t.Text.Load(t.Maker.Text,
		t.Maker.TextFont,
		t.Maker.TextCharacterSize,
		t.Maker.TextLetterSpacing,
		t.Maker.TextLineSpacing,
		t.Maker.TextColor,
		t.Maker.TextLineOffsets[:],
	)

	sprite := graphics.SfSprite_create()
	defer graphics.SfSprite_destroy(sprite)
	graphics.SfSprite_setTexture(sprite, t.Text.ToTexture().Texture, 1)
	graphics.SfSprite_setPosition(sprite, graphic.GetVector2f(t.Maker.TextX, t.Maker.TextY))
	graphics.SfSprite_setRotation(sprite, t.Maker.TextRotation)
	t.Canvas.DrawSprite(&graphic.Sprite{
		Sprite: sprite,
	})

	t.Canvas.Finalize()
}
