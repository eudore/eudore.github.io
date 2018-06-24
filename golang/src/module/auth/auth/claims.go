package auth

import (
	"encoding/json"
)


type AuthClaims struct {
	Path	string		`json:"path"`
	Sub		[]string	`json:"sub"`
	Method	uint		`json:"method"`
	Expries	int64		`json:"expries"`
}

func (c *AuthClaims) Valid() error {	
	return nil
}

func (c *AuthClaims) Marshal() ([]byte, error) {
	return json.Marshal(c)
}
func (c *AuthClaims) Unmarshal(data []byte) error {
	return json.Unmarshal(data,c)
}

type UserClaims struct {
	Name	string
	Uid		int
	Expries	int64		`json:"expries"`
}