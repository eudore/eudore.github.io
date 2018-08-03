package user;

import (
	"fmt"
	"time"
	"strings"
	"net/http"
	"encoding/json"
	"public/log"
	"public/token"
)


func Login(w http.ResponseWriter,r *http.Request) {
	// check referer
	if !strings.HasPrefix(r.Header.Get("Referer"),"https://w") {
//		return
	}
	// get auth token
	r.ParseForm()
	authtoken := r.Form.Get(AuthToken)
	redirect := r.Form.Get(AuthRedirect)
	log.Info(redirect)
	data := make(map[string]interface{})
	if authtoken != "" {
		cl := &token.MapClaims{}
		to, err := token.ParseWithClaims(authtoken, cl, func(to *token.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := to.Method.(*token.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", to.Header["alg"])
			}
			// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
			return []byte("secret"), nil
		})
		if !to.Valid || err != nil {
			log.Error(err)
		}
		data = *cl
		log.Info(authtoken)
		log.Json(data)
	}
	// get user
	uid,ok := data["uid"]
	if !ok {
		return
	}
	user := NewUser(int(uid.(float64)))

	// set user token
	hmacSampleSecret := []byte("secret")
	to := token.NewWithClaims( token.SigningMethodHS256, &token.MapClaims{
		"id":		user.Uid,
		"name":		user.Name,
		"policy":	user.Policy,
		"expires":	time.Now().Add(1000 * time.Second).Unix(),
	})
	tokenString, _ := to.SignedString(hmacSampleSecret)
	cookie := &http.Cookie{
		Name:		"t",
		Value:		tokenString,
		Path:		"/",
		HttpOnly:	true,
		Secure:		true,
		Domain:		"wejass.com",
		Expires:	time.Now().Add(1000 * time.Second),
	}
	http.SetCookie(w, cookie)
	// return repoonse
	switch r.Form.Get("format") {
	case "json":
		responseBody,_ := json.Marshal(map[string]interface{}{
			AuthRedirect: redirect,
		})
		w.Write(responseBody)
	case "redirect":
		log.Info("redirect", redirect)
		http.Redirect(w, r, redirect, http.StatusFound)
	default:
		http.Error(w , http.StatusText(405), http.StatusMethodNotAllowed)
	}
}
