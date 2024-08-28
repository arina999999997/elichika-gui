package memo_maker

import (
	"elichika/gui/graphic"
	"elichika/gui/graphic/textbox"

	"fmt"
	"errors"
	"strconv"
)
const MaxTextLine = 6

type TextPicker struct {
	// this is the background of the memo, so no attaching pin included
	Parent graphic.Object
	Canvas *graphic.Canvas

	TextLabel *graphic.Text
	TextTextbox *textbox.RectTextbox

	XLabel *graphic.Text
	XTextbox *textbox.RectTextbox
	YLabel *graphic.Text
	YTextbox *textbox.RectTextbox
	RotationLabel *graphic.Text
	RotationTextbox *textbox.RectTextbox
	FontLabel *graphic.Text
	FontTextbox *textbox.RectTextbox

	CharacterSizeLabel *graphic.Text
	CharacterSizeTextbox *textbox.RectTextbox
	LetterSpacingLabel *graphic.Text
	LetterSpacingTextbox *textbox.RectTextbox
	LineSpacingLabel *graphic.Text
	LineSpacingTextbox *textbox.RectTextbox
	ColorLabel *graphic.Text
	ColorTextbox *textbox.RectTextbox

	DownScalingLabel *graphic.Text
	DownScalingTextbox *textbox.RectTextbox
	LineOffsetLabel *graphic.Text
	LineOffsetTextboxes []*textbox.RectTextbox

	Objects []graphic.Object
	ObjectX map[graphic.Object]int
	ObjectY map[graphic.Object]int

	// line breaks with \n
	Text string
	X int
	Y int
	Rotation float32
	CharacterSize int
	DownscalingFactor int
	LetterSpacing float32
	LineSpacing float32
	Color uint
	LineOffsets []int
	Font *graphic.Font
}

func (tp *TextPicker) Load() error {
	var err error
	var f64 float64
	var ui64 uint64
	tp.Text = tp.TextTextbox.TextContent

	
	tp.X, err = strconv.Atoi(tp.XTextbox.TextContent)
	if err != nil {
		return errors.New(fmt.Sprint("Error parsing text X: ", err))
	}
	tp.Y, err = strconv.Atoi(tp.YTextbox.TextContent)
	if err != nil {
		return errors.New(fmt.Sprint("Error parsing text Y: ", err))
	}
	
	f64, err = strconv.ParseFloat(tp.RotationTextbox.TextContent, 32)
	if err != nil {
		return errors.New(fmt.Sprint("Error parsing text Rotation: ", err))
	} else {
		tp.Rotation = float32(f64)
	}
	
	tp.CharacterSize, err = strconv.Atoi(tp.CharacterSizeTextbox.TextContent)
	if err != nil {
		return errors.New(fmt.Sprint("Error parsing text CharacterSize: ", err))
	}
	
	tp.DownscalingFactor, err = strconv.Atoi(tp.DownScalingTextbox.TextContent)
	if err != nil {
		return errors.New(fmt.Sprint("Error parsing text DownscalingFactor: ", err))
	}

	f64, err = strconv.ParseFloat(tp.LetterSpacingTextbox.TextContent, 32)
	if err != nil {
		return errors.New(fmt.Sprint("Error parsing text LetterSpacing: ", err))
	} else {
		tp.LetterSpacing = float32(f64)
	}
	

	f64, err = strconv.ParseFloat(tp.LineSpacingTextbox.TextContent, 32)
	if err != nil {
		return errors.New(fmt.Sprint("Error parsing text LineSpacing: ", err))
	} else {
		tp.LineSpacing = float32(f64)
	}
	
	ui64, err = strconv.ParseUint(tp.ColorTextbox.TextContent, 0, 32)
	if err != nil {
		return errors.New(fmt.Sprint("Error parsing text Color: ", err))
	}
	fmt.Println("New color: ", tp.Color)
	tp.Color = uint(ui64)

	for i := 0; i < MaxTextLine; i++ {
		tp.LineOffsets[i], err = strconv.Atoi(tp.LineOffsetTextboxes[i].TextContent)
		if err != nil {
			return errors.New(fmt.Sprint("Error parsing line offset: ", err))
		}
	}

	font := graphic.GetFont(tp.FontTextbox.TextContent)
	if font == nil {
		return errors.New(fmt.Sprint("Error loading font: ", tp.FontTextbox.TextContent))
	} else {
		tp.Font = font
	}
	return nil
}

// graphic.Object
func (*TextPicker) GetWidth() int {
	return MemoMakerWidth
}
func (*TextPicker) GetHeight() int {
	return 200
}

func (tp *TextPicker) InvalidateRenderCache() bool {
	return tp.Canvas.InvalidateRenderCache()
}

