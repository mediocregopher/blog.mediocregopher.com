package mailinglist

import (
	"io"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStore(t *testing.T) {
	tmpFile, err := ioutil.TempFile(os.TempDir(), "mediocre-blog-mailinglist-store-test-")
	if err != nil {
		t.Fatal("Cannot create temporary file", err)
	}
	tmpFilePath := tmpFile.Name()
	tmpFile.Close()

	t.Logf("using temporary sqlite file at %q", tmpFilePath)

	t.Cleanup(func() {
		if err := os.Remove(tmpFilePath); err != nil {
			panic(err)
		}
	})

	store, err := NewStore(tmpFilePath)
	assert.NoError(t, err)

	t.Cleanup(func() {
		assert.NoError(t, store.Close())
	})

	now := func() time.Time {
		return time.Now().Truncate(time.Second)
	}

	assertGet := func(t *testing.T, email Email) {
		t.Helper()

		gotEmail, err := store.Get(email.Email)
		assert.NoError(t, err)
		assert.Equal(t, email, gotEmail)

		gotEmail, err = store.GetBySubToken(email.SubToken)
		assert.NoError(t, err)
		assert.Equal(t, email, gotEmail)

		if email.UnsubToken != "" {
			gotEmail, err = store.GetByUnsubToken(email.UnsubToken)
			assert.NoError(t, err)
			assert.Equal(t, email, gotEmail)
		}
	}

	assertNotFound := func(t *testing.T, email string) {
		t.Helper()
		_, err := store.Get(email)
		assert.ErrorIs(t, err, ErrNotFound)
	}

	// now start actual tests

	// GetAll should not do anything, there's no data
	_, err = store.GetAll()()
	assert.ErrorIs(t, err, io.EOF)

	emailFoo := Email{
		Email:     "foo",
		SubToken:  "subTokenFoo",
		CreatedAt: now(),
	}

	// email isn't stored yet, shouldn't exist
	assertNotFound(t, emailFoo.Email)

	// Set an email, now it should exist
	assert.NoError(t, store.Set(emailFoo))
	assertGet(t, emailFoo)

	// Update the email with an unsub token
	emailFoo.UnsubToken = "unsubTokenFoo"
	emailFoo.VerifiedAt = now()
	assert.NoError(t, store.Set(emailFoo))
	assertGet(t, emailFoo)

	// GetAll should now only return that email
	iter := store.GetAll()
	gotEmail, err := iter()
	assert.NoError(t, err)
	assert.Equal(t, emailFoo, gotEmail)
	_, err = iter()
	assert.ErrorIs(t, err, io.EOF)

	// Delete the email, it should be gone
	assert.NoError(t, store.Delete(emailFoo.Email))
	assertNotFound(t, emailFoo.Email)
	_, err = store.GetAll()()
	assert.ErrorIs(t, err, io.EOF)
}
