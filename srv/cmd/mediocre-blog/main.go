package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	staticDir := flag.String("static-dir", "", "Directory from which static files are served")
	//redisAddr := flag.String("redis-addr", "127.0.0.1:6379", "Address which redis is listening on")
	listenAddr := flag.String("listen-addr", ":4000", "Address to listen for HTTP requests on")
	flag.Parse()

	if *staticDir == "" {
		log.Fatal("-static-dir is required")
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(*staticDir)))

	log.Printf("listening on %q", *listenAddr)
	log.Fatal(http.ListenAndServe(*listenAddr, mux))
}
