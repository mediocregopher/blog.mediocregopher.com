package http

import (
	"net"
	"net/http"
	"time"

	"github.com/mediocregopher/blog.mediocregopher.com/srv/http/apiutil"
	"github.com/mediocregopher/mediocre-go-lib/v2/mctx"
	"github.com/mediocregopher/mediocre-go-lib/v2/mlog"
)

type middleware func(http.Handler) http.Handler

func applyMiddlewares(h http.Handler, middlewares ...middleware) http.Handler {
	for _, m := range middlewares {
		h = m(h)
	}
	return h
}

func addResponseHeadersMiddleware(headers map[string]string) middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			for k, v := range headers {
				rw.Header().Set(k, v)
			}
			h.ServeHTTP(rw, r)
		})
	}
}

func setLoggerMiddleware(logger *mlog.Logger) middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

			type logCtxKey string

			ip, _, _ := net.SplitHostPort(r.RemoteAddr)

			ctx := r.Context()
			ctx = mctx.Annotate(ctx,
				logCtxKey("remote_ip"), ip,
				logCtxKey("url"), r.URL,
				logCtxKey("method"), r.Method,
			)

			r = r.WithContext(ctx)
			r = apiutil.SetRequestLogger(r, logger)
			h.ServeHTTP(rw, r)
		})
	}
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

func logReqMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

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

		apiutil.GetRequestLogger(r).Info(ctx, "handled HTTP request")
	})
}

func disallowGetMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		// we allow websockets to be GETs because, well, they must be
		if r.Method != "GET" || r.Header.Get("Upgrade") == "websocket" {
			h.ServeHTTP(rw, r)
			return
		}

		apiutil.GetRequestLogger(r).WarnString(r.Context(), "method not allowed")
		rw.WriteHeader(405)
	})
}
