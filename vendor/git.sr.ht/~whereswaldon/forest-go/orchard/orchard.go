// Package orchard implements a boltdb backed on-disk node store, satisfying the
// forest.Store interface.
//
// This database is a single file that serializes nodes into buckets.
// Various indexes are used to accelerate important queries.
//
// As a result of using boltdb, Orchard prefers read heavy workloads.
package orchard

import (
	"fmt"
	"io"
	"sync"
	"time"

	"git.sr.ht/~whereswaldon/forest-go"
	"git.sr.ht/~whereswaldon/forest-go/fields"
	"git.sr.ht/~whereswaldon/forest-go/orchard/internal/codec"
	"git.sr.ht/~whereswaldon/forest-go/store"
	bolt "go.etcd.io/bbolt"
)

// Bucket contains the name of a bucket in bytes.
type Bucket []byte

func (b Bucket) String() string {
	return string(b)
}

var (
	// Buckets are the storage primitive used by bolt, that contain a sorted
	// list of key-value pairs. Each bucket is homogeneous.
	BucketReply     Bucket = Bucket("Reply")
	BucketIdentity  Bucket = Bucket("Identity")
	BucketCommunity Bucket = Bucket("Community")

	// Indexes are used to speed up queries.
	IndexAge      Bucket = Bucket("Age")
	IndexType     Bucket = Bucket("Type")
	IndexChildren Bucket = Bucket("Children")

	Buckets []Bucket = []Bucket{
		BucketReply,
		BucketIdentity,
		BucketCommunity,
	}

	Indexes []Bucket = []Bucket{
		IndexAge,
		IndexType,
		IndexChildren,
	}
)

// Orchard is a database-backed node store for `forest.Node`.
// Nodes are persisted as schema entities and can be queried as such.
type Orchard struct {
	*bolt.DB
	cacheLock sync.RWMutex
	readCache *store.MemoryStore

	// codec zero value defaults to Arbor serialization.
	// Can only be overriden by interal code.
	codec codec.Default
}

// Option specifies an option on the Orchard.
type Option func(*Orchard)

// Unsafe uses the unsafe codec, meaning nodes do not get validated. For testing purposes only.
func Unsafe(o *Orchard) {
	o.codec.Inner = codec.Unsafe{}
}

// Open a database file at the given path using the standard OS filesystem.
func Open(path string, opts ...Option) (*Orchard, error) {
	db, err := bolt.Open(path, 0660, nil)
	if err != nil {
		return nil, fmt.Errorf("opening database file: %w", err)
	}
	return Using(db, opts...)
}

