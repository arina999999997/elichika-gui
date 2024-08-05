package content

// this package load texture based on content
// it will show the content as if in a reward list or in present box

import (
	_ "modernc.org/sqlite"

	"elichika/client"
	"elichika/enum"
	"elichika/utils"

	"elichika/gui/graphic"
	"elichika/gui/sifas/asset"

	"fmt"

	"xorm.io/xorm"
)

var contentToThumbnail = map[int32]map[int32]string{}
var contentTypeToThumbnail = map[int32]string{}
var engine *xorm.Engine

// TODO(gui): This doesn't render the amount
func LoadTexture(content client.Content) (texture *graphic.Texture, err error) {
	assetPath := func() string {
		idMap, exist := contentToThumbnail[content.ContentType]
		if exist {
			_, exist = idMap[content.ContentId]
			if !exist {
				fmt.Println("failed to find content id, defaulting to content type common icon: ", content)
			} else {
				return idMap[content.ContentId]
			}
		}
		_, exist = contentTypeToThumbnail[content.ContentType]
		if exist {
			return contentTypeToThumbnail[content.ContentType]
		} else {
			fmt.Println("failed to find asset for content: ", content)
			return "null"
		}
	}()
	texture, err = asset.LoadTexture(assetPath)
	return texture, err
}

func SafeLoadTexture(content client.Content) *graphic.Texture {
	texture, err := LoadTexture(content)
	if err == nil {
		return texture
	} else {
		fmt.Println(err, "\nUsing default texture")
		return graphic.DefaultTexture()
	}
}

func loadTable(contentType int32, table, contentIdColumn, assetPathColumn string) {
	var err error
	contentIds := []int32{}
	assetPaths := []string{}
	err = engine.Table(table).OrderBy(contentIdColumn).Cols(contentIdColumn).Find(&contentIds)
	utils.CheckErr(err)
	err = engine.Table(table).OrderBy(contentIdColumn).Cols(assetPathColumn).Find(&assetPaths)
	utils.CheckErr(err)
	contentToThumbnail[contentType] = map[int32]string{}
	for i := range contentIds {
		contentToThumbnail[contentType][contentIds[i]] = assetPaths[i]
	}
}

func init() {
	var err error
	engine, err = xorm.NewEngine("sqlite", "assets/db/jp/masterdata.db")
	utils.CheckErr(err)
	// loadTable(enum.ContentTypeSnsCoin, "", "", "")
	loadTable(enum.ContentTypeCard, "m_card_appearance", "card_m_id", "thumbnail_asset_path") // TODO(gui): there's no outer frame
	// loadTable(enum.ContentTypeCardExp, "", "", "")
	loadTable(enum.ContentTypeGachaPoint, "m_gacha_point", "id", "thumbnail_asset_path")
	loadTable(enum.ContentTypeLessonEnhancingItem, "m_lesson_enhancing_item", "id", "thumbnail_asset_path")
	loadTable(enum.ContentTypeSuit, "m_suit", "id", "thumbnail_image_asset_path")
	// loadTable(enum.ContentTypeVoice, "", "", "")
	loadTable(enum.ContentTypeGachaTicket, "m_gacha_ticket", "id", "thumbnail_asset_path")
	// loadTable(enum.ContentTypeGameMoney, "", "", "")
	loadTable(enum.ContentTypeTrainingMaterial, "m_training_material", "id", "thumbnail_asset_path")
	loadTable(enum.ContentTypeGradeUpper, "m_grade_upper", "id", "thumbnail_asset_path")
	// loadTable(enum.ContentTypeGiftBox, "", "id", "") // unused
	loadTable(enum.ContentTypeEmblem, "m_emblem", "id", "emblem_asset_path") // TODO(gui): these things can have 2 assets stacked on top of each other
	loadTable(enum.ContentTypeRecoveryAp, "m_recovery_ap", "id", "thumbnail_asset_path")
	loadTable(enum.ContentTypeRecoveryLp, "m_recovery_lp", "id", "thumbnail_asset_path")
	// loadTable(enum.ContentTypeStorySide, "", "", "")
	// loadTable(enum.ContentTypeStoryMember, "", "", "")
	loadTable(enum.ContentTypeExchangeEventPoint, "m_exchange_event_point", "id", "thumbnail_asset_path")
	loadTable(enum.ContentTypeAccessory, "m_accessory", "id", "thumbnail_asset_path")
	loadTable(enum.ContentTypeAccessoryLevelUp, "m_accessory_level_up_item", "id", "thumbnail_asset_path")
	loadTable(enum.ContentTypeAccessoryRarityUp, "m_accessory_rarity_up_item", "id", "thumbnail_asset_path")
	loadTable(enum.ContentTypeCustomBackground, "m_custom_background", "id", "thumbnail_asset_path")
	loadTable(enum.ContentTypeEventMarathonBooster, "m_event_marathon_booster_item", "id", "thumbnail_asset_path")
	loadTable(enum.ContentTypeLiveSkipTicket, "m_live_skip_ticket", "id", "thumbnail_asset_path")
	// loadTable(enum.ContentTypeEventMiningBooster, "", "", "") // unused
	loadTable(enum.ContentTypeStoryEventUnlock, "m_story_event_unlock_item", "id", "thumbnail_asset_path")
	loadTable(enum.ContentTypeRecoveryTowerCardUsedCount, "m_recovery_tower_card_used_count_item", "id", "thumbnail_asset_path")
	// loadTable(enum.ContentTypeSubscriptionCoin, "", "", "")
	loadTable(enum.ContentTypeMemberGuildSupport, "m_member_guild_support_item", "id", "thumbnail_asset_path")

	// these icons are hardcoded in binary, but can be found elsewhere
	contentTypeToThumbnail[enum.ContentTypeSnsCoin] = "FH"
	contentTypeToThumbnail[enum.ContentTypeCardExp] = "1("
	contentTypeToThumbnail[enum.ContentTypeGameMoney] = ".C"
	contentTypeToThumbnail[enum.ContentTypeSubscriptionCoin] = "DN8"

	// TODO(gui): these icon do exist but they are only referenced in binary AFAIK
	// they are not too important to show as they never appear as content
	contentTypeToThumbnail[enum.ContentTypeVoice] = "null"
	contentTypeToThumbnail[enum.ContentTypeStorySide] = "null"
	contentTypeToThumbnail[enum.ContentTypeStoryMember] = "null"
}
