package main

import (
	"flag"
	"github.com/go-chi/chi"
	"log"
	"net/http"
)

func main() {

	bindAddr := flag.String("-bind", ":8080", "specifies the listen address for this server")
	flag.Parse()

	r := chi.NewRouter()

	err := http.ListenAndServe(*bindAddr, r)
	if err != nil {
		log.Printf("unable to bind server to address %s", *bindAddr)
	}

}
