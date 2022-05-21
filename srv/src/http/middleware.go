package http

import (
	"bytes"
	"net"
	"net/http"
	"path/filepath"
	"sync"
	"time"

	lru "github.com/hashicorp/golang-lru"
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

type wrappedResponseWriter struct {
	http.ResponseWriter
	http.Hijacker
	statusCode int
}

func newWrappedResponseWriter(rw http.ResponseWriter) *wrappedResponseWriter {
	h, _ := rw.(http.Hijacker)
	return &wrappedResponseWriter{
		ResponseWriter: rw,
		Hijacker:       h,
		statusCode:     200,
	}
}

func (rw *wrappedResponseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func logReqMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		wrw := newWrappedResponseWriter(rw)

		started := time.Now()
		h.ServeHTTP(wrw, r)
		took := time.Since(started)

		type logCtxKey string

		ctx := r.Context()
		ctx = mctx.Annotate(ctx,
			logCtxKey("took"), took.String(),
			logCtxKey("response_code"), wrw.statusCode,
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

type cacheResponseWriter struct {
	*wrappedResponseWriter
	buf *bytes.Buffer
}

func newCacheResponseWriter(rw http.ResponseWriter) *cacheResponseWriter {
	return &cacheResponseWriter{
		wrappedResponseWriter: newWrappedResponseWriter(rw),
		buf:                   new(bytes.Buffer),
	}
}

func (rw *cacheResponseWriter) Write(b []byte) (int, error) {
	if _, err := rw.buf.Write(b); err != nil {
		panic(err)
	}
	return rw.wrappedResponseWriter.Write(b)
}

func cacheMiddleware(cache *lru.Cache) middleware {

	type entry struct {
		body      []byte
		createdAt time.Time
	}

	pool := sync.Pool{
		New: func() interface{} { return new(bytes.Reader) },
	}

	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

			id := r.URL.RequestURI()

			if val, ok := cache.Get(id); ok {

				entry := val.(entry)

				reader := pool.Get().(*bytes.Reader)
				defer pool.Put(reader)

				reader.Reset(entry.body)

				http.ServeContent(
					rw, r, filepath.Base(r.URL.Path), entry.createdAt, reader,
				)
				return
			}

			cacheRW := newCacheResponseWriter(rw)
			h.ServeHTTP(cacheRW, r)

			if cacheRW.statusCode == 200 {
				cache.Add(id, entry{
					body:      cacheRW.buf.Bytes(),
					createdAt: time.Now(),
				})
			}
		})
	}
}

func purgeCacheOnOKMiddleware(cache *lru.Cache) middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

			wrw := newWrappedResponseWriter(rw)
			h.ServeHTTP(wrw, r)

			if wrw.statusCode == 200 {
				apiutil.GetRequestLogger(r).Info(r.Context(), "purging cache!")
				cache.Purge()
			}
		})
	}
}
