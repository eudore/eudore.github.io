package auth;

import (  
	"fmt"
    "net/http"
	"github.com/golang/glog"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
) 

func logshow(w http.ResponseWriter, r *http.Request) {
	sess,_ := globalSessions.SessionStart(w, r)
	defer sess.SessionRelease(w)
	name := sess.Get("name")
	glog.Info("show login to ",name)
	defer glog.Flush()
	if name == nil{
		http.ServeFile(w, r,"/data/web/templates/auth/login.html")
	}else {
		http.Redirect(w,r,fmt.Sprintf("/auth/home.html?name=%s",name),http.StatusFound)	
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	sess,_ := globalSessions.SessionStart(w, r)
	defer sess.SessionRelease(w)
	if db, err := sql.Open("mysql","root:@/Jass");err==nil {
		defer db.Close()
		defer glog.Flush()
		r.ParseForm()
		user := r.PostFormValue("user")
		pass := r.PostFormValue("pass")
		var uid int
		var name string
		var level int
		db.QueryRow("CALL pro_auth_login(?,?,?,@out_uid,@out_name,@out_level)",user,pass,0).Scan(&uid,&name,&level)
		if uid > 0 {
			sess.Set("uid",uid)
			sess.Set("name",name)
			sess.Set("level",level)
			glog.Info("login name: ",name," uid:",uid)	
			w.Write([]byte(fmt.Sprintf("{\"result\":\"%d\",\"url\":\"/auth/home.html?name=%s\"}", level,name)))
			rows, err := db.Query("SELECT CONCAT_WS('/',`UName`,`PName`) `Path`,`Level` FROM v_auth_authorized WHERE UID=?;",uid)
			if err == nil {
				glog.Info("auth to: ",name)	
				auth := make(map[string]int)
				for rows.Next(){
					rows.Scan(&name,&level)
					auth[name]=level
					glog.Info("auth path: ",name," level: ",level)
				}
				sess.Set("auth",auth)
			}
			return
		}else{
			http.Error(w,"erre",http.StatusUnauthorized)
		}
	}
	w.Write([]byte(fmt.Sprintf("{\"result\":\"%d\"}", -1)))
}

func logout(w http.ResponseWriter, r *http.Request) {
	sess,_ := globalSessions.SessionStart(w, r)
	defer sess.SessionRelease(w)
	sess.SessionRelease(w)
	http.RedirectHandler("/auth/",http.StatusFound)	
}

