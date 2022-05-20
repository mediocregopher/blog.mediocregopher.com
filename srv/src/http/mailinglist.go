package http

import (
	"errors"
	"net/http"
	"strings"

	"github.com/mediocregopher/blog.mediocregopher.com/srv/http/apiutil"
	"github.com/mediocregopher/blog.mediocregopher.com/srv/mailinglist"
)

func (a *api) mailingListSubscribeHandler() http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		email := r.PostFormValue("email")
		if parts := strings.Split(email, "@"); len(parts) != 2 ||
			parts[0] == "" ||
			parts[1] == "" ||
			len(email) >= 512 {
			apiutil.BadRequest(rw, r, errors.New("invalid email"))
			return
		}

		err := a.params.MailingList.BeginSubscription(email)

		if errors.Is(err, mailinglist.ErrAlreadyVerified) {
			// just eat the error, make it look to the user like the
			// verification email was sent.
		} else if err != nil {
			apiutil.InternalServerError(rw, r, err)
			return
		}

		apiutil.JSONResult(rw, r, struct{}{})
	})
}

func (a *api) mailingListFinalizeHandler() http.Handler {
	var errInvalidSubToken = errors.New("invalid subToken")

	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		subToken := r.PostFormValue("subToken")
		if l := len(subToken); l == 0 || l > 128 {
			apiutil.BadRequest(rw, r, errInvalidSubToken)
			return
		}

		err := a.params.MailingList.FinalizeSubscription(subToken)

		if errors.Is(err, mailinglist.ErrNotFound) {
			apiutil.BadRequest(rw, r, errInvalidSubToken)
			return

		} else if errors.Is(err, mailinglist.ErrAlreadyVerified) {
			// no problem

		} else if err != nil {
			apiutil.InternalServerError(rw, r, err)
			return
		}

		apiutil.JSONResult(rw, r, struct{}{})
	})
}

func (a *api) mailingListUnsubscribeHandler() http.Handler {
	var errInvalidUnsubToken = errors.New("invalid unsubToken")

	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		unsubToken := r.PostFormValue("unsubToken")
		if l := len(unsubToken); l == 0 || l > 128 {
			apiutil.BadRequest(rw, r, errInvalidUnsubToken)
			return
		}

		err := a.params.MailingList.Unsubscribe(unsubToken)

		if errors.Is(err, mailinglist.ErrNotFound) {
			apiutil.BadRequest(rw, r, errInvalidUnsubToken)
			return

		} else if err != nil {
			apiutil.InternalServerError(rw, r, err)
			return
		}

		apiutil.JSONResult(rw, r, struct{}{})
	})
}
