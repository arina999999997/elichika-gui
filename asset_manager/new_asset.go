package asset_manager

import (
	_ "modernc.org/sqlite"

	"elichika/gui/locale"
	"elichika/utils"

	"errors"
	"fmt"
	"os"

	"xorm.io/xorm"
)

type NewAsset struct {
	AssetPath string `xorm:"pk 'asset_path'"`
	Language  string `xorm:"pk 'language'"`
	Platform  string `xorm:"pk 'platform'"`
	FilePath  string `xorm:"'file_path'"`
}

var Engine *xorm.Engine

func init() {
	var err error
	// TODO(hardcoded)
	Engine, err = xorm.NewEngine("sqlite", "new_asset.db")
	utils.CheckErr(err)
	Engine.SetMaxOpenConns(1)
	Engine.SetMaxIdleConns(10)

	exist, err := Engine.Table("texture").IsTableExist("texture")
	utils.CheckErr(err)

	if !exist {
		fmt.Println("Creating new table: texture")
		err = Engine.Table("texture").CreateTable(&NewAsset{})
		utils.CheckErr(err)
	}
}

// load asset from the new asset database
// TODO(hardcoded)
func loadNewAsset(assetPath string) (bytes []byte, err error) {
	filePath := ""
	var exist bool
	Engine.ShowSQL(true)
	exist, err = Engine.Table("texture").Where("asset_path = ? AND language = ? AND platform = ?",
		assetPath, locale.Language, locale.Platform).Cols("file_path").Get(&filePath)
	if err != nil {
		return
	}
	if !exist {
		err = errors.New("Asset doesn't exist: " + assetPath + ", Language: " + locale.Language + ", Platform: " + locale.Platform)
		return
	}

	stat, err := os.Stat(filePath)
	if err != nil {
		return
	}
	size := stat.Size()
	bytes = make([]byte, size)

	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		return
	}
	_, err = file.ReadAt(bytes, 0)
	return
}
