package main

import (
	"elichika/client"
	"elichika/utils"

	"fmt"
	"strings"
)

// minimal event marathon point reward structure
// read from a text files with 4 columns:
// - each row is 1 record, with the first header row discarded
// - the cols are (required point, content type, content id, content amount)

type EventMarathonPointReward struct {
	RequiredPoint int32
	RewardContent client.Content
}

func ReadPointRewardListFromFile(path string) []EventMarathonPointReward {
	res := []EventMarathonPointReward{}
	text := utils.ReadAllText(path)
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		reward := EventMarathonPointReward{}
		n, err := fmt.Sscanf(line, "%d%d%d%d", &reward.RequiredPoint, &reward.RewardContent.ContentType, &reward.RewardContent.ContentId, &reward.RewardContent.ContentAmount)
		if (err != nil) || (n != 4) {
			fmt.Println("Error: ", err, "\nLine: ", line)
		} else {
			res = append(res, reward)
		}
	}
	return res
}
