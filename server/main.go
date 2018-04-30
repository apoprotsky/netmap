package main

import (
	// standard
	"log"
	// application
	"./netmap"
	// third party
	"github.com/go-ini/ini"
)

func main() {

	// Globals
	globals := &netmap.Globals{}

	// Read config file
	config, err := ini.Load("map.ini")
    netmap.HandleError(err);
	log.Println("Config loaded")

	// Init
	globals.Bind = config.Section("").Key("bind").String()
	globals.Ping = config.Section("").Key("ping").MustUint(10)
	globals.Timeout = config.Section("").Key("timeout").MustUint(1)
	globals.Free = config.Section("").Key("free").MustUint(0)
	globals.Map = &netmap.Map{}
	globals.Clients = &netmap.Clients{}
	globals.Map.Load(
		config.Section("").Key("user").String() + ":" +
		config.Section("").Key("password").String() + "@tcp(" +
		config.Section("").Key("host").String() + ":" +
		config.Section("").Key("port").String() + ")/" +
		config.Section("").Key("database").String(),
	)
	globals.Network = &netmap.Network{}

	// Run
	go netmap.FreeMemoryLoop(globals.Free)
	go globals.Map.Run(globals)
	go globals.Clients.Run(globals)
	globals.Network.Run(globals)

}
