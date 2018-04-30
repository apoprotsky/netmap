package netmap

import (
	// standard
	"log"
	// third party
	"github.com/gorilla/websocket"
)

type Client struct {
	globals *Globals
	connection *websocket.Conn
	send chan []byte
}

func (this *Client) recieveLoop(globals *Globals) {
	this.globals = globals
	defer func() {
		globals.Clients.del <- this
		this.connection.Close()
		close(this.send)
	}()
	for {
		_, _, err := this.connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("error: %v", err)
				}
			break
		}
	}
}

func (this *Client) sendLoop() {
	for {
		select {
			case packet := <- this.send:
				err := this.connection.WriteMessage(websocket.TextMessage, packet)
				if err != nil {
					return
				}
			break
		}
	}
}
