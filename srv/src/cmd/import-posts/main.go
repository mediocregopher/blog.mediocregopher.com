package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/adrg/frontmatter"
	"github.com/mediocregopher/blog.mediocregopher.com/srv/cfg"
	cfgpkg "github.com/mediocregopher/blog.mediocregopher.com/srv/cfg"
	"github.com/mediocregopher/blog.mediocregopher.com/srv/post"
	"github.com/mediocregopher/mediocre-go-lib/v2/mctx"
	"github.com/mediocregopher/mediocre-go-lib/v2/mlog"
)

type postFrontmatter struct {
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
	Tags        string `yaml:"tags"`
	Series      string `yaml:"series"`
	Updated     string `yaml:"updated"`
}

func parseDate(str string) (time.Time, error) {
	const layout = "2006-01-02"
	return time.Parse(layout, str)
}

var postNameRegexp = regexp.MustCompile(`(20..-..-..)-([^.]+).md`)

func importPost(postStore post.Store, path string) (post.StoredPost, error) {

	fileName := filepath.Base(path)
	fileNameMatches := postNameRegexp.FindStringSubmatch(fileName)

	if len(fileNameMatches) != 3 {
		return post.StoredPost{}, fmt.Errorf("file name %q didn't match regex", fileName)
	}

	publishedAtStr := fileNameMatches[1]
	publishedAt, err := parseDate(publishedAtStr)
	if err != nil {
		return post.StoredPost{}, fmt.Errorf("parsing publish date %q: %w", publishedAtStr, err)
	}

	postID := fileNameMatches[2]

	f, err := os.Open(path)
	if err != nil {
		return post.StoredPost{}, fmt.Errorf("opening file: %w", err)
	}
	defer f.Close()

	var matter postFrontmatter

	body, err := frontmatter.Parse(f, &matter)

	if err != nil {
		return post.StoredPost{}, fmt.Errorf("parsing frontmatter: %w", err)
	}

	// if there is already a post for this ID, delete it, we're overwriting
	if err := postStore.Delete(postID); err != nil {
		return post.StoredPost{}, fmt.Errorf("deleting post id %q: %w", postID, err)
	}

	p := post.Post{
		ID:          postID,
		Title:       matter.Title,
		Description: matter.Description,
		Tags:        strings.Fields(matter.Tags),
		Series:      matter.Series,
		Body:        string(body),
	}

	if err := postStore.Set(p, publishedAt); err != nil {
		return post.StoredPost{}, fmt.Errorf("storing post id %q: %w", p.ID, err)
	}

	if matter.Updated != "" {

		lastUpdatedAt, err := parseDate(matter.Updated)
		if err != nil {
			return post.StoredPost{}, fmt.Errorf("parsing updated date %q: %w", matter.Updated, err)
		}

		// as a hack, we store the post again with the updated date as now. This
		// will update the LastUpdatedAt field in the Store.
		if err := postStore.Set(p, lastUpdatedAt); err != nil {
			return post.StoredPost{}, fmt.Errorf("updating post id %q: %w", p.ID, err)
		}
	}

	storedPost, err := postStore.GetByID(p.ID)
	if err != nil {
		return post.StoredPost{}, fmt.Errorf("retrieving stored post by id %q: %w", p.ID, err)
	}

	return storedPost, nil
}

func main() {

	ctx := context.Background()

	cfg := cfg.NewBlogCfg(cfg.Params{})

	var dataDir cfgpkg.DataDir
	dataDir.SetupCfg(cfg)
	defer dataDir.Close()
	ctx = mctx.WithAnnotator(ctx, &dataDir)

	paths := cfg.Args("<post file paths...>")

	// initialization
	err := cfg.Init(ctx)

	logger := mlog.NewLogger(nil)
	defer logger.Close()

	logger.Info(ctx, "process started")
	defer logger.Info(ctx, "process exiting")

	if err != nil {
		logger.Fatal(ctx, "initializing", err)
	}

	if len(*paths) == 0 {
		logger.FatalString(ctx, "no paths given")
	}

	postStore, err := post.NewStore(post.StoreParams{
		DataDir: dataDir,
	})
	if err != nil {
		logger.Fatal(ctx, "initializing post store", err)
	}
	defer postStore.Close()

	for _, path := range *paths {

		ctx := mctx.Annotate(ctx, "postPath", path)

		storedPost, err := importPost(postStore, path)
		if err != nil {
			logger.Error(ctx, "importing post", err)
		}

		ctx = mctx.Annotate(ctx,
			"postID", storedPost.ID,
			"postTitle", storedPost.Title,
			"postDescription", storedPost.Description,
			"postTags", storedPost.Tags,
			"postSeries", storedPost.Series,
			"postPublishedAt", storedPost.PublishedAt,
		)

		if !storedPost.LastUpdatedAt.IsZero() {
			ctx = mctx.Annotate(ctx,
				"postLastUpdatedAt", storedPost.LastUpdatedAt)
		}

		logger.Info(ctx, "post stored")
	}
}
