package orchard

import (
	"bytes"
	"fmt"

	"git.sr.ht/~whereswaldon/forest-go/fields"
	bolt "go.etcd.io/bbolt"
)

// index maps keys to node IDs for accelerated queries.
// On collision, all values for the collided key are placed into a bucket.
// Must be used inside a transaction.
type index struct {
	*bolt.Bucket
}

// Put a kv pair into the index.
// If an entry for k exists, all values are inserted into a bucket for that
// key.
func (idx index) Put(k []byte, v []byte) error {
	nextID := func(b *bolt.Bucket) ([]byte, error) {
		// Generate a unique ID for this sub-entry.
		// ID is not used other than by the database for purposes of
		// sorting and iterating.
		kint, err := b.NextSequence()
		if err != nil {
			return nil, err
		}
		k := make([]byte, 8)
		fields.MultiByteSerializationOrder.PutUint64(k, kint)
		return k, nil
	}
	if collision := idx.Bucket.Get(k); collision != nil {
		// Move any existing value into a collision bucket.
		if err := idx.Bucket.Delete(k); err != nil {
			return fmt.Errorf("failed deleting colliding entry: %w", err)
		}
		collisions, err := idx.Bucket.CreateBucketIfNotExists(k)
		if err != nil {
			return fmt.Errorf("failed creating new bucket for key: %w", err)
		}
		k, err := nextID(collisions)
		if err != nil {
			return fmt.Errorf("failed generating next id: %w", err)
		}
		if err := collisions.Put(k, collision); err != nil {
			return fmt.Errorf("failed inserting into new bucket: %w", err)
		}
	}
	if collisions := idx.Bucket.Bucket(k); collisions != nil {
		// Place the new value into the collisions bucket.
		k, err := nextID(collisions)
		if err != nil {
			return fmt.Errorf("failed generating next id: %w", err)
		}
		if err := collisions.Put(k, v); err != nil {
			return fmt.Errorf("failed inserting into collision bucket: %w", err)
		}
		return nil
	}
	// No collision, place value into the root bucket.

	if err := idx.Bucket.Put(k, v); err != nil {
		return fmt.Errorf("failed inserting into root bucket: %w", err)
	}
	return nil
}

func (idx index) Delete(k []byte) error {
	bucket := idx.Bucket.Bucket(k)
	if bucket == nil {
		// the key has one (or zero values)
		if err := idx.Bucket.Delete(k); err != nil {
			return fmt.Errorf("failed deleting key %v: %w", k, err)
		}
		return nil
	}
	if err := idx.Bucket.DeleteBucket(k); err != nil {
		return fmt.Errorf("failed deleting bucket for key %v: %w", k, err)
	}
	return nil
}

// DeleteByValue removes the value v from the key k. Note that if k is associated with
// multiple values, only the first occurrence of v within k will be removed. If
// v is the only value associated with k, k will be deleted.
func (idx index) DeleteByValue(k, v []byte) (err error) {
	bucket := idx.Bucket.Bucket(k)
	if bucket == nil {
		// the key has one (or zero values)
		if err := idx.Bucket.Delete(k); err != nil {
			return fmt.Errorf("failed deleting key %v: %w", k, err)
		}
		return nil
	}
	defer func() {
		if err == nil && bucket.Stats().KeyN == 0 {
			if deleteErr := idx.Bucket.DeleteBucket(k); deleteErr != nil {
				err = fmt.Errorf("failed deleting empty bucket for key %v: %w", k, deleteErr)
			}
		}
	}()
	// the key maps to a bucket from which v must be deleted
	c := bucket.Cursor()
	for key, value := c.First(); key != nil; key, value = c.Next() {
		if bytes.Equal(v, value) {
			if err := c.Delete(); err != nil {
				return fmt.Errorf("failed deleting value %v from bucket: %w", v, err)
			}
			return nil
		}
	}
	return fmt.Errorf("key %v did not contain value %v", k, v)
}
