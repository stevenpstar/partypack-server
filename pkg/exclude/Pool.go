package exclude

// +exclude

// package websocket

// import (
// 	"encoding/json"
// 	"fake.com/pkg/logic"
// 	"fake.com/pkg/types"
// 	"fmt"
// 	"strconv"
// )

// func NewPool(code string, game logic.GameLogic) *Pool {
// 	return &types.Pool{
// 		Code:       code,
// 		Register:   make(chan *Client),
// 		Unregister: make(chan *Client),
// 		Clients:    make(map[*Client]bool),
// 		Broadcast:  make(chan Message),
// 		State:      0,
// 		Game:       game,
// 	}
// }

// func contains(s []int, e int) bool {
// 	for _, a := range s {
// 		if a == e {
// 			return true
// 		}
// 	}
// 	return false
// }

// func GetUnusedId(pool *types.Pool) int {

// 	var lastUnused = -1

// 	for id := 1; id < 9; id++ {
// 		var unused = true
// 		for k := range pool.Clients {
// 			if k.ID == id {
// 				unused = false
// 			}
// 		}
// 		if unused {
// 			return id
// 		}
// 	}

// 	return lastUnused
// }

// func ShuffleId(pool *types.Pool, id int) {

// 	if id < 8 && id != 0 {

// 		for n := (id); n <= 8; n++ {
// 			for k := range pool.Clients {
// 				if k.ID == n && k.ID != 0 {
// 					k.ID -= 1
// 					break
// 				}
// 			}
// 		}
// 	}
// }

// func (pool *types.Pool) Start() {

// 	for {
// 		select {
// 		case client := <-pool.Register:

// 			pool.Game.HandleConnection(client, pool)

// 			fmt.Println("Size of Connection Pool: ", len(pool.Clients))
// 			if client.Type == "Host" {

// 				var actions []Action
// 				var action Action
// 				action.Name = "RoomCode"
// 				action.Data = pool.Code
// 				actions = append(actions, action)

// 				fmt.Println("Hello I am sorry is this not working?")
// 				fmt.Printf("Room Code %s", pool.Code)
// 				client.Conn.WriteJSON(Message{Type: 2, Body: MessageBody{Body: actions, Command: "RC"}})
// 			} else {
// 				for c, _ := range pool.Clients {
// 					fmt.Println(c.Type)
// 					if c.Type == "Host" {
// 						fmt.Printf("C.Name")

// 						var actions []Action

// 						var playerId Action
// 						playerId.Name = "PlayerId"
// 						playerId.Data = strconv.Itoa(client.ID)

// 						actions = append(actions, playerId)

// 						var playerName Action
// 						playerName.Name = "PlayerName"
// 						playerName.Data = client.Name

// 						actions = append(actions, playerName)

// 						b, err := json.Marshal(actions)
// 						if err != nil {
// 							fmt.Printf("Error")
// 						}
// 						fmt.Print("Entire Message")
// 						fmt.Printf(string(b))
// 						c.Conn.WriteJSON(Message{Type: 2, Body: MessageBody{Body: actions, Command: "PC"}})
// 					}
// 				}
// 			}
// 			break
// 		case client := <-pool.Unregister:
// 			id := client.ID
// 			delete(pool.Clients, client)
// 			//ShuffleId(pool, id)
// 			fmt.Println("Size of Connection Pool: ", len(pool.Clients))
// 			for cl := range pool.Clients {
// 				if cl.ID == 0 {

// 					var actions []Action

// 					var playerDisconnected Action
// 					playerDisconnected.Name = "PlayerDisconnected"
// 					playerDisconnected.Data = strconv.Itoa(id)

// 					actions = append(actions, playerDisconnected)

// 					cl.Conn.WriteJSON(Message{Type: 3,
// 						Body: MessageBody{Body: actions, Command: "PD"}})
// 				}
// 			}

// 			break
// 		case message := <-pool.Broadcast:
// 			fmt.Println("Sending message to all clients in Pool")

// 			fmt.Println(message)

// 			for client, _ := range pool.Clients {
// 				if err := client.Conn.WriteJSON(message); err != nil {
// 					fmt.Println(err)
// 					return
// 				}
// 			}
// 		}
// 	}
// }
