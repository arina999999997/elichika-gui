package memo_maker

import (
	"elichika/gui/graphic/button"
	"elichika/gui/graphic/textbox"
)

func (m *MemoMaker) LoadOutput() error {
	m.MemoDirectory = m.MemoDirectoryTextbox.TextContent
	return nil
}

func (m *MemoMaker) SetupOutput() {
	m.VisualScalingLabel, m.VisualScalingTextbox =
		textbox.NewLabelAndRectTextbox(m, MemoMakerWidth, 50, "Output scaling factor (preview only): ")
	m.VisualScalingTextbox.SetFloatSettingTextbox(1.2, 1, 4, 0.1, 1, &m.VisualScalingFactor)

	m.AddObject(m.VisualScalingLabel, true, 0, 0)
	m.AddObject(m.VisualScalingTextbox, false, 0, 0)

	m.MemoDirectoryLabel, m.MemoDirectoryTextbox =
		textbox.NewLabelAndRectTextbox(m, MemoMakerWidth-400-10, 50, "Memo directory: ")
	m.MemoDirectoryTextbox.SetStringSettingTextbox("memo_maker_output", &m.MemoDirectory)

	m.MemoExportButton = button.NewButton(m, 200, 50, "Export", func() {
		m.Export()
	}, nil)
	m.MemoImportConfigButton = button.NewButton(m, 200, 50, "Import", func() {
		m.ImportConfig()
	}, nil)

	m.AddObject(m.MemoDirectoryLabel, true, 0, 0)
	m.AddObject(m.MemoDirectoryTextbox, false, 5, 0)
	m.AddObject(m.MemoExportButton, false, 5, 0)
	m.AddObject(m.MemoImportConfigButton, false, 0, 0)
}
