package provider

import (
	"crypto/ed25519"

	"github.com/centrifuge/go-substrate-rpc-client/v3/types"
	"github.com/pkg/errors"
	substrate "github.com/threefoldtech/substrate-client"
)

type AccountWithTwinID struct {
	types.AccountInfo
	twinID uint32
}
type SubstrateClient interface {
	GetAccount(identity *substrate.Identity) (info types.AccountInfo, err error)
	CreateNodeContract(identity *substrate.Identity, node uint32, body []byte, hash string, publicIPs uint32) (uint64, error)
	CreateNameContract(identity *substrate.Identity, name string) (uint64, error)
	UpdateNodeContract(identity *substrate.Identity, contract uint64, body []byte, hash string) (uint64, error)
	CancelContract(identity *substrate.Identity, contract uint64) error
	GetContract(id uint64) (*substrate.Contract, error)
	GetContractIDByNameRegistration(name string) (uint64, error)
	GetNode(id uint32) (*substrate.Node, error)
	GetTwinByPubKey(pk []byte) (uint32, error)
	GetTwin(twinID uint32) (*substrate.Twin, error)
}

type SubstrateMock struct {
	accounts         map[string]AccountWithTwinID
	twins            map[uint32]substrate.Twin
	contracts        map[uint64]substrate.Contract
	contractIDByName map[string]uint64
	nodes            map[uint32]substrate.Node
	contractCounter  uint64
}

func NewSubstrateMock(identity substrate.Identity) (*SubstrateMock, error) {
	sk, err := identity.SecureKey()
	if err != nil {
		return nil, err
	}
	return &SubstrateMock{
		accounts: map[string]AccountWithTwinID{
			string(sk.Public().(ed25519.PublicKey)): {
				twinID: 1,
			},
		},
		twins: map[uint32]substrate.Twin{
			1: {},
		},
		contracts:        make(map[uint64]substrate.Contract),
		contractIDByName: make(map[string]uint64),
		nodes:            make(map[uint32]substrate.Node),
		contractCounter:  0,
	}, nil
}
func (s *SubstrateMock) GetAccount(identity *substrate.Identity) (info types.AccountInfo, err error) {
	if _, ok := s.accounts[string(identity.PublicKey)]; !ok {
		err = substrate.ErrNotFound
	}
	return
}

func (s *SubstrateMock) CreateNodeContract(identity *substrate.Identity, node uint32, body []byte, hash string, publicIPs uint32) (uint64, error) {
	sk, err := identity.SecureKey()
	if err != nil {
		return 0, err
	}
	twinID, ok := s.accounts[string(sk.Public().(ed25519.PublicKey))]
	if !ok {
		return 0, errors.New("identity not registered")
	}
	s.contractCounter++
	s.contracts[s.contractCounter] = substrate.Contract{
		Versioned:  substrate.Versioned{Version: 2},
		State:      substrate.ContractState{IsCreated: true},
		ContractID: types.U64(s.contractCounter),
		TwinID:     types.U32(twinID.twinID),
		ContractType: substrate.ContractType{IsNodeContract: true, NodeContract: substrate.NodeContract{
			Node:           types.U32(node),
			DeploymentData: body,
			DeploymentHash: hash,
			PublicIPsCount: types.U32(publicIPs),
			PublicIPs:      []substrate.PublicIP{
				// TODO: add some ips
			},
		}},
	}
	return s.contractCounter, nil
}

func (s *SubstrateMock) CreateNameContract(identity *substrate.Identity, name string) (uint64, error) {
	if _, ok := s.contractIDByName[name]; ok {
		return 0, errors.New("alerady registered")
	}
	sk, err := identity.SecureKey()
	if err != nil {
		return 0, err
	}
	twinID, ok := s.accounts[string(sk.Public().(ed25519.PublicKey))]
	if !ok {
		return 0, errors.New("identity not registered")
	}

	s.contractCounter++
	s.contracts[s.contractCounter] = substrate.Contract{
		Versioned:  substrate.Versioned{Version: 2},
		State:      substrate.ContractState{IsCreated: true},
		ContractID: types.U64(s.contractCounter),
		TwinID:     types.U32(twinID.twinID),
		ContractType: substrate.ContractType{IsNameContract: true, NameContract: substrate.NameContract{
			Name: name,
		}},
	}
	return s.contractCounter, nil
}

func (s *SubstrateMock) UpdateNodeContract(identity *substrate.Identity, contract uint64, body []byte, hash string) (uint64, error) {
	c, ok := s.contracts[contract]
	if !ok {
		return 0, substrate.ErrNotFound
	}
	if !c.ContractType.IsNodeContract {
		// TODO: error
		return 0, errors.New("not a node contract")
	}
	s.contracts[contract] = c
	return contract, nil
}

func (s *SubstrateMock) CancelContract(identity *substrate.Identity, contract uint64) error {
	c, ok := s.contracts[contract]
	if !ok {
		return substrate.ErrNotFound
	}
	c.State.IsCreated = false
	c.State.IsDeleted = true
	s.contracts[contract] = c
	return nil
}

func (s *SubstrateMock) GetContract(id uint64) (*substrate.Contract, error) {
	c, ok := s.contracts[id]
	if !ok {
		return nil, substrate.ErrNotFound
	}
	return &c, nil
}

func (s *SubstrateMock) GetContractIDByNameRegistration(name string) (uint64, error) {
	c, ok := s.contractIDByName[name]
	if !ok {
		return 0, substrate.ErrNotFound
	}
	return c, nil
}

func (s *SubstrateMock) GetNode(id uint32) (*substrate.Node, error) {
	c, ok := s.nodes[id]
	if !ok {
		return nil, substrate.ErrNotFound
	}
	return &c, nil
}
func (s *SubstrateMock) GetTwin(twinID uint32) (*substrate.Twin, error) {
	c, ok := s.twins[twinID]
	if !ok {
		return nil, substrate.ErrNotFound
	}
	return &c, nil
}

func (s *SubstrateMock) GetTwinByPubKey(pk []byte) (uint32, error) {
	c, ok := s.accounts[string(pk)]
	if !ok {
		return 0, substrate.ErrNotFound
	}
	return c.twinID, nil
}
