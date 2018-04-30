package netmap

import (
	// standard
	"database/sql"
	"log"
	"net"
	"time"
	// third party
	_"github.com/go-sql-driver/mysql"
)

type Map struct {
	globals *Globals
	Nodes map[uint32]*Node
}

func (this *Map) Load(dsn string) {
	// Connect to database
	db, err := sql.Open("mysql", dsn)
	HandleError(err);
	defer db.Close()
	// Test database connection
	err = db.Ping()
	HandleError(err);
	log.Println("Connected to database")
	// Init
	this.Nodes = make(map[uint32]*Node)
	// Load nodes
	rows, err := db.Query("SELECT id, IFNULL(ip, 0), status FROM nodes");
	HandleError(err);
	for rows.Next() {
		var ip uint32
		var node Node
		err = rows.Scan(&node.Id, &ip, &node.State)
		HandleError(err)
		node.Ip = net.IPAddr{IP: Long2IP(ip)}
		node.Last = time.Now()
		node.Online = true
		this.Nodes[ip] = &node
	}
	log.Println("Nodes loaded")
}

func (this *Map) Run(globals *Globals) {
	this.globals = globals
}
