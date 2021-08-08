package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path"

	"github.com/mediocregopher/blog.mediocregopher.com/srv/cfg"
	"github.com/mediocregopher/blog.mediocregopher.com/srv/mailinglist"
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

	var mailerParams mailinglist.MailerParams
	mailerParams.SetupCfg(cfg)
	ctx = mctx.WithAnnotator(ctx, &mailerParams)

	var mlParams mailinglist.Params
	mlParams.SetupCfg(cfg)
	ctx = mctx.WithAnnotator(ctx, &mlParams)

	// initialization
	err := cfg.Init(ctx)

	logger := mlog.NewLogger(nil)
	defer logger.Close()

	logger.Info(ctx, "process started")
	defer logger.Info(ctx, "process exiting")

	if err != nil {
		loggerFatalErr(ctx, logger, "initializing", err)
	}

	clock := clock.Realtime()

	var mailer mailinglist.Mailer
	if mailerParams.SMTPAddr == "" {
		logger.Info(ctx, "-smtp-addr not given, using NullMailer")
		mailer = mailinglist.NullMailer
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
	_ = ml

	args := cfg.Args()
	if len(args) == 0 {
		args = append(args, "")
	}

	action, args := args[0], args[1:]

	switch action {
	case "list":
		for it := mlStore.GetAll(); ; {
			email, err := it()
			if errors.Is(err, io.EOF) {
				break
			} else if err != nil {
				loggerFatalErr(ctx, logger, "retrieving next email", err)
			}

			ctx := mctx.Annotate(context.Background(),
				"email", email.Email,
				"createdAt", email.CreatedAt,
				"verifiedAt", email.VerifiedAt,
			)

			logger.Info(ctx, "next")
		}

	default:
		logger.Fatal(ctx, "invalid action")
	}
}
