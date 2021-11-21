package mugshots

import (
	"encoding/json"
	"fmt"
	"strconv"

	"fake.com/pkg/logic"
	"fake.com/pkg/types"
)

// Message Names (indicating type of data associated with this action)
const (
  PLAYER_NAME string = "PlayerName"
  IMAGE string = "image"
  PROMPTONE string = "prompt1"
  PROMPTTWO string = "prompt2"
)

func HandlePlayerImage(clients map[*types.Client]bool, body []types.Action, pImage []PlayerImage) {
  var player_id = "" // init empty string
  var image_data = ""
  for _, element := range body {
    if element.Name == PLAYER_NAME {
      player_id = element.Data
    }
  }

  for _, element2 := range body {
    if element2.Name == IMAGE {
      image_data = element2.Data
    }
  }

  //player image struct
  // playername string
  // image string

  var PImage PlayerImage
  PImage.PlayerName = player_id
  PImage.Image = image_data
  pImage = append(pImage, PImage)

  AddImageToPlayer(clients, player_id, image_data) 
}

func AddImageToPlayer(clients map[*types.Client]bool, player_id string, image_data string) {
  for client, _ := range clients {
    if client.Name == player_id && client.State == START {
      client.Images = append(client.Images, image_data)
      client.State = FIRST_IMG
      break
    }
  }
}

func AddPromptsToPlayer(clients map[*types.Client]bool, player_id string, prompt_one string, prompt_two string) {
  for client, _ := range clients {
    if client.Name == player_id && client.State == START {
      client.Prompts = append(client.Prompts, prompt_one)
      client.Prompts = append(client.Prompts, prompt_two)
      client.State = FIRST_PROMPT
      break
    }
  }
}

func HandlePlayerPrompt(clients map[*types.Client]bool, body []types.Action) {
  var player_id = ""
  var prompt_one = ""
  var prompt_two = ""
  for _, element := range body {
    if element.Name == PLAYER_NAME {
      player_id = element.Data
    }
    if element.Name == PROMPTONE {
      prompt_one = element.Data
    }
    if element.Name == PROMPTTWO {
      prompt_two = element.Data
    }
  }
//  for _, element := range body {
//    if element.Name == PROMPTONE {
//      prompt_one = element.Data
//    }
//  }
//
//  for _, element := range body {
//    if element.Name == PROMPTTWO {
//      prompt_two = element.Data
//    }
//  }

  AddPromptsToPlayer(clients, player_id, prompt_one, prompt_two)
}

func CreateImageMessage(allData logic.AllData, mappedClients map[int]*types.Client,
  player int, image_index int, command string) types.MessageBody {
  var actions []types.Action

  var image1 types.Action
  image1.Name = "MSImage1"
  image1.Data = mappedClients[allData.Players[player].Images[0]].Images[image_index]

  // b, b_err := json.Marshal(mappedClients[allData.Players[player].Images[0]].Images)
  // if b_err != nil { 
  //   fmt.Println(b_err)
  // }
  // fmt.Println("What is b?")
  // fmt.Println(string(b))

  actions = append(actions, image1)

  var image2 types.Action
  image2.Name = "MSImage2"
  image2.Data = mappedClients[allData.Players[player].Images[1]].Images[image_index]

  actions = append(actions, image2)

  return types.MessageBody{Body: actions, Command: command}
}

func CreatePromptMessage(promptData map[int][]PeePrompt, mugshotCount int, 
  mappedClients map[int]*types.Client, image_index int, command string) types.MessageBody {
  var image_data = "" // image to be sent to front end again

  var player_one_id = ""
  var player_one_prompt = ""

  var player_two_id = ""
  var player_two_prompt = ""
  var counter = 0
  for client, _ := range promptData {
    if client == 0 {
      // it should never get here but just in case
      continue
    }
    if counter != mugshotCount {
      counter++
      continue
    }
    image_data = mappedClients[client].Images[image_index]

    //player one data
    player_one_id = strconv.Itoa(promptData[client][0].PlayerID)
    player_one_prompt = promptData[client][0].Prompt
    //player two data
    player_two_id = strconv.Itoa(promptData[client][1].PlayerID)
    player_two_prompt = promptData[client][1].Prompt
  }

  var actions []types.Action

  var image types.Action
  image.Name = "MSPIMAGE"
  image.Data = image_data

  var prompt1_id types.Action
  prompt1_id.Name = "MSPROMPTONE_ID"
  prompt1_id.Data = player_one_id

  var prompt1 types.Action
  prompt1.Name = "MSPROMPTONE"
  prompt1.Data = player_one_prompt

  var prompt2_id types.Action
  prompt2_id.Name = "MSPROMPTTWO_ID"
  prompt2_id.Data = player_two_id

  var prompt2 types.Action
  prompt2.Name = "MSPROMPTTWO"
  prompt2.Data = player_two_prompt

  actions = append(actions, image)
  actions = append(actions, prompt1_id)
  actions = append(actions, prompt1)
  actions = append(actions, prompt2_id)
  actions = append(actions, prompt2)

  return types.MessageBody{Body: actions, Command: command}
}

