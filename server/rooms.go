package server

import (
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Participant describes a single entry in the hashmap
type Participant struct {
	Host bool
	Conn *websocket.Conn
}

// RoomMap is the main hashmap [roomID string] -> []Participant
type RoomMap struct {
	Mutex sync.RWMutex
	Map   map[string][]Participant
}

// Init initializes the RoomMap struct
func (r *RoomMap) Init() {
	r.Map = make(map[string][]Participant)
}

// Get returns the array of participants in a room
func (r *RoomMap) Get(roomID string) []Participant {
	r.Mutex.RLock()
	defer r.Mutex.RUnlock()
	return r.Map[roomID]
}

// CreateRoom generates a unique room ID, inserts it into the RoomMap, and returns the ID.
func (r *RoomMap) CreateRoom() string {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	rand.Seed(time.Now().UnixNano())
	letters := []rune("qwertyuiopasdfghjklzxcvbnm1234567890")

	const idLength = 8
	var roomID string

	for {
		b := make([]rune, idLength)
		for i := range b {
			b[i] = letters[rand.Intn(len(letters))]
		}
		roomID = string(b)
		if _, exists := r.Map[roomID]; !exists {
			break
		}
	}

	r.Map[roomID] = []Participant{}
	return roomID
}

// InsertIntoRoom inserts a participant into a room
func (r *RoomMap) InsertIntoRoom(roomID string, host bool, conn *websocket.Conn) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	p := Participant{Host: host, Conn: conn}
	log.Println("Inserting into Room with RoomID:", roomID)

	r.Map[roomID] = append(r.Map[roomID], p)
}

// DeleteRoom deletes the room with the given roomID
func (r *RoomMap) DeleteRoom(roomID string) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	delete(r.Map, roomID)
}
