package event_marathon

import (
	"elichika/client"

	"elichika/gui/graphic"
	"elichika/gui/sifas/asset"
	"elichika/gui/sifas/scene"
)

const (
	TopBackgroundX      = 0
	TopBackgroundY      = 0
	TopBackgroundWidth  = scene.GameWidth
	TopBackgroundHeight = scene.GameHeight
	TopBoardX           = 120
	TopBoardY           = 97
	TopBoardWidth       = 931
	TopBoardHeight      = 10000 // this extend out of the thingy
	TopTitleX           = 1150
	TopTitleY           = 134
	TopTitleWidth       = 460
	TopTitleHeight      = 163
)

// TODO(extra): these naming might differ from the name in code
// the objects are all exposed
type EventMarathonTopScene struct {
	BackgroundImage *graphic.Texture
	TitleImage      *graphic.Texture

	Board EventMarathonBoard

	// gui backend
	Canvas *graphic.Canvas
}

// own functions
func (emts *EventMarathonTopScene) Load(status client.EventMarathonTopStatus) error {
	var err error
	if status.TitleImagePath.V.HasValue {
		emts.TitleImage, err = asset.LoadTexture(status.TitleImagePath.V.Value)
		if err != nil {
			return err
		}
	} else {
		emts.TitleImage = nil
	}

	if status.BackgroundImagePath.V.HasValue {
		emts.BackgroundImage, err = asset.LoadTexture(status.BackgroundImagePath.V.Value)
		if err != nil {
			return err
		}
	} else {
		emts.BackgroundImage = nil
	}

	emts.Board.Load(status.BoardStatus)
	graphic.InvalidateRenderCache(emts)
	return nil
}

// graphic.Object

func (emts *EventMarathonTopScene) GetWidth() int {
	return scene.GameWidth
}

func (emts *EventMarathonTopScene) GetHeight() int {
	return scene.GameHeight
}

func (emts *EventMarathonTopScene) InvalidateRenderCache() bool {
	return emts.Canvas.InvalidateRenderCache()
}

func (emts *EventMarathonTopScene) Draw() {
	if emts.Canvas.IsRendered() {
		return
	}
	if emts.Canvas == nil {
		emts.Canvas = graphic.NewCanvas(emts)
	}

	if emts.BackgroundImage != nil {
		emts.Canvas.DrawTexture(emts.BackgroundImage, TopBackgroundX, TopBackgroundY, TopBackgroundWidth, TopBackgroundHeight)
	}

	if emts.TitleImage != nil {
		emts.Canvas.DrawTexture(emts.TitleImage, TopTitleX, TopTitleY, TopTitleWidth, TopTitleHeight)
	}

	emts.Canvas.DrawObject(&emts.Board, TopBoardX, TopBoardY, TopBoardWidth, TopBoardHeight)

	emts.Canvas.Finalize()
}

func (emts *EventMarathonTopScene) ToTexture() *graphic.Texture {
	emts.Draw()
	return emts.Canvas.AsTexture()
}
