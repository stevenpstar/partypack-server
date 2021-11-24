package logic

import (
	"math/rand"
	"time"
)

type PlayerData struct {
	PlayerId int   `json:"id"`
	Images   []int `json:"images"`
}

type PairData struct {
	Data AllData `json:"data"`
}

type ImageData struct {
	PlayerId int   `json:"playerId"`
	Paired   []int `json:"paired"`
}

type AllData struct {
	Players []PlayerData `json:"players"`
	Images  []ImageData  `json:"images"`
}


func Pair(data AllData) {
	players := data.Players
	images := data.Images

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(players), func(i, j int) {
		players[i], players[j] = players[j], players[i]
	})

	//shuffle images
	rand.Shuffle(len(images), func(i, j int) {
		images[i], images[j] = images[j], images[i]
	})

	var imgIndex = 0
	var index = 0
	var totalAttempts = len(players) * 10
	var attempt = 0
	for !allPaired(images) {
		attempt += 1
		if attempt > totalAttempts {
			for p := range data.Players {
				data.Players[p].Images = nil
			}

			for i := range data.Images {
				data.Images[i].Paired = nil
			}
			rand.Seed(time.Now().UnixNano())
			rand.Shuffle(len(players), func(i, j int) {
				players[i], players[j] = players[j], players[i]
			})

			//shuffle images
			rand.Shuffle(len(images), func(i, j int) {
				images[i], images[j] = images[j], images[i]
			})
			attempt = 0
		}
		if len(images[imgIndex].Paired) < 2 {
			for p := index; p < len(players); p++ {

				if images[imgIndex].PlayerId != players[p].PlayerId {
					if !Contains(players[p].Images, images[imgIndex].PlayerId) {
						// attaching player to image
						images[imgIndex].Paired = append(images[imgIndex].Paired, players[p].PlayerId)
						// attaching image to player
						players[p].Images = append(players[p].Images, images[imgIndex].PlayerId)

						//check if we have completed loop of players
						index = p + 1
						if index > len(players)-1 {
							index = 0
						}

						break
					}
				}
			}

			//check if we have completed loop of images

			imgIndex += 1
			if imgIndex > len(images)-1 {
				imgIndex = 0
			}
		}
	}
}

func Contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func allPaired(images []ImageData) bool {
	var paired = true
	for p := range images {
		if len(images[p].Paired) != 2 {
			paired = false
			break
		}
	}

	return paired
}
