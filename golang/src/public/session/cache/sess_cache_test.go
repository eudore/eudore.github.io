package cache

import (
	"testing"
	"fmt"
	"time"
	"net/http"
	"public/cache"
	"public/session"
	_ "public/cache/memcache"
	_ "public/cache/sync"
	_ "public/session/cache"
	_ "public/session/memcache"
)

var bm cache.Cache
var globalSessions *session.Manager;

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("helloHandler:"+r.URL.Path)
	w.Write([]byte("hello: "+req.URL.Path))
}

func setHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(bm.Put("astax", []byte("22"), 10 * time.Second))
	sess, _ := globalSessions.SessionStart(w, r)
	defer sess.SessionRelease(w)
	sess.Set("name", "sssshuve")
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(string(bm.Get("astax")))
	sess, _ := globalSessions.SessionStart(w, r)
	defer sess.SessionRelease(w)
	fmt.Println(sess.Get("name"))
}

func TestCache(t *testing.T) {
	ch, err := cache.NewCache("memcache",`{"conn":"127.0.0.1:12003"}`)	
	bm=ch
	t.Log(err)
	sessionConfig := &session.ManagerConfig{CookieName: "token",EnableSetCookie: true,
		Gclifetime: 3600, Maxlifetime: 3600, Secure: true, CookieLifeTime: 3600,
		ProviderConfig: `{"cache": "memcache","conn":"127.0.0.1:12003"}`}
	globalSessions, err = session.NewManager("cache", sessionConfig)
	go globalSessions.GC()
	t.Log(err)
	http.HandleFunc("/go", helloHandler)
	http.HandleFunc("/set", setHandler)
	http.HandleFunc("/get", getHandler)
	http.ListenAndServe(":8084", nil)
}
