package graphic

import (
	"errors"
	"fmt"

	"github.com/telroshan/go-sfml/v2/graphics"
	"github.com/telroshan/go-sfml/v2/window"
)

// this is the window to render things
// or more precisely, it's used for event polling and real-time displaying
// all the render are done by canvas.
// the window will draw an object and that's it.
// the window have its own size, the object will be rescaled to display on the window:
// - the rescale always keep the aspect ratio of the object
// - it will keep the height or the width of the window, whichever result in smaller dimentions
// the window will display the object's texture as is, it won't follow any style the object's texture might prefer.

type Window struct {
	title          string
	videoMode      window.SfVideoMode
	contextSetting window.SfContextSettings
	renderWindow   graphics.Struct_SS_sfRenderWindow
	object         Object

	focusObject Focusable
}

func NewWindow(title string) (*Window, error) {
	w := Window{
		title: title,
	}
	w.contextSetting = GetContextSetting()
	err := recover()
	if err == nil {
		return &w, nil
	} else {
		return nil, errors.New(fmt.Sprint(err))
	}
}

func (w *Window) SetObject(object Object) {
	w.object = object
	w.UpdateNativeSize()
}

func (w *Window) SetFocusObject(focus Focusable) {
	if w.focusObject == focus {
		return
	}
	if w.focusObject != nil {
		w.focusObject.UnsetFocus()
	}
	w.focusObject = focus
	w.focusObject.SetFocus()
}

func (w *Window) UpdateNativeSize() {
	if (w.videoMode != nil) && (int(w.videoMode.GetWidth()) == w.object.GetWidth()) && (int(w.videoMode.GetHeight()) == w.object.GetHeight()) {
		return
	}
	if w.videoMode != nil {
		window.DeleteSfVideoMode(w.videoMode)
	}
	w.videoMode = window.NewSfVideoMode()
	w.videoMode.SetBitsPerPixel(32)
	w.videoMode.SetWidth(uint(w.object.GetWidth()))
	w.videoMode.SetHeight(uint(w.object.GetHeight()))
	w.renderWindow = graphics.SfRenderWindow_create(w.videoMode, w.title, uint(window.SfTitlebar)|uint(window.SfResize), w.contextSetting)
	graphics.SfRenderWindow_setFramerateLimit(w.renderWindow, 60)
}

// this is the actual size of the window in the OS
func (w *Window) GetSize() (int, int) {
	vec := graphics.SfRenderWindow_getSize(w.renderWindow)
	return int(vec.GetX()), int(vec.GetY())
}

// all the drawing and scaling of object should follow this size
func (w *Window) GetNativeSize() (int, int) {
	return int(w.videoMode.GetWidth()), int(w.videoMode.GetHeight())
}

// update the width to keep aspect ratio while keeping the height
func (w *Window) UpdateWidthFromHeight() {
	width, height := w.GetSize()
	nWidth, nHeight := w.GetNativeSize()
	width = (nWidth*height-1)/nHeight + 1
	graphics.SfRenderWindow_setSize(w.renderWindow, getVector2u(width, height))
}

// update the width to keep aspect ratio while keeping the height
func (w *Window) UpdateHeightFromWidth() {
	width, height := w.GetSize()
	nWidth, nHeight := w.GetNativeSize()
	height = (nHeight*width-1)/nWidth + 1
	graphics.SfRenderWindow_setSize(w.renderWindow, getVector2u(width, height))
}

// update either width or height to keep the aspect ratio, while keeping the other the same
// this can only reduce the window size, i.e. keep whatever size is "smaller"
func (w *Window) UpdateWidthOrHeight() {
	width, height := w.GetSize()
	nWidth, nHeight := w.GetNativeSize()
	if height*nWidth > width*nHeight {
		height = (nHeight*width-1)/nWidth + 1
	} else {
		width = (nWidth*height-1)/nHeight + 1
	}
	graphics.SfRenderWindow_setSize(w.renderWindow, getVector2u(width, height))
}

