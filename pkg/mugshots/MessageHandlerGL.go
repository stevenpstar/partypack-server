package mugshots

import (
	"strconv"
	"fake.com/pkg/logic"
	"fake.com/pkg/types"
)

// Message Names (indicating type of data associated with this action)
const (
  PLAYER_NAME string = "PlayerName"
  IMAGE              = "image"
  PROMPTONE          = "prompt1"
  PROMPTTWO          = "prompt2"
)

const (
  MsgPlayerImage string = "MSPIMAGE"
  MsgPromptOneID        = "MSPROMPTONE_ID"
  MsgPromptOne          = "MSPROMPTONE"
  MsgPromptTwoID        = "MSPROMPTTWO_ID"
  MsgPromptTwo          = "MSPROMPTTWO"
)

func HandlePlayerImage(clients map[*types.Client]bool, 
  body []types.Action, 
  pImage []PlayerImage) {

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
    if client.Name == player_id && client.State == FIRST_IMG {
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

  AddPromptsToPlayer(clients, player_id, prompt_one, prompt_two)
}

func CreateImageMessage(allData logic.AllData, mappedClients map[int]*types.Client,
  player int, image_index int, command string) types.MessageBody {
  var actions []types.Action

  var image1 types.Action
  image1.Name = "MSImage1"
  image1.Data = mappedClients[allData.Players[player].Images[0]].Images[image_index]

  actions = append(actions, image1)

  var image2 types.Action
  image2.Name = "MSImage2"
  image2.Data = mappedClients[allData.Players[player].Images[1]].Images[image_index]

  actions = append(actions, image2)

  return types.MessageBody{Body: actions, Command: command}
}

func createAction(name string, data string) types.Action {
  return types.Action {
    Name: name,
    Data: data,
  }
} 

// This seems to be creating the message that will be sent to the front end
// During the voting round, a single image and two (different) player prompts.
// * We are changing this to send all vote information for the round... because
// why the fuck not? *
func CreatePromptMessage(voteData []VotePrompt, mappedClients map[int]*types.Client,
  image_index int, command string) types.MessageBody {

  var vote VotePrompt
  for _, v := range voteData {
    if !v.Voted {
      vote = v
      break
    }
  }

  var actions []types.Action

  img := createAction(MsgPlayerImage, mappedClients[vote.ID].Images[image_index])
  p1_id := createAction(MsgPromptOneID, strconv.Itoa(vote.PromptOne.PlayerID))
  p1 := createAction(MsgPromptOne, vote.PromptOne.Prompt)
  p2_id := createAction(MsgPromptTwoID, strconv.Itoa(vote.PromptTwo.PlayerID))
  p2 := createAction(MsgPromptTwo, vote.PromptTwo.Prompt)

  actions = append(actions, img, p1_id, p1, p2_id, p2)

  return types.MessageBody{Body: actions, Command: command}
}

