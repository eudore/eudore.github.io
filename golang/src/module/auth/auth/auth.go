package auth

import (
	"time"
	"strings"
	"public/token"
)

const (
	NULL	uint = 1 << iota
	GET
	POST
	PUT
	DELETE
	HEAD
	PATCH
)

func Valid(t string,uri string,method uint) bool{
	cl := &AuthClaims{}
	_, err := token.ParseWithClaims(t,cl,func(to *token.Token) (interface{}, error) {
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(to.Claims.(*AuthClaims).Path), nil
	})
	path := cl.Path
	if err == nil && cl.Expries > time.Now().Unix() &&(method & cl.Method)==method && (( path[len(path)-1] == 47 && strings.HasPrefix(uri,path) ) || uri == path ){
		return true
	}
	return false
}

func Auth(uri string,method uint,t int64) (string,error) {
	hmacSampleSecret := []byte(uri)
	return token.NewWithClaims(token.SigningMethodHS256,&AuthClaims{
			Path: uri,
			Method: method,
			Expries: t,
		}).SignedString(hmacSampleSecret)
}
