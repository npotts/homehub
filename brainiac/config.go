/*
 GNU GENERAL PUBLIC LICENSE
                       Version 3, 29 June 2007

 Copyright (C) 2007 Free Software Foundation, Inc. <http://fsf.org/>
 Everyone is permitted to copy and distribute verbatim copies
 of this license document, but changing it is not allowed.*/

package brainiac

import (
	"github.com/alecthomas/kingpin"
	// "os"
)

/*BrainiacConfig is a configuration structure for a brianiac instance*/
type BrainiacConfig struct {
	Driver  string
	Source  string
	PIDLock string

	//http backend
	ListenHTTP   bool
	HTTPUser     string
	HTTPPassword string
	HTTPPort     int

	//Listen on a UDP socket
	ListenUDP bool
	UDPPort   int

	//ZMQ settings
	ListenZmq bool
	ZmqAllow  []string
	ZmqListen string
}

var (
	app          = kingpin.New("brianiac", "HomeHub's Database Agent")
	driver       = app.Flag("driver", "Specify the drive to use.  One of 'sqlite3', 'postgres' or 'mysql'").Short('d').Default("sqlite3").String()
	source       = app.Flag("source", `Source of the database, usually something like "file.db", ":memory:", "postgres://user:password@server/database"`).Short('s').Default("brainiac.db").String()
	pidlock      = app.Flag("pid", `Path to PID file`).Short('i').Default("brianiac.pid").String()
	listenHTTP   = app.Flag("http", `Listen for requests over HTTP`).Short('H').Default("False").Bool()
	httpUser     = app.Flag("user", `Username to require for over HTTP`).Short('l').Default("brainiac").String()
	httpPassword = app.Flag("password", `Password for login over HTTP`).Short('p').Default("brainiac").String()
	httpPort     = app.Flag("http-port", `Port to start HTTP daemon on`).Short('P').Default("8080").Int()
	listenUDP    = app.Flag("udp", `Listen for requests over UDP`).Short('u').Default("False").Bool()
	udpPort      = app.Flag("udp-port", `Port to listen for incoming UDP packets on`).Short('U').Default("8080").Int()
	listenZmq    = app.Flag("zmq", `Listen for requests over ZMQ`).Short('z').Default("False").Bool()
	zmqAllow     = app.Flag("zmq-allow", `Allow ZMQ access from these hosts alone`).Short('a').Default("").Strings()
	zmqListen    = app.Flag("zmq-listen", `Port to listen for incoming ZMQ clients`).Short('Z').Default("tcp://*:8081").String()
)

/*ConfigForArgs returns a Brainiac config for a given set of arguments.*/
func ConfigForArgs(args []string) *BrainiacConfig {
	kingpin.MustParse(app.Parse(args))
	return &BrainiacConfig{
		Driver:       *driver,
		Source:       *source,
		PIDLock:      *pidlock,
		ListenHTTP:   *listenHTTP,
		HTTPUser:     *httpUser,
		HTTPPassword: *httpPassword,
		HTTPPort:     *httpPort,
		ListenUDP:    *listenUDP,
		UDPPort:      *udpPort,
		ListenZmq:    *listenZmq,
		ZmqAllow:     *zmqAllow,
		ZmqListen:    *zmqListen,
	}
}
