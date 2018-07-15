package note

import (
	"net/http"
	"public/log"
	"module/global"
	//"module/tools"
)


func getnote(w http.ResponseWriter, r *http.Request) {
	sess,_ := global.Session.SessionStart(w, r)
	//defer sess.SessionRelease(w)
	n := NewNote(r.URL.Path)
	n.LoadData()
	au := sess.Get("user")//.(user.User)
	log.Json(au)
	if len(n.EditTime)!= 0 {
		n.Show()
		err := global.Template(w,"note/content.html",map[string]interface{}{"Uri": r.RequestURI,"Note": n,"Level": 0,"User": au})
		log.Info(err)
	}else{
		w.WriteHeader(http.StatusNotFound)
	}
}


func upnote(w http.ResponseWriter, r *http.Request) {
}