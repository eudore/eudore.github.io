package oauth2

import (
	"fmt"
	"time"
	"net/http"
	"public/log"
	"public/router"
	"module/global"
	"module/tools"
)


func GetRouter() (*router.Mux,*router.Mux) {
	return rlogin,rcallback
}

func loadoauth2() error {
	rows, err := global.Sql.Query("SELECT Name,ClientID,ClientSecret FROM tb_auth_oauth2_source;")
	if err != nil {
		return err
	}
	var name,id,secret string
	for rows.Next(){
		rows.Scan(&name,&id,&secret)
		o,err := NewOuath2(name)
		if err != nil {
			fmt.Println(err,name)
			continue
		}
		fmt.Println("init oauth2",name)
		// use default config
		cf := o.Config(nil)
		cf.RedirectURL = "https://wejass.com:8081/auth/oauth2/callback/"+name		//access domain ?
		cf.ClientID = id
		cf.ClientSecret = secret
		rlogin.GetFunc("/"+name,redirectfunc(o))
		rcallback.GetFunc("/"+name,callbackfunc(o))
	}
	return nil
}

func redirectfunc(o Oauth2) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter,r *http.Request) {
			// random stats
			oauthState := getRandomString()
			cookie := &http.Cookie{
				Name:		"oauth2",
				Value:		oauthState,
				Path:		"/auth/oauth2/",
				HttpOnly:	true,
				Secure:		true,
				Domain:		"wejass.com",
				Expires:	time.Now().Add(100 * time.Second),
			}
			http.SetCookie(w,cookie)
			// redirect oauth2
			url := o.Redirect(oauthState)
			http.Redirect(w, r, url, http.StatusFound)
			fmt.Println("redirect:",url)
		}
}

func callbackfunc(o Oauth2) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter,r *http.Request) {
			// check stats
			state := r.FormValue("state")
			cookie, err := r.Cookie("oauth2")
			if state != cookie.Value || err != nil {
				return
			}
			log.Info(tools.GetRealClientIP(r))
			log.Info(state)
			// load uid
			au, err := o.Callback(r)
			err = au.getuid()
			log.Info(err)
			log.Json(au)

			sess,_ := global.Session.SessionStart(w, r)
			sess.Set("oauth2",*au)
			sess.SessionRelease(w)
			if au.Uid != -1 {
				// login
				http.Redirect(w, r, "/auth/user/login", http.StatusFound)
			}else {
				// sign up
				http.Redirect(w, r, "/auth/user/signup", http.StatusFound)
			}
			//http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		}
}
