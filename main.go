package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

var (
	port   string
	client bool
)

func init() {
	flag.BoolVar(&client, "c", false, "a bool")
	flag.StringVar(&port, "port", "9999", "a string var")
	flag.Parse()
}

func main() {
	r := chi.NewRouter()
	wsUpgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	// Initialise the match pool
	pool := &MatchPool{playerQueue: make(chan *Player, 1), active: true}

	// If the client option was selected, statically host client files
	if client {
		fmt.Println("hosting client")
		r.Handle("/*", http.FileServer(http.Dir("./dist")))
	}

	// Accept websocket connections and queue player
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := wsUpgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		pool.QueuePlayer(conn)
	})

	// Begin matchmaking loop
	go pool.Matchmake()

	fmt.Printf("listening on port :%s\n", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
}
