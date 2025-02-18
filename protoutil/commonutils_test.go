/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package protoutil

import (
	"bytes"
	"errors"
	"testing"

	cb "github.com/arogyaGurkha/fabric-protos-go/common"
	"github.com/golang/protobuf/proto"
	pb "github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric/common/crypto"
	"github.com/hyperledger/fabric/protoutil/fakes"
	"github.com/stretchr/testify/assert"
)

//go:generate counterfeiter -o fakes/signer_serializer.go --fake-name SignerSerializer . signerSerializer

type signerSerializer interface {
	Signer
}

func TestNonceRandomness(t *testing.T) {
	n1, err := CreateNonce()
	if err != nil {
		t.Fatal(err)
	}
	n2, err := CreateNonce()
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Equal(n1, n2) {
		t.Fatalf("Expected nonces to be different, got %x and %x", n1, n2)
	}
}

func TestNonceLength(t *testing.T) {
	n, err := CreateNonce()
	if err != nil {
		t.Fatal(err)
	}
	actual := len(n)
	expected := crypto.NonceSize
	if actual != expected {
		t.Fatalf("Expected nonce to be of size %d, got %d instead", expected, actual)
	}

}

func TestUnmarshalPayload(t *testing.T) {
	var payload *cb.Payload
	good, _ := proto.Marshal(&cb.Payload{
		Data: []byte("payload"),
	})
	payload, err := UnmarshalPayload(good)
	assert.NoError(t, err, "Unexpected error unmarshaling payload")
	assert.NotNil(t, payload, "Payload should not be nil")
	payload = UnmarshalPayloadOrPanic(good)
	assert.NotNil(t, payload, "Payload should not be nil")

	bad := []byte("bad payload")
	assert.Panics(t, func() {
		_ = UnmarshalPayloadOrPanic(bad)
	}, "Expected panic unmarshaling malformed payload")

}

func TestUnmarshalSignatureHeader(t *testing.T) {
	t.Run("invalid header", func(t *testing.T) {
		sighdrBytes := []byte("invalid signature header")
		_, err := UnmarshalSignatureHeader(sighdrBytes)
		assert.Error(t, err, "Expected unmarshaling error")
	})

	t.Run("valid empty header", func(t *testing.T) {
		sighdr := &cb.SignatureHeader{}
		sighdrBytes := MarshalOrPanic(sighdr)
		sighdr, err := UnmarshalSignatureHeader(sighdrBytes)
		assert.NoError(t, err, "Unexpected error unmarshaling signature header")
		assert.Nil(t, sighdr.Creator)
		assert.Nil(t, sighdr.Nonce)
	})

	t.Run("valid header", func(t *testing.T) {
		sighdr := &cb.SignatureHeader{
			Creator: []byte("creator"),
			Nonce:   []byte("nonce"),
		}
		sighdrBytes := MarshalOrPanic(sighdr)
		sighdr, err := UnmarshalSignatureHeader(sighdrBytes)
		assert.NoError(t, err, "Unexpected error unmarshaling signature header")
		assert.Equal(t, []byte("creator"), sighdr.Creator)
		assert.Equal(t, []byte("nonce"), sighdr.Nonce)
	})
}

func TestUnmarshalSignatureHeaderOrPanic(t *testing.T) {

	t.Run("panic due to invalid header", func(t *testing.T) {
		sighdrBytes := []byte("invalid signature header")
		assert.Panics(t, func() {
			UnmarshalSignatureHeaderOrPanic(sighdrBytes)
		}, "Expected panic with invalid header")
	})

	t.Run("no panic as the header is valid", func(t *testing.T) {
		sighdr := &cb.SignatureHeader{}
		sighdrBytes := MarshalOrPanic(sighdr)
		sighdr = UnmarshalSignatureHeaderOrPanic(sighdrBytes)
		assert.Nil(t, sighdr.Creator)
		assert.Nil(t, sighdr.Nonce)
	})
}

