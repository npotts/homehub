/*
 GNU GENERAL PUBLIC LICENSE
                       Version 3, 29 June 2007

 Copyright (C) 2007 Free Software Foundation, Inc. <http://fsf.org/>
 Everyone is permitted to copy and distribute verbatim copies
 of this license document, but changing it is not allowed.*/

package main

import (
	"fmt"
	"github.com/alecthomas/kingpin"
	"github.com/vrecan/death"
	"os"
	"syscall"

	"github.com/npotts/homehub/attendants/http"
	"github.com/npotts/homehub/backends/sql"
)

var (
	app     = kingpin.New("brianiac", "HomeHub's Database Agent")
	pidlock = app.Flag("pid", `Path to PID file`).Short('i').Default("brianiac.pid").String()

	dbdriver = app.Flag("driver", "Specify the drive to use.  One of 'sqlite3', 'postgres' or 'mysql'").Short('d').Default("sqlite3").String()
	dbsource = app.Flag("source", `Source of the database, usually something like "file.db", ":memory:", "postgres://user:password@server/database"`).Short('s').Default("brainiac.db").String()

	// listenHTTP   = app.Flag("http", `Listen for requests over HTTP`).Short('H').Default("False").Bool()
	httpUser     = app.Flag("user", `Username to require for over HTTP.  Empty string means disable`).Short('l').Default("").String()
	httpPassword = app.Flag("password", `Password for login over HTTP`).Short('p').Default("").String()
	httpListen   = app.Flag("http-listen", `Listen Dial address.  Usually something like "localhost:9090" or "*:2442"`).Short('P').Default(":8080").TCP()

	// listenUDP = app.Flag("udp", `Listen for requests over UDP`).Short('u').Default("False").Bool()
	// udpPort   = app.Flag("udp-port", `Port to listen for incoming UDP packets on`).Short('U').Default("8080").Int()
	// listenZmq = app.Flag("zmq", `Listen for requests over ZMQ`).Short('z').Default("False").Bool()

	// zmqAllow  = app.Flag("zmq-allow", `Allow ZMQ access from these hosts alone`).Short('a').Strings()
	// zmqListen = app.Flag("zmq-listen", `Port to listen for incoming ZMQ clients`).Short('Z').Default("tcp://*:8081").String()
)

func main() {
	kingpin.MustParse(app.Parse(os.Args[1:]))
	be, err := sql.New(*dbdriver, *dbsource)
	if err != nil {
		fmt.Printf("Unable to initialize database:%v\n", err)
		os.Exit(1)
	}
	h, err := http.Attendant((*httpListen).String(), *httpUser, *httpPassword)
	if err != nil {
		fmt.Printf("Unable to initialize attendant:%v\n", err)
		os.Exit(1)
	}
	h.Use(be)

	file, err := os.Create(*pidlock)
	if err != nil {
		fmt.Printf("Unable to write PID file:%v\n", err)
		os.Exit(1)
	}
	file.Close()
	file.WriteString(fmt.Sprintf("%d", os.Getpid()))
	file.Close()

	closer := func() {
		h.Stop()
		be.Stop()
	}

	death.NewDeath(syscall.SIGINT, syscall.SIGTERM).WaitForDeathWithFunc(closer)
}
