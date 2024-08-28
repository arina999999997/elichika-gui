package locale

import (
	"elichika/gui/graphic"
	"elichika/gui/graphic/button"
)

// A menu to change the locale
// it's a button that cycle through the available option
const (
	LocaleSelectorWidth  = 200
	LocaleSelectorHeight = 50
)

// TODO(gui): this doesn't invalidate every instance, although it might not be necessary if we only allow locale selecting on only one window.
func SelectorText() string {
	res := Language
	if Platform == "a" {
		res += " android"
	} else {
		res += " ios"
	}
	return res
}

func GetLocaleSelector(parent graphic.Object, onUpdate func()) *button.RectButton {
	selector := button.RectButton{
		Parent:  parent,
		Width:   LocaleSelectorWidth,
		Height:  LocaleSelectorHeight,
		Texture: graphic.RGBATexture(0x7f7f7fff),
	}
	selector.Text = graphic.NewText(&selector, SelectorText())
	selector.LeftClickHandler = func() {
		graphic.InvalidateRenderCache(&selector)
		for i := range Languages {
			if Language == Languages[i] {
				SetLanguage(Languages[(i+1)%len(Languages)])
				break
			}
		}
		selector.SetTextString(SelectorText())
		onUpdate()
	}
	selector.RightClickHandler = func() {
		graphic.InvalidateRenderCache(&selector)
		for i := range Platforms {
			if Platform == Platforms[i] {
				SetPlatform(Platforms[(i+1)%len(Platforms)])
				break
			}
		}
		selector.SetTextString(SelectorText())
		onUpdate()
	}
	return &selector
}
