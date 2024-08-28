package memo_maker
import (
	"elichika/gui/graphic"
	
	"github.com/telroshan/go-sfml/v2/graphics"
)
type Memo struct {
	Parent *MemoMaker
	Canvas *graphic.Canvas

	Background *graphic.Texture

	Pin *graphic.Texture
	PinX int
	PinY int

	Text *MultilineCenteredText
}

// graphic.Object
func (*Memo) GetHeight() int {
	return 200 * 4
}

func (*Memo) GetWidth() int {
	return 200 * 4
}

func (m *Memo) InvalidateRenderCache() bool {
	return m.Canvas.InvalidateRenderCache()
}

func  (m *Memo) ToTexture() *graphic.Texture {
	m.Draw()
	return m.Canvas.AsTexture()
}

func (m *Memo) Draw() {
	if m.Canvas.IsRendered() {
		return
	}
	if m.Canvas == nil {
		m.Canvas = graphic.NewCanvas(m)
	}
	
	m.Background.StyleType = graphic.StyleTypeNone
	{
		
		sprite := graphics.SfSprite_create()
		defer graphics.SfSprite_destroy(sprite)
		graphics.SfSprite_setTexture(sprite, m.Background.Texture, 1)
		graphics.SfSprite_setScale(sprite, graphic.GetVector2ff(4, 4))
		m.Canvas.DrawSprite(&graphic.Sprite{
			Sprite: sprite,
		})
	}
	// m.Canvas.DrawTexture(m.Background, 0, 0, 200 * 4, 200 * 4)
	// m.Canvas.DrawTexture(m.Background, 0, 0, 200 * 4, 200 * 4)
	// TODO(pin)
	// m.Pin.StyleType = graphic.StyleTypeNone
	// m.Canvas.DrawTexture(m.Pin, m.PinX, m.PinY, 200 * 4, 200 * 4)

	
	sprite := graphics.SfSprite_create()
	defer graphics.SfSprite_destroy(sprite)
	graphics.SfSprite_setTexture(sprite, m.Text.ToTexture().Texture, 1)
	graphics.SfSprite_setPosition(sprite, graphic.GetVector2f(m.Parent.TextPicker.X, m.Parent.TextPicker.Y))
	graphics.SfSprite_setRotation(sprite, m.Parent.TextPicker.Rotation)
	scale := float32(1) / float32(m.Parent.TextPicker.DownscalingFactor)
	graphics.SfSprite_setScale(sprite, graphic.GetVector2ff(scale, scale))
	m.Canvas.DrawSprite(&graphic.Sprite{
		Sprite: sprite,
	})

	m.Canvas.Finalize()
}

func (m *Memo) GetParent() graphic.Object {
	return m.Parent
}

func (m *Memo) Load(parent *MemoMaker) {
	m.Parent = parent
	m.Background = parent.BackgroundPicker.Texture
	m.Pin = parent.PinPicker.Texture
	m.PinX = parent.PinPicker.X
	m.PinY = parent.PinPicker.Y
	if m.Text == nil {
		m.Text = &MultilineCenteredText{}
	}
	m.Text.Load(parent.TextPicker.Text, 
		parent.TextPicker.Font,
		parent.TextPicker.CharacterSize,
		parent.TextPicker.LetterSpacing,
		parent.TextPicker.LineSpacing,
		parent.TextPicker.Color,
		parent.TextPicker.LineOffsets,
	)
}