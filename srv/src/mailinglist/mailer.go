package mailinglist

import (
	"context"
	"errors"
	"strings"

	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
	"github.com/mediocregopher/blog.mediocregopher.com/srv/cfg"
	"github.com/mediocregopher/mediocre-go-lib/v2/mctx"
	"github.com/mediocregopher/mediocre-go-lib/v2/mlog"
)

// Mailer is used to deliver emails to arbitrary recipients.
type Mailer interface {
	Send(to, subject, body string) error
}

type logMailer struct {
	logger *mlog.Logger
}

// NewLogMailer returns a Mailer instance which will not actually send any
// emails, it will only log to the given Logger when Send is called.
func NewLogMailer(logger *mlog.Logger) Mailer {
	return &logMailer{logger: logger}
}

func (l *logMailer) Send(to, subject, body string) error {
	ctx := mctx.Annotate(context.Background(),
		"to", to,
		"subject", subject,
	)
	l.logger.Info(ctx, "would have sent email")
	return nil
}

// NullMailer acts as a Mailer but actually just does nothing.
var NullMailer = nullMailer{}

type nullMailer struct{}

func (nullMailer) Send(to, subject, body string) error {
	return nil
}

// MailerParams are used to initialize a new Mailer instance.
type MailerParams struct {
	SMTPAddr string

	// Optional, if not given then no auth is attempted.
	SMTPAuth sasl.Client

	// The sending email address to use for all emails being sent.
	SendAs string
}

// SetupCfg implement the cfg.Cfger interface.
func (m *MailerParams) SetupCfg(cfg *cfg.Cfg) {

	cfg.StringVar(&m.SMTPAddr, "ml-smtp-addr", "", "Address of SMTP server to use for sending emails for the mailing list")
	smtpAuthStr := cfg.String("ml-smtp-auth", "", "user:pass to use when authenticating with the mailing list SMTP server. The given user will also be used as the From address.")

	cfg.OnInit(func(ctx context.Context) error {
		if m.SMTPAddr == "" {
			return nil
		}

		smtpAuthParts := strings.SplitN(*smtpAuthStr, ":", 2)
		if len(smtpAuthParts) < 2 {
			return errors.New("invalid -ml-smtp-auth")
		}

		m.SMTPAuth = sasl.NewPlainClient("", smtpAuthParts[0], smtpAuthParts[1])
		m.SendAs = smtpAuthParts[0]

		return nil
	})
}

// Annotate implements mctx.Annotator interface.
func (m *MailerParams) Annotate(a mctx.Annotations) {
	if m.SMTPAddr == "" {
		return
	}

	a["smtpAddr"] = m.SMTPAddr
	a["smtpSendAs"] = m.SendAs
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
