package oauth2

import (
	"fmt"
	"errors"
	"net/http"
	"math/rand"
	"public/router"
	"golang.org/x/oauth2"
)


type Oauth2 interface {
	Config(*oauth2.Config) *oauth2.Config
	Redirect(string) string
	Callback(*http.Request) (*AuthInfo,error)
}

var (
	ErrOauthState		=	errors.New("invalid oauth state")
	ErrOauthCode		=	errors.New("Code exchange failed")
)


var rlogin *router.Mux
var rcallback *router.Mux

func init() {
	rlogin = router.New()
	rcallback = router.New()
	fmt.Println("init oauth2 router")
}

func getRandomString() string {
	letters := []rune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXY")
	result := make([]rune, 16)
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}

