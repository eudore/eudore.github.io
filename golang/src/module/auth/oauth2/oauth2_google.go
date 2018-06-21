package oauth2

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	sourcegoogle		=	getsource("google")
)

type UserGoogle struct {
	Id		string
	Name	string
	Email	string
}

func (u *UserGoogle) ToUser() *AuthInfo{
	return &AuthInfo{
		Source:	sourcegoogle,
		Id:		u.Id,
		Name:	u.Name,
		Email:	u.Email,
	}
}

type Oauth2GoogleHandle struct {
	config		*oauth2.Config
}

func newOauth2Google() Oauth2 {
	return &Oauth2GoogleHandle{}
}

func (o *Oauth2GoogleHandle) Config(config *oauth2.Config) *oauth2.Config {
	if config == nil {
		o.config = &oauth2.Config{
			Scopes: []string{"https://www.googleapis.com/auth/userinfo.profile",
				"https://www.googleapis.com/auth/userinfo.email"},
			Endpoint: google.Endpoint,
		}
	}else{
		o.config = config	
	}
	return o.config
}

func (o *Oauth2GoogleHandle) Redirect(stats string) string{
	return o.config.AuthCodeURL(stats)
}

func (o *Oauth2GoogleHandle) Callback(r *http.Request) (*AuthInfo,error) {
	code := r.FormValue("code")
	token, err := o.config.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil,ErrOauthCode
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	var ug UserGoogle
	json.Unmarshal(contents,&ug)
	return ug.ToUser(),nil
}