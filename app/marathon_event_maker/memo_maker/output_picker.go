package memo_maker
import (
	"elichika/gui/graphic"
	"elichika/gui/graphic/textbox"
	"elichika/gui/graphic/button"

	"fmt"
	"errors"
	"strconv"

	
	"github.com/telroshan/go-sfml/v2/graphics"
)
type OutputPicker struct {
	Parent graphic.Object
	Canvas *graphic.Canvas

	Memo *Memo
	OutputX int
	OutputY int 


	VisualScalingLabel *graphic.Text
	VisualScalingTextbox *textbox.RectTextbox
	
	ExportToLabel *graphic.Text
	ExportToTextbox *textbox.RectTextbox
	ExportButton *button.RectButton



	Objects []graphic.Object
	ObjectX map[graphic.Object]int
	ObjectY map[graphic.Object]int

	ExportToPath string
	VisualScalingFactor float32

	
}

func (op* OutputPicker) Load(parent *MemoMaker) error {
	var err error
	var f64 float64
	op.Memo = parent.Memo
	op.ExportToPath = op.ExportToTextbox.TextContent

	f64, err = strconv.ParseFloat(op.VisualScalingTextbox.TextContent, 32)
	if err != nil {
		return errors.New(fmt.Sprint("Error parsing visual scaling: ", err))
	}
	op.VisualScalingFactor = float32(f64)
	return nil
}


// graphic.Object
func (*OutputPicker) GetWidth() int {
	return MemoMakerWidth
}
func (*OutputPicker) GetHeight() int {
	return MemoMakerHeight
}

func (op *OutputPicker) InvalidateRenderCache() bool {
	return op.Canvas.InvalidateRenderCache()
}


func (op *OutputPicker) ToTexture() *graphic.Texture {
	op.Draw()
	return op.Canvas.AsTexture()
}

func (op *OutputPicker) Draw() {
	if op.Canvas.IsRendered() {
		return
	}
	if op.Canvas == nil {
		op.Canvas = graphic.NewCanvas(op)
	}

	for _, object := range op.Objects {
		op.Canvas.DrawObject(object, op.ObjectX[object], op.ObjectY[object], object.GetWidth(), op.GetHeight())
	}
	// draw the output, if possible
	if op.Memo != nil {
		sprite := graphics.SfSprite_create()
		defer graphics.SfSprite_destroy(sprite)
		graphics.SfSprite_setTexture(sprite, op.Memo.ToTexture().Texture, 1)
		graphics.SfSprite_setPosition(sprite, graphic.GetVector2f(op.OutputX, op.OutputY))
		graphics.SfSprite_setScale(sprite, graphic.GetVector2ff(op.VisualScalingFactor, op.VisualScalingFactor))
		op.Canvas.DrawSprite(&graphic.Sprite{
			Sprite: sprite,
		})
	}
	op.Canvas.Finalize()
}

// graphic.CompositeObject
func (op *OutputPicker) ForEach(f func(graphic.Object)) {
	for _, object := range op.Objects {
		f(object)
	}
}

// graphic.PollingCompositeObject
func (op *OutputPicker) MapEvent(e graphic.Event, child graphic.Object) graphic.Event {
	switch e.(type) {
	case graphic.MouseButtonDownEvent:
		mouseEvent := graphic.MouseButtonDownEvent{}
		mouseEvent = e.(graphic.MouseButtonDownEvent)
		mouseEvent.X -= op.ObjectX[child]
		mouseEvent.Y -= op.ObjectY[child]
		return mouseEvent
	default:
		return e
	}
}

// graphic.ChildObject
func (op *OutputPicker) GetParent() graphic.Object {
	return op.Parent
}

func NewOutputPicker(parent *MemoMaker) *OutputPicker{
	picker := &OutputPicker{
		Parent: parent,
	}
	
	picker.ObjectX = map[graphic.Object]int{}
	picker.ObjectY = map[graphic.Object]int{}
	x := 0
	y := 0
	
	addObject := func(object graphic.Object, newLine bool, paddingX int) {
		if newLine {
			y += 50
			x = 0
		}
		picker.Objects = append(picker.Objects, object)
		picker.ObjectX[object] = x
		picker.ObjectY[object] = y

		x += object.GetWidth() + paddingX
	}

	picker.VisualScalingLabel, picker.VisualScalingTextbox = 
		textbox.NewLabelAndRectTextbox(picker, MemoMakerWidth, 50, "Output scaling factor (preview only): ")
	picker.VisualScalingTextbox.SetText("1.2")

	addObject(picker.VisualScalingLabel, false, 0)
	addObject(picker.VisualScalingTextbox, false, 0)

	picker.ExportToLabel, picker.ExportToTextbox = 
		textbox.NewLabelAndRectTextbox(picker, MemoMakerWidth -200, 50, "Export to: ")
	picker.ExportToTextbox.SetText("output.png")

	picker.ExportButton = button.NewButton(picker, 200, 50, "Export", func(){
		picker.Draw()
		picker.Memo.ToTexture().SaveToImage(picker.ExportToPath)
	}, nil)

	addObject(picker.ExportToLabel, true, 0)
	addObject(picker.ExportToTextbox, false, 0)
	addObject(picker.ExportButton, false, 0)

	picker.OutputX = 0
	picker.OutputY = y + 50

	picker.VisualScalingTextbox.OnTextUpdateFunc = func() {
		picker.VisualScalingTextbox.ForceFloatRange(0.5, 4)
	}

	picker.VisualScalingTextbox.SetOnKeyFunc(graphic.KeyEventUp, func() {
		picker.VisualScalingTextbox.SetNextFloat(0.1, 1)
	})
	picker.VisualScalingTextbox.SetOnKeyFunc(graphic.KeyEventDown, func() {
		picker.VisualScalingTextbox.SetNextFloat(-0.1, 1)
	})
	
	return picker
}
