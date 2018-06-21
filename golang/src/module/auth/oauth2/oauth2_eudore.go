package oauth2

import (
	"strconv"
	"net/http"
	"net/url"
	"io/ioutil"
	"encoding/json"
	"golang.org/x/oauth2"
)

var (
	sourceeudore		=	getsource("eudore")
)

type UserEudore struct {
	Id		int
	Login	string
	Email	string
}

func (u *UserEudore) ToUser() *AuthInfo{
	return &AuthInfo{
		Source:	sourceeudore,
		Id:		strconv.Itoa(u.Id),
		Name:	u.Login,
		Email:	u.Email,
	}
}

type Oauth2EudoreHandle struct {
	config		*oauth2.Config
}

func newOauth2Eudore() Oauth2 {
	return &Oauth2EudoreHandle{}
}

func (o *Oauth2EudoreHandle) Config(config *oauth2.Config) *oauth2.Config {
	if config == nil {
		o.config = &oauth2.Config{
			Scopes: []string{"user:email"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://wejass.com:8081/auth/user/auth",
				TokenURL: "https://wejass.com:8081/auth/user/token",
			},
		}
	}else{
		o.config = config	
	}
	return o.config
}

func (o *Oauth2EudoreHandle) Redirect(stats string) string{
	return o.config.AuthCodeURL(stats)
}

func (o *Oauth2EudoreHandle) Callback(r *http.Request) (*AuthInfo,error){
	// get code
	code := r.FormValue("code")
	response, _ := http.PostForm(o.config.Endpoint.TokenURL, url.Values{
		"client_id":     {o.config.ClientID},
		"client_secret": {o.config.ClientSecret},
		"code":          {code},
	})
	defer response.Body.Close()
	// get user info
	contents, _ := ioutil.ReadAll(response.Body)
	res, _ := http.Get("https://wejass.com:8081/auth/user?" + string(contents))
	defer res.Body.Close()
	con, _ := ioutil.ReadAll(res.Body)
	var ue UserEudore
	json.Unmarshal(con,&ue)
	return ue.ToUser(),nil
}
