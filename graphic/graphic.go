package graphic

// This package is an abstraction layer of various graphic concepts, just so we can switch to a different graphic library if necessary
// We follow the conventions:
// - Coordinate x refer to left - right, the higher x is, the more to right a point is.
// - Coordinate y refer to up - down, the higher y is, the more to the bottom a point is.
// - When refering to size, the width (x wise) come before the height (y wise).
// - Things that need to be drawn are refered to as objects.
// - Each object is treated as a rectangle, and has an internal resolution of HEIGHT * WIDTH.
// - Each object follow the (graphic) Object interface:
//   - Has method for drawing on a Canvas at a specific position.
//   -
// - Each object can also follow one of the interfaces for polling events:
//   - For example, object that follow the Clickable interface will response to click event.
// 	 - The event handling is intended to make each event trigger the handling function for only one object.
// - Generally, there's nothing stopping each object from implementing a complicated drawing function.
// - However, we try to follow the tree model:
//   - Each objects can either be its own drawable thing
//   - Or it can contain other objects as a part of it
//   - When drawing, draw the current "object" as a background, then recursively draw the child objects, so the child objects should be on top of the parent objects.
//   - When handling event, try to handle event using the child objects first (recursively), then try to handle the event with the current object if it's still not handled.

// TODO(gui): Maybe implement a prerender step that traverse through the tree top-down, allowing us to set things up before the bottom-up drawing process.
import (
	"runtime"
)

func init() { runtime.LockOSThread() }

// A lot of freedom, it's up to the object to figure out how and where it need to do things to

type Object interface {
	// get the native width and height
	GetWidth() int
	GetHeight() int

	// invalidate the cache, return true if NEWLY invalidated
	InvalidateRenderCache() bool

	// draw to its own canvas
	Draw()

	// create a Texture that is the whole object
	ToTexture() *Texture
}

type ChildObject interface {
	Object
	GetParent() Object
}

type Clickable interface {
	Object

	// Return true if the click event is handled
	OnClick(*Window, MouseButtonDownEvent) bool
}

type Focusable interface {
	Clickable

	SetFocus()
	UnsetFocus()
	HasFocus() bool
}

type Inputable interface {
	Focusable
	OnText(TextEvent) bool
	OnPaste(PasteEvent) bool
}

type Enterable interface {
	Inputable
	OnEnter() bool
}

type Keyable interface {
	Focusable
	OnKey(KeyEvent) bool
}

type CompositeObject interface {
	ForEach(func(Object))
}

type PollingCompositeObject interface {
	CompositeObject
	MapEvent(e Event, o Object) Event
}

// type Traceable interface {
// 	Object
// 	OnMouseMove()
// }

// type Scrollable interface {
// 	Object
// 	OnScrollUp()
// 	OnScrollDown()
// }

// TODO(memory): This thing have lots of memory leak in it, so don't use it for anything that can result in leaks beyond controls
