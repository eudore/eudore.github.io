package user;

import (
	"net/http"
	"module/global"
)

func Logout(w http.ResponseWriter,r *http.Request) {
	sess,_ := global.Session.SessionStart(w, r)
	defer sess.SessionRelease(w)
	sess.Flush()
	// flush cookies
}
