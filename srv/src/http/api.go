// Package api implements the HTTP-based api for the mediocre-blog.
package http

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net"
	"net/http"
	"os"

	"github.com/mediocregopher/blog.mediocregopher.com/srv/cfg"
	"github.com/mediocregopher/blog.mediocregopher.com/srv/chat"
	"github.com/mediocregopher/blog.mediocregopher.com/srv/http/apiutil"
	"github.com/mediocregopher/blog.mediocregopher.com/srv/mailinglist"
	"github.com/mediocregopher/blog.mediocregopher.com/srv/post"
	"github.com/mediocregopher/blog.mediocregopher.com/srv/pow"
	"github.com/mediocregopher/mediocre-go-lib/v2/mctx"
	"github.com/mediocregopher/mediocre-go-lib/v2/mlog"
)

//go:embed static
var staticFS embed.FS

// Params are used to instantiate a new API instance. All fields are required
// unless otherwise noted.
type Params struct {
	Logger     *mlog.Logger
	PowManager pow.Manager

	// PathPrefix, if given, will be prefixed to all url paths which are
	// rendered by the API's templating system.
	PathPrefix string

	PostStore      post.Store
	PostAssetStore post.AssetStore

	MailingList mailinglist.MailingList

	GlobalRoom       chat.Room
	UserIDCalculator *chat.UserIDCalculator

	// ListenProto and ListenAddr are passed into net.Listen to create the
	// API's listener. Both "tcp" and "unix" protocols are explicitly
	// supported.
	ListenProto, ListenAddr string

	// AuthUsers keys are usernames which are allowed to edit server-side data,
	// and the values are the password hash which accompanies those users. The
	// password hash must have been produced by NewPasswordHash.
	AuthUsers map[string]string
}

// SetupCfg implement the cfg.Cfger interface.
func (p *Params) SetupCfg(cfg *cfg.Cfg) {
	cfg.StringVar(&p.ListenProto, "listen-proto", "tcp", "Protocol to listen for HTTP requests with")
	cfg.StringVar(&p.ListenAddr, "listen-addr", ":4000", "Address/path to listen for HTTP requests on")

	httpAuthUsersStr := cfg.String("http-auth-users", "{}", "JSON object with usernames as values and password hashes (produced by the hash-password binary) as values. Denotes users which are able to edit server-side data")

	cfg.OnInit(func(context.Context) error {
		if err := json.Unmarshal([]byte(*httpAuthUsersStr), &p.AuthUsers); err != nil {
			return fmt.Errorf("unmarshaling -http-auth-users: %w", err)
		}
		return nil
	})
}

// Annotate implements mctx.Annotator interface.
func (p *Params) Annotate(a mctx.Annotations) {
	a["listenProto"] = p.ListenProto
	a["listenAddr"] = p.ListenAddr
}

// API will listen on the port configured for it, and serve HTTP requests for
// the mediocre-blog.
type API interface {
	Shutdown(ctx context.Context) error
}

type api struct {
	params Params
	srv    *http.Server

	redirectTpl *template.Template
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

	a.redirectTpl = a.mustParseTpl("redirect.html")

	a.srv = &http.Server{Handler: a.handler()}

	go func() {

		err := a.srv.Serve(l)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			ctx := mctx.Annotate(context.Background(), a.params)
			params.Logger.Fatal(ctx, "serving http server", err)
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

	requirePow := func(h http.Handler) http.Handler {
		return a.requirePowMiddleware(h)
	}

	formMiddleware := func(h http.Handler) http.Handler {
		h = checkCSRFMiddleware(h)
		h = disallowGetMiddleware(h)
		h = logReqMiddleware(h)
		h = addResponseHeaders(map[string]string{
			"Cache-Control": "no-store, max-age=0",
			"Pragma":        "no-cache",
			"Expires":       "0",
		}, h)
		return h
	}

	auther := NewAuther(a.params.AuthUsers)

	mux := http.NewServeMux()

	{
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

		apiMux.Handle("/chat/global/", http.StripPrefix("/chat/global", newChatHandler(
			a.params.GlobalRoom,
			a.params.UserIDCalculator,
			a.requirePowMiddleware,
		)))

		mux.Handle("/api/", http.StripPrefix("/api", formMiddleware(apiMux)))
	}

	mux.Handle("/posts/", http.StripPrefix("/posts",
		apiutil.MethodMux(map[string]http.Handler{
			"GET":  a.renderPostHandler(),
			"EDIT": a.editPostHandler(),
			"POST": authMiddleware(auther,
				formMiddleware(a.postPostHandler()),
			),
			"DELETE": authMiddleware(auther,
				formMiddleware(a.deletePostHandler()),
			),
			"PREVIEW": authMiddleware(auther,
				formMiddleware(a.previewPostHandler()),
			),
		}),
	))

	mux.Handle("/assets/", http.StripPrefix("/assets",
		apiutil.MethodMux(map[string]http.Handler{
			"GET": a.getPostAssetHandler(),
			"POST": authMiddleware(auther,
				formMiddleware(a.postPostAssetHandler()),
			),
			"DELETE": authMiddleware(auther,
				formMiddleware(a.deletePostAssetHandler()),
			),
		}),
	))

	mux.Handle("/static/", http.FileServer(http.FS(staticFS)))
	mux.Handle("/follow", a.renderDumbTplHandler("follow.html"))
	mux.Handle("/", a.renderIndexHandler())

	var globalHandler http.Handler = mux
	globalHandler = setCSRFMiddleware(globalHandler)
	globalHandler = setLoggerMiddleware(a.params.Logger, globalHandler)

	return globalHandler
}
