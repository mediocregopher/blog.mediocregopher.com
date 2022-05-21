package http

import (
	"context"
	"net/http"
	"time"

	"github.com/mediocregopher/blog.mediocregopher.com/srv/http/apiutil"
	"golang.org/x/crypto/bcrypt"
)

// NewPasswordHash returns the hash of the given plaintext password, for use
// with Auther.
func NewPasswordHash(plaintext string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plaintext), 13)
	if err != nil {
		panic(err)
	}
	return string(hashedPassword)
}

// Auther determines who can do what.
type Auther interface {
	Allowed(ctx context.Context, username, password string) bool
	Close() error
}

type auther struct {
	users  map[string]string
	ticker *time.Ticker
}

// NewAuther initializes and returns an Auther will which allow the given
// username and password hash combinations. Password hashes must have been
// created using NewPasswordHash.
func NewAuther(users map[string]string, ratelimit time.Duration) Auther {
	return &auther{
		users:  users,
		ticker: time.NewTicker(ratelimit),
	}
}

func (a *auther) Close() error {
	a.ticker.Stop()
	return nil
}

func (a *auther) Allowed(ctx context.Context, username, password string) bool {

	select {
	case <-ctx.Done():
		return false
	case <-a.ticker.C:
	}

	hashedPassword, ok := a.users[username]
	if !ok {
		return false
	}

	err := bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword), []byte(password),
	)

	return err == nil
}

func authMiddleware(auther Auther) middleware {

	respondUnauthorized := func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("WWW-Authenticate", `Basic realm="NOPE"`)
		rw.WriteHeader(http.StatusUnauthorized)
		apiutil.GetRequestLogger(r).WarnString(r.Context(), "unauthorized")
	}

	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

			username, password, ok := r.BasicAuth()

			if !ok {
				respondUnauthorized(rw, r)
				return
			}

			if !auther.Allowed(r.Context(), username, password) {
				respondUnauthorized(rw, r)
				return
			}

			h.ServeHTTP(rw, r)
		})
	}
}
