package main

import (
	"fmt"
	"log"
	"net/http"
	"path"
	"strings"
)

type server struct {
	kvs    KVS
	hasher Hasher

	servStats *serverStats
}

// newServer creates a new server, with the specified hasher and key-value store.
func newServer(kvs KVS, hasher Hasher) *server {
	return &server{
		kvs:    kvs,
		hasher: hasher,

		// create a redirect counter
		servStats: newServerStats(),
	}
}

// shorten is an HTTP handler that shortens the URL string specified as 'url' param.
func (s *server) shorten(w http.ResponseWriter, r *http.Request) {
	// extract original URL.
	org := r.URL.Query().Get("url")
	if org == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "empty url field")
		return
	}

	// hash original url.
	short := s.hasher.Hash(org)
	log.Printf("/shorten org=%s short=%s", org, short)

	// store short url/original url pair in the key-value store.
	err := s.kvs.Store(short, org)
	if err != nil {
		// handle error
		log.Printf("/shorten error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "interval server error: %v", err)
		return
	}

	// print hostname and shortened url in the response.
	fullshort := path.Join(r.Host, short)
	fmt.Fprintf(w, "shortened URL: %s\n", fullshort)
}

// redirect is an HTTP handler that redirect a known short URL to its original URL.
func (s *server) redirect(w http.ResponseWriter, r *http.Request) {
	// consider we received a short url, extract it.
	short := strings.Replace(r.URL.RequestURI(), "/", "", 1)

	// load original url from the key-value store.
	org, err := s.kvs.Load(short)
	if err != nil {
		log.Printf("/ error: loading short=%s", short)
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "URL not found: %v", err)
		return
	}

	// forge the redirection URL.
	log.Printf("/ redirecting short=%s org=%s", short, org)

	// increment redirect counter for that url
	s.servStats.IncrRedirects(org)
	http.Redirect(w, r, org, 307)
}

// stats is an HTTP handler printing current stats.
func (s *server) stats(w http.ResponseWriter, r *http.Request) {
	// print the redirect counter into the response.
	err := s.servStats.Show(w)
	if err != nil {
		log.Println("/stats: error showing stats:", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "interval server error: %v", err)
	}
}
