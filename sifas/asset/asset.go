package asset

import (
	_ "modernc.org/sqlite"

	"elichika/assetdata"
	"elichika/gui/locale"
	"elichika/utils"
	hwdecrypt "github.com/arina999999997/gohwdecrypt"

	"fmt"
	"os"

	"xorm.io/xorm"
)

var cdnMasterVersionMapping = map[string]string{}

func init() {
	cdnMasterVersionMapping["en"] = "2d61e7b4e89961c7"
	cdnMasterVersionMapping["ko"] = "2d61e7b4e89961c7"
	cdnMasterVersionMapping["zh"] = "2d61e7b4e89961c7"

	// TODO(cdn): make this change globally
	cdnMasterVersionMapping["en"] = "b66ec2295e9a00aa"
	cdnMasterVersionMapping["ko"] = "b66ec2295e9a00aa"
	cdnMasterVersionMapping["zh"] = "b66ec2295e9a00aa"
	cdnMasterVersionMapping["ja"] = "b66ec2295e9a00aa"
}

type Asset struct {
	AssetPath string `xorm:"asset_path"`
	PackName  string `xorm:"pack_name"`
	Head      int    `xorm:"head"`
	Size      int    `xorm:"size"`
	Key1      uint32 `xorm:"key1"`
	Key2      uint32 `xorm:"key2"`
}

func (a Asset) LoadUnencrypted() []byte {
	data := assetdata.GetDownloadData(a.PackName)
	actualFile := fmt.Sprintf("static/%s/%s", cdnMasterVersionMapping[data.Locale], data.File)
	fmt.Println(actualFile)
	actualStart := data.Start + a.Head
	file, err := os.Open(actualFile)
	utils.CheckErr(err)
	defer file.Close()
	output := make([]byte, a.Size)
	_, err = file.ReadAt(output, int64(actualStart))
	utils.CheckErr(err)
	// keying
	hwdecrypt.DecryptBuffer(&hwdecrypt.HwdKeyset{
		Key1: a.Key1,
		Key2: a.Key2,
		Key3: 12345,
	}, output)
	return output
}

func loadAssets() {
	// assume the db is at assets/db/gl/asset_<Platform>_<Locale>.db
	// or assets/db/jp/asset_<Platform>_<Locale>.db
	engine, err := xorm.NewEngine("sqlite", fmt.Sprintf("assets/db/%s/asset_%s_%s.db", locale.AppVersion(), locale.Platform, locale.Language))
	utils.CheckErr(err)
	engine.SetMaxOpenConns(50)
	engine.SetMaxIdleConns(10)
	assetdata.Init(locale.Language, engine)
	session := engine.NewSession()
	defer session.Close()
	assets := []Asset{}
	err = session.Table("texture").Find(&assets)
	utils.CheckErr(err)

	for _, asset := range assets {
		AssetMap[locale.Locale()][asset.AssetPath] = asset
	}
}
