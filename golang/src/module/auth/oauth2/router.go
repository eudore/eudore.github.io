package oauth2

import (
	"fmt"
	"time"
	"net/http"
	"net/url"
	"public/log"
	"public/router"
	"module/global"
)

const (
	// state cookie name
	CookieState		=	"oauth2state"
	CookieRedirect	=	"oauth2redirect"
	CookiePath		=	"/auth/oauth2/"
	CookieDomain	=	"www.wejass.com"
	TokenOauth2		=	"authtoken"
	TokenRediect	=	"redirect"
	CallbackUrl		=	"https://www.wejass.com/auth/oauth2/callback/"
)


// Get Ouath2 login and callback router
func GetRouter() (*router.Mux,*router.Mux) {
	return rlogin, rcallback
}

func loadrouter() error {
	rows, err := global.Sql.Query("SELECT Name,ClientID,ClientSecret FROM tb_auth_oauth2_source;")
	if err != nil {
		return err
	}
	var name, id, secret string
	for rows.Next(){
		rows.Scan(&name, &id, &secret)
		o,err := NewOuath2(name)
		if err != nil {
			fmt.Println(err,name)
			continue
		}
		fmt.Println("init oauth2", name)
		// use default config
		cf := o.Config(nil)
		cf.RedirectURL = CallbackUrl + name
		cf.ClientID = id
		cf.ClientSecret = secret
		rlogin.GetFunc("/" + name, redirectfunc(o))
		rcallback.GetFunc("/" + name, callbackfunc(o))
	}
	return nil
}

// Oauth2 login redirect func
func redirectfunc(o Oauth2) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter,r *http.Request) {
			// random stats
			oauthState := getRandomString()
			cookie := &http.Cookie{
				Name:		CookieState,
				Value:		oauthState,
				Path:		CookiePath,
				HttpOnly:	true,
				Secure:		true,
				Domain:		CookieDomain,
				Expires:	time.Now().Add(100 * time.Second),
			}
			http.SetCookie(w,cookie)
			// save redirect
			redirect := r.FormValue("redirect")
			cookie = &http.Cookie{
				Name:		CookieRedirect,
				Value:		redirect,
				Path:		CookiePath,
				HttpOnly:	true,
				Secure:		true,
				Domain:		CookieDomain,
				Expires:	time.Now().Add(100 * time.Second),
			}
			http.SetCookie(w, cookie)
			// redirect oauth2
			url := o.Redirect(oauthState)
			http.Redirect(w, r, url, http.StatusFound)
			fmt.Println("redirect:", url)
		}
}

// Oauth2 callback func
func callbackfunc(o Oauth2) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter,r *http.Request) {
			// check stats
			state := r.FormValue("state")
			cookie, err := r.Cookie(CookieState)
			if state != cookie.Value || err != nil {
				return
			}
			log.Info("state ",state)
			// load uid
			au, err := o.Callback(r)
			err = au.getuid()
			oauth2_token, err := au.GetJwt()
			if err != nil {
				http.Error(w,http.StatusText(403),http.StatusUnauthorized)
			}
			log.Info("oauth2_token ",oauth2_token)

			// redirect
			redirect, _ := r.Cookie(CookieRedirect)
			data := url.Values{
				TokenOauth2:		{oauth2_token},
				TokenRediect:		{redirect.Value},
				"format":			{"redirect"},
			}
			uri := "/auth/user/login?"
			if au.Uid == -1 {
				// sign up
				uri = "/auth/user/signup?"
			}
			http.Redirect(w, r, uri + data.Encode(), http.StatusFound)
		}
}
