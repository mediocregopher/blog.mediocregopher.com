package api

import (
	"net/http"

	"github.com/mediocregopher/blog.mediocregopher.com/srv/api/apiutil"
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
	Allowed(username, password string) bool
}

type auther struct {
	users map[string]string
}

// NewAuther initializes and returns an Auther will which allow the given
// username and password hash combinations. Password hashes must have been
// created using NewPasswordHash.
func NewAuther(users map[string]string) Auther {
	return &auther{users: users}
}

func (a *auther) Allowed(username, password string) bool {

	hashedPassword, ok := a.users[username]
	if !ok {
		return false
	}

	err := bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword), []byte(password),
	)

	return err == nil
}

func authMiddleware(auther Auther, h http.Handler) http.Handler {

	respondUnauthorized := func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("WWW-Authenticate", `Basic realm="NOPE"`)
		rw.WriteHeader(http.StatusUnauthorized)
		apiutil.GetRequestLogger(r).WarnString(r.Context(), "unauthorized")
	}

	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		username, password, ok := r.BasicAuth()

		if !ok {
			respondUnauthorized(rw, r)
			return
		}

		if !auther.Allowed(username, password) {
			respondUnauthorized(rw, r)
			return
		}

		h.ServeHTTP(rw, r)
	})
}
