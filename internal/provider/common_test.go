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
 *  create deployments: node1: dl1. node2: dl1, node3: dl1
 *  expected: all created
 *  update deployments: node1: dl1 (changed), node2: dl1 (same), delete node3
 *  expected: node1 dl1 updated, node1 dl2 deleted, node2 d1 not updated, node3 deleted
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
		{
			Name: "network",
			Type: zos.NetworkType,
			Data: gridtypes.MustMarshal(zos.Network{
				NetworkIPRange: gridtypes.MustParseIPNet("10.0.0.0/16"),
				Subnet:         gridtypes.MustParseIPNet("10.0.1.0/24"),
				WGPrivateKey:   "1",
			}),
		},
		{
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
		{
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
		{
			Name: "network",
			Type: zos.NetworkType,
			Data: gridtypes.MustMarshal(zos.Network{
				NetworkIPRange: gridtypes.MustParseIPNet("10.0.0.0/16"),
				Subnet:         gridtypes.MustParseIPNet("10.0.1.0/24"),
				WGPrivateKey:   "1",
			}),
		},
		{
			Name: "gateway",
			Type: zos.GatewayFQDNProxyType,
			Data: gridtypes.MustMarshal(zos.GatewayFQDNProxy{
				FQDN:     "updated.fqdn.com",
				Backends: []zos.Backend{zos.Backend("http://1.1.1.1:123")},
			}),
		},
		{
			Name: "disk",
			Type: zos.ZMountType,
			Data: gridtypes.MustMarshal(zos.ZMount{
				Size: gridtypes.Gigabyte,
			}),
		},
	},
}

var node3_dl1 = node2_dl1_first
var node4_dl1 = node2_dl1_first

func checkEqualDeployments(dl1 gridtypes.Deployment, dl2 gridtypes.Deployment) bool {
	c1, e1 := dl1.ChallengeHash()
	c2, e2 := dl2.ChallengeHash()
	return e1 == nil && e2 == nil && bytes.Equal(c1, c2)
}

func checkEqualDeploymentsWithoutVersions(dl1 gridtypes.Deployment, dl2 gridtypes.Deployment) bool {
	dl1 = client.CloneDeployment(&dl1)
	dl1.Version = 0
	for i := range dl1.Workloads {
		dl1.Workloads[i].Version = 0
	}
	dl2 = client.CloneDeployment(&dl1)
	dl2.Version = 0
	for i := range dl2.Workloads {
		dl2.Workloads[i].Version = 0
	}
	c1, _ := dl1.ChallengeHash()
	c2, _ := dl2.ChallengeHash()
	printDeployments(map[uint32]gridtypes.Deployment{
		1: dl1,
		2: dl2,
	})
	b1 := strings.Builder{}
	b2 := strings.Builder{}
	dl1.Challenge(&b1)
	dl2.Challenge(&b2)
	log.Println(b1.String())
	log.Println(b2.String())
	return bytes.Equal(c1, c2)
}

