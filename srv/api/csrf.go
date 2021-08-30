package api

import (
	"errors"
	"net/http"
)

const (
	csrfTokenCookieName = "csrf_token"
	csrfTokenHeaderName = "X-CSRF-Token"
)

func setCSRFMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		csrfTok, err := getCookie(r, csrfTokenCookieName, "")

		if err != nil {
			internalServerError(rw, r, err)
			return

		} else if csrfTok == "" {
			http.SetCookie(rw, &http.Cookie{
				Name:   csrfTokenCookieName,
				Value:  randStr(32),
				Secure: true,
			})
		}

		h.ServeHTTP(rw, r)
	})
}

func checkCSRFMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		csrfTok, err := getCookie(r, csrfTokenCookieName, "")

		if err != nil {
			internalServerError(rw, r, err)
			return

		} else if csrfTok == "" || r.Header.Get(csrfTokenHeaderName) != csrfTok {
			badRequest(rw, r, errors.New("invalid CSRF token"))
			return
		}

		h.ServeHTTP(rw, r)
	})
}
