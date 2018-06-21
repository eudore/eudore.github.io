package auth;

import (  
    "fmt" 
    "net/http"
	"public/router"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)  


func usernew(w http.ResponseWriter, r *http.Request) {
	sess,_ := globalSessions.SessionStart(w, r)
	defer sess.SessionRelease(w)
	name := router.GetValue(r, "name")
	if db, err := sql.Open("mysql","root:@/Jass");err==nil {
		defer db.Close()
		r.ParseForm()
		user := r.PostFormValue("user")
		pass := r.PostFormValue("pass")
		if stmt, err := db.Prepare("INSERT `tb_auth_user_info`(Name,User,Pass) VALUES(?,?,?);");err==nil{
			_, err = stmt.Exec(name,user,pass)	
			w.Write([]byte(fmt.Sprintf("{\"result\":%t}",err==nil)))
			return
		}
	}	
	w.Write([]byte(fmt.Sprintf("{\"result\":%t}",false)))
}
