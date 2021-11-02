package provider

import (
	"bytes"
	"context"
	"log"
	"net"
	"strings"
	"testing"

	"github.com/threefoldtech/substrate-client"
	client "github.com/threefoldtech/terraform-provider-grid/internal/node"
	"github.com/threefoldtech/zos/pkg/gridtypes"
	"github.com/threefoldtech/zos/pkg/gridtypes/zos"
)

/*
 *  create deployments: node1: dl1, dl2. node2: dl1
 *  expected: all created
 *  update deployments: node1: dl1 (changed), node2: dl1 (same)
 *  expected: node1 dl1 updated, node1 dl2 deleted, node2 d1 not updated
 */
var node2_dl1_first = gridtypes.Deployment{
	TwinID: 1,
	SignatureRequirement: gridtypes.SignatureRequirement{
		WeightRequired: 1,
		Requests: []gridtypes.SignatureRequest{
			{
				TwinID: 1,
				Weight: 1,
			},
		},
	},
	Workloads: []gridtypes.Workload{
		gridtypes.Workload{
			Name: "network",
			Type: zos.NetworkType,
			Data: gridtypes.MustMarshal(zos.Network{
				NetworkIPRange: gridtypes.MustParseIPNet("10.0.0.0/16"),
				Subnet:         gridtypes.MustParseIPNet("10.0.1.0/24"),
				WGPrivateKey:   "1",
			}),
		},
		gridtypes.Workload{
			Name: "vm",
			Type: zos.ZMachineType,
			Data: gridtypes.MustMarshal(zos.ZMachine{
				ComputeCapacity: zos.MachineCapacity{CPU: 1, Memory: gridtypes.Gigabyte},
				Network: zos.MachineNetwork{
					Interfaces: []zos.MachineInterface{
						{
							Network: gridtypes.Name("network"),
							IP:      net.ParseIP("10.0.1.2"),
						},
					},
				},
			}),
		},
		gridtypes.Workload{
			Name: "gateway",
			Type: zos.GatewayFQDNProxyType,
			Data: gridtypes.MustMarshal(zos.GatewayFQDNProxy{
				FQDN:     "my.fqdn.com",
				Backends: []zos.Backend{zos.Backend("http://1.1.1.1:123")},
			}),
		},
	},
}

var node2_dl1_second = gridtypes.Deployment{
	TwinID: 2,
	SignatureRequirement: gridtypes.SignatureRequirement{
		WeightRequired: 1,
		Requests: []gridtypes.SignatureRequest{
			{
				TwinID: 1,
				Weight: 1,
			},
		},
	},
	Workloads: []gridtypes.Workload{
		gridtypes.Workload{
			Name: "network",
			Type: zos.NetworkType,
			Data: gridtypes.MustMarshal(zos.Network{
				NetworkIPRange: gridtypes.MustParseIPNet("10.0.0.0/16"),
				Subnet:         gridtypes.MustParseIPNet("10.0.1.0/24"),
				WGPrivateKey:   "1",
			}),
		},
		gridtypes.Workload{
			Name: "gateway",
			Type: zos.GatewayFQDNProxyType,
			Data: gridtypes.MustMarshal(zos.GatewayFQDNProxy{
				FQDN: "updated.fqdn.com",
			}),
		},
		gridtypes.Workload{
			Name: "qsfs",
			Type: zos.QuantumSafeFSType,
			Data: gridtypes.MustMarshal(zos.QuantumSafeFS{}),
		},
	},
}

var node3_dl1 = node2_dl1_first

func checkEqualDeployments(dl1 gridtypes.Deployment, dl2 gridtypes.Deployment) bool {
	// printDeployments(map[uint32]gridtypes.Deployment{
	// 	1: dl1,
	// 	2: dl2,
	// })
	b1 := strings.Builder{}
	dl1.Challenge(&b1)
	log.Println(b1.String())
	b2 := strings.Builder{}
	dl2.Challenge(&b2)
	log.Println(b2.String())
	c1, _ := dl1.ChallengeHash()

	c2, _ := dl2.ChallengeHash()
	log.Println(c1)
	log.Println(c2)
	return bytes.Equal(c1, c2)
}

func TestDeployConsistentDeployments(t *testing.T) {
	nc2 := client.NewNodeClientMock(client.PublicConfig{}, nil, nil)
	nc3 := client.NewNodeClientMock(client.PublicConfig{}, nil, nil)
	nodes := NewNodeClientCollectionMock(map[uint32]client.NodeClientMock{
		2: nc2,
		3: nc3,
	})
	identity, err := substrate.IdentityFromPhrase("include earth wine leave core become that kiss alarm try student seminar")
	if err != nil {
		t.Fatal(err)
	}
	userSK, err := identity.SecureKey()
	if err != nil {
		t.Fatal(err)
	}
	sub, err := NewSubstrateMock(identity)
	if err != nil {
		t.Fatal(err)
	}
	apiClient := apiClient{
		identity: &identity,
		userSK:   userSK,
		sub:      sub,
	}
	cur, err := deployConsistentDeployments(context.Background(), map[uint32]uint64{}, map[uint32]gridtypes.Deployment{
		2: node2_dl1_first,
		3: node3_dl1,
	}, &nodes, &apiClient)
	if err != nil {
		t.Fatalf("failed to deploy deployments %s", err.Error())
	}
	node2Result, err := nc2.DeploymentGet(context.Background(), cur[2])
	if err != nil {
		t.Fatalf("getting the deployment failed %s", err.Error())
	}
	if !checkEqualDeployments(node2Result, node2_dl1_first) {
		t.Fatalf("passed and retrieved mismatch")
	}
	node3Result, err := nc3.DeploymentGet(context.Background(), cur[3])
	if err != nil {
		t.Fatalf("getting the deployment failed %s", err.Error())
	}
	if !checkEqualDeployments(node3Result, node3_dl1) {
		t.Fatalf("passed and retrieved mismatch")
	}

}
