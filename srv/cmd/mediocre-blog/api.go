package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func internalServerError(rw http.ResponseWriter, r *http.Request, err error) {
	http.Error(rw, "internal server error", 500)
	log.Printf("%s %s: internal server error: %v", r.Method, r.URL, err)
}

func jsonResult(rw http.ResponseWriter, r *http.Request, v interface{}) {
	b, err := json.Marshal(v)
	if err != nil {
		internalServerError(rw, r, err)
		return
	}
	b = append(b, '\n')

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(b)
}
