package netmap

import (
	// standard
	"fmt"
	"log"
	"net"
	"time"
)

type NodeState int8

const (
	NodeState_Up      = iota
	NodeState_Notice  = iota
	NodeState_Warning = iota
	NodeState_Down    = iota
	NodeState_NoLink  = iota
)

var NodeState_Timeout float64

type Node struct {
	Id     uint32
	Ip     net.IPAddr
	State  NodeState
	Last   time.Time
	Online bool
}

func (this *Node) GetJson() string {
	return fmt.Sprintf("{\"code\":0,\"data\":{\"id\":%d,\"state\":%d}}", this.Id, this.State)
}

func (this *Node) ChangeState(delta NodeState) bool {
	oldState := this.State
	this.State += delta
	if this.State < NodeState_Up {
		this.State = NodeState_Up
	}
	if this.State > NodeState_Down {
		this.State = NodeState_Down
	}
	if this.State != oldState {
		if this.State == NodeState_Down {
			if this.Online {
				log.Println("Node down id:", this.Id, "ip:", this.Ip.String())
			}
			this.Online = false
		}
		if this.State == NodeState_Up {
			if !this.Online {
				log.Println("Node up id:", this.Id, "ip:", this.Ip.String())
			}
			this.Online = true
		}
	}
	return this.State != oldState
}

func (this *Node) SetLast(last time.Time) bool {
	oldLast := this.Last
	this.Last = last
	if last.Sub(oldLast).Seconds() < NodeState_Timeout {
		return this.ChangeState(-1)
	} else {
		return false
	}
}
