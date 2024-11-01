package main

import (
	"elichika/gui/graphic"

	hwdecrypt "github.com/arina999999997/gohwdecrypt"

	"bytes"
	"errors"
	"fmt"
	"image/jpeg"
	"image/png"
	"os"
)

func GetImageSize(data []byte) (int32, int32, error) {

	image, err := png.Decode(bytes.NewReader(data))

	if err != nil {
		image, err = jpeg.Decode(bytes.NewReader(data))
	}

	if err != nil {
		return 0, 0, err
	}
	bound := image.Bounds()
	w := int32(bound.Max.X - bound.Min.X)
	h := int32(bound.Max.Y - bound.Min.Y)
	return w, h, nil
}

type RawTexture struct {
	PackName string  `xorm:"'pack_name'"`
	Head     int32   `xorm:"'head'"`
	Size     int32   `xorm:"'size'"`
	Key1     int32   `xorm:"'key1'"`
	Key2     int32   `xorm:"'key2'"`
	Width    int32   `xorm:"'width'"`
	Height   int32   `xorm:"'height'"`
	Error    *string `xorm:"'error'"`
}

type Texture struct {
	AssetPath string `xorm:"pk 'asset_path'"`
	PackName  string `xorm:"'pack_name'"`
	Head      int32  `xorm:"'head'"`
	Size      int32  `xorm:"'size'"`
	Key1      int32  `xorm:"'key1'"`
	Key2      int32  `xorm:"'key2'"`
}

type DetailedTexture struct {
	Language  string  `xorm:"'language'"` // the language of the asset
	Version   string  `xorm:"'version'"`  // the first version where this texture was first inserted
	AssetPath string  `xorm:"'asset_path'"`
	PackName  string  `xorm:"'pack_name'"`
	Head      int32   `xorm:"'head'"`
	Size      int32   `xorm:"'size'"`
	Key1      int32   `xorm:"'key1'"`
	Key2      int32   `xorm:"'key2'"`
	TimeStamp int64   `xorm:"'time_stamp'"` // the time this was first inserted
	Width     int32   `xorm:"'width'"`      // 0 means not calculated, -1 mean the pack is lost
	Height    int32   `xorm:"'height'"`     // 0 means not calculated, -1 mean the pack is lost
	Error     *string `xorm:"'error'"`
}

func (t *Texture) RawTexture() RawTexture {
	return RawTexture{
		PackName: t.PackName,
		Head:     t.Head,
		Size:     t.Size,
		Key1:     t.Key1,
		Key2:     t.Key2,
	}
}

func LoadUnencrypted(path string, texture RawTexture) []byte {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	unencrypted := make([]byte, texture.Size)
	_, err = file.ReadAt(unencrypted, int64(texture.Head))
	if err != nil {
		panic(err)
	}
	// keying
	hwdecrypt.DecryptBuffer(&hwdecrypt.HwdKeyset{
		Key1: uint32(texture.Key1),
		Key2: uint32(texture.Key2),
		Key3: 12345,
	}, unencrypted)
	return unencrypted
}

func LoadImageDirect(textureData RawTexture) (width, height int32, err error) {
	defer func() {
		r := recover()
		if r != nil {
			width = 0
			height = 0
			err = errors.New(fmt.Sprint(err))
		}
	}()
	file := GetStaticFile(textureData.PackName)
	if file == "" {
		err = errors.New(fmt.Sprintf("asset pack doesn't exist: %s", textureData.PackName))
		return
	}
	width, height, err = GetImageSize(LoadUnencrypted(file, textureData))
	return
}

func LoadTextureDirect(textureData RawTexture) (texture *graphic.Texture, err error) {
	defer func() {
		r := recover()
		if r != nil {
			texture.Free()
			texture = nil
			err = errors.New(fmt.Sprint(err))
		}
		// SaveTextureAdditionalDetail(textureData, texture)
	}()
	file := GetStaticFile(textureData.PackName)
	if file == "" {
		err = errors.New(fmt.Sprintf("asset pack doesn't exist: %s", textureData.PackName))
		return
	}
	texture = &graphic.Texture{}
	texture.LoadFromMemory(LoadUnencrypted(file, textureData))
	return
}

func LoadTextureData(assetPath string) (textureData Texture, err error) {

	tables := []string{"en", "ja", "zh", "ko", "th"}
	for _, table := range tables {
		exist, err2 := session.Table("texture_"+table).Where("asset_path = ?", assetPath).Get(&textureData)
		if err2 != nil {
			panic(err2)
		}
		if exist {
			return
		}
	}
	err = errors.New("can't find asset path: " + assetPath)
	return
}

