package api

import (
	"errors"
	"net/http"

	"github.com/mediocregopher/blog.mediocregopher.com/srv/api/apiutils"
)

const (
	csrfTokenCookieName = "csrf_token"
	csrfTokenHeaderName = "X-CSRF-Token"
)

func setCSRFMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		csrfTok, err := apiutils.GetCookie(r, csrfTokenCookieName, "")

		if err != nil {
			apiutils.InternalServerError(rw, r, err)
			return

		} else if csrfTok == "" {
			http.SetCookie(rw, &http.Cookie{
				Name:   csrfTokenCookieName,
				Value:  apiutils.RandStr(32),
				Secure: true,
			})
		}

		h.ServeHTTP(rw, r)
	})
}

func checkCSRFMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		csrfTok, err := apiutils.GetCookie(r, csrfTokenCookieName, "")

		if err != nil {
			apiutils.InternalServerError(rw, r, err)
			return
		}

		givenCSRFTok := r.Header.Get(csrfTokenHeaderName)
		if givenCSRFTok == "" {
			givenCSRFTok = r.FormValue("csrfToken")
		}

		if csrfTok == "" || givenCSRFTok != csrfTok {
			apiutils.BadRequest(rw, r, errors.New("invalid CSRF token"))
			return
		}

		h.ServeHTTP(rw, r)
	})
}
