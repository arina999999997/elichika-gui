package asset

// load sifas texture into graphic.Texture or graphic.Object
import (
	"elichika/gui/graphic"
	"elichika/gui/sifas/locale"

	"errors"
	"fmt"
)

func LoadTexture(assetPath string) (texture *graphic.Texture, err error) {
	defer func() {
		r := recover()
		if r != nil {
			texture = nil
			err = errors.New(fmt.Sprint(err))
		}
	}()
	loadLocale()
	rawData := AssetMap[locale.Locale()][assetPath].LoadUnencrypted()
	texture = &graphic.Texture{}
	texture.LoadFromMemory(rawData)
	return
}
