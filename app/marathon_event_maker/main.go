package main

import (
	"elichika/gui/asset_manager"
	"elichika/gui/graphic"
	"elichika/gui/sifas/scene"
	// "time"
)

// quick and dirty reward list display that show items from a list

func main() {

	window, err := graphic.NewWindow("Marathon event maker")
	if err != nil {
		panic(err)
	}

	// but it could be other things too
	eventTop := scene.GetEventMarathonTopScene()
	eventTop.Background.LoadTexture(asset_manager.SafeLoadTexture("Ie6"))
	eventTop.Title.LoadTexture(asset_manager.SafeLoadTexture(".#+"))
	// eventTop.Title.LoadTexture(sifas.SafeLoadTexture("Ho5"))
	eventTop.Board.LoadTexture(asset_manager.SafeLoadTexture("r0|"))
	window.RenderRootObject(eventTop.Root)

	/* Create the main window */
	// cs := window.NewSfContextSettings()
	// defer window.DeleteSfContextSettings(cs)
	// w := graphics.SfRenderWindow_create(vm, "SFML window", uint(window.SfResize|window.SfClose), cs)
	// defer window.SfWindow_destroy(w)

	// ev := window.NewSfEvent()
	// defer window.DeleteSfEvent(ev)

	/* Start the game loop */
	// for window.SfWindow_isOpen(w) > 0 {
	// 	/* Process events */
	// 	for window.SfWindow_pollEvent(w, ev) > 0 {
	// 		/* Close window: exit */
	// 		if ev.GetEvType() == window.SfEventType(window.SfEvtClosed) {
	// 			return
	// 		}
	// 	}
	// 	graphics.SfRenderWindow_clear(w, graphics.GetSfRed())
	// 	graphics.SfRenderWindow_display(w)
	// }
}