func TestUnmarshalEnvelope(t *testing.T) {
	var env *cb.Envelope
	good, _ := proto.Marshal(&cb.Envelope{})
	env, err := UnmarshalEnvelope(good)
	assert.NoError(t, err, "Unexpected error unmarshaling envelope")
	assert.NotNil(t, env, "Envelope should not be nil")
	env = UnmarshalEnvelopeOrPanic(good)
	assert.NotNil(t, env, "Envelope should not be nil")

	bad := []byte("bad envelope")
	assert.Panics(t, func() {
		_ = UnmarshalEnvelopeOrPanic(bad)
	}, "Expected panic unmarshaling malformed envelope")

}

func TestUnmarshalBlock(t *testing.T) {
	var env *cb.Block
	good, _ := proto.Marshal(&cb.Block{})
	env, err := UnmarshalBlock(good)
	assert.NoError(t, err, "Unexpected error unmarshaling block")
	assert.NotNil(t, env, "Block should not be nil")
	env = UnmarshalBlockOrPanic(good)
	assert.NotNil(t, env, "Block should not be nil")

	bad := []byte("bad block")
	assert.Panics(t, func() {
		_ = UnmarshalBlockOrPanic(bad)
	}, "Expected panic unmarshaling malformed block")

}

func TestUnmarshalEnvelopeOfType(t *testing.T) {
	env := &cb.Envelope{}

	env.Payload = []byte("bad payload")
	_, err := UnmarshalEnvelopeOfType(env, cb.HeaderType_CONFIG, nil)
	assert.Error(t, err, "Expected error unmarshaling malformed envelope")

	payload, _ := proto.Marshal(&cb.Payload{
		Header: nil,
	})
	env.Payload = payload
	_, err = UnmarshalEnvelopeOfType(env, cb.HeaderType_CONFIG, nil)
	assert.Error(t, err, "Expected error with missing payload header")

	payload, _ = proto.Marshal(&cb.Payload{
		Header: &cb.Header{
			ChannelHeader: []byte("bad header"),
		},
	})
	env.Payload = payload
	_, err = UnmarshalEnvelopeOfType(env, cb.HeaderType_CONFIG, nil)
	assert.Error(t, err, "Expected error for malformed channel header")

	chdr, _ := proto.Marshal(&cb.ChannelHeader{
		Type: int32(cb.HeaderType_CHAINCODE_PACKAGE),
	})
	payload, _ = proto.Marshal(&cb.Payload{
		Header: &cb.Header{
			ChannelHeader: chdr,
		},
	})
	env.Payload = payload
	_, err = UnmarshalEnvelopeOfType(env, cb.HeaderType_CONFIG, nil)
	assert.Error(t, err, "Expected error for wrong channel header type")

	chdr, _ = proto.Marshal(&cb.ChannelHeader{
		Type: int32(cb.HeaderType_CONFIG),
	})
	payload, _ = proto.Marshal(&cb.Payload{
		Header: &cb.Header{
			ChannelHeader: chdr,
		},
		Data: []byte("bad data"),
	})
	env.Payload = payload
	_, err = UnmarshalEnvelopeOfType(env, cb.HeaderType_CONFIG, &cb.ConfigEnvelope{})
	assert.Error(t, err, "Expected error for malformed payload data")

	chdr, _ = proto.Marshal(&cb.ChannelHeader{
		Type: int32(cb.HeaderType_CONFIG),
	})
	configEnv, _ := proto.Marshal(&cb.ConfigEnvelope{})
	payload, _ = proto.Marshal(&cb.Payload{
		Header: &cb.Header{
			ChannelHeader: chdr,
		},
		Data: configEnv,
	})
	env.Payload = payload
	_, err = UnmarshalEnvelopeOfType(env, cb.HeaderType_CONFIG, &cb.ConfigEnvelope{})
	assert.NoError(t, err, "Unexpected error unmarshaling envelope")

}

func TestExtractEnvelopeNilData(t *testing.T) {
	block := &cb.Block{}
	_, err := ExtractEnvelope(block, 0)
	assert.Error(t, err, "Nil data")
}

func TestExtractEnvelopeWrongIndex(t *testing.T) {
	block := testBlock()
	if _, err := ExtractEnvelope(block, len(block.GetData().Data)); err == nil {
		t.Fatal("Expected envelope extraction to fail (wrong index)")
	}
}

func TestExtractEnvelopeWrongIndexOrPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Expected envelope extraction to panic (wrong index)")
		}
	}()

	block := testBlock()
	ExtractEnvelopeOrPanic(block, len(block.GetData().Data))
}

func TestExtractEnvelope(t *testing.T) {
	if envelope, err := ExtractEnvelope(testBlock(), 0); err != nil {
		t.Fatalf("Expected envelop extraction to succeed: %s", err)
	} else if !proto.Equal(envelope, testEnvelope()) {
		t.Fatal("Expected extracted envelope to match test envelope")
	}
}

func TestExtractEnvelopeOrPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatal("Expected envelope extraction to succeed")
		}
	}()

	if !proto.Equal(ExtractEnvelopeOrPanic(testBlock(), 0), testEnvelope()) {
		t.Fatal("Expected extracted envelope to match test envelope")
	}
}

func TestExtractPayload(t *testing.T) {
	if payload, err := UnmarshalPayload(testEnvelope().Payload); err != nil {
		t.Fatalf("Expected payload extraction to succeed: %s", err)
	} else if !proto.Equal(payload, testPayload()) {
		t.Fatal("Expected extracted payload to match test payload")
	}
}

func TestExtractPayloadOrPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatal("Expected payload extraction to succeed")
		}
	}()

	if !proto.Equal(UnmarshalPayloadOrPanic(testEnvelope().Payload), testPayload()) {
		t.Fatal("Expected extracted payload to match test payload")
	}
}

func TestUnmarshalChaincodeID(t *testing.T) {
	ccname := "mychaincode"
	ccversion := "myversion"
	ccidbytes, _ := proto.Marshal(&pb.ChaincodeID{
		Name:    ccname,
		Version: ccversion,
	})
	ccid, err := UnmarshalChaincodeID(ccidbytes)
	assert.NoError(t, err)
	assert.Equal(t, ccname, ccid.Name, "Expected ccid names to match")
	assert.Equal(t, ccversion, ccid.Version, "Expected ccid versions to match")

	_, err = UnmarshalChaincodeID([]byte("bad chaincodeID"))
	assert.Error(t, err, "Expected error marshaling malformed chaincode ID")
}

func TestNewSignatureHeaderOrPanic(t *testing.T) {
	var sigHeader *cb.SignatureHeader

	id := &fakes.SignerSerializer{}
	id.SerializeReturnsOnCall(0, []byte("serialized"), nil)
	id.SerializeReturnsOnCall(1, nil, errors.New("serialize failed"))
	sigHeader = NewSignatureHeaderOrPanic(id)
	assert.NotNil(t, sigHeader, "Signature header should not be nil")

	assert.Panics(t, func() {
		_ = NewSignatureHeaderOrPanic(nil)
	}, "Expected panic with nil signer")

	assert.Panics(t, func() {
		_ = NewSignatureHeaderOrPanic(id)
	}, "Expected panic with signature header error")
}

func TestSignOrPanic(t *testing.T) {
	msg := []byte("sign me")
	signer := &fakes.SignerSerializer{}
	signer.SignReturnsOnCall(0, msg, nil)
	signer.SignReturnsOnCall(1, nil, errors.New("bad signature"))
	sig := SignOrPanic(signer, msg)
	// mock signer returns message to be signed
	assert.Equal(t, msg, sig, "Signature does not match expected value")

	assert.Panics(t, func() {
		_ = SignOrPanic(nil, []byte("sign me"))
	}, "Expected panic with nil signer")

	assert.Panics(t, func() {
		_ = SignOrPanic(signer, []byte("sign me"))
	}, "Expected panic with sign error")
}

// Helper functions

func testPayload() *cb.Payload {
	return &cb.Payload{
		Header: MakePayloadHeader(
			MakeChannelHeader(cb.HeaderType_MESSAGE, int32(1), "test", 0),
			MakeSignatureHeader([]byte("creator"), []byte("nonce"))),
		Data: []byte("test"),
	}
}

func testEnvelope() *cb.Envelope {
	// No need to set the signature
	return &cb.Envelope{Payload: MarshalOrPanic(testPayload())}
}

func testBlock() *cb.Block {
	// No need to set the block's Header, or Metadata
	return &cb.Block{
		Data: &cb.BlockData{
			Data: [][]byte{MarshalOrPanic(testEnvelope())},
		},
	}
}

