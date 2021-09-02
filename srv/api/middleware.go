package api

import (
	"net"
	"net/http"
	"time"

	"github.com/mediocregopher/blog.mediocregopher.com/srv/api/apiutils"
	"github.com/mediocregopher/mediocre-go-lib/v2/mctx"
	"github.com/mediocregopher/mediocre-go-lib/v2/mlog"
)

func addResponseHeaders(headers map[string]string, h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		for k, v := range headers {
			rw.Header().Set(k, v)
		}
		h.ServeHTTP(rw, r)
	})
}

func annotateMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		type reqInfoKey string

		ip, _, _ := net.SplitHostPort(r.RemoteAddr)

		ctx := r.Context()
		ctx = mctx.Annotate(ctx,
			reqInfoKey("remote_ip"), ip,
			reqInfoKey("url"), r.URL,
			reqInfoKey("method"), r.Method,
		)

		r = r.WithContext(ctx)
		h.ServeHTTP(rw, r)
	})
}

type logResponseWriter struct {
	http.ResponseWriter
	http.Hijacker
	statusCode int
}

func newLogResponseWriter(rw http.ResponseWriter) *logResponseWriter {
	h, _ := rw.(http.Hijacker)
	return &logResponseWriter{
		ResponseWriter: rw,
		Hijacker:       h,
		statusCode:     200,
	}
}

func (lrw *logResponseWriter) WriteHeader(statusCode int) {
	lrw.statusCode = statusCode
	lrw.ResponseWriter.WriteHeader(statusCode)
}

func logMiddleware(logger *mlog.Logger, h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		r = apiutils.SetRequestLogger(r, logger)

		lrw := newLogResponseWriter(rw)

		started := time.Now()
		h.ServeHTTP(lrw, r)
		took := time.Since(started)

		type logCtxKey string

		ctx := r.Context()
		ctx = mctx.Annotate(ctx,
			logCtxKey("took"), took.String(),
			logCtxKey("response_code"), lrw.statusCode,
		)

		logger.Info(ctx, "handled HTTP request")
	})
}

func postOnlyMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		// we allow websockets to not be POSTs because, well, they can't be
		if r.Method == "POST" || r.Header.Get("Upgrade") == "websocket" {
			h.ServeHTTP(rw, r)
			return
		}

		apiutils.GetRequestLogger(r).WarnString(r.Context(), "method not allowed")
		rw.WriteHeader(405)
	})
}
