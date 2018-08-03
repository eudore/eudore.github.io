package global


import (
	"encoding/json"
	"database/sql"
	"public/cache"
	"public/session"
	"public/router"
	"public/log"
)

// config
var Config *config;
// Singleton
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
}

func Reload() (err error) {
	// cache
	Cache,err = cache.NewCache(Config.App.Cache.Provider, Config.App.Cache.Config)
	log.Info("Cache: ",err)
	// session
	sessionConfig := &session.ManagerConfig{}
	json.Unmarshal([]byte(Config.App.Session),sessionConfig)
	Session, err = session.NewManager("memcache", sessionConfig)
	log.Info("Session: ",err)
	// sql
	Sql,err = sql.Open("mysql",Config.App.Mysql)
	Sql.SetMaxIdleConns(0)
	log.Info("Sql: ",err)
	return
}


