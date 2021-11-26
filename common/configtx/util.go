/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package configtx

import (
	gurkhaB "github.com/arogyaGurkha/fabric-protos-go/common"
	"github.com/golang/protobuf/proto"
	cb "github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric/protoutil"
)

// UnmarshalConfig attempts to unmarshal bytes to a *cb.Config
func UnmarshalConfig(data []byte) (*gurkhaB.Config, error) {
	config := &gurkhaB.Config{}
	err := proto.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

// UnmarshalConfigOrPanic attempts to unmarshal bytes to a *cb.Config or panics on error
func UnmarshalConfigOrPanic(data []byte) *gurkhaB.Config {
	result, err := UnmarshalConfig(data)
	if err != nil {
		panic(err)
	}
	return result
}

// UnmarshalConfigUpdate attempts to unmarshal bytes to a *cb.ConfigUpdate
func UnmarshalConfigUpdate(data []byte) (*gurkhaB.ConfigUpdate, error) {
	configUpdate := &gurkhaB.ConfigUpdate{}
	err := proto.Unmarshal(data, configUpdate)
	if err != nil {
		return nil, err
	}
	return configUpdate, nil
}

// UnmarshalConfigUpdateOrPanic attempts to unmarshal bytes to a *cb.ConfigUpdate or panics on error
func UnmarshalConfigUpdateOrPanic(data []byte) *gurkhaB.ConfigUpdate {
	result, err := UnmarshalConfigUpdate(data)
	if err != nil {
		panic(err)
	}
	return result
}

// UnmarshalConfigUpdateEnvelope attempts to unmarshal bytes to a *cb.ConfigUpdate
func UnmarshalConfigUpdateEnvelope(data []byte) (*gurkhaB.ConfigUpdateEnvelope, error) {
	configUpdateEnvelope := &gurkhaB.ConfigUpdateEnvelope{}
	err := proto.Unmarshal(data, configUpdateEnvelope)
	if err != nil {
		return nil, err
	}
	return configUpdateEnvelope, nil
}

// UnmarshalConfigUpdateEnvelopeOrPanic attempts to unmarshal bytes to a *cb.ConfigUpdateEnvelope or panics on error
func UnmarshalConfigUpdateEnvelopeOrPanic(data []byte) *gurkhaB.ConfigUpdateEnvelope {
	result, err := UnmarshalConfigUpdateEnvelope(data)
	if err != nil {
		panic(err)
	}
	return result
}

// UnmarshalConfigEnvelope attempts to unmarshal bytes to a *cb.ConfigEnvelope
func UnmarshalConfigEnvelope(data []byte) (*gurkhaB.ConfigEnvelope, error) {
	configEnv := &gurkhaB.ConfigEnvelope{}
	err := proto.Unmarshal(data, configEnv)
	if err != nil {
		return nil, err
	}
	return configEnv, nil
}

// UnmarshalConfigEnvelopeOrPanic attempts to unmarshal bytes to a *cb.ConfigEnvelope or panics on error
func UnmarshalConfigEnvelopeOrPanic(data []byte) *gurkhaB.ConfigEnvelope {
	result, err := UnmarshalConfigEnvelope(data)
	if err != nil {
		panic(err)
	}
	return result
}

// UnmarshalConfigUpdateFromPayload unmarshals configuration update from given payload
func UnmarshalConfigUpdateFromPayload(payload *cb.Payload) (*gurkhaB.ConfigUpdate, error) {
	configEnv, err := UnmarshalConfigEnvelope(payload.Data)
	if err != nil {
		return nil, err
	}
	configUpdateEnv, err := protoutil.EnvelopeToConfigUpdate(configEnv.LastUpdate)
	if err != nil {
		return nil, err
	}

	return UnmarshalConfigUpdate(configUpdateEnv.ConfigUpdate)
}
