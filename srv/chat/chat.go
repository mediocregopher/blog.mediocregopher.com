// Package chat implements a simple chatroom system.
package chat

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/mediocregopher/mediocre-go-lib/v2/mctx"
	"github.com/mediocregopher/mediocre-go-lib/v2/mlog"
	"github.com/mediocregopher/radix/v4"
)

// ErrInvalidArg is returned from methods in this package when a call fails due
// to invalid input.
type ErrInvalidArg struct {
	Err error
}

func (e ErrInvalidArg) Error() string {
	return fmt.Sprintf("invalid argument: %v", e.Err)
}

var (
	errInvalidMessageID = ErrInvalidArg{Err: errors.New("invalid Message ID")}
)

// Message describes a message which has been posted to a Room.
type Message struct {
	ID        string `json:"id"`
	UserID    UserID `json:"userID"`
	Body      string `json:"body"`
	CreatedAt int64  `json:"createdAt,omitempty"`
}

func msgFromStreamEntry(entry radix.StreamEntry) (Message, error) {

	// NOTE this should probably be a shortcut in radix
	var bodyStr string
	for _, field := range entry.Fields {
		if field[0] == "json" {
			bodyStr = field[1]
			break
		}
	}

	if bodyStr == "" {
		return Message{}, errors.New("no 'json' field")
	}

	var msg Message
	if err := json.Unmarshal([]byte(bodyStr), &msg); err != nil {
		return Message{}, fmt.Errorf(
			"json unmarshaling body %q: %w", bodyStr, err,
		)
	}

	msg.ID = entry.ID.String()
	msg.CreatedAt = int64(entry.ID.Time / 1000)
	return msg, nil
}

// MessageIterator returns a sequence of Messages which may or may not be
// unbounded.
type MessageIterator interface {

	// Next blocks until it returns the next Message in the sequence, or the
	// context error if the context is cancelled, or io.EOF if the sequence has
	// been exhausted.
	Next(context.Context) (Message, error)

	// Close should always be called once Next has returned an error or the
	// MessageIterator will no longer be used.
	Close() error
}

// HistoryOpts are passed into Room's History method in order to affect its
// result. All fields are optional.
type HistoryOpts struct {
	Limit  int    // defaults to, and is capped at, 100.
	Cursor string // If not given then the most recent Messages are returned.
}

func (o HistoryOpts) sanitize() (HistoryOpts, error) {
	if o.Limit <= 0 || o.Limit > 100 {
		o.Limit = 100
	}

	if o.Cursor != "" {
		id, err := parseStreamEntryID(o.Cursor)
		if err != nil {
			return HistoryOpts{}, fmt.Errorf("parsing Cursor: %w", err)
		}
		o.Cursor = id.String()
	}

	return o, nil
}

// Room implements functionality related to a single, unique chat room.
type Room interface {

	// Append accepts a new Message and stores it at the end of the room's
	// history. The original Message is returned with any relevant fields (e.g.
	// ID) updated.
	Append(context.Context, Message) (Message, error)

	// Returns a cursor and the list of historical Messages in time descending
	// order. The cursor can be passed into the next call to History to receive
	// the next set of Messages.
	History(context.Context, HistoryOpts) (string, []Message, error)

	// Listen returns a MessageIterator which will return all Messages appended
	// to the Room since the given ID. Once all existing messages are iterated
	// through then the MessageIterator will begin blocking until a new Message
	// is posted.
	Listen(ctx context.Context, sinceID string) (MessageIterator, error)

	// Delete deletes a Message from the Room.
	Delete(ctx context.Context, id string) error

	// Close is used to clean up all resources created by the Room.
	Close() error
}

// RoomParams are used to instantiate a new Room. All fields are required unless
// otherwise noted.
type RoomParams struct {
	Logger      *mlog.Logger
	Redis       radix.Client
	ID          string
	MaxMessages int
}

func (p RoomParams) streamKey() string {
	return fmt.Sprintf("chat:{%s}:stream", p.ID)
}

type room struct {
	params RoomParams

	closeCtx    context.Context
	closeCancel context.CancelFunc
	wg          sync.WaitGroup

	listeningL      sync.Mutex
	listening       map[chan Message]struct{}
	listeningLastID radix.StreamEntryID
}

