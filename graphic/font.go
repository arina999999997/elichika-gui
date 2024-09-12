package graphic

import (
	"os"

	"github.com/telroshan/go-sfml/v2/graphics"
)

type sfFont = graphics.Struct_SS_sfFont
type Font struct {
	font sfFont
}

var defaultFont = *GetFont("gui/fonts/FOT-SkipStd-B.otf")
var loadedFont = map[string]*Font{}

func GetFont(path string) *Font {
	font, exist := loadedFont[path]
	if exist {
		return font
	}
	_, err := os.Stat(path)
	if err != nil {
		return nil
	}

	font = &Font{}
	font.font = graphics.SfFont_createFromFile(path)
	loadedFont[path] = font
	return font
}

func GetDefaultFont() *Font {
	return &defaultFont
}
