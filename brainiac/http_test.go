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

var hdatam Datam

func reg(datam Datam) error {
	hdatam = datam
	return nil
}

func store(datam Datam) error {
	hdatam = datam
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

	w = httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", fmt.Sprintf("http://%s/", hconfig.HTTPListen), strings.NewReader(""))
	req.SetBasicAuth(h.user, h.pass)
	h.auth(w, req, next)
	if !accessed {
		t.Errorf("Did not call next function")
	}
}

func TestHTTP_handleJSON(t *testing.T) {
	h, e := newHTTP(hconfig, reg, store)
	if e != nil {
		t.Fatalf("Unable to start: %v", e)
	}
	defer h.stop()
	//various ways to mess with it

	type x struct {
		route  string
		method string
		json   string
		length int64
		datam  Datam
		err    error
	}

	called := false
	callback := func(datam Datam) error {
		called = true
		fmt.Println("Called")
		return nil
	}

	reqForX := func(x x) *http.Request {
		called = false
		r, _ := http.NewRequest(x.method, fmt.Sprintf("http://%s%s", hconfig.HTTPListen, x.route), strings.NewReader(x.json))
		r.SetBasicAuth(h.user, h.pass)
		if x.length > 0 {
			r.ContentLength = x.length
		}
		return r
	}

	tests := []x{
		x{method: "GET", route: "/", json: `not json, but bad length`, length: 120223, err: errHttp},
		x{method: "GET", route: "/", json: `{"json":"but bad format"}`, err: errNotValid},
		x{method: "GET", route: "/", json: `{"table":"table", "data": {"field1": [1,2,3]}}`, err: errFormat},
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
	h, e := newHTTP(hconfig, reg, store)
	if e != nil {
		t.Fatalf("Unable to start: %v", e)
	}
	defer h.stop()
	client := http.Client{}
	getReqs := func(body string) (put *http.Request, post *http.Request) {
		url := fmt.Sprintf("http://localhost%s/", hconfig.HTTPListen)
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
