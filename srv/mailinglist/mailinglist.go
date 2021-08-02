// Package mailinglist manages the list of subscribed emails and allows emailing
// out to them.
package mailinglist

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io"
	"strings"

	"github.com/google/uuid"
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

	// URL of the page which should be navigated to in order to finalize a
	// subscription.
	FinalizeSubURL string

	// URL of the page which should be navigated to in order to remove a
	// subscription.
	UnsubURL string
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
		SubLink: fmt.Sprintf("%s?subToken=%s", m.params.FinalizeSubURL, emailRecord.SubToken),
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
			UnsubURL:  fmt.Sprintf("%s?unsubToken=%s", m.params.UnsubURL, emailRecord.UnsubToken),
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