// Using allocates an Orchard using the provided database handle.
func Using(db *bolt.DB, opts ...Option) (*Orchard, error) {
	if err := db.Update(func(tx *bolt.Tx) error {
		for _, b := range append(Buckets, Indexes...) {
			_, err := tx.CreateBucketIfNotExists(b)
			if err != nil {
				return fmt.Errorf("init %s: %w", b, err)
			}
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("init database: %w", err)
	}
	o := Orchard{
		DB:        db,
		readCache: store.NewMemoryStore(),
	}
	for _, opt := range opts {
		opt(&o)
	}
	return &o, nil
}

// Add inserts the node into the orchard. If the given node is already in the
// orchard, Add will do nothing. It is not an error to insert a node more than
// once.
func (o *Orchard) Add(node forest.Node) error {
	if _, ok, err := o.Get(node.ID()); err != nil {
		return fmt.Errorf("checking existence of node: %w", err)
	} else if ok {
		return nil
	}
	id, err := o.codec.Encode(node.ID())
	if err != nil {
		return fmt.Errorf("serializing node ID: %w", err)
	}
	v, err := o.codec.Encode(node)
	if err != nil {
		return fmt.Errorf("serializing node: %w", err)
	}
	typeIndex := func(tx *bolt.Tx, nt fields.NodeType) error {
		v, err := o.codec.Encode(nt)
		if err != nil {
			return fmt.Errorf("failed encoding node type: %w", err)
		}
		return index{Bucket: tx.Bucket(IndexType)}.Put(id, v)
	}
	ageIndex := func(tx *bolt.Tx, nt fields.NodeType, ts fields.Timestamp) error {
		k, err := o.codec.Encode(ts)
		if err != nil {
			return fmt.Errorf("failed encoding timestamp: %w", err)
		}
		b, err := tx.Bucket(IndexAge).CreateBucketIfNotExists(bucketFromNodeType(nt))
		if err != nil {
			return fmt.Errorf("failed creating bucket: %w", err)
		}
		return index{Bucket: b}.Put(k, id)
	}
	childIndex := func(tx *bolt.Tx) error {
		k, err := o.codec.Encode(node.ParentID())
		if err != nil {
			return fmt.Errorf("failed encoding parent ID: %w", err)
		}
		return index{Bucket: tx.Bucket(IndexChildren)}.Put(k, id)
	}
	return o.DB.Update(func(tx *bolt.Tx) error {
		switch n := node.(type) {
		case *forest.Reply:
			if err := tx.Bucket(BucketReply).Put(id, v); err != nil {
				return fmt.Errorf("updating bucket: %w", err)
			}
			if err := typeIndex(tx, n.Type); err != nil {
				return fmt.Errorf("updating Type index: %w", err)
			}
			if err := ageIndex(tx, n.Type, n.Created); err != nil {
				return fmt.Errorf("updating Age index: %w", err)
			}

		case *forest.Identity:
			if err := tx.Bucket(BucketIdentity).Put(id, v); err != nil {
				return fmt.Errorf("updating bucket: %w", err)
			}
			if err := typeIndex(tx, n.Type); err != nil {
				return fmt.Errorf("updating Type index: %w", err)
			}
			if err := ageIndex(tx, n.Type, n.Created); err != nil {
				return fmt.Errorf("updating Age index: %w", err)
			}
		case *forest.Community:
			if err := tx.Bucket(BucketCommunity).Put(id, v); err != nil {
				return fmt.Errorf("updating bucket: %w", err)
			}
			if err := typeIndex(tx, n.Type); err != nil {
				return fmt.Errorf("updating Type index: %w", err)
			}
			if err := ageIndex(tx, n.Type, n.Created); err != nil {
				return fmt.Errorf("updating Age index: %w", err)
			}
		}
		return childIndex(tx)
	})
}

// Get searches for a node with the given id.
// Present indicates whether the node exists, err indicates a failure to load it.
func (o *Orchard) Get(nodeID *fields.QualifiedHash) (node forest.Node, present bool, err error) {
	o.cacheLock.RLock()
	if n, ok, _ := o.readCache.Get(nodeID); ok {
		o.cacheLock.RUnlock()
		return n, ok, nil
	}
	o.cacheLock.RUnlock()
	defer func() {
		o.cacheLock.Lock()
		defer o.cacheLock.Unlock()
		if err == nil && node != nil {
			_ = o.readCache.Add(node)
		}
	}()
	var (
		nt fields.NodeType
		u  union
	)
	id, err := o.codec.Encode(nodeID)
	if err != nil {
		return nil, false, fmt.Errorf("serializing node ID: %w", err)
	}
	return node, present, o.DB.View(func(tx *bolt.Tx) error {
		v := tx.Bucket(IndexType).Get(id)
		if v == nil {
			node = nil
			present = false
			err = nil
			return nil
		}
		present = true
		if err := o.codec.Decode(v, &nt); err != nil && err != io.EOF {
			return fmt.Errorf("loading node type: %w", err)
		}
		switch nt {
		case fields.NodeTypeCommunity:
			node = &u.Community
		case fields.NodeTypeIdentity:
			node = &u.Identity
		case fields.NodeTypeReply:
			node = &u.Reply
		default:
			return fmt.Errorf("unknown node type %d", nt)
		}
		if err := o.codec.Decode(tx.Bucket(bucketFromNodeType(nt)).Get(id), node); err != nil {
			return fmt.Errorf("deserializing node: %w", err)
		}
		return nil
	})
}

func (o *Orchard) GetIdentity(id *fields.QualifiedHash) (forest.Node, bool, error) {
	return o.Get(id)
}

func (o *Orchard) GetCommunity(id *fields.QualifiedHash) (forest.Node, bool, error) {
	return o.Get(id)
}

func (o *Orchard) GetConversation(
	community *fields.QualifiedHash,
	id *fields.QualifiedHash,
) (forest.Node, bool, error) {
	return o.Get(id)
}

func (o *Orchard) GetReply(
	community *fields.QualifiedHash,
	conversation *fields.QualifiedHash,
	id *fields.QualifiedHash,
) (forest.Node, bool, error) {
	return o.Get(id)
}

// Children returns the IDs of all known child nodes of the specified ID.
func (o *Orchard) Children(
	parent *fields.QualifiedHash,
) (ch []*fields.QualifiedHash, err error) {
	k, err := o.codec.Encode(parent)
	if err != nil {
		return nil, fmt.Errorf("serializing parent ID: %w", err)
	}
	return ch, o.DB.View(func(tx *bolt.Tx) error {
		if child := tx.Bucket(IndexChildren).Get(k); child != nil {
			var id fields.QualifiedHash
			if err := o.codec.Decode(child, &id); err != nil {
				return fmt.Errorf("deserializing node ID: %w", err)
			}
			ch = append(ch, &id)
			return nil
		}
		if b := tx.Bucket(IndexChildren).Bucket(k); b != nil {
			c := b.Cursor()
			for k, v := c.First(); k != nil; k, v = c.Next() {
				var id fields.QualifiedHash
				if err := o.codec.Decode(v, &id); err != nil {
					return fmt.Errorf("deserializing node ID: %w", err)
				}
				ch = append(ch, &id)
			}
		}
		return nil
	})
}

// Recent returns a slice of nodes of a given type ordered by recency, youngest
// first.
//
// NOTE: this function may return both a valid slice of nodes and an error
// in the case that some nodes failed to be unmarshaled from disk, but others
// were successful. Calling code should always check whether the node list is
// empty before throwing it away.
func (o *Orchard) Recent(nt fields.NodeType, n int) (nodes []forest.Node, err error) {
	var (
		data   = make([][]byte, 0, n)
		b      = bucketFromNodeType(nt)
		ii     = 0
		errors Errors
	)
	if err := o.DB.View(func(tx *bolt.Tx) error {
		c := indexCursor{
			Parent:  tx.Bucket(b),
			Index:   tx.Bucket(IndexAge).Bucket(b),
			Reverse: true,
		}
		for k, v := c.First(); k != nil && ii < n; k, v = c.Next() {
			data = append(data, make([]byte, len(v)))
			copy(data[ii], v)
			ii++
		}
		return nil
	}); err != nil {
		return nodes, fmt.Errorf("copying out node data: %w", err)
	}
	for ii := range data {
		switch nt {
		case fields.NodeTypeReply:
			var n forest.Reply
			if err := o.codec.Decode(data[ii], &n); err != nil {
				errors = append(errors, fmt.Errorf("deserializing reply: %w", err))
			} else {
				nodes = append(nodes, &n)
			}
		case fields.NodeTypeIdentity:
			var n forest.Identity
			if err := o.codec.Decode(data[ii], &n); err != nil {
				errors = append(errors, fmt.Errorf("deserializing identity: %w", err))
			} else {
				nodes = append(nodes, &n)
			}
		case fields.NodeTypeCommunity:
			var n forest.Community
			if err := o.codec.Decode(data[ii], &n); err != nil {
				errors = append(errors, fmt.Errorf("deserializing community: %w", err))
			} else {
				nodes = append(nodes, &n)
			}
		}
	}
	if len(errors) > 0 {
		return nodes, errors
	}
	return nodes, nil
}

// CopyInto copies all nodes from the store into the provided store.
func (o *Orchard) CopyInto(other forest.Store) error {
	var data [][]byte
	if err := o.DB.View(func(tx *bolt.Tx) error {
		for _, b := range Buckets {
			c := tx.Bucket(b).Cursor()
			for k, v := c.First(); k != nil; k, v = c.Next() {
				by := make([]byte, len(v))
				copy(by, v)
				data = append(data, by)
			}
		}
		return nil
	}); err != nil {
		return fmt.Errorf("copying out node data: %w", err)
	}
	for _, by := range data {
		var node forest.Node
		nt, err := forest.NodeTypeOf(by)
		if err != nil {
			return err
		}
		switch nt {
		case fields.NodeTypeReply:
			var n forest.Reply
			if err := o.codec.Decode(by, &n); err != nil {
				return fmt.Errorf("deserializing reply: %w", err)
			}
			node = &n
		case fields.NodeTypeIdentity:
			var n forest.Identity
			if err := o.codec.Decode(by, &n); err != nil {
				return fmt.Errorf("deserializing identity: %w", err)
			}
			node = &n
		case fields.NodeTypeCommunity:
			var n forest.Community
			if err := o.codec.Decode(by, &n); err != nil {
				return fmt.Errorf("deserializing community: %w", err)
			}
			node = &n
		}
		if err := other.Add(node); err != nil {
			return fmt.Errorf("copying node: %w", err)
		}
	}
	return nil
}

// RemoveSubtree removes the subtree rooted at the node with the provided ID
// from the orchard.
func (o *Orchard) RemoveSubtree(id *fields.QualifiedHash) error {
	node, ok, err := o.Get(id)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}
	children, err := o.Children(id)
	if err != nil {
		return err
	}
	o.cacheLock.Lock()
	defer o.cacheLock.Unlock()
	for _, child := range children {
		if err := o.RemoveSubtree(child); err != nil {
			return err
		}
	}
	return o.delete(node)
}

func (o *Orchard) delete(n forest.Node) error {
	if err := o.readCache.RemoveSubtree(n.ID()); err != nil {
		return fmt.Errorf("failed removing node %v from cache: %w", n.ID(), err)
	}
	id, err := o.codec.Encode(n.ID())
	if err != nil {
		return fmt.Errorf("serializing node ID: %w", err)
	}
	parent, err := o.codec.Encode(n.ParentID())
	if err != nil {
		return fmt.Errorf("serializing parent ID: %w", err)
	}
	ts, err := o.codec.Encode(fields.TimestampFrom(n.CreatedAt()))
	if err != nil {
		return fmt.Errorf("failed encoding timestamp: %w", err)
	}
	b := bucketFromNode(n)
	return o.DB.Update(func(tx *bolt.Tx) error {
		if err := tx.Bucket(b).Delete(id); err != nil {
			return fmt.Errorf("failed removing from node type index: %w", err)
		}
		if err := (index{tx.Bucket(IndexAge).Bucket(b)}).DeleteByValue(ts, id); err != nil {
			return fmt.Errorf("failed removing from age index: %w", err)
		}
		if err := tx.Bucket(IndexType).Delete(id); err != nil {
			return fmt.Errorf("failed removing from type index: %w", err)
		}
		if err := (index{tx.Bucket(IndexChildren)}).Delete(id); err != nil {
			return fmt.Errorf("failed removing from children index: %w", err)
		}
		if err := (index{tx.Bucket(IndexChildren)}).DeleteByValue(parent, id); err != nil {
			return fmt.Errorf("failed removing from parent's children index: %w", err)
		}
		return nil
	})
}

// union lays out memory big enough to fit all three forest node types.
//
// Note(jfm): there may be a better way of allocating hetrogeneous data.
type union struct {
	forest.Reply
	forest.Identity
	forest.Community
}

func (u *union) Node() forest.Node {
	if len(u.Reply.Identifier) > 0 {
		return &u.Reply
	}
	if len(u.Identity.Identifier) > 0 {
		return &u.Identity
	}
	if len(u.Community.Identifier) > 0 {
		return &u.Community
	}
	return nil
}

func (u *union) UnmarshalBinary(b []byte) error {
	buf := make([]byte, len(b))
	copy(buf, b)
	n, err := forest.UnmarshalBinaryNode(buf)
	if err != nil {
		return err
	}
	switch n := n.(type) {
	case *forest.Reply:
		u.Reply = *n
	case *forest.Identity:
		u.Identity = *n
	case *forest.Community:
		u.Community = *n
	}
	return nil
}

// bucketFromNodeType returns the corresponding bucket for a node type.
func bucketFromNodeType(nt fields.NodeType) Bucket {
	switch nt {
	case fields.NodeTypeReply:
		return BucketReply
	case fields.NodeTypeIdentity:
		return BucketIdentity
	case fields.NodeTypeCommunity:
		return BucketCommunity
	}
	return nil
}

// bucketFromNode returns the corresponding bucket for a node.
func bucketFromNode(n forest.Node) Bucket {
	switch n.(type) {
	case *forest.Reply:
		return BucketReply
	case *forest.Identity:
		return BucketIdentity
	case *forest.Community:
		return BucketCommunity
	}
	return nil
}

// Errors wraps multiple errors into a single return value.
type Errors []error

func (e Errors) Error() string {
	return fmt.Sprintf("%v", []error(e))
}

// RecentReplies returns up to `q` (quantity) replies older than the timestamp.
func (o *Orchard) RecentReplies(
	ts fields.Timestamp,
	q int,
) (replies []forest.Reply, err error) {
	return replies, o.DB.View(func(tx *bolt.Tx) error {
		c := ReplyCursor{
			Inner: &indexCursor{
				Parent:  tx.Bucket(BucketReply),
				Index:   tx.Bucket(IndexAge).Bucket(BucketReply),
				Reverse: true,
			},
			Codec: &o.codec,
		}
		var ii = 0
		for r, err := c.First(); err != io.EOF; r, err = c.Next() {
			if ts > 0 && r.Created < ts && ii < q {
				replies = append(replies, r)
				ii++
			}
			if ii == q {
				break
			}
		}
		return nil
	})
}

// RepliesAfter returns up to `q` (quantity) replies newer than the timestamp.
func (o *Orchard) RepliesAfter(
	ts fields.Timestamp,
	q int,
) (replies []forest.Reply, err error) {
	return replies, o.DB.View(func(tx *bolt.Tx) error {
		c := ReplyCursor{
			Inner: &indexCursor{
				Parent:  tx.Bucket(BucketReply),
				Index:   tx.Bucket(IndexAge).Bucket(BucketReply),
				Reverse: false,
			},
			Codec: &o.codec,
		}
		var ii = 0
		for r, err := c.First(); err != io.EOF; r, err = c.Next() {
			if ts > 0 && r.Created > ts && ii < q {
				replies = append(replies, r)
				ii++
			}
			if ii == q {
				break
			}
		}
		return nil
	})
}

// RecentIdentities returns up to `q` (quantity) identities older than the timestamp.
func (o *Orchard) RecentIdentities(
	ts fields.Timestamp,
	q int,
) (identities []forest.Identity, err error) {
	return identities, o.DB.View(func(tx *bolt.Tx) error {
		c := IdentityCursor{
			Inner: &indexCursor{
				Parent:  tx.Bucket(BucketIdentity),
				Index:   tx.Bucket(IndexAge).Bucket(BucketIdentity),
				Reverse: true,
			},
			Codec: &o.codec,
		}
		var ii = 0
		for r, err := c.First(); err != io.EOF; r, err = c.Next() {
			if ts > 0 && r.Created < ts && ii < q {
				identities = append(identities, r)
				ii++
			}
			if ii == q {
				break
			}
		}
		return nil
	})
}

// RecentCommunities returns up to `q` (quantity) communities older than the timestamp.
func (o *Orchard) RecentCommunities(
	ts fields.Timestamp,
	q int,
) (communities []forest.Community, err error) {
	return communities, o.DB.View(func(tx *bolt.Tx) error {
		c := CommunityCursor{
			Inner: &indexCursor{
				Parent:  tx.Bucket(BucketCommunity),
				Index:   tx.Bucket(IndexAge).Bucket(BucketCommunity),
				Reverse: true,
			},
			Codec: &o.codec,
		}
		var ii = 0
		for r, err := c.First(); err != io.EOF; r, err = c.Next() {
			if ts > 0 && r.Created < ts && ii < q {
				communities = append(communities, r)
				ii++
			}
			if ii == q {
				break
			}
		}
		return nil
	})
}

// RecentFrom queries up to `q` (quantity) nodes of type `nt` that occur after specified timestamp.
// To page through, pass in the next oldest timestamp from the returned nodes.
//
// NOTE(jfm): There's semantic edge cases around what it means to pass in timestamp of 0.
// Theoretically, that would mean "iterate from 0 at the earliest", which should always return no
// nodes.
//
// PERF(jfm): Performance analysis pending.
func (o *Orchard) RecentFrom(
	nt fields.NodeType,
	ts time.Time,
	q int,
) (nodes []forest.Node, err error) {
	switch nt {
	case fields.NodeTypeReply:
		replies, err := o.RecentReplies(fields.TimestampFrom(ts), q)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", nt, err)
		}
		for ii := range replies {
			if ii >= q {
				break
			}
			nodes = append(nodes, &replies[ii])
		}
	case fields.NodeTypeIdentity:
		identities, err := o.RecentIdentities(fields.TimestampFrom(ts), q)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", nt, err)
		}
		for ii := range identities {
			if ii >= q {
				break
			}
			nodes = append(nodes, &identities[ii])
		}
	case fields.NodeTypeCommunity:
		communities, err := o.RecentCommunities(fields.TimestampFrom(ts), q)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", nt, err)
		}
		for ii := range communities {
			if ii >= q {
				break
			}
			nodes = append(nodes, &communities[ii])
		}
	}
	return nodes, nil
}

