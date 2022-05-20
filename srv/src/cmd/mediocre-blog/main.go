package main

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"syscall"
	"time"

	cfgpkg "github.com/mediocregopher/blog.mediocregopher.com/srv/cfg"
	"github.com/mediocregopher/blog.mediocregopher.com/srv/chat"
	"github.com/mediocregopher/blog.mediocregopher.com/srv/http"
	"github.com/mediocregopher/blog.mediocregopher.com/srv/mailinglist"
	"github.com/mediocregopher/blog.mediocregopher.com/srv/post"
	"github.com/mediocregopher/blog.mediocregopher.com/srv/pow"
	"github.com/mediocregopher/mediocre-go-lib/v2/mctx"
	"github.com/mediocregopher/mediocre-go-lib/v2/mlog"
	"github.com/tilinna/clock"
)

func main() {

	ctx := context.Background()

	cfg := cfgpkg.NewBlogCfg(cfgpkg.Params{})

	var dataDir cfgpkg.DataDir
	dataDir.SetupCfg(cfg)
	defer dataDir.Close()
	ctx = mctx.WithAnnotator(ctx, &dataDir)

	var powMgrParams pow.ManagerParams
	powMgrParams.SetupCfg(cfg)
	ctx = mctx.WithAnnotator(ctx, &powMgrParams)

	var mailerParams mailinglist.MailerParams
	mailerParams.SetupCfg(cfg)
	ctx = mctx.WithAnnotator(ctx, &mailerParams)

	var mlParams mailinglist.Params
	mlParams.SetupCfg(cfg)
	ctx = mctx.WithAnnotator(ctx, &mlParams)

	var httpParams http.Params
	httpParams.SetupCfg(cfg)
	ctx = mctx.WithAnnotator(ctx, &httpParams)

	var radixClient cfgpkg.RadixClient
	radixClient.SetupCfg(cfg)
	defer radixClient.Close()
	ctx = mctx.WithAnnotator(ctx, &radixClient)

	chatGlobalRoomMaxMsgs := cfg.Int("chat-global-room-max-messages", 1000, "Maximum number of messages the global chat room can retain")
	chatUserIDCalcSecret := cfg.String("chat-user-id-calc-secret", "", "Secret to use when calculating user ids")

	pathPrefix := cfg.String("path-prefix", "", "Prefix which is optionally applied to all URL paths rendered by the blog")

	httpAuthUsersStr := cfg.String("http-auth-users", "{}", "JSON object with usernames as values and password hashes (produced by the hash-password binary) as values. Denotes users which are able to edit server-side data")

	// initialization
	err := cfg.Init(ctx)

	logger := mlog.NewLogger(nil)
	defer logger.Close()

	logger.Info(ctx, "process started")
	defer logger.Info(ctx, "process exiting")

	if err != nil {
		logger.Fatal(ctx, "initializing", err)
	}

	ctx = mctx.Annotate(ctx,
		"chatGlobalRoomMaxMsgs", *chatGlobalRoomMaxMsgs,
	)

	if *pathPrefix != "" {
		ctx = mctx.Annotate(ctx, "pathPrefix", *pathPrefix)
	}

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

	mlStore, err := mailinglist.NewStore(dataDir)
	if err != nil {
		logger.Fatal(ctx, "initializing mailing list storage", err)
	}
	defer mlStore.Close()

	mlParams.Store = mlStore
	mlParams.Mailer = mailer
	mlParams.Clock = clock

	ml := mailinglist.New(mlParams)

	chatGlobalRoom, err := chat.NewRoom(ctx, chat.RoomParams{
		Logger:      logger.WithNamespace("global-chat-room"),
		Redis:       radixClient.Client,
		ID:          "global",
		MaxMessages: *chatGlobalRoomMaxMsgs,
	})
	if err != nil {
		logger.Fatal(ctx, "initializing global chat room", err)
	}
	defer chatGlobalRoom.Close()

	chatUserIDCalc := chat.NewUserIDCalculator([]byte(*chatUserIDCalcSecret))

	postSQLDB, err := post.NewSQLDB(dataDir)
	if err != nil {
		logger.Fatal(ctx, "initializing sql db for post data", err)
	}
	defer postSQLDB.Close()

	postStore := post.NewStore(postSQLDB)
	postAssetStore := post.NewAssetStore(postSQLDB)

	var httpAuthUsers map[string]string
	if err := json.Unmarshal([]byte(*httpAuthUsersStr), &httpAuthUsers); err != nil {
		logger.Fatal(ctx, "unmarshaling -http-auth-users", err)
	}

	httpParams.Logger = logger.WithNamespace("http")
	httpParams.PowManager = powMgr
	httpParams.PathPrefix = *pathPrefix
	httpParams.PostStore = postStore
	httpParams.PostAssetStore = postAssetStore
	httpParams.MailingList = ml
	httpParams.GlobalRoom = chatGlobalRoom
	httpParams.UserIDCalculator = chatUserIDCalc
	httpParams.AuthUsers = httpAuthUsers

	logger.Info(ctx, "listening")
	httpAPI, err := http.New(httpParams)
	if err != nil {
		logger.Fatal(ctx, "initializing http api", err)
	}
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		if err := httpAPI.Shutdown(shutdownCtx); err != nil {
			logger.Fatal(ctx, "shutting down http api", err)
		}
	}()

	// wait

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	// let the defers begin
}
