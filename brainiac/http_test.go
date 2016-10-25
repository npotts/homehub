/*
 GNU GENERAL PUBLIC LICENSE
                       Version 3, 29 June 2007

 Copyright (C) 2007 Free Software Foundation, Inc. <http://fsf.org/>
 Everyone is permitted to copy and distribute verbatim copies
 of this license document, but changing it is not allowed.*/

package brainiac

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var _ = fmt.Fscan

var hconfig = Config{
	HTTPListen:   ":4549",
	HTTPUser:     "brainiac",
	HTTPPassword: "brainiac",
}

var htable string
var hdatam Datam

func reg(table string, datam Datam) error {
	htable, hdatam = table, datam
	return nil
}

func store(table string, datam Datam) error {
	htable, hdatam = table, datam
	return nil
}

func Test_newHTTP(t *testing.T) {
	h, e := newHTTP(hconfig, reg, store)
	if e != nil {
		t.Fatalf("Unable to start: %v", e)
	}
	for i := 0; i < 100; i++ {
		defer h.stop()
	}

	h2, e2 := newHTTP(hconfig, reg, store)
	if e2 == nil {
		defer h2.stop()
		t.Errorf("Should not be able to start 2 servers on same port")
	}
}

func TestHTTP_Auth(t *testing.T) {
	h, e := newHTTP(hconfig, reg, store)
	if e != nil {
		t.Fatalf("Unable to start: %v", e)
	}
	defer h.stop()

	accessed := false
	next := func(w http.ResponseWriter, r *http.Request) {
		accessed = true
	}
	//access routes
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", fmt.Sprintf("http://%s/", hconfig.HTTPListen), strings.NewReader(""))

	h.auth(w, r, next)
	if w.Code != http.StatusUnauthorized || accessed {
		t.Errorf("Should get Unauthorized Status: [%d] next() called: %v ", w.Code, accessed)
	}
	accessed = false

	//TODO: Added auth backend

}
