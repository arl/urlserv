package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/arl/urlserv/demo"
)

// Hasher is the interface implemented by objects having a Hash method.
type Hasher interface {

	// Hash generates the hash value of v.
	Hash(v string) string
}

// KVS is key value store.
type KVS interface {

	// Store stores a key value pair.
	Store(k, v string) error

	// Load returns the value associated with a given key.
	Load(k string) (string, error)
}

func main() {
	addr := flag.String("addr", "localhost:7070", "listening address")
	flag.Parse()

	// create the server.
	s := newServer(demo.NewKVS(), demo.NewIDGenerator())

	// install HTTP handlers.
	http.HandleFunc("/shorten/", s.shorten)
	http.HandleFunc("/", s.redirect)
	http.HandleFunc("/stats", s.stats)

	log.Println("listening on", *addr)
	log.Fatalln(http.ListenAndServe(*addr, nil))
}
