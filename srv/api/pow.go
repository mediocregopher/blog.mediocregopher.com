package api

import (
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"

	"github.com/mediocregopher/blog.mediocregopher.com/srv/api/apiutils"
)

func (a *api) newPowChallengeHandler() http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		challenge := a.params.PowManager.NewChallenge()

		apiutils.JSONResult(rw, r, struct {
			Seed   string `json:"seed"`
			Target uint32 `json:"target"`
		}{
			Seed:   hex.EncodeToString(challenge.Seed),
			Target: challenge.Target,
		})
	})
}

func (a *api) requirePowMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		seedHex := r.FormValue("powSeed")
		seed, err := hex.DecodeString(seedHex)
		if err != nil || len(seed) == 0 {
			apiutils.BadRequest(rw, r, errors.New("invalid powSeed"))
			return
		}

		solutionHex := r.FormValue("powSolution")
		solution, err := hex.DecodeString(solutionHex)
		if err != nil || len(seed) == 0 {
			apiutils.BadRequest(rw, r, errors.New("invalid powSolution"))
			return
		}

		err = a.params.PowManager.CheckSolution(seed, solution)

		if err != nil {
			apiutils.BadRequest(rw, r, fmt.Errorf("checking proof-of-work solution: %w", err))
			return
		}

		h.ServeHTTP(rw, r)
	})
}
