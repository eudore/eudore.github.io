package file;

import (
	//"fmt"
	"time"
    "strings"
    "io/ioutil"
    "net/http"
    "html/template"
    "public/log"
)

type FileInfo struct{
	Name 	string
	Size 	int64
	Dir 	bool
	ModTime	time.Time
}
func file_list(w http.ResponseWriter, r *http.Request) {
    log.Info(r.RequestURI)
    //解析模板文件
    dir := "/data/web/upload/"+strings.SplitN(r.URL.Path,"/",3)[2]
    files,err := ioutil.ReadDir(dir)
    var fs []FileInfo = make([]FileInfo,len(files))
	for i, fi := range files {  
		fs[i]=FileInfo{	Name:	fi.Name(),Size: fi.Size(),Dir:fi.IsDir(),	ModTime: fi.ModTime()	}
    }  
    t, err := template.ParseFiles("/data/web/templates/file/file.html");
    log.Info(err)
    t.Execute(w, map[string]interface{}{"files": fs});
}