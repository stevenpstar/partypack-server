package main

import (
	"encoding/json"
	"fake.com/pkg/mugshots"
	"fake.com/pkg/types"
	"fake.com/pkg/websocket"
  "fake.com/pkg/logger"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

var router *mux.Router
var rooms map[string]types.Pool

func playerWs(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Upgrade(w, r)
	if err != nil {
		logger.LogError("SOMETHING WENT WRONG HERE OH NO")
		fmt.Fprintf(w, "%+v\n", err)
	}

	querycode := strings.Split(r.URL.String(), "/")
	code := querycode[len(querycode)-2]
	name := querycode[len(querycode)-1]

	fmt.Println(querycode)
	logger.Log("Player is uhhh joining?")

	if val, ok := rooms[code]; ok {
		fmt.Println(val)
		client := &types.Client{
			Type:  "Player",
			Name:  name,
			Conn:  conn,
			Pool:  &val,
			State: 0,
		}
		rooms[code].Register <- client
		logger.Log("joined!")
		client.Read()
	} else {
		logger.LogError("We have errrrrrr")
	}

}

func serveWs(w http.ResponseWriter, r *http.Request) {
	log.Println("Logging here")
	conn, err := websocket.Upgrade(w, r)
	var mapRoom = generateRoomCode()
	var g mugshots.MugshotsGL
	pool := types.NewPool(mapRoom, g)
	rooms[mapRoom] = *pool

	go pool.Start()

	if err != nil {
		fmt.Fprintf(w, "%+v\n", err)
	}

	client := &types.Client{
		Type:  "Host",
		Conn:  conn,
		Pool:  pool,
		State: 0,
	}

	pool.Register <- client
	client.Read()
}

type Code struct {
	Code string `json:"code"`
}

func joinRoom(w http.ResponseWriter, r *http.Request) {
	if (r).Method == "OPTIONS" {
		fmt.Println("in")
		return
	}
	decoder := json.NewDecoder(r.Body)
	var c Code
	d := decoder.Decode(&c)
	if d != nil {
		log.Println(d)
	}
	log.Println("Attempting to join", c.Code)
	if value, ok := rooms[c.Code]; ok {
		fmt.Println(value)
		w.Write([]byte(c.Code))
	} else {
		w.Write([]byte("Not Found"))
	}
	for k := range rooms {
		log.Println(k + ", ")
	}
}

func generateRoomCode() string {
	//generate room code
	rand.Seed(time.Now().UnixNano())

	abc := [26]string{"A", "B", "C", "D", "E",
		"F", "G", "H", "I", "J",
		"K", "L", "M", "N", "O",
		"P", "Q", "R", "S", "T",
		"U", "V", "W", "X", "Y",
		"Z"}

	numbers := [9]string{"1", "2", "3",
		"4", "5", "6", "7",
		"8", "9"}

	var code = ""

	for {

		for i := 0; i < 4; i++ {
			var n = rand.Intn(2)
			if n == 0 {
				code += abc[rand.Intn(len(abc))]
			} else {
				code += numbers[rand.Intn(len(numbers))]
			}
		}

		if _, ok := rooms[code]; !ok {
			return code
    }

	}
}

func preFlight(w http.ResponseWriter, r *http.Request) {
	log.Println("Preflight")
}

func main() {
  logger.Log("Starting Server")
	//rooms
	rooms = make(map[string]types.Pool)
	//router
	router = mux.NewRouter()
	router.HandleFunc("/join", joinRoom).Methods("POST")
	router.HandleFunc("/join", preFlight).Methods("OPTIONS")
	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		logger.Log("Loggin here")
		serveWs(w, r)
	})
	router.HandleFunc("/cws/{key}/{name}", func(w http.ResponseWriter, r *http.Request) {
		logger.Log("Player Loggin here")
		playerWs(w, r)
	})
	//HandleFuncEx("/join", joinRoom)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PATCH", "OPTIONS"},
	})

	handler := c.Handler(router)
	//room map

	//setupRoutes()
	http.ListenAndServe(":8080", handler)
}

func HandleFuncEx(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	log.Println("handled function", pattern)
	router.HandleFunc(pattern, handler)
}
