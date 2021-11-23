package mugshots

import (
	"fmt"
	"strconv"

	"fake.com/pkg/logic"
	"fake.com/pkg/types"
)

// Game State "Enum"
const (
	START         = 0
	FIRST_IMG     = 1
	FIRST_PROMPT  = 2
	SECOND_PROMPT = 3
	SECOND_IMG    = 4
	THIRD_PROMPT  = 5
	FOURTH_PROMPT = 6
)

const (
  PLAYER_IMAGE  string = "PI"
  PLAYER_PROMPT string = "PD"
  PROMPT_VOTE   string = "PV"
)

// Prompt State sub-state
const (
  WAITING = 0
  VOTING = 1
)

//var prompt_ones PromptOnes;
var prompt_ones Prompts
var pImages []PlayerImage
var pairedAllData logic.AllData
var mugshotCount int
// voting vars to keep track of voting rounds
var prompt_state = WAITING
var voting_count = 0
// all clients (id) and if they have voted
var votedClients map[int]bool 
// these are the ids of the clients whose prompts are up for voting
var nonVotingClients []int 

var currentMessage types.MessageBody

type Prompts struct {
  Prompts []Prompt
}

type PPrompt struct {
  PlayerID string
  Prompt string
}

type Prompt struct {
  PromptOne PPrompt
  PromptTwo PPrompt
  Image     PlayerImage
}

type ImagePrompt struct {
  PlayerName string
  Prompts []string
}

type PlayerImage struct {
  PlayerName string
  Image string
}

type PromptOnes struct {
  prompts []ImagePrompt
}

type MugshotsGL struct {
}

// check all players are at the pools state
func (m MugshotsGL) CheckAllPlayers(pool *types.Pool) bool {
	var allPlayers = true
	for c := range pool.Clients {
		if c.State != pool.State {
			allPlayers = false
			break
		}
	}
	return allPlayers
}

func (m MugshotsGL) CheckAllPlayersAreState(pool *types.Pool, state int) bool {
	var allPlayers = true
	for c := range pool.Clients {
		if c.State != state && c.ID != 0 {
      fmt.Println("Player not at state")
      fmt.Println(c.Name)
			allPlayers = false
			break
		}
	}
	return allPlayers
}

func (m MugshotsGL) GetGameState(pool *types.Pool) int {
	return pool.State
}

func (m MugshotsGL) HandleMessage(message types.Message, pool *types.Pool) {

  fmt.Println(message);

  var allData logic.AllData
  var mappedClients map[int]*types.Client

  switch message.Body.Command {
    case PLAYER_IMAGE:
      HandlePlayerImage(pool.Clients, message.Body.Body, pImages)

      if m.CheckAllPlayersAreState(pool, FIRST_IMG) {
        allData = GetPlayerImageData(pool.Clients)
        mugshotCount = len(allData.Images)
        pairedAllData = allData
        mappedClients = MapClients(pool.Clients)

        for p := range allData.Players {
          var message = CreateImageMessage(allData, mappedClients, p,  0, "PRMPT1")
     
          var _C = mappedClients[allData.Players[p].PlayerId]
          _C.Conn.WriteJSON(types.Message{Type: 2,
            Body: message})
        }
      } else if m.CheckAllPlayersAreState(pool, SECOND_IMG) {
        // returns paired data
        allData = GetPlayerImageData(pool.Clients)
        mappedClients = MapClients(pool.Clients)

        for p := range allData.Players {
          var message = CreateImageMessage(allData, mappedClients, p, 1, "PRMPT2")

          var _C = mappedClients[allData.Players[p].PlayerId]
          _C.Conn.WriteJSON(types.Message{Type: 2,
            Body: message})
        }
      }
      break
    case PLAYER_PROMPT:
      HandlePlayerPrompt(pool.Clients, message.Body.Body)
      fmt.Println("Received player prompt")
      if m.CheckAllPlayersAreState(pool, FIRST_PROMPT) {
          fmt.Println("All players promptin")
          // we need to set the sub-state to voting
          prompt_state = WAITING
          voting_count = 0
 //       allData = GetPlayerPromptData(pool.Clients)
          mappedClients = MapClients(pool.Clients)
          votedClients = ResetVotedClients(pool.Clients)
          var promptData = GetPlayerPromptData(mappedClients, pairedAllData)
          nonVotingClients = AddNonVotingClients(promptData)
          // send message to alert clients / players that it is time to vote!
          // this should actually just send a message to the game client to
          // start displaying the first prompt
         // for p := range allData.Players {
          if (prompt_state == WAITING) {
            currentMessage = CreatePromptMessage(promptData, voting_count, mappedClients,
              0, "VOTESEND")
            mappedClients[0].Conn.WriteJSON(types.Message{Type: 2,
              Body: currentMessage})
          }
          //}

      }
      break
    case PROMPT_VOTE:
      // just in case this happens somehow, don't actually do anything
      if prompt_state != VOTING {
        break
      }
      
      break
  }
  for client, _ := range pool.Clients {
    if err := client.Conn.WriteJSON(message); err != nil {
      return
    }
  }
}

