/*
Copyright (c) 2016 Nick Potts
Licensed to You under the GNU GPLv3
See the LICENSE file at github.com/npotts/homehub/LICENSE

This file is part of the HomeHub project
*/

/*Package mango holds various wrappers around various protocols
in order to expose easier access to storage mechanisms.

Mostly attendants monitor some sort of input from remote clients and
via callbacks, send data to some sort of backend, which would be SQL,
NOSQL, binary blobs, or whatever else is desired.

Built Ins

Sub-package http provide a HTTP server that listens on a particular port
and optionally requires a username/password via Basic Auth before forwarding
data on.  Likewise, udp performs a similar service, but over UDP spockets.  It
is planned to also have a ZMQv4 setup
*/
package mango

import (
	"github.com/go-mangos/mangos"
	"github.com/go-mangos/mangos/protocol/rep"
	"github.com/go-mangos/mangos/protocol/req"
	"github.com/go-mangos/mangos/transport/ipc"
	"github.com/go-mangos/mangos/transport/tcp"
	"github.com/npotts/go-patterns/stoppable"

	"github.com/npotts/homehub"
)

/*Rep listens for incoming data feeds over mango sockets*/
type Rep struct {
	sock mangos.Socket
	be   homehub.Backend
	stpr stoppable.Halter
}

/*New returns a initialzied and running Rep*/
func New(url string) (r *Rep, err error) {
	r = &Rep{stpr: stoppable.NewStopable()}
	if r.sock, err = rep.NewSocket(); err != nil {
		return nil, err
	}
	r.sock.AddTransport(ipc.NewTransport())
	r.sock.AddTransport(tcp.NewTransport())

	if err = r.sock.Listen(url); err != nil {
		return nil, err
	}

	cerr := make(chan error)
	defer close(cerr)
	go r.monitor(cerr)
	return r, <-cerr
}

func (r *Rep) monitor(cerr chan error) {
	cerr <- nil
	for r.stpr.Alive() {
		msg, err := r.sock.Recv()

	}
}
