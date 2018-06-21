package user;

import (
	"strings"
	"net/http"
	"public/log"
	"module/global"
	"module/auth/oauth2"
)

func Login(w http.ResponseWriter,r *http.Request) {
	// check referer
	if !strings.HasPrefix(r.Header.Get("Referer"),"https://w") {
		return
	}
	sess,_ := global.Session.SessionStart(w, r)
	defer sess.SessionRelease(w)
	// load oauth2 user info
	au := sess.Get("oauth2")
	if au != nil {
		sess.Set("user", *(NewUser(au.(oauth2.AuthInfo).Uid)))
		sess.Delete("oauth2")
		log.Json(sess.Get("user"))
	}
	// redirect
	http.Redirect(w, r, sess.Get("redirect").(string), http.StatusFound)
	sess.Delete("redirect")
}
