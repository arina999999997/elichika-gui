package main

import (
	"elichika/gui/graphic"
	"elichika/gui/graphic/button"
	"elichika/gui/graphic/textbox"
	"elichika/gui/sifas/asset"
	"elichika/gui/sifas/locale"

	"fmt"
	// "os"
	// "time"
)

// this is mostly a testing program
// take args being the asset paths of texture
// clicking will cycle the assets

type AssetDisplay struct {
	AssetLabel    *graphic.Text
	AssetTextbox  *textbox.RectTextbox
	LoadButton    *button.RectButton
	LocaleButton  *button.RectButton
	ExtractButton *button.RectButton

	AssetTexture *graphic.Texture

	Canvas *graphic.Canvas
}

func (ad *AssetDisplay) InvalidateRenderCache() bool {
	return ad.Canvas.InvalidateRenderCache()
}

func (ad *AssetDisplay) ToTexture() *graphic.Texture {
	ad.Draw()
	return ad.Canvas.AsTexture()
}

func (ad *AssetDisplay) GetWidth() int {
	return 1800
}

func (ad *AssetDisplay) GetHeight() int {
	return 1000
}

func (ad *AssetDisplay) Draw() {
	if ad.Canvas.IsRendered() {
		return
	}
	if ad.Canvas == nil {
		ad.Canvas = graphic.NewCanvas(ad)
	}

	ad.Canvas.DrawObject(ad.AssetLabel, 0, 0, ad.AssetLabel.GetWidth(), 50)
	ad.Canvas.DrawObject(ad.AssetTextbox, ad.AssetLabel.GetWidth(), 0, ad.AssetTextbox.GetWidth(), 50)
	ad.Canvas.DrawObject(ad.LoadButton, 405, 0, ad.LoadButton.GetWidth(), 50)
	ad.Canvas.DrawObject(ad.LocaleButton, 610, 0, ad.LocaleButton.GetWidth(), 50)
	ad.Canvas.DrawObject(ad.ExtractButton, 815, 0, ad.ExtractButton.GetWidth(), 50)

	if ad.AssetTexture != nil {
		w, h := ad.AssetTexture.GetSize()
		ad.Canvas.DrawTexture(ad.AssetTexture, 0, 100, w, h)
	}

	ad.Canvas.Finalize()

}

func (ad *AssetDisplay) ForEach(f func(graphic.Object)) {
	f(ad.AssetLabel)
	f(ad.AssetTextbox)
	f(ad.LoadButton)
	f(ad.LocaleButton)
	f(ad.ExtractButton)
}

func (ad *AssetDisplay) MapEvent(e graphic.Event, child graphic.Object) graphic.Event {
	switch e.(type) {
	case graphic.MouseButtonDownEvent:
		mouseEvent := graphic.MouseButtonDownEvent{}
		mouseEvent = e.(graphic.MouseButtonDownEvent)
		if child == ad.AssetLabel {
			mouseEvent.X -= 0
		} else if child == ad.AssetTextbox {
			mouseEvent.X -= ad.AssetLabel.GetWidth()
		} else if child == ad.LoadButton {
			mouseEvent.X -= 405
		} else if child == ad.LocaleButton {
			mouseEvent.X -= 610
		} else {
			mouseEvent.X -= 815
		}
		return mouseEvent
	default:
		return e
	}
}

func main() {
	window, err := graphic.NewWindow("Texture viewer")
	if err != nil {
		panic(err)
	}
	object := &AssetDisplay{}
	object.AssetLabel, object.AssetTextbox = textbox.NewLabelAndRectTextbox(object, 400, 50, "Asset path: ")

	object.AssetTextbox.Texture = graphic.RGBATexture(0x0f0f0fff)
	object.AssetTextbox.FocusTexture = graphic.RGBATexture(0x7f7f7fff)

	object.AssetTextbox.OnTextUpdateFunc = func() {
		if len(object.AssetTextbox.TextContent) > 3 {
			graphic.InvalidateRenderCache(object.AssetTextbox)
			object.AssetTextbox.TextContent = object.AssetTextbox.TextContent[:3]
		}
	}

	displayAsset := func() {
		texture, err := asset.LoadTexture(object.AssetTextbox.TextContent)
		graphic.InvalidateRenderCache(object)
		if err != nil {
			// TODO(gui): Dialog box
			fmt.Println("Failed to load asset: ", object.AssetTextbox.TextContent, ", error: ", err)
			return
		}
		object.AssetTexture = texture
	}

	object.LocaleButton = locale.GetLocaleSelector(object, displayAsset)

	object.LoadButton = &button.RectButton{
		Width:   200,
		Height:  50,
		Texture: graphic.RGBATexture(0x7f7f7fff),
	}
	object.LoadButton.Text = graphic.NewText(object.LoadButton, "Load asset")
	object.LoadButton.LeftClickHandler = displayAsset
	object.AssetTextbox.OnEnterFunc = displayAsset

	object.ExtractButton = &button.RectButton{
		Width:   200,
		Height:  50,
		Texture: graphic.RGBATexture(0x7f7f7fff),
	}
	object.ExtractButton.Text = graphic.NewText(object.ExtractButton, "Extract")
	object.ExtractButton.LeftClickHandler = func() {
		if object.AssetTexture != nil {
			object.AssetTexture.SaveToImage(locale.SelectorText() + ".png")
		}
	}

	window.SetObject(object)
	window.Display()
}
