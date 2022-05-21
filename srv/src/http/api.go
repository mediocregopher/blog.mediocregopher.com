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
	"net/url"
	"os"
	"strings"
	"time"

	lru "github.com/hashicorp/golang-lru"
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

	PostStore      post.Store
	PostAssetStore post.AssetStore

	MailingList mailinglist.MailingList

	GlobalRoom       chat.Room
	UserIDCalculator *chat.UserIDCalculator

	// PublicURL is the base URL which site visitors can navigate to.
	PublicURL *url.URL

	// ListenProto and ListenAddr are passed into net.Listen to create the
	// API's listener. Both "tcp" and "unix" protocols are explicitly
	// supported.
	ListenProto, ListenAddr string

	// AuthUsers keys are usernames which are allowed to edit server-side data,
	// and the values are the password hash which accompanies those users. The
	// password hash must have been produced by NewPasswordHash.
	AuthUsers map[string]string

	// AuthRatelimit indicates how much time must pass between subsequent auth
	// attempts.
	AuthRatelimit time.Duration
}

// SetupCfg implement the cfg.Cfger interface.
func (p *Params) SetupCfg(cfg *cfg.Cfg) {

	publicURLStr := cfg.String("http-public-url", "http://localhost:4000", "URL this service is accessible at")

	cfg.StringVar(&p.ListenProto, "http-listen-proto", "tcp", "Protocol to listen for HTTP requests with")
	cfg.StringVar(&p.ListenAddr, "http-listen-addr", ":4000", "Address/path to listen for HTTP requests on")

	httpAuthUsersStr := cfg.String("http-auth-users", "{}", "JSON object with usernames as values and password hashes (produced by the hash-password binary) as values. Denotes users which are able to edit server-side data")

	httpAuthRatelimitStr := cfg.String("http-auth-ratelimit", "5s", "Minimum duration which must be waited between subsequent auth attempts")

	cfg.OnInit(func(context.Context) error {

		err := json.Unmarshal([]byte(*httpAuthUsersStr), &p.AuthUsers)

		if err != nil {
			return fmt.Errorf("unmarshaling -http-auth-users: %w", err)
		}

		if p.AuthRatelimit, err = time.ParseDuration(*httpAuthRatelimitStr); err != nil {
			return fmt.Errorf("unmarshaling -http-auth-ratelimit: %w", err)
		}

		*publicURLStr = strings.TrimSuffix(*publicURLStr, "/")
		if p.PublicURL, err = url.Parse(*publicURLStr); err != nil {
			return fmt.Errorf("parsing -http-public-url: %w", err)
		}

		return nil
	})
}

// Annotate implements mctx.Annotator interface.
func (p *Params) Annotate(a mctx.Annotations) {
	a["httpPublicURL"] = p.PublicURL
	a["httpListenProto"] = p.ListenProto
	a["httpListenAddr"] = p.ListenAddr
	a["httpAuthRatelimit"] = p.AuthRatelimit
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
	auther      Auther
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
		auther: NewAuther(params.AuthUsers, params.AuthRatelimit),
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
	defer a.auther.Close()
	if err := a.srv.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}

func (a *api) apiHandler() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/pow/challenge", a.newPowChallengeHandler())
	mux.Handle("/pow/check",
		a.requirePowMiddleware(
			http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {}),
		),
	)

	mux.Handle("/mailinglist/subscribe", a.requirePowMiddleware(a.mailingListSubscribeHandler()))
	mux.Handle("/mailinglist/finalize", a.mailingListFinalizeHandler())
	mux.Handle("/mailinglist/unsubscribe", a.mailingListUnsubscribeHandler())

	mux.Handle("/chat/global/", http.StripPrefix("/chat/global", newChatHandler(
		a.params.GlobalRoom,
		a.params.UserIDCalculator,
		a.requirePowMiddleware,
	)))

	// disallowGetMiddleware is used rather than a MethodMux because it has an
	// exception for websockets, which is needed for chat.
	return disallowGetMiddleware(mux)

}

func (a *api) blogHandler() http.Handler {

	cache, err := lru.New(5000)

	// instantiating the lru cache can't realistically fail
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()

	mux.Handle("/posts/", http.StripPrefix("/posts",
		apiutil.MethodMux(map[string]http.Handler{
			"GET":     a.renderPostHandler(),
			"POST":    a.postPostHandler(),
			"DELETE":  a.deletePostHandler(),
			"PREVIEW": a.previewPostHandler(),
		}),
	))

	mux.Handle("/assets/", http.StripPrefix("/assets",
		apiutil.MethodMux(map[string]http.Handler{
			"GET":    a.getPostAssetHandler(),
			"POST":   a.postPostAssetHandler(),
			"DELETE": a.deletePostAssetHandler(),
		}),
	))

	mux.Handle("/static/", http.FileServer(http.FS(staticFS)))
	mux.Handle("/follow", a.renderDumbTplHandler("follow.html"))
	mux.Handle("/mailinglist/unsubscribe", a.renderDumbTplHandler("unsubscribe.html"))
	mux.Handle("/mailinglist/finalize", a.renderDumbTplHandler("finalize.html"))
	mux.Handle("/feed.xml", a.renderFeedHandler())
	mux.Handle("/", a.renderIndexHandler())

	h := apiutil.MethodMux(map[string]http.Handler{
		"GET": applyMiddlewares(
			mux,
			logReqMiddleware, // only log GETs on cache miss
			cacheMiddleware(cache),
		),
		"*": applyMiddlewares(
			mux,
			purgeCacheOnOKMiddleware(cache),
			authMiddleware(a.auther),
		),
	})

	return h
}

func (a *api) handler() http.Handler {

	mux := http.NewServeMux()

	mux.Handle("/api/", http.StripPrefix("/api", a.apiHandler()))
	mux.Handle("/", a.blogHandler())

	h := apiutil.MethodMux(map[string]http.Handler{
		"GET": applyMiddlewares(
			mux,
			setCSRFMiddleware,
		),
		"*": applyMiddlewares(
			mux,
			checkCSRFMiddleware,
			addResponseHeadersMiddleware(map[string]string{
				"Cache-Control": "no-store, max-age=0",
				"Pragma":        "no-cache",
				"Expires":       "0",
			}),
			logReqMiddleware,
		),
	})

	h = setLoggerMiddleware(a.params.Logger)(h)

	return h
}
