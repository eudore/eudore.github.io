package note;

import (
	"fmt"
//	"bytes"
	"strings"
	"net/http"
//	"strings"
	"io/ioutil" 
//	"html/template"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"public/router"
	"module/global"
//	"module/tools"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)


var (
	stmtQueryPathHash	*sql.Stmt
	stmtQueryNoteData	*sql.Stmt
	stmtUpdatePathHash 	*sql.Stmt
	stmtUpdateNoteData	*sql.Stmt
	stmtUpdateNoteFormat	*sql.Stmt
	)	

var globalDB *sql.DB

func init() {
	// router
	mux := router.New()
	// api
	mux.GetFunc("/api/corn/*",getnote)
	mux.GetFunc("/api/share/*",getnote)
	mux.GetFunc("/api/format/:format",getnote)

	rhash := "/#index^[0-9a-z]{32,32}$/content"
	mux.GetFunc(rhash, getcontent)
	mux.PostFunc(rhash, postcontent)
	mux.PutFunc(rhash, putcontent)
	mux.DeleteFunc(rhash, delcontent)
	mux.DeleteFunc("/api/#index^[0-9a-z]{32,32}$/content", delcontent)
	mux.GetFunc("/#index^[0-9a-z]{32,32}$/index", getindex)
	mux.PostFunc("/#index^[0-9a-z]{32,32}$/index", getindex)
	// note
	mux.GetFunc(":user/*", getnote)
	mux.GetFunc(":user", getnote)
	mux.PostFunc(":user/*", upnote)
	mux.PostFunc(":user", upnote)
	global.Router.SubRoute("/note",mux)
}

func Reload() error {
	globalDB = global.Sql
	stmtQueryPathHash = global.Stmt("SELECT Path,Hash,PHash FROM tb_note_save;")
	stmtQueryNoteData = global.Stmt("SELECT EditTime,Title,Format,Content FROM tb_note_save WHERE Hash=?;")
	stmtUpdatePathHash = global.Stmt("UPDATE tb_note_save SET Hash=?,PHash=? WHERE Hash=?;")
	stmtUpdateNoteData = global.Stmt("UPDATE tb_note_save SET Content=?,Format=? WHERE Hash=?;")
	stmtUpdateNoteFormat = global.Stmt("UPDATE tb_note_save SET Content=?,Format=? WHERE Hash=?;")
	RenewHash()
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


func PPathHash(path string) string{
	n := strings.LastIndex(path,"/")
	if n == -1 {
		n = 0
	}
	return PathHash(path[:n])
}

func PathHash(path string) string{
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(path))
	return hex.EncodeToString(md5Ctx.Sum(nil))
}

func GetAction(r *http.Request) string {
	return ""
}

func GetResource(r *http.Request) string {
	if strings.HasPrefix(r.URL.Path,"/api/") {
		return r.URL.Path[5:]
	}
	return r.URL.Path
}
/*
action:
	readnote 	get
	createnote 	put
	updatenote  post
	deletenote  delete
	sharenote

	updateformat
	formathtml

	flushhash
	flushhome
	cleanfile

role:
	admin: [file:*]
	editer: [file:readfile,file:sharenote,file:createnote,file:updatenote,file:deletenote,file:updateformat]
	reader: [file:readfile,file:sharenote]
*/