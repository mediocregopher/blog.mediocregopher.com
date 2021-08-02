package main

import (
	"errors"
	"net/http"
	"strings"

	"github.com/mediocregopher/blog.mediocregopher.com/srv/mailinglist"
)

func mailingListSubscribeHandler(ml mailinglist.MailingList) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		email := r.PostFormValue("email")
		if parts := strings.Split(email, "@"); len(parts) != 2 ||
			parts[0] == "" ||
			parts[1] == "" ||
			len(email) >= 512 {
			badRequest(rw, r, errors.New("invalid email"))
		}

		if err := ml.BeginSubscription(email); errors.Is(err, mailinglist.ErrAlreadyVerified) {
			// just eat the error, make it look to the user like the
			// verification email was sent.
		} else if err != nil {
			internalServerError(rw, r, err)
		}
	})
}

func mailingListFinalizeHandler(ml mailinglist.MailingList) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		subToken := r.PostFormValue("subToken")
		if l := len(subToken); l == 0 || l > 128 {
			badRequest(rw, r, errors.New("invalid subToken"))
			return
		}

		err := ml.FinalizeSubscription(subToken)
		if errors.Is(err, mailinglist.ErrNotFound) ||
			errors.Is(err, mailinglist.ErrAlreadyVerified) {
			badRequest(rw, r, err)
			return
		} else if err != nil {
			internalServerError(rw, r, err)
			return
		}
	})
}

func mailingListUnsubscribeHandler(ml mailinglist.MailingList) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		unsubToken := r.PostFormValue("unsubToken")
		if l := len(unsubToken); l == 0 || l > 128 {
			badRequest(rw, r, errors.New("invalid unsubToken"))
			return
		}

		err := ml.Unsubscribe(unsubToken)
		if errors.Is(err, mailinglist.ErrNotFound) {
			badRequest(rw, r, err)
			return
		} else if err != nil {
			internalServerError(rw, r, err)
			return
		}
	})
}
