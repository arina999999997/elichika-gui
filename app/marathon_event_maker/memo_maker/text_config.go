package memo_maker

import (
	"elichika/gui/graphic"
	"elichika/gui/graphic/button"
	"elichika/gui/graphic/textbox"

	"errors"
	"fmt"
)

func (m *MemoMaker) LoadText() error {
	font := graphic.GetFont(fmt.Sprintf("gui/fonts/%s", m.Font))
	if font == nil {
		return errors.New(fmt.Sprint("Error loading font: ", m.Font))
	} else {
		m.TextFont = font
	}
	return nil
}

func (m *MemoMaker) SetupText() {
	m.TextLabel, m.TextTextbox =
		textbox.NewLabelAndRectTextbox(m, MemoMakerWidth, 50, "Memo text: ")
	m.TextTextbox.SetStringSettingTextbox("Our surprise\\party was a\\success!", &m.Text)
	// m.TextTextbox.SetStringSettingTextbox("Text goes here, use \\ for line breaks", &m.Text)

	m.AddObject(m.TextLabel, true, 0, 0)
	m.AddObject(m.TextTextbox, false, 0, 0)

	m.TextXLabel, m.TextXTextbox =
		textbox.NewLabelAndRectTextbox(m, MemoMakerWidth/5, 50, "Text X: ")
	m.TextYLabel, m.TextYTextbox =
		textbox.NewLabelAndRectTextbox(m, MemoMakerWidth/5, 50, "Text Y: ")
	textbox.SetCoordinateTextboxes(m.TextXTextbox, m.TextYTextbox, 140, 310, 0, 0, 800, 800, &m.TextX, &m.TextY)
	m.TextRotationLabel, m.TextRotationTextbox =
		textbox.NewLabelAndRectTextbox(m, MemoMakerWidth/5, 50, "Rotation: ")
	m.TextRotationTextbox.SetFloatSettingTextbox(-3.8, -180, 180, 0.1, 1, &m.TextRotation)
	m.TextFontLabel, m.TextFontTextbox =
		textbox.NewLabelAndRectTextbox(m, MemoMakerWidth-(MemoMakerWidth/5*3)-15, 50, "Font: ")
	m.TextFontTextbox.SetStringSettingTextbox("FOT-SkipStd-B.otf", &m.Font)

	m.AddObject(m.TextXLabel, true, 0, 0)
	m.AddObject(m.TextXTextbox, false, 5, 0)
	m.AddObject(m.TextYLabel, false, 0, 0)
	m.AddObject(m.TextYTextbox, false, 5, 0)
	m.AddObject(m.TextRotationLabel, false, 0, 0)
	m.AddObject(m.TextRotationTextbox, false, 5, 0)
	m.AddObject(m.TextFontLabel, false, 0, 0)
	m.AddObject(m.TextFontTextbox, false, 5, 0)

	m.TextCharacterSizeLabel, m.TextCharacterSizeTextbox =
		textbox.NewLabelAndRectTextbox(m, MemoMakerWidth/4, 50, "Character Size: ")
	m.TextCharacterSizeTextbox.SetIntSettingTextbox(72, 1, 256, &m.TextCharacterSize)
	m.TextLetterSpacingLabel, m.TextLetterSpacingTextbox =
		textbox.NewLabelAndRectTextbox(m, MemoMakerWidth/4, 50, "Letter Spacing: ")
	m.TextLetterSpacingTextbox.SetFloatSettingTextbox(1.32, 0.5, 2, 0.01, 2, &m.TextLetterSpacing)
	m.TextLineSpacingLabel, m.TextLineSpacingTextbox =
		textbox.NewLabelAndRectTextbox(m, MemoMakerWidth-(MemoMakerWidth/4*3)-15, 50, "Line Spacing: ")
	m.TextLineSpacingTextbox.SetFloatSettingTextbox(0.72, 0.5, 2, 0.01, 2, &m.TextLineSpacing)
	m.TextColorLabel, m.TextColorTextbox =
		textbox.NewLabelAndRectTextbox(m, MemoMakerWidth/4, 50, "Color: ")
	m.TextColorTextbox.SetHexSettingTextbox(0x494949ff, &m.TextColor)

	m.AddObject(m.TextCharacterSizeLabel, true, 0, 0)
	m.AddObject(m.TextCharacterSizeTextbox, false, 5, 0)
	m.AddObject(m.TextLetterSpacingLabel, false, 0, 0)
	m.AddObject(m.TextLetterSpacingTextbox, false, 5, 0)
	m.AddObject(m.TextLineSpacingLabel, false, 0, 0)
	m.AddObject(m.TextLineSpacingTextbox, false, 5, 0)
	m.AddObject(m.TextColorLabel, false, 0, 0)
	m.AddObject(m.TextColorTextbox, false, 5, 0)

	m.TextScalingFactorLabel, m.TextScalingFactorTextbox = textbox.NewLabelAndRectTextbox(
		m, MemoMakerWidth/4-5, 50, "Text scaling factor: ")
	m.TextScalingFactorTextbox.SetIntSettingTextbox(4, 1, 8, &m.TextScalingFactor)
	m.TextLineOffsetLabel, m.TextLineOffsetTextboxes = textbox.NewLabelAndRectTextboxes(
		m, (MemoMakerWidth+1)/2, 50, "Line offsets: ", MaxTextLine, (MaxTextLine-1)*5)

	m.TextScaleBackground = true
	m.TextScaleBackgroundButton = button.NewButton(m, MemoMakerWidth/4-5, 50, "Background scaling on", func() {
		graphic.InvalidateRenderCache(m.TextScaleBackgroundButton)
		if m.TextScaleBackground {
			m.TextScaleBackground = false
			m.TextScaleBackgroundButton.SetTextString("Background scaling off")
		} else {
			m.TextScaleBackground = true
			m.TextScaleBackgroundButton.SetTextString("Background scaling on")
		}
	}, nil)
	m.TextScaleBackgroundButton.SyncFunc = func() {
		if m.TextScaleBackground {
			m.TextScaleBackgroundButton.SetTextString("Background scaling on")
		} else {
			m.TextScaleBackgroundButton.SetTextString("Background scaling off")
		}
	}
	m.AddObject(m.TextScalingFactorLabel, true, 0, 0)
	m.AddObject(m.TextScalingFactorTextbox, false, 5, 0)
	m.AddObject(m.TextScaleBackgroundButton, false, 5, 0)
	m.AddObject(m.TextLineOffsetLabel, false, 0, 0)
	for i := 0; i < MaxTextLine; i++ {
		textbox := m.TextLineOffsetTextboxes[i]
		m.AddObject(textbox, false, 5, 0)
		textbox.SetIntSettingTextbox(0, -100, 100, &m.TextLineOffsets[i])
	}
}
