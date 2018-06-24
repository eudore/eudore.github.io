package token

import (
	"bytes"
	"encoding/json"
)

type Claims interface {
	Valid() error
	Marshal() ([]byte, error)
	Unmarshal(data []byte) error
}


type MapClaims map[string]interface{}

func (c *MapClaims) Valid() error {
	return nil
}

func (c *MapClaims) Marshal() ([]byte, error) {
	return json.Marshal(c)
}

func (c *MapClaims) Unmarshal(data []byte) error {
	dec := json.NewDecoder(bytes.NewBuffer(data))
	return dec.Decode(c)
}