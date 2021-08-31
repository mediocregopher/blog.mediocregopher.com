package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"unicode"

	"github.com/mediocregopher/blog.mediocregopher.com/srv/api/apiutils"
	"github.com/mediocregopher/blog.mediocregopher.com/srv/chat"
)

type chatHandler struct {
	*http.ServeMux

	room       chat.Room
	userIDCalc *chat.UserIDCalculator
}

func newChatHandler(
	room chat.Room, userIDCalc *chat.UserIDCalculator,
	requirePowMiddleware func(http.Handler) http.Handler,
) http.Handler {
	c := &chatHandler{
		ServeMux:   http.NewServeMux(),
		room:       room,
		userIDCalc: userIDCalc,
	}

	c.Handle("/history", c.historyHandler())
	c.Handle("/user-id", requirePowMiddleware(c.userIDHandler()))
	c.Handle("/append", requirePowMiddleware(c.appendHandler()))

	return c
}

func (c *chatHandler) historyHandler() http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		limit, err := apiutils.StrToInt(r.PostFormValue("limit"), 0)
		if err != nil {
			apiutils.BadRequest(rw, r, fmt.Errorf("invalid limit parameter: %w", err))
			return
		}

		cursor := r.PostFormValue("cursor")

		cursor, msgs, err := c.room.History(r.Context(), chat.HistoryOpts{
			Limit:  limit,
			Cursor: cursor,
		})

		if argErr := (chat.ErrInvalidArg{}); errors.As(err, &argErr) {
			apiutils.BadRequest(rw, r, argErr.Err)
			return
		} else if err != nil {
			apiutils.InternalServerError(rw, r, err)
		}

		apiutils.JSONResult(rw, r, struct {
			Cursor   string         `json:"cursor"`
			Messages []chat.Message `json:"messages"`
		}{
			Cursor:   cursor,
			Messages: msgs,
		})
	})
}

func (c *chatHandler) userID(r *http.Request) (chat.UserID, error) {
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

	return c.userIDCalc.Calculate(name, password), nil
}

func (c *chatHandler) userIDHandler() http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		userID, err := c.userID(r)
		if err != nil {
			apiutils.BadRequest(rw, r, err)
			return
		}

		apiutils.JSONResult(rw, r, struct {
			UserID chat.UserID `json:"userID"`
		}{
			UserID: userID,
		})
	})
}

func (c *chatHandler) appendHandler() http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		userID, err := c.userID(r)
		if err != nil {
			apiutils.BadRequest(rw, r, err)
			return
		}

		body := r.PostFormValue("body")

		if l := len(body); l == 0 {
			apiutils.BadRequest(rw, r, errors.New("body is required"))
			return

		} else if l > 300 {
			apiutils.BadRequest(rw, r, errors.New("body too long"))
			return
		}

		msg, err := c.room.Append(r.Context(), chat.Message{
			UserID: userID,
			Body:   body,
		})

		if err != nil {
			apiutils.InternalServerError(rw, r, err)
			return
		}

		apiutils.JSONResult(rw, r, struct {
			MessageID string `json:"messageID"`
		}{
			MessageID: msg.ID,
		})
	})
}
