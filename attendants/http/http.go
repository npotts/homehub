/*
Copyright (c) 2016 Nick Potts
Licensed to You under the GNU GPLv3
See the LICENSE file at github.com/npotts/homehub/LICENSE

This file is part of the HomeHub project
*/

package brainiac

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/npotts/go-patterns/stoppable"
	"github.com/pkg/errors"
	"github.com/tylerb/graceful"
	"github.com/urfave/negroni"
	"io"
	"net/http"
	"time"

	"github.com/npotts/homehub"
)

//SHA1HashedPassword computes and returns a SHA1 hashed password string that can be used in HTTP Auth routines
func SHA1HashedPassword(pass string) (hashed string) {
	hasher := sha1.New()
	io.WriteString(hasher, pass)
	hashed = "{SHA}" + base64.StdEncoding.EncodeToString(hasher.Sum(nil))
	return
}

var _ = fmt.Println

/*Config is a configuration structure used for a HTTPd attendant*/
type Config struct {
	HTTPUser     string
	HTTPPassword string
	HTTPListen   string
}

/*HTTPd is a HTTP based object that listens for incoming JSON messages
via PUT or POST on / at the listening address.*/
type HTTPd struct {
	httpd            *graceful.Server //stoppable server
	mux              *mux.Router      //http router
	negroni          *negroni.Negroni //middelware
	regFxn, storeFxn homehub.RegStore //callback fxns
	user             string           //http info
	pass             string           //password
	stopper          stoppable.Halter //atomic halter
}

func newHTTP(cfg Config, reg, store homehub.RegStore) (*HTTPd, error) {
	err := make(chan error)
	defer close(err)
	neg := negroni.Classic()
	h := &HTTPd{
		stopper: stoppable.NewStopable(),
		mux:     mux.NewRouter(),
		negroni: neg,
		httpd: &graceful.Server{
			Timeout: 100 * time.Millisecond, //no timeout, which has its own set of issues
			Server: &http.Server{
				Addr:           cfg.HTTPListen,
				Handler:        neg,
				ReadTimeout:    1 * time.Second,
				WriteTimeout:   1 * time.Second,
				MaxHeaderBytes: 1024 * 1024 * 1024 * 10, //10meg
			},
		},
		regFxn:   reg,
		storeFxn: store,
		user:     cfg.HTTPUser,
		pass:     SHA1HashedPassword(cfg.HTTPPassword),
	}
	h.mux.HandleFunc("/", h.put).Methods("PUT")
	h.mux.HandleFunc("/", h.post).Methods("POST")
	// h.mux.HandleFunc("/", h.get).Methods("GET") //Version info eventually?
	if h.pass != "" && h.user != "" {
		h.negroni.UseFunc(h.auth)
	}
	h.negroni.UseHandler(h.mux)

	go h.monitor(err)
	return h, <-err
}

/*monitor starts the HTTP server and attempts to keep it going*/
func (h *HTTPd) monitor(startup chan error) {
	ecc := make(chan error)
	go func() { ecc <- h.httpd.ListenAndServe(); close(ecc) }() //start daemon

	select {
	case <-time.After(100 * time.Millisecond):
		startup <- nil
	case e := <-ecc:
		startup <- e
	}
}

/*Stop kills the service*/
func (h *HTTPd) Stop() {
	defer h.stopper.Die()
	if h.stopper.Alive() {
		c := h.httpd.StopChan()
		go func() { h.httpd.Stop(100 * time.Millisecond) }()
		<-c
	}
	return
}

/*fill in with basic authentication validator*/
func (h *HTTPd) auth(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if _, _, ok := r.BasicAuth(); ok {
		next(w, r)
		return
	}
	w.WriteHeader(http.StatusUnauthorized)
}

var errHTTP = errors.New("Invalid HTTP data")
var errNotValid = errors.New("Invalid JSON Structure")

/*handleJSON breaks up json data*/
func (h *HTTPd) handleJSON(r *http.Request, fxn homehub.RegStore) error {
	data := make([]byte, r.ContentLength)
	if n, err := r.Body.Read(data); int64(n) != r.ContentLength || (err != nil && err != io.EOF) {
		return errHTTP
	}
	m := homehub.Datam{}

	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}

	if m.Valid() {
		return fxn(m)
	}
	return errNotValid

}

/*put handles incoming data formats to register*/
func (h *HTTPd) put(w http.ResponseWriter, r *http.Request) {
	if err := h.handleJSON(r, h.regFxn); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

/*post handles 'inserting' actual data*/
func (h *HTTPd) post(w http.ResponseWriter, r *http.Request) {
	if err := h.handleJSON(r, h.storeFxn); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

/*post handles 'inserting' actual data*/
// func (h *HTTPd) get(w http.ResponseWriter, r *http.Request) {
// 	w.WriteHeader(200)
// 	fmt.Fprintf(w, "Hello there")
// }