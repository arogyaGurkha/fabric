/*
Copyright IBM Corp. 2017 All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package genesis

import (
	gurkhaB "github.com/arogyaGurkha/fabric-protos-go/common"
	"github.com/hyperledger/fabric/protoutil"
)

const (
	msgVersion = int32(1)

	// These values are fixed for the genesis block.
	epoch = 0
)

// Factory facilitates the creation of genesis blocks.
type Factory interface {
	// Block returns a genesis block for a given channel ID.
	Block(channelID string) *gurkhaB.Block
}

type factory struct {
	channelGroup *gurkhaB.ConfigGroup
}

// NewFactoryImpl creates a new Factory.
func NewFactoryImpl(channelGroup *gurkhaB.ConfigGroup) Factory {
	return &factory{channelGroup: channelGroup}
}

// Block constructs and returns a genesis block for a given channel ID.
func (f *factory) Block(channelID string) *gurkhaB.Block {
	payloadChannelHeader := protoutil.MakeChannelHeader(gurkhaB.HeaderType_CONFIG, msgVersion, channelID, epoch)
	payloadSignatureHeader := protoutil.MakeSignatureHeader(nil, protoutil.CreateNonceOrPanic())
	protoutil.SetTxID(payloadChannelHeader, payloadSignatureHeader)
	payloadHeader := protoutil.MakePayloadHeader(payloadChannelHeader, payloadSignatureHeader)
	payload := &gurkhaB.Payload{Header: (*gurkhaB.Header)(payloadHeader), Data: protoutil.MarshalOrPanic(&gurkhaB.ConfigEnvelope{Config: &gurkhaB.Config{ChannelGroup: f.channelGroup}})}
	envelope := &gurkhaB.Envelope{Payload: protoutil.MarshalOrPanic(payload), Signature: nil}

	block := protoutil.NewBlock(0, nil)
	block.Data = &gurkhaB.BlockData{Data: [][]byte{protoutil.MarshalOrPanic(envelope)}}
	block.Extension = &gurkhaB.BlockExtension{ExtensionData: [][]byte{protoutil.MarshalOrPanic(envelope)}}
	block.Header.DataHash = protoutil.BlockDataHash(block.Data)
	block.Metadata.Metadata[gurkhaB.BlockMetadataIndex_LAST_CONFIG] = protoutil.MarshalOrPanic(&gurkhaB.Metadata{
		Value: protoutil.MarshalOrPanic(&gurkhaB.LastConfig{Index: 0}),
	})
	block.Metadata.Metadata[gurkhaB.BlockMetadataIndex_SIGNATURES] = protoutil.MarshalOrPanic(&gurkhaB.Metadata{
		Value: protoutil.MarshalOrPanic(&gurkhaB.OrdererBlockMetadata{
			LastConfig: &gurkhaB.LastConfig{Index: 0},
		}),
	})
	return block
}