// NewRoom initializes and returns a new Room instance.
func NewRoom(ctx context.Context, params RoomParams) (Room, error) {

	params.Logger = params.Logger.WithNamespace("chat-room")

	r := &room{
		params:    params,
		listening: map[chan Message]struct{}{},
	}

	r.closeCtx, r.closeCancel = context.WithCancel(context.Background())

	// figure out the most recent message, if any.
	lastEntryID, err := r.mostRecentMsgID(ctx)
	if err != nil {
		return nil, fmt.Errorf("discovering most recent entry ID in stream: %w", err)
	}
	r.listeningLastID = lastEntryID

	r.wg.Add(1)
	go func() {
		defer r.wg.Done()
		r.readStreamLoop(r.closeCtx)
	}()

	return r, nil
}

func (r *room) Close() error {
	r.closeCancel()
	r.wg.Wait()
	return nil
}

func (r *room) mostRecentMsgID(ctx context.Context) (radix.StreamEntryID, error) {

	var entries []radix.StreamEntry
	err := r.params.Redis.Do(ctx, radix.Cmd(
		&entries,
		"XREVRANGE", r.params.streamKey(), "+", "-", "COUNT", "1",
	))

	if err != nil || len(entries) == 0 {
		return radix.StreamEntryID{}, err
	}

	return entries[0].ID, nil
}

func (r *room) Append(ctx context.Context, msg Message) (Message, error) {
	msg.ID = "" // just in case

	b, err := json.Marshal(msg)
	if err != nil {
		return Message{}, fmt.Errorf("json marshaling Message: %w", err)
	}

	key := r.params.streamKey()
	maxLen := strconv.Itoa(r.params.MaxMessages)
	body := string(b)

	var id radix.StreamEntryID

	err = r.params.Redis.Do(ctx, radix.Cmd(
		&id, "XADD", key, "MAXLEN", "=", maxLen, "*", "json", body,
	))

	if err != nil {
		return Message{}, fmt.Errorf("posting message to redis: %w", err)
	}

	msg.ID = id.String()
	msg.CreatedAt = int64(id.Time / 1000)
	return msg, nil
}

const zeroCursor = "0-0"

func (r *room) History(ctx context.Context, opts HistoryOpts) (string, []Message, error) {
	opts, err := opts.sanitize()
	if err != nil {
		return "", nil, err
	}

	key := r.params.streamKey()
	end := opts.Cursor
	if end == "" {
		end = "+"
	}
	start := "-"
	count := strconv.Itoa(opts.Limit)

	msgs := make([]Message, 0, opts.Limit)
	streamEntries := make([]radix.StreamEntry, 0, opts.Limit)

	err = r.params.Redis.Do(ctx, radix.Cmd(
		&streamEntries,
		"XREVRANGE", key, end, start, "COUNT", count,
	))

	if err != nil {
		return "", nil, fmt.Errorf("calling XREVRANGE: %w", err)
	}

	var oldestEntryID radix.StreamEntryID

	for _, entry := range streamEntries {
		oldestEntryID = entry.ID

		msg, err := msgFromStreamEntry(entry)
		if err != nil {
			return "", nil, fmt.Errorf(
				"parsing stream entry %q: %w", entry.ID, err,
			)
		}
		msgs = append(msgs, msg)
	}

	if len(msgs) < opts.Limit {
		return zeroCursor, msgs, nil
	}

	cursor := oldestEntryID.Prev()
	return cursor.String(), msgs, nil
}

