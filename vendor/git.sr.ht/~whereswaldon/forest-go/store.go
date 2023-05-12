package forest

import (
	"time"

	"git.sr.ht/~whereswaldon/forest-go/fields"
)

// Store describes a collection of `forest.Node`.
type Store interface {
	// Get retrieves a node by ID.
	Get(*fields.QualifiedHash) (Node, bool, error)
	// GetIdentity retrieves an identity node by ID.
	GetIdentity(*fields.QualifiedHash) (Node, bool, error)
	// GetCommunity retrieves a community node by ID.
	GetCommunity(*fields.QualifiedHash) (Node, bool, error)
	// GetConversation retrieves a conversation node by ID.
	GetConversation(communityID, conversationID *fields.QualifiedHash) (Node, bool, error)
	// GetReply retrieves a reply node by ID.
	GetReply(communityID, conversationID, replyID *fields.QualifiedHash) (Node, bool, error)
	// Children returns a list of child nodes for the given node ID.
	Children(*fields.QualifiedHash) ([]*fields.QualifiedHash, error)
	// Recent returns recently-created (as per the timestamp in the node) nodes.
	// It may return both a slice of nodes and an error if some nodes in the
	// store were unreadable.
	Recent(nodeType fields.NodeType, quantity int) ([]Node, error)
	// Add inserts a node into the store. It is *not* an error to insert a node which is already
	// stored. Implementations must not return an error in this case.
	Add(Node) error
	// RemoveSubtree from the store.
	RemoveSubtree(*fields.QualifiedHash) error
}

// Copiable stores can copy themselves into another store.
type Copiable interface {
	CopyInto(Store) error
}

// Paginated stores can page through nodes with a series of queries.
type Paginated interface {
	// ChildrenBatched allows paging through the children of the node in fixed-size batches.
	// Iteration order is youngest first.
	ChildrenBatched(
		root *fields.QualifiedHash,
		quantity, offset int,
	) (batch []*fields.QualifiedHash, total int, err error)
	// RecentFrom allows paging through nodes in fixed-size batches, starting from the given
	// timestamp.
	// Iteration order is youngest first.
	RecentFrom(nt fields.NodeType, ts time.Time, quantity int) ([]Node, error)
}