func TestDeployConsistentDeployments(t *testing.T) {
	nc2 := client.NewNodeClientMock(client.PublicConfig{}, nil, nil)
	nc3 := client.NewNodeClientMock(client.PublicConfig{}, nil, nil)
	nc4 := client.NewNodeClientMock(client.PublicConfig{}, nil, nil)
	nodes := NewNodeClientCollectionMock(map[uint32]client.NodeClientMock{
		2: nc2,
		3: nc3,
		4: nc4,
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
		4: node4_dl1,
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
	node4Result, err := nc4.DeploymentGet(context.Background(), cur[4])
	if err != nil {
		t.Fatalf("getting the deployment failed %s", err.Error())
	}
	if !checkEqualDeployments(node4Result, node3_dl1) {
		t.Fatalf("passed and retrieved mismatch")
	}
	oldNode4Deployment := cur[4]
	// update

	cur, err = deployConsistentDeployments(context.Background(), cur, map[uint32]gridtypes.Deployment{
		2: node2_dl1_second,
		3: node3_dl1,
	}, &nodes, &apiClient)
	if err != nil {
		t.Fatalf("failed to deploy deployments %s", err.Error())
	}
	node2Result, err = nc2.DeploymentGet(context.Background(), cur[2])
	if err != nil {
		t.Fatalf("getting the deployment failed %s", err.Error())
	}
	node2_dl1_second.Version = 1
	node2_dl1_second.Workloads[1].Version = 1
	node2_dl1_second.Workloads[2].Version = 1
	if !checkEqualDeployments(node2Result, node2_dl1_second) {
		t.Fatalf("passed and retrieved mismatch")
	}
	node3Result, err = nc3.DeploymentGet(context.Background(), cur[3])
	if err != nil {
		t.Fatalf("getting the deployment failed %s", err.Error())
	}
	if !checkEqualDeployments(node3Result, node3_dl1) {
		t.Fatalf("passed and retrieved mismatch")
	}
	node4Result, err = nc4.DeploymentGet(context.Background(), oldNode4Deployment)
	if err != nil {
		t.Fatalf("getting the deployment failed %s", err.Error())
	}
	if node4Result.Workloads[0].Result.State != gridtypes.StateDeleted {
		t.Fatalf("node 4 contract not deleted")
	}
}

func TestRevertingDeployments(t *testing.T) {
	nc2 := client.NewNodeClientMock(client.PublicConfig{}, nil, nil)
	nc3 := client.NewNodeClientMock(client.PublicConfig{}, nil, nil)
	nc4 := client.NewNodeClientMock(client.PublicConfig{}, nil, nil)
	nodes := NewNodeClientCollectionMock(map[uint32]client.NodeClientMock{
		2: nc2,
		3: nc3,
		4: nc4,
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
		4: node4_dl1,
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
	node4Result, err := nc4.DeploymentGet(context.Background(), cur[4])
	if err != nil {
		t.Fatalf("getting the deployment failed %s", err.Error())
	}
	if !checkEqualDeployments(node4Result, node3_dl1) {
		t.Fatalf("passed and retrieved mismatch")
	}
	oldNode3Deployment := client.CloneDeployment(&node3_dl1)
	node3_dl1.Workloads[0].Name = node3_dl1.Workloads[1].Name
	d, _ := node3_dl1.Workloads[0].WorkloadData()
	d.(*zos.Network).WGPrivateKey = "2" // change it to update
	node3_dl1.Workloads[0].Data = gridtypes.MustMarshal(d)
	cur, err = deployDeployments(context.Background(), cur, map[uint32]gridtypes.Deployment{
		2: node2_dl1_second,
		3: node3_dl1,
	}, &nodes, &apiClient, true)

	if err == nil {
		t.Fatalf("this deployment should fail")
	}
	node2Result, err = nc2.DeploymentGet(context.Background(), cur[2])
	if err != nil {
		t.Fatalf("getting the deployment failed %s", err.Error())
	}
	if !checkEqualDeploymentsWithoutVersions(node2Result, node2_dl1_second) {
		t.Fatalf("passed and retrieved mismatch")
	}
	node3Result, err = nc3.DeploymentGet(context.Background(), cur[3])
	if err != nil {
		t.Fatalf("getting the deployment failed %s", err.Error())
	}
	if !checkEqualDeploymentsWithoutVersions(node3Result, oldNode3Deployment) {
		t.Fatalf("passed and retrieved mismatch")
	}
	node4Result, err = nc4.DeploymentGet(context.Background(), cur[4])
	if err != nil {
		t.Fatalf("getting the deployment failed %s", err.Error())
	}
	if !checkEqualDeploymentsWithoutVersions(node4Result, node3_dl1) {
		t.Fatalf("passed and retrieved mismatch")
	}
}

func TestSameWorkloadName(t *testing.T) {

}
