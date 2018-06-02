package config;
 
import (
	"fmt"
	"os"
	"strings"
	"encoding/json"
)
 
type Config struct {
	Config		string      `comment:"config path"`
	Command		string      `comment:"start command"`
	Workdir		string      `comment:"Current working directory"`
	Tempdir		string      `comment:"Template file dir"`
	Pidfile		string      `comment:"Pid file path"`
	Logfile		string      `comment:"Log file path"`
	Ip			string      `comment:"Listen Ip Addr" json:"IP"`
	Port		int         `comment:"Server use port"`
	Listen		*Listen		`comment:"Listen Info"`
	App 		*App 		`comment:"app config"`
	Dbconfig	string      `comment:"MariaDB connect info"`
	Memcache	string      `comment:"Memcached connect addr and port"`
	Session		string
	Cache		string
	Const		map[string]*string
	Enable		[]string
	Disable		[]string
	Flag		map[string]interface{}	`json:"-"`
	Mode		map[string]interface{}
}

type Listen struct {
	Ip			string		`comment:"Listen Ip Addr" json:"IP"`
	Port		int			`comment:"Server use port"`
	Https		bool		`comment:"is https"`
	Html2		bool		`comment:"is html2"`
	Certfile	string		`comment:"cert file"`
	Keyfile		string		`comment:"key file"`
}

type App struct {
	Mysql 		string		`comment:"Mysql"`
	Memcache 	string		`comment:"Memcached"`
	Session		string		`comment:"Session"`
	Cache		string		`comment:"Cache"`
}

func (c *Config) setconf() {
	// set flag workdir 
	for _,value := range os.Args[1:] {
		if strings.HasPrefix(value, "--workdir=") || strings.HasPrefix(value, "--config=") {
			c.set(value[2:])
		}
	}
	os.Chdir(conf.Workdir)
	// get config value
	if cf, err := readconfig(c.Config);err == nil {
		err := json.Unmarshal(cf,c)
		if err != nil {
			fmt.Println("配置解析失败:", err)
		}
	}else {
		fmt.Println("配置读取失败:", err)
	}
}

func (c *Config) setmode() {	
	// set flag mode
	for _,value := range os.Args[1:] {
		if strings.HasPrefix(value, "--enable=") || strings.HasPrefix(value, "--disable=") {
			c.set(value[2:])
		}
	}

	var d []string
	for _,ve := range conf.Enable {
		b := true
		for _,vd := range c.Disable {
			if ve==vd {
				b = false
				break
			}
		}
		if b{
			d=append(d,ve)
		}
	}
	conf.Enable = d
	
	// set mode config
	for _,v := range conf.Enable {
		if b,err := json.Marshal(conf.Mode[v]);err == nil && conf.Mode[v] != nil {
			json.Unmarshal(b, &conf)
		}
	}
	os.Chdir(conf.Workdir) 
}

func (c *Config) Reload() {
	c.setconf()
	c.setmode()
	c.readenv()
	c.readflag()
}

func NewConfig() *Config {
	return &Config {
		Workdir:	"/data/web",
		Tempdir:	"template",
		Config:		"config/conf.json",
		Pidfile:	"/var/run/index.pid",
		Logfile:	"logs",
		Port:		80,
		Dbconfig:	"root:@/Jass",
	}
}

var conf *Config

func init() {
	conf = NewConfig()
	conf.Reload()
}



func (c *Config) IsMode(m string) bool {
	for _,v := range c.Enable {
		if v==m {
			return true
		}
	}
	return false
}


func (c *Config) Getconst(k string) *string {
	if v,ok := c.Const[k]; ok {
		return v
	}
	return nil
}


func Instance() *Config {
	return conf
}

func Reload() {
	conf.Reload()
}

func IsMode(m string) bool {
	return conf.IsMode(m)
}
	
func Getconst(k string) *string {
	return conf.Getconst(k)
}
