package global


import (
	"encoding/json"
	"public/config"
	"public/cache"
	"public/session"
	"public/router"
	"public/log"
	"database/sql"
)

var Config *config.Config;
var Cache cache.Cache
var Session *session.Manager;
var Router *router.Mux
var Sql *sql.DB

func init(){
	var err error
	// config
	Config = config.Instance()
	// cache
	//Cache,_ = cache.NewCache(Config.App.Cache)
	// session
	sessionConfig := &session.ManagerConfig{}
	json.Unmarshal([]byte(Config.App.Session),sessionConfig)
	Session, err = session.NewManager("memcache", sessionConfig)
	log.Info("Session: ",err)
	// router
	Router = router.Instance()
	// sql
	Sql,err = sql.Open("mysql",Config.App.Mysql)
	log.Info("Sql: ",err)
}