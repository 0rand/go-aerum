// Copyright 2017 The go-aerum Authors
// This file is part of the go-aerum library.
//
// The go-aerum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-aerum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-aerum library. If not, see <http://www.gnu.org/licenses/>.

package atmos

import (
	"github.com/AERUMTechnology/go-aerum/common"
	"github.com/AERUMTechnology/go-aerum/consensus"
	"github.com/AERUMTechnology/go-aerum/core/types"
	"github.com/AERUMTechnology/go-aerum/rpc"
)

// API is a user facing RPC API to allow controlling the signer and voting
// mechanisms of the proof-of-authority scheme.
type API struct {
	chain consensus.ChainReader
	atmos *Atmos
}

// GetSnapshot retrieves the state snapshot at a given block.
func (api *API) GetSnapshot(number *rpc.BlockNumber) (*Snapshot, error) {
	// Retrieve the requested block number (or current if none requested)
	var header *types.Header
	if number == nil || *number == rpc.LatestBlockNumber {
		header = api.chain.CurrentHeader()
	} else {
		header = api.chain.GetHeaderByNumber(uint64(number.Int64()))
	}
	// Ensure we have an actually valid block and return its snapshot
	if header == nil {
		return nil, errUnknownBlock
	}
	return api.atmos.snapshot(api.chain, header.Number.Uint64(), header.Hash(), nil)
}

// GetSnapshotAtHash retrieves the state snapshot at a given block.
func (api *API) GetSnapshotAtHash(hash common.Hash) (*Snapshot, error) {
	header := api.chain.GetHeaderByHash(hash)
	if header == nil {
		return nil, errUnknownBlock
	}
	return api.atmos.snapshot(api.chain, header.Number.Uint64(), header.Hash(), nil)
}

// GetSigners retrieves the list of authorized signers at the specified block.
func (api *API) GetSigners(number *rpc.BlockNumber) ([]common.Address, error) {
	// Retrieve the requested block number (or current if none requested)
	var header *types.Header
	if number == nil || *number == rpc.LatestBlockNumber {
		header = api.chain.CurrentHeader()
	} else {
		header = api.chain.GetHeaderByNumber(uint64(number.Int64()))
	}
	// Ensure we have an actually valid block and return the signers from its snapshot
	if header == nil {
		return nil, errUnknownBlock
	}
	snap, err := api.atmos.snapshot(api.chain, header.Number.Uint64(), header.Hash(), nil)
	if err != nil {
		return nil, err
	}
	return snap.signers(), nil
}

// GetSignersAtHash retrieves the state snapshot at a given block.
func (api *API) GetSignersAtHash(hash common.Hash) ([]common.Address, error) {
	header := api.chain.GetHeaderByHash(hash)
	if header == nil {
		return nil, errUnknownBlock
	}
	snap, err := api.atmos.snapshot(api.chain, header.Number.Uint64(), header.Hash(), nil)
	if err != nil {
		return nil, err
	}
	return snap.signers(), nil
}

/// // Proposals returns the current proposals the node tries to uphold and vote on.
/// func (api *API) Proposals() map[common.Address]bool {
/// 	api.atmos.lock.RLock()
/// 	defer api.atmos.lock.RUnlock()
///
/// 	proposals := make(map[common.Address]bool)
/// 	for address, auth := range api.atmos.proposals {
/// 		proposals[address] = auth
/// 	}
/// 	return proposals
/// }

/// // Propose injects a new authorization proposal that the signer will attempt to
/// // push through.
/// func (api *API) Propose(address common.Address, auth bool) {
/// 	api.atmos.lock.Lock()
/// 	defer api.atmos.lock.Unlock()
///
/// 	api.atmos.proposals[address] = auth
/// }

/// // Discard drops a currently running proposal, stopping the signer from casting
/// // further votes (either for or against).
/// func (api *API) Discard(address common.Address) {
/// 	api.atmos.lock.Lock()
/// 	defer api.atmos.lock.Unlock()
///
/// 	delete(api.atmos.proposals, address)
/// }
