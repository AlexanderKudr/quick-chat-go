package main

import (
	"log"
	"net/http"

	"quick-chat-go/client/pkg"
	tmpl "quick-chat-go/client/template"
)

const port = ":8080"

func handleTemplatess() {
	http.Handle("/ws", tmpl.Template("index-ws.html"))
	http.Handle("/sse", tmpl.Template("index-sse.html"))
}

func handleWebsockets() {
	wsRoom := pkg.NewRoom()
	http.Handle("/ws-room", wsRoom)
	go wsRoom.Run()
}

func handleSse() {
	http.HandleFunc("/sse-room", pkg.SSEHandler)
	http.HandleFunc("/send-message", pkg.ReceiveEvent)
}

func serve() {
	log.Println("Serving at", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func main() {
	handleTemplatess()
	handleWebsockets()
	handleSse()
	serve()
}
