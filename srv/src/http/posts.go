package http

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/mediocregopher/blog.mediocregopher.com/srv/http/apiutil"
	"github.com/mediocregopher/blog.mediocregopher.com/srv/post"
)

type postTplPayload struct {
	post.StoredPost
	SeriesPrevious, SeriesNext *post.StoredPost
	Body                       template.HTML
}

func (a *api) postToPostTplPayload(storedPost post.StoredPost) (postTplPayload, error) {

	bodyTpl, err := a.parseTpl(storedPost.Body)
	if err != nil {
		return postTplPayload{}, fmt.Errorf("parsing post body as template: %w", err)
	}

	bodyBuf := new(bytes.Buffer)
	if err := bodyTpl.Execute(bodyBuf, nil); err != nil {
		return postTplPayload{}, fmt.Errorf("executing post body as template: %w", err)
	}

	parserExt := parser.CommonExtensions | parser.AutoHeadingIDs
	parser := parser.NewWithExtensions(parserExt)

	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	htmlRenderer := html.NewRenderer(html.RendererOptions{Flags: htmlFlags})

	renderedBody := markdown.ToHTML(bodyBuf.Bytes(), parser, htmlRenderer)

	tplPayload := postTplPayload{
		StoredPost: storedPost,
		Body:       template.HTML(renderedBody),
	}

	if series := storedPost.Series; series != "" {

		seriesPosts, err := a.params.PostStore.GetBySeries(series)
		if err != nil {
			return postTplPayload{}, fmt.Errorf(
				"fetching posts for series %q: %w", series, err,
			)
		}

		var foundThis bool

		for i := range seriesPosts {

			seriesPost := seriesPosts[i]

			if seriesPost.ID == storedPost.ID {
				foundThis = true
				continue
			}

			if !foundThis {
				tplPayload.SeriesNext = &seriesPost
				continue
			}

			tplPayload.SeriesPrevious = &seriesPost
			break
		}
	}

	return tplPayload, nil
}

func (a *api) renderPostHandler() http.Handler {

	tpl := a.mustParseBasedTpl("post.html")
	renderIndexHandler := a.renderPostsIndexHandler()

	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		id := strings.TrimSuffix(filepath.Base(r.URL.Path), ".html")

		if id == "/" {
			renderIndexHandler.ServeHTTP(rw, r)
			return
		}

		storedPost, err := a.params.PostStore.GetByID(id)

		if errors.Is(err, post.ErrPostNotFound) {
			http.Error(rw, "Post not found", 404)
			return
		} else if err != nil {
			apiutil.InternalServerError(
				rw, r, fmt.Errorf("fetching post with id %q: %w", id, err),
			)
			return
		}

		tplPayload, err := a.postToPostTplPayload(storedPost)

		if err != nil {
			apiutil.InternalServerError(
				rw, r, fmt.Errorf(
					"generating template payload for post with id %q: %w",
					id, err,
				),
			)
			return
		}

		executeTemplate(rw, r, tpl, tplPayload)
	})
}

func (a *api) renderPostsIndexHandler() http.Handler {

	tpl := a.mustParseBasedTpl("posts.html")
	const pageCount = 20

	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		page, err := apiutil.StrToInt(r.FormValue("p"), 0)
		if err != nil {
			apiutil.BadRequest(
				rw, r, fmt.Errorf("invalid page number: %w", err),
			)
			return
		}

		posts, hasMore, err := a.params.PostStore.Get(page, pageCount)
		if err != nil {
			apiutil.InternalServerError(
				rw, r, fmt.Errorf("fetching page %d of posts: %w", page, err),
			)
			return
		}

		tplPayload := struct {
			Posts              []post.StoredPost
			PrevPage, NextPage int
		}{
			Posts:    posts,
			PrevPage: -1,
			NextPage: -1,
		}

		if page > 0 {
			tplPayload.PrevPage = page - 1
		}

		if hasMore {
			tplPayload.NextPage = page + 1
		}

		executeTemplate(rw, r, tpl, tplPayload)
	})
}

func (a *api) editPostHandler() http.Handler {

	tpl := a.mustParseBasedTpl("edit-post.html")

	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		id := filepath.Base(r.URL.Path)

		var storedPost post.StoredPost

		if id != "/" {

			var err error
			storedPost, err = a.params.PostStore.GetByID(id)

			if errors.Is(err, post.ErrPostNotFound) {
				http.Error(rw, "Post not found", 404)
				return
			} else if err != nil {
				apiutil.InternalServerError(
					rw, r, fmt.Errorf("fetching post with id %q: %w", id, err),
				)
				return
			}
		}

		executeTemplate(rw, r, tpl, storedPost)
	})
}

func postFromPostReq(r *http.Request) (post.Post, error) {

	p := post.Post{
		ID:          r.PostFormValue("id"),
		Title:       r.PostFormValue("title"),
		Description: r.PostFormValue("description"),
		Tags:        strings.Fields(r.PostFormValue("tags")),
		Series:      r.PostFormValue("series"),
	}

	// textareas encode newlines as CRLF for historical reasons
	p.Body = strings.ReplaceAll(p.Body, "\r\n", "\n")
	p.Body = strings.TrimSpace(r.PostFormValue("body"))

	if p.ID == "" ||
		p.Title == "" ||
		p.Description == "" ||
		p.Body == "" ||
		len(p.Tags) == 0 {
		return post.Post{}, errors.New("ID, Title, Description, Tags, and Body are all required")
	}

	return p, nil
}

func (a *api) postPostHandler() http.Handler {

	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		p, err := postFromPostReq(r)
		if err != nil {
			apiutil.BadRequest(rw, r, err)
			return
		}

		if err := a.params.PostStore.Set(p, time.Now()); err != nil {
			apiutil.InternalServerError(
				rw, r, fmt.Errorf("storing post with id %q: %w", p.ID, err),
			)
			return
		}

		redirectPath := fmt.Sprintf("posts/%s?method=edit", p.ID)

		a.executeRedirectTpl(rw, r, redirectPath)
	})
}

func (a *api) deletePostHandler() http.Handler {

	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		id := filepath.Base(r.URL.Path)

		if id == "/" {
			apiutil.BadRequest(rw, r, errors.New("id is required"))
			return
		}

		err := a.params.PostStore.Delete(id)

		if errors.Is(err, post.ErrPostNotFound) {
			http.Error(rw, "Post not found", 404)
			return
		} else if err != nil {
			apiutil.InternalServerError(
				rw, r, fmt.Errorf("deleting post with id %q: %w", id, err),
			)
			return
		}

		a.executeRedirectTpl(rw, r, "posts/")

	})
}

func (a *api) previewPostHandler() http.Handler {

	tpl := a.mustParseBasedTpl("post.html")

	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		p, err := postFromPostReq(r)
		if err != nil {
			apiutil.BadRequest(rw, r, err)
			return
		}

		storedPost := post.StoredPost{
			Post:        p,
			PublishedAt: time.Now(),
		}

		tplPayload, err := a.postToPostTplPayload(storedPost)

		if err != nil {
			apiutil.InternalServerError(
				rw, r, fmt.Errorf("generating template payload: %w", err),
			)
			return
		}

		executeTemplate(rw, r, tpl, tplPayload)
	})
}
