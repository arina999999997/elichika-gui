package graphic

// The texture type need to implement the following function:
// - LoadFromFile: Load the texture from a file on disk. (TODO: maybe this could be an url too)
//   - it should panic if the file isn't present
// - LoadFromMemory: Load the texture from a chunk of memory (passed in as a byte array):
//   - this allow us to load from memory, and to do preprocessing like decrypting assets.
// a texture by itself can't be drawn, it will need a managing object

import (
	"os"
	"unsafe"

	"github.com/telroshan/go-sfml/v2/graphics"
)

type sfTexture = graphics.Struct_SS_sfTexture

type Texture struct {
	texture sfTexture
}

func (t *Texture) LoadFromFile(path string) {
	_, err := os.Stat(path)
	if err != nil {
		panic(err)
	}
	t.texture = graphics.SfTexture_createFromFile(path, graphics.NewSfIntRect())
}

func (t *Texture) LoadFromMemory(data []byte) {
	t.texture = graphics.SfTexture_createFromMemory(uintptr(unsafe.Pointer(&data[0])), int64(len(data)), graphics.NewSfIntRect())
}

func (t *Texture) GetSize() (int, int) {
	vector := graphics.SfTexture_getSize(t.texture)
	return int(vector.GetX()), int(vector.GetY())
}

var defaultTexture *Texture

func DefaultTexture() *Texture {
	if defaultTexture == nil {
		defaultTexture = &Texture{}
		defaultTexture.LoadFromFile("gui/graphic/missing.png")
	}
	return defaultTexture

}
