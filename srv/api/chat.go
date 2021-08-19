package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"unicode"

	"github.com/mediocregopher/blog.mediocregopher.com/srv/chat"
)

func (a *api) chatHistoryHandler() http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		limit, err := strToInt(r.FormValue("limit"), 0)
		if err != nil {
			badRequest(rw, r, fmt.Errorf("invalid limit parameter: %w", err))
			return
		}

		cursor := r.FormValue("cursor")

		cursor, msgs, err := a.params.GlobalRoom.History(r.Context(), chat.HistoryOpts{
			Limit:  limit,
			Cursor: cursor,
		})

		if argErr := (chat.ErrInvalidArg{}); errors.As(err, &argErr) {
			badRequest(rw, r, argErr.Err)
			return
		} else if err != nil {
			internalServerError(rw, r, err)
		}

		jsonResult(rw, r, struct {
			Cursor   string         `json:"cursor"`
			Messages []chat.Message `json:"messages"`
		}{
			Cursor:   cursor,
			Messages: msgs,
		})
	})
}

func (a *api) getUserID(r *http.Request) (chat.UserID, error) {
	name := r.PostFormValue("name")
	if l := len(name); l == 0 {
		return chat.UserID{}, errors.New("name is required")
	} else if l > 16 {
		return chat.UserID{}, errors.New("name too long")
	}

	nameClean := strings.Map(func(r rune) rune {
		if !unicode.IsPrint(r) {
			return -1
		}
		return r
	}, name)

	if nameClean != name {
		return chat.UserID{}, errors.New("name contains invalid characters")
	}

	password := r.PostFormValue("password")
	if l := len(password); l == 0 {
		return chat.UserID{}, errors.New("password is required")
	} else if l > 128 {
		return chat.UserID{}, errors.New("password too long")
	}

	return a.params.UserIDCalculator.Calculate(name, password), nil
}

func (a *api) getUserIDHandler() http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		userID, err := a.getUserID(r)
		if err != nil {
			badRequest(rw, r, err)
			return
		}

		jsonResult(rw, r, struct {
			UserID chat.UserID `json:"userID"`
		}{
			UserID: userID,
		})
	})
}
