package provider

import (
	"context"
	"net"
	"testing"

	"github.com/pkg/errors"
	client "github.com/threefoldtech/terraform-provider-grid/internal/node"
	"github.com/threefoldtech/zos/pkg/gridtypes"
)

type NodeClientCollectionMock struct {
	cls map[uint32]client.NodeClientMock
}

func NewNodeClientCollectionMock(cls map[uint32]client.NodeClientMock) NodeClientCollectionMock {
	return NodeClientCollectionMock{
		cls,
	}
}

func (nc *NodeClientCollectionMock) getNodeClient(nodeID uint32) (client.NodeClientInterface, error) {
	cl, ok := nc.cls[nodeID]
	if !ok {
		return nil, errors.New("node client is not added to the mock")
	}
	return &cl, nil
}

func TestCorrectPublicConfig(t *testing.T) {
	nc := client.NewNodeClientMock(client.PublicConfig{
		IPv4: gridtypes.MustParseIPNet("50.50.50.50/24"),
	}, nil, nil)
	nodes := NewNodeClientCollectionMock(map[uint32]client.NodeClientMock{
		1: nc,
	})
	if err := validatePublicNode(context.Background(), 1, &nodes); err != nil {
		t.Error("a correct public configuration didn't pass the validation")
	}
}

func TestLocalPublicConfig(t *testing.T) {
	nc := client.NewNodeClientMock(client.PublicConfig{
		IPv4: gridtypes.MustParseIPNet("192.168.123.50/24"),
	}, nil, nil)
	nodes := NewNodeClientCollectionMock(map[uint32]client.NodeClientMock{
		1: nc,
	})
	if err := validatePublicNode(context.Background(), 1, &nodes); err == nil {
		t.Error("a public configuration with local ipv4 passed the validation")
	}
}

func TestMissingIPv4PublicConfig(t *testing.T) {
	nc := client.NewNodeClientMock(client.PublicConfig{}, nil, nil)
	nodes := NewNodeClientCollectionMock(map[uint32]client.NodeClientMock{
		1: nc,
	})
	if err := validatePublicNode(context.Background(), 1, &nodes); err == nil {
		t.Error("a public configuration with missing ipv4 passed the validation")
	}
}

func TestCorrectZOSIPv6GetNodeEndpoint(t *testing.T) {
	cIP := net.ParseIP("123:342::")
	nc := client.NewNodeClientMock(client.PublicConfig{}, map[string][]net.IP{
		"zos": {net.ParseIP("127.0.0.1"), net.ParseIP("fe80::1"), cIP},
	}, nil)
	endpoint, err := getNodeEndpoint(context.TODO(), &nc)
	if err != nil {
		t.Error("zos interface contains a public ipv6 that wasn't used")
	}
	if cIP.String() != endpoint.String() {
		t.Errorf("incorrect ip returned %s", endpoint.String())
	}
}

func TestCorrectZOSIPv4GetNodeEndpoint(t *testing.T) {
	cIP := net.ParseIP("50.50.50.50")
	nc := client.NewNodeClientMock(client.PublicConfig{}, map[string][]net.IP{
		"zos": {net.ParseIP("127.0.0.1"), net.ParseIP("fe80::1"), cIP},
	}, nil)
	endpoint, err := getNodeEndpoint(context.TODO(), &nc)
	if err != nil {
		t.Error("zos interface contains a public ipv6 that wasn't used")
	}
	if cIP.String() != endpoint.String() {
		t.Errorf("incorrect ip returned %s", endpoint.String())
	}
}

func TestCorrectPublicConfigIPv4GetNodeEndpoint(t *testing.T) {
	cIP := gridtypes.MustParseIPNet("50.50.50.50/24")
	nc := client.NewNodeClientMock(client.PublicConfig{
		IPv4: cIP,
	}, map[string][]net.IP{
		"zos": {net.ParseIP("127.0.0.1"), net.ParseIP("fe80::1"), net.ParseIP("123.123.132.132")},
	}, nil)
	endpoint, err := getNodeEndpoint(context.TODO(), &nc)
	if err != nil {
		t.Error("zos interface contains a public ipv6 that wasn't used")
	}
	if cIP.IP.String() != endpoint.String() {
		t.Errorf("incorrect ip returned %s", endpoint.String())
	}
}

func TestCorrectPublicConfigIPv6GetNodeEndpoint(t *testing.T) {
	cIP := gridtypes.MustParseIPNet("53::1/64")
	nc := client.NewNodeClientMock(client.PublicConfig{
		IPv6: cIP,
		IPv4: gridtypes.MustParseIPNet("127.0.0.1/8"),
	}, map[string][]net.IP{
		"zos": {net.ParseIP("127.0.0.1"), net.ParseIP("fe80::1"), net.ParseIP("127.0.0.1")},
	}, nil)
	endpoint, err := getNodeEndpoint(context.TODO(), &nc)
	if err != nil {
		t.Errorf("zos interface contains a public ipv6 that wasn't used: %s", err.Error())
	}
	if cIP.IP.String() != endpoint.String() {
		t.Errorf("incorrect ip returned %s", endpoint.String())
	}
}
func TestNodeFreeWGPort(t *testing.T) {
	l := make([]uint16, 0)
	i := 65535 - 1
	chosen := 32561
	for i >= 1 {
		if i != chosen {
			l = append(l, uint16(i))
		}
		i -= 1
	}
	nc := client.NewNodeClientMock(client.PublicConfig{}, nil, l)
	p, err := getNodeFreeWGPort(context.Background(), &nc)
	if err != nil {
		t.Errorf("couldn't select a free wg port %s", err.Error())
	}
	if p != chosen {
		t.Errorf("chose the wrong port: %d, expected: %d", p, chosen)
	}
}
