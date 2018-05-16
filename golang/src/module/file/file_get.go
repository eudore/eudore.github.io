package file;


import (
	"fmt"
	"net/http"
	"html/template"
//	"public/log"
	"public/router"
	"module/file/filestore"
)


func file_get(w http.ResponseWriter, r *http.Request) {
	
}

func file_list(w http.ResponseWriter, r *http.Request) {
	fmt.Println(filestore.PathHash(r.URL.Path[6:]))
	if filestore.IsFile(r.URL.Path[6:]) {
		user := router.GetValue(r, "user")
		zone := router.GetValue(r, "zone")
		fs,err := filestore.Getstore(user+"/"+zone)
		if err!=nil {
			w.WriteHeader(http.StatusNotFound) 
			return
		}
		fs.Load(w,r)
		return
	}
	fs := filestore.List(r.URL.Path)
	t, _ := template.ParseFiles("/data/web/templates/file/file.html");
	t.Execute(w, map[string]interface{}{"url": r.URL.Path,"files": fs});
}