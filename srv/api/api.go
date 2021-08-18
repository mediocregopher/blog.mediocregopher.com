// Package api implements the HTTP-based api for the mediocre-blog.
package api

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/mediocregopher/blog.mediocregopher.com/srv/cfg"
	"github.com/mediocregopher/blog.mediocregopher.com/srv/chat"
	"github.com/mediocregopher/blog.mediocregopher.com/srv/mailinglist"
	"github.com/mediocregopher/blog.mediocregopher.com/srv/pow"
	"github.com/mediocregopher/mediocre-go-lib/v2/mctx"
	"github.com/mediocregopher/mediocre-go-lib/v2/mlog"
)

// Params are used to instantiate a new API instance. All fields are required
// unless otherwise noted.
type Params struct {
	Logger      *mlog.Logger
	PowManager  pow.Manager
	MailingList mailinglist.MailingList
	GlobalRoom  chat.Room

	// ListenProto and ListenAddr are passed into net.Listen to create the
	// API's listener. Both "tcp" and "unix" protocols are explicitly
	// supported.
	ListenProto, ListenAddr string

	// StaticDir and StaticProxy are mutually exclusive.
	//
	// If StaticDir is set then that directory on the filesystem will be used to
	// serve the static site.
	//
	// Otherwise if StaticProxy is set all requests for the static site will be
	// reverse-proxied there.
	StaticDir   string
	StaticProxy *url.URL
}

// SetupCfg implement the cfg.Cfger interface.
func (p *Params) SetupCfg(cfg *cfg.Cfg) {

	cfg.StringVar(&p.ListenProto, "listen-proto", "tcp", "Protocol to listen for HTTP requests with")
	cfg.StringVar(&p.ListenAddr, "listen-addr", ":4000", "Address/path to listen for HTTP requests on")

	cfg.StringVar(&p.StaticDir, "static-dir", "", "Directory from which static files are served (mutually exclusive with -static-proxy-url)")
	staticProxyURLStr := cfg.String("static-proxy-url", "", "HTTP address from which static files are served (mutually exclusive with -static-dir)")

	cfg.OnInit(func(ctx context.Context) error {
		if *staticProxyURLStr != "" {
			var err error
			if p.StaticProxy, err = url.Parse(*staticProxyURLStr); err != nil {
				return fmt.Errorf("parsing -static-proxy-url: %w", err)
			}

		} else if p.StaticDir == "" {
			return errors.New("-static-dir or -static-proxy-url is required")
		}

		return nil
	})
}

// Annotate implements mctx.Annotator interface.
func (p *Params) Annotate(a mctx.Annotations) {
	a["listenProto"] = p.ListenProto
	a["listenAddr"] = p.ListenAddr

	if p.StaticProxy != nil {
		a["staticProxy"] = p.StaticProxy.String()
		return
	}

	a["staticDir"] = p.StaticDir
}

// API will listen on the port configured for it, and serve HTTP requests for
// the mediocre-blog.
type API interface {
	Shutdown(ctx context.Context) error
}

type api struct {
	params Params
	srv    *http.Server
}

// New initializes and returns a new API instance, including setting up all
// listening ports.
func New(params Params) (API, error) {

	l, err := net.Listen(params.ListenProto, params.ListenAddr)
	if err != nil {
		return nil, fmt.Errorf("creating listen socket: %w", err)
	}

	if params.ListenProto == "unix" {
		if err := os.Chmod(params.ListenAddr, 0777); err != nil {
			return nil, fmt.Errorf("chmod-ing unix socket: %w", err)
		}
	}

	a := &api{
		params: params,
	}

	a.srv = &http.Server{Handler: a.handler()}

	go func() {

		err := a.srv.Serve(l)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			ctx := mctx.Annotate(context.Background(), a.params)
			params.Logger.Fatal(ctx, fmt.Sprintf("%s: %v", "serving http server", err))
		}
	}()

	return a, nil
}

func (a *api) Shutdown(ctx context.Context) error {
	if err := a.srv.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}

func (a *api) handler() http.Handler {

	var staticHandler http.Handler
	if a.params.StaticDir != "" {
		staticHandler = http.FileServer(http.Dir(a.params.StaticDir))
	} else {
		staticHandler = httputil.NewSingleHostReverseProxy(a.params.StaticProxy)
	}

	// sugar
	requirePow := func(h http.Handler) http.Handler {
		return a.requirePowMiddleware(h)
	}

	mux := http.NewServeMux()

	mux.Handle("/", staticHandler)

	apiMux := http.NewServeMux()
	apiMux.Handle("/pow/challenge", a.newPowChallengeHandler())
	apiMux.Handle("/pow/check",
		requirePow(
			http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {}),
		),
	)

	apiMux.Handle("/mailinglist/subscribe", requirePow(a.mailingListSubscribeHandler()))
	apiMux.Handle("/mailinglist/finalize", a.mailingListFinalizeHandler())
	apiMux.Handle("/mailinglist/unsubscribe", a.mailingListUnsubscribeHandler())

	apiHandler := logMiddleware(a.params.Logger, apiMux)
	apiHandler = annotateMiddleware(apiHandler)
	apiHandler = addResponseHeaders(map[string]string{
		"Cache-Control": "no-store, max-age=0",
		"Pragma":        "no-cache",
		"Expires":       "0",
	}, apiHandler)

	mux.Handle("/api/", http.StripPrefix("/api", apiHandler))

	return mux
}
