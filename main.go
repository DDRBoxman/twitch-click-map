package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var clicks chan click

func main() {
	clicks = make(chan click, 100)

	r := mux.NewRouter()

	r.HandleFunc("/click", clickHandler)
	r.HandleFunc("/ws", wsHandler)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./public/")))

	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}

type click struct {
	X int
	Y int
}

func newClick(x, y int) *click {
	return &click{x, y}
}

func clickHandler(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var c click
	err := decoder.Decode(&c)
	if err != nil {
		return
	}

	clicks <- c
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	for {
		c := <-clicks
		if err = conn.WriteJSON(c); err != nil {
			log.Print(err)
			return
		}
	}

	conn.Close()
}
