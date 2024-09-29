package main

import (
	"log"
	"net/http"
	lib "quick-chat-go/client/lib"
)

func main() {

	wsRoom := lib.NewRoom()
	http.Handle("/", &lib.TemplateHandler{Filename: "index.html"})
	http.Handle("/ws-room", wsRoom)
	go wsRoom.Run()

	port := ":8080"
	log.Println("Serving at", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
