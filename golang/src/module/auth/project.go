package auth;

import (  
    "fmt" 
    "net/http"
	"public/router"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)  

func projectnew(w http.ResponseWriter, r *http.Request) {
	sess,_ := globalSessions.SessionStart(w, r)
	defer sess.SessionRelease(w)
	uid := sess.Get("uid")
	pname := router.GetValue(r, "name")
	if db, err := sql.Open("mysql","root:@/Jass");err==nil {
		defer db.Close()
		r.ParseForm()
		rang := r.PostFormValue("range")
		if stmt, err := db.Prepare("INSERT `tb_auth_project`(Name,UID,Range) VALUES(?,?,?);");err==nil{
			_, err = stmt.Exec(pname,uid,rang)	
			w.Write([]byte(fmt.Sprintf("{\"result\":%t}",err==nil)))
			return
		}
	}	
	w.Write([]byte(fmt.Sprintf("{\"result\":%t}",false)))
}