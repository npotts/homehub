/*
 GNU GENERAL PUBLIC LICENSE
                       Version 3, 29 June 2007

 Copyright (C) 2007 Free Software Foundation, Inc. <http://fsf.org/>
 Everyone is permitted to copy and distribute verbatim copies
 of this license document, but changing it is not allowed.*/

package brainiac

import (
	"github.com/alecthomas/kingpin"
)

type brainiacConfig struct {
	Driver *string
	Source *string
	PIDLoc *string

	//http backend
	ListenHTTP   *bool
	HTTPUser     *string
	HTTPPassword *string
	HTTPPort     *int

	//Listen on a UDP socket
	ListenUDP *bool
	UDPPort   *int

	//ZMQ settings
	ListenZmq *bool
	ZmqAllow  *[]string
	ZmqListen *string
}

var (
	app = kingpin.New("brianiac", "HomeHub's Database Agent")
	// http  = app.F
	config = &brainiacConfig{
		Driver: app.Flag("driver", "Specify the drive to use.  One of 'sqlite3', 'postgres' or 'mysql'").Short('d').Default("sqlite3").String(),
		Source: app.Flag("source", `Source of the database, usually something like "file.db", ":memory:", "postgres://user:password@server/database"`).Short('s').Default("brainiac.db").String(),
		PIDLoc: app.Flag("pid", `Path to PID file`).Short('i').Default("brianiac.pid").String(),

		ListenHTTP:   app.Flag("http", `Listen for requests over HTTP`).Short('H').Default("False").Bool(),
		HTTPUser:     app.Flag("user", `Username to require for over HTTP`).Short('l').Default("brainiac").String(),
		HTTPPassword: app.Flag("password", `Password for login over HTTP`).Short('p').Default("brainiac").String(),
		HTTPPort:     app.Flag("http-port", `Port to start HTTP daemon on`).Short('P').Default("8080").Int(),

		ListenUDP: app.Flag("udp", `Listen for requests over UDP`).Short('u').Default("False").Bool(),
		UDPPort:   app.Flag("udp-port", `Port to listen for incoming UDP packets on`).Short('U').Default("8080").Int(),

		ListenZmq: app.Flag("zmq", `Listen for requests over ZMQ`).Short('z').Default("False").Bool(),
		ZmqAllow:  app.Flag("zmq-allow", `Allow ZMQ access from these hosts alone`).Short('a').Default("").Strings(),
		ZmqListen: app.Flag("zmq-listen", `Port to listen for incoming ZMQ clients`).Short('Z').Default("tcp://*:8081").String(),
	}
)

func parse(args []string) {
	kingpin.MustParse(app.Parse(args))
}
