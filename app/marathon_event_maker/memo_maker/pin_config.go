package memo_maker

import (
	"elichika/gui/graphic"
	"elichika/gui/graphic/textbox"

	"errors"
	"fmt"
)

func (m *MemoMaker) LoadPin() error {
	var err error
	if (m.PinPath != m.PinLoadedPath) || (m.PinTexture == nil) {
		m.PinLoadedPath = m.PinPath
		m.PinTexture, err = graphic.FileTexture(m.PinPath)
		if err != nil {
			return errors.New(fmt.Sprint("Error while loading pin texture: ", err))
		}
	}
	return nil
}

func (m *MemoMaker) SetupPin() {
	m.PinTextureLabel, m.PinTextureTextbox =
		textbox.NewLabelAndRectTextbox(m, MemoMakerWidth*6/8-10, 50, "Memo pin: ")
	m.PinTextureTextbox.SetStringSettingTextbox("gui/app/marathon_event_maker/memo_maker/pink_pin.png", &m.PinPath)

	m.PinXLabel, m.PinXTextbox =
		textbox.NewLabelAndRectTextbox(m, MemoMakerWidth*1/8, 50, "Pin X: ")
	m.PinYLabel, m.PinYTextbox =
		textbox.NewLabelAndRectTextbox(m, MemoMakerWidth*1/8, 50, "Pin Y: ")

	textbox.SetCoordinateTextboxes(m.PinXTextbox, m.PinYTextbox, 0, 0, -200, -200, 200, 200, &m.PinX, &m.PinY)

	m.AddObject(m.PinTextureLabel, true, 0, 0)
	m.AddObject(m.PinTextureTextbox, false, 5, 0)
	m.AddObject(m.PinXLabel, false, 0, 0)
	m.AddObject(m.PinXTextbox, false, 5, 0)
	m.AddObject(m.PinYLabel, false, 0, 0)
	m.AddObject(m.PinYTextbox, false, 5, 0)
}
