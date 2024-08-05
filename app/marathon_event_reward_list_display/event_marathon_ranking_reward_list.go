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
// - the cols are (highest rank, content type, content id, content amount)
// - the ranking divisions are from this highest rank to the next highest rank - 1
// - multiple rewards require multiple row, depending on display order

type EventMarathonRankingReward struct {
	HighestRank   int32
	RewardContent client.Content
}

func ReadRankingRewardListFromFile(path string) []EventMarathonRankingReward {
	res := []EventMarathonRankingReward{}
	text := utils.ReadAllText(path)
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		reward := EventMarathonRankingReward{}
		n, err := fmt.Sscanf(line, "%d%d%d%d", &reward.HighestRank, &reward.RewardContent.ContentType, &reward.RewardContent.ContentId, &reward.RewardContent.ContentAmount)
		if (err != nil) || (n != 4) {
			fmt.Println("Error: ", err, "\nLine: ", line)
		} else {
			res = append(res, reward)
		}
	}
	return res
}
