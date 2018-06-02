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
	_ "module/home"
	_ "module/auth"
	_ "module/note"
	_ "module/file"
	_ "module/chat"
	_ "module/tools"
)

var globalSessions *session.Manager;
var conf *config.Config;

func init() {
	conf = config.Instance()
	bm,err := cache.NewCache("memcache",`{"conn":"127.0.0.1:12001"}`)
	if(err==nil){
		bm.Put("weer/public",[]byte("file"),8640000 * time.Second)
	}
	sessionConfig := &session.ManagerConfig{
		CookieName: "token",
		EnableSetCookie: true,
		Gclifetime: 3600,
		Maxlifetime: 3600,
		Secure: true,
		CookieLifeTime: 3600,
		ProviderConfig: conf.Memcache,
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
	log.Json(conf)
	// router
	mux := router.Instance()
	static := http.FileServer(http.Dir("/data/web/static"))
	mux.Handle("/js/", static)
	mux.Handle("/css/", static)
	mux.Handle("/favicon.ico", static)
	mux.HandleFunc("/", test);
	// set reload
	server.SetReload(config.Reload)
	server.SetOut(log.Info)
	// start
	err := server.Resolve(conf.Command,conf.Pidfile, func() error {
		if conf.Listen.Https {
			return server.ListenAndServeTLS(fmt.Sprintf("%s:%d",conf.Ip,conf.Port),conf.Listen.Certfile,conf.Listen.Keyfile, mux)
		}else {
			return server.ListenAndServe(fmt.Sprintf("%s:%d",conf.Ip,conf.Port), mux)
		}
		
	})
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
