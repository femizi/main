package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)
var playerCount  []int

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type Message struct {
	Name  string
	Count int
	keypressed string
	score int
}
type ServerResponse struct {
	Name string
	Randomness string
	Paddle1Position int
	Paddle2Position int
}

func createRandomness() string {
	rand.Seed(time.Now().UnixNano())
	randomnumber := rand.Float64()

	return fmt.Sprintf("%f", randomnumber)

}

func homePage(w http.ResponseWriter, r *http.Request) {

}
func reader(conn *websocket.Conn) {
	
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		// if string(p)== "reset"{
		// 	conn.WriteMessage(messageType, []byte(createRandomness()))

		// }
		var message Message
		json.Unmarshal([]byte(p), &message)
		if message.Name == "connected" {
			if len(playerCount ) != 2{
				playerCount = append(playerCount, message.Count)
				fmt.Println(playerCount)
			}

		}

		log.Println(string(p))

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)

		}
		if len(playerCount) == 2 {
			conn.WriteMessage(messageType, []byte("enough players"))
			conn.WriteMessage(messageType, []byte(createRandomness()))
			
		}

	}
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Client CONNECTED SUCCESFULLY")
	reader(ws)
}

func setupRoutes() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/ws", wsEndpoint)
}

func main() {
	fmt.Println("Go websockets")
	fmt.Println(createRandomness())
	setupRoutes()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
