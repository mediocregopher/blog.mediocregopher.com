// Package mailinglist manages the list of subscribed emails and allows emailing
// out to them.
package mailinglist

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"html/template"
	"io"
	"net/url"
	"strings"

	"github.com/google/uuid"
	"github.com/mediocregopher/blog.mediocregopher.com/srv/cfg"
	"github.com/mediocregopher/mediocre-go-lib/v2/mctx"
	"github.com/tilinna/clock"
)

var (
	// ErrAlreadyVerified is used when the email is already fully subscribed.
	ErrAlreadyVerified = errors.New("email is already subscribed")
)

// MailingList is able to subscribe, unsubscribe, and iterate through emails.
type MailingList interface {

	// May return ErrAlreadyVerified.
	BeginSubscription(email string) error

	// May return ErrNotFound or ErrAlreadyVerified.
	FinalizeSubscription(subToken string) error

	// May return ErrNotFound.
	Unsubscribe(unsubToken string) error

	// Publish blasts the mailing list with an update about a new blog post.
	Publish(postTitle, postURL string) error
}

// Params are parameters used to initialize a new MailingList. All fields are
// required unless otherwise noted.
type Params struct {
	Store  Store
	Mailer Mailer
	Clock  clock.Clock

	// PublicURL is the base URL which site visitors can navigate to.
	// MailingList will generate links based on this value.
	PublicURL *url.URL
}

// SetupCfg implement the cfg.Cfger interface.
func (p *Params) SetupCfg(cfg *cfg.Cfg) {
	publicURLStr := cfg.String("ml-public-url", "http://localhost:4000", "URL this service is accessible at")

	cfg.OnInit(func(ctx context.Context) error {
		var err error
		*publicURLStr = strings.TrimSuffix(*publicURLStr, "/")
		if p.PublicURL, err = url.Parse(*publicURLStr); err != nil {
			return fmt.Errorf("parsing -ml-public-url: %w", err)
		}

		return nil
	})
}

// Annotate implements mctx.Annotator interface.
func (p *Params) Annotate(a mctx.Annotations) {
	a["mlPublicURL"] = p.PublicURL
}

// New initializes and returns a MailingList instance using the given Params.
func New(params Params) MailingList {
	return &mailingList{params: params}
}

type mailingList struct {
	params Params
}

var beginSubTpl = template.Must(template.New("beginSub").Parse(`
Welcome to the Mediocre Blog mailing list! By subscribing to this mailing list
you are signing up to receive an email everytime a new blog post is published.

In order to complete your subscription please navigate to the following link:

{{ .SubLink }}

This mailing list is built and run using my own hardware and software, and I
solemnly swear that you'll never receive an email from it unless there's a new
blog post.

If you did not initiate this email, and/or do not wish to subscribe to the
mailing list, then simply delete this email and pretend that nothing ever
happened.

- Brian
`))

func (m *mailingList) BeginSubscription(email string) error {

	emailRecord, err := m.params.Store.Get(email)

	if errors.Is(err, ErrNotFound) {
		emailRecord = Email{
			Email:     email,
			SubToken:  uuid.New().String(),
			CreatedAt: m.params.Clock.Now(),
		}

		if err := m.params.Store.Set(emailRecord); err != nil {
			return fmt.Errorf("storing pending email: %w", err)
		}

	} else if err != nil {
		return fmt.Errorf("finding existing email record: %w", err)

	} else if !emailRecord.VerifiedAt.IsZero() {
		return ErrAlreadyVerified
	}

	body := new(bytes.Buffer)
	err = beginSubTpl.Execute(body, struct {
		SubLink string
	}{
		SubLink: fmt.Sprintf(
			"%s/mailinglist/finalize?subToken=%s",
			m.params.PublicURL.String(),
			emailRecord.SubToken,
		),
	})

	if err != nil {
		return fmt.Errorf("executing beginSubTpl: %w", err)
	}

	err = m.params.Mailer.Send(
		email,
		"Mediocre Blog - Please verify your email address",
		body.String(),
	)

	if err != nil {
		return fmt.Errorf("sending email: %w", err)
	}

	return nil
}

func (m *mailingList) FinalizeSubscription(subToken string) error {
	emailRecord, err := m.params.Store.GetBySubToken(subToken)

	if err != nil {
		return fmt.Errorf("retrieving email record: %w", err)

	} else if !emailRecord.VerifiedAt.IsZero() {
		return ErrAlreadyVerified
	}

	emailRecord.VerifiedAt = m.params.Clock.Now()
	emailRecord.UnsubToken = uuid.New().String()

	if err := m.params.Store.Set(emailRecord); err != nil {
		return fmt.Errorf("storing verified email: %w", err)
	}

	return nil
}

func (m *mailingList) Unsubscribe(unsubToken string) error {
	emailRecord, err := m.params.Store.GetByUnsubToken(unsubToken)

	if err != nil {
		return fmt.Errorf("retrieving email record: %w", err)
	}

	if err := m.params.Store.Delete(emailRecord.Email); err != nil {
		return fmt.Errorf("deleting email record: %w", err)
	}

	return nil
}

var publishTpl = template.Must(template.New("publish").Parse(`
A new post has been published to the Mediocre Blog!

{{ .PostTitle }}
{{ .PostURL }}

If you're interested then please check it out!

If you'd like to unsubscribe from this mailing list then visit the following
link instead:

{{ .UnsubURL }}

- Brian
`))

type multiErr []error

func (m multiErr) Error() string {
	if len(m) == 0 {
		panic("multiErr with no members")
	}

	b := new(strings.Builder)
	fmt.Fprintln(b, "The following errors were encountered:")
	for _, err := range m {
		fmt.Fprintf(b, "\t- %s\n", err.Error())
	}

	return b.String()
}

func (m *mailingList) Publish(postTitle, postURL string) error {

	var mErr multiErr

	iter := m.params.Store.GetAll()
	for {
		emailRecord, err := iter()
		if errors.Is(err, io.EOF) {
			break

		} else if err != nil {
			mErr = append(mErr, fmt.Errorf("iterating through email records: %w", err))
			break

		} else if emailRecord.VerifiedAt.IsZero() {
			continue
		}

		body := new(bytes.Buffer)
		err = publishTpl.Execute(body, struct {
			PostTitle string
			PostURL   string
			UnsubURL  string
		}{
			PostTitle: postTitle,
			PostURL:   postURL,
			UnsubURL: fmt.Sprintf(
				"%s/mailinglist/unsubscribe?unsubToken=%s",
				m.params.PublicURL.String(),
				emailRecord.UnsubToken,
			),
		})

		if err != nil {
			mErr = append(mErr, fmt.Errorf("rendering publish email template for %q: %w", emailRecord.Email, err))
			continue
		}

		err = m.params.Mailer.Send(
			emailRecord.Email,
			fmt.Sprintf("Mediocre Blog - New Post! - %s", postTitle),
			body.String(),
		)

		if err != nil {
			mErr = append(mErr, fmt.Errorf("sending email to %q: %w", emailRecord.Email, err))
			continue
		}
	}

	if len(mErr) > 0 {
		return mErr
	}

	return nil
}
