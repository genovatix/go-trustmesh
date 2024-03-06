package network

import (
	"time"
	"trustmesh/types"
)

type P2PNetwork interface {
}

type p2pNetwork struct{}

func (p p2pNetwork) Network() types.Network {
	//TODO implement me
	panic("implement me")
}

func (p p2pNetwork) Nodes() map[string]*types.Node {
	//TODO implement me
	panic("implement me")
}

func (p p2pNetwork) BootNodes() []*types.Node {
	//TODO implement me
	panic("implement me")
}

func (p p2pNetwork) UTCTime() time.Time {
	//TODO implement me
	panic("implement me")
}

func (p p2pNetwork) UpTime() time.Duration {
	//TODO implement me
	panic("implement me")
}

func NewNetwork() types.TrustMesh {
	return &p2pNetwork{}
}
