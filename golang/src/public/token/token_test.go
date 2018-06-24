package token

import (
	"fmt"
	"time"
	"testing"
	"public/token"
)

var ts string
func TestDecode(t *testing.T) {
	hmacSampleSecret := []byte("secret")
	to := token.NewWithClaims( token.SigningMethodHS256, &token.MapClaims{
		"foo": "bar",
		"nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	})
/*	to := token.Token{
		Header:	map[string]interface{}{
			"typ": "JWT",
			"alg": token.SigningMethodHS256.Alg(),
		},
		Claims:	&token.MapClaims{
			"foo": "bar",
			"nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
		},
		Method: token.SigningMethodHS256,
	}*/
	tokenString, err := to.SignedString(hmacSampleSecret)
	ts = tokenString
	t.Log(tokenString,err)
}

func TestEecode(t *testing.T) {
	tokenString := ts
	cl := &token.MapClaims{}
	to, err := token.ParseWithClaims(tokenString,cl,func(to *token.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := to.Method.(*token.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", to.Header["alg"])
		}
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte("secret"), nil
	})

	if to.Valid {
		t.Log(cl)
	} else {
		t.Log(err)
	}
}