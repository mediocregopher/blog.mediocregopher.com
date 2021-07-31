package mailinglist

import (
	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
)

// Mailer is used to deliver emails to arbitrary recipients.
type Mailer interface {
	Send(to, subject, body string) error
}

// MailerParams are used to initialize a new Mailer instance
type MailerParams struct {
	SMTPAddr string

	// Optional, if not given then no auth is attempted.
	SMTPAuth sasl.Client

	// The sending email address to use for all emails being sent.
	SendAs string
}

type mailer struct {
	params MailerParams
}

// NewMailer initializes and returns a Mailer which will use an external SMTP
// server to deliver email.
func NewMailer(params MailerParams) Mailer {
	return &mailer{
		params: params,
	}
}

func (m *mailer) Send(to, subject, body string) error {

	msg := []byte("From: " + m.params.SendAs + "\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n\r\n" +
		body + "\r\n")

	c, err := smtp.Dial(m.params.SMTPAddr)
	if err != nil {
		return err
	}
	defer c.Close()

	if err = c.Auth(m.params.SMTPAuth); err != nil {
		return err
	}

	if err = c.Mail(m.params.SendAs, nil); err != nil {
		return err
	}

	if err = c.Rcpt(to); err != nil {
		return err
	}

	w, err := c.Data()
	if err != nil {
		return err
	}

	if _, err = w.Write(msg); err != nil {
		return err
	}

	if err = w.Close(); err != nil {
		return err
	}

	return c.Quit()
}
