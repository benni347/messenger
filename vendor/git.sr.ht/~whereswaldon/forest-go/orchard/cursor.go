package orchard

import (
	"fmt"
	"io"

	"git.sr.ht/~whereswaldon/forest-go"
	"git.sr.ht/~whereswaldon/forest-go/orchard/internal/codec"
	bolt "go.etcd.io/bbolt"
)

// ByteCursor is a cursor that iterates over []byte kv pairs.
type ByteCursor interface {
	First() ([]byte, []byte)
	Next() ([]byte, []byte)
}

// indexCursor uses an index bucket to iterate over the parent bucket.
// Collided index entries are flattened in the order they were inserted.
type indexCursor struct {
	// Parent bucket contains the real values being indexed.
	Parent *bolt.Bucket
	// Index bucket contains the index values, which are some arbitrary key
	// mapped to entries in the Parent bucket.
	Index *bolt.Bucket
	// Reverse iterates backward, if true.
	Reverse bool

	// cursor holds the current state of the bucket cursor.
	cursor *bolt.Cursor
	// collisions holds the current state of the collision bucket cursor.
	// If collisions is not nil, we are iterating over collided entries for a
	// given key.
	collisions *bolt.Cursor
}

// First returns the first k,v pair.
// This method initializes the underlying cursor and must be called before calls
// to `indexCursor.Next`.
func (c *indexCursor) First() ([]byte, []byte) {
	if c.Index == nil {
		return nil, nil
	}
	c.cursor = c.Index.Cursor()
	k, v := c.first(c.cursor)
	if k == nil {
		return nil, nil
	}
	if v != nil {
		id := v
		return id, c.Parent.Get(id)
	}
	if b := c.Index.Bucket(k); b != nil {
		c.collisions = b.Cursor()
		_, id := c.first(c.collisions)
		return id, c.Parent.Get(id)
	}
	return nil, nil
}

// Next returns the next k,v pair until k is nil.
func (c *indexCursor) Next() ([]byte, []byte) {
	var cursor = c.cursor
	if c.collisions != nil {
		cursor = c.collisions
	}
	k, v := c.next(cursor)
	if k == nil {
		if c.collisions != nil {
			c.collisions = nil
			_, id := c.next(c.cursor)
			return id, c.Parent.Get(id)
		}
		return nil, nil
	}
	if v != nil {
		id := v
		return id, c.Parent.Get(id)
	}
	if b := c.Index.Bucket(k); b != nil {
		c.collisions = b.Cursor()
		_, id := c.first(c.collisions)
		return id, c.Parent.Get(id)
	}
	return nil, nil
}

func (c *indexCursor) next(cursor *bolt.Cursor) ([]byte, []byte) {
	if c.Reverse {
		return cursor.Prev()
	}
	return cursor.Next()
}

func (c *indexCursor) first(cursor *bolt.Cursor) ([]byte, []byte) {
	if c.Reverse {
		return cursor.Last()
	}
	return cursor.First()
}

/*
	NOTE(jfm): The following cursor types are structurally identical, but with different return types.
*/

// ReplyCursor wraps byte cursor and decodes into `forest.Reply` values.
type ReplyCursor struct {
	Inner ByteCursor
	Codec *codec.Default
}

// First returns the first reply or EOF if there are no.
func (c *ReplyCursor) First() (r forest.Reply, err error) {
	k, v := c.Inner.First()
	if k == nil {
		return r, io.EOF
	}
	data := make([]byte, len(v))
	copy(data, v)
	if err := c.Codec.Decode(data, &r); err != nil {
		return r, fmt.Errorf("decoding reply: %w", err)
	}
	return r, nil
}

// Next returns the next reply or EOF if there are no more replies.
func (c *ReplyCursor) Next() (r forest.Reply, err error) {
	k, v := c.Inner.Next()
	if k == nil {
		return r, io.EOF
	}
	data := make([]byte, len(v))
	copy(data, v)
	if err := c.Codec.Decode(data, &r); err != nil {
		return r, fmt.Errorf("decoding reply: %w", err)
	}
	return r, nil
}

// IdentityCursor wraps byte cursor and decodes into `forest.Identity` values.
type IdentityCursor struct {
	Inner ByteCursor
	Codec *codec.Default
}

// First returns the first identity or EOF if there are no.
func (c *IdentityCursor) First() (r forest.Identity, err error) {
	k, v := c.Inner.First()
	if k == nil {
		return r, io.EOF
	}
	data := make([]byte, len(v))
	copy(data, v)
	if err := c.Codec.Decode(data, &r); err != nil {
		return r, fmt.Errorf("decoding identity: %w", err)
	}
	return r, nil
}

// Next returns the next identity or EOF if there are no more.
func (c *IdentityCursor) Next() (r forest.Identity, err error) {
	k, v := c.Inner.Next()
	if k == nil {
		return r, io.EOF
	}
	data := make([]byte, len(v))
	copy(data, v)
	if err := c.Codec.Decode(data, &r); err != nil {
		return r, fmt.Errorf("decoding identity: %w", err)
	}
	return r, nil
}

// CommunityCursor wraps byte cursor and decodes into `forest.Community` values.
type CommunityCursor struct {
	Inner ByteCursor
	Codec *codec.Default
}

// First returns the first community or EOF if there are no.
func (c *CommunityCursor) First() (r forest.Community, err error) {
	k, v := c.Inner.First()
	if k == nil {
		return r, io.EOF
	}
	data := make([]byte, len(v))
	copy(data, v)
	if err := c.Codec.Decode(data, &r); err != nil {
		return r, fmt.Errorf("decoding community: %w", err)
	}
	return r, nil
}

// Next returns the next community or EOF if there are no more.
func (c *CommunityCursor) Next() (r forest.Community, err error) {
	k, v := c.Inner.Next()
	if k == nil {
		return r, io.EOF
	}
	data := make([]byte, len(v))
	copy(data, v)
	if err := c.Codec.Decode(data, &r); err != nil {
		return r, fmt.Errorf("decoding community: %w", err)
	}
	return r, nil
}
