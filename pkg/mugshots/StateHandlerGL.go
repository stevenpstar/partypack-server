package mugshots

import (
	"fake.com/pkg/logic"
	"fake.com/pkg/types"
  "fmt"
)

func GetPlayerImageData(clients map[*types.Client]bool) logic.AllData {
  var allData = GetAllData(clients)
  logic.Pair(allData)
//  b, err := json.Marshal(allData)
//  if err != nil {
//    fmt.Println("Error printing allData")
//  }
//  fmt.Println(string(b))
  return allData
}

//TODO rename this to something better once we have prompts worked out lmao
type PeePrompt struct {
  PlayerID int
  Prompt   string
}

type VotePrompt struct {
  ID int
  PromptOne PeePrompt
  PromptTwo PeePrompt
  Voted bool
}

func SetVoted(voteData []VotePrompt) {
  for _, v := range voteData {
    if !v.Voted {
      v.Voted = true
      break
    }
  }
  fmt.Println("End of round?")
}

func CreateRoundVoteData(promptData map[int][]PeePrompt) []VotePrompt {
    var voteData []VotePrompt
    for id, prmpt := range promptData {
      // check "peeprompt" (terrible name) is of len(2)
      if len(prmpt) != 2 {
        fmt.Println("Error with pee prompt")
        break
      }
      if id == 0 {
        // This is the game client, they will not vote.
        // TODO also exclude users who created a prompt this round (future)
        continue
      }
      vPrompt := VotePrompt{
        ID: id,
        PromptOne: prmpt[0],
        PromptTwo: prmpt[1],
        Voted: false,
      }
      voteData = append(voteData, vPrompt)
    }
    return voteData
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
      var peePrompt = createPeePrompt(clients, allData.Players[p].PlayerId, imageId)
      promptMap[id] = append(promptMap[id], peePrompt)
    }
  }

  // TODO remove this temporary printing of the prompt map
  for pr, arr := range promptMap {
    // this might be the player image id I think?
    fmt.Printf("Player id: %v\n", pr)
    for _, e := range arr {
      fmt.Printf("id: %v, prompt: %s\n", e.PlayerID, e.Prompt)
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

