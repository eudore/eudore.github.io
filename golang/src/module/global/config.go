package global

import (
	"fmt"
)

type config struct {
	Config		string		`comment:"config path"`
	Command		string		`comment:"start command"`
	Workdir		string		`comment:"Current working directory"`
	Tempdir		string		`comment:"Template file dir"`
	Pidfile		string		`comment:"Pid file path"`
	Logfile		string		`comment:"Log file path"`
	Listen		*listenconfig	`comment:"Listen Info"`
	App 		*appconfig		`comment:"App config"`
}

type appconfig struct {
	Mysql 		string		`comment:"Mysql"`
	Etcd		string		`comment:"Etcd Addr"`
	Memcache 	string		`comment:"Memcached Addr"`
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