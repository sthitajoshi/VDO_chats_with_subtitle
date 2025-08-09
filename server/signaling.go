package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// all rooms is the global hasmap for the server
var AllRooms RoomMap

// create a room and return room id
func CreateRoomRequestHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-control-Allow-Origin", "*")
	roomID := AllRooms.CreateRoom()

	type resp struct {
		RoomID string `json:"room_id"`
	}

	log.Println(AllRooms.Map)

	json.NewEncoder(w).Encode(resp{RoomID: roomID})
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type broadcastMsg struct {
	Message map[string]interface{}
	RoomID  string
	client  *websocket.Conn
}

var broadcast = make(chan broadcastMsg) //when msg come here we read that msg

func broadcaster() { //pass it to here and then it will broardcast it

	for {
		msg := <-broadcast
		for _, client := range AllRooms.Map[msg.RoomID] {
			if client.Conn != msg.client {
				err := client.Conn.WriteJSON(msg.Message)

				if err != nil {
					log.Fatal(err)
					client.Conn.Close()
				}
			}
		}
	}
}

// join roomRequrestHandaler this will join the client
func JoinRoomRequestHandler(w http.ResponseWriter, r *http.Request) {
	roomID, ok := r.URL.Query()["roomID"]

	if !ok {
		log.Println("RoomID missing in the URl Parameters")
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("web Socket Upgrade Error", err)
	}

	AllRooms.InsertIntoRoom(roomID[0], false, ws)

	go broadcaster()

	for {
		var msg broadcastMsg

		err := ws.ReadJSON(&msg.Message)
		if err != nil {
			log.Fatal("read ERROR:", err)
		}
		msg.client = ws
		msg.RoomID = roomID[0]

		broadcast <- msg
	}
}
