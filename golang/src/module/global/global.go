package global


import (
	"fmt"
	"encoding/json"
	"public/config"
	"public/cache"
	"public/session"
	"public/router"
	"public/log"
	"database/sql"
)

// config
var Config *config.Config;
var Listen *listenconfig
var App *appconfig

// Singleton
var Cache cache.Cache
var Session *session.Manager;
var Router *router.Mux
var Sql *sql.DB


func init(){
	Config = config.Instance()
	App = &appconfig{}
	Listen = &listenconfig{}
	Config.App = App
	Config.Listen = Listen
}

func Reload() (err error) {
	// cache
	//Cache,_ = cache.NewCache(Config.App.Cache)
	// session
	sessionConfig := &session.ManagerConfig{}
	json.Unmarshal([]byte(App.Session),sessionConfig)
	Session, err = session.NewManager("memcache", sessionConfig)
	log.Info("Session: ",err)
	// router
	Router = router.Instance()
	// sql
	Sql,err = sql.Open("mysql",App.Mysql)
	log.Info("Sql: ",err)
	return
}



type appconfig struct {
	Mysql 		string		`comment:"Mysql"`
	Memcache 	string		`comment:"Memcached"`
	Session		string		`comment:"Session"`
	Cache		string		`comment:"Cache"`
}


type listenconfig struct {
	Ip			string		`comment:"Listen Ip Addr" json:"IP"`
	Port		int			`comment:"Server use port"`
	Https		bool		`comment:"is https"`
	Html2		bool		`comment:"is html2"`
	Certfile	string		`comment:"cert file"`
	Keyfile		string		`comment:"key file"`
}

func (l *listenconfig) Addr() string {
 	return fmt.Sprintf("%s:%d",l.Ip,l.Port)
 } 