package exclude

// package websocket

// import (
// 	"fmt"
// 	"log"

// 	"encoding/json"
// 	"fake.com/pkg/types"
// 	"github.com/gorilla/websocket"
// )

// // type MessageBody struct {
// // 	Body    string `json:"body"`
// // 	Command string `json:"command"`
// // }

// func (c *types.Client) Read() {
// 	defer func() {
// 		if c.Pool != nil {
// 			c.Pool.Unregister <- c
// 			c.Conn.Close()
// 		}
// 	}()

// 	for {
// 		messageType, p, err := c.Conn.ReadMessage()

// 		if err != nil {
// 			log.Println(err)
// 			return
// 		}

// 		fmt.Printf(string(p))

// 		var actions []Action
// 		var messageBody MessageBody
// 		if messageError := json.Unmarshal(p, &messageBody); messageError != nil {
// 			fmt.Println(messageError)
// 			var err Action
// 			err.Name = "Error"
// 			err.Data = "There was an error"
// 			actions = append(actions, err)
// 		} else {
// 			actions = messageBody.Body

// 		}

// 		fmt.Println(string(messageBody.Body[0].Data))

// 		// var actions []Action
// 		// if actionError := json.Unmarshal(p, &actions); actionError != nil {
// 		// 	fmt.Println(actionError)
// 		// 	var err Action
// 		// 	err.Name = "Error"
// 		// 	err.Data = "There was an error"
// 		// 	actions = append(actions, err)
// 		// }

// 		message := Message{
// 			Type: messageType,
// 			Body: MessageBody{Body: actions, Command: messageBody.Command},
// 		}
// 		c.Pool.Broadcast <- message
// 		fmt.Printf("Message Received: %+v\n", message)

// 	}
// }
