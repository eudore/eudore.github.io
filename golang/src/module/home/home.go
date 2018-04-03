package note;

import (
//	"fmt"
	"bytes"
//	"reflect"
	"net/http"
//	"strings"
//	"io/ioutil" 
//	"encoding/json"
	"html/template"
	"public/router"
	"public/session"
	"public/cache"
//	"module/tools"
//	"github.com/golang/glog"
//	"database/sql"
//	_ "github.com/go-sql-driver/mysql"
)

var globalSessions *session.Manager;
var globalCache cache.Cache;

func init() {
	sessionConfig := &session.ManagerConfig{CookieName: "gosessionid",Gclifetime: 3600,ProviderConfig: "127.0.0.1:12001"}
	globalSessions, _ = session.NewManager("memcache", sessionConfig)
	go globalSessions.GC()
	//var err error;
	globalCache,_ = cache.NewCache("memcache",`{"conn":"127.0.0.1:12001"}`)
	mux := router.Instance()
	mux.GetFunc("/home/:name", home)
	mux.GetFunc("/home/:name/:zone", home)
}



func home(w http.ResponseWriter, r *http.Request) {
	sess,_ := globalSessions.SessionStart(w, r)
	defer sess.SessionRelease(w)
	//glog.Info("pro to: ",sess.Get("auth") )

	name := router.GetValue(r, "name")
	zone := router.GetValue(r, "zone")
	tmp, err := template.ParseFiles("/data/web/templates/auth/auth.html","/data/web/templates/base.html")
	if err == nil {	
		var doc bytes.Buffer
		var data map[string]int
		if value, ok := sess.Get("auth").(map[string]int); ok {
			data = value
		}
		tmp.Execute(&doc,map[string]interface{}{"data": data,"zone": globalCache.Get(name+zone)})
		w.Write([]byte(doc.String()))
	}
}