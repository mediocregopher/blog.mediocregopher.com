package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	"github.com/mediocregopher/blog.mediocregopher.com/srv/api"
	"github.com/mediocregopher/blog.mediocregopher.com/srv/cfg"
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
	cfg := cfg.New()

	dataDir := cfg.String("data-dir", ".", "Directory to use for long term storage")

	var powMgrParams pow.ManagerParams
	powMgrParams.SetupCfg(cfg)
	ctx = mctx.WithAnnotator(ctx, &powMgrParams)

	var mailerParams mailinglist.MailerParams
	mailerParams.SetupCfg(cfg)
	ctx = mctx.WithAnnotator(ctx, &mailerParams)

	var mlParams mailinglist.Params
	mlParams.SetupCfg(cfg)
	ctx = mctx.WithAnnotator(ctx, &mlParams)

	var apiParams api.Params
	apiParams.SetupCfg(cfg)
	ctx = mctx.WithAnnotator(ctx, &apiParams)

	// initialization
	err := cfg.Init(ctx)

	logger := mlog.NewLogger(nil)
	defer logger.Close()

	logger.Info(ctx, "process started")
	defer logger.Info(ctx, "process exiting")

	if err != nil {
		loggerFatalErr(ctx, logger, "initializing", err)
	}

	ctx = mctx.Annotate(ctx,
		"dataDir", *dataDir,
	)

	clock := clock.Realtime()

	powStore := pow.NewMemoryStore(clock)
	defer powStore.Close()

	powMgrParams.Store = powStore
	powMgrParams.Clock = clock

	powMgr := pow.NewManager(powMgrParams)

	var mailer mailinglist.Mailer
	if mailerParams.SMTPAddr == "" {
		logger.Info(ctx, "-smtp-addr not given, using a fake Mailer")
		mailer = mailinglist.NewLogMailer(logger.WithNamespace("fake-mailer"))
	} else {
		mailer = mailinglist.NewMailer(mailerParams)
	}

	mlStore, err := mailinglist.NewStore(path.Join(*dataDir, "mailinglist.sqlite3"))
	if err != nil {
		loggerFatalErr(ctx, logger, "initializing mailing list storage", err)
	}
	defer mlStore.Close()

	mlParams.Store = mlStore
	mlParams.Mailer = mailer
	mlParams.Clock = clock

	ml := mailinglist.New(mlParams)

	apiParams.Logger = logger.WithNamespace("api")
	apiParams.PowManager = powMgr
	apiParams.MailingList = ml

	logger.Info(ctx, "listening")
	a, err := api.New(apiParams)
	if err != nil {
		loggerFatalErr(ctx, logger, "initializing api", err)
	}
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		if err := a.Shutdown(shutdownCtx); err != nil {
			loggerFatalErr(ctx, logger, "shutting down api", err)
		}
	}()

	// wait

	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	// let the defers begin
}
