package netmap

import (
	// standard
	"log"
	"net"
	"net/http"
	"time"
	// third party
	"github.com/gorilla/websocket"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

type Network struct {
	globals *Globals
	socket *icmp.PacketConn
	upgrader *websocket.Upgrader
}

func (this *Network) Run(globals *Globals) {
	var err error
	this.globals = globals
	NodeState_Timeout = float64(globals.Ping + globals.Timeout)
	// ICMP
	this.socket, err = icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	HandleError(err)
	defer this.socket.Close()
	// Ping
	go this.pingLoop()
	go this.pongLoop()
	go this.stateLoop()
	// WebSockets
	this.upgrader = &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	http.HandleFunc("/", this.handler)
	log.Fatal(http.ListenAndServe(this.globals.Bind, nil))
}

// Connection handler
func (this *Network) handler(response http.ResponseWriter, request *http.Request) {
	connection, err := this.upgrader.Upgrade(response, request, nil)
	HandleError(err);
	client := &Client{
		connection: connection,
		send: make(chan []byte, 1024),
	}
	this.globals.Clients.add <- client
	for _, node := range this.globals.Map.Nodes {
		if node.State == NodeState_Up {
			continue
		}
		connection.WriteMessage(websocket.TextMessage, []byte(node.GetJson()))
	}
	go client.recieveLoop(this.globals)
	go client.sendLoop()
}

func (this *Network) pingLoop() {
	message := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{ ID: 1, Seq: 1 },
	}
	packet, err := message.Marshal(nil)
	HandleError(err)
	for {
		for _, node := range this.globals.Map.Nodes {
			if node.Ip.String() == "0.0.0.0" {
				continue
			}
			_, err := this.socket.WriteTo(packet, &node.Ip);
			HandleError(err)
			time.Sleep(time.Millisecond)
		}
		time.Sleep(time.Duration(this.globals.Ping) * time.Second)
	}
}

func (this *Network) pongLoop() {
	buffer := make([]byte, 1500)
	for {
		length, peer, err := this.socket.ReadFrom(buffer)
		if err != nil {
			log.Fatal(err)
			break
		}
		// iana.ProtocolICMP == 1
		message, err := icmp.ParseMessage(1, buffer[:length])
		if err != nil {
			log.Fatal(err)
			continue
		}
		if message.Type == ipv4.ICMPTypeEchoReply {
			ip := IP2Long(net.ParseIP(peer.String()))
			if node, ok := this.globals.Map.Nodes[ip]; ok {
				changed := node.SetLast(time.Now())
				if changed {
					this.globals.Clients.Broadcast(node.GetJson())
				}
			}
		}
	}
}

func (this *Network) stateLoop() {
	for {
		now := time.Now()
		for _, node := range this.globals.Map.Nodes {
			if now.Sub(node.Last).Seconds() > NodeState_Timeout {
				changed := node.ChangeState(1)
				if changed {
					this.globals.Clients.Broadcast(node.GetJson())
				}
			}
		}
		time.Sleep(time.Duration(this.globals.Ping) * time.Second)
	}
}
