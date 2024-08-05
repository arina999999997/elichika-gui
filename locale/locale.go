package locale

// This package handle the shared locale in all the GUI apps

var Languages = []string{"ja", "en", "ko", "zh"} // Japanese, English, Korean, Chinese
var Platforms = []string{"a", "i"}               // Android / IOS

var Language = "en"
var Platform = "a"

func SetLanguage(lang string) {
	for _, accept := range Languages {
		if accept == lang {
			Language = accept
			return
		}
	}
	panic("Unknown language: " + lang)
}

func SetPlatform(plat string) {
	for _, accept := range Languages {
		if accept == plat {
			Platform = plat
			return
		}
	}
	panic("Unknown platform: " + plat)
}

// return an unique key to this locale
func Locale() string {
	return Language + "_" + Platform
}

func AppVersion() string {
	if Language == "ja" {
		return "jp"
	} else {
		return "gl"
	}
}
