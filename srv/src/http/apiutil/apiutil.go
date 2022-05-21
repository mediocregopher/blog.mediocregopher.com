// Package apiutil contains utilities which are useful for implementing api
// endpoints.
package apiutil

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/mediocregopher/mediocre-go-lib/v2/mlog"
)

type loggerCtxKey int

// SetRequestLogger sets the given Logger onto the given Request's Context,
// returning a copy.
func SetRequestLogger(r *http.Request, logger *mlog.Logger) *http.Request {
	ctx := r.Context()
	ctx = context.WithValue(ctx, loggerCtxKey(0), logger)
	return r.WithContext(ctx)
}

// GetRequestLogger returns the Logger which was set by SetRequestLogger onto
// this Request, or nil.
func GetRequestLogger(r *http.Request) *mlog.Logger {
	ctx := r.Context()
	logger, _ := ctx.Value(loggerCtxKey(0)).(*mlog.Logger)
	if logger == nil {
		logger = mlog.Null
	}
	return logger
}

// JSONResult writes the JSON encoding of the given value as the response body.
func JSONResult(rw http.ResponseWriter, r *http.Request, v interface{}) {
	b, err := json.Marshal(v)
	if err != nil {
		InternalServerError(rw, r, err)
		return
	}
	b = append(b, '\n')

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(b)
}

// BadRequest writes a 400 status and a JSON encoded error struct containing the
// given error as the response body.
func BadRequest(rw http.ResponseWriter, r *http.Request, err error) {
	GetRequestLogger(r).Warn(r.Context(), "bad request", err)

	rw.WriteHeader(400)
	JSONResult(rw, r, struct {
		Error string `json:"error"`
	}{
		Error: err.Error(),
	})
}

// InternalServerError writes a 500 status and a JSON encoded error struct
// containing a generic error as the response body (though it will log the given
// one).
func InternalServerError(rw http.ResponseWriter, r *http.Request, err error) {
	GetRequestLogger(r).Error(r.Context(), "internal server error", err)

	rw.WriteHeader(500)
	JSONResult(rw, r, struct {
		Error string `json:"error"`
	}{
		Error: "internal server error",
	})
}

// StrToInt parses the given string as an integer, or returns the given default
// integer if the string is empty.
func StrToInt(str string, defaultVal int) (int, error) {
	if str == "" {
		return defaultVal, nil
	}
	return strconv.Atoi(str)
}

// GetCookie returns the namd cookie's value, or the given default value if the
// cookie is not set.
//
// This will only return an error if there was an unexpected error parsing the
// Request's cookies.
func GetCookie(r *http.Request, cookieName, defaultVal string) (string, error) {
	c, err := r.Cookie(cookieName)
	if errors.Is(err, http.ErrNoCookie) {
		return defaultVal, nil
	} else if err != nil {
		return "", fmt.Errorf("reading cookie %q: %w", cookieName, err)
	}

	return c.Value, nil
}

// RandStr returns a human-readable random string with the given number of bytes
// of randomness.
func RandStr(numBytes int) string {
	b := make([]byte, numBytes)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}

// MethodMux will take the request method (GET, POST, etc...) and handle the
// request using the corresponding Handler in the given map.
//
// If no Handler is defined for a method then a 405 Method Not Allowed error is
// returned.
func MethodMux(handlers map[string]http.Handler) http.Handler {

	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		method := strings.ToUpper(r.Method)
		formMethod := strings.ToUpper(r.FormValue("method"))

		if method == "POST" && formMethod != "" {
			method = formMethod
		}

		handler, ok := handlers[method]

		if !ok {
			http.Error(rw, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		handler.ServeHTTP(rw, r)
	})
}
