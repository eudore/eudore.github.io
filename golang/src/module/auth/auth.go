package auth;

import ( 
	"bytes"
	"net/http"
	"encoding/json"
	"html/template"
	"public/config"
	"public/router"
	"public/session"
	"public/log"
)  

var conf *config.Config;
var globalSessions *session.Manager;

func init() {
	conf = config.Instance()
	sessionConfig := &session.ManagerConfig{}
	json.Unmarshal([]byte(conf.Session),sessionConfig)
	globalSessions, _ = session.NewManager("memcache", sessionConfig)
	go globalSessions.GC()
	mux := router.Instance()
	rlogin := "/auth/login"
	mux.GetFunc(rlogin, logshow)
	mux.PostFunc(rlogin, login)
	mux.DeleteFunc(rlogin, logout)
	mux.PostFunc("/auth/user/:name", usernew)
	mux.GetFunc("/auth/projetc/new", projectnew)
	mux.PostFunc("/auth/projetc/:name", projectnew)
	mux.PostFunc("/auth/projetc/:name/snippets", projectnew)
	mux.PostFunc("/auth/projetc/:name/members", projectnew)
	mux.PostFunc("/auth/projetc/:name/members/:user", projectnew)
	mux.PostFunc("/auth/manage/:name", projectnew)
	mux.PostFunc("/auth/share/:name", projectnew)
	mux.GetFunc("/auth/", auth)
}

func auth(w http.ResponseWriter, r *http.Request) {
	sess,_ := globalSessions.SessionStart(w, r)
	defer sess.SessionRelease(w)
	log.Info("pro to: ",sess.Get("auth") )
	tmp, err := template.ParseFiles("/data/web/templates/auth/auth.html","/data/web/templates/base.html")
	if err == nil {	
		var doc bytes.Buffer
		var data map[string]int
		if value, ok := sess.Get("auth").(map[string]int); ok {
			data = value
		}
		tmp.Execute(&doc,map[string]interface{}{"data": data})
		w.Write([]byte(doc.String()))
	}
}
