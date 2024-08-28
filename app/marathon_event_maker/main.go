package main

import (

	"elichika/gui/app/marathon_event_maker/memo_maker"
	
	"elichika/gui/graphic"
	// "elichika/gui/sifas/locale"
	"elichika/gui/sifas/scene/event_marathon"
	// "time"

	"elichika/client"
	"elichika/enum"
	"elichika/generic"

	"fmt"
)

// quick and dirty reward list display that show items from a list

func main() {

	window, err := graphic.NewWindow("Marathon event maker")
	if err != nil {
		panic(err)
	}
	topWindow, err := graphic.NewWindow("Event Marathon Top Scene")
	if err != nil {
		panic(err)
	}

	eventMarathonTopScene := &event_marathon.EventMarathonTopScene{}
	topStatus := client.EventMarathonTopStatus{}
	topStatus.BackgroundImagePath.V = generic.NewNullable(":Y")
	topStatus.TitleImagePath.V = generic.NewNullable("?0=")
	topStatus.BoardStatus.BoardBaseImagePath.V = generic.NewNullable("+6I")
	topStatus.BoardStatus.BoardDecoImagePath.V = generic.NewNullable("br{")
	for i := int32(1); i <= 9; i++ {
		topStatus.BoardStatus.BoardThingMasterRows.Append(client.EventMarathonBoardMemorialThingsMasterRow{
			EventMarathonBoardPositionType: enum.EventMarathonBoardPositionTypeMemo,
			Position:                       i,
			ImageThumbnailAssetPath: client.TextureStruktur{
				V: generic.NewNullable("q*O"),
			},
		})
	}
	// topStatus.BoardStatus.BoardThingMasterRows.Append(client.EventMarathonBoardMemorialThingsMasterRow{
	// 	EventMarathonBoardPositionType: enum.EventMarathonBoardPositionTypePicture,
	// 	Position: 1,
	// 	ImageThumbnailAssetPath: client.TextureStruktur{
	// 		V: generic.NewNullable("P44"),
	// 	},
	// 	Priority: 1,
	// })
	// topStatus.BoardStatus.BoardThingMasterRows.Append(client.EventMarathonBoardMemorialThingsMasterRow{
	// 	EventMarathonBoardPositionType: enum.EventMarathonBoardPositionTypePicture,
	// 	Position: 2,
	// 	ImageThumbnailAssetPath: client.TextureStruktur{
	// 		V: generic.NewNullable("M32"),
	// 	},
	// 	Priority: 2,
	// })
	// topStatus.BoardStatus.BoardThingMasterRows.Append(client.EventMarathonBoardMemorialThingsMasterRow{
	// 	EventMarathonBoardPositionType: enum.EventMarathonBoardPositionTypePicture,
	// 	Position: 3,
	// 	ImageThumbnailAssetPath: client.TextureStruktur{
	// 		V: generic.NewNullable("*.x"),
	// 	},
	// 	Priority: 3,
	// })
	// topStatus.BoardStatus.BoardThingMasterRows.Append(client.EventMarathonBoardMemorialThingsMasterRow{
	// 	EventMarathonBoardPositionType: enum.EventMarathonBoardPositionTypePicture,
	// 	Position: 4,
	// 	ImageThumbnailAssetPath: client.TextureStruktur{
	// 		V: generic.NewNullable("8R$"),
	// 	},
	// 	Priority: 4,
	// })

	// topStatus.BoardStatus.BoardThingMasterRows.Append(client.EventMarathonBoardMemorialThingsMasterRow{
	// 	EventMarathonBoardPositionType: enum.EventMarathonBoardPositionTypePicture,
	// 	Position: 5,
	// 	ImageThumbnailAssetPath: client.TextureStruktur{
	// 		V: generic.NewNullable(":\\L"),
	// 	},
	// 	Priority: 7,
	// })
	// topStatus.BoardStatus.BoardThingMasterRows.Append(client.EventMarathonBoardMemorialThingsMasterRow{
	// 	EventMarathonBoardPositionType: enum.EventMarathonBoardPositionTypePicture,
	// 	Position: 6,
	// 	ImageThumbnailAssetPath: client.TextureStruktur{
	// 		V: generic.NewNullable("r}}"),
	// 	},
	// 	Priority: 5,
	// })
	// topStatus.BoardStatus.BoardThingMasterRows.Append(client.EventMarathonBoardMemorialThingsMasterRow{
	// 	EventMarathonBoardPositionType: enum.EventMarathonBoardPositionTypePicture,
	// 	Position: 7,
	// 	ImageThumbnailAssetPath: client.TextureStruktur{
	// 		V: generic.NewNullable("iWh"),
	// 	},
	// 	Priority: 6,
	// })
	err = eventMarathonTopScene.Load(topStatus)
	if err != nil {
		fmt.Println("error: ", err)
	}
	topWindow.SetObject(eventMarathonTopScene)
	topWindow.DisplayNoPoll()

	// selector := locale.GetLocaleSelector(nil, func() {
	// 	err = eventMarathonTopScene.Load(topStatus)
	// 	if err != nil {
	// 		fmt.Println("error: ", err)
	// 	}
	// 	topWindow.DisplayNoPoll()
	// })

	// window.SetObject(selector)
	backgroundPicker := memo_maker.NewMemoMaker(nil)
	window.SetObject(backgroundPicker)
	window.Display()
}
