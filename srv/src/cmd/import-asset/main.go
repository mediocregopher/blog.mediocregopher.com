package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	cfgpkg "github.com/mediocregopher/blog.mediocregopher.com/srv/cfg"
	"github.com/mediocregopher/blog.mediocregopher.com/srv/post"
	"github.com/mediocregopher/mediocre-go-lib/v2/mctx"
	"github.com/mediocregopher/mediocre-go-lib/v2/mlog"
)

func importAsset(assetStore post.AssetStore, id, path string) error {

	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("opening file: %w", err)
	}
	defer f.Close()

	if err := assetStore.Set(id, f); err != nil {
		return fmt.Errorf("setting into asset store: %w", err)
	}

	return nil
}

func main() {

	ctx := context.Background()

	cfg := cfgpkg.NewBlogCfg(cfgpkg.Params{})

	var dataDir cfgpkg.DataDir
	dataDir.SetupCfg(cfg)
	defer dataDir.Close()
	ctx = mctx.WithAnnotator(ctx, &dataDir)

	id := cfg.String("id", "", "ID the asset will be stored under")
	path := cfg.String("path", "", "Path the asset should be imported from")

	fromStdin := cfg.Bool("from-stdin", false, "If set, ignore id and path, read space separated id/path pairs from stdin")

	// initialization
	err := cfg.Init(ctx)

	logger := mlog.NewLogger(nil)
	defer logger.Close()

	if !*fromStdin && (*id == "" || *path == "") {
		logger.FatalString(ctx, "-id and -path are required if -from-stdin is not given")
	}

	logger.Info(ctx, "process started")
	defer logger.Info(ctx, "process exiting")

	if err != nil {
		logger.Fatal(ctx, "initializing", err)
	}

	postDB, err := post.NewSQLDB(dataDir)
	if err != nil {
		logger.Fatal(ctx, "initializing post sql db", err)
	}
	defer postDB.Close()

	assetStore := post.NewAssetStore(postDB)

	if !*fromStdin {

		ctx := mctx.Annotate(ctx, "id", *id, "path", *path)

		if err := importAsset(assetStore, *id, *path); err != nil {
			logger.Fatal(ctx, "failed to import asset", err)
		}

		logger.Info(ctx, "asset stored")

		return
	}

	for stdin := bufio.NewReader(os.Stdin); ; {

		line, err := stdin.ReadString('\n')

		if errors.Is(err, io.EOF) {
			return
		} else if err != nil {
			logger.Fatal(ctx, "reading from stdin", err)
		}

		fields := strings.Fields(line)

		if len(fields) < 2 {
			ctx := mctx.Annotate(ctx, "line", line)
			logger.FatalString(ctx, "cannot process line with fewer than 2 fields")
		}

		id, path := fields[0], fields[1]

		ctx := mctx.Annotate(ctx, "id", id, "path", path)

		if err := importAsset(assetStore, id, path); err != nil {
			logger.Fatal(ctx, "failed to import asset", err)
		}

		logger.Info(ctx, "asset stored")
	}
}