func TestChannelHeader(t *testing.T) {
	makeEnvelope := func(payload *cb.Payload) *cb.Envelope {
		return &cb.Envelope{
			Payload: MarshalOrPanic(payload),
		}
	}

	_, err := ChannelHeader(makeEnvelope(&cb.Payload{
		Header: &cb.Header{
			ChannelHeader: MarshalOrPanic(&cb.ChannelHeader{
				ChannelId: "foo",
			}),
		},
	}))
	assert.NoError(t, err, "Channel header was present")

	_, err = ChannelHeader(makeEnvelope(&cb.Payload{
		Header: &cb.Header{},
	}))
	assert.Error(t, err, "ChannelHeader was missing")

	_, err = ChannelHeader(makeEnvelope(&cb.Payload{}))
	assert.Error(t, err, "Header was missing")

	_, err = ChannelHeader(&cb.Envelope{})
	assert.Error(t, err, "Payload was missing")
}

func TestIsConfigBlock(t *testing.T) {
	newBlock := func(env *cb.Envelope) *cb.Block {
		return &cb.Block{
			Data: &cb.BlockData{
				Data: [][]byte{MarshalOrPanic(env)},
			},
		}
	}

	newConfigEnv := func(envType int32) *cb.Envelope {
		return &cb.Envelope{
			Payload: MarshalOrPanic(&cb.Payload{
				Header: &cb.Header{
					ChannelHeader: MarshalOrPanic(&cb.ChannelHeader{
						Type:      envType,
						ChannelId: "test-chain",
					}),
				},
				Data: []byte("test bytes"),
			}), // common.Payload
		} // LastUpdate
	}

	// scenario 1: CONFIG envelope
	envType := int32(cb.HeaderType_CONFIG)
	env := newConfigEnv(envType)
	block := newBlock(env)

	result := IsConfigBlock(block)
	assert.True(t, result, "IsConfigBlock returns true for blocks with CONFIG envelope")

	// scenario 2: ORDERER_TRANSACTION envelope
	envType = int32(cb.HeaderType_ORDERER_TRANSACTION)
	env = newConfigEnv(envType)
	block = newBlock(env)

	result = IsConfigBlock(block)
	assert.True(t, result, "IsConfigBlock returns true for blocks with ORDERER_TRANSACTION envelope")

	// scenario 3: MESSAGE envelope
	envType = int32(cb.HeaderType_MESSAGE)
	env = newConfigEnv(envType)
	block = newBlock(env)

	result = IsConfigBlock(block)
	assert.False(t, result, "IsConfigBlock returns false for blocks with MESSAGE envelope")
}

func TestEnvelopeToConfigUpdate(t *testing.T) {
	makeEnv := func(data []byte) *cb.Envelope {
		return &cb.Envelope{
			Payload: MarshalOrPanic(&cb.Payload{
				Header: &cb.Header{
					ChannelHeader: MarshalOrPanic(&cb.ChannelHeader{
						Type:      int32(cb.HeaderType_CONFIG_UPDATE),
						ChannelId: "test-chain",
					}),
				},
				Data: data,
			}), // common.Payload
		} // LastUpdate
	}

	// scenario 1: for valid envelopes
	configUpdateEnv := &cb.ConfigUpdateEnvelope{}
	env := makeEnv(MarshalOrPanic(configUpdateEnv))
	result, err := EnvelopeToConfigUpdate(env)

	assert.NoError(t, err, "EnvelopeToConfigUpdate runs without error for valid CONFIG_UPDATE envelope")
	assert.Equal(t, configUpdateEnv, result, "Correct configUpdateEnvelope returned")

	// scenario 2: for invalid envelopes
	env = makeEnv([]byte("test bytes"))
	_, err = EnvelopeToConfigUpdate(env)

	assert.Error(t, err, "EnvelopeToConfigUpdate fails with error for invalid CONFIG_UPDATE envelope")
}

func TestGetRandomNonce(t *testing.T) {
	key1, err := getRandomNonce()
	assert.NoErrorf(t, err, "error getting random bytes")
	assert.Len(t, key1, crypto.NonceSize)
}
