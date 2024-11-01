package textbox

import (
	"elichika/gui/graphic"

	// "fmt"
	"strconv"
	// "math"
)

type RectTextbox struct {
	Parent graphic.Object

	Width  int
	Height int

	IsFocused bool

	Texture      *graphic.Texture
	FocusTexture *graphic.Texture

	Text          *graphic.Text
	TextContent   string
	TextStyleType int

	Canvas *graphic.Canvas

	OnTextUpdateFunc func()
	OnEnterFunc      func()
	OnKeyFunc        map[graphic.KeyEvent]func()

	// Optional function to sync the value displayed in the textbox from some value in memory
	SyncFunc func()
}

func (rt *RectTextbox) SetText(s string) {
	graphic.InvalidateRenderCache(rt)
	rt.TextContent = s
}

func NewLabelAndRectTextbox(parent graphic.Object, width, height int, label string) (*graphic.Text, *RectTextbox) {
	text := graphic.NewText(parent, label)
	text.SetHeight(height)
	textbox := &RectTextbox{
		Parent:       parent,
		Width:        width - text.GetWidth(),
		Height:       height,
		Texture:      graphic.RGBATexture(0x2f2f2fff),
		FocusTexture: graphic.RGBATexture(0x7f7f7fff),
	}
	if textbox.Width <= 0 {
		panic("label is too long for desired size")
	}
	return text, textbox
}

func NewLabelAndRectTextboxes(parent graphic.Object, width, height int, label string, textboxCount, reservedGap int) (*graphic.Text, []*RectTextbox) {
	text := graphic.NewText(parent, label)
	text.SetHeight(height)
	textboxes := []*RectTextbox{}
	for i := 0; i < textboxCount; i++ {
		textbox := &RectTextbox{
			Parent:       parent,
			Width:        (width - text.GetWidth() - reservedGap) / textboxCount,
			Height:       height,
			Texture:      graphic.RGBATexture(0x2f2f2fff),
			FocusTexture: graphic.RGBATexture(0x7f7f7fff),
		}
		if textbox.Width <= 0 {
			panic("label and reserved is too long for desired size")
		}
		textboxes = append(textboxes, textbox)
	}
	return text, textboxes
}

// Object
func (rt *RectTextbox) GetWidth() int {
	return rt.Width
}

func (rt *RectTextbox) GetHeight() int {
	return rt.Height
}

func (rt *RectTextbox) InvalidateRenderCache() bool {
	return rt.Canvas.InvalidateRenderCache()
}

func (rt *RectTextbox) Draw() {
	if rt.Canvas.IsRendered() {
		return
	}
	if rt.Canvas == nil {
		rt.Canvas = graphic.NewCanvas(rt)
	}

	texture := rt.Texture
	if rt.IsFocused && (rt.FocusTexture != nil) {
		texture = rt.FocusTexture
	}
	if texture != nil {
		rt.Canvas.DrawTexture(texture, 0, 0, rt.Width, rt.Height)
	}
	if rt.TextContent != "" {
		if rt.Text == nil {
			rt.Text = graphic.NewText(rt, rt.TextContent)
		} else {
			rt.Text.SetText(rt.TextContent)
		}
		rt.Text.SetHeight(rt.Height)
		textTexture := rt.Text.ToTexture()
		textTexture.StyleType = graphic.StyleTypeNone
		tW, tH := textTexture.GetSize()
		rt.Canvas.DrawTexture(textTexture, 0, 0, tW, tH)
	}
	rt.Canvas.Finalize()
}

func (rt *RectTextbox) ToTexture() *graphic.Texture {
	rt.Draw()
	texture := rt.Canvas.AsTexture()
	return texture
}

// Child Object
func (rt *RectTextbox) GetParent() graphic.Object {
	return rt.Parent
}

// Focusable
func (rt *RectTextbox) SetFocus() {
	rt.IsFocused = true
	// fmt.Println("textbox gained focus")
	graphic.InvalidateRenderCache(rt)
}

func (rt *RectTextbox) UnsetFocus() {
	rt.IsFocused = false
	// fmt.Println("textbox lost focus")
	graphic.InvalidateRenderCache(rt)
}

func (rt *RectTextbox) HasFocus() bool {
	return rt.IsFocused
}

// Clickable
func (rt *RectTextbox) OnClick(w *graphic.Window, event graphic.MouseButtonDownEvent) bool {
	if (event.X < 0) || (event.X >= rt.Width) {
		return false
	}
	if (event.Y < 0) || (event.Y >= rt.Height) {
		return false
	}
	w.SetFocusObject(rt)
	return true
}

// Inputable
func (rt *RectTextbox) OnText(event graphic.TextEvent) bool {
	if !rt.IsFocused {
		return false
	}
	if !UpdateInputText(&rt.TextContent, event.Rune) {
		return false
	}
	// fmt.Println("new text: ", rt.TextContent)
	graphic.InvalidateRenderCache(rt)
	if rt.OnTextUpdateFunc != nil {
		rt.OnTextUpdateFunc()
	}
	return true
}

