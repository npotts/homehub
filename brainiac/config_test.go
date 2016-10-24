/*
 GNU GENERAL PUBLIC LICENSE
                       Version 3, 29 June 2007

 Copyright (C) 2007 Free Software Foundation, Inc. <http://fsf.org/>
 Everyone is permitted to copy and distribute verbatim copies
 of this license document, but changing it is not allowed.*/

package brainiac

import (
	"testing"
)

func Test_ConfigForArgs(t *testing.T) {
	type x struct {
		cfg  Config
		Args []string
	}

	tests := map[string]x{
		"defaults": x{
			cfg: Config{
				Driver:       "sqlite3",
				Source:       "brainiac.db",
				PIDLock:      "brianiac.pid",
				ListenHTTP:   false,
				HTTPUser:     "brainiac",
				HTTPPassword: "brainiac",
				HTTPListen:   "*:8080",
				ListenUDP:    false,
				UDPPort:      8080,
				ListenZmq:    false,
				ZmqAllow:     []string{},
				ZmqListen:    "tcp://*:8081",
			},
		},
		"multiple zmq allows": x{
			Args: []string{"-a", "host1", "-a", "host2"},
			cfg: Config{
				Driver:       "sqlite3",
				Source:       "brainiac.db",
				PIDLock:      "brianiac.pid",
				ListenHTTP:   false,
				HTTPUser:     "brainiac",
				HTTPPassword: "brainiac",
				HTTPListen:   "*:8080",
				ListenUDP:    false,
				UDPPort:      8080,
				ListenZmq:    false,
				ZmqAllow:     []string{"host1", "host2"},
				ZmqListen:    "tcp://*:8081",
			},
		},
	}
	for tname, x := range tests {
		if c := ConfigForArgs(x.Args); !c.Equals(x.cfg) {
			t.Logf("Checking %q", tname)
			t.Logf("Got: %v", c)
			t.Logf("Wnt: %v", &x.cfg)
			t.Errorf("Didnt get what I expected")
		}
	}
}
