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

type Message struct {
	Name       string
	Count      int
	keypressed string
	score      int
}
type ServerResponse struct {
	Name            string `json:"name"`
	Randomness      string `json:"randomness"`
	Paddle1Position int    `json:"paddle1position"`
	Paddle2Position int    `json:"paddle2position"`
}
var client [] *websocket.Conn
var playerCount []int
var message Message
var serverResponse ServerResponse
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
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
		fmt.Println(messageType)

		json.Unmarshal([]byte(p), &message)
		switch message.Name {
		case "connected":
			if len(playerCount) != 2 {
				playerCount = append(playerCount, message.Count)
				client = append(client, conn)
			
				if len(playerCount) == 2 {
					serverResponse.Name = "enough players"
					serverResponse.Randomness = createRandomness()
					serverResponse.Paddle1Position = 300
					serverResponse.Paddle2Position = 300
					for _, element:= range client{
						
						element.WriteJSON(serverResponse)

					} 
					

				}
			}
		case "reset":
			serverResponse.Name = "Ball reset"
			serverResponse.Randomness = createRandomness()
			serverResponse.Paddle1Position = 150
			serverResponse.Paddle2Position = 150
			for _, element:= range client{
				
				element.WriteJSON(serverResponse)

			} 
			
		case "key pressed w":
			serverResponse.Name = "update Paddle1 up"
			serverResponse.Paddle1Position = 20
			for _, element:= range client{
				
				element.WriteJSON(serverResponse)

			} 
		case "key pressed s":
			serverResponse.Name = "update Paddle1 down"
			serverResponse.Paddle1Position = 20
			for _, element:= range client{
				
				element.WriteJSON(serverResponse)

			} 
		


		}

		log.Println(string(p))

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
