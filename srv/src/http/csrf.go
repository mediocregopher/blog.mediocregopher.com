package http

import (
	"errors"
	"net/http"

	"github.com/mediocregopher/blog.mediocregopher.com/srv/http/apiutil"
)

const (
	csrfTokenCookieName = "csrf_token"
	csrfTokenHeaderName = "X-CSRF-Token"
	csrfTokenFormName   = "csrfToken"
)

func setCSRFMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		csrfTok, err := apiutil.GetCookie(r, csrfTokenCookieName, "")

		if err != nil {
			apiutil.InternalServerError(rw, r, err)
			return

		} else if csrfTok == "" {
			http.SetCookie(rw, &http.Cookie{
				Name:   csrfTokenCookieName,
				Value:  apiutil.RandStr(32),
				Secure: true,
			})
		}

		h.ServeHTTP(rw, r)
	})
}

func checkCSRFMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		csrfTok, err := apiutil.GetCookie(r, csrfTokenCookieName, "")

		if err != nil {
			apiutil.InternalServerError(rw, r, err)
			return
		}

		givenCSRFTok := r.Header.Get(csrfTokenHeaderName)
		if givenCSRFTok == "" {
			givenCSRFTok = r.FormValue(csrfTokenFormName)
		}

		if csrfTok == "" || givenCSRFTok != csrfTok {
			apiutil.BadRequest(rw, r, errors.New("invalid CSRF token"))
			return
		}

		h.ServeHTTP(rw, r)
	})
}