func (tp *TextPicker) ToTexture() *graphic.Texture {
	tp.Draw()
	return tp.Canvas.AsTexture()
}

func (tp *TextPicker) Draw() {
	if tp.Canvas.IsRendered() {
		return
	}
	if tp.Canvas == nil {
		tp.Canvas = graphic.NewCanvas(tp)
	}

	for _, object := range tp.Objects {
		tp.Canvas.DrawObject(object, tp.ObjectX[object], tp.ObjectY[object], object.GetWidth(), object.GetHeight())
	}
	tp.Canvas.Finalize()
}

// graphic.CompositeObject
func (tp *TextPicker) ForEach(f func(graphic.Object)) {
	for _, object := range tp.Objects {
		f(object)
	}
}

// graphic.PollingCompositeObject
func (tp *TextPicker) MapEvent(e graphic.Event, child graphic.Object) graphic.Event {
	switch e.(type) {
	case graphic.MouseButtonDownEvent:
		mouseEvent := graphic.MouseButtonDownEvent{}
		mouseEvent = e.(graphic.MouseButtonDownEvent)
		mouseEvent.X -= tp.ObjectX[child]
		mouseEvent.Y -= tp.ObjectY[child]
		return mouseEvent
	default:
		return e
	}
}

// graphic.ChildObject
func (tp *TextPicker) GetParent() graphic.Object {
	return tp.Parent
}

