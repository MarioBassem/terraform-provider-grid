package provider

import (
	"github.com/pkg/errors"
	client "github.com/threefoldtech/terraform-provider-grid/internal/node"
	"github.com/threefoldtech/zos/pkg/rmb"
)

type NodeClientPool struct {
	nodeClients map[uint32]*client.NodeClient

	sub SubstrateClient
	rmb rmb.Client
}

func NewNodeClient(sub SubstrateClient, rmb rmb.Client) *NodeClientPool {
	return &NodeClientPool{
		nodeClients: make(map[uint32]*client.NodeClient),
		rmb:         rmb,
		sub:         sub,
	}
}

func (k *NodeClientPool) getNodeClient(nodeID uint32) (client.NodeClientInterface, error) {
	cl, ok := k.nodeClients[nodeID]
	if ok {
		return cl, nil
	}
	nodeInfo, err := k.sub.GetNode(nodeID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get node")
	}

	cl = client.NewNodeClient(uint32(nodeInfo.TwinID), k.rmb)
	k.nodeClients[nodeID] = cl
	return cl, nil
}

func isIn(l []uint16, i uint16) bool {
	for _, x := range l {
		if i == x {
			return true
		}
	}
	return false
}

func isInByte(l []byte, i byte) bool {
	for _, x := range l {
		if i == x {
			return true
		}
	}
	return false
}

func isInUint32(l []uint32, i uint32) bool {
	for _, x := range l {
		if i == x {
			return true
		}
	}
	return false
}

func isInStr(l []string, i string) bool {
	for _, x := range l {
		if i == x {
			return true
		}
	}
	return false
}
