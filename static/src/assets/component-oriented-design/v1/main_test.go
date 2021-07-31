package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

type nullLogger struct{}

func (nullLogger) Printf(string, ...interface{}) {}

////////////////////////////////////////////////////////////////////////////////
// Test scoreboard component

type fileStub struct {
	*bytes.Buffer
}

func newFileStub(init string) *fileStub {
	return &fileStub{Buffer: bytes.NewBufferString(init)}
}

func (fs *fileStub) Truncate(i int64) error {
	fs.Buffer.Truncate(int(i))
	return nil
}

func (fs *fileStub) Seek(i int64, whence int) (int64, error) {
	return i, nil
}

func TestScoreboard(t *testing.T) {
	newScoreboard := func(t *testing.T, fileStub *fileStub, saveTicker <-chan time.Time) *scoreboard {
		t.Helper()
		scoreboard, err := newScoreboard(fileStub, saveTicker, nullLogger{})
		if err != nil {
			t.Errorf("unexpected error checking saved scored: %v", err)
		}
		return scoreboard
	}

	assertScores := func(t *testing.T, expScores, gotScores map[string]int) {
		t.Helper()
		if !reflect.DeepEqual(expScores, gotScores) {
			t.Errorf("expected scores of %+v, but instead got %+v", expScores, gotScores)
		}
	}

	assertSavedScores := func(t *testing.T, expScores map[string]int, fileStub *fileStub) {
		t.Helper()
		fileStubCp := newFileStub(fileStub.String())
		tmpScoreboard := newScoreboard(t, fileStubCp, nil)
		assertScores(t, expScores, tmpScoreboard.scores())
	}

	t.Run("loading", func(t *testing.T) {
		// make sure loading scoreboards with various file contents works
		assertSavedScores(t, map[string]int{}, newFileStub(""))
		assertSavedScores(t, map[string]int{"foo": 1}, newFileStub(`{"foo":1}`))
		assertSavedScores(t, map[string]int{"foo": 1, "bar": -2}, newFileStub(`{"foo":1,"bar":-2}`))
	})

	t.Run("tracking", func(t *testing.T) {
		scoreboard := newScoreboard(t, newFileStub(""), nil)
		assertScores(t, map[string]int{}, scoreboard.scores()) // sanity check

		scoreboard.guessedCorrect("foo")
		assertScores(t, map[string]int{"foo": 1000}, scoreboard.scores())

		scoreboard.guessedIncorrect("bar")
		assertScores(t, map[string]int{"foo": 1000, "bar": -1}, scoreboard.scores())

		scoreboard.guessedIncorrect("foo")
		assertScores(t, map[string]int{"foo": 999, "bar": -1}, scoreboard.scores())
	})

	t.Run("saving", func(t *testing.T) {
		// this test tests scoreboard's periodic save feature using a ticker
		// channel which will be written to manually. The saveLoopWaitCh is used
		// here to ensure that each ticker has been fully processed.
		ticker := make(chan time.Time)
		fileStub := newFileStub("")
		scoreboard := newScoreboard(t, fileStub, ticker)

		tick := func() {
			ticker <- time.Time{}
			scoreboard.saveLoopWaitCh <- struct{}{}
		}

		// this should not effect the save file at first
		scoreboard.guessedCorrect("foo")
		assertSavedScores(t, map[string]int{}, fileStub)

		// after the ticker the new score should get saved
		tick()
		assertSavedScores(t, map[string]int{"foo": 1000}, fileStub)

		// ticker again after no changes should save the same thing as before
		tick()
		assertSavedScores(t, map[string]int{"foo": 1000}, fileStub)

		// buffer a bunch of changes, shouldn't get saved till after tick
		scoreboard.guessedCorrect("foo")
		scoreboard.guessedCorrect("bar")
		scoreboard.guessedCorrect("bar")
		assertSavedScores(t, map[string]int{"foo": 1000}, fileStub)
		tick()
		assertSavedScores(t, map[string]int{"foo": 2000, "bar": 2000}, fileStub)
	})
}

////////////////////////////////////////////////////////////////////////////////
// Test httpHandlers component

type mockScoreboard map[string]int

func (mockScoreboard) guessedCorrect(name string) int { return 1 }

func (mockScoreboard) guessedIncorrect(name string) int { return -1 }

func (m mockScoreboard) scores() map[string]int { return m }

type mockRandSrc struct{}

func (m mockRandSrc) Int() int { return 666 }

func TestHTTPHandlers(t *testing.T) {
	mockScoreboard := mockScoreboard{"foo": 1, "bar": 2}
	httpHandlers := newHTTPHandlers(mockScoreboard, mockRandSrc{}, nullLogger{})

	assertRequest := func(t *testing.T, expCode int, expBody string, r *http.Request) {
		t.Helper()
		rw := httptest.NewRecorder()
		httpHandlers.ServeHTTP(rw, r)
		if rw.Code != expCode {
			t.Errorf("expected HTTP response code %d, got %d", expCode, rw.Code)
		} else if rw.Body.String() != expBody {
			t.Errorf("expected HTTP response body %q, got %q", expBody, rw.Body.String())
		}
	}

	r := httptest.NewRequest("GET", "/guess?name=foo&n=665", nil)
	assertRequest(t, 400, "Try higher. Your score is now -1\n", r)

	r = httptest.NewRequest("GET", "/guess?name=foo&n=667", nil)
	assertRequest(t, 400, "Try lower. Your score is now -1\n", r)

	r = httptest.NewRequest("GET", "/guess?name=foo&n=666", nil)
	assertRequest(t, 200, "Correct! Your score is now 1\n", r)

	r = httptest.NewRequest("GET", "/scores", nil)
	assertRequest(t, 200, "bar: 2\nfoo: 1\n", r)
}

////////////////////////////////////////////////////////////////////////////////
//
// httpServer is NOT tested, for the following reasons:
// * It depends on a `net.Listener`, which is not trivial to mock.
// * It does very little besides passing an httpHandlers along to an http.Server
//   and managing cleanup.
// * It isn't likely to be changed often.
// * If it were to break it would be very apparent in subsequent testing stages.
//
