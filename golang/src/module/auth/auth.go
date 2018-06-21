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
	"module/auth/oauth2"
	"module/auth/user"
)  

var conf *config.Config;
var globalSessions *session.Manager;

func init() {
	conf = config.Instance()
	sessionConfig := &session.ManagerConfig{}
	json.Unmarshal([]byte(conf.App.Session),sessionConfig)
	globalSessions, _ = session.NewManager("memcache", sessionConfig)
	go globalSessions.GC()
	
	mux := router.New()
	rlogin := "/login"
	mux.GetFunc(rlogin, logshow)
	mux.PostFunc(rlogin, login)
	mux.DeleteFunc(rlogin, logout)
	mux.PostFunc("/user/:name", usernew)
	mux.GetFunc("/projetc/new", projectnew)
	mux.PostFunc("/projetc/:name", projectnew)
	mux.PostFunc("/projetc/:name/snippets", projectnew)
	mux.PostFunc("/projetc/:name/members", projectnew)
	mux.PostFunc("/projetc/:name/members/:user", projectnew)
	mux.PostFunc("/manage/:name", projectnew)
	mux.PostFunc("/share/:name", projectnew)
	mux.GetFunc("/", auth)

	mux.GetFunc("/user/auth",user.Auth)
	mux.PostFunc("/user/auth",user.Authpass)
	mux.GetFunc("/user/signup",user.Signup)
	mux.PostFunc("/user/signup",user.SignupSubmit)
	mux.GetFunc("/user/login",user.Login)
	mux.GetFunc("/user/logout",user.Logout)
	r1,r2 := oauth2.GetRouter()
	mux.SubRoute("/oauth2/login",r1)
	mux.SubRoute("/oauth2/callback",r2)

	router.Instance().SubRoute("/auth",mux)
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
