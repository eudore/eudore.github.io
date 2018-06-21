package oauth2


import (
	"errors"
)

var (
	// ErrUnknownOauth2 unknown oauth2 error
	ErrUnknownOauth2 = errors.New("unknow oauth2")
)

const (
	Oauth2Eudore	=	"eudore"
	Oauth2Github	=	"github"
	Oauth2Google	=	"google"
	Oauth2Gitlab	=	"gitlab"
)
var oauth2source [4]string = [4]string{"eudore","github","google","gitlab"}

func NewOuath2(name string) (Oauth2,error) {
	switch name {
	case Oauth2Eudore:
		return newOauth2Eudore(), nil
	case Oauth2Github:
		return newOauth2Github(), nil
	case Oauth2Google:
		return newOauth2Google(), nil
	case Oauth2Gitlab:
		return newOauth2Gitlab(), nil
	default:
		return nil, ErrUnknownOauth2
	}
}

func getsource(s string) int {
	for i,v := range oauth2source {
		if v == s {
			return i
		}
	}
	return -1
}