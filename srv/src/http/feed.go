package http

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gorilla/feeds"
	"github.com/mediocregopher/blog.mediocregopher.com/srv/http/apiutil"
	"github.com/mediocregopher/blog.mediocregopher.com/srv/post"
)

func (a *api) renderFeedHandler() http.Handler {

	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		tag := r.FormValue("tag")

		var (
			posts []post.StoredPost
			err   error
		)

		if tag == "" {
			posts, _, err = a.params.PostStore.Get(0, 20)
		} else {
			posts, err = a.params.PostStore.GetByTag(tag)
		}

		if err != nil {
			apiutil.InternalServerError(rw, r, fmt.Errorf("fetching recent posts: %w", err))
			return
		}

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

		for _, post := range posts {

			if post.PublishedAt.After(feed.Updated) {
				feed.Updated = post.PublishedAt
			}

			if post.LastUpdatedAt.After(feed.Updated) {
				feed.Updated = post.LastUpdatedAt
			}

			postURL := publicURL + filepath.Join("/posts", post.ID)

			item := &feeds.Item{
				Title:       post.Title,
				Link:        &feeds.Link{Href: postURL},
				Author:      author,
				Description: post.Description,
				Id:          postURL,
				Created:     post.PublishedAt,
			}

			feed.Items = append(feed.Items, item)
		}

		if err := feed.WriteAtom(rw); err != nil {
			apiutil.InternalServerError(rw, r, fmt.Errorf("writing atom feed: %w", err))
			return
		}
	})
}
