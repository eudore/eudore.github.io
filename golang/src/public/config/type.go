package config

import (
	"fmt"
)


type Listen struct {
	Ip			string		`comment:"Listen Ip Addr" json:"IP"`
	Port		int			`comment:"Server use port"`
	Https		bool		`comment:"is https"`
	Html2		bool		`comment:"is html2"`
	Certfile	string		`comment:"cert file"`
	Keyfile		string		`comment:"key file"`
}

func (l *Listen) Addr() string {
 	return fmt.Sprintf("%s:%d",l.Ip,l.Port)
 } 

type App struct {
	Mysql 		string		`comment:"Mysql"`
	Memcache 	string		`comment:"Memcached"`
	Session		string		`comment:"Session"`
	Cache		string		`comment:"Cache"`
}

type Auth struct {
	RedirectURL string
}