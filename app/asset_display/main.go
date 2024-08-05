package main

import (
	"elichika/gui/graphic"
	"elichika/gui/sifas/asset"

	"fmt"
	"os"
	// "time"
)

// this is mostly a testing program
// take args being the asset paths of texture
// clicking will cycle the assets
func main() {
	window, err := graphic.NewWindow("Texture viewer")
	if err != nil {
		panic(err)
	}
	objects := []*graphic.Object{}
	for i := 1; i < len(os.Args); i++ {
		object, err := asset.LoadTextureToObject(os.Args[i])
		if err != nil {
			fmt.Println("Failed to load asset: ", os.Args[i], ", error: ", err)
		} else {
			objects = append(objects, object)
		}
	}
	if len(objects) == 0 {
		panic("No textures to show")
	}
	i := 0
	for {
		if i == len(objects) {
			i = 0
		}
		fmt.Println(i)
		window.RenderRootObject(objects[i])
		i++
		// time.Sleep(time.Second * 10)
	}
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
