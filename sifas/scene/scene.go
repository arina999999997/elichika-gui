package scene

// Scene refer to a specific scene in the game
// for example, the training menu is a scene, the event menu is a scene
// even the homescreen is just the home scene

// in our term, scenes are TEXTURELESS composition of textures.
// application can then set the relevant texture and render the scene, as if they were in the game (ideally)

// for now, the scenes are only for development purpose, and are all modeled on a 1800x900 canvas as
// that is the native resolution for background textures (cards and events)
// ideally the scenes would follow the game's way of placing items for different resolution size, but that's a project for another time.

const (
	GameWidth  = 1800
	GameHeight = 900
)
