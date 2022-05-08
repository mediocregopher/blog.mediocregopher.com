package post

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

type assetTestHarness struct {
	store AssetStore
}

func newAssetTestHarness(t *testing.T) assetTestHarness {

	db := NewInMemSQLDB()
	t.Cleanup(func() { db.Close() })

	store := NewAssetStore(db)

	return assetTestHarness{
		store: store,
	}
}

func (h *assetTestHarness) assertGet(t *testing.T, exp, id string) {
	t.Helper()
	buf := new(bytes.Buffer)
	err := h.store.Get(id, buf)
	assert.NoError(t, err)
	assert.Equal(t, exp, buf.String())
}

func (h *assetTestHarness) assertNotFound(t *testing.T, id string) {
	t.Helper()
	err := h.store.Get(id, io.Discard)
	assert.ErrorIs(t, ErrAssetNotFound, err)
}

func TestAssetStore(t *testing.T) {

	h := newAssetTestHarness(t)

	h.assertNotFound(t, "foo")
	h.assertNotFound(t, "bar")

	err := h.store.Set("foo", bytes.NewBufferString("FOO"))
	assert.NoError(t, err)

	h.assertGet(t, "FOO", "foo")
	h.assertNotFound(t, "bar")

	err = h.store.Set("foo", bytes.NewBufferString("FOOFOO"))
	assert.NoError(t, err)

	h.assertGet(t, "FOOFOO", "foo")
	h.assertNotFound(t, "bar")

	assert.NoError(t, h.store.Delete("foo"))
	h.assertNotFound(t, "foo")
	h.assertNotFound(t, "bar")

	assert.NoError(t, h.store.Delete("bar"))
	h.assertNotFound(t, "foo")
	h.assertNotFound(t, "bar")
}
