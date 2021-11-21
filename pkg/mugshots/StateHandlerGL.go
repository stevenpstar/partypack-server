package mugshots

import (
	"encoding/json"
	"fmt"

	"fake.com/pkg/logic"
	"fake.com/pkg/types"
)

func GetPlayerImageData(clients map[*types.Client]bool) logic.AllData {
  var allData = GetAllData(clients)
  logic.Pair(allData)
  b, err := json.Marshal(allData)
  if err != nil {
    fmt.Println("Error printing allData")
  }
  fmt.Println(string(b))
  return allData
}

//TODO rename this to something better once we have prompts worked out lmao
type PeePrompt struct {
  PlayerID int
  Prompt   string
}

func GetPlayerPromptData(clients map[int]*types.Client, allData logic.AllData) map[int][]PeePrompt {
  // prompt map structure
  // {id of client image: [{playerId: number, prompt: string}, {playerId: number, prompt: string}]
  var promptMap = make(map[int][]PeePrompt)
  for p, _ := range allData.Players {
    var pImages = allData.Players[p].Images
    for imageId, _ := range pImages {
      // id of the client (player) this image belongs to
      var id = pImages[imageId]
      // if the id exists in the map then add to it, otherwise create it
     // if _, exists := promptMap[id]; exists {
     //   var peePrompt = createPeePrompt(clients, allData.Players[p].PlayerId, p)
     //   promptMap[id] = append(promptMap[id], peePrompt)
     // } 
     var peePrompt = createPeePrompt(clients, allData.Players[p].PlayerId, p)
     promptMap[id] = append(promptMap[id], peePrompt)
    }
  }

  return promptMap
}

func createPeePrompt(clients map[int]*types.Client, playerId int, promptIndex int) PeePrompt {
  var peePrompt PeePrompt
  peePrompt.PlayerID = playerId
  peePrompt.Prompt = clients[playerId].Prompts[promptIndex]
  return peePrompt
}

func GetAllData(clients map[*types.Client]bool) logic.AllData {
  var allData logic.AllData
  for client, _ := range clients {
    if client.ID == 0 {
      continue
    }
    var playerData logic.PlayerData
    playerData.PlayerId = client.ID
    allData.Players = append(allData.Players, playerData)

    var imgData logic.ImageData
    imgData.PlayerId = client.ID
    allData.Images = append(allData.Images, imgData)
  }

  return allData
}

type PlayerPrompt struct {
  PlayerId int    `json:"playerId"`
  Prompt   string `json:"prompt"`
}

type PromptData struct {
  PlayerOne      PlayerPrompt `json:"playerOne"`
  PlayerTwo      PlayerPrompt `json:"playerTwo"`
  PlayerImage    string       `json:"playerImage"`
  PlayersImageId string       `json:"playersImageId"`
}

//func CreatePromptArray(clients map[*types.Client]bool, pairedData logic.AllData) {
//  var prompts []PromptData
//  for client, _ := range clients {
//    // skipping the host (game client)
//    if client.ID == 0 {
//      continue
//    }
//  }
//}

