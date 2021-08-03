package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path"
	"strconv"
	"strings"

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

	logger := mlog.NewLogger(nil)

	publicURLStr := flag.String("public-url", "http://localhost:4000", "URL this service is accessible at")
	listenAddr := flag.String("listen-addr", ":4000", "Address to listen for HTTP requests on")
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
		logger.Fatal(context.Background(), "-static-dir or -static-proxy-url is required")
	case *powSecret == "":
		logger.Fatal(context.Background(), "-pow-secret is required")
	case *smtpAddr == "":
		logger.Fatal(context.Background(), "-ml-smtp-addr is required")
	case *smtpAuthStr == "":
		logger.Fatal(context.Background(), "-ml-smtp-auth is required")
	}

	publicURL, err := url.Parse(*publicURLStr)
	if err != nil {
		loggerFatalErr(context.Background(), logger, "parsing -public-url", err)
	}

	var staticProxyURL *url.URL
	if *staticProxyURLStr != "" {
		var err error
		if staticProxyURL, err = url.Parse(*staticProxyURLStr); err != nil {
			loggerFatalErr(context.Background(), logger, "parsing -static-proxy-url", err)
		}
	}

	powTargetUint, err := strconv.ParseUint(*powTargetStr, 0, 32)
	if err != nil {
		loggerFatalErr(context.Background(), logger, "parsing -pow-target", err)
	}
	powTarget := uint32(powTargetUint)

	smtpAuthParts := strings.SplitN(*smtpAuthStr, ":", 2)
	if len(smtpAuthParts) < 2 {
		logger.Fatal(context.Background(), "invalid -ml-smtp-auth")
	}
	smtpAuth := sasl.NewPlainClient("", smtpAuthParts[0], smtpAuthParts[1])
	smtpSendAs := smtpAuthParts[0]

	// initialization

	ctx := mctx.Annotate(context.Background(),
		"publicURL", publicURL.String(),
		"listenAddr", *listenAddr,
		"dataDir", *dataDir,
		"powTarget", fmt.Sprintf("%x", powTarget),
		"smtpAddr", *smtpAddr,
		"smtpSendAs", smtpSendAs,
	)

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

	mailer := mailinglist.NewMailer(mailinglist.MailerParams{
		SMTPAddr: *smtpAddr,
		SMTPAuth: smtpAuth,
		SendAs:   smtpSendAs,
	})

	mlStore, err := mailinglist.NewStore(path.Join(*dataDir, "mailinglist.sqlite3"))
	if err != nil {
		loggerFatalErr(ctx, logger, "initializing mailing list storage", err)
	}
	defer mlStore.Close()

	ml := mailinglist.New(mailinglist.Params{
		Store:          mlStore,
		Mailer:         mailer,
		Clock:          clock,
		FinalizeSubURL: path.Join(publicURL.String(), "/mailinglist/finalize.html"),
		UnsubURL:       path.Join(publicURL.String(), "/mailinglist/unsubscribe.html"),
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
	mux.Handle("/api/", http.StripPrefix("/api", apiHandler))

	// run

	logger.Info(ctx, "listening")

	// TODO graceful shutdown
	err = http.ListenAndServe(*listenAddr, mux)
	loggerFatalErr(ctx, logger, "listening", err)
}
