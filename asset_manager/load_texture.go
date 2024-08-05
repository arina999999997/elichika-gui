package asset_manager

import (
	"elichika/gui/graphic"
	"elichika/gui/sifas/asset"

	"errors"
	"fmt"
)

// Try to load a texture from the game's database
// If it is not present, then load from the new asset database instead
func LoadTexture(assetPath string) (texture *graphic.Texture, err error) {
	texture, err = asset.LoadTexture(assetPath)
	if err == nil {
		return
	} else {
		fmt.Println(err)
	}
	texture, err = LoadNewTexture(assetPath)
	if err == nil {
		return
	} else {
		fmt.Println(err)
	}
	return nil, errors.New("Asset doesn't exist in either game or new database")
}

func LoadNewTexture(assetPath string) (texture *graphic.Texture, err error) {
	rawData, err := loadNewAsset(assetPath)
	if err != nil {
		return
	}
	texture = &graphic.Texture{}
	texture.LoadFromMemory(rawData)
	return
}

// Loading a texture by assets path, but if it is not present, then just return a missing texture
func SafeLoadTexture(assetPath string) *graphic.Texture {
	texture, err := LoadTexture(assetPath)
	if err == nil {
		return texture
	} else {
		fmt.Println(err, "\nUsing default texture")
		return graphic.DefaultTexture()
	}
}
