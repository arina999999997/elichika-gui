package memo_maker
import (
	"elichika/gui/graphic"
	"elichika/gui/graphic/textbox"

	"fmt"
	"errors"
	"strconv"
)
type PinPicker struct {
	Parent graphic.Object
	Canvas *graphic.Canvas

	// default to gui/app/marathon_event_maker/memo_maker/pink_pin.png
	TextureLabel *graphic.Text
	TextureTextbox *textbox.RectTextbox

	XLabel *graphic.Text
	XTextbox *textbox.RectTextbox

	YLabel *graphic.Text
	YTextbox *textbox.RectTextbox

	Objects []graphic.Object
	ObjectX map[graphic.Object]int

	LoadedPath string
	Texture *graphic.Texture
	X int
	Y int
}

func (pp* PinPicker) Load() error {
	var err error
	if pp.TextureTextbox.TextContent != pp.LoadedPath {
		pp.Texture, err = graphic.FileTexture(pp.TextureTextbox.TextContent)
		if err != nil {
			return errors.New(fmt.Sprint("Error while loading pin texture: ", err))
		}
	}

	pp.X, err = strconv.Atoi(pp.XTextbox.TextContent)
	if err != nil {
		return errors.New(fmt.Sprint("Error parsing pin X: ", err))
	}
	pp.Y, err = strconv.Atoi(pp.YTextbox.TextContent)
	if err != nil {
		return errors.New(fmt.Sprint("Error parsing pin Y: ", err))
	}
	return nil
}


// graphic.Object
func (*PinPicker) GetWidth() int {
	return MemoMakerWidth
}
func (*PinPicker) GetHeight() int {
	return 50
}

func (pp *PinPicker) InvalidateRenderCache() bool {
	return pp.Canvas.InvalidateRenderCache()
}


func (pp *PinPicker) ToTexture() *graphic.Texture {
	pp.Draw()
	return pp.Canvas.AsTexture()
}

func (pp *PinPicker) Draw() {
	if pp.Canvas.IsRendered() {
		return
	}
	if pp.Canvas == nil {
		pp.Canvas = graphic.NewCanvas(pp)
	}

	for _, object := range pp.Objects {
		pp.Canvas.DrawObject(object, pp.ObjectX[object], 0, object.GetWidth(), pp.GetHeight())
	}

	pp.Canvas.Finalize()
}


// graphic.CompositeObject
func (pp *PinPicker) ForEach(f func(graphic.Object)) {
	for _, object := range pp.Objects {
		f(object)
	}
}

// graphic.PollingCompositeObject
func (pp *PinPicker) MapEvent(e graphic.Event, child graphic.Object) graphic.Event {
	switch e.(type) {
	case graphic.MouseButtonDownEvent:
		mouseEvent := graphic.MouseButtonDownEvent{}
		mouseEvent = e.(graphic.MouseButtonDownEvent)
		mouseEvent.X -= pp.ObjectX[child]
		return mouseEvent
	default:
		return e
	}
}

// graphic.ChildObject
func (pp *PinPicker) GetParent() graphic.Object {
	return pp.Parent
}


func NewPinPicker(parent graphic.Object) *PinPicker{
	picker := &PinPicker{
		Parent: parent,
	}
	
	picker.TextureLabel, picker.TextureTextbox = 
		textbox.NewLabelAndRectTextbox(picker, MemoMakerWidth * 6 / 8 - 10, 50, "Memo pin: ")
	picker.TextureTextbox.SetText("gui/app/marathon_event_maker/memo_maker/pink_pin.png")

	picker.XLabel, picker.XTextbox = 
		textbox.NewLabelAndRectTextbox(picker, MemoMakerWidth * 1 / 8, 50, "Pin X: ")
	picker.XTextbox.SetText("0")
	picker.YLabel, picker.YTextbox = 
		textbox.NewLabelAndRectTextbox(picker, MemoMakerWidth * 1 / 8, 50, "Pin Y: ")
	picker.YTextbox.SetText("0")

	picker.Objects = []graphic.Object{picker.TextureLabel, picker.TextureTextbox, picker.XLabel, picker.XTextbox, picker.YLabel, picker.YTextbox}
	picker.ObjectX = map[graphic.Object]int{}
	x := 0
	for i, object := range picker.Objects {
		picker.ObjectX[object] = x
		x += object.GetWidth()
		if i %2 == 1 {
			x += 5
		}
	}

	picker.XTextbox.OnTextUpdateFunc = func() {
		picker.XTextbox.ForceIntRange(-200, 200)
	}
	picker.YTextbox.OnTextUpdateFunc = func() {
		picker.YTextbox.ForceIntRange(-200, 200)
	}

	leftKeyFunc := func(){
		picker.XTextbox.SetPrevInt()
		picker.XTextbox.OnTextUpdateFunc()
	}
	rightKeyFunc := func(){
		picker.XTextbox.SetNextInt()
		picker.XTextbox.OnTextUpdateFunc()
	}
	upKeyFunc := func(){
		picker.YTextbox.SetPrevInt()
		picker.YTextbox.OnTextUpdateFunc()
	}
	downKeyFunc := func(){
		picker.YTextbox.SetNextInt()
		picker.YTextbox.OnTextUpdateFunc()
	}

	picker.XTextbox.SetOnKeyFunc(graphic.KeyEventLeft, leftKeyFunc)
	picker.XTextbox.SetOnKeyFunc(graphic.KeyEventRight, rightKeyFunc)
	picker.XTextbox.SetOnKeyFunc(graphic.KeyEventUp, upKeyFunc)
	picker.XTextbox.SetOnKeyFunc(graphic.KeyEventDown, downKeyFunc)
	picker.YTextbox.SetOnKeyFunc(graphic.KeyEventLeft, leftKeyFunc)
	picker.YTextbox.SetOnKeyFunc(graphic.KeyEventRight, rightKeyFunc)
	picker.YTextbox.SetOnKeyFunc(graphic.KeyEventUp, upKeyFunc)
	picker.YTextbox.SetOnKeyFunc(graphic.KeyEventDown, downKeyFunc)
	
	return picker
}