func NewTextPicker(parent graphic.Object) *TextPicker{
	picker := &TextPicker{
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

	picker.TextLabel, picker.TextTextbox = 
		textbox.NewLabelAndRectTextbox(picker, MemoMakerWidth, 50, "Memo text: ")
	// picker.TextTextbox.SetText("Text goes here, use \\ for line breaks")
	picker.TextTextbox.SetText("Our surprise\\party was a\\success!")
	
	addObject(picker.TextLabel, false, 0)
	addObject(picker.TextTextbox, false, 0)

	picker.XLabel, picker.XTextbox = 
		textbox.NewLabelAndRectTextbox(picker, MemoMakerWidth / 5 , 50, "Text X: ")
	picker.XTextbox.SetText("35")
	
	picker.YLabel, picker.YTextbox = 
		textbox.NewLabelAndRectTextbox(picker, MemoMakerWidth / 5 , 50, "Text Y: ")
	picker.YTextbox.SetText("77")

	picker.RotationLabel, picker.RotationTextbox = 
		textbox.NewLabelAndRectTextbox(picker, MemoMakerWidth / 5 , 50, "Rotation: ")
	picker.RotationTextbox.SetText("-4")

	picker.FontLabel, picker.FontTextbox = 
		textbox.NewLabelAndRectTextbox(picker, MemoMakerWidth - (MemoMakerWidth / 5 * 3) - 15 , 50, "Font: ")
	picker.FontTextbox.SetText("gui/fonts/FOT-SkipStd-M.otf")
	
	addObject(picker.XLabel, true, 0)
	addObject(picker.XTextbox, false, 5)
	addObject(picker.YLabel, false, 0)
	addObject(picker.YTextbox, false, 5)
	addObject(picker.RotationLabel, false, 0)
	addObject(picker.RotationTextbox, false, 5)
	addObject(picker.FontLabel, false, 0)
	addObject(picker.FontTextbox, false, 5)
	

	picker.CharacterSizeLabel, picker.CharacterSizeTextbox = 
		textbox.NewLabelAndRectTextbox(picker, MemoMakerWidth / 4, 50, "Character Size: ")
	picker.CharacterSizeTextbox.SetText("18")
	picker.LetterSpacingLabel, picker.LetterSpacingTextbox = 
		textbox.NewLabelAndRectTextbox(picker, MemoMakerWidth / 4, 50, "Letter Spacing: ")
	picker.LetterSpacingTextbox.SetText("1.32")
	picker.LineSpacingLabel, picker.LineSpacingTextbox = 
		textbox.NewLabelAndRectTextbox(picker, MemoMakerWidth - (MemoMakerWidth / 4 * 3) - 15, 50, "Line Spacing: ")
	picker.LineSpacingTextbox.SetText("0.72")
	picker.ColorLabel, picker.ColorTextbox = 
		textbox.NewLabelAndRectTextbox(picker, MemoMakerWidth / 4, 50, "Color: ")
	picker.ColorTextbox.SetText("0x494949ff")

	addObject(picker.CharacterSizeLabel, true, 0)
	addObject(picker.CharacterSizeTextbox, false, 5)
	addObject(picker.LetterSpacingLabel, false, 0)
	addObject(picker.LetterSpacingTextbox, false, 5)
	addObject(picker.LineSpacingLabel, false, 0)
	addObject(picker.LineSpacingTextbox, false, 5)
	addObject(picker.ColorLabel, false, 0)
	addObject(picker.ColorTextbox, false, 5)
	

	
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

	picker.CharacterSizeTextbox.SetOnKeyFunc(graphic.KeyEventUp, func(){
		picker.CharacterSizeTextbox.SetNextInt()
	})
	picker.CharacterSizeTextbox.SetOnKeyFunc(graphic.KeyEventDown, func(){
		picker.CharacterSizeTextbox.SetPrevInt()
	})

	picker.RotationTextbox.SetOnKeyFunc(graphic.KeyEventUp, func(){
		picker.RotationTextbox.SetNextFloat(0.1, 1)
		picker.RotationTextbox.OnTextUpdateFunc()
	})
	picker.RotationTextbox.SetOnKeyFunc(graphic.KeyEventDown, func(){
		picker.RotationTextbox.SetNextFloat(-0.1, 1)
		picker.RotationTextbox.OnTextUpdateFunc()
	})

	picker.LetterSpacingTextbox.SetOnKeyFunc(graphic.KeyEventUp, func(){
		picker.LetterSpacingTextbox.SetNextFloat(0.01, 2)
	})
	picker.LetterSpacingTextbox.SetOnKeyFunc(graphic.KeyEventDown, func(){
		picker.LetterSpacingTextbox.SetNextFloat(-0.01, 2)
	})

	picker.LineSpacingTextbox.SetOnKeyFunc(graphic.KeyEventUp, func(){
		picker.LineSpacingTextbox.SetNextFloat(0.01, 2)
	})
	picker.LineSpacingTextbox.SetOnKeyFunc(graphic.KeyEventDown, func(){
		picker.LineSpacingTextbox.SetNextFloat(-0.01, 2)
	})

	
	picker.XTextbox.OnTextUpdateFunc = func() {
		picker.XTextbox.ForceIntRange(-200, 200)
	}
	picker.YTextbox.OnTextUpdateFunc = func() {
		picker.YTextbox.ForceIntRange(-200, 200)
	}
	picker.RotationTextbox.OnTextUpdateFunc = func(){
		picker.RotationTextbox.ForceFloatRange(-180, 180)
	}
	picker.CharacterSizeTextbox.OnTextUpdateFunc = func() {
		picker.CharacterSizeTextbox.ForceIntRange(1, 64)
	}
	picker.LetterSpacingTextbox.OnTextUpdateFunc = func() {
		picker.LetterSpacingTextbox.ForceFloatRange(-3, 3)
	}
	picker.LineSpacingTextbox.OnTextUpdateFunc = func() {
		picker.LineSpacingTextbox.ForceFloatRange(0, 3)
	}

	picker.DownScalingLabel, picker.DownScalingTextbox = textbox.NewLabelAndRectTextbox(
		picker, MemoMakerWidth / 2 - 5, 50, "Text down scaling: ")
	picker.DownScalingTextbox.SetText("4")
	picker.DownScalingTextbox.OnTextUpdateFunc = func() {
		picker.DownScalingTextbox.ForceIntRange(1, 10)
	}
	picker.DownScalingTextbox.SetOnKeyFunc(graphic.KeyEventUp, func(){
		picker.DownScalingTextbox.SetNextInt()
		picker.DownScalingTextbox.OnTextUpdateFunc()
	})
	picker.DownScalingTextbox.SetOnKeyFunc(graphic.KeyEventDown, func(){
		picker.DownScalingTextbox.SetPrevInt()
		picker.DownScalingTextbox.OnTextUpdateFunc()
	})


	picker.LineOffsetLabel, picker.LineOffsetTextboxes = textbox.NewLabelAndRectTextboxes(
		picker,	(MemoMakerWidth + 1) / 2,	50, "Line offsets: ", MaxTextLine, (MaxTextLine - 1) * 5)
	addObject(picker.DownScalingLabel, true, 0)
	addObject(picker.DownScalingTextbox, false, 0)
	addObject(picker.LineOffsetLabel, false, 0)
	for i := 0; i < MaxTextLine; i++ {
		object := picker.LineOffsetTextboxes[i]
		addObject(object, false, 5)
		object.SetText("0")
		picker.LineOffsets = append(picker.LineOffsets, 0)
		object.OnTextUpdateFunc = func() {
			object.ForceIntRange(-10, 10)
		}
		object.SetOnKeyFunc(graphic.KeyEventUp, func(){
			object.SetNextInt()
		})
		object.SetOnKeyFunc(graphic.KeyEventDown, func(){
			object.SetPrevInt()
		})
		
	}

	return picker
}
