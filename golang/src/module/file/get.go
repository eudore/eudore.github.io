package file;


import (
	"fmt"
	"strconv"
	"net/http"
	"html/template"
	"public/log"
	"public/router"
	"module/file/store"
)


func fileget(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.RequestURI,r.URL.Path,store.PathHash(r.URL.Path))
	user := router.GetValue(r, "user")
	zone := router.GetValue(r, "zone")
	fs,err := store.Getstore(user+"/"+zone)
	//log.Json(fs)
	if err!=nil {
		w.WriteHeader(http.StatusNotFound) 
		return
	}
	if r.URL.RawQuery == "signed" {
		length,_ := strconv.ParseInt(r.Header.Get("length"), 10, 64)
		p := &store.Policy{
			Host: "https://" + r.Host,
			Directory: Prefix + r.URL.Path,
			Method: r.Header.Get("method"),
			Length: length,
		}
		response := fs.Signed(p)
		log.Json(p)
		w.Header().Set("Access-Control-Allow-Methods", r.Header.Get("method"))
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
		return
	}
	if store.IsFile(r.URL.Path) {
		if err!=nil {
			w.WriteHeader(http.StatusNotFound) 
			return
		}
		fmt.Println("file load file:",r.URL.Path,fs.Load(w,r.URL.Path))
		return
	}
	f := store.List(r.URL.Path)
	tmp, _ := template.ParseFiles("/data/web/templates/file/file.html","/data/web/templates/base.html")
	tmp.Execute(w, map[string]interface{}{"url": r.URL.Path,"files": f});
	w.WriteHeader(http.StatusOK)
}