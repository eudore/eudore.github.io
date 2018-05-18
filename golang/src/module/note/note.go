package note;

import (
	"fmt"
	"bytes"
	"net/http"
	"strings"
	"io/ioutil" 
	"html/template"
	"encoding/json"
	"public/config"
	"public/router"
	"public/session"
	"module/tools"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var conf *config.Config;
var globalSessions *session.Manager;
var dbconfig string = "root:@/Jass"

func init() {
	conf = config.Instance()
	sessionConfig := &session.ManagerConfig{}
	json.Unmarshal([]byte(conf.Session),sessionConfig)
	globalSessions, _ = session.NewManager("memcache", sessionConfig)
	go globalSessions.GC()
	mux := router.Instance()
	rhash := "/note/#index^[0-9a-z]{32,32}$/content"
	mux.GetFunc(rhash, getcontent)
	mux.PostFunc(rhash, postcontent)
	mux.PutFunc(rhash, putcontent)
	mux.DeleteFunc(rhash, delcontent)
	mux.DeleteFunc("/api/note/#index^[0-9a-z]{32,32}$/content", delcontent)
	
	mux.GetFunc("/note/#index^[0-9a-z]{32,32}$/index", getindex)
	mux.PostFunc("/note/#index^[0-9a-z]{32,32}$/index", getindex)
	mux.GetFunc("/note/:user/*", note)
	mux.GetFunc("/note/:user", note)
	mux.GetFunc("/note/content/", note)
}

func postcontent(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()  
	data, _ := ioutil.ReadAll(r.Body) //获取post的数据  
	index := router.GetValue(r, "index")
	if db, err := sql.Open("mysql",dbconfig);err==nil {
		defer db.Close()
		if stmt, err := db.Prepare("INSERT tb_note_save(Content,Hash) VALUES(?,?);");err==nil{
			_, err = stmt.Exec(data, index)	
			responseBody,_ := json.Marshal(map[string]interface{}{"result": err==nil})
			w.Write(responseBody)
			return
		}
	}	
	//w.Write([]byte(fmt.Sprintf("{\"result\":%t}",false)))
}

func putcontent(w http.ResponseWriter, r *http.Request) {
	//r.ParseForm() 
	//r.PostFormValue("data")//获取post的数据
	defer r.Body.Close()
	data, _ := ioutil.ReadAll(r.Body)
	index := router.GetValue(r, "index")
	if db, err := sql.Open("mysql",dbconfig);err==nil {
		defer db.Close()
		if stmt, err := db.Prepare("UPDATE tb_note_save SET Content=? WHERE Hash=?;");err==nil{
			_, err = stmt.Exec(data, index)
			responseBody,_ := json.Marshal(map[string]interface{}{"result": err==nil})
			w.Write(responseBody)
		}
	}	
}

func delcontent(w http.ResponseWriter, r *http.Request) {
	index := router.GetValue(r, "index")
	if db, err := sql.Open("mysql",dbconfig);err==nil {
		defer db.Close()
		if stmt, err := db.Prepare("DELETE FROM tb_note_save WHERE Hash=?;");err==nil{
			_, err = stmt.Exec(index)
			responseBody,_ := json.Marshal(map[string]interface{}{"result": err==nil})
			w.Write(responseBody)
		}
	}	
}



func note(w http.ResponseWriter, r *http.Request) {
	md5 :=tools.Md5(strings.Replace(r.URL.Path[6:],"/",string(0),-1))
	var t string
	var val string
	//user := router.GetValue(r, "user")
	//index := router.GetValue(r, "index")
	if db, err := sql.Open("mysql",dbconfig);err==nil {
		defer db.Close()
		db.QueryRow("SELECT EditTime,Content FROM tb_note_save WHERE Hash=?;",md5).Scan(&t,&val)
	}	
	var doc bytes.Buffer
	tmp, err := template.ParseFiles("/data/web/templates/note/content.html")
	if err == nil {
		tmp.Execute(&doc,map[string]interface{}{"Content": template.HTML(val),"Edittime": t})
		w.Write([]byte(doc.String()))
	}
}




func user(w http.ResponseWriter, r *http.Request) {
	val := router.GetValue(r, "user")
	w.Write([]byte(fmt.Sprintf("note user: %s", val)))
}


func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
