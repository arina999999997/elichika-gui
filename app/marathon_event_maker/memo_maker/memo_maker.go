package memo_maker
import (
	"elichika/gui/graphic"

	"fmt"
)
const (
	MemoMakerWidth = 1920
	MemoMakerHeight = 1080
)


type MemoMaker struct {
	Parent graphic.Object
	Canvas *graphic.Canvas

	BackgroundPicker *BackgroundPicker
	PinPicker *PinPicker
	TextPicker *TextPicker
	OutputPicker *OutputPicker

	Memo *Memo

	Objects []graphic.Object
	ObjectX map[graphic.Object]int
	ObjectY map[graphic.Object]int
}

func (*MemoMaker) GetWidth() int {
	return MemoMakerWidth
}
func (*MemoMaker) GetHeight() int {
	return MemoMakerHeight
}

func (mm *MemoMaker) InvalidateRenderCache() bool {
	result := mm.Canvas.InvalidateRenderCache()
	graphic.InvalidateRenderCache(mm.Memo)
	if mm.Memo != nil {
		graphic.InvalidateRenderCache(mm.Memo)
		graphic.InvalidateRenderCache(mm.OutputPicker)
	}
	return result
}

func (mm *MemoMaker) ToTexture() *graphic.Texture {
	mm.Draw()
	return mm.Canvas.AsTexture()
}

func (mm*MemoMaker) Draw() {
	if mm.Canvas.IsRendered() {
		return
	}
	if mm.Canvas == nil {
		mm.Canvas = graphic.NewCanvas(mm)
	}

	mm.Canvas.DrawObject(mm.BackgroundPicker, 0, mm.ObjectY[mm.BackgroundPicker], mm.BackgroundPicker.GetWidth(), mm.BackgroundPicker.GetHeight())
	mm.Canvas.DrawObject(mm.PinPicker, 0, mm.ObjectY[mm.PinPicker], mm.PinPicker.GetWidth(), mm.PinPicker.GetHeight())
	mm.Canvas.DrawObject(mm.TextPicker, 0, mm.ObjectY[mm.TextPicker], mm.TextPicker.GetWidth(), mm.TextPicker.GetHeight())
	err := mm.Load()
	if err != nil {
		fmt.Println("Error generating output: ", err)
	}
	mm.Canvas.DrawObject(mm.OutputPicker, 0, mm.ObjectY[mm.OutputPicker], mm.OutputPicker.GetWidth(), mm.OutputPicker.GetHeight())

	mm.Canvas.Finalize()
}

func (mm *MemoMaker) ForEach(f func (graphic.Object)) {
	for _, object := range mm.Objects {
		f(object)
	}
}

func (mm* MemoMaker) MapEvent(e graphic.Event, child graphic.Object) graphic.Event {
	switch e.(type) {
	case graphic.MouseButtonDownEvent:
		mouseEvent := graphic.MouseButtonDownEvent{}
		mouseEvent = e.(graphic.MouseButtonDownEvent)
		mouseEvent.Y -= mm.ObjectY[child]
		return mouseEvent
	default: 
	return e
	}
}

func (mm *MemoMaker) Load() error {
	var err error
	err = mm.BackgroundPicker.Load()
	if err != nil {
		return err
	}
	err = mm.PinPicker.Load()
	if err != nil {
		return err
	}
	err = mm.TextPicker.Load()
	if err != nil {
		return err
	}
	if mm.Memo == nil {
		mm.Memo = &Memo{}
	}
	graphic.InvalidateRenderCache(mm.Memo)
	mm.Memo.Load(mm)
	err = mm.OutputPicker.Load(mm)
	if err != nil {
		return err
	}
	return nil
}

func NewMemoMaker(parent graphic.Object) *MemoMaker {
	memoMaker := &MemoMaker{
		Parent: parent,
	}

	memoMaker.BackgroundPicker = NewBackgroundPicker(memoMaker)
	memoMaker.PinPicker = NewPinPicker(memoMaker)
	memoMaker.TextPicker = NewTextPicker(memoMaker)
	memoMaker.OutputPicker = NewOutputPicker(memoMaker)
	
	y := 0
	memoMaker.ObjectY = map[graphic.Object]int{}
	addObject := func(object graphic.Object) {
		memoMaker.Objects = append(memoMaker.Objects, object)
		memoMaker.ObjectY[object] = y
		y += object.GetHeight()
	}

	addObject(memoMaker.BackgroundPicker)
	addObject(memoMaker.PinPicker)
	addObject(memoMaker.TextPicker)
	addObject(memoMaker.OutputPicker)
	
	return memoMaker
}