package main

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
	"public/middlewares"
)

func init() {
	config.SetDf(global.Config)
	config.ReloadAll()
	bm, err := cache.NewCache("memcache", `{"conn":"127.0.0.1:12001"}`)
	if err == nil {
		bm.Put("weer/public", []byte("file"), 8640000*time.Second)
	}
}

func main() {
	// set default router
	mux := global.Router
	static := http.FileServer(http.Dir("/data/web/static"))
	mux.Handle("/js/", static)
	mux.Handle("/css/", static)
	mux.Handle("/favicon.ico", static)
	mux.HandleFunc("/test/test", test)
	mux.HandleFunc("/test/echo", echo)
	// set reload config function
	config.SetReload("global", 0x100, global.Reload)
	config.SetReload("auth", 0x200, auth.Reload)
	config.SetReload("note", 0x300, note.Reload)
	config.SetReload("file", 0x400, file.Reload)
	config.SetReload("middlewares", 0x500, middlewares.Reload)
	// server set logout reload start
	server.SetOut(log.Info)
	server.SetReload(func() error {
		return config.ReloadAll()
	})
	server.Parse(global.Config.Command, global.Config.Pidfile, func() error {
		// load all config
		config.ReloadAll()
		listen := global.Config.Listen
		if listen.Http2 {
			os.Setenv("LISTEN_HTTP2", "1")
		}
		if listen.Https {
			return server.ListenAndServeTLS(listen.Addr(), listen.Certfile, listen.Keyfile, middlewares.Md)
		} else {
			return server.ListenAndServe(listen.Addr(), middlewares.Md)
		}

	})
}

func test(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:     "test",
		Value:    "true",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Domain:   "wejass.com",
		Expires:  time.Now().Add(10000 * time.Second),
	}
	http.SetCookie(w, cookie)
	w.Write([]byte("Enable test mode !\n"))
}

func echo(w http.ResponseWriter, r *http.Request) {
	log.Info("---echo---")
	log.Info(r.RemoteAddr)
	log.Json(r.URL)
	log.Json(r.Header)
}
