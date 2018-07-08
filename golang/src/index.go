package main;

import (
	"os"
	"time"
	"net/http"
	_ "github.com/go-sql-driver/mysql"
	_ "public/cache/memcache"
	_ "public/session/memcache"
	"public/config"
	"public/log"
	"public/server"
	"public/cache"
	"module/global"
	// "module/tools"
	"module/auth"
	"module/note"
	"module/file"
	// _ "module/chat"
	// _ "module/home"
	// _ "module/tools"
)

func init() {
	config.Reload()
	bm,err := cache.NewCache("memcache",`{"conn":"127.0.0.1:12001"}`)
	if(err==nil){
		bm.Put("weer/public",[]byte("file"),8640000 * time.Second)
	}
}




func main() {
	// set default router
	mux := global.Router
	static := http.FileServer(http.Dir("/data/web/static"))
	mux.Handle("/js/", static)
	mux.Handle("/css/", static)
	mux.Handle("/favicon.ico", static)
	mux.HandleFunc("/test", test)
	// set reload config function
	config.SetReload("global", 100, global.Reload)
	config.SetReload("auth", 200, auth.Reload)
	config.SetReload("note", 300, note.Reload)
	config.SetReload("file", 400, file.Reload)
	// server set logout reload start
	server.SetOut(log.Info)
	server.SetReload(func() error {
		return config.Reload()
	})
	conf := config.Instance()
	server.Parse(conf.Command,conf.Pidfile, func() error {
		// load all config
		config.Reload()
		if global.Listen.Html2 {
			os.Setenv("LISTEN_HTML2","1")
		}
		if global.Listen.Https {
			return server.ListenAndServeTLS(global.Listen.Addr(), global.Listen.Certfile, global.Listen.Keyfile, global.Gw)
		}else {
			return server.ListenAndServe(global.Listen.Addr(), global.Gw)
		}
		
	})
}


func test(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:		"test",
		Value:		"true",
		Path:		"/",
		HttpOnly:	true,
		Secure:		true,
		Domain:		"wejass.com",
		Expires:	time.Now().Add(10000 * time.Second),
	}
	http.SetCookie(w,cookie)
	log.Info("---echo---")
	log.Info(r.RemoteAddr)
	log.Json(r.URL)
	log.Json(r.Header)
}
