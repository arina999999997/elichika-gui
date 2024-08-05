package graphic

import (
	"errors"
	"fmt"

	"github.com/telroshan/go-sfml/v2/graphics"
	"github.com/telroshan/go-sfml/v2/window"
)

// this is the window to render things in
// this object is just a wrapper for whatever window handler types we use
// it will have no bearing on the actual object compositions.
type Window struct {
	title string
	vm    window.SfVideoMode
	cs    window.SfContextSettings
	rw    graphics.Struct_SS_sfRenderWindow
	rt    graphics.Struct_SS_sfRenderTexture
}

func NewWindow(title string) (*Window, error) {
	w := Window{
		title: title,
	}
	w.vm = window.NewSfVideoMode()
	// defer window.DeleteSfVideoMode(vm)
	w.vm.SetBitsPerPixel(32)
	w.cs = window.NewSfContextSettings()
	// defer window.DeleteSfContextSettings(cs)
	w.rw = graphics.SfRenderWindow_create(w.vm, title, uint(window.SfTitlebar), w.cs)
	graphics.SfRenderWindow_setFramerateLimit(w.rw, 60)
	// defer window.SfWindow_destroy(w)

	err := recover()
	if err == nil {
		return &w, nil
	} else {
		return nil, errors.New(fmt.Sprint(err))
	}
}

func (w *Window) RenderRootObject(obj *Object) {
	if (w.vm.GetWidth() != uint(obj.Width)) || (w.vm.GetHeight() != uint(obj.Height)) {
		w.vm.SetWidth(uint(obj.Width))
		w.vm.SetHeight(uint(obj.Height))
		window.SfWindow_destroy(w.rw)
		w.rw = graphics.SfRenderWindow_create(w.vm, w.title, uint(window.SfTitlebar), w.cs)
		w.rt = graphics.SfRenderTexture_create(w.vm.GetWidth(), w.vm.GetHeight(), 0)
	}
	ev := window.NewSfEvent()
	defer window.DeleteSfEvent(ev)

	// for window.SfWindow_pollEvent(w.rw, ev) > 0 {
	// 		/* Close window: exit */
	// if ev.GetEvType() == window.SfEventType(window.SfEvtClosed) {
	// return
	// }
	// }
	// graphics.SfRenderWindow_clear(w.rw, graphics.GetSfTransparent())
	graphics.SfRenderWindow_clear(w.rw, graphics.GetSfTransparent())
	graphics.SfRenderTexture_clear(w.rt, graphics.GetSfTransparent())
	obj.draw(w, 0, 0, obj.Width, obj.Height)
	w.SaveToImage("output.png")
	graphics.SfRenderWindow_display(w.rw)
	// TODO(debug): For now, this show the image until closed
	for graphics.SfRenderWindow_waitEvent(w.rw, ev) > 0 {
		if ev.GetEvType() == window.SfEventType(window.SfEvtClosed) {
			break
		}
	}
}

func (w *Window) SaveToImage(path string) {
	graphics.SfRenderTexture_display(w.rt)
	texture := graphics.SfRenderTexture_getTexture(w.rt)
	image := graphics.SfTexture_copyToImage(texture)
	graphics.SfImage_saveToFile(image, path)
}
