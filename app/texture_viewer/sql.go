package main

import (
	"elichika/gui/graphic"
	"elichika/gui/graphic/button"
	"elichika/gui/graphic/textbox"

	"fmt"
)

// define sql functionality:
// - using SELECT query, select some assets from the texture db to be shown
// - Note that there is no checking on SQL

func (t *TextureDisplay) SetupSQL(w *graphic.Window) {

	t.SQLLabel, t.SQLTextbox = textbox.NewLabelAndRectTextbox(t, t.GetWidth()-200, 50, "SQL: ")

	t.SQLTextbox.Texture = graphic.RGBATexture(0x0f0f0fff)
	t.SQLTextbox.FocusTexture = graphic.RGBATexture(0x7f7f7fff)

	t.SQLButton = &button.RectButton{
		Width:   200,
		Height:  50,
		Texture: graphic.RGBATexture(0x7f7f7fff),
	}
	t.SQLButton.Text = graphic.NewText(t.ExtractButton, "Run SQL")

	displayAssetSql := func() {
		if t.SQLLength == 0 {
			statusChannel <- fmt.Sprint("There is no SQL result")
			return
		}
		t.SQLIndex %= t.SQLLength
		if t.SQLIndex < 0 {
			t.SQLIndex += t.SQLLength
		}
		statusChannel <- fmt.Sprintf("%s from %s (index %d out of %d)", t.SQLResult[t.SQLIndex].AssetPath, t.SQLResult[t.SQLIndex].PackName, t.SQLIndex, t.SQLLength)
		texture, err := LoadTextureDirect(t.SQLResult[t.SQLIndex].RawTexture())
		if err != nil {
			// TODO(gui): Dialog box
			statusChannel <- fmt.Sprint("Failed to load asset: ", t.SQLResult[t.SQLIndex], ", error: ", err)
			return
		}

		w.InternalEvent(func() {
			graphic.InvalidateRenderCache(t)
			t.AssetTexture.Free()
			t.AssetTexture = texture
		})
	}

	runSQL := func() {
		go func() {
			sql := t.SQLTextbox.TextContent
			if len(sql) <= 3 {
				sql = fmt.Sprint("SELECT * FROM texture_pack WHERE asset_path = \"", sql, "\"")
			}
			results, err := session.Query(sql)
			if err != nil {
				statusChannel <- fmt.Sprint(err)
				return
			}
			t.SQLResult = []Texture{}

			for _, result := range results {
				_, exists := result["key1"]
				if exists {
					assetPath := string(result["asset_path"])
					if assetPath == "" {
						assetPath = string(result["pack_name"]) + "_" + string(result["head"])
					}
					t.SQLResult = append(t.SQLResult, Texture{
						AssetPath: assetPath,
						PackName:  string(result["pack_name"]),
						Head:      GetInt32(string(result["head"])),
						Size:      GetInt32(string(result["size"])),
						Key1:      GetInt32(string(result["key1"])),
						Key2:      GetInt32(string(result["key2"])),
					})
				} else {
					assetPath, exists := result["asset_path"]
					if exists {
						texture, err := LoadTextureData(string(assetPath))
						if err == nil {
							t.SQLResult = append(t.SQLResult, texture)

						}
					}
				}
			}
			t.SQLLength = len(t.SQLResult)
			t.SQLIndex = 0
			displayAssetSql()
		}()
	}

	t.SQLButton.LeftClickHandler = runSQL
	t.SQLTextbox.OnEnterFunc = runSQL
	t.AddObject(t.SQLLabel, 0)
	t.AddObject(t.SQLTextbox, 0)
	t.AddObject(t.SQLButton, 0)

	t.NewLine()

	t.SQLNextButton = &button.RectButton{
		Width:   200,
		Height:  50,
		Texture: graphic.RGBATexture(0x7f7f7fff),
	}
	t.SQLNextButton.Text = graphic.NewText(t.SQLNextButton, "Next SQL")
	t.SQLNextButton.LeftClickHandler = func() {
		t.SQLIndex++
		go displayAssetSql()
	}
	t.SQLPreviousButton = &button.RectButton{
		Width:   200,
		Height:  50,
		Texture: graphic.RGBATexture(0x7f7f7fff),
	}
	t.SQLPreviousButton.LeftClickHandler = func() {
		t.SQLIndex--
		go displayAssetSql()
	}
	t.SQLPreviousButton.Text = graphic.NewText(t.SQLPreviousButton, "Previous SQL")

	t.AddObject(t.SQLNextButton, 5)
	t.AddObject(t.SQLPreviousButton, 5)
	// TODO(gui): add a text box to allow for abitrary indexing
	t.NewLine()
}
