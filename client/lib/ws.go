package ws

import (
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/gorilla/websocket"
)

type Client struct {
	socket   *websocket.Conn
	receiver chan []byte

	room *room
}

func (c *Client) read() {
	defer c.socket.Close()
	for {
		_, message, err := c.socket.ReadMessage()

		if err != nil {
			return
		}

		c.room.forward <- message
	}
}

func (c *Client) write() {
	defer c.socket.Close()
	for message := range c.receiver {
		err := c.socket.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			return
		}
	}
}

type room struct {
	clients map[*Client]bool

	join    chan *Client
	leave   chan *Client
	forward chan []byte
}

func NewRoom() *room {
	return &room{
		clients: make(map[*Client]bool),
		join:    make(chan *Client),
		leave:   make(chan *Client),
		forward: make(chan []byte),
	}
}

func (r *room) Run() {
	for {
		select {
		case client := <-r.join:
			r.clients[client] = true
		case client := <-r.leave:
			delete(r.clients, client)
			close(client.receiver)
		case message := <-r.forward:
			for client := range r.clients {
				select {
				case client.receiver <- message:
				default:
					close(client.receiver)
					delete(r.clients, client)
				}
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("Serve room:", err)
		return
	}

	client := &Client{
		socket:   socket,
		receiver: make(chan []byte, messageBufferSize),
		room:     r,
	}

	r.join <- client
	defer func() { r.leave <- client }()
	go client.write()
	client.read()
}

type TemplateHandler struct {
	once     sync.Once
	Filename string
	templ    *template.Template
}

func (t *TemplateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("template", t.Filename)))
	})

	if t.templ == nil {
		log.Fatal("Failed to parse template")
	}

	t.templ.Execute(w, r)
}
