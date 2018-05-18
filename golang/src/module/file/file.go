//file.go
package file;

import (
	"net/http"
	"encoding/json"
	"public/config"
//	"public/log"
	"public/router"
	"public/session"
	"module/file/filestore"
	_ "module/file/filestore/disk"
	_ "module/file/filestore/oss"
)


var globalSessions *session.Manager;
var conf *config.Config;

func init() {
	filestore.Reload()
	conf = config.Instance()
	sessionConfig := &session.ManagerConfig{}
	//sessionConfig := &session.ManagerConfig{CookieName: "token",EnableSetCookie: true, Gclifetime: 3600, Maxlifetime: 3600, Secure: true, CookieLifeTime: 3600, ProviderConfig: "127.0.0.1:12001"}
	json.Unmarshal([]byte(conf.Session),sessionConfig)
	globalSessions, _ = session.NewManager("memcache", sessionConfig)
	go globalSessions.GC()
	mux := router.Instance()
	mux.GetFunc("/file/:user/:zone/*",fileget)
	mux.PostFunc("/file/:user/:zone/*",fileup)
//	mux.PutFunc("/file/:user/:zone/*",file_up)
	mux.GetFunc("/file/file",file)
}


func file(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("file:"+r.URL.Path))
}