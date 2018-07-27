package netmap

import (
	// standard
	"log"
)

type Clients struct {
	globals   *Globals
	clients   map[*Client]bool
	add       chan *Client
	del       chan *Client
	broadcast chan string
}

func (this *Clients) Run(globals *Globals) {
	this.globals = globals
	this.clients = make(map[*Client]bool)
	this.add = make(chan *Client)
	this.del = make(chan *Client)
	this.broadcast = make(chan string)
	for {
		select {
		case client := <-this.add:
			this.clients[client] = true
			log.Println("Client connected", client.connection.RemoteAddr().String())
			break
		case client := <-this.del:
			_, ok := this.clients[client]
			if ok {
				log.Println("Client disonnected", client.connection.RemoteAddr().String())
				delete(this.clients, client)
			}
			break
		case message := <-this.broadcast:
			this.Broadcast(message)
			break
		}
	}
}

func (this *Clients) Broadcast(message string) {
	bytes := []byte(message)
	for client := range this.clients {
		select {
		case client.send <- bytes:
			break
		}
	}
}