func GetUnusedId(pool *types.Pool) int {

	var lastUnused = -1

	for id := 1; id < 9; id++ {
		var unused = true
		for k := range pool.Clients {
			if k.ID == id {
				unused = false
			}
		}
		if unused {
			return id
		}
	}

	return lastUnused
}

func MapClients(clients map[*types.Client]bool) map[int]*types.Client {
		var mappedClients = make(map[int]*types.Client)
		for client, _ := range clients {
			mappedClients[client.ID] = client
		}
    return mappedClients
}

func ResetVotedClients(clients map[*types.Client]bool) map[int]bool {
  var mappedClients = make(map[int]bool)
  for client, _ := range clients {
    mappedClients[client.ID] = false
  }
  return mappedClients
}

func AddNonVotingClients(clients map[int][]PeePrompt) []int {
  var nonVotingClients []int
  var count = 0
  for client, _ := range clients {
    if count != mugshotCount {
      count++
      continue
    }
    for player, _ := range clients[client] {
      nonVotingClients = append(nonVotingClients, clients[client][player].PlayerID)
    }
    break
  }

  return nonVotingClients
}

func (m MugshotsGL) HandleConnection(client *types.Client, pool *types.Pool) {

	// assigning ID

	if client.Type == "Host" {
		client.ID = 0
		pool.Clients[client] = true
	} else {
		id := GetUnusedId(pool)
		if id != -1 {
			client.ID = id
			pool.Clients[client] = true
		}
	}

	if client.Type == "Host" {

		var actions []types.Action
		var action types.Action
		action.Name = "RoomCode"
		action.Data = pool.Code
		actions = append(actions, action)

		client.Conn.WriteJSON(types.Message{Type: 2, Body: types.MessageBody{Body: actions, Command: "RC"}})

	} else {
		for c, _ := range pool.Clients {

			// on player join that isn't a host, notify the host client
			if c.Type == "Host" {

				var actions []types.Action
				var playerID types.Action
				playerID.Name = "PlayerId"
				playerID.Data = strconv.Itoa(client.ID)
				actions = append(actions, playerID)

				var playerName types.Action
				playerName.Name = "PlayerName"
				playerName.Data = client.Name
				actions = append(actions, playerName)

				c.Conn.WriteJSON(types.Message{Type: 2, Body: types.MessageBody{Body: actions, Command: "PC"}})
			}

		}
	}

}

func (m MugshotsGL) HandleDisconnect(client *types.Client, pool *types.Pool) {
	id := client.ID
	delete(pool.Clients, client)
	for cl := range pool.Clients {
		if cl.ID == 0 {

			var actions []types.Action

			var playerDisconnected types.Action
			playerDisconnected.Name = "PlayerDisconnected"
			playerDisconnected.Data = strconv.Itoa(id)

			actions = append(actions, playerDisconnected)

			cl.Conn.WriteJSON(types.Message{Type: 3,
				Body: types.MessageBody{Body: actions, Command: "PD"}})
		}
	}
}
