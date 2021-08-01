package main

import (
	"encoding/hex"
	"net/http"

	"github.com/mediocregopher/blog.mediocregopher.com/srv/pow"
)

func newPowChallengeHandler(mgr pow.Manager) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		challenge := mgr.NewChallenge()

		jsonResult(rw, r, struct {
			Seed   string `json:"seed"`
			Target uint32 `json:"target"`
		}{
			Seed:   hex.EncodeToString(challenge.Seed),
			Target: challenge.Target,
		})
	})
}
