package qdht

// Node represents a single node in the qDHT.
type Node interface {
	ID() string      // Unique identifier for the node.
	Address() string // Network address or connection info for the node.
}

// DataItem represents a data item stored in the qDHT.
type DataItem interface {
	Key() string   // Unique key identifying the data item.
	Value() []byte // The data itself.
}

// QDHT defines the interface for a quantum-resistant distributed hash table.
type QDHT interface {
	// Join allows a node to join the network.
	Join(node Node) error

	// Leave gracefully removes a node from the network.
	Leave(node Node) error

	// Put stores a data item in the qDHT.
	Put(item DataItem) error

	// Get retrieves a data item from the qDHT by its key.
	Get(key string) (DataItem, error)

	// Remove deletes a data item from the qDHT by its key.
	Remove(key string) error

	// Nodes returns a list of nodes currently in the qDHT.
	Nodes() ([]Node, error)

	// Close shuts down the qDHT gracefully.
	Close() error
}

// SimpleNode Example implementations of Node and DataItem could be straightforward structs:
type SimpleNode struct {
	id      string
	address string
}

func (n SimpleNode) ID() string {
	return n.id
}

func (n SimpleNode) Address() string {
	return n.address
}

type SimpleDataItem struct {
	key   string
	value []byte
}

func (d SimpleDataItem) Key() string {
	return d.key
}

func (d SimpleDataItem) Value() []byte {
	return d.value
}
