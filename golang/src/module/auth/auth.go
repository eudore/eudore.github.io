package auth;

import ( 
	"bytes"
	"net/http"
	"html/template"
	"public/router"
	"public/session"
	"public/log"
)  

var globalSessions *session.Manager;

func init() {
	sessionConfig := &session.ManagerConfig{CookieName: "token",EnableSetCookie: true, Gclifetime: 3600, Maxlifetime: 3600, Secure: true, CookieLifeTime: 3600, ProviderConfig: "127.0.0.1:12001"}
	globalSessions, _ = session.NewManager("memcache", sessionConfig)
	go globalSessions.GC()
	mux := router.Instance()
	mux.GetFunc("/auth/login", logshow)
	mux.PostFunc("/auth/login", login)
	mux.DeleteFunc("/auth/login", logout)
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
