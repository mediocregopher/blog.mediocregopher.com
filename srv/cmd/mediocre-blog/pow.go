package main

import (
	"encoding/hex"
	"errors"
	"fmt"
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

func requirePowMiddleware(mgr pow.Manager, h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		seedHex := r.PostFormValue("powSeed")
		seed, err := hex.DecodeString(seedHex)
		if err != nil || len(seed) == 0 {
			badRequest(rw, r, errors.New("invalid powSeed"))
			return
		}

		solutionHex := r.PostFormValue("powSolution")
		solution, err := hex.DecodeString(solutionHex)
		if err != nil || len(seed) == 0 {
			badRequest(rw, r, errors.New("invalid powSolution"))
			return
		}

		if err := mgr.CheckSolution(seed, solution); err != nil {
			badRequest(rw, r, fmt.Errorf("checking proof-of-work solution: %w", err))
			return
		}

		h.ServeHTTP(rw, r)
	})
}
