package graphic

import (
	"unicode/utf8"
	"unsafe"

	"github.com/telroshan/go-sfml/v2/graphics"
)

type sfText = graphics.Struct_SS_sfText
type sfFloatRect = graphics.SfFloatRect

// TODO(text): For now, Text is single line only, ideally we would want it to fit into a box, and use multiple lines if necessary

type Text struct {
	height      int
	textContent string

	text sfText
	font *Font

	bounds *sfFloatRect

	Parent Object

	canvas *Canvas

	LineSpacingFactor float32
}

func (t *Text) LoadUTF32(s string) {
	if t.textContent == s {
		return
	}
	InvalidateRenderCache(t)
	t.textContent = s
	// very cursed but this is the only way to get around cgo shitty memory leak detection
	runeCount := utf8.RuneCountInString(s)
	dummy := string("a")
	for i := 0; i < runeCount; i++ {
		dummy += "a"
	}
	graphics.SfText_setString(t.text, dummy)
	mem := graphics.SfText_getUnicodeString(t.text)
	ptr := uintptr(unsafe.Pointer(mem))
	for _, r := range s {
		*(*uint)(unsafe.Pointer(ptr)) = uint(r)
		ptr += 4
	}
	*(*uint)(unsafe.Pointer(ptr)) = uint(0)
	graphics.SfText_setUnicodeString(t.text, mem)
}

// own functions
func NewText(parent Object, s string) *Text {
	text := Text{
		Parent: parent,
	}
	text.text = graphics.SfText_create()
	text.LoadUTF32(s)
	text.SetFont(GetDefaultFont())
	return &text
}

func (t *Text) SetText(s string) {
	t.LoadUTF32(s)
}

func (t *Text) SetFont(f *Font) {
	if t.font == f {
		return
	}
	InvalidateRenderCache(t)
	t.font = f
	graphics.SfText_setFont(t.text, t.font.font)
}

func (t *Text) SetCharacterSize(characterSize int) {
	InvalidateRenderCache(t)
	graphics.SfText_setCharacterSize(t.text, uint(characterSize))
}

func (t *Text) SetLetterSpacing(LetterSpacing float32) {
	InvalidateRenderCache(t)
	graphics.SfText_setLetterSpacing(t.text, LetterSpacing)
}

func (t *Text) SetLineSpacingFactor(lineSpacingFactor float32) {
	InvalidateRenderCache(t)
	t.LineSpacingFactor = lineSpacingFactor
}

func (t *Text) GetLineSpacing() int {
	return int(t.LineSpacingFactor * float32(graphics.SfFont_getLineSpacing(t.font.font, graphics.SfText_getCharacterSize(t.text))))
}

// rgba from high bytes to low bytes
func (t *Text) SetColor(color uint) {
	InvalidateRenderCache(t)
	graphics.SfText_setColor(t.text, graphics.SfColor_fromInteger(color))
}

// set the height of the text so it fit in a box of this height
// the actual height of the text is still decided by other factors.
var hasRune = map[rune]bool{}
var adaptiveString = "|"

func (t *Text) SetHeight(height int) {
	if height >= 256 {
		panic("text height is too big")
	}
	// if t.height == height {
	// 	return
	// }
	// t.height = height

	InvalidateRenderCache(t)

	oldText := t.textContent

	// this is kind of a hack, we try to render a string that is all character we have seen before
	// there's probably a better way of doing this
	for _, r := range oldText {
		if hasRune[r] {
			continue
		}
		hasRune[r] = true
		adaptiveString += string(r)
	}

	t.SetText(adaptiveString)

	// the formula for going from character size to desired height change depending on the font, so we use a binary search here
	size := uint(1)
	for size <= 256 { // fail save
		InvalidateRenderCache(t)
		graphics.SfText_setCharacterSize(t.text, size)
		if t.GetHeight() > height {
			break
		}
		size *= 2
	}
	size /= 2
	finalSize := uint(0)
	for size > 0 {
		InvalidateRenderCache(t)
		graphics.SfText_setCharacterSize(t.text, finalSize+size)
		if t.GetHeight() <= height {
			finalSize += size
		}
		size /= 2
	}
	InvalidateRenderCache(t)
	graphics.SfText_setCharacterSize(t.text, finalSize)
	t.SetText(oldText)
}

func (t *Text) getBounds() {
	if t.bounds != nil {
		return
	}
	rect := graphics.SfText_getLocalBounds(t.text)
	t.bounds = &rect
}

// Object interface

// these are hidden offset in sfml:
// - text is aligned at the base line of the text
// - but when drawing, the draw using the "top left" of the text string
// - but we don't want the text to move up and down if we use higher letter
// - so there is an offset.
// - GetTop() + GetHeight() will always result in the highest possible character
// - which is what we want to draw the text properly.
// - this apply to other things, and width as well, although a lot of the time GetLeft() and GetTop() would be 0
// also note that these return the actual size of the texture, not the desired height, which might not match perfectly.
func (t *Text) GetWidth() int {
	t.getBounds()
	return int((*t.bounds).GetLeft() + (*t.bounds).GetWidth())
}
func (t *Text) GetHeight() int {
	t.getBounds()
	return int((*t.bounds).GetTop() + (*t.bounds).GetHeight())
}

func (t *Text) InvalidateRenderCache() bool {
	DeleteCanvas(&t.canvas)
	if t.bounds != nil {
		graphics.DeleteSfFloatRect(*t.bounds)
		t.bounds = nil
	}
	return true
}

func (t *Text) Draw() {
	if t.canvas.IsRendered() {
		return
	}
	if t.canvas == nil {
		t.canvas = NewCanvas(t)
	}
	t.canvas.DrawText(t)
	t.canvas.Finalize()
}

func (t *Text) ToTexture() *Texture {
	t.Draw()
	texture := t.canvas.AsTexture()
	return texture
}

// ChildObject interface
func (t *Text) GetParent() Object {
	return t.Parent
}
