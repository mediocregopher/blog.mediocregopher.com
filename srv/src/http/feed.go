package http

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gorilla/feeds"
	"github.com/mediocregopher/blog.mediocregopher.com/srv/http/apiutil"
)

func (a *api) renderFeedHandler() http.Handler {

	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		author := &feeds.Author{
			Name: "mediocregopher",
		}

		publicURL := a.params.PublicURL.String()

		feed := feeds.Feed{
			Title:       "Mediocre Blog",
			Link:        &feeds.Link{Href: publicURL + "/"},
			Description: "A mix of tech, art, travel, and who knows what else.",
			Author:      author,
		}

		recentPosts, _, err := a.params.PostStore.WithOrderDesc().Get(0, 20)
		if err != nil {
			apiutil.InternalServerError(rw, r, fmt.Errorf("fetching recent posts: %w", err))
			return
		}

		for _, post := range recentPosts {

			if post.PublishedAt.After(feed.Updated) {
				feed.Updated = post.PublishedAt
			}

			if post.LastUpdatedAt.After(feed.Updated) {
				feed.Updated = post.LastUpdatedAt
			}

			postURL := publicURL + filepath.Join("/posts", post.ID)

			feed.Items = append(feed.Items, &feeds.Item{
				Title:       post.Title,
				Link:        &feeds.Link{Href: postURL},
				Author:      author,
				Description: post.Description,
				Id:          postURL,
				Updated:     post.LastUpdatedAt,
				Created:     post.PublishedAt,
			})
		}

		if err := feed.WriteAtom(rw); err != nil {
			apiutil.InternalServerError(rw, r, fmt.Errorf("writing atom feed: %w", err))
			return
		}
	})
}
