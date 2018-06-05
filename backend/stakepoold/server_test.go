// Copyright (c) 2017 The Decred developers
// Copyright (c) 2018 The EXCCoin team
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.
package main

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	mrand "math/rand"
	"strconv"
	"testing"

	"github.com/EXCCoin/exccd/chaincfg"
	"github.com/EXCCoin/exccd/chaincfg/chainhash"
	"github.com/EXCCoin/exccstakepool/backend/stakepoold/userdata"
)

func TestCalculateFeeAddresses(t *testing.T) {
	xpubStr := "xpub6BpcvynN5SeJND9vFj1noMPwb4SWqjKGmvPq65GToQnob9UbhirCetBWzXm7c3oXPzG9r5VcaPWcaFw6hgizKn6rFk29sLShxwKYFrf5LHR"
	firstAddrs := []string{
		"22tv7nd31sMmD8BpcVRJAWQLqYCjaCuqpWpz",
		"22u4VtFDzXDT517DFfCLRM8i5t814pZXePBK",
		"22u8LM7siis3LUTKfJSnhoHAgh1d1igVdoZi",
		"22u4Qur3XjYLpHJyxxvrmFXRTzCmANFTDUYX",
		"22u7vimxVWBb3axtC2uHyCTprXs74pxC6n2o",
		"22tt57MfvPr3aotqKwQNKiMc12jhSPCVBbzr",
		"22u1VbEvX9Kapqukc7oshGcnt2ZCCeaYEDCP",
		"22tyuR7Ais8NZirdkToEo7qruMZ6maPqXSG1",
	}
	params := &chaincfg.MainNetParams

	// calculateFeeAddresses is currently hard-coded to return 10,000 addresses
	numAddr := 10000
	addrs, err := calculateFeeAddresses(xpubStr, params)
	if err != nil {
		t.Error("calculateFeeAddresses failed with ", err)
	}
	if len(addrs) != numAddr {
		t.Errorf("expected %d addresses, got %d", numAddr, len(addrs))
	}

	// Check that the first few addresses are in the map. NOTE: don't even think
	// about doing a range over the map as the order is random
	for _, addr := range firstAddrs {
		if _, ok := addrs[addr]; !ok {
			t.Errorf("Did not find address %s in derived address map", addr)
		}
	}

	// empty (i.e. invalid) xpubStr
	addrs, err = calculateFeeAddresses("", params)
	if err == nil {
		t.Error("calculateFeeAddresses did not error with empty extended key")
	}
	if len(addrs) != 0 {
		t.Errorf("expected empty map, actual length %d", len(addrs))
	}

	// wrong network
	expectedErr := fmt.Errorf("extended public key is for wrong network")
	addrs, err = calculateFeeAddresses(xpubStr, &chaincfg.TestNet2Params)
	if err == nil {
		t.Error("calculateFeeAddresses did not error with wrong network params")
	}
	if err.Error() != expectedErr.Error() {
		t.Errorf("expected error %v, got %v", expectedErr, err)
	}
	if len(addrs) != 0 {
		t.Errorf("expected empty map, actual length %d", len(addrs))
	}
}

func randomBytes(length int) []byte {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		panic(err.Error())
	}
	return b
}

func randString(n int) string {
	b := make([]byte, n)
	const addressLetters = "123456789abcdefghijkmnpQRSTUVWXYZ"
	for i := range b {
		b[i] = addressLetters[mrand.Intn(len(addressLetters))]
	}
	return string(b)
}

var (
	c  *appContext
	wt WinningTicketsForBlock
)

func init() {

	c = &appContext{
		liveTicketsMSA: make(map[chainhash.Hash]string),
		votingConfig: &VotingConfig{
			VoteBits:         1,
			VoteBitsExtended: "05000000",
			VoteVersion:      5,
		},
		userVotingConfig: make(map[string]userdata.UserVotingConfig),
		testing:          true,
	}

	// Create users
	userCount := 10000
	// leave out last 5, as they will be inserted when tickets are generated
	for i := 0; i < userCount-5; i++ {
		msa := "Tc" + randString(33)
		c.userVotingConfig[msa] = userdata.UserVotingConfig{
			Userid:          int64(i),
			MultiSigAddress: msa,
			VoteBits:        c.votingConfig.VoteBits,
			VoteBitsVersion: c.votingConfig.VoteVersion,
		}
	}

	// Create a pool of tickets around expected size
	ticketCount := 49000
	for i := 0; i < ticketCount; i++ {
		b := randomBytes(4)
		uid := int(binary.LittleEndian.Uint32(b)) % userCount
		msa := strconv.Itoa(uid | 1<<31)
		ticket := &chainhash.Hash{b[0], b[1], b[2], b[3]}

		// use ticket as the key
		c.liveTicketsMSA[*ticket] = msa

		// last 5 tickets win
		if i > ticketCount-6 {
			wt.winningTickets = append(wt.winningTickets, ticket)
			c.userVotingConfig[msa] = userdata.UserVotingConfig{}
		}
	}
}

func BenchmarkProcessWinningTickets(b *testing.B) {
	for n := 0; n < b.N; n++ {
		c.processWinningTickets(wt)
	}
}
