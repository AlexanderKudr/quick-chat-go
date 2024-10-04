package pkg

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sync"
)

var (
	connections = make([]http.ResponseWriter, 0)
	mu          sync.Mutex
)

func flusher(w http.ResponseWriter) {
	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush()
	}
}

func SSEHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher(w)

	mu.Lock()
	connections = append(connections, w)
	mu.Unlock()

	select {}
}

func ReceiveEvent(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading request body:", err)
		return
	}
	defer r.Body.Close()

	var eventData map[string]string
	err = json.Unmarshal(body, &eventData)
	if err != nil {
		log.Println("Error unmarshaling JSON data:", err)
		return
	}

	message := eventData["data"]
	log.Println("Received event from client:", message)

	SendToClients(message)

}

func SendToClients(message string) {
	mu.Lock()
	defer mu.Unlock()

	for _, conn := range connections {
		if _, err := conn.Write([]byte("data: " + message + "\n\n")); err != nil {
			log.Println("Error sending message to client:", err)
			continue
		}

		flusher(conn)
	}
}