func (w *Window) Render() {
	// this function always rerender the window
	// the object might have render cache
	// only render when necessary
	w.UpdateNativeSize()
	texture := w.object.ToTexture()
	sprite := graphics.SfSprite_create()
	defer graphics.SfSprite_destroy(sprite)
	graphics.SfSprite_setTexture(sprite, texture.Texture, 1)
	graphics.SfRenderWindow_clear(w.renderWindow, graphics.GetSfTransparent())
	graphics.SfRenderWindow_drawSprite(w.renderWindow, sprite, (graphics.SfRenderStates)(graphics.SwigcptrSfRenderStates(0)))
}

// display a window without polling for event, used to display derived windows
// maybe we can still poll them but it'll need to be in a different thread or something
func (w *Window) DisplayNoPoll() {
	w.Render()
	graphics.SfRenderWindow_display(w.renderWindow)
}

func (w *Window) Display() {
	w.Render()
	graphics.SfRenderWindow_display(w.renderWindow)
	event := window.NewSfEvent()
	defer window.DeleteSfEvent(event)
	for graphics.SfRenderWindow_waitEvent(w.renderWindow, event) > 0 {
		switch event.GetEvType() {
		case window.SfEventType(window.SfEvtClosed):
			return
		case window.SfEventType(window.SfEvtResized):
			w.UpdateWidthOrHeight()
		case window.SfEventType(window.SfEvtMouseButtonPressed):
			mb := event.GetMouseButton()
			wWidth, wHeight := w.GetSize()
			nWidth, nHeight := w.GetNativeSize()

			x, y := MapCoordinateIfInside(mb.GetX(), mb.GetY(), 0, 0, wWidth, wHeight, nWidth, nHeight)
			buttonDownEvent := MouseButtonDownEvent{
				X: x,
				Y: y,
			}
			switch int(mb.GetButton()) {
			case window.SfMouseLeft:
				buttonDownEvent.Button = MouseButtonLeft
			case window.SfMouseRight:
				buttonDownEvent.Button = MouseButtonRight
			case window.SfMouseMiddle:
				buttonDownEvent.Button = MouseButtonMiddle
			default:
				continue
			}
			if !HandleEvent(w, w.object, buttonDownEvent) {
				continue
			}
		case window.SfEventType(window.SfEvtTextEntered):
			textEvent := TextEvent{
				Rune: rune(event.GetText().GetUnicode()),
			}
			if !HandleEvent(w, w.object, textEvent) {
				fmt.Println("skipped text event: ", textEvent)
				if textEvent.Rune == 3 {
					w.SaveToImage("screenshot.png")
				}
				continue
			}
		case window.SfEventType(window.SfEvtKeyPressed):
			if (event.GetKey().GetControl() != 0) && (event.GetKey().GetCode() == window.SfKeyCode(window.SfKeyV)) {
				pasteEvent := PasteEvent{
					Clipboard: StringFromUTF32(window.SfClipboard_getUnicodeString()),
				}

				if !HandleEvent(w, w.object, pasteEvent) {
					fmt.Println("skipped paste event: ", pasteEvent)
					continue
				}
			} else if event.GetKey().GetCode() == window.SfKeyCode(window.SfKeyEnter) {
				if !HandleEvent(w, w.object, EnterEvent{}) {
					fmt.Println("skipped enter event")
					continue
				}
			} else {
				keyEvent, isAllowed := allowedKeyEvent[event.GetKey().GetCode()]
				if isAllowed {
					if !HandleEvent(w, w.object, keyEvent) {
						fmt.Println("skipped key event: ", keyEvent)
						continue
					}
				}
			}
		default:
			// fmt.Println("unhandled event type")
			continue
		}
		// if something changed then we redraw
		w.Render()
		graphics.SfRenderWindow_display(w.renderWindow)
	}
}

func (w *Window) SaveToImage(path string) {
	w.object.ToTexture().SaveToImage(path)
}

func (w *Window) ChangeFocusObject(o Focusable) {
	if o != w.focusObject {
		if w.focusObject != nil {
			w.focusObject.UnsetFocus()
		}
		w.focusObject = o
		o.SetFocus()
	} else {
		fmt.Println("already focused!")
	}
}
