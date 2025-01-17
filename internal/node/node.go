// Package client provides a simple RMB interface to work with the node.
//
// # Requirements
//
// 1. A msgbusd instance must be running on the node. this client uses RMB (message bus)
// to send messages to nodes, and get the repspons.
// 2. A valid ed25519 key pair. this key is used to sign deployments and MUST be the same
// key used to configure the local twin on substrate.
//
// # Simple deployment
//
// create an instance from the default rmb client.
// ```
// cl, err := rmb.Default()
//
//	if err != nil {
//		panic(err)
//	}
//
// ```
// then create an instance of the node client
// ```
// node := client.NewNodeClient(NodeTwinID, cl)
// ```
// define your deployment object
// ```
//
//	dl := gridtypes.Deployment{
//		Version: Version,
//		TwinID:  Twin, //LocalTwin,
//		// this contract id must match the one on substrate
//		Workloads: []gridtypes.Workload{
//			network(), // network workload definition
//			zmount(), // zmount workload definition
//			publicip(), // public ip definition
//			zmachine(), // zmachine definition
//		},
//		SignatureRequirement: gridtypes.SignatureRequirement{
//			WeightRequired: 1,
//			Requests: []gridtypes.SignatureRequest{
//				{
//					TwinID: Twin,
//					Weight: 1,
//				},
//			},
//		},
//	}
//
// ```
// compute hash
// ```
// hash, err := dl.ChallengeHash()
//
//	if err != nil {
//		panic("failed to create hash")
//	}
//
// fmt.Printf("Hash: %x\n", hash)
// ```
// create the contract and ge the contract id
// then
// “
// dl.ContractID = 11 // from substrate
// ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
// defer cancel()
// err = node.DeploymentDeploy(ctx, dl)
//
//	if err != nil {
//		panic(err)
//	}
//
// ```
package client

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/threefoldtech/terraform-provider-grid/pkg/subi"
	"github.com/threefoldtech/zos/pkg/capacity/dmi"
	"github.com/threefoldtech/zos/pkg/gridtypes"
	"github.com/threefoldtech/zos/pkg/rmb"
)

// IfaceType define the different public interface supported
type IfaceType string

// PublicConfig is the configuration of the interface
// that is connected to the public internet
type PublicConfig struct {
	// Type define if we need to use
	// the Vlan field or the MacVlan
	Type IfaceType `json:"type"`
	// Vlan int16     `json:"vlan"`
	// Macvlan net.HardwareAddr

	IPv4 gridtypes.IPNet `json:"ipv4"`
	IPv6 gridtypes.IPNet `json:"ipv6"`

	GW4 net.IP `json:"gw4"`
	GW6 net.IP `json:"gw6"`

	// Domain is the node domain name like gent01.devnet.grid.tf
	// or similar
	Domain string `json:"domain"`
}

// NodeClient struct
type NodeClient struct {
	nodeTwin uint32
	bus      rmb.Client
}

type args map[string]interface{}

// NewNodeClient creates a new node RMB client. This client then can be used to
// communicate with the node over RMB.
func NewNodeClient(nodeTwin uint32, bus rmb.Client) *NodeClient {
	return &NodeClient{nodeTwin, bus}
}

// DeploymentDeploy sends the deployment to the node for processing.
func (n *NodeClient) DeploymentDeploy(ctx context.Context, dl gridtypes.Deployment) error {
	const cmd = "zos.deployment.deploy"
	return n.bus.Call(ctx, n.nodeTwin, cmd, dl, nil)
}

// DeploymentUpdate update the given deployment. deployment must be a valid update for
// a deployment that has been already created via DeploymentDeploy
func (n *NodeClient) DeploymentUpdate(ctx context.Context, dl gridtypes.Deployment) error {
	const cmd = "zos.deployment.update"
	return n.bus.Call(ctx, n.nodeTwin, cmd, dl, nil)
}

// DeploymentGet gets a deployment via contract ID
func (n *NodeClient) DeploymentGet(ctx context.Context, contractID uint64) (dl gridtypes.Deployment, err error) {
	const cmd = "zos.deployment.get"
	in := args{
		"contract_id": contractID,
	}

	if err = n.bus.Call(ctx, n.nodeTwin, cmd, in, &dl); err != nil {
		return dl, err
	}

	return dl, nil
}

// DeploymentDelete deletes a deployment, the node will make sure to decomission all deployments
// and set all workloads to deleted. A call to Get after delete is valid
func (n *NodeClient) DeploymentDelete(ctx context.Context, contractID uint64) error {
	const cmd = "zos.deployment.delete"
	in := args{
		"contract_id": contractID,
	}

	return n.bus.Call(ctx, n.nodeTwin, cmd, in, nil)
}

// Counters returns some node statistics. Including total and available cpu, memory, storage, etc...
func (n *NodeClient) Counters(ctx context.Context) (total gridtypes.Capacity, used gridtypes.Capacity, err error) {
	const cmd = "zos.statistics.get"
	var result struct {
		Total gridtypes.Capacity `json:"total"`
		Used  gridtypes.Capacity `json:"used"`
	}
	if err = n.bus.Call(ctx, n.nodeTwin, cmd, nil, &result); err != nil {
		return
	}

	return result.Total, result.Used, nil
}

