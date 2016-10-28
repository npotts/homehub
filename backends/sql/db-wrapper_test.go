/*
Copyright (c) 2016 Nick Potts
Licensed to You under the GNU GPLv3
See the LICENSE file at github.com/npotts/homehub/LICENSE

This file is part of the HomeHub project
*/

package sql

import (
	"testing"

	"github.com/npotts/homehub"
)

func Test_New(t *testing.T) {
	_, err := New("postgres", ":memory:")
	if err == nil {
		t.Errorf("Should error out when we do not have a correct connection string")
	}

	n, err := New("sqlite3", ":memory:")
	defer n.Stop()
	if err != nil {
		t.Errorf("Should not fail with memory database: %v", err)
	}

	nn, err := Backend("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("Should not fail with memory database: %v", err)
	} else {
		nn.Stop()
	}
}

func TestSQLBackend_RegisterStore(t *testing.T) {
	bad, good := homehub.Datam{}, homehub.GoodSample

	q, e := New("sqlite3", ":memory:")
	if e != nil {
		t.Fatalf("Couldnt start instance: %v", e)
	}

	if e := q.Register(bad); e == nil {
		t.Errorf("Should get an error registering garbage")
	}
	if e := q.Store(bad); e == nil {
		t.Errorf("Should get an error storing garbage")
	}

	if e := q.Register(good); e != nil {
		t.Errorf("Should get an error registering garbage")
	}
	if e := q.Store(good); e != nil {
		t.Errorf("Should get an error storing garbage")
	}

}
