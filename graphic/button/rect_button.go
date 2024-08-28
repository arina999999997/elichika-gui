package button

import (
	"elichika/gui/graphic"
)

// TODO(refactor): Move to separate file
// a rectangle button
// so click on the bounding box to click the button
type RectButton struct {
	Parent graphic.Object

	Width  int
	Height int

	Texture *graphic.Texture
	Text    *graphic.Text

	Canvas *graphic.Canvas

	LeftClickHandler  func()
	RightClickHandler func()
}

// Own functions

func (rb *RectButton) SetTexture(texture *graphic.Texture) {
	if texture == rb.Texture {
		return
	}
	rb.Texture = texture
	graphic.InvalidateRenderCache(rb)
}

func (rb *RectButton) SetText(text *graphic.Text) {
	if text == rb.Text {
		return
	}
	rb.Text = text
	rb.Text.Parent = rb
	graphic.InvalidateRenderCache(rb)
}

func (rb *RectButton) SetTextString(text string) {
	rb.Text.SetText(text)
}

// Focusable

func (*RectButton) SetFocus() {
}

func (*RectButton) UnsetFocus() {
}

func (*RectButton) HasFocus() bool {
	return false
}

// Clickable interface
func (rb *RectButton) OnClick(w *graphic.Window, event graphic.MouseButtonDownEvent) bool {
	if (event.X < 0) || (event.X >= rb.Width) {
		return false
	}
	if (event.Y < 0) || (event.Y >= rb.Height) {
		return false
	}

	// button doesn't have a focus functionality for now, but we still take away focus from other object
	w.SetFocusObject(rb)

	// no matter if the click is handled or not, the event will be considered handled

	if (event.Button == graphic.MouseButtonLeft) && (rb.LeftClickHandler != nil) {
		rb.LeftClickHandler()
	}

	if (event.Button == graphic.MouseButtonRight) && (rb.RightClickHandler != nil) {
		rb.RightClickHandler()
	}

	return true
}

// Button interface
func (rb *RectButton) GetWidth() int {
	return rb.Width
}

func (rb *RectButton) GetHeight() int {
	return rb.Height
}

// Object interface
func (rb *RectButton) InvalidateRenderCache() bool {
	return rb.Canvas.InvalidateRenderCache()
}

func (rb *RectButton) Draw() {
	if rb.Canvas.IsRendered() {
		return
	}
	if rb.Canvas == nil {
		rb.Canvas = graphic.NewCanvas(rb)
	}

	// draw the background, then the text
	if rb.Texture != nil {
		rb.Canvas.DrawTexture(rb.Texture, 0, 0, rb.Width, rb.Height)
	}

	if rb.Text != nil {
		rb.Text.SetHeight(rb.Height)
		textTexture := rb.Text.ToTexture()
		textTexture.StyleType = graphic.StyleTypeAutoCentered
		// textTexture.StyleType = graphic.StyleTypeNone
		// rb.Canvas.DrawObject(rb.Text, 0, 0, rb.Width, rb.Height)
		rb.Canvas.DrawTexture(textTexture, 0, 0, rb.Width, rb.Height)
	}
	rb.Canvas.Finalize()
}

func (rb *RectButton) ToTexture() *graphic.Texture {
	rb.Draw()
	texture := rb.Canvas.AsTexture()
	return texture
}

func NewButton(parent graphic.Object, width, height int, text string, leftClickHandler func(), rightClickHandler func()) *RectButton {
	button := &RectButton{
		Parent: parent,
		Width: width,
		Height: height,
		Texture: graphic.RGBATexture(0x7f7f7fff),
		LeftClickHandler: leftClickHandler,
		RightClickHandler: rightClickHandler,
	}
	button.Text = graphic.NewText(button, text)
	return button
}