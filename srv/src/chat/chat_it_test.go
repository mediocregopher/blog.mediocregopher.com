//go:build integration

package chat

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/mediocregopher/mediocre-go-lib/v2/mlog"
	"github.com/mediocregopher/radix/v4"
	"github.com/stretchr/testify/assert"
)

const roomTestHarnessMaxMsgs = 10

type roomTestHarness struct {
	ctx     context.Context
	room    Room
	allMsgs []Message
}

func (h *roomTestHarness) newMsg(t *testing.T) Message {
	msg, err := h.room.Append(h.ctx, Message{
		UserID: UserID{
			Name: uuid.New().String(),
			Hash: "0000",
		},
		Body: uuid.New().String(),
	})
	assert.NoError(t, err)

	t.Logf("appended message %s", msg.ID)

	h.allMsgs = append([]Message{msg}, h.allMsgs...)

	if len(h.allMsgs) > roomTestHarnessMaxMsgs {
		h.allMsgs = h.allMsgs[:roomTestHarnessMaxMsgs]
	}

	return msg
}

func newRoomTestHarness(t *testing.T) *roomTestHarness {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	t.Cleanup(cancel)

	redis, err := radix.Dial(ctx, "tcp", "127.0.0.1:6379")
	assert.NoError(t, err)
	t.Cleanup(func() { redis.Close() })

	roomParams := RoomParams{
		Logger:      mlog.NewLogger(nil),
		Redis:       redis,
		ID:          uuid.New().String(),
		MaxMessages: roomTestHarnessMaxMsgs,
	}

	t.Logf("creating test Room %q", roomParams.ID)
	room, err := NewRoom(ctx, roomParams)
	assert.NoError(t, err)

	t.Cleanup(func() {
		err := redis.Do(context.Background(), radix.Cmd(
			nil, "DEL", roomParams.streamKey(),
		))
		assert.NoError(t, err)
	})

	return &roomTestHarness{ctx: ctx, room: room}
}

func TestRoom(t *testing.T) {
	t.Run("history", func(t *testing.T) {

		tests := []struct {
			numMsgs int
			limit   int
		}{
			{numMsgs: 0, limit: 1},
			{numMsgs: 1, limit: 1},
			{numMsgs: 2, limit: 1},
			{numMsgs: 2, limit: 10},
			{numMsgs: 9, limit: 2},
			{numMsgs: 9, limit: 3},
			{numMsgs: 9, limit: 4},
			{numMsgs: 15, limit: 3},
		}

		for i, test := range tests {
			t.Run(strconv.Itoa(i), func(t *testing.T) {
				t.Logf("test: %+v", test)

				h := newRoomTestHarness(t)

				for j := 0; j < test.numMsgs; j++ {
					h.newMsg(t)
				}

				var gotMsgs []Message
				var cursor string

				for {

					var msgs []Message
					var err error
					cursor, msgs, err = h.room.History(h.ctx, HistoryOpts{
						Cursor: cursor,
						Limit:  test.limit,
					})

					assert.NoError(t, err)
					assert.NotEmpty(t, cursor)

					if len(msgs) == 0 {
						break
					}

					gotMsgs = append(gotMsgs, msgs...)
				}

				assert.Equal(t, h.allMsgs, gotMsgs)
			})
		}
	})

	assertNextMsg := func(
		t *testing.T, expMsg Message,
		ctx context.Context, it MessageIterator,
	) {
		t.Helper()
		gotMsg, err := it.Next(ctx)
		assert.NoError(t, err)
		assert.Equal(t, expMsg, gotMsg)
	}

	t.Run("listen/already_populated", func(t *testing.T) {
		h := newRoomTestHarness(t)

		msgA, msgB, msgC := h.newMsg(t), h.newMsg(t), h.newMsg(t)
		_ = msgA
		_ = msgB

		itFoo, err := h.room.Listen(h.ctx, msgC.ID)
		assert.NoError(t, err)
		defer itFoo.Close()

		itBar, err := h.room.Listen(h.ctx, msgA.ID)
		assert.NoError(t, err)
		defer itBar.Close()

		msgD := h.newMsg(t)

		// itBar should get msgB and msgC before anything else.
		assertNextMsg(t, msgB, h.ctx, itBar)
		assertNextMsg(t, msgC, h.ctx, itBar)

		// now both iterators should give msgD
		assertNextMsg(t, msgD, h.ctx, itFoo)
		assertNextMsg(t, msgD, h.ctx, itBar)

		// timeout should be honored
		{
			timeoutCtx, timeoutCancel := context.WithTimeout(h.ctx, 1*time.Second)
			_, errFoo := itFoo.Next(timeoutCtx)
			_, errBar := itBar.Next(timeoutCtx)
			timeoutCancel()

			assert.ErrorIs(t, errFoo, context.DeadlineExceeded)
			assert.ErrorIs(t, errBar, context.DeadlineExceeded)
		}

		// new message should work
		{
			expMsg := h.newMsg(t)

			timeoutCtx, timeoutCancel := context.WithTimeout(h.ctx, 1*time.Second)
			gotFooMsg, errFoo := itFoo.Next(timeoutCtx)
			gotBarMsg, errBar := itBar.Next(timeoutCtx)
			timeoutCancel()

			assert.Equal(t, expMsg, gotFooMsg)
			assert.NoError(t, errFoo)
			assert.Equal(t, expMsg, gotBarMsg)
			assert.NoError(t, errBar)
		}

	})

	t.Run("listen/empty", func(t *testing.T) {
		h := newRoomTestHarness(t)

		it, err := h.room.Listen(h.ctx, "")
		assert.NoError(t, err)
		defer it.Close()

		msg := h.newMsg(t)
		assertNextMsg(t, msg, h.ctx, it)
	})
}
