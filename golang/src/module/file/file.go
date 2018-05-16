//file.go
package file;

import (
	"net/http"
	"public/router"
	"public/session"
	"module/file/filestore"
	_ "module/file/filestore/disk"
	_ "module/file/filestore/oss"
)


var globalSessions *session.Manager;

func init() {
	filestore.Reload()
	sessionConfig := &session.ManagerConfig{CookieName: "token",EnableSetCookie: true, Gclifetime: 3600, Maxlifetime: 3600, Secure: true, CookieLifeTime: 3600, ProviderConfig: "127.0.0.1:12001"}
	globalSessions, _ = session.NewManager("memcache", sessionConfig)
	go globalSessions.GC()
	mux := router.Instance()
	mux.GetFunc("/file/:user/:zone/*",file_list)
	mux.PostFunc("/file/:user/:zone/*",file_up)
//	mux.PutFunc("/file/:user/:zone/*",file_up)
	mux.GetFunc("/file/file",file)
}


func file(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("file"))
}