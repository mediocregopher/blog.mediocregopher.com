package main

import (
	"net"
	"net/http"
	"time"

	"github.com/mediocregopher/mediocre-go-lib/v2/mctx"
	"github.com/mediocregopher/mediocre-go-lib/v2/mlog"
)

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
	statusCode int
}

func newLogResponseWriter(rw http.ResponseWriter) *logResponseWriter {
	return &logResponseWriter{
		ResponseWriter: rw,
		statusCode:     200,
	}
}

func (lrw *logResponseWriter) WriteHeader(statusCode int) {
	lrw.statusCode = statusCode
	lrw.ResponseWriter.WriteHeader(statusCode)
}

func logMiddleware(logger *mlog.Logger, h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		r = setRequestLogger(r, logger)

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
