//file.go
package file;

import (
	"net/http"
	// "encoding/json"
	// "public/config"
	// "public/log"
	"public/router"
	"public/session"
	"module/global"
	"module/file/store"
	_ "module/file/disk"
	_ "module/file/oss"
)


var globalSessions *session.Manager;

const (
	Prefix		=	"/file"
)

func init() {
	mux := router.New()
	mux.GetFunc("/:user/:zone/*",fileget)
	mux.PostFunc("/:user/:zone/*",fileup)
//	mux.PutFunc("/:user/:zone/*",file_up)
	mux.GetFunc("/file",file)

	global.Router.SubRoute(Prefix,mux)
}


func Reload() error {
	globalSessions = global.Session
	return store.Reload()
}

func file(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("file:"+r.URL.Path))
}