func LoadTexture(assetPath string) (texture *graphic.Texture, err error) {
	defer func() {
		r := recover()
		if r != nil {
			texture = nil
			errors.New(fmt.Sprint(err))
		}
	}()
	tables := []string{"en", "ja", "zh", "ko", "th"}
	for _, table := range tables {
		textureData := Texture{}
		exist, err2 := session.Table("texture_"+table).Where("asset_path = ?", assetPath).Get(&textureData)
		if err2 != nil {
			panic(err2)
		}
		if exist {
			texture, err = LoadTextureDirect(textureData.RawTexture())
			if texture != nil {
				return
			}
		}
	}
	err = errors.New(fmt.Sprintf("asset path doesn't exist: %s", assetPath))
	return
}

func startupRawTextureDbUpdate(rawTextureChannel chan RawTexture, syncChannel chan struct{}) {
	cnt := 0
	for {
		rawTexture := <-rawTextureChannel
		if rawTexture.PackName == "finished" {
			flush()
			break
		}
		cnt++
		fmt.Println(cnt, rawTexture)
		_, err := session.Table("raw_texture").Where("pack_name = ? AND head = ? AND size = ? AND key1 = ? AND key2 = ?",
			rawTexture.PackName, rawTexture.Head, rawTexture.Size, rawTexture.Key1, rawTexture.Key2).AllCols().Update(&rawTexture)
		if err != nil {
			panic(err)
		}
		if cnt%1024 == 0 {
			flush()
		}
	}
	syncChannel <- struct{}{}
}

func CalculateAdditionalDetail(rawTexture RawTexture, rawTextureChannel chan RawTexture) {
	// always calculate and update the table
	width, height, err := LoadImageDirect(rawTexture)
	if err != nil {
		rawTexture.Error = new(string)
		*rawTexture.Error = fmt.Sprint(err)
	} else {
		rawTexture.Width = width
		rawTexture.Height = height
	}
	rawTextureChannel <- rawTexture
}

var textures = []RawTexture{}
var n int

func startupCalculateRawTexture(start, step int, rawTextureChannel chan RawTexture, syncChannel chan struct{}) {
	for ; start < n; start += step {
		CalculateAdditionalDetail(textures[start], rawTextureChannel)
	}
	syncChannel <- struct{}{}
}

func GetAllDetail() {
	// everytime this program is start, we calculate the details for textures that has not been calculated
	rawTextureChannel := make(chan RawTexture)
	syncChannel := make(chan struct{})
	go startupRawTextureDbUpdate(rawTextureChannel, syncChannel)
	err := session.Table("raw_texture").Where("width == 0 AND height == 0 AND error IS NULL").Find(&textures)
	if err != nil {
		panic(err)
	}
	n = len(textures)
	threadCount := 8
	for i := 0; i < threadCount; i++ {
		go startupCalculateRawTexture(i, threadCount, rawTextureChannel, syncChannel)
	}
	for i := 0; i < threadCount; i++ {
		<-syncChannel
	}
	rawTextureChannel <- RawTexture{
		PackName: "finished",
	}
	<-syncChannel
	updated := int64(0)
	for {
		texture := DetailedTexture{}
		exist, err := session.Table("detailed_texture").Where("width == 0 AND height == 0 AND error IS NULL").Get(&texture)
		if err != nil {
			panic(err)
		}
		if !exist {
			break
		}
		rawTexture := RawTexture{}
		exist, err = session.Table("raw_texture").Where("pack_name = ? AND head = ? AND size = ? AND key1 = ? AND key2 = ?",
			texture.PackName, texture.Head, texture.Size, texture.Key1, texture.Key2).Get(&rawTexture)
		if err != nil {
			panic(err)
		}
		if !exist {
			panic(fmt.Sprint("detaield texture isn't included in raw texture: ", texture))
		}

		texture.Width = rawTexture.Width
		texture.Height = rawTexture.Height
		texture.Error = rawTexture.Error
		affected, err := session.Table("detailed_texture").
			Where("pack_name = ? AND head = ? AND size = ? AND key1 = ? AND key2 = ?",
				rawTexture.PackName,
				rawTexture.Head,
				rawTexture.Size,
				rawTexture.Key1,
				rawTexture.Key2).Cols("width, height, error").
			Update(&texture)
		if err != nil {
			panic(err)
		}
		updated += affected
		if updated > 1024 {
			flush()
			updated = 0
		}
	}
	flush()

}
