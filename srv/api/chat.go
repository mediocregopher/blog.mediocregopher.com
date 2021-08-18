package api

import (
	"errors"
	"fmt"
	"net/http"

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
