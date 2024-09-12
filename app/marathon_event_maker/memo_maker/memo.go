package memo_maker

import (
	"elichika/gui/graphic"
)

// The text and background might be rendered at a higher resolution, but the output is always 200x200
type Memo struct {
	Maker  *MemoMaker
	Canvas *graphic.Canvas
}

// graphic.Object
func (m *Memo) GetHeight() int {
	return MemoWidth
}

func (m *Memo) GetWidth() int {
	return MemoHeight
}

func (m *Memo) InvalidateRenderCache() bool {
	return m.Canvas.InvalidateRenderCache()
}

func (m *Memo) ToTexture() *graphic.Texture {
	m.Draw()
	return m.Canvas.AsTexture()
}

func (m *Memo) Draw() {
	if m.Canvas.IsRendered() {
		return
	}
	graphic.DeleteCanvas(&m.Canvas) // TODO(gui): we can do some stuff to not have to do this every time
	if m.Canvas == nil {
		m.Canvas = graphic.NewCanvas(m)
	}

	if !m.Maker.TextScaleBackground {
		m.Canvas.DrawTexture(m.Maker.BackgroundTexture, 0, 0, m.GetWidth(), m.GetHeight())
	}
	// m.Canvas.DrawObject(m.Maker.TextOutput, 0, 0, m.GetWidth(), m.GetHeight())
	m.Canvas.DrawObject(m.Maker.TextOutput, 0, 0, 200, 200)

	m.Maker.PinTexture.StyleType = graphic.StyleTypeNone
	m.Canvas.DrawTexture(m.Maker.PinTexture, m.Maker.PinX, m.Maker.PinY, 1, 1)

	m.Canvas.Finalize()
}

func (m *Memo) GetParent() graphic.Object {
	return m.Maker
}
