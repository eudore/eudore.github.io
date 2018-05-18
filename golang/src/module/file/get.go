package file;


import (
	"fmt"
	"net/http"
	"html/template"
//	"public/log"
	"public/router"
	"module/file/filestore"
)


func fileget(w http.ResponseWriter, r *http.Request) {
	fmt.Println(filestore.PathHash(r.RequestURI[6:]))
	if filestore.IsFile(r.RequestURI[6:]) {
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
	fs := filestore.List(r.RequestURI)
	//t, _ := template.ParseFiles("/data/web/templates/file/file.html");
	tmp, _ := template.ParseFiles("/data/web/templates/file/file.html","/data/web/templates/base.html")
	tmp.Execute(w, map[string]interface{}{"url": r.RequestURI,"files": fs});
}