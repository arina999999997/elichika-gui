package graphic
import (
	"github.com/telroshan/go-sfml/v2/window"
)
// const (
// 	EventTypeClick      int = 0
// 	EventTypeRightClick int = 1
// 	EventTypeScrollUp   int = 2
// 	EventTypeScrollDown int = 3
// )

type Event interface {
}

const (
	MouseButtonLeft   int = 0
	MouseButtonRight  int = 1
	MouseButtonMiddle int = 2
)

type MouseButtonDownEvent struct {
	Button int
	X      int
	Y      int
}

type TextEvent struct {
	Rune rune
}

type PasteEvent struct {
	Clipboard string
}

type EnterEvent struct {
}

type KeyEvent = int

var allowedKeyEvent = map[window.SfKeyCode]KeyEvent{}

const (
	KeyEventUp KeyEvent = 0
	KeyEventDown KeyEvent = 1
	KeyEventLeft KeyEvent = 2
	KeyEventRight KeyEvent = 3
)

func init() {
	allowedKeyEvent[window.SfKeyCode(window.SfKeyUp)] = KeyEventUp
	allowedKeyEvent[window.SfKeyCode(window.SfKeyDown)] = KeyEventDown
	allowedKeyEvent[window.SfKeyCode(window.SfKeyLeft)] = KeyEventLeft
	allowedKeyEvent[window.SfKeyCode(window.SfKeyRight)] = KeyEventRight
}
// generic handle event function, return true if the event is handled
// the event is always mapped to the native resolution of the object
//
// if a type implement the composite object interface, then its events will be handled using the following rule:
// - Invoke ForEach
// - Call MapEvent to map the event to child object, if necessary
func HandleEvent(w *Window, o Object, e Event) bool {
	_, isPollingComposite := o.(PollingCompositeObject)
	if isPollingComposite {
		result := false
		pco := o.(PollingCompositeObject)
		pco.ForEach(func(child Object) {
			if result {
				return
			}
			result = HandleEvent(w, child, pco.MapEvent(e, child))
		})
		if result {
			return true
		}
	}
	switch e.(type) {
		case MouseButtonDownEvent:
			_, isClickable := o.(Clickable)
			if isClickable {
				return o.(Clickable).OnClick(w, e.(MouseButtonDownEvent))
			}
		case TextEvent:
			_, isInputable := o.(Inputable)
			if isInputable {
				if o.(Inputable).HasFocus() {
					return o.(Inputable).OnText(e.(TextEvent))
				}
			}
		case PasteEvent:
			_, isInputable := o.(Inputable)
			if isInputable {
				if o.(Inputable).HasFocus() {
					return o.(Inputable).OnPaste(e.(PasteEvent))
				}
			}
		case EnterEvent:
			_, isEnterable := o.(Enterable)
			if isEnterable {
				if o.(Enterable).HasFocus() {
					return o.(Enterable).OnEnter()
				}
			}
		case KeyEvent:
			_, isKeyable := o.(Keyable)
			if isKeyable {
				if o.(Keyable).HasFocus() {
					return o.(Keyable).OnKey(e.(KeyEvent))
				}
			}
		default:
			panic("Unsupported event type")
	}
	return false
}
