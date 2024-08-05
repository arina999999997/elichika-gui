package asset

import (
	"elichika/gui/locale"
)

var AssetMap = map[string]map[string]Asset{}

func loadLocale() {
	if AssetMap[locale.Locale()] == nil {
		AssetMap[locale.Locale()] = map[string]Asset{}
		loadAssets()
	}
}
