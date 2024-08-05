package scene

import (
	"elichika/gui/graphic"
)

const (
	EventMarathonTopBackgroundX      = 0
	EventMarathonTopBackgroundY      = 0
	EventMarathonTopBackgroundWidth  = GameWidth
	EventMarathonTopBackgroundHeight = GameHeight
	EventMarathonTopBoardX           = 120
	EventMarathonTopBoardY           = 97
	EventMarathonTopBoardWidth       = 931
	EventMarathonTopBoardHeight      = 1293 // this extend out of the thingy
	EventMarathonTopTitleX           = 1150
	EventMarathonTopTitleY           = 133
	EventMarathonTopTitleWidth       = 460
	EventMarathonTopTitleHeight      = 163
)

// TODO(extra): these naming might differ from the name in code
// the objects are all exposed
type EventMarathonTopScene struct {
	Root       *graphic.Object
	Background *graphic.Object
	Board      *graphic.Object // The board is kinda its own scene in and of itself
	Title      *graphic.Object
}

func GetEventMarathonTopScene() EventMarathonTopScene {
	s := EventMarathonTopScene{}
	s.Root = &graphic.Object{
		Width:  GameWidth,
		Height: GameHeight,
	}

	s.Background = &graphic.Object{}
	s.Root.AnchorChild(s.Background, EventMarathonTopBackgroundX, EventMarathonTopBackgroundY,
		EventMarathonTopBackgroundWidth, EventMarathonTopBackgroundHeight)

	s.Board = &graphic.Object{}
	s.Root.AnchorChild(s.Board, EventMarathonTopBoardX, EventMarathonTopBoardY,
		EventMarathonTopBoardWidth, EventMarathonTopBoardHeight)

	s.Title = &graphic.Object{}
	s.Root.AnchorChild(s.Title, EventMarathonTopTitleX, EventMarathonTopTitleY,
		EventMarathonTopTitleWidth, EventMarathonTopTitleHeight)

	return s
}
