package types

import (
	"encoding/json"
	"github.com/gorilla/websocket"
)

type GameLogic interface {
	CheckAllPlayers(pool *Pool) bool
	CheckAllPlayersAreState(pool *Pool, state int) bool
	GetGameState(pool *Pool) int
	HandleConnection(client *Client, pool *Pool)
	HandleDisconnect(client *Client, pool *Pool)
	HandleMessage(message Message, pool *Pool)
}

type Message struct {
	Type int         `json:"type"`
	Body MessageBody `json:"body"`
}

type MessageBody struct {
	Body    []Action `json:"body"`
	Command string   `json:"command"`
}

type Action struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type Pool struct {
	Code       string
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
	Broadcast  chan Message
	State      int
	Game       GameLogic
}

func NewPool(code string, game GameLogic) *Pool {
	return &Pool{
		Code:       code,
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan Message),
		State:      0,
		Game:       game,
	}
}

func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.Register:

			pool.Game.HandleConnection(client, pool)

			break
		case client := <-pool.Unregister:

			pool.Game.HandleDisconnect(client, pool)

			break
		case message := <-pool.Broadcast:
			pool.Game.HandleMessage(message, pool)

			break
		}
	}
}

// PlayerID is the player whose image is used I believe
type TestPrompt struct {
  PlayerID int
  Prompt string
}

type Client struct {
	ID     int
	Name   string
	Type   string
	Conn   *websocket.Conn
	Pool   *Pool
	State  int
	Images []string
  Prompts []string
}

func (c *Client) Read() {
	defer func() {
		if c.Pool != nil {
			c.Pool.Unregister <- c
			c.Conn.Close()
		}
	}()

	for {
		messageType, p, err := c.Conn.ReadMessage()

		if err != nil {
			return
		}

		var actions []Action
		var messageBody MessageBody
		if messageError := json.Unmarshal(p, &messageBody); messageError != nil {
			var err Action
			err.Name = "Error"
			err.Data = "There was an error"
			actions = append(actions, err)
		} else {
			actions = messageBody.Body
		}

		// these messages may also be handled by game logic layer

		message := Message{
			Type: messageType,
			Body: MessageBody{Body: actions, Command: messageBody.Command},
		}
		c.Pool.Broadcast <- message

	}
}
