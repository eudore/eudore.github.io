package main;

import (
	"os"
	"fmt"
	"time"
	"net/http"
	_ "github.com/go-sql-driver/mysql"
	_ "public/cache/memcache"
	_ "public/session/memcache"
	"public/config"
	"public/log"
	"public/server"
	"public/session"
	"public/cache"
	"public/router"
	"module/global"
	// _ "module/home"
	"module/auth"
	"module/note"
	// "module/file"
	// _ "module/chat"
	// _ "module/tools"
)

var globalSessions *session.Manager;
var conf *config.Config;

func init() {
	conf = config.Instance()
	config.Reload("config")
	bm,err := cache.NewCache("memcache",`{"conn":"127.0.0.1:12001"}`)
	if(err==nil){
		bm.Put("weer/public",[]byte("file"),8640000 * time.Second)
	}
	sessionConfig := &session.ManagerConfig{
		CookieName: "s",
		EnableSetCookie: true,
		Gclifetime: 3600,
		Maxlifetime: 3600,
		Secure: true,
		CookieLifeTime: 3600,
		ProviderConfig: global.App.Memcache,
	}
	globalSessions,_ = session.NewManager("memcache", sessionConfig)
	go globalSessions.GC()
}



func test(w http.ResponseWriter, r *http.Request) {
	sess,_ := globalSessions.SessionStart(w, r)
	defer sess.SessionRelease(w)
	username := sess.Get("username")
	fmt.Println( r.URL.Path," ",username)
}

func main() {
	// set default router
	mux := router.Instance()
	static := http.FileServer(http.Dir("/data/web/static"))
	mux.Handle("/js/", static)
	mux.Handle("/css/", static)
	mux.Handle("/favicon.ico", static)
	mux.HandleFunc("/", test);
	// set reload config function
	config.SetReload("auth", auth.Reload)
	config.SetReload("note", note.Reload)
	// config.SetReload("file", file.Reload)
	// server set logout reload start
	server.SetOut(log.Info)
	server.SetReload(func() error {
		return config.Reload()
	})
	server.Parse(conf.Command,conf.Pidfile, func() error {
		// load all config
		config.Reload()
		if global.Listen.Html2 {
			os.Setenv("LISTEN_HTML2","1")
		}
		if global.Listen.Https {
			return server.ListenAndServeTLS(global.Listen.Addr(),global.Listen.Certfile,global.Listen.Keyfile, mux)
		}else {
			return server.ListenAndServe(global.Listen.Addr(), mux)
		}
		
	})
}
