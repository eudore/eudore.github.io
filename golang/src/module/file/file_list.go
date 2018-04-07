package file;

import (
	"fmt"
    "strings"
    "io/ioutil"
    "net/http"
    "html/template"
    "public/log"
    //"database/sql"
    _ "github.com/go-sql-driver/mysql"
    //"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type FileInfo struct{
	Name 	string
	Size 	string
	Dir 	bool
	ModTime	string
}
func local_list(w http.ResponseWriter, r *http.Request) {
    log.Info(r.RequestURI)
    //解析模板文件
    dir := "/data/web/upload/"+strings.SplitN(r.URL.Path,"/",3)[2]
    files,_ := ioutil.ReadDir(dir)
    var fs []FileInfo = make([]FileInfo,len(files))
	for i, fi := range files {  
        fs[i]=FileInfo{ Name:   fi.Name(),Size: get_size(fi.Size()),Dir:fi.IsDir(),   ModTime: fi.ModTime().Format("2006-01-02 15:04")   }
        if fi.IsDir(){
            fs[i].Size="-"
        }
    }  
    t, _ := template.ParseFiles("/data/web/templates/file/file.html");
    t.Execute(w, map[string]interface{}{"files": fs});
}

func file_list(w http.ResponseWriter, r *http.Request) []FileInfo{
/*    path := strings.SplitN(r.RequestURI,"/",3)[2]
    rows, err := db.Query("SELECT Name,Size,ModTime FROM tb_file_save WHERE Hash=?;",path)
    if err == nil {
        log.Info("auth to: ",name)  
        auth := make(map[string]int)
        for rows.Next(){
            rows.Scan(&name,&level)
            auth[name]=level
            log.Info("auth path: ",name," level: ",level)
        }
    }*/
    return nil
}


func oss_list(w http.ResponseWriter, r *http.Request) {
    
}


func get_size(file_bytes int64) string {
    var i     int
    var units = [6]string{"B", "K", "M", "G", "T", "P"}
    i = 0
    for {
        if file_bytes < 1024 {
            return fmt.Sprintf("%d", file_bytes) + units[i]
        }
        file_bytes = file_bytes >> 10
        i++
    }
}