func (r *room) readStream(ctx context.Context) error {

	r.listeningL.Lock()
	lastEntryID := r.listeningLastID
	r.listeningL.Unlock()

	redisAddr := r.params.Redis.Addr()
	redisConn, err := radix.Dial(ctx, redisAddr.Network(), redisAddr.String())
	if err != nil {
		return fmt.Errorf("creating redis connection: %w", err)
	}
	defer redisConn.Close()

	streamReader := (radix.StreamReaderConfig{}).New(
		redisConn,
		map[string]radix.StreamConfig{
			r.params.streamKey(): {After: lastEntryID},
		},
	)

	for {
		dlCtx, dlCtxCancel := context.WithTimeout(ctx, 10*time.Second)
		_, streamEntry, err := streamReader.Next(dlCtx)
		dlCtxCancel()

		if errors.Is(err, radix.ErrNoStreamEntries) {
			continue
		} else if err != nil {
			return fmt.Errorf("fetching next entry from stream: %w", err)
		}

		msg, err := msgFromStreamEntry(streamEntry)
		if err != nil {
			return fmt.Errorf("parsing stream entry %q: %w", streamEntry, err)
		}

		r.listeningL.Lock()

		var dropped int
		for ch := range r.listening {
			select {
			case ch <- msg:
			default:
				dropped++
			}
		}

		if dropped > 0 {
			ctx := mctx.Annotate(ctx, "msgID", msg.ID, "dropped", dropped)
			r.params.Logger.WarnString(ctx, "some listening channels full, messages dropped")
		}

		r.listeningLastID = streamEntry.ID

		r.listeningL.Unlock()
	}
}

func (r *room) readStreamLoop(ctx context.Context) {
	for {
		err := r.readStream(ctx)
		if errors.Is(err, context.Canceled) {
			return
		} else if err != nil {
			r.params.Logger.Error(ctx, "reading from redis stream", err)
		}
	}
}

type listenMsgIterator struct {
	ch           <-chan Message
	missedMsgs   []Message
	sinceEntryID radix.StreamEntryID
	cleanup      func()
}

func (i *listenMsgIterator) Next(ctx context.Context) (Message, error) {

	if len(i.missedMsgs) > 0 {
		msg := i.missedMsgs[0]
		i.missedMsgs = i.missedMsgs[1:]
		return msg, nil
	}

	for {
		select {
		case <-ctx.Done():
			return Message{}, ctx.Err()
		case msg := <-i.ch:

			entryID, err := parseStreamEntryID(msg.ID)
			if err != nil {
				return Message{}, fmt.Errorf("parsing Message ID %q: %w", msg.ID, err)

			} else if !i.sinceEntryID.Before(entryID) {
				// this can happen if someone Appends a Message at the same time
				// as another calls Listen. The Listener might have already seen
				// the Message by calling History prior to the stream reader
				// having processed it and updating listeningLastID.
				continue
			}

			return msg, nil
		}
	}
}

func (i *listenMsgIterator) Close() error {
	i.cleanup()
	return nil
}

func (r *room) Listen(
	ctx context.Context, sinceID string,
) (
	MessageIterator, error,
) {

	var sinceEntryID radix.StreamEntryID

	if sinceID != "" {
		var err error
		if sinceEntryID, err = parseStreamEntryID(sinceID); err != nil {
			return nil, fmt.Errorf("parsing sinceID: %w", err)
		}
	}

	ch := make(chan Message, 32)

	r.listeningL.Lock()
	lastEntryID := r.listeningLastID
	r.listening[ch] = struct{}{}
	r.listeningL.Unlock()

	cleanup := func() {
		r.listeningL.Lock()
		defer r.listeningL.Unlock()
		delete(r.listening, ch)
	}

	key := r.params.streamKey()
	start := sinceEntryID.Next().String()
	end := "+"
	if lastEntryID != (radix.StreamEntryID{}) {
		end = lastEntryID.String()
	}

	var streamEntries []radix.StreamEntry

	err := r.params.Redis.Do(ctx, radix.Cmd(
		&streamEntries,
		"XRANGE", key, start, end,
	))

	if err != nil {
		cleanup()
		return nil, fmt.Errorf("retrieving missed stream entries: %w", err)
	}

	missedMsgs := make([]Message, len(streamEntries))

	for i := range streamEntries {

		msg, err := msgFromStreamEntry(streamEntries[i])
		if err != nil {
			cleanup()
			return nil, fmt.Errorf(
				"parsing stream entry %q: %w", streamEntries[i].ID, err,
			)
		}

		missedMsgs[i] = msg
	}

	return &listenMsgIterator{
		ch:           ch,
		missedMsgs:   missedMsgs,
		sinceEntryID: sinceEntryID,
		cleanup:      cleanup,
	}, nil
}

func (r *room) Delete(ctx context.Context, id string) error {
	return r.params.Redis.Do(ctx, radix.Cmd(
		nil, "XDEL", r.params.streamKey(), id,
	))
}
