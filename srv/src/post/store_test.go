package post

import (
	"sort"
	"strconv"
	"testing"
	"time"

	"github.com/mediocregopher/blog.mediocregopher.com/srv/cfg"
	"github.com/stretchr/testify/assert"
	"github.com/tilinna/clock"
)

func testPost(i int) Post {
	istr := strconv.Itoa(i)
	return Post{
		ID:          istr,
		Title:       istr,
		Description: istr,
		Body:        istr,
	}
}

type storeTestHarness struct {
	clock *clock.Mock
	store Store
}

func newStoreTestHarness(t *testing.T) storeTestHarness {

	var dataDir cfg.DataDir

	if err := dataDir.Init(); err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() { dataDir.Close() })

	clock := clock.NewMock(time.Now().UTC().Truncate(1 * time.Hour))

	store, err := NewStore(StoreParams{
		DataDir: dataDir,
	})
	assert.NoError(t, err)

	t.Cleanup(func() {
		assert.NoError(t, store.Close())
	})

	return storeTestHarness{
		clock: clock,
		store: store,
	}
}

func (h *storeTestHarness) testStoredPost(i int) StoredPost {
	post := testPost(i)
	return StoredPost{
		Post:        post,
		PublishedAt: h.clock.Now(),
	}
}

func TestStore(t *testing.T) {

	assertPostEqual := func(t *testing.T, exp, got StoredPost) {
		t.Helper()
		sort.Strings(exp.Tags)
		sort.Strings(got.Tags)
		assert.Equal(t, exp, got)
	}

	assertPostsEqual := func(t *testing.T, exp, got []StoredPost) {
		t.Helper()

		if !assert.Len(t, got, len(exp), "exp:%+v\ngot: %+v", exp, got) {
			return
		}

		for i := range exp {
			assertPostEqual(t, exp[i], got[i])
		}
	}

	t.Run("not_found", func(t *testing.T) {
		h := newStoreTestHarness(t)

		_, err := h.store.GetByID("foo")
		assert.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("set_get_delete", func(t *testing.T) {
		h := newStoreTestHarness(t)

		now := h.clock.Now().UTC()

		post := testPost(0)
		post.Tags = []string{"foo", "bar"}

		assert.NoError(t, h.store.Set(post, now))

		gotPost, err := h.store.GetByID(post.ID)
		assert.NoError(t, err)

		assertPostEqual(t, StoredPost{
			Post:        post,
			PublishedAt: now,
		}, gotPost)

		// we will now try updating the post on a different day, and ensure it
		// updates properly

		h.clock.Add(24 * time.Hour)
		newNow := h.clock.Now().UTC()

		post.Title = "something else"
		post.Series = "whatever"
		post.Body = "anything"
		post.Tags = []string{"bar", "baz"}

		assert.NoError(t, h.store.Set(post, newNow))

		gotPost, err = h.store.GetByID(post.ID)
		assert.NoError(t, err)

		assertPostEqual(t, StoredPost{
			Post:          post,
			PublishedAt:   now,
			LastUpdatedAt: newNow,
		}, gotPost)

		// delete the post, it should go away
		assert.NoError(t, h.store.Delete(post.ID))

		_, err = h.store.GetByID(post.ID)
		assert.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("get", func(t *testing.T) {
		h := newStoreTestHarness(t)

		now := h.clock.Now().UTC()

		posts := []StoredPost{
			h.testStoredPost(0),
			h.testStoredPost(1),
			h.testStoredPost(2),
			h.testStoredPost(3),
		}

		for _, post := range posts {
			assert.NoError(t, h.store.Set(post.Post, now))
		}

		gotPosts, hasMore, err := h.store.Get(0, 2)
		assert.NoError(t, err)
		assert.True(t, hasMore)
		assertPostsEqual(t, posts[:2], gotPosts)

		gotPosts, hasMore, err = h.store.Get(1, 2)
		assert.NoError(t, err)
		assert.False(t, hasMore)
		assertPostsEqual(t, posts[2:4], gotPosts)

		posts = append(posts, h.testStoredPost(4))
		assert.NoError(t, h.store.Set(posts[4].Post, now))

		gotPosts, hasMore, err = h.store.Get(1, 2)
		assert.NoError(t, err)
		assert.True(t, hasMore)
		assertPostsEqual(t, posts[2:4], gotPosts)

		gotPosts, hasMore, err = h.store.Get(2, 2)
		assert.NoError(t, err)
		assert.False(t, hasMore)
		assertPostsEqual(t, posts[4:], gotPosts)
	})

	t.Run("get_by_series", func(t *testing.T) {
		h := newStoreTestHarness(t)

		now := h.clock.Now().UTC()

		posts := []StoredPost{
			h.testStoredPost(0),
			h.testStoredPost(1),
			h.testStoredPost(2),
			h.testStoredPost(3),
		}

		posts[0].Series = "foo"
		posts[1].Series = "bar"
		posts[2].Series = "bar"

		for _, post := range posts {
			assert.NoError(t, h.store.Set(post.Post, now))
		}

		fooPosts, err := h.store.GetBySeries("foo")
		assert.NoError(t, err)
		assertPostsEqual(t, posts[:1], fooPosts)

		barPosts, err := h.store.GetBySeries("bar")
		assert.NoError(t, err)
		assertPostsEqual(t, posts[1:3], barPosts)

		bazPosts, err := h.store.GetBySeries("baz")
		assert.NoError(t, err)
		assert.Empty(t, bazPosts)
	})

	t.Run("get_by_tag", func(t *testing.T) {

		h := newStoreTestHarness(t)

		now := h.clock.Now().UTC()

		posts := []StoredPost{
			h.testStoredPost(0),
			h.testStoredPost(1),
			h.testStoredPost(2),
			h.testStoredPost(3),
		}

		posts[0].Tags = []string{"foo"}
		posts[1].Tags = []string{"foo", "bar"}
		posts[2].Tags = []string{"bar"}

		for _, post := range posts {
			assert.NoError(t, h.store.Set(post.Post, now))
		}

		fooPosts, err := h.store.GetByTag("foo")
		assert.NoError(t, err)
		assertPostsEqual(t, posts[:2], fooPosts)

		barPosts, err := h.store.GetByTag("bar")
		assert.NoError(t, err)
		assertPostsEqual(t, posts[1:3], barPosts)

		bazPosts, err := h.store.GetByTag("baz")
		assert.NoError(t, err)
		assert.Empty(t, bazPosts)
	})
}
