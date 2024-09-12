package memo_maker

import (
	"elichika/gui/graphic"

	"strings"
)

// receive lines as a slice of string
// each line is rendered to its own object, then the lines are stacked together and centered
type MultilineCenteredText struct {
	LineCount int
	Lines     []string
	TextLines []*graphic.Text

	Font              *graphic.Font
	CharacterSize     int
	LetterSpacing     float32
	LineSpacingFactor float32

	Canvas      *graphic.Canvas
	Width       int
	Height      int
	Color       uint
	LineOffsets []int

	Built bool
}

func (mct *MultilineCenteredText) GetWidth() int {
	mct.Build()
	return mct.Width
}

func (mct *MultilineCenteredText) GetHeight() int {
	mct.Build()
	return mct.Height
}

func (mct *MultilineCenteredText) InvalidateRenderCache() bool {
	return mct.Canvas.InvalidateRenderCache()
}

func (mct *MultilineCenteredText) ToTexture() *graphic.Texture {
	mct.Draw()
	return mct.Canvas.AsTexture()
}

func (mct *MultilineCenteredText) Draw() {
	mct.Build()
	if mct.Canvas.IsRendered() {
		return
	}
	if mct.Canvas != nil {
		graphic.DeleteCanvas(&mct.Canvas)
	}
	mct.Canvas = graphic.NewCanvas(mct)

	y := 0
	for i := 0; i < mct.LineCount; i++ {
		mct.Canvas.DrawObject(mct.TextLines[i], (mct.Width-mct.TextLines[i].GetWidth())/2+mct.LineOffsets[i], y,
			mct.TextLines[i].GetWidth(), mct.TextLines[i].GetHeight())
		y += mct.TextLines[i].GetLineSpacing()
	}
	mct.Canvas.Finalize()
}

func (mct *MultilineCenteredText) Build() {
	if mct.Built {
		return
	}
	mct.LineCount = len(mct.Lines)
	for i := 0; i < mct.LineCount; i++ {
		if len(mct.TextLines) == i {
			mct.TextLines = append(mct.TextLines, graphic.NewText(mct, mct.Lines[i]))
		} else {
			mct.TextLines[i].SetText(mct.Lines[i])
		}
	}

	// set the line spacing and characte spacing

	mct.Height = 0
	mct.Width = 0
	for i := 0; i < mct.LineCount; i++ {
		mct.TextLines[i].SetCharacterSize(mct.CharacterSize)
		mct.TextLines[i].SetLetterSpacing(mct.LetterSpacing)
		mct.TextLines[i].SetLineSpacingFactor(mct.LineSpacingFactor)
		mct.TextLines[i].SetFont(mct.Font)
		mct.TextLines[i].SetColor(mct.Color)
		if i+1 < mct.LineCount {
			mct.Height += mct.TextLines[i].GetLineSpacing()
		} else {
			mct.Height += mct.TextLines[i].GetHeight()
		}
		if mct.Width < mct.TextLines[i].GetWidth() {
			mct.Width = mct.TextLines[i].GetWidth()
		}
	}
	mct.Built = true
}

func (mct *MultilineCenteredText) Load(lines string,
	font *graphic.Font,
	characterSize int,
	letterSpacing,
	lineSpacingFactor float32,
	color uint,
	lineOffsets []int) {
	mct.Lines = strings.Split(lines, "\\")
	mct.Font = font
	mct.CharacterSize = characterSize
	mct.LetterSpacing = letterSpacing
	mct.LineSpacingFactor = lineSpacingFactor
	mct.Color = color
	mct.LineOffsets = lineOffsets
	mct.Built = false
}
