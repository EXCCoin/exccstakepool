// Copyright (c) 2013-2014 The btcsuite developers
// Copyright (c) 2015-2016 The Decred developers
// Copyright (c) 2018 The EXCCoin team
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	"github.com/EXCCoin/exccd/chaincfg"
	"github.com/EXCCoin/exccd/wire"
)

// activeNetParams is a pointer to the parameters specific to the
// currently active excc network.
var activeNetParams = &mainNetParams

// params is used to group parameters for various networks such as the main
// network and test networks.
type params struct {
	*chaincfg.Params
	ExccdRPCServerPort  string
	RPCServerPort       string
	WalletRPCServerPort string
}

// mainNetParams contains parameters specific to the main network
// (wire.MainNet).  NOTE: The RPC port is intentionally different than the
// reference implementation because exccd does not handle wallet requests.  The
// separate wallet process listens on the well-known port and forwards requests
// it does not handle on to exccd.  This approach allows the wallet process
// to emulate the full reference implementation RPC API.
var mainNetParams = params{
	Params:              &chaincfg.MainNetParams,
	ExccdRPCServerPort:  "9109",
	RPCServerPort:       "9113",
	WalletRPCServerPort: "9110",
}

// testNetParams contains parameters specific to the test network (version 0)
// (wire.TestNet).  NOTE: The RPC port is intentionally different than the
// reference implementation - see the mainNetParams comment for details.
var testNetParams = params{
	Params:              &chaincfg.TestNetParams,
	ExccdRPCServerPort:  "19109",
	RPCServerPort:       "19113",
	WalletRPCServerPort: "19110",
}

// simNetParams contains parameters specific to the simulation test network
// (wire.SimNet).
var simNetParams = params{
	Params:              &chaincfg.SimNetParams,
	ExccdRPCServerPort:  "19556",
	RPCServerPort:       "19560",
	WalletRPCServerPort: "19557",
}

// netName returns the name used when referring to a excc network.  At the
// time of writing, exccd currently places blocks for testnet version 0 in the
// data and log directory "testnet", which does not match the Name field of the
// chaincfg parameters.  This function can be used to override this directory name
// as "testnet" when the passed active network matches wire.TestNet.
//
// A proper upgrade to move the data and log directories for this network to
// "testnet" is planned for the future, at which point this function can be
// removed and the network parameter's name used instead.
func netName(chainParams *params) string {
	switch chainParams.Net {
	case wire.TestNet:
		return "testnet"
	default:
		return chainParams.Name
	}
}
