/*
Copyright (c) 2016 Nick Potts
Licensed to You under the GNU GPLv3
See the LICENSE file at github.com/npotts/homehub/LICENSE

This file is part of the HomeHub project
*/

package http

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/npotts/homehub"
)

type fake struct{}

func (fake) Register(datam homehub.Datam) error {
	return nil
}
func (fake) Store(datam homehub.Datam) error {
	return nil
}
func (fake) Stop() {}

var (
	_        = fmt.Fscan
	listen   = ":4549"
	user     = "brainiac"
	password = "brainiac"
	faker    = &fake{}
)

var hdatam homehub.Datam

func getter() (h *HTTPd, e error) {
	h, e = new(listen, user, password)
	h.Use(faker)
	return h, e
}

func Test_newHTTP(t *testing.T) {
	h, e := getter()
	if e != nil {
		t.Fatalf("Unable to start: %v", e)
	}
	for i := 0; i < 100; i++ {
		defer h.Stop()
	}

	h2, e2 := Attendant(listen, user, password)
	if e2 == nil {
		defer h2.Stop()
		t.Errorf("Should not be able to start 2 servers on same port")
	}
}

func TestHTTP_Auth(t *testing.T) {
	h, e := getter()
	if e != nil {
		t.Fatalf("Unable to start: %v", e)
	}
	defer h.Stop()

	accessed := false
	next := func(w http.ResponseWriter, r *http.Request) {
		accessed = true
	}
	//access routes
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", fmt.Sprintf("http://%s/", listen), strings.NewReader(""))

	h.auth(w, r, next)
	if w.Code != http.StatusUnauthorized || accessed {
		t.Errorf("Should get Unauthorized Status: [%d] next() called: %v ", w.Code, accessed)
	}
	accessed = false

	w = httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", fmt.Sprintf("http://%s/", listen), strings.NewReader(""))
	req.SetBasicAuth(h.user, h.pass)
	h.auth(w, req, next)
	if !accessed {
		t.Errorf("Did not call next function")
	}
}

func TestHTTP_handleJSON(t *testing.T) {
	h, e := getter()
	if e != nil {
		t.Fatalf("Unable to start: %v", e)
	}
	defer h.Stop()
	//various ways to mess with it

	type x struct {
		route  string
		method string
		json   string
		length int64
		datam  homehub.Datam
		err    error
	}

	called := false
	callback := func(datam homehub.Datam) error {
		called = true
		fmt.Println("Called")
		return nil
	}

	reqForX := func(x x) *http.Request {
		called = false
		r, _ := http.NewRequest(x.method, fmt.Sprintf("http://%s%s", listen, x.route), strings.NewReader(x.json))
		r.SetBasicAuth(h.user, h.pass)
		if x.length > 0 {
			r.ContentLength = x.length
		}
		return r
	}

	tests := []x{
		x{method: "GET", route: "/", json: `not json, but bad length`, length: 120223, err: errHTTP},
		x{method: "GET", route: "/", json: `{"json":"but bad format"}`, err: errNotValid},
		// x{method: "GET", route: "/", json: `{"table":"table", "data": {"field1": [1,2,3]}}`, err: errFormat},
		x{method: "GET", route: "/", json: `{"table":"table", "data": {"field": 1.0}}`, err: nil},
	}

	for i, x := range tests {
		t.Logf("Running check #%d", i)
		if e := h.handleJSON(reqForX(x), callback); e != x.err {
			t.Logf(" Got:%v", e)
			t.Logf("Want:%v", x.err)
			t.Errorf("Errored out")
		}
	}
}

func TestHTTP_putpost(t *testing.T) {
	h, e := getter()
	if e != nil {
		t.Fatalf("Unable to start: %v", e)
	}
	defer h.Stop()
	client := http.Client{}
	getReqs := func(body string) (put *http.Request, post *http.Request) {
		url := fmt.Sprintf("http://localhost%s/", listen)
		fmt.Println("url=", url)
		if put, _ = http.NewRequest("PUT", url, strings.NewReader(body)); e != nil {
			t.Error(e)
			t.FailNow()
		}
		if post, e = http.NewRequest("POST", url, strings.NewReader(body)); e != nil {
			t.Error(e)
			t.FailNow()
		}
		put.SetBasicAuth(h.user, h.pass)
		post.SetBasicAuth(h.user, h.pass)
		return
	}
	process := func(body string, code int) {
		t.Logf("Running Checks with %s [expect %d]", body, code)
		put, post := getReqs(body)
		if r, e := client.Do(put); e != nil || r.StatusCode != code {
			t.Fatalf("PUT failed: %d:%v:%s", r.StatusCode, e, r.Status)
		}
		if r, e := client.Do(post); e != nil || r.StatusCode != code {
			t.Fatalf("POST failed: %d:%v:%s", r.StatusCode, e, r.Status)
		}
	}

	good, bad := `{"table":"table", "data": {"field": 1.0}}`, `{not json}`

	process(good, 200)
	process(bad, http.StatusBadRequest)
	// <-time.After(200 * time.Second)
}