// NetworkListWGPorts return a list of all "taken" ports on the node. A new deployment
// should be careful to use a free port for its network setup.
func (n *NodeClient) NetworkListWGPorts(ctx context.Context) ([]uint16, error) {
	const cmd = "zos.network.list_wg_ports"
	var result []uint16

	if err := n.bus.Call(ctx, n.nodeTwin, cmd, nil, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// NetworkListInterfaces return a map of all interfaces and their ips
func (n *NodeClient) NetworkListInterfaces(ctx context.Context) (map[string][]net.IP, error) {
	const cmd = "zos.network.interfaces"
	var result map[string][]net.IP

	if err := n.bus.Call(ctx, n.nodeTwin, cmd, nil, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// DeploymentChanges return changes of a deployment via contract ID
func (n *NodeClient) DeploymentChanges(ctx context.Context, contractID uint64) (changes []gridtypes.Workload, err error) {
	const cmd = "zos.deployment.changes"
	in := args{
		"contract_id": contractID,
	}

	if err = n.bus.Call(ctx, n.nodeTwin, cmd, in, &changes); err != nil {
		return changes, err
	}

	return changes, nil
}

// RandomFreePort query the node for used ports, then it tries to find a ramdom
// port that is in not in the "taken" ports list, this can be used to set up
// network wireguard ports
// func (n *NodeClient) RandomFreePort(ctx context.Context) (uint16, error) {
// 	used, err := n.NetworkListWGPorts(ctx)
// 	if err != nil {
// 		return 0, err
// 	}
// 	//rand.
// }

// NetworkListIPs list taken public IPs on the node
func (n *NodeClient) NetworkListIPs(ctx context.Context) ([]string, error) {
	const cmd = "zos.network.list_public_ips"
	var result []string

	if err := n.bus.Call(ctx, n.nodeTwin, cmd, nil, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// NetworkGetPublicConfig returns the current public node network configuration. A node with a
// public config can be used as an access node for wireguard.
func (n *NodeClient) NetworkGetPublicConfig(ctx context.Context) (cfg PublicConfig, err error) {
	const cmd = "zos.network.public_config_get"

	if err = n.bus.Call(ctx, n.nodeTwin, cmd, nil, &cfg); err != nil {
		return
	}

	return
}

// NetworkGetPublicConfig returns the current public node network configuration. A node with a
// public config can be used as an access node for wireguard.
func (n *NodeClient) NetworkSetPublicConfig(ctx context.Context, cfg PublicConfig) error {
	const cmd = "zos.network.public_config_set"

	if err := n.bus.Call(ctx, n.nodeTwin, cmd, cfg, nil); err != nil {
		return err
	}

	return nil
}

// SystemDMI executes dmidecode to get dmidecode output
func (n *NodeClient) SystemDMI(ctx context.Context) (result dmi.DMI, err error) {
	const cmd = "zos.system.dmi"

	if err = n.bus.Call(ctx, n.nodeTwin, cmd, nil, &result); err != nil {
		return
	}

	return
}

// SystemHypervisor executes hypervisor cmd
func (n *NodeClient) SystemHypervisor(ctx context.Context) (result string, err error) {
	const cmd = "zos.system.hypervisor"

	if err = n.bus.Call(ctx, n.nodeTwin, cmd, nil, &result); err != nil {
		return
	}

	return
}

// Version is ZOS version
type Version struct {
	ZOS   string `json:"zos"`
	ZInit string `json:"zinit"`
}

// SystemVersion executes system version cmd
func (n *NodeClient) SystemVersion(ctx context.Context) (ver Version, err error) {
	const cmd = "zos.system.version"

	if err = n.bus.Call(ctx, n.nodeTwin, cmd, nil, &ver); err != nil {
		return
	}

	return
}

// IsNodeUp checks if the node is up
func (n *NodeClient) IsNodeUp(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := n.SystemVersion(ctx)
	if err != nil {
		return err
	}

	return nil
}

// AreNodesUp checks if nodes are up
func AreNodesUp(ctx context.Context, sub subi.SubstrateExt, nodes []uint32, nc NodeClientGetter) (err error) {
	var wg sync.WaitGroup

	for _, node := range nodes {

		wg.Add(1)
		go func(node uint32) {

			defer wg.Done()
			cl, clientErr := nc.GetNodeClient(sub, node)
			if clientErr != nil {
				err = multierror.Append(err, fmt.Errorf("couldn't get node %d client: %w", node, clientErr))
				return
			}
			if clientErr := cl.IsNodeUp(ctx); clientErr != nil {
				err = multierror.Append(err, fmt.Errorf("couldn't reach node %d: %w", node, clientErr))
			}

		}(node)
	}

	wg.Wait()
	return
}
