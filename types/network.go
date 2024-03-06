package types

import (
	"context"
	"io"
	"net"
	"time"
)

type Network struct {
	Conn    net.Conn
	Addr    net.TCPAddr
	Self    *Node
	Version string
}

type NBucket map[int32][]*Node

type Node struct {
	ID         NodeID
	PeerInfo   *Peer
	Host, Port string
}

type NodeID [20]byte

type Peer struct {
	Keys [][]byte
}

type Protocol struct {
	Name   string `json:"protocol_name"`
	ID     []byte `json:"protocol_id"`
	Opcode byte   `json:",omitempty"`
}

type ProtocolHandler func(path string, cb func(ctx context.Context, rw io.ReadWriteCloser) error)

type TrustMesh interface {
	Network() Network
	Nodes() map[string]*Node
	BootNodes() []*Node
	UTCTime() time.Time
	UpTime() time.Duration
}
