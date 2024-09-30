package main

import (
	"log"
	"net/http"

	"quick-chat-go/client/pkg"
	"quick-chat-go/client/template"
)

const port = ":8080"

func main() {

	wsRoom := ws.NewRoom()
	http.Handle("/ws", &serve.TemplateHandler{Filename: "index.html"})
	http.Handle("/ws-room", wsRoom)
	go wsRoom.Run()

	log.Println("Serving at", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
