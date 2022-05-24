package http

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/mediocregopher/blog.mediocregopher.com/srv/http/apiutil"
)

func (a *api) checkCSRFMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		refererURL, err := url.Parse(r.Referer())
		if err != nil {
			apiutil.BadRequest(rw, r, errors.New("invalid Referer"))
			return
		}

		if refererURL.Scheme != a.params.PublicURL.Scheme ||
			refererURL.Host != a.params.PublicURL.Host {
			apiutil.BadRequest(rw, r, errors.New("invalid Referer"))
			return
		}

		h.ServeHTTP(rw, r)
	})
}
