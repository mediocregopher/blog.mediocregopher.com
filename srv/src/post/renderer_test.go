package post

import (
	"bytes"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMarkdownBodyToHTML(t *testing.T) {

	tests := []struct {
		body string
		exp  string
	}{
		{
			body: `
# Foo
`,
			exp: `<h1 id="foo">Foo</h1>`,
		},
		{
			body: `
this is a body

this is another
`,
			exp: `
<p>this is a body</p>

<p>this is another</p>`,
		},
		{
			body: `this is a [link](somewhere.html)`,
			exp:  `<p>this is a <a href="somewhere.html" target="_blank">link</a></p>`,
		},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {

			outB := mdBodyToHTML([]byte(test.body))
			out := string(outB)

			// just to make the tests nicer
			out = strings.TrimSpace(out)
			test.exp = strings.TrimSpace(test.exp)

			assert.Equal(t, test.exp, out)
		})
	}
}

func TestMarkdownToHTMLRenderer(t *testing.T) {

	r := NewMarkdownToHTMLRenderer()

	post := RenderablePost{
		StoredPost: StoredPost{
			Post: Post{
				ID:          "foo",
				Title:       "Foo",
				Description: "Bar.",
				Body:        "This is the body.",
				Series:      "baz",
			},
			PublishedAt: time.Now(),
		},

		SeriesPrevious: &StoredPost{
			Post: Post{
				ID:    "foo-prev",
				Title: "Foo Prev",
			},
		},

		SeriesNext: &StoredPost{
			Post: Post{
				ID:    "foo-next",
				Title: "Foo Next",
			},
		},
	}

	buf := new(bytes.Buffer)
	err := r.Render(buf, post)
	assert.NoError(t, err)
	t.Log(buf.String())
}
