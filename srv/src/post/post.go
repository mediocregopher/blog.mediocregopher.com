// Package post deals with the storage and rending of blog post.
package post

import (
	"regexp"
	"strings"
)

var titleCleanRegexp = regexp.MustCompile(`[^a-z ]`)

// NewID generates a (hopefully) unique ID based on the given title.
func NewID(title string) string {
	title = strings.ToLower(title)
	title = titleCleanRegexp.ReplaceAllString(title, "")
	title = strings.ReplaceAll(title, " ", "-")
	return title
}

// Post contains all information having to do with a blog post.
type Post struct {
	ID          string
	Title       string
	Description string
	Tags        []string
	Series      string
	Body        string
}
