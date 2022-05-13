package post

import (
	_ "embed"
	"fmt"
	"html/template"
	"io"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/mediocregopher/blog.mediocregopher.com/srv/tpl"
)

// RenderablePost is a Post wrapped with extra information necessary for
// rendering.
type RenderablePost struct {
	StoredPost
	SeriesPrevious, SeriesNext *StoredPost
}

// NewRenderablePost wraps an existing Post such that it can be rendered.
func NewRenderablePost(store Store, post StoredPost) (RenderablePost, error) {

	renderablePost := RenderablePost{
		StoredPost: post,
	}

	if post.Series != "" {

		seriesPosts, err := store.GetBySeries(post.Series)
		if err != nil {
			return RenderablePost{}, fmt.Errorf(
				"fetching posts for series %q: %w",
				post.Series, err,
			)
		}

		var foundThis bool

		for i := range seriesPosts {

			seriesPost := seriesPosts[i]

			if seriesPost.ID == post.ID {
				foundThis = true
				continue
			}

			if !foundThis {
				renderablePost.SeriesPrevious = &seriesPost
				continue
			}

			renderablePost.SeriesNext = &seriesPost
			break
		}
	}

	return renderablePost, nil
}

// Renderer takes a Post and renders it to some encoding.
type Renderer interface {
	Render(io.Writer, RenderablePost) error
}

func mdBodyToHTML(body []byte) []byte {
	parserExt := parser.CommonExtensions | parser.AutoHeadingIDs
	parser := parser.NewWithExtensions(parserExt)

	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	htmlRenderer := html.NewRenderer(html.RendererOptions{Flags: htmlFlags})

	return markdown.ToHTML(body, parser, htmlRenderer)
}

type mdHTMLRenderer struct{}

// NewMarkdownToHTMLRenderer renders Posts from markdown to HTML.
func NewMarkdownToHTMLRenderer() Renderer {
	return mdHTMLRenderer{}
}

func (r mdHTMLRenderer) Render(into io.Writer, post RenderablePost) error {

	data := struct {
		RenderablePost
		Body template.HTML
	}{
		RenderablePost: post,
		Body:           template.HTML(mdBodyToHTML([]byte(post.Body))),
	}

	return tpl.HTML.ExecuteTemplate(into, "post.html", data)
}
