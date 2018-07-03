package note

import (
	"strings"
	"net/http"
	"html/template"
	"public/log"
	"module/global"
	"module/tools"
)

func note(w http.ResponseWriter, r *http.Request) {
	sess,_ := global.Session.SessionStart(w, r)
	//defer sess.SessionRelease(w)
	var t, til, val string
	md5 := tools.Md5(strings.Replace(r.URL.Path[1:],"/",string(0),-1))
	globalDB.QueryRow("SELECT EditTime,Title,Content FROM tb_note_save WHERE Hash=?;",md5).Scan(&t,&til,&val)
	log.Info(md5,t, til)
	au := sess.Get("user")//.(user.User)
	log.Json(au)
	if len(t)!= 0 {
		err := tools.Template(w,"note/content.html",map[string]interface{}{"Uri": r.RequestURI,"Title": til,"Content": template.HTML(val),"Edittime": t[:10],"Level": 0,"User": au})
		log.Info(err)
	}else{
		w.WriteHeader(http.StatusNotFound)
	}
}
