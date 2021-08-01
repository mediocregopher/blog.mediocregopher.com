package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"

	"github.com/mediocregopher/blog.mediocregopher.com/srv/pow"
	"github.com/tilinna/clock"
)

func main() {
	staticDir := flag.String("static-dir", "", "Directory from which static files are served")
	listenAddr := flag.String("listen-addr", ":4000", "Address to listen for HTTP requests on")
	powTargetStr := flag.String("pow-target", "0x000FFFF", "Proof-of-work target, lower is more difficult")
	powSecret := flag.String("pow-secret", "", "Secret used to sign proof-of-work challenge seeds")

	// parse config

	flag.Parse()

	switch {
	case *staticDir == "":
		log.Fatal("-static-dir is required")
	case *powSecret == "":
		log.Fatal("-pow-secret is required")
	}

	powTargetUint, err := strconv.ParseUint(*powTargetStr, 0, 32)
	if err != nil {
		log.Fatalf("parsing -pow-target: %v", err)
	}
	powTarget := uint32(powTargetUint)

	// initialization

	clock := clock.Realtime()

	powStore := pow.NewMemoryStore(clock)
	defer powStore.Close()

	mgr := pow.NewManager(pow.ManagerParams{
		Clock:  clock,
		Store:  powStore,
		Secret: []byte(*powSecret),
		Target: powTarget,
	})

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(*staticDir)))
	mux.Handle("/api/pow/challenge", newPowChallengeHandler(mgr))

	// run

	log.Printf("listening on %q", *listenAddr)
	log.Fatal(http.ListenAndServe(*listenAddr, mux))
}
