// Code generated by github.com/atombender/go-jsonschema, DO NOT EDIT.

package types

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type Bridge struct {
	// AgentAddress corresponds to the JSON schema field "agentAddress".
	AgentAddress *string `json:"agentAddress,omitempty" yaml:"agentAddress,omitempty"`

	// DestChain corresponds to the JSON schema field "destChain".
	DestChain BridgeChain `json:"destChain" yaml:"destChain" binding:"required"`

	// SrcChain corresponds to the JSON schema field "srcChain".
	SrcChain BridgeChain `json:"srcChain" yaml:"srcChain" binding:"required"`

	// Tokens corresponds to the JSON schema field "tokens".
	Tokens []Token `json:"tokens" yaml:"tokens" binding:"required"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *Invoke) UnmarshalJSON(b []byte) error {
	type Plain Invoke
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	if err := checkRequired(plain);err != nil {
		return err
	}
	*j = Invoke(plain)
	return nil
}

type BridgeChain struct {
	// ChainId corresponds to the JSON schema field "chain_id".
	ChainId string `json:"chain_id" yaml:"chain_id" binding:"required"`

	// Channel corresponds to the JSON schema field "channel".
	Channel *string `json:"channel,omitempty" yaml:"channel,omitempty"`

	// IsTele corresponds to the JSON schema field "is_tele".
	IsTele *bool `json:"is_tele,omitempty" yaml:"is_tele,omitempty"`

	// Name corresponds to the JSON schema field "name".
	Name string `json:"name" yaml:"name" binding:"required"`

	// Endpoint corresponds to the JSON schema field "endpoint".
	Endpoint *Invoke `json:"endpoint,omitempty" yaml:"endpoint,omitempty"`

	// Packet corresponds to the JSON schema field "packet".
	Packet Invoke `json:"packet" yaml:"packet"`

	// Proxy corresponds to the JSON schema field "proxy".
	Proxy *Invoke `json:"proxy,omitempty" yaml:"proxy,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *BridgeChain) UnmarshalJSON(b []byte) error {
	type Plain BridgeChain
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	if err := checkRequired(plain);err != nil {
		return err
	}
	*j = BridgeChain(plain)
	return nil
}

type Token struct {
	// Name corresponds to the JSON schema field "name".
	Name string `json:"name" yaml:"name"`

	// PairType  corresponds to the JSON schema field "name".
	PairType string `json:"pairType" yaml:"pairType"`

	// DestToken corresponds to the JSON schema field "destToken".
	DestToken string `json:"destToken" yaml:"destToken" binding:"required"`

	// RelayToken corresponds to the JSON schema field "relayToken".
	RelayToken *string `json:"relayToken,omitempty" yaml:"relayToken,omitempty"`

	// SrcToken corresponds to the JSON schema field "srcToken".
	SrcToken string `json:"srcToken" yaml:"srcToken" binding:"required"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *Token) UnmarshalJSON(b []byte) error {
	type Plain Token
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	if err := checkRequired(plain);err != nil {
		return err
	}
	*j = Token(plain)
	return nil
}

type Invoke struct {
	// Abi corresponds to the JSON schema field "abi".
	Abi string `json:"abi" yaml:"abi"`

	// Address corresponds to the JSON schema field "address".
	Address string `json:"address" yaml:"address"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *Bridge) UnmarshalJSON(b []byte) error {
	type Plain Bridge
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	if err := checkRequired(plain);err != nil {
		return err
	}
	*j = Bridge(plain)
	return nil
}

type Bridges struct {
	// Bridges corresponds to the JSON schema field "bridges".
	Bridges []Bridge `json:"bridges" yaml:"bridges"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *Bridges) UnmarshalJSON(b []byte) error {
	type Plain Bridges
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	if err := checkRequired(plain);err != nil {
		return err
	}
	*j = Bridges(plain)
	return nil
}

func checkRequired(str interface{})error {
	t := reflect.TypeOf(str)
	v := reflect.ValueOf(str)
	for k := 0; k < t.NumField(); k++ {
		fieldType := v.Field(k).Kind()
		if fieldType == reflect.Struct {
			checkRequired(v.Field(k).Interface())
		}
		if t.Field(k).Tag.Get("binding") == "required" {
			if v.Field(k).IsZero() {
				panic(fmt.Sprintf("field %v in Bridges: required", t.Field(k).Name))
			}
		}
	}
	return nil
}