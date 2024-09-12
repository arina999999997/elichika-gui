package memo_maker

import (
	"elichika/gui/graphic"
	"elichika/gui/graphic/button"
	"elichika/gui/graphic/textbox"

	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

const (
	MemoMakerWidth  = 1920
	MemoMakerHeight = 1080
	MaxTextLine     = 6

	MemoWidth  = 200
	MemoHeight = 200

	GIMP_CONSOLE = `C:\Program Files\GIMP 2\bin\gimp-console-2.10`
)

type MmeoMakerConfig struct {
}

type MemoMaker struct {
	Parent graphic.Object  `json:"-"`
	Canvas *graphic.Canvas `json:"-"`
	// special canvas used for exporting output with gimp
	GimpCanvas *graphic.Canvas `json:"-"`

	// values
	BackgroundPath       string           `json:"background_path"`
	BackgroundLoadedPath string           `json:"-"`
	BackgroundTexture    *graphic.Texture `json:"-"`

	PinPath       string           `json:"pin_path"`
	PinLoadedPath string           `json:"-"`
	PinTexture    *graphic.Texture `json:"-"`
	PinX          int              `json:"pin_x"`
	PinY          int              `json:"pin_y"`

	Text                string           `json:"text"`
	TextX               int              `json:"text_x"`
	TextY               int              `json:"text_y"`
	TextRotation        float32          `json:"text_rotation"`
	TextCharacterSize   int              `json:"text_character_size"`
	TextScalingFactor   int              `json:"text_scaling_factor"`
	TextLetterSpacing   float32          `json:"text_letter_spacing"`
	TextLineSpacing     float32          `json:"text_line_spacing"`
	TextColor           uint             `json:"text_color"`
	TextLineOffsets     [MaxTextLine]int `json:"text_line_offsets"`
	Font                string           `json:"font"`
	TextFont            *graphic.Font    `json:"-"`
	TextScaleBackground bool             `json:"text_scale_background"`

	MemoDirectory       string  `json:"memo_directory"`
	VisualScalingFactor float32 `json:"visual_scaling_factor"`

	// interfaces
	BackgroundLabel   *graphic.Text        `json:"-"`
	BackgroundTextbox *textbox.RectTextbox `json:"-"`

	PinTextureLabel   *graphic.Text        `json:"-"`
	PinTextureTextbox *textbox.RectTextbox `json:"-"`

	PinXLabel   *graphic.Text        `json:"-"`
	PinXTextbox *textbox.RectTextbox `json:"-"`
	PinYLabel   *graphic.Text        `json:"-"`
	PinYTextbox *textbox.RectTextbox `json:"-"`

	TextLabel   *graphic.Text        `json:"-"`
	TextTextbox *textbox.RectTextbox `json:"-"`

	TextXLabel          *graphic.Text        `json:"-"`
	TextXTextbox        *textbox.RectTextbox `json:"-"`
	TextYLabel          *graphic.Text        `json:"-"`
	TextYTextbox        *textbox.RectTextbox `json:"-"`
	TextRotationLabel   *graphic.Text        `json:"-"`
	TextRotationTextbox *textbox.RectTextbox `json:"-"`
	TextFontLabel       *graphic.Text        `json:"-"`
	TextFontTextbox     *textbox.RectTextbox `json:"-"`

	TextCharacterSizeLabel   *graphic.Text        `json:"-"`
	TextCharacterSizeTextbox *textbox.RectTextbox `json:"-"`
	TextLetterSpacingLabel   *graphic.Text        `json:"-"`
	TextLetterSpacingTextbox *textbox.RectTextbox `json:"-"`
	TextLineSpacingLabel     *graphic.Text        `json:"-"`
	TextLineSpacingTextbox   *textbox.RectTextbox `json:"-"`
	TextColorLabel           *graphic.Text        `json:"-"`
	TextColorTextbox         *textbox.RectTextbox `json:"-"`

	TextScalingFactorLabel    *graphic.Text          `json:"-"`
	TextScalingFactorTextbox  *textbox.RectTextbox   `json:"-"`
	TextScaleBackgroundButton *button.RectButton     `json:"-"`
	TextLineOffsetLabel       *graphic.Text          `json:"-"`
	TextLineOffsetTextboxes   []*textbox.RectTextbox `json:"-"`

	VisualScalingLabel   *graphic.Text        `json:"-"`
	VisualScalingTextbox *textbox.RectTextbox `json:"-"`

	MemoDirectoryLabel     *graphic.Text        `json:"-"`
	MemoDirectoryTextbox   *textbox.RectTextbox `json:"-"`
	MemoExportButton       *button.RectButton   `json:"-"`
	MemoImportConfigButton *button.RectButton   `json:"-"`

	// drawing objects
	TextOutput *TextOutput `json:"-"`
	Memo       *Memo       `json:"-"`

	// internals
	lastX         int                    `json:"-"`
	lastY         int                    `json:"-"`
	currentHeight int                    `json:"-"`
	Objects       []graphic.Object       `json:"-"`
	ObjectX       map[graphic.Object]int `json:"-"`
	ObjectY       map[graphic.Object]int `json:"-"`
}

func (m *MemoMaker) AddObject(object graphic.Object, newLine bool, paddingX, paddingY int) {
	if newLine {
		m.lastY += m.currentHeight
		m.lastX = 0
		m.currentHeight = 0
	}
	m.Objects = append(m.Objects, object)
	m.ObjectX[object] = m.lastX
	m.ObjectY[object] = m.lastY

	m.lastX += object.GetWidth() + paddingX
	if m.currentHeight < object.GetHeight()+paddingY {
		m.currentHeight = object.GetHeight() + paddingY
	}
}

func (*MemoMaker) GetWidth() int {
	return MemoMakerWidth
}
func (*MemoMaker) GetHeight() int {
	return MemoMakerHeight
}

func (m *MemoMaker) InvalidateRenderCache() bool {
	result := m.Canvas.InvalidateRenderCache()
	// graphic.InvalidateRenderCache(m.Memo)
	if m.Memo != nil {
		graphic.InvalidateRenderCache(m.Memo)
	}
	return result
}

func (m *MemoMaker) ToTexture() *graphic.Texture {
	m.Draw()
	return m.Canvas.AsTexture()
}

func (m *MemoMaker) Draw() {
	if m.Canvas.IsRendered() {
		return
	}
	if m.Canvas == nil {
		m.Canvas = graphic.NewCanvas(m)
	}

	err := m.Load()
	if err != nil {
		fmt.Println("Error generating output: ", err)
	}

	m.ForEach(func(object graphic.Object) {
		if object == m.Memo {
			return
		}
		m.Canvas.DrawObject(object, m.ObjectX[object], m.ObjectY[object], object.GetWidth(), object.GetHeight())
	})
	m.Canvas.DrawObject(m.Memo, m.ObjectX[m.Memo], m.ObjectY[m.Memo],
		int(float32(m.Memo.GetWidth())*m.VisualScalingFactor), int(float32(m.Memo.GetHeight())*m.VisualScalingFactor))
	m.Canvas.Finalize()
}

func (m *MemoMaker) ForEach(f func(graphic.Object)) {
	for _, object := range m.Objects {
		f(object)
	}
}

func (m *MemoMaker) MapEvent(e graphic.Event, child graphic.Object) graphic.Event {
	switch e.(type) {
	case graphic.MouseButtonDownEvent:
		mouseEvent := graphic.MouseButtonDownEvent{}
		mouseEvent = e.(graphic.MouseButtonDownEvent)
		mouseEvent.X -= m.ObjectX[child]
		mouseEvent.Y -= m.ObjectY[child]
		return mouseEvent
	default:
		return e
	}
}

func (m *MemoMaker) Load() error {
	var err error
	err = m.LoadBackground()
	if err != nil {
		return err
	}
	err = m.LoadPin()
	if err != nil {
		return err
	}
	err = m.LoadText()
	if err != nil {
		return err
	}

	graphic.InvalidateRenderCache(m.Memo)
	graphic.InvalidateRenderCache(m.TextOutput)

	err = m.LoadOutput()
	if err != nil {
		return err
	}
	return nil
}

func NewMemoMaker(parent graphic.Object) *MemoMaker {
	memoMaker := &MemoMaker{
		Parent:  parent,
		ObjectX: map[graphic.Object]int{},
		ObjectY: map[graphic.Object]int{},
	}
	memoMaker.Memo = &Memo{
		Maker: memoMaker,
	}
	memoMaker.TextOutput = &TextOutput{
		Maker: memoMaker,
	}

	memoMaker.SetupBackground()
	memoMaker.SetupPin()
	memoMaker.SetupText()
	memoMaker.SetupOutput()

	memoMaker.AddObject(memoMaker.Memo, true, 0, 0)
	return memoMaker
}

func (m *MemoMaker) ExportConfig(path string) {
	jsonOutput, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(path, jsonOutput, 0644)
	if err != nil {
		panic(err)
	}
}

func (m *MemoMaker) ImportConfig() {
	jsonInput, err := os.ReadFile(m.MemoDirectory + "/memo_maker_config.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(jsonInput, m)
	if err != nil {
		panic(err)
	}
	for _, object := range m.Objects {
		switch object.(type) {
		case *textbox.RectTextbox:
			object.(*textbox.RectTextbox).TrySync()
		case *button.RectButton:
			object.(*button.RectButton).TrySync()
		case *graphic.Text:
		case *Memo:
		default:
			fmt.Printf("unsynced: %T\n", object)
		}
	}
}

func (m *MemoMaker) Export() {
	err := os.MkdirAll(m.MemoDirectory, os.ModePerm)
	if err != nil {
		panic(err)
	}
	m.ExportConfig(m.MemoDirectory + "/memo_maker_config.json")
	m.Memo.ToTexture().SaveToImage(m.MemoDirectory + "/output.png")

	// gimp export
	gimpOutputPath := m.MemoDirectory + "/output_gimp.png"
	defer func() {
		err := recover()
		if err == nil {
			return
		}
		fmt.Println("Error trying to generate gimp output: ", err)
		err = os.Remove(gimpOutputPath)
		if err != nil {
			panic(err)
		}
	}()

	// first export the text, maybe along with the background
	m.TextOutput.ToTexture().SaveToImage(gimpOutputPath)
	// then use gimp to resize the text
	gimpCommand := exec.Command(GIMP_CONSOLE, `-i`, `-b`,
		fmt.Sprintf(`(memo-scale "%s")`, gimpOutputPath), `-b`, `(gimp-quit 0)`)
	err = gimpCommand.Run()
	if err != nil {
		panic(err)
	}

	// finally draw the pin on top of the generated text
	if m.GimpCanvas == nil {
		m.GimpCanvas = graphic.NewCanvas(m.Memo)
	} else {
		m.GimpCanvas.InvalidateRenderCache()
	}
	texture, err := graphic.FileTexture(gimpOutputPath)
	if err != nil {
		panic(err)
	}
	m.GimpCanvas.DrawTexture(texture, 0, 0, 200, 200)
	m.PinTexture.StyleType = graphic.StyleTypeNone
	m.GimpCanvas.DrawTexture(m.PinTexture, m.PinX, m.PinY, 1, 1)
	m.GimpCanvas.Finalize()
	m.GimpCanvas.AsTexture().SaveToImage(gimpOutputPath)
}
