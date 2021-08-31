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
	"github.com/mediocregopher/blog.mediocregopher.com/srv/chat"
	"github.com/mediocregopher/blog.mediocregopher.com/srv/mailinglist"
	"github.com/mediocregopher/blog.mediocregopher.com/srv/pow"
	"github.com/mediocregopher/mediocre-go-lib/v2/mctx"
	"github.com/mediocregopher/mediocre-go-lib/v2/mlog"
	"github.com/mediocregopher/radix/v4"
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

	redisProto := cfg.String("redis-proto", "tcp", "Network protocol to connect to redis over, can be tcp or unix")
	redisAddr := cfg.String("redis-addr", "127.0.0.1:6379", "Address redis is expected to listen on")
	redisPoolSize := cfg.Int("redis-pool-size", 5, "Number of connections in the redis pool to keep")

	chatGlobalRoomMaxMsgs := cfg.Int("chat-global-room-max-messages", 1000, "Maximum number of messages the global chat room can retain")
	chatUserIDCalcSecret := cfg.String("chat-user-id-calc-secret", "", "Secret to use when calculating user ids")

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
		"redisProto", *redisProto,
		"redisAddr", *redisAddr,
		"redisPoolSize", *redisPoolSize,
		"chatGlobalRoomMaxMsgs", *chatGlobalRoomMaxMsgs,
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

	redis, err := (radix.PoolConfig{
		Size: *redisPoolSize,
	}).New(
		ctx, *redisProto, *redisAddr,
	)

	if err != nil {
		loggerFatalErr(ctx, logger, "initializing redis pool", err)
	}
	defer redis.Close()

	chatGlobalRoom, err := chat.NewRoom(ctx, chat.RoomParams{
		Logger:      logger.WithNamespace("global-chat-room"),
		Redis:       redis,
		ID:          "global",
		MaxMessages: *chatGlobalRoomMaxMsgs,
	})
	if err != nil {
		loggerFatalErr(ctx, logger, "initializing global chat room", err)
	}
	defer chatGlobalRoom.Close()

	chatUserIDCalc := chat.NewUserIDCalculator([]byte(*chatUserIDCalcSecret))

	apiParams.Logger = logger.WithNamespace("api")
	apiParams.PowManager = powMgr
	apiParams.MailingList = ml
	apiParams.GlobalRoom = chatGlobalRoom
	apiParams.UserIDCalculator = chatUserIDCalc

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