func (rt *RectTextbox) OnPaste(event graphic.PasteEvent) bool {
	if !rt.IsFocused {
		return false
	}
	updated := false
	for _, r := range event.Clipboard {
		if UpdateInputText(&rt.TextContent, r) {
			updated = true
		}
	}
	if !updated {
		return false
	}
	graphic.InvalidateRenderCache(rt)
	if rt.OnTextUpdateFunc != nil {
		rt.OnTextUpdateFunc()
	}
	return true
}

func (rt *RectTextbox) OnEnter() bool {
	if !rt.IsFocused {
		return false
	}
	if rt.OnEnterFunc != nil {
		rt.OnEnterFunc()
	}
	return true
}

func (rt *RectTextbox) OnKey(key graphic.KeyEvent) bool {
	f, exist := rt.OnKeyFunc[key]
	if !exist {
		return false
	}
	f()
	return true
}

func (rt *RectTextbox) SetOnKeyFunc(key graphic.KeyEvent, f func()) {
	if rt.OnKeyFunc == nil {
		rt.OnKeyFunc = map[graphic.KeyEvent]func(){}
	}
	rt.OnKeyFunc[key] = f
}

// typed getters / setters , assuming the values are correct
func (rt *RectTextbox) GetInt() int {
	value, _ := strconv.Atoi(rt.TextContent)
	return value
}
func (rt *RectTextbox) SetInt(value int) {
	graphic.InvalidateRenderCache(rt)
	rt.TextContent = strconv.Itoa(value)
}

func (rt *RectTextbox) GetHex() uint {
	value, _ := strconv.ParseUint(rt.TextContent, 16, 64)
	return uint(value)
}
func (rt *RectTextbox) SetHex(value uint) {
	rt.TextContent = strconv.FormatUint(uint64(value), 16)
}

func (rt *RectTextbox) GetFloat() float32 {
	value, _ := strconv.ParseFloat(rt.TextContent, 32)
	return float32(value)
}
func (rt *RectTextbox) SetFloat(value float32) {
	graphic.InvalidateRenderCache(rt)
	rt.TextContent = strconv.FormatFloat(float64(value), 'f', -1, 32)
}

// value transforming functions
func (rt *RectTextbox) ForceInts() {
	for i, r := range rt.TextContent {
		good := (r >= '0') && (r <= '9')
		good = good || ((r == '-') && (i == 0))
		if !good {
			rt.TextContent = rt.TextContent[:i]
			break
		}
	}
}

func (rt *RectTextbox) ForceIntRange(low, high int) {
	rt.ForceInts()
	value := rt.GetInt()
	if value > high {
		value = high
	} else if value < low {
		value = low
	}
	rt.SetInt(value)
}

func (rt *RectTextbox) ForceFloats() {
	hasDot := false
	for i, r := range rt.TextContent {
		good := (r >= '0') && (r <= '9')
		good = good || ((r == '-') && (i == 0))
		good = good || ((r == '.') && (hasDot == false))
		hasDot = hasDot || (r == '.')
		if !good {
			rt.TextContent = rt.TextContent[:i]
			break
		}
	}
}

func (rt *RectTextbox) ForceFloatRange(low, high float32) {
	rt.ForceFloats()
	if rt.TextContent != "" {
		f64, _ := strconv.ParseFloat(rt.TextContent, 32)
		value := float32(f64)
		if value > high {
			value = high
		} else if value < low {
			value = low
		}
		rt.SetFloat(value)
	}
}

// Helper functions for special inputs
func (rt *RectTextbox) SetNextInt() {
	rt.ForceInts()
	rt.SetInt(rt.GetInt() + 1)
}

func (rt *RectTextbox) SetPrevInt() {
	rt.ForceInts()
	rt.SetInt(rt.GetInt() - 1)
}

func (rt *RectTextbox) SetNextFloat(delta float32, significant int) {
	graphic.InvalidateRenderCache(rt)
	rt.ForceFloats()
	f64 := rt.GetFloat()
	// set float use max significant, so we need to do this
	rt.TextContent = strconv.FormatFloat(float64(f64+delta), 'f', significant, 32)
}

// Helpers for setting up commonly used type of textboxes

func (rt *RectTextbox) SetStringSettingTextbox(defaultValue string, value *string) {
	rt.TextContent = defaultValue
	if value != nil {
		*value = defaultValue
		rt.OnTextUpdateFunc = func() {
			*value = rt.TextContent
		}
		rt.SyncFunc = func() {
			rt.TextContent = *value
		}
	}

}

