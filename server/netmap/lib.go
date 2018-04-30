package netmap

import (
	// standard
	"encoding/binary"
	"net"
	"runtime/debug"
	"time"
)

type Globals struct {
	Bind string
	Ping uint
	Timeout uint
	Free uint
	Map *Map
	Clients *Clients
	Network *Network
}

func FreeMemoryLoop(period uint) {
	if period == 0 {
		return
	}
	for {
		time.Sleep(time.Duration(period) * time.Second)
		debug.FreeOSMemory()
	}
}

func HandleError(e error) {
	if e != nil {
		panic(e)
	}
}

func Long2IP(long uint32) net.IP {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, long)
	return ip
}

func IP2Long(ip net.IP) uint32 {
	ip = ip.To4()
	return binary.BigEndian.Uint32(ip)
}
