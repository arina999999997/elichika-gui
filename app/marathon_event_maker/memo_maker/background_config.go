package memo_maker

import (
	"elichika/gui/graphic"
	"elichika/gui/graphic/textbox"

	"errors"
	"fmt"
)

func (m *MemoMaker) LoadBackground() error {
	var err error
	if (m.BackgroundPath != m.BackgroundLoadedPath) || (m.BackgroundTexture == nil) {
		m.BackgroundLoadedPath = m.BackgroundPath
		m.BackgroundTexture, err = graphic.FileTexture(m.BackgroundPath)
		if err != nil {
			return errors.New(fmt.Sprint("Error while loading background texture: ", err))
		}
	}
	return nil
}

func (m *MemoMaker) SetupBackground() {
	m.BackgroundLabel, m.BackgroundTextbox = textbox.NewLabelAndRectTextbox(m, MemoMakerWidth, 50, "Memo background: ")
	m.BackgroundTextbox.SetStringSettingTextbox("gui/app/marathon_event_maker/memo_maker/common_note.png", &m.BackgroundPath)
	m.AddObject(m.BackgroundLabel, true, 0, 0)
	m.AddObject(m.BackgroundTextbox, false, 0, 0)
}
