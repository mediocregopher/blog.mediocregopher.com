package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"path"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/emersion/go-sasl"
	"github.com/mediocregopher/blog.mediocregopher.com/srv/mailinglist"
	"github.com/mediocregopher/blog.mediocregopher.com/srv/pow"
	"github.com/mediocregopher/mediocre-go-lib/v2/mctx"
	"github.com/mediocregopher/mediocre-go-lib/v2/mlog"
	"github.com/tilinna/clock"
)

func loggerFatalErr(ctx context.Context, logger *mlog.Logger, descr string, err error) {
	logger.Fatal(ctx, fmt.Sprintf("%s: %v", descr, err))
}

func main() {

	ctx := context.Background()

	logger := mlog.NewLogger(nil)
	defer logger.Close()

	logger.Info(ctx, "process started")
	defer logger.Info(ctx, "process exiting")

	publicURLStr := flag.String("public-url", "http://localhost:4000", "URL this service is accessible at")
	listenProto := flag.String("listen-proto", "tcp", "Protocol to listen for HTTP requests with")
	listenAddr := flag.String("listen-addr", ":4000", "Address/path to listen for HTTP requests on")
	dataDir := flag.String("data-dir", ".", "Directory to use for long term storage")

	staticDir := flag.String("static-dir", "", "Directory from which static files are served (mutually exclusive with -static-proxy-url)")
	staticProxyURLStr := flag.String("static-proxy-url", "", "HTTP address from which static files are served (mutually exclusive with -static-dir)")

	powTargetStr := flag.String("pow-target", "0x0000FFFF", "Proof-of-work target, lower is more difficult")
	powSecret := flag.String("pow-secret", "", "Secret used to sign proof-of-work challenge seeds")

	smtpAddr := flag.String("ml-smtp-addr", "", "Address of SMTP server to use for sending emails for the mailing list")
	smtpAuthStr := flag.String("ml-smtp-auth", "", "user:pass to use when authenticating with the mailing list SMTP server. The given user will also be used as the From address.")

	// parse config

	flag.Parse()

	switch {
	case *staticDir == "" && *staticProxyURLStr == "":
		logger.Fatal(ctx, "-static-dir or -static-proxy-url is required")
	case *powSecret == "":
		logger.Fatal(ctx, "-pow-secret is required")
	}

	publicURL, err := url.Parse(*publicURLStr)
	if err != nil {
		loggerFatalErr(ctx, logger, "parsing -public-url", err)
	}

	var staticProxyURL *url.URL
	if *staticProxyURLStr != "" {
		var err error
		if staticProxyURL, err = url.Parse(*staticProxyURLStr); err != nil {
			loggerFatalErr(ctx, logger, "parsing -static-proxy-url", err)
		}
	}

	powTargetUint, err := strconv.ParseUint(*powTargetStr, 0, 32)
	if err != nil {
		loggerFatalErr(ctx, logger, "parsing -pow-target", err)
	}
	powTarget := uint32(powTargetUint)

	var mailerCfg mailinglist.MailerParams

	if *smtpAddr != "" {
		mailerCfg.SMTPAddr = *smtpAddr
		smtpAuthParts := strings.SplitN(*smtpAuthStr, ":", 2)
		if len(smtpAuthParts) < 2 {
			logger.Fatal(ctx, "invalid -ml-smtp-auth")
		}
		mailerCfg.SMTPAuth = sasl.NewPlainClient("", smtpAuthParts[0], smtpAuthParts[1])
		mailerCfg.SendAs = smtpAuthParts[0]

		ctx = mctx.Annotate(ctx,
			"smtpAddr", mailerCfg.SMTPAddr,
			"smtpSendAs", mailerCfg.SendAs,
		)
	}

	ctx = mctx.Annotate(ctx,
		"publicURL", publicURL.String(),
		"listenProto", *listenProto,
		"listenAddr", *listenAddr,
		"dataDir", *dataDir,
		"powTarget", fmt.Sprintf("%x", powTarget),
	)

	// initialization

	if *staticDir != "" {
		ctx = mctx.Annotate(ctx, "staticDir", *staticDir)
	} else {
		ctx = mctx.Annotate(ctx, "staticProxyURL", *staticProxyURLStr)
	}

	clock := clock.Realtime()

	powStore := pow.NewMemoryStore(clock)
	defer powStore.Close()

	powMgr := pow.NewManager(pow.ManagerParams{
		Clock:  clock,
		Store:  powStore,
		Secret: []byte(*powSecret),
		Target: powTarget,
	})

	// sugar
	requirePow := func(h http.Handler) http.Handler { return requirePowMiddleware(powMgr, h) }

	var mailer mailinglist.Mailer
	if *smtpAddr == "" {
		logger.Info(ctx, "-smtp-addr not given, using NullMailer")
		mailer = mailinglist.NullMailer
	} else {
		mailer = mailinglist.NewMailer(mailerCfg)
	}

	mlStore, err := mailinglist.NewStore(path.Join(*dataDir, "mailinglist.sqlite3"))
	if err != nil {
		loggerFatalErr(ctx, logger, "initializing mailing list storage", err)
	}
	defer mlStore.Close()

	ml := mailinglist.New(mailinglist.Params{
		Store:          mlStore,
		Mailer:         mailer,
		Clock:          clock,
		FinalizeSubURL: publicURL.String() + "/mailinglist/finalize.html",
		UnsubURL:       publicURL.String() + "/mailinglist/unsubscribe.html",
	})

	mux := http.NewServeMux()

	var staticHandler http.Handler
	if *staticDir != "" {
		staticHandler = http.FileServer(http.Dir(*staticDir))
	} else {
		staticHandler = httputil.NewSingleHostReverseProxy(staticProxyURL)
	}

	mux.Handle("/", staticHandler)

	apiMux := http.NewServeMux()
	apiMux.Handle("/pow/challenge", newPowChallengeHandler(powMgr))
	apiMux.Handle("/pow/check",
		requirePow(
			http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {}),
		),
	)

	apiMux.Handle("/mailinglist/subscribe", requirePow(mailingListSubscribeHandler(ml)))
	apiMux.Handle("/mailinglist/finalize", mailingListFinalizeHandler(ml))
	apiMux.Handle("/mailinglist/unsubscribe", mailingListUnsubscribeHandler(ml))

	apiHandler := logMiddleware(logger.WithNamespace("api"), apiMux)
	apiHandler = annotateMiddleware(apiHandler)
	apiHandler = addResponseHeaders(map[string]string{
		"Cache-Control": "no-store, max-age=0",
		"Pragma":        "no-cache",
		"Expires":       "0",
	}, apiHandler)

	mux.Handle("/api/", http.StripPrefix("/api", apiHandler))

	// run

	logger.Info(ctx, "listening")

	l, err := net.Listen(*listenProto, *listenAddr)
	if err != nil {
		loggerFatalErr(ctx, logger, "creating listen socket", err)
	}

	if *listenProto == "unix" {
		if err := os.Chmod(*listenAddr, 0777); err != nil {
			loggerFatalErr(ctx, logger, "chmod-ing unix socket", err)
		}
	}

	srv := &http.Server{Handler: mux}
	go func() {
		if err := srv.Serve(l); err != nil && !errors.Is(err, http.ErrServerClosed) {
			loggerFatalErr(ctx, logger, "serving http server", err)
		}
	}()

	defer func() {
		closeCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		logger.Info(ctx, "beginning graceful shutdown of http server")

		if err := srv.Shutdown(closeCtx); err != nil {
			loggerFatalErr(ctx, logger, "gracefully shutting down http server", err)
		}
	}()

	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	// let the defers begin
}
