/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package blockledger

import (
	//cb "github.com/hyperledger/fabric-protos-go/common"
	gurkhaB "github.com/arogyaGurkha/fabric-protos-go/common"
	ab "github.com/hyperledger/fabric-protos-go/orderer"
)

// Factory retrieves or creates new ledgers by channelID
type Factory interface {
	// GetOrCreate gets an existing ledger (if it exists)
	// or creates it if it does not
	GetOrCreate(channelID string) (ReadWriter, error)

	// ChannelIDs returns the channel IDs the Factory is aware of
	ChannelIDs() []string

	// Close releases all resources acquired by the factory
	Close()
}

// Iterator is useful for a chain Reader to stream blocks as they are created
type Iterator interface {
	// Next blocks until there is a new block available, or returns an error if
	// the next block is no longer retrievable
	Next() (*gurkhaB.Block, gurkhaB.Status)
	// Close releases resources acquired by the Iterator
	Close()
}

// Reader allows the caller to inspect the ledger
type Reader interface {
	// Iterator returns an Iterator, as specified by an ab.SeekInfo message, and
	// its starting block number
	Iterator(startType *ab.SeekPosition) (Iterator, uint64)
	// Height returns the number of blocks on the ledger
	Height() uint64
	// retrieve blockByNumber
	RetrieveBlockByNumber(blockNumber uint64) (*gurkhaB.Block, error)
}

// Writer allows the caller to modify the ledger
type Writer interface {
	// Append a new block to the ledger
	Append(block *gurkhaB.Block) error
}

// ReadWriter encapsulates the read/write functions of the ledger
type ReadWriter interface {
	Reader
	Writer
}
