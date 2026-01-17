package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	_ "github.com/glebarez/go-sqlite" // Pure-Go SQLite driver
)

// Upgrader handles the transition from HTTP to WebSockets
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// ChatMessage represents the data structure shared between Go and React
type ChatMessage struct {
	Username  string `json:"username"`
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"`
}

// RealmHub manages active connections and message broadcasting
type RealmHub struct {
	clients   map[*websocket.Conn]bool
	broadcast chan ChatMessage
	mutex     sync.Mutex
	db        *sql.DB
}

func newHub(db *sql.DB) *RealmHub {
	return &RealmHub{
		clients:   make(map[*websocket.Conn]bool),
		broadcast: make(chan ChatMessage),
		db:        db,
	}
}

// Run listens for incoming messages and blasts them to all connected users
func (h *RealmHub) Run() {
	for {
		msg := <-h.broadcast

		// Save the message to the database
		_, err := h.db.Exec("INSERT INTO messages (username, content, timestamp) VALUES (?, ?, ?)",
			msg.Username, msg.Content, msg.Timestamp)
		if err != nil {
			log.Println("Database save error:", err)
		}

		// Marshal back to JSON to send to clients
		payload, _ := json.Marshal(msg)

		h.mutex.Lock()
		for client := range h.clients {
			err := client.WriteMessage(websocket.TextMessage, payload)
			if err != nil {
				log.Printf("Error sending to client: %v", err)
				client.Close()
				delete(h.clients, client)
			}
		}
		h.mutex.Unlock()
	}
}

func main() {
	// 1. Initialize SQLite (creates realm.db file if it doesn't exist)
	db, err := sql.Open("sqlite", "./realm.db")
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}
	defer db.Close()

	// 2. Create the messages table if it doesn't exist
	createTableSQL := `CREATE TABLE IF NOT EXISTS messages (
		id INTEGER PRIMARY KEY AUTOINCREMENT, 
		username TEXT, 
		content TEXT, 
		timestamp TEXT
	);`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal("Failed to create table:", err)
	}

	hub := newHub(db)
	go hub.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Upgrade error:", err)
			return
		}

		// 3. Load Chat History: Fetch last 50 messages for the new user
		rows, err := db.Query("SELECT username, content, timestamp FROM messages ORDER BY id DESC LIMIT 50")
		if err == nil {
			var history []ChatMessage
			for rows.Next() {
				var m ChatMessage
				rows.Scan(&m.Username, &m.Content, &m.Timestamp)
				// Prepend to keep chronological order (oldest at top)
				history = append([]ChatMessage{m}, history...)
			}
			rows.Close()

			for _, m := range history {
				payload, _ := json.Marshal(m)
				conn.WriteMessage(websocket.TextMessage, payload)
			}
		}

		// Register the new client
		hub.mutex.Lock()
		hub.clients[conn] = true
		hub.mutex.Unlock()

		// Listen for messages from this client
		go func() {
			defer func() {
				hub.mutex.Lock()
				delete(hub.clients, conn)
				hub.mutex.Unlock()
				conn.Close()
			}()

			for {
				_, payload, err := conn.ReadMessage()
				if err != nil {
					break
				}

				var incomingMsg ChatMessage
				if err := json.Unmarshal(payload, &incomingMsg); err != nil {
					log.Println("JSON unmarshal error:", err)
					continue
				}

				hub.broadcast <- incomingMsg
			}
		}()
	})

	fmt.Println("REALM SERVER IS LIVE")
	log.Fatal(http.ListenAndServe(":8080", nil))
}