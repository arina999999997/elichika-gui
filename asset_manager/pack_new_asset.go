package asset_manager

import (
	"elichika/assetdata"
	"elichika/gui/locale"
	"elichika/gui/sifas/asset"
	"elichika/utils"
	hwdecrypt "github.com/arina999999997/gohwdecrypt"

	"fmt"
	"math/rand"
	"os"
	"sort"
	"strings"

	"xorm.io/xorm"
)

// generate the following files:
// - static/new/<random name>:
//   - these are the files client will download and use
//
// - asset_<platform>_<language>.sql:
//   - these are the files necessary
//
// this doesn't affect the new asset database, so a manual clean up is necessary
var usedNewName = map[string]bool{}

func getUniquePackName() string {
	check := func(name string) bool {
		_, exist := assetdata.Pack[name]
		if exist {
			return false
		}
		_, exist = assetdata.Metapack[name]
		if exist {
			return false
		}
		return !usedNewName[name]
	}
	randomName := func() string {
		runes := []rune("abcdefghijklmnopqrstuvwxyz01234566789")
		// 36^6 = 2176782336 different names
		name := [6]rune{}
		for i := 0; i < 6; i++ {
			name[i] = runes[rand.Intn(36)]
		}
		return string(name[:])
	}
	for {
		name := randomName()
		if check(name) {
			usedNewName[name] = true
			return name
		}
	}
}

func PackNewAssets() {
	// initialise the database
	for _, language := range locale.Languages {
		appVersion := "gl"
		if language == "ja" {
			appVersion = "jp"
		}
		for _, platform := range locale.Platforms {
			engine, err := xorm.NewEngine("sqlite", fmt.Sprintf("assets/db/%s/asset_%s_%s.db", appVersion, platform, language))
			utils.CheckErr(err)
			assetdata.Init(language, engine)
		}
	}

	// for now, assume the package key is main
	packageKey := "main"

	newAssets := []NewAsset{}
	session := Engine.NewSession()
	defer session.Close()
	session.Begin()
	err := session.Table("texture").Find(&newAssets)
	utils.CheckErr(err)

	fileToVersions := map[string][]string{}
	for _, asset := range newAssets {
		fileToVersions[asset.FilePath] = append(fileToVersions[asset.FilePath], asset.Language+","+asset.Platform)
	}

	versionToFiles := map[string][]string{}
	for filePath, fileVersions := range fileToVersions {
		sort.Slice(fileVersions, func(i, j int) bool {
			return fileVersions[i] < fileVersions[j]
		})
		versionString := ""
		for i, version := range fileVersions {
			if i == 0 {
				versionString += version
			} else if fileVersions[i-1] != version {
				versionString += "|" + version
			}
		}
		versionToFiles[versionString] = append(versionToFiles[versionString], filePath)
	}

	fileToAsset := map[string]asset.Asset{}
	versionSql := map[string][]string{}
	// pack all the files with the same versions together
	os.MkdirAll("static/new", 0777)
	for versionString, files := range versionToFiles {
		packName := getUniquePackName()
		pack, err := os.Create("static/new/" + packName)
		defer pack.Close()
		utils.CheckErr(err)
		for _, file := range files {
			key1 := rand.Uint32()
			key2 := rand.Uint32()
			bytes, err := LoadFileOnDisk(file)
			utils.CheckErr(err)
			hwdecrypt.DecryptBuffer(&hwdecrypt.HwdKeyset{
				Key1: key1,
				Key2: key2,
				Key3: 12345,
			}, bytes)
			head, err := pack.Seek(0, os.SEEK_CUR)
			utils.CheckErr(err)
			size, err := pack.Write(bytes)
			utils.CheckErr(err) // err != nil if size != actual size
			fileToAsset[file] = asset.Asset{
				PackName: packName,
				Head:     int(head),
				Size:     size,
				Key1:     key1,
				Key2:     key2,
			}
		}
		packSize, err := pack.Seek(0, os.SEEK_CUR)
		utils.CheckErr(err)

		// insert the pack itself
		versions := strings.Split(versionString, "|")
		for _, version := range versions {
			sql := fmt.Sprintf(`INSERT INTO m_asset_package_mapping VALUES ("%s", "%s", "%d", NULL, 0, 8)`, packageKey, packName, packSize)
			versionSql[version] = append(versionSql[version], sql)
			sql = fmt.Sprintf(`INSERT INTO m_asset_pack VALUES ("%s", 0)`, packageKey)
			versionSql[version] = append(versionSql[version], sql)
		}
	}

	// update the count for each package key
	// update the version number

	randomHash := func() string {
		runes := []rune("0123456789abcdef")
		// 36^6 = 2176782336 different names
		name := [40]rune{}
		for i := 0; i < 6; i++ {
			name[i] = runes[rand.Intn(16)]
		}
		return string(name[:])
	}

	// this should be per package_key, but we only need 1 because we assume the package key will me main
	versionHash := randomHash()

	for version := range versionSql {
		sql := fmt.Sprintf(`UPDATE m_asset_package SET version = "%s", pack_num = (SELECT COUNT(*) FROM m_asset_package_mapping WHERE package_key="%s") WHERE package_key = "%s"`,
			versionHash, packageKey, packageKey)
		versionSql[version] = append(versionSql[version], sql)
	}

	// insert the asset path mapping for the packs
	for _, asset := range newAssets {
		assetObject := fileToAsset[asset.FilePath]
		sql := fmt.Sprintf(`INSERT INTO texture VALUES ("%s", "%s", %d, %d, %d, %d)`,
			sqlEscape(asset.AssetPath), assetObject.PackName, assetObject.Head, assetObject.Size, assetObject.Key1, assetObject.Key2)

		versionSql[asset.Language+","+asset.Platform] = append(versionSql[asset.Language+","+asset.Platform], sql)
	}

	// generate the sql files
	for version, lines := range versionSql {
		tokens := strings.Split(version, ",")
		language := tokens[0]
		platform := tokens[1]
		file := fmt.Sprintf("asset_%s_%s.db.sql", platform, language)
		f, err := os.Create(file)
		defer f.Close()
		utils.CheckErr(err)
		for _, line := range lines {
			_, err := f.WriteString(line + "\n")
			utils.CheckErr(err)
		}
	}

}
