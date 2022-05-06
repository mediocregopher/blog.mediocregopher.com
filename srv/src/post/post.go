// Package post deals with the storage and rending of blog post.
package post

import (
	"fmt"
	"path"
	"regexp"
	"strings"
	"time"
)

// Date represents a calendar date with no timezone information attached.
type Date struct {
	Year  int
	Month time.Month
	Day   int
}

// DateFromTime converts a Time into a Date, truncating all non-date
// information.
func DateFromTime(t time.Time) Date {
	return Date{
		Year:  t.Year(),
		Month: t.Month(),
		Day:   t.Day(),
	}
}

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

	PublishedAt   Date
	LastUpdatedAt Date

	Body string
}

// URL returns the relative URL of the Post.
func (p Post) URL() string {
	return path.Join(
		fmt.Sprintf(
			"%d/%0d/%0d",
			p.PublishedAt.Year,
			p.PublishedAt.Month,
			p.PublishedAt.Day,
		),
		p.ID+".html",
	)
}
