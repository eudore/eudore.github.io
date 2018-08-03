package user;

import (
	"time"
	"net/http"
	"module/global"
)

func Logout(w http.ResponseWriter,r *http.Request) {
	// session
	sess,_ := global.Session.SessionStart(w, r)
	defer sess.SessionRelease(w)
	sess.Flush()
	// token
	cookie := &http.Cookie{
		Name:		"t",
		Value:		"",
		Path:		"/",
		HttpOnly:	true,
		Secure:		true,
		Domain:		"wejass.com",
		Expires:	time.Now().Add(-1 * time.Second),
	}
	http.SetCookie(w,cookie)
}


