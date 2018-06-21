package oauth2

import (
	"strconv"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/gitlab"
)

var (
	sourcegitlab		=	getsource("gitlab")
)

type UserGitlab struct {
	Id			int
	Username	string
	Email		string
}

func (u *UserGitlab) ToUser() *AuthInfo{
	return &AuthInfo{
		Source:	sourcegitlab,
		Id:		strconv.Itoa(u.Id),
		Name:	u.Username,
		Email:	u.Email,
	}
}

type Oauth2GitlabHandle struct {
	config		*oauth2.Config
}

func newOauth2Gitlab() Oauth2 {
	return &Oauth2GitlabHandle{}
}

func (o *Oauth2GitlabHandle) Config(config *oauth2.Config) *oauth2.Config {
	if config == nil {
		o.config = &oauth2.Config{
			Scopes: []string{"read_user"},
			Endpoint: gitlab.Endpoint,
		}
	}else{
		o.config = config	
	}
	return o.config
}

func (o *Oauth2GitlabHandle) Redirect(stats string) string{
	return o.config.AuthCodeURL(stats)
}

func (o *Oauth2GitlabHandle) Callback(r *http.Request) (*AuthInfo,error) {
	// get code
	code := r.FormValue("code")
	token, err := o.config.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil,ErrOauthCode
	}
	// get user info
	response, err := http.Get("https://gitlab.com/api/v4/user?access_token=" + token.AccessToken)
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	var ug UserGitlab
	json.Unmarshal(contents,&ug)
	return ug.ToUser(),nil
}