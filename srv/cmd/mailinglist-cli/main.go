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

	subCmd := cfg.SubCmd()
	ctx = mctx.Annotate(ctx, "subCmd", subCmd)

	switch subCmd {
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

	case "publish":

		title := cfg.String("title", "", "Title of the post which was published")
		url := cfg.String("url", "", "URL of the post which was published")
		cfg.Init(ctx)

		if *title == "" {
			logger.Fatal(ctx, "-title is required")

		} else if *url == "" {
			logger.Fatal(ctx, "-url is required")
		}

		err := ml.Publish(*title, *url)
		if err != nil {
			loggerFatalErr(ctx, logger, "publishing", err)
		}

	default:
		logger.Fatal(ctx, "invalid sub-command, must be list|publish")
	}
}
