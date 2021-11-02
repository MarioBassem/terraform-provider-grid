package client

import (
	"context"
	"log"
	"net"

	"github.com/pkg/errors"
	"github.com/threefoldtech/zos/pkg/gridtypes"
)

type NodeClientInterface interface {
	DeploymentDeploy(ctx context.Context, dl gridtypes.Deployment) error
	DeploymentUpdate(ctx context.Context, dl gridtypes.Deployment) error
	DeploymentGet(ctx context.Context, contractID uint64) (dl gridtypes.Deployment, err error)
	DeploymentDelete(ctx context.Context, contractID uint64) error
	NetworkListWGPorts(ctx context.Context) ([]uint16, error)
	NetworkListInterfaces(ctx context.Context) (map[string][]net.IP, error)
	NetworkGetPublicConfig(ctx context.Context) (cfg PublicConfig, err error)
}
type NodeClientMock struct {
	deployments  map[uint64]gridtypes.Deployment
	publicConfig PublicConfig
	ifs          map[string][]net.IP
	wgPorts      []uint16
}

func NewNodeClientMock(publicConfig PublicConfig, ifs map[string][]net.IP, wgPorts []uint16) NodeClientMock {
	return NodeClientMock{
		deployments:  make(map[uint64]gridtypes.Deployment),
		publicConfig: publicConfig,
		ifs:          ifs,
		wgPorts:      wgPorts,
	}
}

func (nc *NodeClientMock) DeploymentDeploy(ctx context.Context, dl gridtypes.Deployment) error {
	log.Printf("contract id: %d\n", dl.ContractID)
	for i := range dl.Workloads {
		dl.Workloads[i].Result.State = gridtypes.StateOk
	}
	nc.deployments[dl.ContractID] = dl
	return nil
}

func (nc *NodeClientMock) DeploymentUpdate(ctx context.Context, dl gridtypes.Deployment) error {
	nc.deployments[dl.ContractID] = dl
	for i := range dl.Workloads {
		dl.Workloads[i].Result.State = gridtypes.StateOk
	}
	return nil
}

func (nc *NodeClientMock) DeploymentGet(ctx context.Context, contractID uint64) (gridtypes.Deployment, error) {
	dl, ok := nc.deployments[contractID]
	if !ok {
		return gridtypes.Deployment{}, errors.New("deployment not found")
	}
	return dl, nil
}

func (nc *NodeClientMock) DeploymentDelete(ctx context.Context, contractID uint64) error {
	dl, ok := nc.deployments[contractID]
	if !ok {
		return errors.New("deployment not found")
	}
	for i := range dl.Workloads {
		dl.Workloads[i].Result.State = gridtypes.StateDeleted
	}
	nc.deployments[dl.ContractID] = dl
	return nil
}

func (nc *NodeClientMock) NetworkListWGPorts(ctx context.Context) ([]uint16, error) {
	return nc.wgPorts, nil
}

func (nc *NodeClientMock) NetworkListInterfaces(ctx context.Context) (map[string][]net.IP, error) {
	return nc.ifs, nil
}

func (nc *NodeClientMock) NetworkGetPublicConfig(ctx context.Context) (cfg PublicConfig, err error) {
	return nc.publicConfig, nil
}
