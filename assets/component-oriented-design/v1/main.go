package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"sort"
	"strconv"
	"sync"
	"time"
)

// Logger describes a simple component used for printing log lines.
type Logger interface {
	Printf(string, ...interface{})
}

////////////////////////////////////////////////////////////////////////////////
// The scoreboard component

// File wraps the standard os.File type.
type File interface {
	io.ReadWriter
	Truncate(int64) error
	Seek(int64, int) (int64, error)
}

// scoreboard loads player scores from a save file, tracks score updates, and
// periodically saves those scores back to the save file.
type scoreboard struct {
	file       File
	scoresM    map[string]int
	scoresLock sync.Mutex

	// this field will only be set in tests, and is used to synchronize with the
	// the for-select loop in saveLoop.
	saveLoopWaitCh chan struct{}
}

// newScoreboard initializes a scoreboard using scores saved in the given File
// (which may be empty). The scoreboard will rewrite the save file with the
// latest scores everytime saveTicker is written to.
func newScoreboard(file File, saveTicker <-chan time.Time, logger Logger) (*scoreboard, error) {
	fileBody, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("reading saved scored: %w", err)
	}

	scoresM := map[string]int{}
	if len(fileBody) > 0 {
		if err := json.Unmarshal(fileBody, &scoresM); err != nil {
			return nil, fmt.Errorf("decoding saved scores: %w", err)
		}
	}

	scoreboard := &scoreboard{
		file:           file,
		scoresM:        scoresM,
		saveLoopWaitCh: make(chan struct{}),
	}

	go scoreboard.saveLoop(saveTicker, logger)

	return scoreboard, nil
}

func (s *scoreboard) guessedCorrect(name string) int {
	s.scoresLock.Lock()
	defer s.scoresLock.Unlock()

	s.scoresM[name] += 1000
	return s.scoresM[name]
}

func (s *scoreboard) guessedIncorrect(name string) int {
	s.scoresLock.Lock()
	defer s.scoresLock.Unlock()

	s.scoresM[name] -= 1
	return s.scoresM[name]
}

func (s *scoreboard) scores() map[string]int {
	s.scoresLock.Lock()
	defer s.scoresLock.Unlock()

	scoresCp := map[string]int{}
	for name, score := range s.scoresM {
		scoresCp[name] = score
	}

	return scoresCp
}

func (s *scoreboard) save() error {
	scores := s.scores()
	if _, err := s.file.Seek(0, 0); err != nil {
		return fmt.Errorf("seeking to start of save file: %w", err)
	} else if err := s.file.Truncate(0); err != nil {
		return fmt.Errorf("truncating save file: %w", err)
	} else if err := json.NewEncoder(s.file).Encode(scores); err != nil {
		return fmt.Errorf("encoding scores to save file: %w", err)
	}
	return nil
}

func (s *scoreboard) saveLoop(ticker <-chan time.Time, logger Logger) {
	for {
		select {
		case <-ticker:
			if err := s.save(); err != nil {
				logger.Printf("error saving scoreboard to file: %v", err)
			}
		case <-s.saveLoopWaitCh:
			// test will unblock, nothing to do here.
		}
	}
}

////////////////////////////////////////////////////////////////////////////////
// The httpHandlers component

// Scoreboard describes the scoreboard component from the point of view of the
// httpHandler component (which only needs a subset of scoreboard's methods).
type Scoreboard interface {
	guessedCorrect(name string) int
	guessedIncorrect(name string) int
	scores() map[string]int
}

// RandSrc describes a randomness component which can produce random integers.
type RandSrc interface {
	Int() int
}

// httpHandlers implements the http.HandlerFuncs used by the httpServer.
type httpHandlers struct {
	scoreboard Scoreboard
	randSrc    RandSrc
	logger     Logger

	mux   *http.ServeMux
	n     int
	nLock sync.Mutex
}

func newHTTPHandlers(scoreboard Scoreboard, randSrc RandSrc, logger Logger) *httpHandlers {
	n := randSrc.Int()
	logger.Printf("first n is %v", n)

	httpHandlers := &httpHandlers{
		scoreboard: scoreboard,
		randSrc:    randSrc,
		logger:     logger,
		mux:        http.NewServeMux(),
		n:          n,
	}

	httpHandlers.mux.HandleFunc("/guess", httpHandlers.handleGuess)
	httpHandlers.mux.HandleFunc("/scores", httpHandlers.handleScores)

	return httpHandlers
}

