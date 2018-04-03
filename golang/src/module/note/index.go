package note;

import (
	"fmt"
	"net/http"
	"encoding/json"
	"public/router"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func getindex(w http.ResponseWriter, r *http.Request) {
	val := router.GetValue(r, "index")
	w.Write([]byte(fmt.Sprintf("note content: %s", val)))
}

func getcontent(w http.ResponseWriter, r *http.Request) {
	var t float64
	var val string
	index := router.GetValue(r, "index")
	if db, err := sql.Open("mysql","root:@/Jass");err==nil {
		defer db.Close()
		db.QueryRow("SELECT UNIX_TIMESTAMP(EditTime),Content FROM tb_note_save WHERE Hash=?;",index).Scan(&t,&val)
	}	
	b, _ := json.Marshal(map[string]interface{}{"result": (val!=""),"edittime": t,"content":val } )
	w.Write([]byte(b))	
}