// ChildrenBatched traverses the children of a node in fixed-sized batches by age, youngest first.
func (o *Orchard) ChildrenBatched(
	parent *fields.QualifiedHash,
	q, offset int,
) (ch []*fields.QualifiedHash, total int, err error) {
	// 1. Seek to the offset.
	// 2. Collect up to `q` IDs and return the collected amount.
	k, err := o.codec.Encode(parent)
	if err != nil {
		return nil, 0, fmt.Errorf("serializing parent ID: %w", err)
	}
	defer func() {
		total = len(ch)
	}()
	return ch, total, o.DB.View(func(tx *bolt.Tx) error {
		if child := tx.Bucket(IndexChildren).Get(k); child != nil {
			// If offset > 0 but we only have a single child, then it must have already been
			// processed.
			if offset > 0 {
				return nil
			}
			var id fields.QualifiedHash
			if err := o.codec.Decode(child, &id); err != nil {
				return fmt.Errorf("deserializing node ID: %w", err)
			}
			ch = append(ch, &id)
			return nil
		}
		if b := tx.Bucket(IndexChildren).Bucket(k); b != nil {
			var (
				c    = b.Cursor()
				k, v = c.First()
				ii   int
			)
			for jj := 0; jj < offset; jj++ {
				k, v = c.Next()
			}
			if k == nil {
				return nil
			}
			for {
				if k == nil || ii >= q {
					break
				}
				var id fields.QualifiedHash
				if err := o.codec.Decode(v, &id); err != nil {
					return fmt.Errorf("deserializing node ID: %w", err)
				}
				ch = append(ch, &id)
				k, v = c.Next()
				ii++
			}
		}
		return nil
	})
}
