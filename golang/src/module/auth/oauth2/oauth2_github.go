package oauth2

import (
	"strconv"
	"net/http"
	"net/url"
	"io/ioutil"
	"encoding/json"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var (
	sourcegithub		=	getsource("github")
)

type UserGithub struct {
	Id		int
	Login	string
	Email	string
}

func (u *UserGithub) ToUser() *AuthInfo{
	return &AuthInfo{
		Source:	sourcegithub,
		Id:		strconv.Itoa(u.Id),
		Name:	u.Login,
		Email:	u.Email,
	}
}

type Oauth2GithubHandle struct {
	config		*oauth2.Config
}

func newOauth2Github() Oauth2 {
	return &Oauth2GithubHandle{}
}

func (o *Oauth2GithubHandle) Config(config *oauth2.Config) *oauth2.Config {
	if config == nil {
		o.config = &oauth2.Config{
			Scopes: []string{"user:email"},
			Endpoint: github.Endpoint,
		}
	}else{
		o.config = config	
	}
	return o.config
}

func (o *Oauth2GithubHandle) Redirect(stats string) string{
	return o.config.AuthCodeURL(stats)
}

func (o *Oauth2GithubHandle) Callback(r *http.Request) (*AuthInfo,error) {
	// get code
	code := r.FormValue("code")
	response, _ := http.PostForm("https://github.com/login/oauth/access_token", url.Values{
		"client_id":     {o.config.ClientID},
		"client_secret": {o.config.ClientSecret},
		"code":          {code},
	})
	defer response.Body.Close()
	// get user info
	contents, _ := ioutil.ReadAll(response.Body)
	res, _ := http.Get("https://api.github.com/user?" + string(contents))
	defer res.Body.Close()
	con, _ := ioutil.ReadAll(res.Body)
	var ug UserGithub
	json.Unmarshal(con,&ug)
	return ug.ToUser(),nil
}