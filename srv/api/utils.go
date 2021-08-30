package api

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/mediocregopher/mediocre-go-lib/v2/mlog"
)

type loggerCtxKey int

func setRequestLogger(r *http.Request, logger *mlog.Logger) *http.Request {
	ctx := r.Context()
	ctx = context.WithValue(ctx, loggerCtxKey(0), logger)
	return r.WithContext(ctx)
}

func getRequestLogger(r *http.Request) *mlog.Logger {
	ctx := r.Context()
	logger, _ := ctx.Value(loggerCtxKey(0)).(*mlog.Logger)
	if logger == nil {
		logger = mlog.Null
	}
	return logger
}

func jsonResult(rw http.ResponseWriter, r *http.Request, v interface{}) {
	b, err := json.Marshal(v)
	if err != nil {
		internalServerError(rw, r, err)
		return
	}
	b = append(b, '\n')

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(b)
}

func badRequest(rw http.ResponseWriter, r *http.Request, err error) {
	getRequestLogger(r).Warn(r.Context(), "bad request", err)

	rw.WriteHeader(400)
	jsonResult(rw, r, struct {
		Error string `json:"error"`
	}{
		Error: err.Error(),
	})
}

func internalServerError(rw http.ResponseWriter, r *http.Request, err error) {
	getRequestLogger(r).Error(r.Context(), "internal server error", err)

	rw.WriteHeader(500)
	jsonResult(rw, r, struct {
		Error string `json:"error"`
	}{
		Error: "internal server error",
	})
}

func strToInt(str string, defaultVal int) (int, error) {
	if str == "" {
		return defaultVal, nil
	}
	return strconv.Atoi(str)
}

func getCookie(r *http.Request, cookieName, defaultVal string) (string, error) {
	c, err := r.Cookie(cookieName)
	if errors.Is(err, http.ErrNoCookie) {
		return defaultVal, nil
	} else if err != nil {
		return "", fmt.Errorf("reading cookie %q: %w", cookieName, err)
	}

	return c.Value, nil
}

func randStr(numBytesEntropy int) string {
	b := make([]byte, numBytesEntropy)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}
