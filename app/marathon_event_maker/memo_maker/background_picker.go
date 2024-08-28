package memo_maker

import (
	"elichika/gui/graphic"
	"elichika/gui/graphic/textbox"

	"fmt"
	"errors"
)

type BackgroundPicker struct {
	// this is the background of the memo, so no attaching pin included
	Parent graphic.Object
	Canvas *graphic.Canvas

	// default to gui/app/marathon_event_maker/memo_maker/common_note.png
	Label *graphic.Text
	Textbox *textbox.RectTextbox

	// the background must have size 200x200, and must already include every transformation necessary on it
	LoadedPath string
	Texture *graphic.Texture


}

func (bp *BackgroundPicker) Load() error {
	if bp.Textbox.TextContent == bp.LoadedPath {
		return nil
	}
	texture, err := graphic.FileTexture(bp.Textbox.TextContent)
	if err != nil {
		return errors.New(fmt.Sprint("Error while loading background texture: ", err))
	}
	bp.LoadedPath = bp.Textbox.TextContent
	bp.Texture = texture
	return nil
}

// graphic.Object
func (*BackgroundPicker) GetWidth() int {
	return MemoMakerWidth
}
func (*BackgroundPicker) GetHeight() int {
	return 50
}

func (bp *BackgroundPicker) InvalidateRenderCache() bool {
	return bp.Canvas.InvalidateRenderCache()
}

func (bp *BackgroundPicker) ToTexture() *graphic.Texture {
	bp.Draw()
	return bp.Canvas.AsTexture()
}

func (bp *BackgroundPicker) Draw() {
	if bp.Canvas.IsRendered() {
		return
	}
	if bp.Canvas == nil {
		bp.Canvas = graphic.NewCanvas(bp)
	}

	bp.Canvas.DrawObject(bp.Label, 0, 0, bp.Label.GetWidth(), 50)
	bp.Canvas.DrawObject(bp.Textbox, bp.Label.GetWidth(), 0, bp.Textbox.GetWidth(), 50)

	bp.Canvas.Finalize()
}

// graphic.CompositeObject
func (bp *BackgroundPicker) ForEach(f func(graphic.Object)) {
	f(bp.Label)
	f(bp.Textbox)
}

// graphic.PollingCompositeObject
func (bp *BackgroundPicker) MapEvent(e graphic.Event, child graphic.Object) graphic.Event {
	switch e.(type) {
	case graphic.MouseButtonDownEvent:
		mouseEvent := graphic.MouseButtonDownEvent{}
		mouseEvent = e.(graphic.MouseButtonDownEvent)
		if child == bp.Label {
			mouseEvent.X -= 0
		} else if child == bp.Textbox {
			mouseEvent.X -= bp.Label.GetWidth()
		} else {
			mouseEvent.X -= bp.Label.GetWidth() + bp.Textbox.GetWidth() + 5
		}
		return mouseEvent
	default:
		return e
	}
}

// graphic.ChildObject
func (bp *BackgroundPicker) GetParent() graphic.Object {
	return bp.Parent
}

func NewBackgroundPicker(parent graphic.Object) *BackgroundPicker{
	picker := &BackgroundPicker{
		Parent: parent,
	}
	picker.Label, picker.Textbox = textbox.NewLabelAndRectTextbox(picker, MemoMakerWidth, 50, "Memo background: ")
	picker.Textbox.SetText("gui/app/marathon_event_maker/memo_maker/common_note.png")
	return picker
}
