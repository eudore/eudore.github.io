package file;

import (
	"strconv"
	"net/http"
	"encoding/json"
	"public/router"
	"public/log"
	"module/file/filestore"
	"database/sql"
)
func file_up(w http.ResponseWriter, r *http.Request) {
	user := router.GetValue(r, "user")
	zone := router.GetValue(r, "zone")
	log.Info(r.URL.Path)
	fs,err := filestore.Getstore(user+"/"+zone)
	if err!=nil {
		w.WriteHeader(http.StatusNotFound) 
		return
	}
	if r.URL.RawQuery == "" {
		var p filestore.PolicyInfo
		p.Directory = r.URL.Path
		p.Method = "POST"
		p.Length = 0
		response := fs.Policy(&p)
		log.Json(p)
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
		// if(true || strings.HasPrefix(r.Header.Get("Referer"),"https://www.wejass.com/file/")){
		// }else{
		// 	w.WriteHeader(http.StatusNotFound)                  
		// }
	}else{
		fs,err := fs.Save(r)
		if err==nil{
			filestore.Add(user+"/"+zone,fs)
			responseBody,_ := json.Marshal(map[string]interface{}{"status":"ok","data":fs})
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Content-Length", strconv.Itoa(len(responseBody)))
			w.WriteHeader(http.StatusOK)
			w.Write(responseBody)
			log.Info("Post Response : 200 OK . uri: ",fs)   
		}
	}
}


func file_insert(path ,source string) error {
	if db, err := sql.Open("mysql","root:@/Jass");err==nil {
		defer db.Close()
		if stmt, err := db.Prepare("INSERT tb_note_save(Content,Hash) VALUES(?,?);");err==nil{
			_, err = stmt.Exec(path, source) 
			//w.Write([]byte(fmt.Sprintf("{\"result\":%t}",err==nil)))
		}
	}   
	return nil
}