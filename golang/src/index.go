package main;

import (
	"fmt"
	"log"
	"net/http"
	_ "public/daemon"
	"public/runtime"
	_ "public/cache/memcache"
	_ "public/session/memcache"
	"public/session"
	"public/router"
	_ "public/init"
	_ "module/home"
	_ "module/auth"
	_ "module/note"
	_ "module/file"
	_ "module/chat"
	_ "module/tools"
)

var globalSessions *session.Manager;

func init() {
	sessionConfig := &session.ManagerConfig{CookieName: "token",EnableSetCookie: true, Gclifetime: 3600, Maxlifetime: 3600, Secure: true, CookieLifeTime: 3600, ProviderConfig: "127.0.0.1:12001"}
	globalSessions, _ = session.NewManager("memcache", sessionConfig)
	go globalSessions.GC()
}



func test(w http.ResponseWriter, r *http.Request) {
	sess,_ := globalSessions.SessionStart(w, r)
	defer sess.SessionRelease(w)
	username := sess.Get("username")
	fmt.Println( r.URL.Path," ",username)
}

func main() {
	//http.HandleFunc("/note", note.Handle);
	//http.HandleFunc("/auth/", auth.Handle);
	//http.HandleFunc("/file/upload", file.Handle);
	//http.Handle("/chat/wss", websocket.Handler(chat.Handle));
	mux := router.Instance()
	mux.HandleFunc("/go/", test);
	mux.Handle("/",http.HandlerFunc(test));
	err := runhttp.ListenAndServe(":8080", mux);
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	select{};//阻塞进程
}
