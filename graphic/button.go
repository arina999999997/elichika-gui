package graphic

// A generic button type:
// - Background texture for the button
// - Text for the text on the button
// - OnClick for the function to call when the button is actually clicked
type Button interface {
	SetTexture(*Texture)
	SetText(*Text)
	Clickable
}
