package event_marathon

import (
	"elichika/client"
	"elichika/enum"

	"elichika/gui/graphic"
	"elichika/gui/sifas/asset"

	"sort"

	"github.com/telroshan/go-sfml/v2/graphics"
)

const (
	// The textures are all rendered on the board natively
	// howerver, the whole things then get scaled for some reason

	// actual board width and height
	BoardWidth  = 776
	BoardHeight = 1078

	BoardDecoWidth  = 708
	BoardDecoHeight = 512
	BoardDecoX      = 34
	BoardDecoY      = 33

	BoardMemoWidth  = 200
	BoardMemoHeight = 200
)

// Memo is prerotated, in fact it can be other shapes too
var (
	BoardMemoX = []int{0, 90, 327, 576, 13, 273, 525, 67, 323, 552}
	BoardMemoY = []int{0, 2, 2, 2, 200, 187, 189, 386, 387, 380}

	BoardPictureX = []int{0, 54, 324, 482, 145, 474, 44, 308}
	BoardPictureY = []int{0, 53, 48, 87, 217, 254, 370, 358}
	// rotation is measured in degree, the angle is counter clockwise from horizontal
	BoardPictureRotation = []float32{0, 4, -2, 4, -2, 4, 4, 4}
)

type EventMarathonBoardThing struct {
	PositionType int32 `enum:"EventMarathonBoardPositionType"`
	Position     int32
	Priority     int32

	Texture *graphic.Texture
}

type EventMarathonBoard struct {
	BoardBaseImage *graphic.Texture
	BoardDecoImage *graphic.Texture

	BoardThings []EventMarathonBoardThing

	// gui backend
	Canvas *graphic.Canvas
}

func (emb *EventMarathonBoard) Load(board client.EventMarathonBoard) error {
	var err error
	if board.BoardBaseImagePath.V.HasValue {
		emb.BoardBaseImage, err = asset.LoadTexture(board.BoardBaseImagePath.V.Value)
		if err != nil {
			return err
		}
	} else {
		emb.BoardBaseImage = nil
	}
	if board.BoardDecoImagePath.V.HasValue {
		emb.BoardDecoImage, err = asset.LoadTexture(board.BoardDecoImagePath.V.Value)
		if err != nil {
			return err
		}
	} else {
		emb.BoardDecoImage = nil
	}

	emb.BoardThings = nil
	for _, thing := range board.BoardThingMasterRows.Slice {
		texture := graphic.DefaultTexture()
		if thing.ImageThumbnailAssetPath.V.HasValue {
			texture, err = asset.LoadTexture(thing.ImageThumbnailAssetPath.V.Value)
			if err != nil {
				return err
			}
		}
		emb.BoardThings = append(emb.BoardThings, EventMarathonBoardThing{
			PositionType: thing.EventMarathonBoardPositionType,
			Position:     thing.Position,
			Texture:      texture,
			Priority:     thing.Priority,
		})
	}
	if emb.BoardThings != nil {
		sort.Slice(emb.BoardThings, func(i, j int) bool {
			return emb.BoardThings[i].Priority < emb.BoardThings[j].Priority
		})
	}
	graphic.InvalidateRenderCache(emb)
	return nil
}

// graphic.Object

func (emb *EventMarathonBoard) GetWidth() int {
	return BoardWidth
}

func (emb *EventMarathonBoard) GetHeight() int {
	return BoardHeight
}

func (emb *EventMarathonBoard) InvalidateRenderCache() bool {
	return emb.Canvas.InvalidateRenderCache()
}

func (emb *EventMarathonBoard) Draw() {
	if emb.Canvas.IsRendered() {
		return
	}
	if emb.Canvas == nil {
		emb.Canvas = graphic.NewCanvas(emb)
	}

	if emb.BoardBaseImage != nil {
		emb.Canvas.DrawTexture(emb.BoardBaseImage, 0, 0, BoardWidth, BoardHeight)
	}
	if emb.BoardDecoImage != nil {
		emb.Canvas.DrawTexture(emb.BoardDecoImage, BoardDecoX, BoardDecoY, BoardDecoWidth, BoardDecoHeight)
	}

	for _, thing := range emb.BoardThings {
		if thing.PositionType == enum.EventMarathonBoardPositionTypeMemo {
			emb.Canvas.DrawTexture(thing.Texture, BoardMemoX[thing.Position], BoardMemoY[thing.Position], BoardMemoWidth, BoardMemoHeight)
		} else {
			sprite := graphics.SfSprite_create()
			defer graphics.SfSprite_destroy(sprite)
			graphics.SfSprite_setTexture(sprite, thing.Texture.Texture, 1)
			graphics.SfSprite_setPosition(sprite, graphic.GetVector2f(BoardPictureX[thing.Position], BoardPictureY[thing.Position]))
			graphics.SfSprite_setRotation(sprite, BoardPictureRotation[thing.Position])
			emb.Canvas.DrawSprite(&graphic.Sprite{
				Sprite: sprite,
			})
		}
	}

	emb.Canvas.Finalize()
}

func (emb *EventMarathonBoard) ToTexture() *graphic.Texture {
	emb.Draw()
	return emb.Canvas.AsTexture()
}