// A textbox where user can use UP and DOWN key to tune the values
// TODO(gui): Maybe scrolling too
// if low < high, then the values are constrained to that range, otherwise the values are not constrained
// if value is not nil, then value is updated every time the text is updated
func (rt *RectTextbox) SetIntSettingTextbox(defaultValue, lowValue, highValue int, value *int) {
	rt.SetInt(defaultValue)
	if lowValue < highValue {
		if (defaultValue < lowValue) || (defaultValue > highValue) {
			panic("wrongly setup int setting text box")
		}
		rt.OnTextUpdateFunc = func() {
			rt.ForceIntRange(lowValue, highValue)
			if value != nil {
				*value = rt.GetInt()
			}
		}
	} else if value != nil {
		rt.OnTextUpdateFunc = func() {
			*value = rt.GetInt()
		}
	}
	if value != nil {
		*value = defaultValue
		rt.SyncFunc = func() {
			rt.SetInt(*value)
		}
	}
	rt.SetOnKeyFunc(graphic.KeyEventUp, func() {
		rt.SetNextInt()
		if rt.OnTextUpdateFunc != nil {
			rt.OnTextUpdateFunc()
		}
	})
	rt.SetOnKeyFunc(graphic.KeyEventDown, func() {
		rt.SetPrevInt()
		if rt.OnTextUpdateFunc != nil {
			rt.OnTextUpdateFunc()
		}
	})
}

// coordinate textboxs are a pair of textboxes where one control x and one control y
// instead of the standard up and down for increase and decrease, we use up down left right to move the coordinate
// in that direction when either of the textboxes are focused
// if x and y is not nil, then those values are updated along with the text as well
func SetCoordinateTextboxes(xTextbox, yTextbox *RectTextbox, xDefault, yDefault, xLow, yLow, xHigh, yHigh int, x, y *int) {
	xTextbox.SetIntSettingTextbox(xDefault, xLow, xHigh, x)
	yTextbox.SetIntSettingTextbox(yDefault, yLow, yHigh, y)
	leftKeyFunc := func() {
		xTextbox.SetPrevInt()
		if xTextbox.OnTextUpdateFunc != nil {
			xTextbox.OnTextUpdateFunc()
		}
	}
	rightKeyFunc := func() {
		xTextbox.SetNextInt()
		if xTextbox.OnTextUpdateFunc != nil {
			xTextbox.OnTextUpdateFunc()
		}
	}
	upKeyFunc := func() {
		yTextbox.SetPrevInt()
		if yTextbox.OnTextUpdateFunc != nil {
			yTextbox.OnTextUpdateFunc()
		}
	}
	downKeyFunc := func() {
		yTextbox.SetNextInt()
		if yTextbox.OnTextUpdateFunc != nil {
			yTextbox.OnTextUpdateFunc()
		}
	}

	xTextbox.SetOnKeyFunc(graphic.KeyEventLeft, leftKeyFunc)
	xTextbox.SetOnKeyFunc(graphic.KeyEventRight, rightKeyFunc)
	xTextbox.SetOnKeyFunc(graphic.KeyEventUp, upKeyFunc)
	xTextbox.SetOnKeyFunc(graphic.KeyEventDown, downKeyFunc)
	yTextbox.SetOnKeyFunc(graphic.KeyEventLeft, leftKeyFunc)
	yTextbox.SetOnKeyFunc(graphic.KeyEventRight, rightKeyFunc)
	yTextbox.SetOnKeyFunc(graphic.KeyEventUp, upKeyFunc)
	yTextbox.SetOnKeyFunc(graphic.KeyEventDown, downKeyFunc)
}

// A textbox where user can use to input hex values, mostly used for RGBA colors
func (rt *RectTextbox) SetHexSettingTextbox(defaultValue uint, value *uint) {
	rt.SetHex(defaultValue)
	if value != nil {
		*value = defaultValue
		rt.SyncFunc = func() {
			rt.SetHex(*value)
		}
		rt.OnTextUpdateFunc = func() {
			*value = rt.GetHex()
		}
	}
}

func (rt *RectTextbox) SetFloatSettingTextbox(defaultValue, lowValue, highValue, step float32, significantDigit int, value *float32) {
	rt.SetFloat(defaultValue)
	if lowValue < highValue {
		if (defaultValue < lowValue) || (defaultValue > highValue) {
			panic("wrongly setup float setting text box")
		}
		rt.OnTextUpdateFunc = func() {
			rt.ForceFloatRange(lowValue, highValue)
			if value != nil {
				*value = rt.GetFloat()
			}
		}
	} else if value != nil {
		rt.OnTextUpdateFunc = func() {
			*value = rt.GetFloat()
		}
	}
	if value != nil {
		*value = defaultValue
		rt.SyncFunc = func() {
			rt.SetFloat(*value)
		}
	}
	rt.SetOnKeyFunc(graphic.KeyEventUp, func() {
		rt.SetNextFloat(step, significantDigit)
		if rt.OnTextUpdateFunc != nil {
			rt.OnTextUpdateFunc()
		}
	})
	rt.SetOnKeyFunc(graphic.KeyEventDown, func() {
		rt.SetNextFloat(-step, significantDigit)
		if rt.OnTextUpdateFunc != nil {
			rt.OnTextUpdateFunc()
		}
	})
}

func (rt *RectTextbox) TrySync() {
	if rt.SyncFunc != nil {
		rt.SyncFunc()
	}
}
