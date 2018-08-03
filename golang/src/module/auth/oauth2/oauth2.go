// golang oauth2 define.
//
package oauth2

import (
	"fmt"
	"errors"
	"net/http"
	"math/rand"
	"database/sql"
	"golang.org/x/oauth2"
	"public/router"
	"module/global"
)

// Ouath2 handle
type Oauth2 interface {
	// Set Ouath2 config
	Config(*oauth2.Config) *oauth2.Config
	// Get redirect Addr
	Redirect(string) string
	// Handle callback request
	Callback(*http.Request) (*AuthInfo,error)
}

var (
	ErrOauthState		=	errors.New("invalid oauth state")
	ErrOauthCode		=	errors.New("Code exchange failed")
)

var (
	stmtQueryOauth2Source		*sql.Stmt
)

var rlogin *router.Mux
var rcallback *router.Mux

func init() {
	rlogin = router.New()
	rcallback = router.New()
}

// Reload oauth2 config
func Reload() error {
	stmtQueryOauth2Source	=	global.Stmt("SELECT Name,ClientID,ClientSecret FROM tb_auth_oauth2_source;")
	fmt.Println("init oauth2 router")
	return loadrouter()
}

func getRandomString() string {
	letters := []rune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXY")
	result := make([]rune, 16)
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}

