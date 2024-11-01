package main

import (
	"elichika/gui/graphic"
	"elichika/gui/graphic/button"
	"elichika/gui/graphic/textbox"

	"fmt"
	"os"
	"strconv"

	_ "modernc.org/sqlite"
	"xorm.io/xorm"
)

var statusChannel = make(chan string)

var engine *xorm.Engine
var session *xorm.Session

func flush() {
	session.Commit()
	err := session.Begin()
	if err != nil {
		panic(err)
	}
}

var staticFileDirs = []string{"static/b66ec2295e9a00aa/", "static/2d61e7b4e89961c7/"}

func init() {
	var err error
	engine, err = xorm.NewEngine("sqlite", "texture.db")
	if err != nil {
		panic(err)
	}
	exist, err := engine.Table("detailed_texture").IsTableExist("detailed_texture")
	if err != nil {
		panic(err)
	}
	if !exist {
		err = engine.Table("detailed_texture").CreateTable(&DetailedTexture{})
		if err != nil {
			panic(err)
		}
	}
	exist, err = engine.Table("static_file_dir").IsTableExist("static_file_dir")
	if err != nil {
		panic(err)
	}
	if exist {
		err = engine.Table("static_file_dir").Cols("dir").Find(&staticFileDirs)
		if err != nil {
			panic(err)
		}
	}

	session = engine.NewSession()
	err = session.Begin()
	if err != nil {
		panic(err)
	}
}

// TODO(file): For now this doesn't handle metapack
func GetStaticFile(fileName string) string {
	for _, dir := range staticFileDirs {
		_, err := os.Stat(dir + fileName)
		if err != nil {
			continue
		}
		return dir + fileName
	}
	return ""
}


type TextureDisplay struct {
	StatusLabel *graphic.Text

	AssetLabel    *graphic.Text
	AssetTextbox  *textbox.RectTextbox
	LoadButton    *button.RectButton
	LocaleButton  *button.RectButton
	ExtractButton *button.RectButton

	SQLLabel          *graphic.Text
	SQLTextbox        *textbox.RectTextbox
	SQLButton         *button.RectButton
	SQLResult         []Texture
	SQLLength         int
	SQLIndex          int
	SQLNextButton     *button.RectButton
	SQLPreviousButton *button.RectButton

	// GapLabel          *graphic.Text
	// GapTextbox        *textbox.RectTextbox
	// GapFillButton     *button.RectButton
	// GapContinueButton *button.RectButton
	// GapStopButton     *button.RectButton
	// CurrentGap        Gap

	AssetTexture *graphic.Texture

	Canvas *graphic.Canvas

	objects []graphic.Object
	x       int
	y       int
	objectX map[graphic.Object]int
	objectY map[graphic.Object]int
}

func NewTextureDisplay() *TextureDisplay {
	t := &TextureDisplay{}
	t.objectX = map[graphic.Object]int{}
	t.objectY = map[graphic.Object]int{}
	return t
}

func (t *TextureDisplay) AddObject(object graphic.Object, gap int) {
	t.objects = append(t.objects, object)
	t.objectX[object] = t.x
	t.objectY[object] = t.y
	t.x += object.GetWidth() + gap
}

func (t *TextureDisplay) NewLine() {
	t.y += 50
	t.x = 0
}

func (t *TextureDisplay) InvalidateRenderCache() bool {
	return t.Canvas.InvalidateRenderCache()
}

func (t *TextureDisplay) ToTexture() *graphic.Texture {
	t.Draw()
	return t.Canvas.AsTexture()
}

func (t *TextureDisplay) GetWidth() int {
	return 1800
}

func (t *TextureDisplay) GetHeight() int {
	return 1000
}

func (t *TextureDisplay) Draw() {
	if t.Canvas.IsRendered() {
		return
	}
	if t.Canvas == nil {
		t.Canvas = graphic.NewCanvas(t)
	}
	for _, o := range t.objects {
		t.Canvas.DrawObject(o, t.objectX[o], t.objectY[o], o.GetWidth(), o.GetHeight())
	}

	if t.AssetTexture != nil {
		w, h := t.AssetTexture.GetSize()
		t.Canvas.DrawTexture(t.AssetTexture, 0, t.y, w, h)
	}

	t.Canvas.Finalize()

}

func (t *TextureDisplay) ForEach(f func(graphic.Object)) {
	for _, o := range t.objects {
		f(o)
	}
}

func (t *TextureDisplay) MapEvent(e graphic.Event, child graphic.Object) graphic.Event {
	switch e.(type) {
	case graphic.MouseButtonDownEvent:
		mouseEvent := graphic.MouseButtonDownEvent{}
		mouseEvent = e.(graphic.MouseButtonDownEvent)
		mouseEvent.X -= t.objectX[child]
		mouseEvent.Y -= t.objectY[child]
		return mouseEvent
	default:
		return e
	}
}

func GetInt32(s string) int32 {
	value, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return int32(value)
}

func main() {
	GetAllDetail()
	window, err := graphic.NewWindow("Texture viewer")
	if err != nil {
		panic(err)
	}
	object := NewTextureDisplay()
	object.AssetLabel, object.AssetTextbox = textbox.NewLabelAndRectTextbox(object, 600, 50, "Asset path: ")
	object.AssetTextbox.Texture = graphic.RGBATexture(0x0f0f0fff)
	object.AssetTextbox.FocusTexture = graphic.RGBATexture(0x7f7f7fff)

	object.StatusLabel = graphic.NewText(object, "Status: ")
	go func() {
		for {
			status := <-statusChannel
			fmt.Println("Status: " + status)
			window.InternalEvent(func() {
				object.StatusLabel.SetText("Status: " + status)
			})
		}
	}()
	object.AddObject(object.StatusLabel, 0)
	object.NewLine()

	object.AddObject(object.AssetLabel, 0)
	object.AddObject(object.AssetTextbox, 5)

	displayAsset := func() {
		texture, err := LoadTexture(object.AssetTextbox.TextContent)
		graphic.InvalidateRenderCache(object)
		if err != nil {
			// TODO(gui): Dialog box
			statusChannel <- fmt.Sprint("Failed to load asset: ", object.AssetTextbox.TextContent, ", error: ", err)
			return
		}
		object.AssetTexture.Free()
		object.AssetTexture = texture
	}

	object.LoadButton = &button.RectButton{
		Width:   200,
		Height:  50,
		Texture: graphic.RGBATexture(0x7f7f7fff),
	}

	object.AddObject(object.LoadButton, 5)

	object.LoadButton.Text = graphic.NewText(object.LoadButton, "Load asset")
	object.LoadButton.LeftClickHandler = displayAsset
	object.AssetTextbox.OnEnterFunc = displayAsset

	object.ExtractButton = &button.RectButton{
		Width:   200,
		Height:  50,
		Texture: graphic.RGBATexture(0x7f7f7fff),
	}
	object.ExtractButton.Text = graphic.NewText(object.ExtractButton, "Extract")
	object.ExtractButton.LeftClickHandler = func() {
		if object.AssetTexture != nil {
			object.AssetTexture.SaveToImage("texture_output.png")
		}
	}
	object.AddObject(object.ExtractButton, 5)
	object.NewLine()
	object.SetupSQL(window)
	// object.SetupGap(window)
	window.SetObject(object)
	window.DisplayWithChannel()
}
