package main

import (
	"elichika/gui/graphic"
	"elichika/gui/sifas/content"
	// "time"
)

const itemWidth = 128
const itemHeight = 128
const buffer = 16

func GetPointRewardObject(rewards []EventMarathonPointReward) *graphic.Object {
	n := len(rewards)
	root := graphic.Object{
		Width:  itemWidth,
		Height: itemHeight*n + buffer*(n-1),
	}
	y := 0
	for _, reward := range rewards {
		texture := content.SafeLoadTexture(reward.RewardContent)
		itemObject := graphic.Object{}
		itemObject.LoadTexture(texture)
		frame := graphic.Object{
			Width:  itemWidth,
			Height: itemHeight,
		}
		frame.AnchorChildCenter(&itemObject)
		root.AnchorChildNative(&frame, 0, y)
		y += itemHeight + buffer
	}
	return &root
}

func GetRankingRewardObject(rewards []EventMarathonRankingReward) *graphic.Object {
	root := graphic.Object{
		Width:  0,
		Height: 0,
	}
	y := 0
	x := 0
	for i, reward := range rewards {
		if (i > 0) && (reward.HighestRank != rewards[i-1].HighestRank) {
			y += itemHeight + buffer
			x = 0
		}
		texture := content.SafeLoadTexture(reward.RewardContent)
		itemObject := graphic.Object{}
		itemObject.LoadTexture(texture)
		frame := graphic.Object{
			Width:  itemObject.Width,
			Height: itemHeight,
		}
		if frame.Width < itemWidth {
			frame.Width = itemWidth
		}
		frame.AnchorChildCenter(&itemObject)
		root.AnchorChildNative(&frame, x, y)
		if root.Width < x+frame.Width {
			root.Width = x + frame.Width
		}
		if root.Height < y+frame.Height {
			root.Height = y + frame.Height
		}
		x += frame.Width + buffer/2
	}
	return &root
}
func main() {

	window, err := graphic.NewWindow("Marathon reward list from file")
	if err != nil {
		panic(err)
	}

	for {
		// rewards := ReadPointRewardListFromFile("rewards.txt")
		// object := GetPointRewardObject(rewards)
		rewards := ReadRankingRewardListFromFile("rewards.txt")
		object := GetRankingRewardObject(rewards)
		window.RenderRootObject(object)
	}
}
