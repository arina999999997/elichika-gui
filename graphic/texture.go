package graphic

// The texture type need to implement the following function:
// - LoadFromFile: Load the texture from a file on disk. (TODO: maybe this could be an url too)
//   - it should panic if the file isn't present
// - LoadFromMemory: Load the texture from a chunk of memory (passed in as a byte array):
//   - this allow us to load from memory, and to do preprocessing like decrypting assets

import (
	"elichika/utils"

	"errors"
	"fmt"
	"os"
	"unsafe"

	"github.com/telroshan/go-sfml/v2/graphics"
)

type sfTexture = graphics.Struct_SS_sfTexture

const (
	StyleTypeFitContainer int = 0 // fit inside the desired size, so apply the lower scale between width and height
	StyleTypeFitWidth     int = 1 // keeping the aspect ratio, the width is scale to the desired size, the height is not considered
	StyleTypeFitHeight    int = 2 // keeping the aspect ratio, the height is scale to the desired size, the width is not considered
	StyleTypeNone         int = 3 // no scaling, draw the texture as is
	StyleTypeAutoCentered int = 4 // scale down to fit, then center the texture
	StyleTypeIndependent  int = 5 // scahe each dimension independently
	StyleTypeRepeat       int = 6 // repeat / tile the desired area at the native resolution

)

type Texture struct {
	Texture   sfTexture
	StyleType int
}

var textureCnt = 0

func (t *Texture) LoadFromFile(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	if info.IsDir() {
		return errors.New("path is directory")
	}
	textureCnt++
	fmt.Println("texture count: ", textureCnt)
	t.Texture = graphics.SfTexture_createFromFile(path, graphics.NewSfIntRect())
	graphics.SfTexture_setSmooth(t.Texture, 1)
	return nil
}

func (t *Texture) LoadFromMemory(data []byte) {
	textureCnt++
	fmt.Println("texture count: ", textureCnt)
	t.Texture = graphics.SfTexture_createFromMemory(uintptr(unsafe.Pointer(&data[0])), int64(len(data)), graphics.NewSfIntRect())
	graphics.SfTexture_setSmooth(t.Texture, 1)
}

func (t *Texture) Free() {
	if t == nil {
		return
	}
	if t.Texture != nil {
		graphics.SfTexture_destroy(t.Texture)
		textureCnt--
		t.Texture = nil
	}
}

func (t *Texture) SetStyleType(styleType int) {
	t.StyleType = styleType
}

func (t *Texture) GetSize() (int, int) {
	vector := graphics.SfTexture_getSize(t.Texture)
	return int(vector.GetX()), int(vector.GetY())
}

func (t *Texture) SaveToImage(path string) {
	image := graphics.SfTexture_copyToImage(t.Texture)
	graphics.SfImage_saveToFile(image, path)
}

var defaultTexture *Texture

func DefaultTexture() *Texture {
	if defaultTexture == nil {
		defaultTexture = &Texture{}
		err := defaultTexture.LoadFromFile("gui/graphic/missing.png")
		utils.CheckErr(err)
		defaultTexture.SetStyleType(StyleTypeRepeat)
	}
	return defaultTexture
}

func (t *Texture) SetSmooth(smooth bool) {
	if smooth {
		graphics.SfTexture_setSmooth(t.Texture, 1)
	} else {
		graphics.SfTexture_setSmooth(t.Texture, 0)
	}
}

// rgba is provided using a number:
// - 0xRRGGBBAA
// where RR, GG, BB, AA is the hex values of each component
func RGBATexture(rgba uint32) *Texture {
	texture := &Texture{}
	texture.Texture = graphics.SfTexture_create(1, 1)
	bytes := make([]byte, 4)
	bytes[0] = byte((rgba >> 24) % 256)
	bytes[1] = byte((rgba >> 16) % 256)
	bytes[2] = byte((rgba >> 8) % 256)
	bytes[3] = byte((rgba >> 0) % 256)
	graphics.SfTexture_updateFromPixels(texture.Texture, &bytes[0], 1, 1, 0, 0)
	texture.SetStyleType(StyleTypeRepeat)
	return texture
}

func FileTexture(path string) (*Texture, error) {
	texture := &Texture{}
	err := texture.LoadFromFile(path)
	return texture, err
}
