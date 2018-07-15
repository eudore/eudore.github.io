package global


import (
	"net/http"
	"encoding/json"
	"database/sql"
	"github.com/NYTimes/gziphandler"
	"public/cache"
	"public/session"
	"public/router"
	"public/log"
)

// config
var Config *config;

// Singleton
var Gw http.Handler 
var Cache cache.Cache
var Session *session.Manager;
var Router *router.Mux
var Sql *sql.DB


func init(){
	Config = &config{
		Config:		"/data/web/config/conf.json",
	}
	// router
	Router = router.New()
	Gw = gziphandler.GzipHandler(&gw{})
}

func Reload() (err error) {
	// cache
	//Cache,_ = cache.NewCache(Config.App.Cache)
	// session
	sessionConfig := &session.ManagerConfig{}
	json.Unmarshal([]byte(Config.App.Session),sessionConfig)
	Session, err = session.NewManager("memcache", sessionConfig)
	log.Info("Session: ",err)
	// sql
	Sql,err = sql.Open("mysql",Config.App.Mysql)
	log.Info("Sql: ",err)
	return
}


