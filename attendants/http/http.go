/*
Copyright (c) 2016 Nick Potts
Licensed to You under the GNU GPLv3
See the LICENSE file at github.com/npotts/homehub/LICENSE

This file is part of the HomeHub project
*/

package http

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
	"log"
	"net/http"
	"os"
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

type logger struct {
	// ALogger implements just enough log.Logger interface to be compatible with other implementations
	logger *log.Logger
}

func (l *logger) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()
	next(rw, r)
	if res, ok := rw.(negroni.ResponseWriter); ok {
		l.logger.Printf("%s %v from %v [%v %s] in %v", r.Method, r.URL.Path, r.Host, res.Status(), http.StatusText(res.Status()), time.Since(start))
		return
	}
	l.logger.Printf("%s %v from %v in %v", r.Method, r.URL.Path, r.Host, time.Since(start))

}

var _ = fmt.Println

/*HTTPd is a HTTP based object that listens for incoming JSON messages
via PUT or POST on / at the listening address.*/
type HTTPd struct {
	httpd   *graceful.Server //stoppable server
	mux     *mux.Router      //http router
	negroni *negroni.Negroni //middelware
	// regFxn, storeFxn homehub.RegStore //callback fxns
	user    string           //http info
	pass    string           //password
	stopper stoppable.Halter //atomic halter
	backend homehub.Backend  //storage backend
	stats   map[homehub.Alphabetic]int
	logger  *logger
}

/*Attendant returns a homehub.Attendant and a nil error*/
func Attendant(listen, user, password string) (homehub.Attendant, error) {
	return new(listen, user, password)
}

/*Use sets the backend*/
func (h *HTTPd) Use(backend homehub.Backend) {
	h.backend = backend
}

func new(listen, user, password string) (*HTTPd, error) {
	err := make(chan error)
	defer close(err)
	recovery := negroni.NewRecovery()
	logger := &logger{logger: log.New(os.Stdout, "[http] ", 0)}
	neg := negroni.New(recovery, logger)
	h := &HTTPd{
		backend: nil,
		stopper: stoppable.NewStopable(),
		mux:     mux.NewRouter(),
		negroni: neg,
		logger:  logger,
		httpd: &graceful.Server{
			Timeout: 100 * time.Millisecond, //no timeout, which has its own set of issues
			Server: &http.Server{
				Addr:           listen,
				Handler:        neg,
				ReadTimeout:    1 * time.Second,
				WriteTimeout:   1 * time.Second,
				MaxHeaderBytes: 1024 * 1024 * 1024 * 10, //10meg
			},
		},
		user:  user,
		pass:  SHA1HashedPassword(password),
		stats: map[homehub.Alphabetic]int{},
	}
	h.mux.HandleFunc("/", h.put).Methods("PUT")
	h.mux.HandleFunc("/", h.post).Methods("POST")
	h.mux.HandleFunc("/", h.get).Methods("GET") //Version info eventually?
	if password != "" && h.user != "" {
		h.logger.logger.Printf("Using Password %s:%s\n", h.user, h.pass)
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
		h.stats[m.Table]++
		return fxn(m)
	}
	return errNotValid

}

/*put handles incoming data formats to register*/
func (h *HTTPd) put(w http.ResponseWriter, r *http.Request) {
	if err := h.handleJSON(r, h.backend.Register); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

/*post handles 'inserting' actual data*/
func (h *HTTPd) post(w http.ResponseWriter, r *http.Request) {
	if err := h.handleJSON(r, h.backend.Store); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

/*post handles 'inserting' actual data*/
func (h *HTTPd) get(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	fmt.Fprintf(w, "Stats: %v", h.stats)
}
