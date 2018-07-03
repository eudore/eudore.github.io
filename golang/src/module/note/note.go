package note;

import (
	"fmt"
//	"bytes"
	"net/http"
//	"strings"
	"io/ioutil" 
//	"html/template"
	"encoding/json"
	"public/router"
	"public/session"
	"module/global"
//	"module/tools"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)


type NoteContent struct {
	Name		string
	Format		string
	Title		string
	EditTime	string
	Content		string
}

var globalSessions *session.Manager;
var globalDB *sql.DB

func init() {
	// router
	mux := router.New()
	rhash := "/#index^[0-9a-z]{32,32}$/content"
	mux.GetFunc(rhash, getcontent)
	mux.PostFunc(rhash, postcontent)
	mux.PutFunc(rhash, putcontent)
	mux.DeleteFunc(rhash, delcontent)
	mux.DeleteFunc("/api/#index^[0-9a-z]{32,32}$/content", delcontent)
	
	mux.GetFunc("/#index^[0-9a-z]{32,32}$/index", getindex)
	mux.PostFunc("/#index^[0-9a-z]{32,32}$/index", getindex)
	mux.GetFunc(":user/*", note)
	mux.GetFunc(":user", note)
	mux.GetFunc("content/", note)
	router.Instance().SubRoute("/note",mux)
}

func Reload() error {
	globalDB = global.Sql
	globalSessions = global.Session
	return nil
}

func postcontent(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()  
	data, _ := ioutil.ReadAll(r.Body) //获取post的数据  
	index := router.GetValue(r, "index")
	if stmt, err := globalDB.Prepare("INSERT tb_note_save(Content,Hash) VALUES(?,?);");err==nil{
		_, err = stmt.Exec(data, index)	
		responseBody,_ := json.Marshal(map[string]interface{}{"result": err==nil})
		w.Write(responseBody)
		return
	}
	//w.Write([]byte(fmt.Sprintf("{\"result\":%t}",false)))
}

func putcontent(w http.ResponseWriter, r *http.Request) {
	//r.ParseForm() 
	//r.PostFormValue("data")//获取post的数据
	defer r.Body.Close()
	data, _ := ioutil.ReadAll(r.Body)
	index := router.GetValue(r, "index")
	if stmt, err := globalDB.Prepare("UPDATE tb_note_save SET Content=? WHERE Hash=?;");err==nil{
		_, err = stmt.Exec(data, index)
		responseBody,_ := json.Marshal(map[string]interface{}{"result": err==nil})
		w.Write(responseBody)
	}
}

func delcontent(w http.ResponseWriter, r *http.Request) {
	index := router.GetValue(r, "index")
	if stmt, err := globalDB.Prepare("DELETE FROM tb_note_save WHERE Hash=?;");err==nil{
		_, err = stmt.Exec(index)
		responseBody,_ := json.Marshal(map[string]interface{}{"result": err==nil})
		w.Write(responseBody)
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
