package textbox

import (
	"elichika/gui/graphic"

	"fmt"
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
	OnKeyFunc            map[graphic.KeyEvent]func()
}

func (rt *RectTextbox) SetText(s string) {
	graphic.InvalidateRenderCache(rt)
	rt.TextContent = s
}

func NewLabelAndRectTextbox(parent graphic.Object, width, height int, label string) (*graphic.Text, *RectTextbox) {
	text := graphic.NewText(parent, label)
	text.SetHeight(height)
	textbox := &RectTextbox{
		Parent: parent,
		Width: width - text.GetWidth(),
		Height: height,
		Texture: graphic.RGBATexture(0x2f2f2fff),
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
			Parent: parent,
			Width: (width - text.GetWidth() - reservedGap) / textboxCount,
			Height: height,
			Texture: graphic.RGBATexture(0x2f2f2fff),
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
		rt.Canvas.DrawTexture(textTexture, 0, 0, tW,  tH)
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
	fmt.Println("textbox gained focus")
	graphic.InvalidateRenderCache(rt)
}

func (rt *RectTextbox) UnsetFocus() {
	rt.IsFocused = false
	fmt.Println("textbox lost focus")
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
	fmt.Println("new text: ", rt.TextContent)
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
		rt.OnKeyFunc= map[graphic.KeyEvent]func(){}
	}
	rt.OnKeyFunc[key] = f
}

// Input constrants functions

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
	if rt.TextContent != "" {
		value, _ := strconv.Atoi(rt.TextContent)
		if value > high {
			rt.TextContent = strconv.Itoa(high)
		} else if value < low {
			rt.TextContent = strconv.Itoa(low)
		}
	}
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

func (rt *RectTextbox) ForceFloatRange(low, high float64) {
	rt.ForceFloats()
	if rt.TextContent != "" {
		value, _ := strconv.ParseFloat(rt.TextContent, 64)
		if value > high {
			rt.TextContent = strconv.FormatFloat(high, 'f', -1, 64)
		} else if value < low {
			rt.TextContent = strconv.FormatFloat(low, 'f', -1, 64)
		}
	}
}

// Helper functions for special inputs

func (rt *RectTextbox) SetNextInt() {
	graphic.InvalidateRenderCache(rt)
	rt.ForceInts()
	value, _ := strconv.Atoi(rt.TextContent)
	rt.TextContent = strconv.Itoa(value + 1)
}

func (rt *RectTextbox) SetPrevInt() {
	graphic.InvalidateRenderCache(rt)
	rt.ForceInts()
	value, _ := strconv.Atoi(rt.TextContent)
	rt.TextContent = strconv.Itoa(value - 1)
}

func (rt *RectTextbox) SetNextFloat(delta float64, significant int) {
	graphic.InvalidateRenderCache(rt)
	rt.ForceFloats()
	f64, _ := strconv.ParseFloat(rt.TextContent, 64)
	rt.TextContent = strconv.FormatFloat(f64 + delta, 'f', significant, 64)
}


// TODO(gui): This need to check for whether a modulo operation is necessary, otherwise it will overwrite the current data and we can't delete
// func (rt *RectTextbox) FloatMod(mod float64, significant int) {
// 	graphic.InvalidateRenderCache(rt)
// 	rt.ForceFloats()
// 	f64, _ := strconv.ParseFloat(rt.TextContent, 64)
// 	f64 = math.Mod(f64, mod)
// 	rt.TextContent = strconv.FormatFloat(f64, 'f', significant, 64)
// }