func (h *httpHandlers) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(rw, r)
}

func (h *httpHandlers) handleGuess(rw http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "text/plain")

	name := r.FormValue("name")
	nStr := r.FormValue("n")
	if name == "" || nStr == "" {
		http.Error(rw, `"name" and "n" GET args are required`, http.StatusBadRequest)
		return
	}

	n, err := strconv.Atoi(nStr)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	h.nLock.Lock()
	defer h.nLock.Unlock()

	if h.n == n {
		newScore := h.scoreboard.guessedCorrect(name)
		h.n = h.randSrc.Int()
		h.logger.Printf("new n is %v", h.n)
		rw.WriteHeader(http.StatusOK)
		fmt.Fprintf(rw, "Correct! Your score is now %d\n", newScore)
		return
	}

	hint := "higher"
	if h.n < n {
		hint = "lower"
	}

	newScore := h.scoreboard.guessedIncorrect(name)
	rw.WriteHeader(http.StatusBadRequest)
	fmt.Fprintf(rw, "Try %s. Your score is now %d\n", hint, newScore)
}

func (h *httpHandlers) handleScores(rw http.ResponseWriter, r *http.Request) {
	r.Header.Set("Content-Type", "text/plain")

	h.nLock.Lock()
	defer h.nLock.Unlock()

	type scoreTup struct {
		name  string
		score int
	}

	scores := h.scoreboard.scores()
	scoresTups := make([]scoreTup, 0, len(scores))
	for name, score := range scores {
		scoresTups = append(scoresTups, scoreTup{name, score})
	}

	sort.Slice(scoresTups, func(i, j int) bool {
		return scoresTups[i].score > scoresTups[j].score
	})

	for _, scoresTup := range scoresTups {
		fmt.Fprintf(rw, "%s: %d\n", scoresTup.name, scoresTup.score)
	}
}

////////////////////////////////////////////////////////////////////////////////
// The httpServer component.

type httpServer struct {
	httpServer *http.Server
	errCh      chan error
}

func newHTTPServer(listener net.Listener, httpHandlers *httpHandlers, logger Logger) *httpServer {
	loggingHandler := http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		logger.Printf("HTTP request -> %s %s %s", ip, r.Method, r.URL.String())
		httpHandlers.ServeHTTP(rw, r)
	})

	server := &httpServer{
		httpServer: &http.Server{
			Handler: loggingHandler,
		},
		errCh: make(chan error, 1),
	}

	go func() {
		err := server.httpServer.Serve(listener)
		if errors.Is(err, http.ErrServerClosed) {
			err = nil
		}
		server.errCh <- err
	}()

	return server
}

////////////////////////////////////////////////////////////////////////////////
// main

const (
	saveFilePath = "./save.json"
	listenAddr   = ":8888"
	saveInterval = 5 * time.Second
)

func main() {
	logger := log.New(os.Stdout, "", log.LstdFlags)

	logger.Printf("opening scoreboard save file %q", saveFilePath)
	file, err := os.OpenFile(saveFilePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		logger.Fatalf("failed to open file %q: %v", saveFilePath, err)
	}

	saveTicker := time.NewTicker(saveInterval)
	randSrc := rand.New(rand.NewSource(time.Now().UnixNano()))

	logger.Printf("initializing scoreboard")
	scoreboard, err := newScoreboard(file, saveTicker.C, logger)
	if err != nil {
		logger.Fatalf("failed to initialize scoreboard: %v", err)
	}

	logger.Printf("listening on %q", listenAddr)
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		logger.Fatalf("failed to listen on %q: %v", listenAddr, err)
	}

	logger.Printf("setting up HTTP handlers")
	httpHandlers := newHTTPHandlers(scoreboard, randSrc, logger)

	logger.Printf("serving HTTP requests")
	newHTTPServer(listener, httpHandlers, logger)

	logger.Printf("initialization done")
	select {} // block forever
}
